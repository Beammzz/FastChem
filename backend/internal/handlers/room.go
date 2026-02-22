package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/takumi/fastchem/internal/middleware"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// RoomHandler handles custom room WebSocket connections.
type RoomHandler struct {
	roomService  *services.RoomService
	matchService *services.RankedMatchService
	connections  sync.Map // userID → *websocket.Conn
}

// NewRoomHandler creates a new room handler.
func NewRoomHandler(roomService *services.RoomService, matchService *services.RankedMatchService) *RoomHandler {
	return &RoomHandler{
		roomService:  roomService,
		matchService: matchService,
	}
}

// HandleWebSocket handles the WebSocket connection for custom rooms.
// Route: GET /api/room/ws?token=<jwt>&action=create  OR  &action=join&code=XXXXXX
func (h *RoomHandler) HandleWebSocket(c *gin.Context) {
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

	action := c.Query("action")
	code := c.Query("code")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("room ws: upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	log.Printf("room ws: user %s (%d) connected, action=%s", username, userID, action)

	switch action {
	case "create":
		h.handleCreate(conn, userID, username)
	case "join":
		h.handleJoin(conn, userID, username, code)
	default:
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Invalid action. Use 'create' or 'join'"},
		})
	}
}

// handleCreate handles a host creating a room and waiting for a guest.
func (h *RoomHandler) handleCreate(conn *websocket.Conn, hostID int64, hostName string) {
	roomCode := h.roomService.CreateRoom(hostID, hostName)

	// Send room code to host
	sendWSMessage(conn, models.WSMessage{
		Type: "ROOM_CREATED",
		Payload: map[string]interface{}{
			"code":    roomCode,
			"message": "รอผู้เล่นเข้าร่วม...",
		},
	})

	// Wait for guest to join (poll) or handle pings
	guestJoined := make(chan struct{})
	readDone := make(chan struct{})

	go func() {
		defer close(readDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// Host disconnected
				h.roomService.RemoveRoom(roomCode)
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

	// Poll for guest joining
	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-readDone:
				return
			case <-ticker.C:
				room, ok := h.roomService.GetRoom(roomCode)
				if !ok {
					return
				}
				if room.GuestID != 0 {
					close(guestJoined)
					return
				}
			}
		}
	}()

	// Wait for guest or disconnection
	select {
	case <-guestJoined:
		// Guest joined!
	case <-readDone:
		// Host disconnected
		return
	}

	room, ok := h.roomService.GetRoom(roomCode)
	if !ok {
		return
	}

	// Start the match
	match, err := h.roomService.StartMatch(roomCode)
	if err != nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "ไม่สามารถเริ่มเกมได้"},
		})
		return
	}

	// Notify host: match found
	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchFound,
		Payload: models.MatchFoundPayload{
			MatchID:    match.MatchID,
			OpponentID: room.GuestID,
			Opponent:   room.GuestName,
		},
	})

	time.Sleep(1 * time.Second)

	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchStart,
		Payload: models.MatchStartPayload{
			MatchID:        match.MatchID,
			TotalQuestions: models.RankedQuestionsPerMatch,
		},
	})

	// Send first question
	h.sendNextQuestion(conn, match.MatchID, hostID)

	// Run the match session (reuse ranked match logic)
	h.handleMatchSession(conn, match.MatchID, hostID, hostName, roomCode)
}

