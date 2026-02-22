package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/middleware"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// RankedHandler handles ranked match WebSocket connections.
type RankedHandler struct {
	MatchSessionHandler
	queue    *services.MatchmakingQueue
	upgrader websocket.Upgrader
	// Track active WebSocket connections by userID
	connections sync.Map // userID → *websocket.Conn
}

// NewRankedHandler creates a new ranked handler.
func NewRankedHandler(matchService *services.RankedMatchService, queue *services.MatchmakingQueue, allowedOrigins map[string]bool) *RankedHandler {
	return &RankedHandler{
		MatchSessionHandler: MatchSessionHandler{MatchService: matchService},
		queue:               queue,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true // non-browser clients
				}
				return allowedOrigins[origin]
			},
		},
	}
}

// HandleWebSocket handles the WebSocket connection for ranked matches.
// Route: GET /api/ranked/ws?token=<jwt>
func (h *RankedHandler) HandleWebSocket(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	userID, username, err := middleware.ParseToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn("ranked ws: upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	slog.Info("ranked ws: user connected", "username", username, "userID", userID)

	rating, _, _, _, dbErr := services.GetUserRankedStats(userID)
	if dbErr != nil {
		rating = models.DefaultRating
	}

	// Start the single pump goroutine that owns conn.ReadMessage() for the
	// entire lifetime of this connection. Never call conn.ReadMessage() elsewhere.
	msgChan := StartMsgPump(conn)

	// Check if user is already in an active match (reconnect)
	if existingMatch, ok := h.MatchService.GetMatchByUserID(userID); ok {
		if existingMatch.Status == "active" {
			reconnected := h.MatchService.HandleReconnect(existingMatch.MatchID, userID)
			if reconnected {
				h.HandleMatchSession(conn, msgChan, existingMatch.MatchID, userID, username, h.sendMatchEnd)
				return
			}
		}
	}

	// Enqueue for matchmaking
	matchedChan := h.queue.Enqueue(userID, username, rating)
	if matchedChan == nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Already in a match"},
		})
		return
	}

	sendWSMessage(conn, models.WSMessage{
		Type:    models.EventQueueJoined,
		Payload: models.QueueJoinedPayload{Message: "Searching for opponent...", Rating: rating},
	})

	// Queue-phase wait: select on inbound messages (ping/disconnect) and the
	// matchmaking channel simultaneously — no deadline tricks, no concurrent readers.
	matchID := int64(0)
queueWait:
	for {
		select {
		case raw, ok := <-msgChan:
			if !ok || raw.Err != nil {
				// Connection closed while queuing.
				h.queue.Dequeue(userID)
				return
			}
			var msg models.WSMessage
			if err := json.Unmarshal(raw.Data, &msg); err == nil && msg.Type == models.EventPing {
				sendWSMessage(conn, models.WSMessage{Type: models.EventPong})
			}
		case id, ok := <-matchedChan:
			if !ok {
				return
			}
			matchID = id
			break queueWait
		}
	}

	// Match found — send pre-session messages and delegate to shared handler
	match, ok := h.MatchService.GetMatch(matchID)
	if !ok {
		return
	}

	var opponentID int64
	var opponentName string
	if userID == match.Player1.UserID {
		opponentID = match.Player2.UserID
		opponentName = match.Player2.Username
	} else {
		opponentID = match.Player1.UserID
		opponentName = match.Player1.Username
	}

	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchFound,
		Payload: models.MatchFoundPayload{
			MatchID:    matchID,
			OpponentID: opponentID,
			Opponent:   opponentName,
		},
	})

	time.Sleep(1 * time.Second)

	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchStart,
		Payload: models.MatchStartPayload{
			MatchID:        matchID,
			TotalQuestions: models.RankedQuestionsPerMatch,
		},
	})

	h.SendNextQuestionWithCallback(conn, matchID, userID, h.sendMatchEnd)
	h.HandleMatchSession(conn, msgChan, matchID, userID, username, h.sendMatchEnd)
}

// sendMatchEnd sends MATCH_END to both players with rating changes.
func (h *RankedHandler) sendMatchEnd(matchID int64) {
	match, ok := h.MatchService.GetMatch(matchID)
	if !ok {
		return
	}

	p1Payload := h.MatchService.GetMatchEndPayload(matchID, match.Player1.UserID)
	p2Payload := h.MatchService.GetMatchEndPayload(matchID, match.Player2.UserID)

	if p1Payload != nil {
		match.SafeSendTimeout(match.P1Send, models.WSMessage{Type: models.EventMatchEnd, Payload: p1Payload}, 5*time.Second)
	}
	if p2Payload != nil {
		match.SafeSendTimeout(match.P2Send, models.WSMessage{Type: models.EventMatchEnd, Payload: p2Payload}, 5*time.Second)
	}

	go func() {
		time.Sleep(5 * time.Second)
		h.MatchService.RemoveMatch(matchID)
	}()
}

