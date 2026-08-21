package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/anticheat"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
	"golang.org/x/crypto/bcrypt"
)

// AdminHandler serves the operator console under /api/admin.
//
// Everything here is read-mostly aggregation over tables other handlers own,
// plus the few writes an operator needs: promote an account, retune an
// anti-cheat rule, remove a player. It holds the live services so the console
// can show what is in memory, not only what reached SQLite.
type AdminHandler struct {
	ranked    *services.RankedMatchService
	queue     *services.MatchmakingQueue
	rooms     *services.RoomService
	rules     *anticheat.Store
	findings  *anticheat.DBSink
	generator *services.QuestionGenerator
	dbPath    string
	startedAt time.Time
}

// NewAdminHandler wires the console to the live services. findings may be nil
// when no database sink is installed; the console then reports no drops.
func NewAdminHandler(
	ranked *services.RankedMatchService,
	queue *services.MatchmakingQueue,
	rooms *services.RoomService,
	rules *anticheat.Store,
	findings *anticheat.DBSink,
	generator *services.QuestionGenerator,
	dbPath string,
) *AdminHandler {
	return &AdminHandler{
		ranked:    ranked,
		queue:     queue,
		rooms:     rooms,
		rules:     rules,
		findings:  findings,
		generator: generator,
		dbPath:    dbPath,
		startedAt: time.Now(),
	}
}

// RequireAdmin rejects anyone whose account is not flagged as an admin.
//
// It runs after middleware.AuthRequired() and reads users.is_admin on every
// request rather than trusting a claim in the token: demoting an account must
// take effect immediately, and tokens live for a week. That is also why the
// flag is not in the JWT — middleware stays database-free, so the check lives
// here with the rest of the admin surface.
func (h *AdminHandler) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var isAdmin bool
		err := database.DB.QueryRowContext(c.Request.Context(),
			"SELECT is_admin FROM users WHERE id = ?", c.GetInt64("user_id"),
		).Scan(&isAdmin)
		if err != nil || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// PromoteAdmins flags the named accounts as admins at startup, so a fresh
// deployment has a way in that does not involve hand-editing SQLite. A name
// with no account yet is skipped — register it, restart, and it takes.
func PromoteAdmins(ctx context.Context, usernames []string) {
	for _, name := range usernames {
		res, err := database.DB.ExecContext(ctx,
			"UPDATE users SET is_admin = 1 WHERE LOWER(username) = LOWER(?)", name)
		if err != nil {
			slog.Error("admin: promoting account failed", "username", name, "error", err)
			continue
		}
		if n, _ := res.RowsAffected(); n == 0 {
			slog.Warn("admin: no such account to promote", "username", name)
			continue
		}
		slog.Info("admin: account promoted", "username", name)
	}
}

// ─── Overview ──────────────────────────────────────────────────

