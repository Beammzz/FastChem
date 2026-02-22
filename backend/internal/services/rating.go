package services

import (
	"math"

	"github.com/takumi/fastchem/internal/models"
)

// RatingService handles ELO rating calculations and database updates.
type RatingService struct{}

// NewRatingService creates a new RatingService.
func NewRatingService() *RatingService {
	return &RatingService{}
}

// ExpectedScore calculates the expected score for a player.
// Formula: E = 1 / (1 + 10^((opponentRating - playerRating) / 400))
func (rs *RatingService) ExpectedScore(playerRating, opponentRating int) float64 {
	exponent := float64(opponentRating-playerRating) / 400.0
	return 1.0 / (1.0 + math.Pow(10, exponent))
}

// CalculateNewRating computes the new ELO rating after a match.
// actualScore: 1.0 for win, 0.0 for loss (no draws).
func (rs *RatingService) CalculateNewRating(currentRating, opponentRating int, won bool) int {
	expected := rs.ExpectedScore(currentRating, opponentRating)

	var actual float64
	if won {
		actual = 1.0
	} else {
		actual = 0.0
	}

	newRating := float64(currentRating) + float64(models.EloKFactor)*(actual-expected)

	// Floor at 0 to prevent negative ratings
	result := int(math.Round(newRating))
	if result < 0 {
		result = 0
	}
	return result
}

// RankedComboMultiplier returns the combo multiplier for ranked matches.
// 2 correct = 1.1x, 3 correct = 1.2x, 5+ correct = 1.5x.
// Reset on wrong answer.
func RankedComboMultiplier(combo int) float64 {
	switch {
	case combo >= 5:
		return 1.5
	case combo >= 3:
		return 1.2
	case combo >= 2:
		return 1.1
	default:
		return 1.0
	}
}

// CalculateRankedScore computes score for a single ranked question.
// score = base + speed_bonus, where speed_bonus = (timeLimit - timeSpent) / timeLimit * base
// Wrong answer → 0 points, combo reset.
func CalculateRankedScore(difficulty string, timeSpent float64, correct bool, combo int) (finalScore int, speedBonus int, comboMult float64) {
	if !correct {
		return 0, 0, 1.0
	}

	cfg := GetDifficultyConfig(difficulty)
	base := float64(cfg.BaseScore)
	timeLimit := float64(cfg.TimeLimit)

	// Timed out → 0
	if timeSpent >= timeLimit {
		return 0, 0, 1.0
	}

	if timeSpent < 0 {
		timeSpent = 0
	}

	// speed_bonus = (timeLimit - timeSpent) / timeLimit * base
	sb := (timeLimit - timeSpent) / timeLimit * base
	speedBonusInt := int(math.Floor(sb))

	rawScore := int(base) + speedBonusInt

	comboMult = RankedComboMultiplier(combo)
	finalScore = int(math.Floor(float64(rawScore) * comboMult))

	return finalScore, speedBonusInt, comboMult
}
