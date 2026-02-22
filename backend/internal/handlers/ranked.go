package handlers

import (
	"encoding/json"
	"log"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from the frontend
		return true
	},
}

// RankedHandler handles ranked match WebSocket connections.
type RankedHandler struct {
	matchService *services.RankedMatchService
	queue        *services.MatchmakingQueue
	// Track active WebSocket connections by userID
	connections sync.Map // userID → *websocket.Conn
}

// NewRankedHandler creates a new ranked handler.
func NewRankedHandler(matchService *services.RankedMatchService, queue *services.MatchmakingQueue) *RankedHandler {
	return &RankedHandler{
		matchService: matchService,
		queue:        queue,
	}
}

// HandleWebSocket handles the WebSocket connection for ranked matches.
// Route: GET /api/ranked/ws?token=<jwt>
func (h *RankedHandler) HandleWebSocket(c *gin.Context) {
	// Authenticate via query parameter (WebSocket can't send headers easily)
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

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ranked ws: upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Store connection
	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	log.Printf("ranked ws: user %s (%d) connected", username, userID)

	// Get user rating
	rating, _, _, _, dbErr := services.GetUserRankedStats(userID)
	if dbErr != nil {
		rating = models.DefaultRating
	}

	// Check if user is already in an active match (reconnect)
	if existingMatch, ok := h.matchService.GetMatchByUserID(userID); ok {
		if existingMatch.Status == "active" {
			reconnected := h.matchService.HandleReconnect(existingMatch.MatchID, userID)
			if reconnected {
				h.handleMatchSession(conn, existingMatch.MatchID, userID, username)
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

	// Notify client they've joined the queue
	sendWSMessage(conn, models.WSMessage{
		Type:    models.EventQueueJoined,
		Payload: map[string]interface{}{"message": "Searching for opponent...", "rating": rating},
	})

	// Wait for a match or disconnection
	// Also handle pings from client while waiting
	matchID := int64(0)
	waitDone := make(chan struct{})

	// Read loop for ping/cancel while in queue
	go func() {
		defer close(waitDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Client disconnected while in queue
				h.queue.Dequeue(userID)
				return
			}
			var msg models.WSMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}
			if msg.Type == models.EventPing {
				sendWSMessage(conn, models.WSMessage{Type: models.EventPong})
			}
		}
	}()

	// Wait for match
	select {
	case id, ok := <-matchedChan:
		if !ok {
			// Channel closed — cancelled
			return
		}
		matchID = id
	case <-waitDone:
		// Client disconnected
		return
	}

	// Match found!
	h.handleMatchSession(conn, matchID, userID, username)
}

// handleMatchSession manages the WebSocket session for an active match.
func (h *RankedHandler) handleMatchSession(conn *websocket.Conn, matchID int64, userID int64, username string) {
	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Match not found"},
		})
		return
	}

	// Determine opponent
	var opponentID int64
	var opponentName string
	var sendChan chan models.WSMessage

	if userID == match.Player1.UserID {
		opponentID = match.Player2.UserID
		opponentName = match.Player2.Username
		sendChan = match.P1Send
	} else {
		opponentID = match.Player1.UserID
		opponentName = match.Player1.Username
		sendChan = match.P2Send
	}

	// Send MATCH_FOUND
	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchFound,
		Payload: models.MatchFoundPayload{
			MatchID:    matchID,
			OpponentID: opponentID,
			Opponent:   opponentName,
		},
	})

	// Brief delay, then send MATCH_START
	time.Sleep(1 * time.Second)

	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchStart,
		Payload: models.MatchStartPayload{
			MatchID:        matchID,
			TotalQuestions: models.RankedQuestionsPerMatch,
		},
	})

	// Send first question
	h.sendNextQuestion(conn, matchID, userID)

	// Start disconnect checker goroutine for this match
	stopChecker := make(chan struct{})
	go h.disconnectChecker(matchID, userID, stopChecker)
	defer close(stopChecker)

	// Start the outbound message pump (for server-initiated messages)
	outDone := make(chan struct{})
	go func() {
		defer close(outDone)
		for msg := range sendChan {
			if err := sendWSMessage(conn, msg); err != nil {
				return
			}
		}
	}()

	// Read loop — process client messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ranked ws: read error for user %d: %v", userID, err)
			h.matchService.HandleDisconnect(matchID, userID)
			break
		}

		var msg models.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case models.EventPing:
			sendWSMessage(conn, models.WSMessage{Type: models.EventPong})

		case models.EventSubmitAnswer:
			h.handleSubmitAnswer(conn, matchID, userID, msg.Payload)

		default:
			sendWSMessage(conn, models.WSMessage{
				Type:    models.EventError,
				Payload: models.ErrorPayload{Message: "Unknown event type"},
			})
		}
	}

	// Wait for outbound pump to finish
	<-outDone
}