// handleJoin handles a guest joining a room by code.
func (h *RoomHandler) handleJoin(conn *websocket.Conn, guestID int64, guestName string, code string) {
	if code == "" {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "กรุณาใส่รหัสห้อง"},
		})
		return
	}

	_, err := h.roomService.JoinRoom(code, guestID, guestName)
	if err != nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: err.Error()},
		})
		return
	}

	// Send waiting message
	sendWSMessage(conn, models.WSMessage{
		Type: "ROOM_JOINED",
		Payload: map[string]interface{}{
			"code":    code,
			"message": "เข้าร่วมห้องสำเร็จ กำลังเริ่มเกม...",
		},
	})

	// Wait for match to start (host triggers start)
	matchReady := make(chan int64)
	readDone := make(chan struct{})

	go func() {
		defer close(readDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
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

	go func() {
		ticker := time.NewTicker(300 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-readDone:
				return
			case <-ticker.C:
				room, ok := h.roomService.GetRoom(code)
				if !ok {
					return
				}
				if room.Status == "active" && room.MatchID != 0 {
					matchReady <- room.MatchID
					return
				}
			}
		}
	}()

	var matchID int64
	select {
	case mid := <-matchReady:
		matchID = mid
	case <-readDone:
		return
	case <-time.After(30 * time.Second):
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "หมดเวลารอ"},
		})
		return
	}

	room, ok := h.roomService.GetRoom(code)
	if !ok {
		return
	}

	// Send match found to guest
	sendWSMessage(conn, models.WSMessage{
		Type: models.EventMatchFound,
		Payload: models.MatchFoundPayload{
			MatchID:    matchID,
			OpponentID: room.HostID,
			Opponent:   room.HostName,
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

	// Send first question
	h.sendNextQuestion(conn, matchID, guestID)

	// Run match session
	h.handleMatchSession(conn, matchID, guestID, guestName, code)
}

// handleMatchSession manages the WebSocket session for an active custom match.
// This is essentially the same as ranked but without rating changes.
func (h *RoomHandler) handleMatchSession(conn *websocket.Conn, matchID int64, userID int64, username string, roomCode string) {
	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Match not found"},
		})
		return
	}

	var sendChan chan models.WSMessage
	if userID == match.Player1.UserID {
		sendChan = match.P1Send
	} else {
		sendChan = match.P2Send
	}

	// Start disconnect checker
	stopChecker := make(chan struct{})
	go h.disconnectChecker(matchID, userID, stopChecker)
	defer close(stopChecker)

	// Start outbound message pump
	outDone := make(chan struct{})
	go func() {
		defer close(outDone)
		for {
			select {
			case msg, ok := <-sendChan:
				if !ok {
					return // channel closed
				}
				if err := sendWSMessage(conn, msg); err != nil {
					return
				}
			case <-match.Done:
				return
			}
		}
	}()

	// Read loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("room ws: read error for user %d: %v", userID, err)
			h.matchService.HandleDisconnect(matchID, userID)
			break
		}

		// If match was removed, stop processing
		select {
		case <-match.Done:
			return
		default:
		}

		var msg models.WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case models.EventPing:
			match.SafeSend(sendChan, models.WSMessage{Type: models.EventPong})

		case models.EventSubmitAnswer:
			h.handleSubmitAnswer(match, sendChan, matchID, userID, msg.Payload)

		default:
			match.SafeSend(sendChan, models.WSMessage{
				Type:    models.EventError,
				Payload: models.ErrorPayload{Message: "Unknown event type"},
			})
		}
	}

	<-outDone
}

// handleSubmitAnswer processes a SUBMIT_ANSWER event from a player in a custom match.
func (h *RoomHandler) handleSubmitAnswer(match *services.ActiveRankedMatch, sendChan chan models.WSMessage, matchID int64, userID int64, rawPayload interface{}) {
	payloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return
	}
	var payload models.SubmitAnswerPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		match.SafeSend(sendChan, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Invalid answer payload"},
		})
		return
	}

	result, playerDone, err := h.matchService.SubmitAnswer(matchID, userID, payload.Index, payload.SelectedIndex)
	if err != nil || result == nil {
		match.SafeSend(sendChan, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Failed to process answer"},
		})
		return
	}

	match.SafeSendTimeout(sendChan, models.WSMessage{Type: "ANSWER_RESULT", Payload: result}, 5*time.Second)

	h.notifyOpponentProgress(matchID, userID)

	if playerDone {
		if h.matchService.FinalizeIfBothDone(matchID) {
			h.sendMatchEnd(matchID)
		}
	} else {
		go h.tryAdvanceBothPlayers(matchID)
	}
}

// sendNextQuestion sends the next question to a single player.
func (h *RoomHandler) sendNextQuestion(conn *websocket.Conn, matchID int64, userID int64) {
	q, idx, ok := h.matchService.GetCurrentQuestion(matchID, userID)
	if !ok {
		return
	}

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

	h.startQuestionTimeout(matchID, userID, idx, q.TimeLimit)
}

// startQuestionTimeout starts a goroutine that auto-submits a timeout answer.
func (h *RoomHandler) startQuestionTimeout(matchID int64, userID int64, questionIdx int, timeLimitSec int) {
	go func() {
		timer := time.NewTimer(time.Duration(timeLimitSec) * time.Second)
		defer timer.Stop()
		<-timer.C

		_, currentIdx, exists := h.matchService.GetCurrentQuestion(matchID, userID)
		if !exists || currentIdx != questionIdx {
			return
		}

		result, playerDone, _ := h.matchService.SubmitAnswer(matchID, userID, questionIdx, -1)
		if result == nil {
			return
		}

		match, ok := h.matchService.GetMatch(matchID)
		if !ok {
			return
		}
		var sendChan chan models.WSMessage
		if userID == match.Player1.UserID {
			sendChan = match.P1Send
		} else {
			sendChan = match.P2Send
		}
		match.SafeSendTimeout(sendChan, models.WSMessage{Type: "ANSWER_RESULT", Payload: result}, 5*time.Second)

		h.notifyOpponentProgress(matchID, userID)

		if playerDone {
			if h.matchService.FinalizeIfBothDone(matchID) {
				h.sendMatchEnd(matchID)
			}
		} else {
			h.tryAdvanceBothPlayers(matchID)
		}
	}()
}