// GetRankedStats handles GET /api/ranked/stats (authenticated).
func (h *RankedHandler) GetRankedStats(c *gin.Context) {
	userID := c.GetInt64("user_id")

	rating, wins, losses, highestRating, err := services.GetUserRankedStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ranked stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rating":        rating,
		"rankedWins":    wins,
		"rankedLosses":  losses,
		"highestRating": highestRating,
	})
}

// GetRankedHistory handles GET /api/ranked/history (authenticated).
func (h *RankedHandler) GetRankedHistory(c *gin.Context) {
	userID := c.GetInt64("user_id")

	rows, err := database.DB.QueryContext(c.Request.Context(), `
		SELECT rm.id, rm.player1_id, rm.player2_id, rm.player1_score, rm.player2_score,
		       rm.winner_id, rm.status, rm.created_at, rm.finished_at,
		       u1.username as p1_name, u2.username as p2_name
		FROM ranked_matches rm
		JOIN users u1 ON u1.id = rm.player1_id
		JOIN users u2 ON u2.id = rm.player2_id
		WHERE (rm.player1_id = ? OR rm.player2_id = ?) AND rm.status IN ('finished', 'abandoned')
		ORDER BY rm.created_at DESC
		LIMIT 20
	`, userID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch match history"})
		return
	}
	defer rows.Close()

	type MatchHistoryEntry struct {
		MatchID      int64      `json:"matchId"`
		Player1ID    int64      `json:"player1Id"`
		Player2ID    int64      `json:"player2Id"`
		Player1Score int        `json:"player1Score"`
		Player2Score int        `json:"player2Score"`
		WinnerID     *int64     `json:"winnerId"`
		Status       string     `json:"status"`
		CreatedAt    time.Time  `json:"createdAt"`
		FinishedAt   *time.Time `json:"finishedAt"`
		Player1Name  string     `json:"player1Name"`
		Player2Name  string     `json:"player2Name"`
		Won          bool       `json:"won"`
	}

	var history []MatchHistoryEntry
	for rows.Next() {
		var e MatchHistoryEntry
		err := rows.Scan(
			&e.MatchID, &e.Player1ID, &e.Player2ID,
			&e.Player1Score, &e.Player2Score,
			&e.WinnerID, &e.Status, &e.CreatedAt, &e.FinishedAt,
			&e.Player1Name, &e.Player2Name,
		)
		if err != nil {
			continue
		}
		if e.WinnerID != nil {
			e.Won = *e.WinnerID == userID
		}
		history = append(history, e)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read match history"})
		return
	}

	if history == nil {
		history = []MatchHistoryEntry{}
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetRankedLeaderboard handles GET /api/ranked/leaderboard.
func (h *RankedHandler) GetRankedLeaderboard(c *gin.Context) {
	rows, err := database.DB.QueryContext(c.Request.Context(), `
		SELECT u.username, u.id, u.rating, u.ranked_wins, u.ranked_losses, u.highest_rating
		FROM users u
		WHERE u.ranked_wins + u.ranked_losses > 0
		ORDER BY u.rating DESC
		LIMIT 50
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ranked leaderboard"})
		return
	}
	defer rows.Close()

	type RankedLeaderboardEntry struct {
		Rank          int    `json:"rank"`
		Username      string `json:"username"`
		UserID        int64  `json:"userId"`
		Rating        int    `json:"rating"`
		Wins          int    `json:"wins"`
		Losses        int    `json:"losses"`
		HighestRating int    `json:"highestRating"`
	}

	var entries []RankedLeaderboardEntry
	rank := 0
	for rows.Next() {
		rank++
		var e RankedLeaderboardEntry
		err := rows.Scan(&e.Username, &e.UserID, &e.Rating, &e.Wins, &e.Losses, &e.HighestRating)
		if err != nil {
			continue
		}
		e.Rank = rank
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read ranked leaderboard"})
		return
	}

	if entries == nil {
		entries = []RankedLeaderboardEntry{}
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

// sendWSMessage sends a WebSocket message (JSON).
func sendWSMessage(conn *websocket.Conn, msg models.WSMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}
