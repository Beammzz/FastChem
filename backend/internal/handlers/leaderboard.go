package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// leaderboardCache holds a cached leaderboard response with TTL.
type leaderboardCache struct {
	mu      sync.RWMutex
	entries []models.LeaderboardEntry
	updated time.Time
	ttl     time.Duration
}

func newLeaderboardCache(ttl time.Duration) *leaderboardCache {
	return &leaderboardCache{ttl: ttl}
}

func (c *leaderboardCache) get() ([]models.LeaderboardEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entries == nil || time.Since(c.updated) > c.ttl {
		return nil, false
	}
	return c.entries, true
}

func (c *leaderboardCache) set(entries []models.LeaderboardEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = entries
	c.updated = time.Now()
}

// LeaderboardHandler handles leaderboard-related HTTP requests.
type LeaderboardHandler struct {
	cache *leaderboardCache
}

// NewLeaderboardHandler creates a new leaderboard handler.
func NewLeaderboardHandler() *LeaderboardHandler {
	return &LeaderboardHandler{
		cache: newLeaderboardCache(30 * time.Second),
	}
}

// SubmitScore handles POST /api/scores
func (h *LeaderboardHandler) SubmitScore(c *gin.Context) {
	userID := c.GetInt64("user_id")
	ctx := c.Request.Context()

	var req models.SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score data"})
		return
	}

	// Sanity checks
	if req.CorrectAnswers > req.TotalAnswered {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Correct answers cannot exceed total answered"})
		return
	}

	// Validate score against the scoring formula (with max combo multiplier)
	maxPerQuestion := services.MaxScorePerQuestion(req.Difficulty)
	maxPossibleScore := req.CorrectAnswers * maxPerQuestion
	if req.Score > maxPossibleScore {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score value"})
		return
	}

	tx, err := database.DB.BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
		return
	}
	defer tx.Rollback()

	// Insert score record
	result, err := tx.ExecContext(ctx,
		"INSERT INTO scores (user_id, score, total_answered, correct_answers, difficulty, time_limit, time_spent) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, req.Score, req.TotalAnswered, req.CorrectAnswers, req.Difficulty, req.TimeLimit, req.TimeSpent,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
		return
	}

	scoreID, _ := result.LastInsertId()

	// Add score to user's total_points
	_, err = tx.ExecContext(ctx, "UPDATE users SET total_points = total_points + ? WHERE id = ?", req.Score, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
		return
	}

	// Get updated total points
	var totalPoints int
	database.DB.QueryRowContext(ctx, "SELECT total_points FROM users WHERE id = ?", userID).Scan(&totalPoints)

	c.JSON(http.StatusCreated, gin.H{
		"id":          scoreID,
		"message":     "Score saved successfully",
		"totalPoints": totalPoints,
	})
}

// GetLeaderboard handles GET /api/leaderboard
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	// Check cache first
	if cached, ok := h.cache.get(); ok {
		c.JSON(http.StatusOK, models.LeaderboardResponse{Entries: cached})
		return
	}

	query := `
		SELECT 
			u.username,
			u.id as user_id,
			u.total_points,
			COUNT(s.id) as total_games
		FROM users u
		LEFT JOIN scores s ON s.user_id = u.id
		WHERE u.total_points > 0
		GROUP BY u.id
		ORDER BY u.total_points DESC
		LIMIT 50
	`

	rows, err := database.DB.QueryContext(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	rank := 0
	for rows.Next() {
		rank++
		var entry models.LeaderboardEntry
		err := rows.Scan(
			&entry.Username,
			&entry.UserID,
			&entry.TotalPoints,
			&entry.TotalGames,
		)
		if err != nil {
			continue
		}
		entry.Rank = rank
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read leaderboard"})
		return
	}

	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	// Store in cache
	h.cache.set(entries)

	c.JSON(http.StatusOK, models.LeaderboardResponse{
		Entries: entries,
	})
}

// GetUserScores handles GET /api/scores/me
func (h *LeaderboardHandler) GetUserScores(c *gin.Context) {
	userID := c.GetInt64("user_id")
	ctx := c.Request.Context()

	rows, err := database.DB.QueryContext(ctx, `
		SELECT id, score, total_answered, correct_answers, COALESCE(difficulty, 'easy'), time_limit, COALESCE(time_spent, 0), played_at
		FROM scores
		WHERE user_id = ?
		ORDER BY played_at DESC
		LIMIT 20
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores"})
		return
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var s models.Score
		s.UserID = userID
		err := rows.Scan(&s.ID, &s.Score, &s.TotalAnswered, &s.CorrectAnswers, &s.Difficulty, &s.TimeLimit, &s.TimeSpent, &s.PlayedAt)
		if err != nil {
			continue
		}
		scores = append(scores, s)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read scores"})
		return
	}

	if scores == nil {
		scores = []models.Score{}
	}

	c.JSON(http.StatusOK, gin.H{"scores": scores})
}

// GetProfile handles GET /api/profile/:username
func (h *LeaderboardHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit
	ctx := c.Request.Context()

	// Get user
	var userID int64
	var totalPoints int
	err := database.DB.QueryRowContext(ctx, "SELECT id, total_points FROM users WHERE LOWER(username) = LOWER(?)", username).Scan(&userID, &totalPoints)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get stats
	var stats models.UserStats
	stats.TotalPoints = totalPoints
	err = database.DB.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total_games,
			COALESCE(MAX(score), 0) as high_score,
			COALESCE(SUM(correct_answers), 0) as total_correct,
			COALESCE(SUM(total_answered), 0) as total_answered,
			COALESCE(AVG(score), 0) as average_score,
			COALESCE(MIN(NULLIF(time_spent, 0)), 0) as fastest_time,
			COALESCE(AVG(NULLIF(time_spent, 0)), 0) as average_time
		FROM scores WHERE user_id = ?
	`, userID).Scan(&stats.TotalGames, &stats.HighScore, &stats.TotalCorrect, &stats.TotalAnswered, &stats.AverageScore, &stats.FastestTime, &stats.AverageTime)

	if err != nil {
		stats = models.UserStats{TotalPoints: totalPoints}
	}

	if stats.TotalAnswered > 0 {
		stats.AverageAccuracy = float64(stats.TotalCorrect) / float64(stats.TotalAnswered) * 100
	}

	// Get total count for pagination
	var total int
	database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM scores WHERE user_id = ?", userID).Scan(&total)

	// Get game history
	rows, err := database.DB.QueryContext(ctx, `
		SELECT id, score, total_answered, correct_answers, COALESCE(difficulty, 'easy'), time_limit, COALESCE(time_spent, 0), played_at
		FROM scores
		WHERE user_id = ?
		ORDER BY played_at DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}
	defer rows.Close()

	var history []models.Score
	for rows.Next() {
		var s models.Score
		s.UserID = userID
		err := rows.Scan(&s.ID, &s.Score, &s.TotalAnswered, &s.CorrectAnswers, &s.Difficulty, &s.TimeLimit, &s.TimeSpent, &s.PlayedAt)
		if err != nil {
			continue
		}
		history = append(history, s)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read history"})
		return
	}
	if history == nil {
		history = []models.Score{}
	}

	c.JSON(http.StatusOK, models.ProfileResponse{
		User: models.UserPublic{
			ID:          userID,
			Username:    username,
			TotalPoints: totalPoints,
		},
		Stats:   stats,
		History: history,
		Page:    page,
		Total:   total,
	})
}
