package services

import "math"

// DifficultyConfig holds scoring parameters for each difficulty level.
type DifficultyConfig struct {
	TimeLimit  int     // seconds per question
	BaseScore  int     // base points for a correct answer
	Multiplier float64 // speed bonus multiplier
}

// DifficultyConfigs maps difficulty names to their scoring configuration.
var DifficultyConfigs = map[string]DifficultyConfig{
	"easy":   {TimeLimit: 30, BaseScore: 5, Multiplier: 1},
	"medium": {TimeLimit: 60, BaseScore: 10, Multiplier: 2},
	"hard":   {TimeLimit: 120, BaseScore: 20, Multiplier: 3},
}

// GetDifficultyConfig returns the config for a difficulty, defaulting to easy.
func GetDifficultyConfig(difficulty string) DifficultyConfig {
	if cfg, ok := DifficultyConfigs[difficulty]; ok {
		return cfg
	}
	return DifficultyConfigs["easy"]
}

// ─── Combo multiplier ────────────────────────────────────────────

// ComboMultiplier returns the multiplier for the given combo streak.
//
//	0–2  → x1.0
//	3–5  → x1.2
//	6–9  → x1.5
//	10+  → x2.0
func ComboMultiplier(combo int) float64 {
	switch {
	case combo >= 10:
		return 2.0
	case combo >= 6:
		return 1.5
	case combo >= 3:
		return 1.2
	default:
		return 1.0
	}
}

// ─── Score calculation ───────────────────────────────────────────

// CalculateScore computes the score for a single question answer.
// Returns (totalScore, speedBonus) BEFORE combo multiplier.
// If incorrect or timed out, returns (0, 0).
func CalculateScore(difficulty string, timeSpent float64, correct bool) (int, int) {
	if !correct {
		return 0, 0
	}

	cfg := GetDifficultyConfig(difficulty)

	// Timed out — no points
	if timeSpent >= float64(cfg.TimeLimit) {
		return 0, 0
	}

	// Clamp timeSpent to minimum 0
	if timeSpent < 0 {
		timeSpent = 0
	}

	speedBonus := int(math.Floor((float64(cfg.TimeLimit) - timeSpent) * cfg.Multiplier))
	total := cfg.BaseScore + speedBonus

	return total, speedBonus
}

// CalculateScoreWithCombo computes the final score including combo multiplier.
// Returns (finalScore, baseScore, speedBonus, comboMultiplier).
func CalculateScoreWithCombo(difficulty string, timeSpent float64, correct bool, combo int) (int, int, int, float64) {
	baseScore, speedBonus := CalculateScore(difficulty, timeSpent, correct)
	if baseScore == 0 {
		return 0, 0, 0, 1.0
	}
	mult := ComboMultiplier(combo)
	finalScore := int(math.Floor(float64(baseScore) * mult))
	return finalScore, baseScore, speedBonus, mult
}

// MaxScorePerQuestion returns the theoretical max score for a difficulty
// (instant answer, max combo). Used for server-side validation.
func MaxScorePerQuestion(difficulty string) int {
	cfg := GetDifficultyConfig(difficulty)
	base := cfg.BaseScore + int(float64(cfg.TimeLimit)*cfg.Multiplier)
	return int(math.Ceil(float64(base) * 2.0)) // max combo = x2
}
