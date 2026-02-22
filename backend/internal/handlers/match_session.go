package handlers

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// MatchEndFunc is called when a match ends to send results to both players.
// Implementations differ: ranked sends rating changes, rooms send zeroes.
type MatchEndFunc func(matchID int64)

// RawMsg carries a single WebSocket message or a read error from the pump
// goroutine to the consumers (wait phase + HandleMatchSession).
type RawMsg struct {
	Data []byte
	Err  error
}

// StartMsgPump starts a goroutine that owns conn.ReadMessage() for the full
// lifetime of the connection. It returns a channel that receives every message
// (or the terminal error). After an error the channel is closed.
// IMPORTANT: only one pump must exist per connection — never call ReadMessage
// from any other goroutine on the same conn.
func StartMsgPump(conn *websocket.Conn) <-chan RawMsg {
	ch := make(chan RawMsg, 8)
	go func() {
		defer close(ch)
		for {
			_, data, err := conn.ReadMessage()
			ch <- RawMsg{Data: data, Err: err}
			if err != nil {
				return
			}
		}
	}()
	return ch
}

// MatchSessionHandler holds shared logic for ranked and room match WebSocket
// sessions. Both RankedHandler and RoomHandler embed this to avoid code
// duplication.
type MatchSessionHandler struct {
	MatchService *services.RankedMatchService
}

// HandleMatchSession manages the WebSocket session for an active match.
// msgChan must be the channel returned by StartMsgPump for this connection —
// it is the sole source of inbound messages; HandleMatchSession must never
// call conn.ReadMessage() directly.
// sendMatchEnd is a callback that varies between ranked (with rating) and room
// (without rating). The function blocks until the session ends.
func (ms *MatchSessionHandler) HandleMatchSession(
	conn *websocket.Conn,
	msgChan <-chan RawMsg,
	matchID int64,
	userID int64,
	username string,
	sendMatchEnd MatchEndFunc,
) {
	match, ok := ms.MatchService.GetMatch(matchID)
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

	// Start disconnect checker goroutine for this match
	stopChecker := make(chan struct{})
	go ms.disconnectChecker(matchID, userID, stopChecker, sendMatchEnd)
	defer close(stopChecker)

	// Start the outbound message pump (only writer to conn after this point)
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

	// Read loop — process client messages from the pump channel
	for {
		raw, ok := <-msgChan
		if !ok || raw.Err != nil {
			if raw.Err != nil {
				slog.Warn("ws: read error", "userID", userID, "error", raw.Err)
			}
			ms.MatchService.HandleDisconnect(matchID, userID)
			break
		}

		message := raw.Data

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
			// Route through channel to avoid concurrent writes
			match.SafeSend(sendChan, models.WSMessage{Type: models.EventPong})

		case models.EventSubmitAnswer:
			ms.handleSubmitAnswer(match, sendChan, matchID, userID, msg.Payload, sendMatchEnd)

		default:
			match.SafeSend(sendChan, models.WSMessage{
				Type:    models.EventError,
				Payload: models.ErrorPayload{Message: "Unknown event type"},
			})
		}
	}

	// Wait for outbound pump to finish
	<-outDone
}

// handleSubmitAnswer processes a SUBMIT_ANSWER event from a player.
func (ms *MatchSessionHandler) handleSubmitAnswer(
	match *services.ActiveRankedMatch,
	sendChan chan models.WSMessage,
	matchID int64,
	userID int64,
	rawPayload interface{},
	sendMatchEnd MatchEndFunc,
) {
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

	result, playerDone, err := ms.MatchService.SubmitAnswer(
		matchID, userID, payload.Index, payload.SelectedIndex,
	)
	if err != nil || result == nil {
		match.SafeSend(sendChan, models.WSMessage{
			Type:    models.EventError,
			Payload: models.ErrorPayload{Message: "Failed to process answer"},
		})
		return
	}

	match.SafeSendTimeout(sendChan, models.WSMessage{
		Type:    "ANSWER_RESULT",
		Payload: result,
	}, 5*time.Second)

	ms.notifyOpponentProgress(matchID, userID)

	if playerDone {
		if ms.MatchService.FinalizeIfBothDone(matchID) {
			sendMatchEnd(matchID)
		}
	} else {
		go ms.tryAdvanceBothPlayers(matchID, sendMatchEnd)
	}
}

// SendNextQuestion sends the next question to a single player (used for the
// first question only, before the outbound pump is running).
func (ms *MatchSessionHandler) SendNextQuestion(conn *websocket.Conn, matchID int64, userID int64) {
	q, idx, ok := ms.MatchService.GetCurrentQuestion(matchID, userID)
	if !ok {
		return
	}

	ms.MatchService.StartQuestion(matchID, userID)

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

	ms.startQuestionTimeout(matchID, userID, idx, q.TimeLimit, func(matchID int64) {
		// For the first question, we need the sendMatchEnd callback.
		// Since this is called before the outbound pump, we pass a no-op.
		// The real sendMatchEnd is used inside tryAdvanceBothPlayers.
	})
}

