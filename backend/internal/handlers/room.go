package handlers

import (
	"encoding/json"
	"log/slog"
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
	MatchSessionHandler
	roomService *services.RoomService
	upgrader    websocket.Upgrader
	connections sync.Map // userID → *websocket.Conn
}

// NewRoomHandler creates a new room handler.
func NewRoomHandler(roomService *services.RoomService, matchService *services.RankedMatchService, allowedOrigins map[string]bool) *RoomHandler {
	return &RoomHandler{
		MatchSessionHandler: MatchSessionHandler{MatchService: matchService},
		roomService:         roomService,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" {
					return true
				}
				return allowedOrigins[origin]
			},
		},
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

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn("room ws: upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	h.connections.Store(userID, conn)
	defer h.connections.Delete(userID)

	slog.Info("room ws: user connected", "username", username, "userID", userID, "action", action)

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

	sendWSMessage(conn, models.WSMessage{
		Type: "ROOM_CREATED",
		Payload: models.RoomCreatedPayload{
			Code:    roomCode,
			Message: "รอผู้เล่นเข้าร่วม...",
		},
	})

	// Single pump goroutine owns conn.ReadMessage() for the connection lifetime.
	msgChan := StartMsgPump(conn)

	// Host wait loop: select on inbound messages and a ticker that polls for a
	// guest. No read-deadline tricks — no concurrent readers.
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
createWait:
	for {
		select {
		case raw, ok := <-msgChan:
			if !ok || raw.Err != nil {
				h.roomService.RemoveRoom(roomCode)
				return
			}
			var msg models.WSMessage
			if err := json.Unmarshal(raw.Data, &msg); err == nil && msg.Type == models.EventPing {
				sendWSMessage(conn, models.WSMessage{Type: models.EventPong})
			}
			// Check immediately after a message in case the guest just joined.
			if room, ok := h.roomService.GetRoom(roomCode); ok && room.GuestID != 0 {
				break createWait
			}
		case <-ticker.C:
			room, ok := h.roomService.GetRoom(roomCode)
			if !ok {
				return
			}
			if room.GuestID != 0 {
				break createWait
			}
		}
	}
	ticker.Stop()

	room, ok := h.roomService.GetRoom(roomCode)
	if !ok {
		return
	}

	match, err := h.roomService.StartMatch(roomCode)
	if err != nil {
		sendWSMessage(conn, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "ไม่สามารถเริ่มเกมได้"},
		})
		return
	}

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

	h.SendNextQuestionWithCallback(conn, match.MatchID, hostID, h.sendMatchEnd)
	h.HandleMatchSession(conn, msgChan, match.MatchID, hostID, hostName, h.sendMatchEnd)
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

	sendWSMessage(conn, models.WSMessage{
		Type: "ROOM_JOINED",
		Payload: models.RoomJoinedPayload{
			Code:    code,
			Message: "เข้าร่วมห้องสำเร็จ กำลังเริ่มเกม...",
		},
	})

	// Single pump goroutine owns conn.ReadMessage() for the connection lifetime.
	msgChan := StartMsgPump(conn)

	// Guest wait loop: select on inbound messages and a ticker that polls for the
	// match to start. No read-deadline tricks — no concurrent readers.
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	waitStart := time.Now()
	var matchID int64
joinWait:
	for {
		if time.Since(waitStart) > 30*time.Second {
			sendWSMessage(conn, models.WSMessage{
				Type:    models.EventError,
				Payload: models.ErrorPayload{Message: "หมดเวลารอ"},
			})
			return
		}
		select {
		case raw, ok := <-msgChan:
			if !ok || raw.Err != nil {
				return
			}
			var msg models.WSMessage
			if err := json.Unmarshal(raw.Data, &msg); err == nil && msg.Type == models.EventPing {
				sendWSMessage(conn, models.WSMessage{Type: models.EventPong})
			}
			// Check immediately after a message.
			if room, ok := h.roomService.GetRoom(code); ok && room.Status == "active" && room.MatchID != 0 {
				matchID = room.MatchID
				break joinWait
			}
		case <-ticker.C:
			room, ok := h.roomService.GetRoom(code)
			if !ok {
				return
			}
			if room.Status == "active" && room.MatchID != 0 {
				matchID = room.MatchID
				break joinWait
			}
		}
	}
	ticker.Stop()

	room, ok := h.roomService.GetRoom(code)
	if !ok {
		return
	}

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

	h.SendNextQuestionWithCallback(conn, matchID, guestID, h.sendMatchEnd)
	h.HandleMatchSession(conn, msgChan, matchID, guestID, guestName, h.sendMatchEnd)
}

// sendMatchEnd sends MATCH_END to both players (no rating changes).
func (h *RoomHandler) sendMatchEnd(matchID int64) {
	match, ok := h.MatchService.GetMatch(matchID)
	if !ok {
		return
	}

	p1Payload := h.getCustomMatchEndPayload(match)
	p2Payload := h.getCustomMatchEndPayload(match)

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

// getCustomMatchEndPayload creates a MatchEndPayload with ratingChange=0.
func (h *RoomHandler) getCustomMatchEndPayload(match *services.ActiveRankedMatch) *models.MatchEndPayload {
	match.Mu().Lock()
	defer match.Mu().Unlock()

	p1 := match.Player1
	p2 := match.Player2

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
		RatingChange: 0,
		NewRating:    0,
	}
}
