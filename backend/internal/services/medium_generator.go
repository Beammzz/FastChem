package services

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/google/uuid"
	"github.com/takumi/fastchem/internal/models"
)

const avogadro = 6.02e23

// MediumGenerator produces Mole-Concept questions (mass↔moles, moles↔particles)
type MediumGenerator struct{}

// Generate creates one medium question
func (g *MediumGenerator) Generate() models.Question {
	switch rand.Intn(3) {
	case 0:
		return g.massToMoles()
	case 1:
		return g.molesToMass()
	default:
		return g.molesToParticles()
	}
}

// ─── mass → moles ────────────────────────────────────────────────
// mol = mass / MolarMass
func (g *MediumGenerator) massToMoles() models.Question {
	elem := elements[rand.Intn(len(elements))]

	// Pick a nice mass (1–500 g, 1 decimal)
	mass := math.Round((1.0+rand.Float64()*499.0)*10) / 10
	moles := mass / elem.MolarMass
	moles = math.Round(moles*100) / 100

	questionText := fmt.Sprintf(
		"ถ้ามี %s (%s) มวล %.1f g จำนวนโมลคือเท่าใด?",
		elem.NameTH, elem.Symbol, mass,
	)

	choices, correctIndex := generateFloatDistractors(moles, "%.2f mol")

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

// ─── moles → mass ────────────────────────────────────────────────
// mass = mol × MolarMass
func (g *MediumGenerator) molesToMass() models.Question {
	elem := elements[rand.Intn(len(elements))]

	moles := math.Round((0.1+rand.Float64()*9.9)*100) / 100 // 0.10–10.00 mol
	mass := moles * elem.MolarMass
	mass = math.Round(mass*100) / 100

	questionText := fmt.Sprintf(
		"%s (%s) จำนวน %.2f mol มีมวลกี่กรัม?",
		elem.NameTH, elem.Symbol, moles,
	)

	choices, correctIndex := generateFloatDistractors(mass, "%.2f g")

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

// ─── moles → particles ──────────────────────────────────────────
// particles = mol × 6.02×10²³
func (g *MediumGenerator) molesToParticles() models.Question {
	elem := elements[rand.Intn(len(elements))]

	moles := math.Round((0.1+rand.Float64()*4.9)*100) / 100 // 0.10–5.00 mol
	particles := moles * avogadro

	questionText := fmt.Sprintf(
		"%s (%s) จำนวน %.2f mol มีจำนวนอนุภาคเท่าใด? (Nₐ = 6.02×10²³)",
		elem.NameTH, elem.Symbol, moles,
	)

	choices, correctIndex := generateSciDistractors(particles)

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("medium").TimeLimit,
		Category:     "mole_concept",
		Difficulty:   "medium",
	}
}

// generateSciDistractors creates 4 unique choices in scientific notation
func generateSciDistractors(correct float64) ([]string, int) {
	correctStr := formatSci(correct)
	choices := make(map[string]bool)
	choices[correctStr] = true

	for attempt := 0; len(choices) < 4 && attempt < distractorAttempts; attempt++ {
		factor := 0.4 + rand.Float64()*1.2 // 0.4–1.6
		if factor > 0.9 && factor < 1.1 {
			factor = 1.5
		}
		d := correct * factor
		dStr := formatSci(d)
		if dStr != correctStr {
			choices[dStr] = true
		}
	}
	fillSciLadder(choices, correct)

	choiceSlice := make([]string, 0, 4)
	for c := range choices {
		choiceSlice = append(choiceSlice, c)
	}
	rand.Shuffle(len(choiceSlice), func(i, j int) {
		choiceSlice[i], choiceSlice[j] = choiceSlice[j], choiceSlice[i]
	})

	correctIndex := 0
	for i, c := range choiceSlice {
		if c == correctStr {
			correctIndex = i
			break
		}
	}
	return choiceSlice, correctIndex
}

// fillSciLadder tops choices up to 4 by shifting the mantissa in 5% steps.
// The mantissa is always in [1, 10), so a 5% shift always changes the second
// decimal — every k yields a distinct string and the loop always terminates.
func fillSciLadder(choices map[string]bool, correct float64) {
	correctStr := formatSci(correct)
	for k := 1; len(choices) < 4; k++ {
		for _, f := range []float64{1 + 0.05*float64(k), 1 - 0.05*float64(k)} {
			if s := formatSci(correct * f); s != correctStr {
				choices[s] = true
			}
			if len(choices) >= 4 {
				break
			}
		}
	}
}

// formatSci formats a number in "a.bb × 10^n" style
func formatSci(v float64) string {
	if v == 0 {
		return "0"
	}
	exp := math.Floor(math.Log10(math.Abs(v)))
	mantissa := v / math.Pow(10, exp)
	mantissa = math.Round(mantissa*100) / 100
	return fmt.Sprintf("%.2f × 10^%.0f", mantissa, exp)
}