// handleSubmitAnswer processes a SUBMIT_ANSWER event from a player.
func (h *RankedHandler) handleSubmitAnswer(conn *websocket.Conn, matchID int64, userID int64, rawPayload interface{}) {
	// Parse payload
	payloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return
	}
	var payload models.SubmitAnswerPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Invalid answer payload"},
		})
		return
	}

	// Submit to match service
	result, playerDone, err := h.matchService.SubmitAnswer(
		matchID, userID, payload.Index, payload.SelectedIndex,
	)
	if err != nil || result == nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Failed to process answer"},
		})
		return
	}

	// Send answer result
	sendWSMessage(conn, models.WSMessage{
		Type:    "ANSWER_RESULT",
		Payload: result,
	})

	// Notify opponent of progress update (via their send channel)
	h.notifyOpponentProgress(matchID, userID, result)

	if playerDone {
		// Check if both players are done
		if h.matchService.FinalizeIfBothDone(matchID) {
			h.sendMatchEnd(matchID)
		}
		// If only this player is done, they wait for opponent
	} else {
		// Send next question after a brief delay
		time.Sleep(500 * time.Millisecond)
		h.sendNextQuestion(conn, matchID, userID)
	}
}

// sendNextQuestion sends the next question to a player.
func (h *RankedHandler) sendNextQuestion(conn *websocket.Conn, matchID int64, userID int64) {
	q, idx, ok := h.matchService.GetCurrentQuestion(matchID, userID)
	if !ok {
		return // no more questions
	}

	// Mark question start time
	h.matchService.StartQuestion(matchID, userID)

	sendWSMessage(conn, models.WSMessage{
		Type: models.EventQuestionStart,
		Payload: models.QuestionStartPayload{
			Index:      idx,
			Question:   q.Question,
			Choices:    q.Choices,
			TimeLimit:  q.TimeLimit,
			Difficulty: q.Difficulty,
			Category:   q.Category,
		},
	})

	// Start a timer for this question — auto-submit if player doesn't answer
	go func() {
		timer := time.NewTimer(time.Duration(q.TimeLimit) * time.Second)
		defer timer.Stop()

		<-timer.C

		// Check if the player has already answered this question
		_, currentIdx, exists := h.matchService.GetCurrentQuestion(matchID, userID)
		if !exists || currentIdx != idx {
			return // already answered or match ended
		}

		// Auto-submit with wrong answer (timeout)
		result, playerDone, _ := h.matchService.SubmitAnswer(matchID, userID, idx, -1)
		if result != nil {
			sendWSMessage(conn, models.WSMessage{
				Type:    "ANSWER_RESULT",
				Payload: result,
			})
			h.notifyOpponentProgress(matchID, userID, result)

			if playerDone {
				if h.matchService.FinalizeIfBothDone(matchID) {
					h.sendMatchEnd(matchID)
				}
			} else {
				time.Sleep(500 * time.Millisecond)
				h.sendNextQuestion(conn, matchID, userID)
			}
		}
	}()
}

// notifyOpponentProgress sends a progress update to the opponent.
func (h *RankedHandler) notifyOpponentProgress(matchID int64, userID int64, result *models.AnswerResultPayload) {
	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		return
	}

	var opponentChan chan models.WSMessage
	if userID == match.Player1.UserID {
		opponentChan = match.P2Send
	} else {
		opponentChan = match.P1Send
	}

	// Non-blocking send
	msg := models.WSMessage{
		Type: "OPPONENT_PROGRESS",
		Payload: map[string]interface{}{
			"answered":   result.OpponentAnswered,
			"totalScore": result.OpponentScore,
		},
	}
	select {
	case opponentChan <- msg:
	default:
		// Channel full, skip
	}
}

// sendMatchEnd sends the MATCH_END event to both players.
func (h *RankedHandler) sendMatchEnd(matchID int64) {
	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		return
	}

	// Get payloads for each player
	p1Payload := h.matchService.GetMatchEndPayload(matchID, match.Player1.UserID)
	p2Payload := h.matchService.GetMatchEndPayload(matchID, match.Player2.UserID)

	if p1Payload != nil {
		select {
		case match.P1Send <- models.WSMessage{Type: models.EventMatchEnd, Payload: p1Payload}:
		default:
		}
	}
	if p2Payload != nil {
		select {
		case match.P2Send <- models.WSMessage{Type: models.EventMatchEnd, Payload: p2Payload}:
		default:
		}
	}

	// Clean up after a delay to let messages flush
	go func() {
		time.Sleep(5 * time.Second)
		h.matchService.RemoveMatch(matchID)
	}()
}

// disconnectChecker periodically checks if the opponent has disconnected too long.
func (h *RankedHandler) disconnectChecker(matchID int64, userID int64, stop chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			timedOut, loserID := h.matchService.CheckDisconnectTimeout(matchID)
			if timedOut {
				// Determine winner
				match, ok := h.matchService.GetMatch(matchID)
				if !ok {
					return
				}
				var winnerID int64
				if loserID == match.Player1.UserID {
					winnerID = match.Player2.UserID
				} else {
					winnerID = match.Player1.UserID
				}

				h.matchService.ForceWin(matchID, winnerID)
				h.sendMatchEnd(matchID)
				return
			}
		}
	}
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

	rows, err := database.DB.Query(`
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

	if history == nil {
		history = []MatchHistoryEntry{}
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// GetRankedLeaderboard handles GET /api/ranked/leaderboard.
func (h *RankedHandler) GetRankedLeaderboard(c *gin.Context) {
	rows, err := database.DB.Query(`
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
