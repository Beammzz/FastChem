package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// MatchHandler handles match-related HTTP requests.
type MatchHandler struct {
	generator *services.QuestionGenerator
}

// NewMatchHandler creates a new match handler.
func NewMatchHandler(gen *services.QuestionGenerator) *MatchHandler {
	return &MatchHandler{generator: gen}
}

// StartMatch handles POST /api/match/start
// Creates a new match in the DB + in-memory store and returns its ID.
func (h *MatchHandler) StartMatch(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.StartMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate difficulty
	cfg := services.GetDifficultyConfig(req.Difficulty)

	// Insert match into DB
	result, err := database.DB.Exec(
		"INSERT INTO matches (user_id, difficulty, total_score, best_combo) VALUES (?, ?, 0, 0)",
		userID, req.Difficulty,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create match"})
		return
	}
	matchID, _ := result.LastInsertId()

	// Register in the active match store
	services.GlobalMatchStore.Create(matchID, userID, req.Difficulty)

	c.JSON(http.StatusCreated, models.StartMatchResponse{
		MatchID:    matchID,
		Difficulty: req.Difficulty,
		Questions:  models.DefaultQuestionsPerMatch,
		TimeLimit:  cfg.TimeLimit,
	})
}

// AnswerMatch handles POST /api/match/answer
// Validates an answer within an active match (combo-aware, anti-cheat).
func (h *MatchHandler) AnswerMatch(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req models.MatchAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify match ownership
	match, ok := services.GlobalMatchStore.Get(req.MatchID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match not found or already ended"})
		return
	}
	if match.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your match"})
		return
	}

	// Look up question from server-side store
	stored, ok := services.GlobalQuestionStore.Get(req.QuestionID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown or expired question"})
		return
	}
	services.GlobalQuestionStore.Delete(req.QuestionID)

	// Calculate server-side time_spent
	timeSpent := time.Since(stored.IssuedAt).Seconds()

	// Record the answer (combo + scoring handled atomically)
	resp, _, err := services.GlobalMatchStore.RecordAnswer(
		req.MatchID, stored, req.QuestionID, req.SelectedIndex, timeSpent,
	)
	if err != nil || resp == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to record answer"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// EndMatch handles POST /api/match/end
// Force-ends a match and returns the final summary.
func (h *MatchHandler) EndMatch(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var body struct {
		MatchID int64 `json:"matchId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	match, ok := services.GlobalMatchStore.Get(body.MatchID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Match not found or already ended"})
		return
	}
	if match.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not your match"})
		return
	}

	resp, err := h.finalizeMatch(body.MatchID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize match"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// finalizeMatch persists match results to DB and cleans up in-memory state.
func (h *MatchHandler) finalizeMatch(matchID, userID int64) (*models.EndMatchResponse, error) {
	match, ok := services.GlobalMatchStore.Get(matchID)
	if !ok {
		return nil, nil
	}

	now := time.Now()

	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update match record
	_, err = tx.Exec(
		"UPDATE matches SET total_score = ?, best_combo = ?, ended_at = ? WHERE id = ?",
		match.TotalScore, match.BestCombo, now, matchID,
	)
	if err != nil {
		return nil, err
	}

	// Insert all question attempts
	for _, a := range match.Attempts {
		_, err = tx.Exec(
			`INSERT INTO question_attempts 
				(match_id, question_snapshot, correct, timed_out, time_spent_ms, score_awarded, combo_at_answer, combo_multiplier)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			matchID, a.QuestionSnapshot, a.Correct, a.TimedOut,
			a.TimeSpentMs, a.ScoreAwarded, a.ComboAtAnswer, a.ComboMultiplier,
		)
		if err != nil {
			return nil, err
		}
	}

	// Calculate actual total time spent from question attempts
	var totalTimeSpentMs int
	for _, a := range match.Attempts {
		totalTimeSpentMs += a.TimeSpentMs
	}
	totalTimeSpentSec := float64(totalTimeSpentMs) / 1000.0

	// Get the per-question time limit for this difficulty
	cfg := services.GetDifficultyConfig(match.Difficulty)

	// Insert into scores table so leaderboard/profile/history track this match
	_, err = tx.Exec(
		`INSERT INTO scores (user_id, score, total_answered, correct_answers, difficulty, time_limit, time_spent)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, match.TotalScore, match.Answered, match.Correct, match.Difficulty,
		cfg.TimeLimit, totalTimeSpentSec,
	)
	if err != nil {
		return nil, err
	}

	// Update user total_points
	_, err = tx.Exec("UPDATE users SET total_points = total_points + ? WHERE id = ?", match.TotalScore, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Get updated total points
	var totalPoints int
	database.DB.QueryRow("SELECT total_points FROM users WHERE id = ?", userID).Scan(&totalPoints)

	resp := &models.EndMatchResponse{
		MatchID:        matchID,
		TotalScore:     match.TotalScore,
		TotalAnswered:  match.Answered,
		CorrectAnswers: match.Correct,
		BestCombo:      match.BestCombo,
		Difficulty:     match.Difficulty,
		Attempts:       match.Attempts,
		TotalPoints:    totalPoints,
	}

	// Clean up in-memory store
	services.GlobalMatchStore.Delete(matchID)

	return resp, nil
}
