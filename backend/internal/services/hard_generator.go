package services

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/google/uuid"
	"github.com/takumi/fastchem/internal/models"
)

// HardGenerator produces advanced calculation questions
type HardGenerator struct{}

// solute data for hard questions
type solute struct {
	Name      string // Thai name
	Formula   string
	MolarMass float64 // g/mol
	VantHoff  int     // van't Hoff factor i
}

var solutes = []solute{
	{Name: "โซเดียมคลอไรด์", Formula: "NaCl", MolarMass: 58.44, VantHoff: 2},
	{Name: "กลูโคส", Formula: "C₆H₁₂O₆", MolarMass: 180.16, VantHoff: 1},
	{Name: "แคลเซียมคลอไรด์", Formula: "CaCl₂", MolarMass: 110.98, VantHoff: 3},
	{Name: "โพแทสเซียมคลอไรด์", Formula: "KCl", MolarMass: 74.55, VantHoff: 2},
	{Name: "ยูเรีย", Formula: "CO(NH₂)₂", MolarMass: 60.06, VantHoff: 1},
	{Name: "ซูโครส", Formula: "C₁₂H₂₂O₁₁", MolarMass: 342.30, VantHoff: 1},
	{Name: "โซเดียมไฮดรอกไซด์", Formula: "NaOH", MolarMass: 40.00, VantHoff: 2},
	{Name: "แมกนีเซียมซัลเฟต", Formula: "MgSO₄", MolarMass: 120.37, VantHoff: 2},
}

// Kf for water in °C·kg/mol
const kfWater = 1.86

// Generate creates one hard question
func (g *HardGenerator) Generate() models.Question {
	switch rand.Intn(3) {
	case 0:
		return g.dilution()
	case 1:
		return g.preparingFromSolid()
	default:
		return g.freezingPointDepression()
	}
}

// ─── Dilution: M₁V₁ = M₂V₂ ─────────────────────────────────────
func (g *HardGenerator) dilution() models.Question {
	s := solutes[rand.Intn(len(solutes))]

	// Pick M1, V1, M2 → solve V2
	m1 := math.Round((0.5+rand.Float64()*5.5)*100) / 100        // 0.50–6.00 M
	v1 := math.Round((10+rand.Float64()*490)*10) / 10           // 10.0–500.0 mL
	m2 := math.Round((0.01+rand.Float64()*(m1-0.02))*100) / 100 // smaller than M1

	if m2 <= 0 {
		m2 = 0.01
	}

	v2 := (m1 * v1) / m2
	v2 = math.Round(v2*10) / 10

	questionText := fmt.Sprintf(
		"ต้องเจือจางสารละลาย %s (%s) ความเข้มข้น %.2f M ปริมาตร %.1f mL ให้เหลือ %.2f M ปริมาตรสุดท้ายคือเท่าใด?",
		s.Name, s.Formula, m1, v1, m2,
	)

	choices, correctIndex := generateFloatDistractors(v2, "%.1f mL")

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "dilution",
		Difficulty:   "hard",
	}
}

// ─── Preparing from Solid: mass = M × V(L) × MW ─────────────────
func (g *HardGenerator) preparingFromSolid() models.Question {
	s := solutes[rand.Intn(len(solutes))]

	molarity := math.Round((0.1+rand.Float64()*2.9)*100) / 100 // 0.10–3.00 M
	volumeML := math.Round((50+rand.Float64()*950)*10) / 10    // 50.0–1000.0 mL
	volumeL := volumeML / 1000.0

	mass := molarity * volumeL * s.MolarMass
	mass = math.Round(mass*100) / 100

	questionText := fmt.Sprintf(
		"ต้องใช้ %s (%s, MW=%.2f g/mol) กี่กรัม เพื่อเตรียมสารละลาย %.2f M ปริมาตร %.1f mL?",
		s.Name, s.Formula, s.MolarMass, molarity, volumeML,
	)

	choices, correctIndex := generateFloatDistractors(mass, "%.2f g")

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "preparing_solution",
		Difficulty:   "hard",
	}
}

// ─── Freezing Point Depression: ΔTf = i × Kf × m ────────────────
func (g *HardGenerator) freezingPointDepression() models.Question {
	s := solutes[rand.Intn(len(solutes))]

	// molality m = mol solute / kg solvent
	molality := math.Round((0.1+rand.Float64()*4.9)*100) / 100 // 0.10–5.00 m

	deltaTf := float64(s.VantHoff) * kfWater * molality
	deltaTf = math.Round(deltaTf*100) / 100

	newFP := 0.0 - deltaTf
	newFP = math.Round(newFP*100) / 100

	questionText := fmt.Sprintf(
		"สารละลาย %s (%s, i=%d) ความเข้มข้น %.2f mol/kg ในน้ำ (Kf=1.86 °C·kg/mol) จุดเยือกแข็งใหม่คือเท่าใด?",
		s.Name, s.Formula, s.VantHoff, molality,
	)

	choices, correctIndex := generateFloatDistractors(newFP, "%.2f °C")

	return models.Question{
		ID:           uuid.New().String(),
		Question:     questionText,
		Choices:      choices,
		CorrectIndex: correctIndex,
		TimeLimit:    GetDifficultyConfig("hard").TimeLimit,
		Category:     "freezing_point",
		Difficulty:   "hard",
	}
}
