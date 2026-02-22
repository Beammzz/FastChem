package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/takumi/fastchem/internal/models"
	"github.com/takumi/fastchem/internal/services"
)

// QuestionHandler handles question-related HTTP requests
type QuestionHandler struct {
	generator *services.QuestionGenerator
}

// NewQuestionHandler creates a new handler instance
func NewQuestionHandler(gen *services.QuestionGenerator) *QuestionHandler {
	return &QuestionHandler{generator: gen}
}

// GetQuestion handles GET /api/question?difficulty=easy|medium|hard&categories=cat1,cat2
func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	difficulty := c.DefaultQuery("difficulty", "easy")
	categoriesStr := c.Query("categories")

	var question models.Question
	if categoriesStr != "" {
		categories := splitCSV(categoriesStr)
		question = h.generator.GenerateQuestionByCategories(categories)
	} else {
		question = h.generator.GenerateQuestion(difficulty)
	}

	// Store question server-side for anti-cheat validation
	services.GlobalQuestionStore.Store(question.ID, question.CorrectIndex, question.Difficulty)

	c.JSON(http.StatusOK, question)
}

func splitCSV(s string) []string {
	var result []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// ValidateAnswer handles POST /api/validate (legacy endpoint)
func (h *QuestionHandler) ValidateAnswer(c *gin.Context) {
	var req models.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.generator.ValidateAnswer(req)
	c.JSON(http.StatusOK, result)
}

// SubmitAnswer handles POST /api/answer
// The server looks up the question from its store, calculates time_spent,
// validates the answer, and returns the score earned.
func (h *QuestionHandler) SubmitAnswer(c *gin.Context) {
	var req models.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Look up question from server-side store
	stored, ok := services.GlobalQuestionStore.Get(req.QuestionID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown or expired question"})
		return
	}

	// Remove from store (each question can only be answered once)
	services.GlobalQuestionStore.Delete(req.QuestionID)

	// Calculate server-side time_spent
	timeSpent := time.Since(stored.IssuedAt).Seconds()

	cfg := services.GetDifficultyConfig(stored.Difficulty)
	timedOut := timeSpent >= float64(cfg.TimeLimit)

	// Determine correctness
	correct := req.SelectedIndex == stored.CorrectIndex && !timedOut

	// Calculate score
	scoreEarned, speedBonus := services.CalculateScore(stored.Difficulty, timeSpent, correct)

	c.JSON(http.StatusOK, models.AnswerResponse{
		Correct:      correct,
		CorrectIndex: stored.CorrectIndex,
		ScoreEarned:  scoreEarned,
		SpeedBonus:   speedBonus,
		TimeSpent:    timeSpent,
		Difficulty:   stored.Difficulty,
		TimedOut:     timedOut,
	})
}
