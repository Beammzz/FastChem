package services

import "fmt"

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

// Cryoscopic and ebullioscopic constants for water, in °C·kg/mol.
// Reaching for Kb where the problem calls for Kf is a standard slip, so it
// earns a distractor.
const (
	kfWater = 1.86
	kbWater = 0.512
)

// ─── Dilution: M₁V₁ = M₂V₂ ──────────────────────────────────────

type dilutionTopic struct{}

func (dilutionTopic) Category() string   { return "dilution" }
func (dilutionTopic) Difficulty() string { return "hard" }

func (t dilutionTopic) Generate(s *Source) Problem {
	const format = "%.1f mL"
	sol := solutes[s.Pick("solutes", len(solutes))]

	m1 := snapTo(0.5+s.Float64()*5.5, 0.01) // 0.50–6.00 M
	v1 := snapTo(10+s.Float64()*490, 0.1)   // 10.0–500.0 mL
	m2 := snapTo(0.01+s.Float64()*(m1-0.02), 0.01)
	if m2 <= 0 {
		m2 = 0.01
	}
	v2 := snapTo(m1*v1/m2, 0.1)

	candidates := floatCandidates(v2, format, []float64{
		v2 - v1,      // gave the volume of water to add, not the final volume
		m2 * v1 / m1, // inverted the concentration ratio
		m1 * v1 * m2, // multiplied through where the rule divides
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, v2), candidates,
		floatTopUp(v2, format))

	return Problem{
		Question: fmt.Sprintf(
			"ต้องเจือจางสารละลาย %s (%s) ความเข้มข้น %.2f M ปริมาตร %.1f mL ให้เหลือ %.2f M ปริมาตรสุดท้ายคือเท่าใด?",
			sol.Name, sol.Formula, m1, v1, m2),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Preparing from solid: mass = M × V(L) × MW ─────────────────

type preparingSolutionTopic struct{}

func (preparingSolutionTopic) Category() string   { return "preparing_solution" }
func (preparingSolutionTopic) Difficulty() string { return "hard" }

func (t preparingSolutionTopic) Generate(s *Source) Problem {
	const format = "%.2f g"
	sol := solutes[s.Pick("solutes", len(solutes))]

	molarity := snapTo(0.1+s.Float64()*2.9, 0.01) // 0.10–3.00 M
	volumeML := snapTo(50+s.Float64()*950, 0.1)   // 50.0–1000.0 mL
	volumeL := volumeML / 1000.0
	mass := snapTo(molarity*volumeL*sol.MolarMass, 0.01)

	candidates := floatCandidates(mass, format, []float64{
		volumeL * sol.MolarMass,             // dropped the concentration
		molarity * sol.MolarMass,            // dropped the volume
		molarity * volumeL / sol.MolarMass,  // divided by the molar mass
		molarity * volumeML * sol.MolarMass, // never converted mL to L
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, mass), candidates,
		floatTopUp(mass, format))

	return Problem{
		Question: fmt.Sprintf(
			"ต้องใช้ %s (%s, MW=%.2f g/mol) กี่กรัม เพื่อเตรียมสารละลาย %.2f M ปริมาตร %.1f mL?",
			sol.Name, sol.Formula, sol.MolarMass, molarity, volumeML),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Freezing point depression: ΔTf = i × Kf × m ────────────────

type freezingPointTopic struct{}

func (freezingPointTopic) Category() string   { return "freezing_point" }
func (freezingPointTopic) Difficulty() string { return "hard" }

func (t freezingPointTopic) Generate(s *Source) Problem {
	const format = "%.2f °C"
	sol := solutes[s.Pick("solutes", len(solutes))]

	molality := snapTo(0.1+s.Float64()*4.9, 0.01) // 0.10–5.00 mol/kg
	i := float64(sol.VantHoff)
	deltaTf := snapTo(i*kfWater*molality, 0.01)
	freezingPoint := snapTo(-deltaTf, 0.01)

	candidates := floatCandidates(freezingPoint, format, []float64{
		deltaTf,                 // reported the depression, not the new freezing point
		-kfWater * molality,     // omitted the van 't Hoff factor
		-i * kbWater * molality, // used Kb for water instead of Kf
		-i * kfWater / molality, // divided by molality where the rule multiplies
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, freezingPoint), candidates,
		floatTopUp(freezingPoint, format))

	return Problem{
		Question: fmt.Sprintf(
			"สารละลาย %s (%s, i=%d) ความเข้มข้น %.2f mol/kg ในน้ำ (Kf=1.86 °C·kg/mol) จุดเยือกแข็งใหม่คือเท่าใด?",
			sol.Name, sol.Formula, sol.VantHoff, molality),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}
