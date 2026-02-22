package models

import "time"

// DefaultQuestionsPerMatch is the fixed number of questions in a match.
const DefaultQuestionsPerMatch = 10

// Match represents a single game session.
type Match struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"userId"`
	Difficulty string     `json:"difficulty"`
	TotalScore int        `json:"totalScore"`
	BestCombo  int        `json:"bestCombo"`
	StartedAt  time.Time  `json:"startedAt"`
	EndedAt    *time.Time `json:"endedAt"`
}

// QuestionAttempt records one answered question within a match.
type QuestionAttempt struct {
	ID               int64     `json:"id"`
	MatchID          int64     `json:"matchId"`
	QuestionSnapshot string    `json:"questionSnapshot"` // JSON blob
	Correct          bool      `json:"correct"`
	TimedOut         bool      `json:"timedOut"`
	TimeSpentMs      int       `json:"timeSpentMs"`
	ScoreAwarded     int       `json:"scoreAwarded"`
	ComboAtAnswer    int       `json:"comboAtAnswer"`
	ComboMultiplier  float64   `json:"comboMultiplier"`
	CreatedAt        time.Time `json:"createdAt"`
}

// StartMatchRequest is sent by the client to begin a new match.
type StartMatchRequest struct {
	Difficulty string `json:"difficulty" binding:"required"`
}

// StartMatchResponse is returned after creating a match.
type StartMatchResponse struct {
	MatchID    int64  `json:"matchId"`
	Difficulty string `json:"difficulty"`
	Questions  int    `json:"questions"`
	TimeLimit  int    `json:"timeLimit"`
}

// MatchAnswerRequest is sent for each answer during a match.
type MatchAnswerRequest struct {
	MatchID       int64  `json:"matchId" binding:"required"`
	QuestionID    string `json:"questionId" binding:"required"`
	SelectedIndex int    `json:"selectedIndex"`
}

// MatchAnswerResponse is returned after validating a match answer.
type MatchAnswerResponse struct {
	Correct         bool    `json:"correct"`
	CorrectIndex    int     `json:"correctIndex"`
	ScoreEarned     int     `json:"scoreEarned"`
	SpeedBonus      int     `json:"speedBonus"`
	ComboMultiplier float64 `json:"comboMultiplier"`
	Combo           int     `json:"combo"`
	TimeSpent       float64 `json:"timeSpent"`
	TimedOut        bool    `json:"timedOut"`
	TotalScore      int     `json:"totalScore"`
	QuestionNumber  int     `json:"questionNumber"`
}

// EndMatchResponse is returned when a match finishes.
type EndMatchResponse struct {
	MatchID        int64             `json:"matchId"`
	TotalScore     int               `json:"totalScore"`
	TotalAnswered  int               `json:"totalAnswered"`
	CorrectAnswers int               `json:"correctAnswers"`
	BestCombo      int               `json:"bestCombo"`
	Difficulty     string            `json:"difficulty"`
	Attempts       []QuestionAttempt `json:"attempts"`
	TotalPoints    int               `json:"totalPoints"`
}