// GetOverview handles GET /api/admin/overview.
func (h *AdminHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()
	var o models.AdminOverview

	database.DB.QueryRowContext(ctx, `
		SELECT COUNT(*),
		       COALESCE(SUM(is_admin), 0),
		       COALESCE(SUM(created_at >= datetime('now', '-7 days')), 0),
		       COALESCE(SUM(total_points), 0)
		FROM users
	`).Scan(&o.Users, &o.Admins, &o.NewUsers7d, &o.TotalPoints)

	var answered, correct int
	database.DB.QueryRowContext(ctx, `
		SELECT COUNT(*),
		       COALESCE(SUM(played_at >= datetime('now', '-1 day')), 0),
		       COUNT(DISTINCT CASE WHEN played_at >= datetime('now', '-7 days') THEN user_id END),
		       COALESCE(SUM(total_answered), 0),
		       COALESCE(SUM(correct_answers), 0)
		FROM scores
	`).Scan(&o.Games, &o.Games24h, &o.ActivePlayers7d, &answered, &correct)
	if answered > 0 {
		o.AverageAccuracy = float64(correct) / float64(answered) * 100
	}

	var attempts, rankedAnswers int
	database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM matches").Scan(&o.SoloMatches)
	database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM ranked_matches").Scan(&o.RankedMatches)
	database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM question_attempts").Scan(&attempts)
	database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM ranked_question_results").Scan(&rankedAnswers)
	o.QuestionsAnswered = answered + attempts + rankedAnswers

	database.DB.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(at >= datetime('now', '-1 day')), 0) FROM anticheat_findings
	`).Scan(&o.Findings, &o.Findings24h)

	for _, r := range h.rules.Snapshot() {
		if r.Enabled && r.Action == anticheat.ActionReject {
			o.RejectingRules++
		}
	}
	if h.findings != nil {
		o.FindingsDropped = h.findings.Dropped()
	}

	o.QueueSize = h.queue.QueueSize()
	o.LiveMatches = len(h.ranked.Snapshot())
	o.OpenRooms = len(h.rooms.Snapshot())
	o.LiveSoloMatches = services.GlobalMatchStore.Len()
	o.StoredQuestions = services.GlobalQuestionStore.Len()
	o.Topics = len(services.TopicCatalog())
	o.UptimeSeconds = time.Since(h.startedAt).Seconds()
	o.Goroutines = runtime.NumGoroutine()
	if info, err := os.Stat(h.dbPath); err == nil {
		o.DatabaseBytes = info.Size()
	}

	c.JSON(http.StatusOK, o)
}

// ─── Analytics ─────────────────────────────────────────────────

// GetAnalytics handles GET /api/admin/analytics.
func (h *AdminHandler) GetAnalytics(c *gin.Context) {
	ctx := c.Request.Context()
	days := clampInt(intQuery(c, "days", 30), 1, 90)

	// Four tables keyed by day, merged into one series so the chart has a row
	// for every day in the window — including the quiet ones.
	dates := make([]string, 0, days)
	byDate := make(map[string]*models.AdminDailyPoint, days)
	for i := days - 1; i >= 0; i-- {
		date := time.Now().UTC().AddDate(0, 0, -i).Format(time.DateOnly)
		dates = append(dates, date)
		byDate[date] = &models.AdminDailyPoint{Date: date}
	}

	window := fmt.Sprintf("-%d days", days)
	daily := func(query string, apply func(p *models.AdminDailyPoint, a, b int)) {
		err := eachRow(ctx, query, []any{window}, func(rows *sql.Rows) {
			var date string
			var a, b int
			if err := rows.Scan(&date, &a, &b); err != nil {
				return
			}
			if p, ok := byDate[date]; ok {
				apply(p, a, b)
			}
		})
		if err != nil {
			slog.Error("admin: daily analytics query failed", "error", err)
		}
	}

	daily(`SELECT date(played_at), COUNT(*), COALESCE(SUM(score), 0)
	       FROM scores WHERE played_at >= datetime('now', ?) GROUP BY 1`,
		func(p *models.AdminDailyPoint, games, points int) { p.Games, p.Points = games, points })
	daily(`SELECT date(played_at), COUNT(DISTINCT user_id), 0
	       FROM scores WHERE played_at >= datetime('now', ?) GROUP BY 1`,
		func(p *models.AdminDailyPoint, players, _ int) { p.Players = players })
	daily(`SELECT date(created_at), COUNT(*), 0
	       FROM ranked_matches WHERE created_at >= datetime('now', ?) GROUP BY 1`,
		func(p *models.AdminDailyPoint, matches, _ int) { p.RankedMatches = matches })
	daily(`SELECT date(at), COUNT(*), 0
	       FROM anticheat_findings WHERE at >= datetime('now', ?) GROUP BY 1`,
		func(p *models.AdminDailyPoint, findings, _ int) { p.Findings = findings })

	out := models.AdminAnalytics{
		Daily:  make([]models.AdminDailyPoint, 0, days),
		Topics: []models.AdminTopicStat{},
		Hours:  make([]models.AdminBucket, 24),
	}
	for _, date := range dates {
		out.Daily = append(out.Daily, *byDate[date])
	}

	// Games and mean accuracy per difficulty.
	out.Difficulty = queryBuckets(ctx, `
		SELECT COALESCE(difficulty, 'easy'), COUNT(*),
		       COALESCE(AVG(CASE WHEN total_answered > 0
		                         THEN correct_answers * 100.0 / total_answered END), 0)
		FROM scores GROUP BY 1 ORDER BY 2 DESC`)

	// Topic performance comes from ranked results, the only table that records
	// which topic each question came from.
	eachRow(ctx, `
		SELECT topic, difficulty, COUNT(*), COALESCE(SUM(correct), 0), COALESCE(AVG(time_spent), 0)
		FROM ranked_question_results WHERE topic <> '' GROUP BY topic, difficulty ORDER BY 3 DESC`,
		nil, func(rows *sql.Rows) {
			var t models.AdminTopicStat
			if err := rows.Scan(&t.Topic, &t.Difficulty, &t.Attempts, &t.Correct, &t.AverageTime); err != nil {
				return
			}
			if t.Attempts > 0 {
				t.Accuracy = float64(t.Correct) / float64(t.Attempts) * 100
			}
			out.Topics = append(out.Topics, t)
		})

	for i := range out.Hours {
		out.Hours[i] = models.AdminBucket{Label: fmt.Sprintf("%02d", i)}
	}
	for _, b := range queryBuckets(ctx, `SELECT strftime('%H', played_at), COUNT(*), 0 FROM scores GROUP BY 1`) {
		if hour, err := strconv.Atoi(b.Label); err == nil && hour >= 0 && hour < 24 {
			out.Hours[hour].Count = b.Count
		}
	}

	out.Ratings = queryBuckets(ctx, `
		SELECT CAST((rating / 100) * 100 AS TEXT), COUNT(*), 0
		FROM users WHERE ranked_wins + ranked_losses > 0 GROUP BY 1 ORDER BY 1`)
	out.Rules = queryBuckets(ctx, `SELECT rule, COUNT(*), 0 FROM anticheat_findings GROUP BY 1 ORDER BY 2 DESC`)

	c.JSON(http.StatusOK, out)
}

// ─── Leaderboards ──────────────────────────────────────────────

// GetLeaderboards handles GET /api/admin/leaderboards. It reads both boards
// live: the public one sits behind a 30-second cache, which is the wrong
// answer for an operator checking whether a correction landed.
func (h *AdminHandler) GetLeaderboards(c *gin.Context) {
	ctx := c.Request.Context()
	limit := clampInt(intQuery(c, "limit", 100), 1, 500)

	out := models.AdminLeaderboards{
		Casual: []models.LeaderboardEntry{},
		Ranked: []models.RankedLeaderboardEntry{},
	}

	err := eachRow(ctx, `
		SELECT u.username, u.id, u.total_points, COUNT(s.id)
		FROM users u LEFT JOIN scores s ON s.user_id = u.id
		WHERE u.total_points > 0
		GROUP BY u.id ORDER BY u.total_points DESC LIMIT ?`,
		[]any{limit}, func(rows *sql.Rows) {
			var e models.LeaderboardEntry
			if err := rows.Scan(&e.Username, &e.UserID, &e.TotalPoints, &e.TotalGames); err != nil {
				return
			}
			e.Rank = len(out.Casual) + 1
			out.Casual = append(out.Casual, e)
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboards"})
		return
	}

	ranked, err := rankedLadder(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ranked ladder"})
		return
	}
	out.Ranked = ranked

	c.JSON(http.StatusOK, out)
}

// ─── Users ─────────────────────────────────────────────────────

// userSortColumns is the whitelist behind the ?sort= parameter. Its values are
// spliced into SQL, so a sort key may only ever come from this map.
var userSortColumns = map[string]string{
	"points":    "u.total_points",
	"rating":    "u.rating",
	"games":     "games",
	"accuracy":  "accuracy",
	"findings":  "findings",
	"created":   "u.created_at",
	"username":  "u.username",
	"lastplay":  "last_played",
	"wins":      "u.ranked_wins",
	"highest":   "u.highest_rating",
	"answered":  "answered",
	"highscore": "high_score",
}

// GetUsers handles GET /api/admin/users.
func (h *AdminHandler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	page := clampInt(intQuery(c, "page", 1), 1, 100000)
	pageSize := clampInt(intQuery(c, "pageSize", 25), 1, 200)
	search := strings.ToLower(strings.TrimSpace(c.Query("search")))

	sortColumn, ok := userSortColumns[strings.ToLower(c.DefaultQuery("sort", "points"))]
	if !ok {
		sortColumn = "u.total_points"
	}
	direction := "DESC"
	if strings.EqualFold(c.Query("order"), "asc") {
		direction = "ASC"
	}

	like := "%" + search + "%"
	var total int
	if err := database.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM users u WHERE ? = '' OR LOWER(u.username) LIKE ?", search, like,
	).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count users"})
		return
	}

	online := h.ranked.ActiveUserIDs()
	for _, q := range h.queue.Snapshot() {
		online[q.UserID] = true
	}

	out := models.AdminUsersResponse{Users: []models.AdminUserRow{}, Total: total, Page: page, PageSize: pageSize}
	err := eachRow(ctx, `
		SELECT u.id, u.username, u.is_admin, u.total_points, u.rating, u.ranked_wins,
		       u.ranked_losses, u.highest_rating, u.created_at,
		       COUNT(s.id) AS games,
		       COALESCE(MAX(s.score), 0) AS high_score,
		       COALESCE(SUM(s.total_answered), 0) AS answered,
		       COALESCE(SUM(s.correct_answers), 0) AS correct,
		       COALESCE(SUM(s.correct_answers), 0) * 100.0 / NULLIF(SUM(s.total_answered), 0) AS accuracy,
		       MAX(s.played_at) AS last_played,
		       (SELECT COUNT(*) FROM anticheat_findings f WHERE f.user_id = u.id) AS findings
		FROM users u LEFT JOIN scores s ON s.user_id = u.id
		WHERE ? = '' OR LOWER(u.username) LIKE ?
		GROUP BY u.id
		ORDER BY `+sortColumn+" "+direction+`, u.id ASC
		LIMIT ? OFFSET ?`,
		[]any{search, like, pageSize, (page - 1) * pageSize},
		func(rows *sql.Rows) {
			var u models.AdminUserRow
			var accuracy sql.NullFloat64
			var lastPlayed sql.NullString
			err := rows.Scan(
				&u.ID, &u.Username, &u.IsAdmin, &u.TotalPoints, &u.Rating, &u.RankedWins,
				&u.RankedLosses, &u.HighestRating, &u.CreatedAt,
				&u.Games, &u.HighScore, &u.TotalAnswered, &u.TotalCorrect, &accuracy, &lastPlayed, &u.Findings,
			)
			if err != nil {
				slog.Error("admin: reading user row failed", "error", err)
				return
			}
			u.Accuracy = accuracy.Float64
			u.LastPlayed = parseSQLiteTime(lastPlayed)
			u.Online = online[u.ID]
			out.Users = append(out.Users, u)
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, out)
}

// UpdateUser handles POST /api/admin/users/update. Omitted fields are left as
// they are, so the console can send one toggle at a time.
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	var req models.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update"})
		return
	}

	// Demoting yourself is how an operator locks everyone out of the console:
	// with the last admin gone, access comes back only by editing SQLite or
	// restarting with ADMIN_USERNAMES set.
	if req.IsAdmin != nil && !*req.IsAdmin && req.UserID == c.GetInt64("user_id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove your own admin access"})
		return
	}

	sets := []string{}
	args := []any{}
	if req.IsAdmin != nil {
		sets = append(sets, "is_admin = ?")
		args = append(args, boolToInt(*req.IsAdmin))
	}
	if req.Rating != nil {
		rating := clampInt(*req.Rating, 0, 100000)
		sets = append(sets, "rating = ?", "highest_rating = MAX(highest_rating, ?)")
		args = append(args, rating, rating)
	}
	if req.TotalPoints != nil {
		sets = append(sets, "total_points = ?")
		args = append(args, clampInt(*req.TotalPoints, 0, 1000000000))
	}
	if req.Password != nil {
		if len(*req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		sets = append(sets, "password_hash = ?")
		args = append(args, string(hash))
	}
	if len(sets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nothing to update"})
		return
	}

	args = append(args, req.UserID)
	res, err := database.DB.ExecContext(c.Request.Context(),
		"UPDATE users SET "+strings.Join(sets, ", ")+" WHERE id = ?", args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	slog.Info("admin: user updated", "by", c.GetString("username"), "user_id", req.UserID, "fields", len(sets))
	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// DeleteUser handles POST /api/admin/users/delete. The cascade itself lives in
// database.DeleteUser, shared with cmd/fastchemctl so the two cannot drift.
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	var req models.AdminUserIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.UserID == c.GetInt64("user_id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	if err := database.DeleteUser(c.Request.Context(), req.UserID); err != nil {
		slog.Error("admin: deleting user failed", "user_id", req.UserID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	slog.Warn("admin: user deleted", "by", c.GetString("username"), "user_id", req.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ─── Anti-cheat ────────────────────────────────────────────────

// GetRules handles GET /api/admin/anticheat/rules. Settings come from the live
// store rather than the table, so the console shows what is actually running.
func (h *AdminHandler) GetRules(c *gin.Context) {
	ctx := c.Request.Context()

	total := map[string]int{}
	recent := map[string]int{}
	for _, b := range queryBuckets(ctx, `
		SELECT rule, COUNT(*), COALESCE(SUM(at >= datetime('now', '-1 day')), 0)
		FROM anticheat_findings GROUP BY 1`) {
		total[b.Label] = b.Count
		recent[b.Label] = int(b.Value)
	}

	out := []models.AdminRule{}
	for _, r := range h.rules.Snapshot() {
		out = append(out, models.AdminRule{
			Name:        r.Name,
			Enabled:     r.Enabled,
			Action:      string(r.Action),
			Params:      r.Params,
			Findings:    total[r.Name],
			Findings24h: recent[r.Name],
		})
	}
	c.JSON(http.StatusOK, gin.H{"rules": out})
}

// UpdateRule handles POST /api/admin/anticheat/rules. The write goes through
// anticheat.Store.Save, which validates and refreshes the cache immediately
// rather than waiting on the 30-second reload ticker.
func (h *AdminHandler) UpdateRule(c *gin.Context) {
	var req models.AdminRuleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule update"})
		return
	}

	action, ok := anticheat.ParseAction(req.Action)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
		return
	}
	if err := h.rules.Save(c.Request.Context(), req.Name, req.Enabled, action, anticheat.Params(req.Params)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Warn("admin: anticheat rule changed",
		"by", c.GetString("username"), "rule", req.Name, "enabled", req.Enabled, "action", req.Action)
	c.JSON(http.StatusOK, gin.H{"message": "Rule updated"})
}

// GetFindings handles GET /api/admin/anticheat/findings.
func (h *AdminHandler) GetFindings(c *gin.Context) {
	ctx := c.Request.Context()
	page := clampInt(intQuery(c, "page", 1), 1, 100000)
	pageSize := clampInt(intQuery(c, "pageSize", 50), 1, 200)

	// One filter expression, reused by the page query and every breakdown, so
	// the summary always describes exactly the rows being listed.
	where := `WHERE (? = '' OR f.rule = ?)
	          AND (? = '' OR f.mode = ?)
	          AND (? = '' OR f.action = ?)
	          AND (? = 0 OR f.user_id = ?)`
	rule, mode, action := c.Query("rule"), c.Query("mode"), c.Query("action")
	userID := int64(intQuery(c, "userId", 0))
	args := []any{rule, rule, mode, mode, action, action, userID, userID}

	out := models.AdminFindingsResponse{Findings: []models.AdminFinding{}, Page: page, PageSize: pageSize}
	if h.findings != nil {
		out.Dropped = h.findings.Dropped()
	}

	if err := database.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM anticheat_findings f "+where, args...,
	).Scan(&out.Total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count findings"})
		return
	}

	err := eachRow(ctx, `
		SELECT f.id, f.at, f.subject, f.user_id, COALESCE(u.username, ''), f.match_id, f.mode,
		       f.question_id, f.difficulty, f.time_spent, f.rule, f.detail, f.action
		FROM anticheat_findings f LEFT JOIN users u ON u.id = f.user_id
		`+where+` ORDER BY f.at DESC, f.id DESC LIMIT ? OFFSET ?`,
		append(append([]any{}, args...), pageSize, (page-1)*pageSize),
		func(rows *sql.Rows) {
			var f models.AdminFinding
			err := rows.Scan(&f.ID, &f.At, &f.Subject, &f.UserID, &f.Username, &f.MatchID, &f.Mode,
				&f.QuestionID, &f.Difficulty, &f.TimeSpent, &f.Rule, &f.Detail, &f.Action)
			if err != nil {
				slog.Error("admin: reading finding failed", "error", err)
				return
			}
			out.Findings = append(out.Findings, f)
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch findings"})
		return
	}

	out.ByRule = queryBuckets(ctx, "SELECT f.rule, COUNT(*), 0 FROM anticheat_findings f "+where+" GROUP BY 1 ORDER BY 2 DESC", args...)
	out.ByMode = queryBuckets(ctx, "SELECT f.mode, COUNT(*), 0 FROM anticheat_findings f "+where+" GROUP BY 1 ORDER BY 2 DESC", args...)
	out.ByAction = queryBuckets(ctx, "SELECT f.action, COUNT(*), 0 FROM anticheat_findings f "+where+" GROUP BY 1 ORDER BY 2 DESC", args...)

	c.JSON(http.StatusOK, out)
}

// ─── Live state ────────────────────────────────────────────────

// GetLive handles GET /api/admin/live: everything the server holds in memory
// right now, which no table can show.
func (h *AdminHandler) GetLive(c *gin.Context) {
	c.JSON(http.StatusOK, models.AdminLive{
		Queue:           h.queue.Snapshot(),
		Matches:         h.ranked.Snapshot(),
		Rooms:           h.rooms.Snapshot(),
		SoloMatches:     services.GlobalMatchStore.Len(),
		StoredQuestions: services.GlobalQuestionStore.Len(),
	})
}

// ─── Recent activity ───────────────────────────────────────────

// GetMatches handles GET /api/admin/matches.
func (h *AdminHandler) GetMatches(c *gin.Context) {
	ctx := c.Request.Context()
	limit := clampInt(intQuery(c, "limit", 25), 1, 200)
	out := models.AdminMatchesResponse{
		Ranked: []models.AdminRankedMatchRow{},
		Solo:   []models.AdminSoloMatchRow{},
		Scores: []models.AdminScoreRow{},
	}

	eachRow(ctx, `
		SELECT m.id, COALESCE(p1.username, '?'), COALESCE(p2.username, '?'),
		       m.player1_score, m.player2_score, COALESCE(w.username, ''), m.status, m.created_at, m.finished_at
		FROM ranked_matches m
		LEFT JOIN users p1 ON p1.id = m.player1_id
		LEFT JOIN users p2 ON p2.id = m.player2_id
		LEFT JOIN users w ON w.id = m.winner_id
		ORDER BY m.created_at DESC LIMIT ?`,
		[]any{limit}, func(rows *sql.Rows) {
			var r models.AdminRankedMatchRow
			var finished sql.NullString
			if err := rows.Scan(&r.MatchID, &r.Player1, &r.Player2, &r.Player1Score, &r.Player2Score,
				&r.Winner, &r.Status, &r.CreatedAt, &finished); err != nil {
				return
			}
			r.FinishedAt = parseSQLiteTime(finished)
			out.Ranked = append(out.Ranked, r)
		})

	eachRow(ctx, `
		SELECT m.id, COALESCE(u.username, '?'), m.difficulty, m.total_score, m.best_combo,
		       (SELECT COUNT(*) FROM question_attempts a WHERE a.match_id = m.id), m.started_at, m.ended_at
		FROM matches m LEFT JOIN users u ON u.id = m.user_id
		ORDER BY m.started_at DESC LIMIT ?`,
		[]any{limit}, func(rows *sql.Rows) {
			var r models.AdminSoloMatchRow
			var ended sql.NullString
			if err := rows.Scan(&r.MatchID, &r.Username, &r.Difficulty, &r.TotalScore, &r.BestCombo,
				&r.Attempts, &r.StartedAt, &ended); err != nil {
				return
			}
			r.EndedAt = parseSQLiteTime(ended)
			out.Solo = append(out.Solo, r)
		})

	eachRow(ctx, `
		SELECT s.id, s.user_id, COALESCE(u.username, '?'), s.score, s.total_answered, s.correct_answers,
		       COALESCE(s.difficulty, 'easy'), COALESCE(s.time_spent, 0), s.played_at
		FROM scores s LEFT JOIN users u ON u.id = s.user_id
		ORDER BY s.played_at DESC LIMIT ?`,
		[]any{limit}, func(rows *sql.Rows) {
			var r models.AdminScoreRow
			if err := rows.Scan(&r.ID, &r.UserID, &r.Username, &r.Score, &r.TotalAnswered,
				&r.CorrectAnswers, &r.Difficulty, &r.TimeSpent, &r.PlayedAt); err != nil {
				return
			}
			out.Scores = append(out.Scores, r)
		})

	c.JSON(http.StatusOK, out)
}

// ─── Question bank ─────────────────────────────────────────────

// GetTopics handles GET /api/admin/topics.
func (h *AdminHandler) GetTopics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"topics": services.TopicCatalog()})
}

// PreviewQuestion handles GET /api/admin/topics/preview?category=…. The
// generated question is a throwaway: never stored, never scored, never served
// to a player, so returning its answer reveals nothing about a live question.
func (h *AdminHandler) PreviewQuestion(c *gin.Context) {
	preview, ok := h.generator.PreviewQuestion(c.Query("category"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown category"})
		return
	}
	c.JSON(http.StatusOK, preview)
}

// ─── Helpers ───────────────────────────────────────────────────

// eachRow runs a query, hands every row to fn, and closes the result set
// before returning.
//
// SQLite is capped at two pooled connections, so a handler holding three
// result sets open at once deadlocks waiting for a third. Every read in this
// file goes through here so no result set outlives the loop that reads it.
func eachRow(ctx context.Context, query string, args []any, fn func(*sql.Rows)) error {
	rows, err := database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		fn(rows)
	}
	return rows.Err()
}

// queryBuckets runs a three-column (label, count, value) query. A failure
// yields an empty series: one broken chart must not take the page down.
func queryBuckets(ctx context.Context, query string, args ...any) []models.AdminBucket {
	out := []models.AdminBucket{}
	err := eachRow(ctx, query, args, func(rows *sql.Rows) {
		var b models.AdminBucket
		if err := rows.Scan(&b.Label, &b.Count, &b.Value); err != nil {
			return
		}
		out = append(out, b)
	})
	if err != nil {
		slog.Error("admin: bucket query failed", "error", err)
	}
	return out
}

func intQuery(c *gin.Context, key string, def int) int {
	if v, err := strconv.Atoi(c.Query(key)); err == nil {
		return v
	}
	return def
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// parseSQLiteTime converts a nullable timestamp read as text. Aggregates like
// MAX(played_at) lose the column's declared type, so the driver hands them
// back as strings rather than time.Time.
func parseSQLiteTime(s sql.NullString) *time.Time {
	if !s.Valid || s.String == "" {
		return nil
	}
	for _, layout := range []string{"2006-01-02 15:04:05.999999999-07:00", time.DateTime, time.RFC3339} {
		if t, err := time.Parse(layout, s.String); err == nil {
			return &t
		}
	}
	return nil
}
