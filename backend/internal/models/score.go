package models

import "time"

// Score represents a game score entry.
type Score struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"userId"`
	Score          int       `json:"score"`
	TotalAnswered  int       `json:"totalAnswered"`
	CorrectAnswers int       `json:"correctAnswers"`
	Difficulty     string    `json:"difficulty"`
	TimeLimit      int       `json:"timeLimit"`
	TimeSpent      float64   `json:"timeSpent"`
	PlayedAt       time.Time `json:"playedAt"`
}

// SubmitScoreRequest is the payload for submitting a game score.
type SubmitScoreRequest struct {
	Score          int     `json:"score" binding:"required,min=0"`
	TotalAnswered  int     `json:"totalAnswered" binding:"required,min=0"`
	CorrectAnswers int     `json:"correctAnswers" binding:"required,min=0"`
	Difficulty     string  `json:"difficulty" binding:"required"`
	TimeLimit      int     `json:"timeLimit" binding:"required,min=1"`
	TimeSpent      float64 `json:"timeSpent"`
}

// LeaderboardEntry represents a single row on the leaderboard (by user total points).
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	Username    string `json:"username"`
	UserID      int64  `json:"userId"`
	TotalPoints int    `json:"totalPoints"`
	TotalGames  int    `json:"totalGames"`
}

// LeaderboardResponse wraps leaderboard data.
type LeaderboardResponse struct {
	Entries []LeaderboardEntry `json:"entries"`
}

// UserStats represents a user's aggregated statistics.
type UserStats struct {
	TotalGames      int     `json:"totalGames"`
	TotalPoints     int     `json:"totalPoints"`
	HighScore       int     `json:"highScore"`
	TotalCorrect    int     `json:"totalCorrect"`
	TotalAnswered   int     `json:"totalAnswered"`
	AverageScore    float64 `json:"averageScore"`
	AverageAccuracy float64 `json:"averageAccuracy"`
	FastestTime     float64 `json:"fastestTime"`
	AverageTime     float64 `json:"averageTime"`
}

// ProfileResponse represents the full profile data.
type ProfileResponse struct {
	User    UserPublic `json:"user"`
	Stats   UserStats  `json:"stats"`
	History []Score    `json:"history"`
	Page    int        `json:"page"`
	Total   int        `json:"total"`
}
