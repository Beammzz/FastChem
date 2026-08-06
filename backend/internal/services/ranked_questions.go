package services

import (
	"github.com/takumi/fastchem/internal/models"
)

// RankedQuestionSet holds a deterministically-generated set of questions for a ranked match.
type RankedQuestionSet struct {
	Questions []RankedQuestion
}

// RankedQuestion is a single question with its metadata.
type RankedQuestion struct {
	Question     string   `json:"question"`
	Choices      []string `json:"choices"`
	CorrectIndex int      `json:"correctIndex"`
	TimeLimit    int      `json:"timeLimit"`
	Category     string   `json:"category"`
	Difficulty   string   `json:"difficulty"`
}

// GenerateRankedQuestions builds the question set for a 1v1 match from a seed.
// The same seed always yields the same questions, in the same order, with the
// choices in the same order — both players see one identical set, and a
// finished match can be rebuilt from the seed stored on its row.
//
// Distribution: 4 easy, 3 medium, 3 hard, shuffled.
func GenerateRankedQuestions(seed int64) *RankedQuestionSet {
	source := NewSeededSource(seed)

	plan := []struct {
		difficulty string
		count      int
	}{
		{"easy", models.RankedEasyCount},
		{"medium", models.RankedMediumCount},
		{"hard", models.RankedHardCount},
	}

	questions := make([]RankedQuestion, 0, models.RankedQuestionsPerMatch)
	for _, step := range plan {
		topics := topicsForDifficulty(step.difficulty)
		for i := 0; i < step.count; i++ {
			topic := topics[source.Intn(len(topics))]
			questions = append(questions, problemToRanked(generateValid(source, topic)))
		}
	}

	source.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	return &RankedQuestionSet{Questions: questions}
}

func problemToRanked(p Problem) RankedQuestion {
	return RankedQuestion{
		Question:     p.Question,
		Choices:      p.Choices,
		CorrectIndex: p.CorrectIndex,
		TimeLimit:    GetDifficultyConfig(p.Difficulty).TimeLimit,
		Category:     p.Category,
		Difficulty:   p.Difficulty,
	}
}
