package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/takumi/fastchem/internal/models"
)

// QuestionGenerator produces questions for casual play — single player and the
// server-scored match endpoints.
//
// It holds one Source for the lifetime of the process, whose history keeps
// consecutive questions from reusing the same element or compound. Ranked and
// room matches build their own seeded Source instead, but run the same topics
// through the same code path, so the two modes cannot drift apart.
type QuestionGenerator struct {
	source *Source
}

// NewQuestionGenerator creates a generator seeded from the clock.
func NewQuestionGenerator() *QuestionGenerator {
	return &QuestionGenerator{source: NewSeededSource(time.Now().UnixNano())}
}

// GenerateQuestion produces a question for the given difficulty. Unknown
// difficulties fall back to easy.
func (g *QuestionGenerator) GenerateQuestion(difficulty string) models.Question {
	topics := topicsForDifficulty(difficulty)
	topic := topics[g.source.Intn(len(topics))]
	return problemToQuestion(generateValid(g.source, topic))
}

// GenerateQuestionByCategories picks one of the requested categories and
// generates a matching question. Unknown category IDs are ignored; if none
// are recognised it falls back to an easy question.
func (g *QuestionGenerator) GenerateQuestionByCategories(categories []string) models.Question {
	pool := make([]Topic, 0, len(categories))
	for _, category := range categories {
		if topic, ok := topicForCategory(category); ok {
			pool = append(pool, topic)
		}
	}
	if len(pool) == 0 {
		return g.GenerateQuestion("easy")
	}
	topic := pool[g.source.Intn(len(pool))]
	return problemToQuestion(generateValid(g.source, topic))
}

// ValidateAnswer checks if the selected answer is correct (legacy)
func (g *QuestionGenerator) ValidateAnswer(req models.ValidateRequest) models.ValidateResponse {
	correct := req.SelectedIndex == req.CorrectIndex
	msg := "ผิด!"
	if correct {
		msg = "ถูกต้อง!"
	}
	return models.ValidateResponse{
		Correct:      correct,
		CorrectIndex: req.CorrectIndex,
		Message:      msg,
	}
}

// problemToQuestion renders a Problem as the API's question shape. The time
// limit always comes from the topic's own difficulty, so a hard topic served
// through a category filter still gets the hard clock.
func problemToQuestion(p Problem) models.Question {
	return models.Question{
		ID:           uuid.New().String(),
		Question:     p.Question,
		Choices:      p.Choices,
		CorrectIndex: p.CorrectIndex,
		TimeLimit:    GetDifficultyConfig(p.Difficulty).TimeLimit,
		Category:     p.Category,
		Difficulty:   p.Difficulty,
	}
}