// tryAdvanceBothPlayers advances both players to the next question.
func (h *RoomHandler) tryAdvanceBothPlayers(matchID int64) {
	time.Sleep(2 * time.Second)

	var ready bool
	var idx int
	var q *services.RankedQuestion
	for attempt := 0; attempt < 3; attempt++ {
		ready, idx, q = h.matchService.TryAdvanceQuestion(matchID)
		if ready && q != nil {
			break
		}
		match, ok := h.matchService.GetMatch(matchID)
		if !ok || match.Status != "active" {
			return
		}
		if attempt < 2 {
			time.Sleep(1 * time.Second)
		}
	}
	if !ready || q == nil {
		return
	}

	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		return
	}

	msg := models.WSMessage{
		Type: models.EventQuestionStart,
		Payload: models.QuestionStartPayload{
			Index:      idx,
			Question:   q.Question,
			Choices:    q.Choices,
			TimeLimit:  q.TimeLimit,
			Difficulty: q.Difficulty,
			Category:   q.Category,
		},
	}

	match.SafeSendTimeout(match.P1Send, msg, 5*time.Second)
	match.SafeSendTimeout(match.P2Send, msg, 5*time.Second)

	h.startQuestionTimeout(matchID, match.Player1.UserID, idx, q.TimeLimit)
	h.startQuestionTimeout(matchID, match.Player2.UserID, idx, q.TimeLimit)
}

// notifyOpponentProgress sends a progress update to the opponent.
func (h *RoomHandler) notifyOpponentProgress(matchID int64, userID int64) {
	answered, totalScore := h.matchService.GetPlayerProgress(matchID, userID)

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

	msg := models.WSMessage{
		Type: "OPPONENT_PROGRESS",
		Payload: map[string]interface{}{
			"answered":   answered,
			"totalScore": totalScore,
		},
	}
	match.SafeSend(opponentChan, msg)
}

// sendMatchEnd sends MATCH_END to both players (no rating changes).
func (h *RoomHandler) sendMatchEnd(matchID int64) {
	match, ok := h.matchService.GetMatch(matchID)
	if !ok {
		return
	}

	p1Payload := h.getCustomMatchEndPayload(match, match.Player1.UserID)
	p2Payload := h.getCustomMatchEndPayload(match, match.Player2.UserID)

	if p1Payload != nil {
		match.SafeSendTimeout(match.P1Send, models.WSMessage{Type: models.EventMatchEnd, Payload: p1Payload}, 5*time.Second)
	}
	if p2Payload != nil {
		match.SafeSendTimeout(match.P2Send, models.WSMessage{Type: models.EventMatchEnd, Payload: p2Payload}, 5*time.Second)
	}

	go func() {
		time.Sleep(5 * time.Second)
		h.matchService.RemoveMatch(matchID)
	}()
}

// getCustomMatchEndPayload creates a MatchEndPayload with ratingChange=0.
func (h *RoomHandler) getCustomMatchEndPayload(match *services.ActiveRankedMatch, userID int64) *models.MatchEndPayload {
	match.Mu().Lock()
	defer match.Mu().Unlock()

	p1 := match.Player1
	p2 := match.Player2

	// Determine winner by score
	var winnerID int64
	var winnerUsername string
	if p1.TotalScore >= p2.TotalScore {
		winnerID = p1.UserID
		winnerUsername = p1.Username
	}
	if p2.TotalScore > p1.TotalScore {
		winnerID = p2.UserID
		winnerUsername = p2.Username
	}

	return &models.MatchEndPayload{
		MatchID:  match.MatchID,
		Winner:   winnerUsername,
		WinnerID: winnerID,
		Player1: models.MatchEndPlayerSummary{
			UserID:         p1.UserID,
			Username:       p1.Username,
			TotalScore:     p1.TotalScore,
			CorrectAnswers: p1.Correct,
			HardScore:      p1.HardScore,
			TotalTime:      p1.TotalTime,
			Combo:          p1.BestCombo,
		},
		Player2: models.MatchEndPlayerSummary{
			UserID:         p2.UserID,
			Username:       p2.Username,
			TotalScore:     p2.TotalScore,
			CorrectAnswers: p2.Correct,
			HardScore:      p2.HardScore,
			TotalTime:      p2.TotalTime,
			Combo:          p2.BestCombo,
		},
		RatingChange: 0, // No rating change for custom rooms
		NewRating:    0, // Not applicable
	}
}

// disconnectChecker periodically checks if the opponent has disconnected.
func (h *RoomHandler) disconnectChecker(matchID int64, userID int64, stop chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			timedOut, loserID := h.matchService.CheckDisconnectTimeout(matchID)
			if timedOut {
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
