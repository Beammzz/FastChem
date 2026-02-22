package models

import "time"

// RankedQuestionsPerMatch is the fixed number of questions in a ranked match.
const RankedQuestionsPerMatch = 10

// Ranked difficulty distribution: 4 easy, 3 medium, 3 hard.
const (
	RankedEasyCount   = 4
	RankedMediumCount = 3
	RankedHardCount   = 3
)

// Default ELO constants.
const (
	DefaultRating = 1200
	EloKFactor    = 24
)

// RankedMatch represents a 1v1 ranked match stored in the database.
type RankedMatch struct {
	ID           int64      `json:"id"`
	Player1ID    int64      `json:"player1Id"`
	Player2ID    int64      `json:"player2Id"`
	Seed         int64      `json:"seed"`
	Player1Score int        `json:"player1Score"`
	Player2Score int        `json:"player2Score"`
	WinnerID     *int64     `json:"winnerId"`
	Status       string     `json:"status"` // "active", "finished", "abandoned"
	CreatedAt    time.Time  `json:"createdAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
}

// RankedQuestionResult records one answered question within a ranked match.
type RankedQuestionResult struct {
	ID            int64   `json:"id"`
	RankedMatchID int64   `json:"rankedMatchId"`
	UserID        int64   `json:"userId"`
	QuestionIndex int     `json:"questionIndex"`
	Difficulty    string  `json:"difficulty"`
	Topic         string  `json:"topic"`
	Correct       bool    `json:"correct"`
	TimeSpent     float64 `json:"timeSpent"` // seconds
	ScoreEarned   int     `json:"scoreEarned"`
}

// ─── WebSocket event types ───────────────────────────────────────

// WSEventType enumerates WebSocket message types.
type WSEventType string

const (
	EventMatchFound    WSEventType = "MATCH_FOUND"
	EventMatchStart    WSEventType = "MATCH_START"
	EventQuestionStart WSEventType = "QUESTION_START"
	EventSubmitAnswer  WSEventType = "SUBMIT_ANSWER"
	EventMatchEnd      WSEventType = "MATCH_END"
	EventError         WSEventType = "ERROR"
	EventPing          WSEventType = "PING"
	EventPong          WSEventType = "PONG"
	EventQueueJoined   WSEventType = "QUEUE_JOINED"
)

// WSMessage is the envelope for all WebSocket messages.
type WSMessage struct {
	Type    WSEventType `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

// MatchFoundPayload is sent when two players are matched.
type MatchFoundPayload struct {
	MatchID    int64  `json:"matchId"`
	OpponentID int64  `json:"opponentId"`
	Opponent   string `json:"opponent"`
}

// MatchStartPayload is sent when the match begins.
type MatchStartPayload struct {
	MatchID        int64 `json:"matchId"`
	TotalQuestions int   `json:"totalQuestions"`
}

// QuestionStartPayload is sent to deliver a question to the player.
type QuestionStartPayload struct {
	Index      int      `json:"index"` // 0-based
	Question   string   `json:"question"`
	Choices    []string `json:"choices"`
	TimeLimit  int      `json:"timeLimit"`
	Difficulty string   `json:"difficulty"`
	Category   string   `json:"category"`
}

// SubmitAnswerPayload is sent by the client to submit an answer.
type SubmitAnswerPayload struct {
	Index         int `json:"index"`         // question index
	SelectedIndex int `json:"selectedIndex"` // chosen answer
}

// AnswerResultPayload is sent after the server validates an answer.
type AnswerResultPayload struct {
	Index           int     `json:"index"`
	Correct         bool    `json:"correct"`
	CorrectIndex    int     `json:"correctIndex"`
	ScoreEarned     int     `json:"scoreEarned"`
	SpeedBonus      int     `json:"speedBonus"`
	ComboMultiplier float64 `json:"comboMultiplier"`
	Combo           int     `json:"combo"`
	TotalScore      int     `json:"totalScore"`
	TimedOut        bool    `json:"timedOut"`
	TimeSpent       float64 `json:"timeSpent"`
	// Opponent progress (limited info)
	OpponentAnswered int `json:"opponentAnswered"`
	OpponentScore    int `json:"opponentScore"`
}

// MatchEndPayload is sent when the match finishes.
type MatchEndPayload struct {
	MatchID      int64                 `json:"matchId"`
	Winner       string                `json:"winner"` // username
	WinnerID     int64                 `json:"winnerId"`
	Player1      MatchEndPlayerSummary `json:"player1"`
	Player2      MatchEndPlayerSummary `json:"player2"`
	RatingChange int                   `json:"ratingChange"` // for THIS player
	NewRating    int                   `json:"newRating"`    // for THIS player
}

// MatchEndPlayerSummary shows per-player summary at match end.
type MatchEndPlayerSummary struct {
	UserID         int64   `json:"userId"`
	Username       string  `json:"username"`
	TotalScore     int     `json:"totalScore"`
	CorrectAnswers int     `json:"correctAnswers"`
	HardScore      int     `json:"hardScore"`
	TotalTime      float64 `json:"totalTime"`
	Combo          int     `json:"bestCombo"`
}

// ErrorPayload is sent on errors.
type ErrorPayload struct {
	Message string `json:"message"`
}

// QueueJoinedPayload is sent when the player joins the matchmaking queue.
type QueueJoinedPayload struct {
	Message string `json:"message"`
	Rating  int    `json:"rating"`
}

// OpponentProgressPayload is sent to notify a player of their opponent's progress.
type OpponentProgressPayload struct {
	Answered   int `json:"answered"`
	TotalScore int `json:"totalScore"`
}

// RoomCreatedPayload is sent to the host when a room is created.
type RoomCreatedPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RoomJoinedPayload is sent to the guest when they join a room.
type RoomJoinedPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