// SendNextQuestionWithCallback is like SendNextQuestion but accepts a
// sendMatchEnd callback so the timeout handler can finalize the match.
func (ms *MatchSessionHandler) SendNextQuestionWithCallback(
	conn *websocket.Conn,
	matchID int64,
	userID int64,
	sendMatchEnd MatchEndFunc,
) {
	q, idx, ok := ms.MatchService.GetCurrentQuestion(matchID, userID)
	if !ok {
		return
	}

	ms.MatchService.StartQuestion(matchID, userID)

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

	ms.startQuestionTimeout(matchID, userID, idx, q.TimeLimit, sendMatchEnd)
}

// startQuestionTimeout starts a goroutine that auto-submits a timeout answer.
func (ms *MatchSessionHandler) startQuestionTimeout(
	matchID int64,
	userID int64,
	questionIdx int,
	timeLimitSec int,
	sendMatchEnd MatchEndFunc,
) {
	go func() {
		timer := time.NewTimer(time.Duration(timeLimitSec) * time.Second)
		defer timer.Stop()
		<-timer.C

		_, currentIdx, exists := ms.MatchService.GetCurrentQuestion(matchID, userID)
		if !exists || currentIdx != questionIdx {
			return
		}

		result, playerDone, _ := ms.MatchService.SubmitAnswer(matchID, userID, questionIdx, -1)
		if result == nil {
			return
		}

		match, ok := ms.MatchService.GetMatch(matchID)
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

		ms.notifyOpponentProgress(matchID, userID)

		if playerDone {
			if ms.MatchService.FinalizeIfBothDone(matchID) {
				sendMatchEnd(matchID)
			}
		} else {
			ms.tryAdvanceBothPlayers(matchID, sendMatchEnd)
		}
	}()
}

// tryAdvanceBothPlayers checks if both players have answered and sends the
// next question to both after a brief delay.
func (ms *MatchSessionHandler) tryAdvanceBothPlayers(matchID int64, sendMatchEnd MatchEndFunc) {
	time.Sleep(2 * time.Second)

	var ready bool
	var idx int
	var q *services.RankedQuestion
	for attempt := 0; attempt < 3; attempt++ {
		ready, idx, q = ms.MatchService.TryAdvanceQuestion(matchID)
		if ready && q != nil {
			break
		}
		match, ok := ms.MatchService.GetMatch(matchID)
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

	match, ok := ms.MatchService.GetMatch(matchID)
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

	// Start timeout timers for both players
	ms.startSyncedQuestionTimeout(matchID, match.Player1.UserID, idx, q.TimeLimit, sendMatchEnd)
	ms.startSyncedQuestionTimeout(matchID, match.Player2.UserID, idx, q.TimeLimit, sendMatchEnd)
}

// startSyncedQuestionTimeout starts an auto-timeout for a synchronized question.
func (ms *MatchSessionHandler) startSyncedQuestionTimeout(
	matchID int64,
	userID int64,
	questionIdx int,
	timeLimitSec int,
	sendMatchEnd MatchEndFunc,
) {
	go func() {
		timer := time.NewTimer(time.Duration(timeLimitSec) * time.Second)
		defer timer.Stop()
		<-timer.C

		_, currentIdx, exists := ms.MatchService.GetCurrentQuestion(matchID, userID)
		if !exists || currentIdx != questionIdx {
			return
		}

		result, playerDone, _ := ms.MatchService.SubmitAnswer(matchID, userID, questionIdx, -1)
		if result == nil {
			return
		}

		match, ok := ms.MatchService.GetMatch(matchID)
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

		ms.notifyOpponentProgress(matchID, userID)

		if playerDone {
			if ms.MatchService.FinalizeIfBothDone(matchID) {
				sendMatchEnd(matchID)
			}
		} else {
			ms.tryAdvanceBothPlayers(matchID, sendMatchEnd)
		}
	}()
}

// notifyOpponentProgress sends a progress update to the opponent.
func (ms *MatchSessionHandler) notifyOpponentProgress(matchID int64, userID int64) {
	answered, totalScore := ms.MatchService.GetPlayerProgress(matchID, userID)

	match, ok := ms.MatchService.GetMatch(matchID)
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
		Payload: models.OpponentProgressPayload{
			Answered:   answered,
			TotalScore: totalScore,
		},
	}
	match.SafeSend(opponentChan, msg)
}

// disconnectChecker periodically checks if the opponent has disconnected.
func (ms *MatchSessionHandler) disconnectChecker(
	matchID int64,
	userID int64,
	stop chan struct{},
	sendMatchEnd MatchEndFunc,
) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			timedOut, loserID := ms.MatchService.CheckDisconnectTimeout(matchID)
			if timedOut {
				match, ok := ms.MatchService.GetMatch(matchID)
				if !ok {
					return
				}
				var winnerID int64
				if loserID == match.Player1.UserID {
					winnerID = match.Player2.UserID
				} else {
					winnerID = match.Player1.UserID
				}

				ms.MatchService.ForceWin(matchID, winnerID)
				sendMatchEnd(matchID)
				return
			}
		}
	}
}
