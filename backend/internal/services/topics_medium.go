package services

import (
	"fmt"
	"math"

	"github.com/takumi/fastchem/internal/models"
)

const avogadro = 6.02e23

// minMoles keeps answers at or above 0.10 mol. Below that a two-decimal
// answer has almost no room for distinct, plausible distractors — the old
// generator hung outright on 0.02 and 0.03 mol.
const minMoles = 0.10

type moleConceptTopic struct{}

func (moleConceptTopic) Category() string   { return "mole_concept" }
func (moleConceptTopic) Difficulty() string { return "medium" }

func (t moleConceptTopic) Generate(s *Source) Problem {
	elem := elements[s.Pick("elements", len(elements))]
	switch s.Intn(3) {
	case 0:
		return t.massToMoles(s, elem)
	case 1:
		return t.molesToMass(s, elem)
	default:
		return t.molesToParticles(s, elem)
	}
}

// mol = mass / molar mass
func (t moleConceptTopic) massToMoles(s *Source, elem models.Element) Problem {
	const format = "%.2f mol"

	minMass := math.Max(1.0, minMoles*elem.MolarMass)
	mass := snapTo(minMass+s.Float64()*(500.0-minMass), 0.1)
	moles := snapTo(mass/elem.MolarMass, 0.01)

	candidates := floatCandidates(moles, format, []float64{
		mass / float64(elem.AtomicNumber), // used the atomic number as the molar mass
		elem.MolarMass / mass,             // inverted the ratio
		mass * elem.MolarMass,             // multiplied where the rule divides
		mass / elem.MolarMass * 1000,      // read the mass as kilograms
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, moles), candidates,
		floatTopUp(moles, format))

	return Problem{
		Question: fmt.Sprintf("ถ้ามี %s (%s) มวล %.1f g จำนวนโมลคือเท่าใด?",
			elem.NameTH, elem.Symbol, mass),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// mass = mol × molar mass
func (t moleConceptTopic) molesToMass(s *Source, elem models.Element) Problem {
	const format = "%.2f g"

	moles := snapTo(0.1+s.Float64()*9.9, 0.01)
	mass := snapTo(moles*elem.MolarMass, 0.01)

	candidates := floatCandidates(mass, format, []float64{
		moles * float64(elem.AtomicNumber), // used the atomic number as the molar mass
		moles / elem.MolarMass,             // divided where the rule multiplies
		elem.MolarMass / moles,             // inverted the ratio
		moles * elem.MolarMass / 1000,      // answered in kilograms
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, mass), candidates,
		floatTopUp(mass, format))

	return Problem{
		Question: fmt.Sprintf("%s (%s) จำนวน %.2f mol มีมวลกี่กรัม?",
			elem.NameTH, elem.Symbol, moles),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// particles = mol × Avogadro's number
func (t moleConceptTopic) molesToParticles(s *Source, elem models.Element) Problem {
	moles := snapTo(0.1+s.Float64()*4.9, 0.01)
	particles := moles * avogadro

	candidates := sciCandidates(particles, []float64{
		particles * 10, // slipped the exponent
		particles / 10,
		moles * avogadro / elem.MolarMass, // divided by the molar mass as well
	})

	choices, idx := assembleChoices(s, formatSci(particles), candidates,
		sciTopUp(particles))

	return Problem{
		Question: fmt.Sprintf("%s (%s) จำนวน %.2f mol มีจำนวนอนุภาคเท่าใด? (Nₐ = 6.02×10²³)",
			elem.NameTH, elem.Symbol, moles),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Molar mass of a compound (บทที่ 4) ───────────────────────────

type molarMassTopic struct{}

func (molarMassTopic) Category() string   { return "molar_mass" }
func (molarMassTopic) Difficulty() string { return "medium" }

func (t molarMassTopic) Generate(s *Source) Problem {
	const format = "%.2f g/mol"
	c := formulaCompounds[s.Pick("formulaCompounds", len(formulaCompounds))]
	mw := snapTo(c.MolarMass(), 0.01)

	// Subscript mistakes, spelled out as the arithmetic they produce.
	noSubscripts, addedSubscripts := 0.0, 0.0
	for _, p := range c.Parts {
		noSubscripts += atomicMasses[p.Symbol]
		addedSubscripts += atomicMasses[p.Symbol] + float64(p.Count)
	}

	candidates := floatCandidates(mw, format, []float64{
		noSubscripts,                 // ignored the subscripts entirely
		addedSubscripts,              // added the subscript instead of multiplying by it
		mw / float64(c.TotalAtoms()), // averaged the mass over the atoms
		mw * float64(len(c.Parts)),   // multiplied through by the number of elements
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, mw), candidates,
		floatTopUp(mw, format))

	return Problem{
		Question:     fmt.Sprintf("มวลต่อโมลของ %s (%s) เท่ากับเท่าใด?", c.NameTH, c.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Percent composition by mass (บทที่ 4) ────────────────────────

type percentCompositionTopic struct{}

func (percentCompositionTopic) Category() string   { return "percent_composition" }
func (percentCompositionTopic) Difficulty() string { return "medium" }

func (t percentCompositionTopic) Generate(s *Source) Problem {
	const format = "%.2f %%"
	c := formulaCompounds[s.Pick("formulaCompounds", len(formulaCompounds))]
	part := c.Parts[s.Intn(len(c.Parts))]

	mw := c.MolarMass()
	percent := snapTo(part.Mass()/mw*100, 0.01)

	candidates := percentCandidates(percent, format, []float64{
		atomicMasses[part.Symbol] / mw * 100,                // dropped the subscript from the numerator
		float64(part.Count) / float64(c.TotalAtoms()) * 100, // counted atoms instead of weighing them
		100 - percent,                          // answered with everything except the element asked for
		part.Mass() / (mw - part.Mass()) * 100, // divided by the rest of the compound, not the whole
		part.Mass() / mw,                       // never multiplied by 100
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, percent), candidates,
		floatTopUp(percent, format))

	return Problem{
		Question: fmt.Sprintf("ร้อยละโดยมวลของธาตุ %s ใน %s (%s, M = %.2f g/mol) เท่ากับเท่าใด?",
			part.Symbol, c.NameTH, c.Formula, snapTo(mw, 0.01)),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Concentration units (บทที่ 5) ────────────────────────────────

type concentrationTopic struct{}

func (concentrationTopic) Category() string   { return "concentration" }
func (concentrationTopic) Difficulty() string { return "medium" }

func (t concentrationTopic) Generate(s *Source) Problem {
	sol := solutes[s.Pick("solutes", len(solutes))]
	switch s.Intn(4) {
	case 0:
		return t.molarity(s, sol)
	case 1:
		return t.percentByMass(s, sol)
	case 2:
		return t.partsPerMillion(s, sol)
	default:
		return t.molality(s, sol)
	}
}

// M = (mass / MW) / volume in litres
func (t concentrationTopic) molarity(s *Source, sol solute) Problem {
	const format = "%.2f M"

	// Keep the answer at or above 0.10 M so two decimals still leave room for
	// distinct distractors, the same reason minMoles exists.
	volumeML := snapTo(100+s.Float64()*900, 1)
	minMass := minMoles * sol.MolarMass * volumeML / 1000
	mass := snapTo(minMass+s.Float64()*(minMass*9), 0.01)
	molarity := snapTo(mass/sol.MolarMass/(volumeML/1000), 0.01)

	candidates := floatCandidates(molarity, format, []float64{
		mass / sol.MolarMass / volumeML, // never converted mL to L
		mass / volumeML * 1000,          // skipped the molar mass
		mass * sol.MolarMass / volumeML * 1000,
		volumeML / 1000 / (mass / sol.MolarMass), // inverted the ratio
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, molarity), candidates,
		floatTopUp(molarity, format))

	return Problem{
		Question: fmt.Sprintf(
			"ละลาย %s (%s, MW=%.2f g/mol) มวล %.2f g ในน้ำจนได้สารละลายปริมาตร %.0f mL ความเข้มข้นเป็นโมลาร์เท่าใด?",
			sol.Name, sol.Formula, sol.MolarMass, mass, volumeML),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// %w/w = mass of solute / mass of solution × 100
func (t concentrationTopic) percentByMass(s *Source, sol solute) Problem {
	const format = "%.2f %%"

	soluteMass := snapTo(5+s.Float64()*95, 0.1)
	solventMass := snapTo(100+s.Float64()*400, 0.1)
	solutionMass := soluteMass + solventMass
	percent := snapTo(soluteMass/solutionMass*100, 0.01)

	candidates := percentCandidates(percent, format, []float64{
		soluteMass / solventMass * 100,   // divided by the solvent instead of the solution
		solventMass / solutionMass * 100, // answered with the solvent's share
		soluteMass / solutionMass,        // never multiplied by 100
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, percent), candidates,
		floatTopUp(percent, format))

	return Problem{
		Question: fmt.Sprintf(
			"ละลาย %s (%s) มวล %.1f g ในน้ำ %.1f g ความเข้มข้นเป็นร้อยละโดยมวลเท่าใด?",
			sol.Name, sol.Formula, soluteMass, solventMass),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ppm = mass of solute / mass of solution × 10⁶
func (t concentrationTopic) partsPerMillion(s *Source, sol solute) Problem {
	const format = "%.0f ppm"

	soluteMG := snapTo(20+s.Float64()*1980, 1)   // 20–2000 mg
	solutionG := snapTo(200+s.Float64()*1800, 1) // 200–2000 g
	ppm := snapTo(soluteMG/1000/solutionG*1e6, 1)

	candidates := floatCandidates(ppm, format, []float64{
		soluteMG / solutionG * 1e6,        // left the solute in milligrams
		soluteMG / 1000 / solutionG * 1e3, // used 10³ (ppt) instead of 10⁶
		soluteMG / 1000 / solutionG * 100, // computed a percentage
		solutionG / (soluteMG / 1000) * 1e6,
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, ppm), candidates,
		floatTopUp(ppm, format))

	return Problem{
		Question: fmt.Sprintf(
			"สารละลาย %s (%s) มวล %.0f g มีตัวละลายอยู่ %.0f mg ความเข้มข้นเป็นกี่ ppm?",
			sol.Name, sol.Formula, solutionG, soluteMG),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// molality = moles of solute / kilograms of solvent
func (t concentrationTopic) molality(s *Source, sol solute) Problem {
	const format = "%.2f mol/kg"

	solventKG := snapTo(0.2+s.Float64()*1.8, 0.01)
	minMass := minMoles * sol.MolarMass * solventKG
	mass := snapTo(minMass+s.Float64()*(minMass*9), 0.01)
	molality := snapTo(mass/sol.MolarMass/solventKG, 0.01)

	candidates := floatCandidates(molality, format, []float64{
		mass / sol.MolarMass / (solventKG * 1000), // used grams of solvent
		mass / solventKG, // skipped the molar mass
		mass / sol.MolarMass / (solventKG + mass/1000), // divided by the solution, not the solvent
		solventKG / (mass / sol.MolarMass),
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, molality), candidates,
		floatTopUp(molality, format))

	return Problem{
		Question: fmt.Sprintf(
			"ละลาย %s (%s, MW=%.2f g/mol) มวล %.2f g ในน้ำ %.2f kg ความเข้มข้นเป็นโมแลลเท่าใด?",
			sol.Name, sol.Formula, sol.MolarMass, mass, solventKG),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Stoichiometry (บทที่ 6) ──────────────────────────────────────

type stoichiometryTopic struct{}

func (stoichiometryTopic) Category() string   { return "stoichiometry" }
func (stoichiometryTopic) Difficulty() string { return "medium" }

// roledSpecies is a species together with the side of the equation it sits on,
// which is what decides whether the Thai says "ต้องใช้" or "จะเกิด".
type roledSpecies struct {
	reactionSpecies
	IsReactant bool
}

// speciesPair draws two different species from a reaction, so a question never
// asks how much A is produced from A.
func speciesPair(s *Source, r reaction) (from, to roledSpecies) {
	all := make([]roledSpecies, 0, len(r.Reactants)+len(r.Products))
	for _, sp := range r.Reactants {
		all = append(all, roledSpecies{sp, true})
	}
	for _, sp := range r.Products {
		all = append(all, roledSpecies{sp, false})
	}

	i := s.Intn(len(all))
	j := s.Intn(len(all) - 1)
	if j >= i {
		j++
	}
	return all[i], all[j]
}

// givenClause and askedClause keep the question grammatical in both
// directions: a reactant is consumed, a product is formed, and asking "how
// much H₂ is formed from ZnCl₂" with one verb for both reads as nonsense.
func givenClause(from roledSpecies, amount string) string {
	if from.IsReactant {
		return fmt.Sprintf("ถ้าใช้ %s %s", from.Formula, amount)
	}
	return fmt.Sprintf("ถ้าเกิด %s %s", from.Formula, amount)
}

func askedClause(to roledSpecies, tail string) string {
	if to.IsReactant {
		return fmt.Sprintf("ต้องใช้ %s %s", to.Formula, tail)
	}
	return fmt.Sprintf("จะเกิด %s %s", to.Formula, tail)
}

func (t stoichiometryTopic) Generate(s *Source) Problem {
	r := reactions[s.Pick("reactions", len(reactions))]
	from, to := speciesPair(s, r)

	if s.Intn(2) == 0 {
		return t.moleToMole(s, r, from, to)
	}
	return t.moleToMass(s, r, from, to)
}

func (t stoichiometryTopic) moleToMole(s *Source, r reaction, from, to roledSpecies) Problem {
	const format = "%.2f mol"

	moles := snapTo(0.2+s.Float64()*9.8, 0.01)
	ratio := float64(to.Coef) / float64(from.Coef)
	answer := snapTo(moles*ratio, 0.01)

	candidates := floatCandidates(answer, format, []float64{
		moles / ratio, // inverted the mole ratio
		moles,         // assumed 1:1 regardless of the coefficients
		moles * float64(to.Coef) * float64(from.Coef), // multiplied the coefficients instead of dividing
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, answer), candidates,
		floatTopUp(answer, format))

	return Problem{
		Question: fmt.Sprintf("จากสมการ %s %s %s",
			r.Equation,
			givenClause(from, fmt.Sprintf("จำนวน %.2f mol", moles)),
			askedClause(to, "กี่โมล?")),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t stoichiometryTopic) moleToMass(s *Source, r reaction, from, to roledSpecies) Problem {
	const format = "%.2f g"

	moles := snapTo(0.2+s.Float64()*4.8, 0.01)
	ratio := float64(to.Coef) / float64(from.Coef)
	answer := snapTo(moles*ratio*to.MolarMass, 0.01)

	candidates := floatCandidates(answer, format, []float64{
		moles * to.MolarMass,           // ignored the mole ratio
		moles / ratio * to.MolarMass,   // inverted the mole ratio
		moles * ratio,                  // stopped at moles and never weighed them
		moles * ratio * from.MolarMass, // weighed with the wrong substance's molar mass
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, answer), candidates,
		floatTopUp(answer, format))

	return Problem{
		Question: fmt.Sprintf("จากสมการ %s %s %s",
			r.Equation,
			givenClause(from, fmt.Sprintf("จำนวน %.2f mol", moles)),
			askedClause(to, fmt.Sprintf("(M = %.2f g/mol) กี่กรัม?", to.MolarMass))),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Gas laws (บทที่ 7) ───────────────────────────────────────────

type gasLawTopic struct{}

func (gasLawTopic) Category() string   { return "gas_law" }
func (gasLawTopic) Difficulty() string { return "medium" }

// molarVolumeSTP is the volume one mole of an ideal gas occupies at STP,
// in litres. 24.0 L/mol — the value at 25 °C — is the classic wrong constant,
// so it earns a distractor.
const (
	molarVolumeSTP = 22.4
	molarVolumeRTP = 24.0
	kelvinOffset   = 273.15
)

func (t gasLawTopic) Generate(s *Source) Problem {
	switch s.Intn(4) {
	case 0:
		return t.boyle(s)
	case 1:
		return t.charles(s)
	case 2:
		return t.combined(s)
	default:
		return t.molarVolume(s)
	}
}

// P₁V₁ = P₂V₂ at constant temperature
func (t gasLawTopic) boyle(s *Source) Problem {
	const format = "%.2f L"

	p1 := snapTo(1+s.Float64()*4, 0.1)
	v1 := snapTo(1+s.Float64()*9, 0.1)
	p2 := snapTo(1+s.Float64()*4, 0.1)
	v2 := snapTo(p1*v1/p2, 0.01)

	candidates := floatCandidates(v2, format, []float64{
		p2 * v1 / p1,   // inverted the pressure ratio
		p1 * v1 * p2,   // multiplied where the law divides
		v1 + (p1 - p2), // treated the relationship as additive
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, v2), candidates,
		floatTopUp(v2, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สปริมาตร %.1f L ที่ความดัน %.1f atm ถูกเปลี่ยนความดันเป็น %.1f atm ที่อุณหภูมิคงที่ ปริมาตรใหม่เป็นเท่าใด?",
			v1, p1, p2),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// V₁/T₁ = V₂/T₂ at constant pressure, with T in kelvin
func (t gasLawTopic) charles(s *Source) Problem {
	const format = "%.2f L"

	v1 := snapTo(1+s.Float64()*9, 0.1)
	c1 := snapTo(20+s.Float64()*60, 1)   // 20–80 °C
	c2 := snapTo(100+s.Float64()*200, 1) // 100–300 °C
	t1 := c1 + kelvinOffset
	t2 := c2 + kelvinOffset
	v2 := snapTo(v1*t2/t1, 0.01)

	candidates := floatCandidates(v2, format, []float64{
		v1 * c2 / c1, // stayed in celsius instead of converting to kelvin
		v1 * t1 / t2, // inverted the temperature ratio
		v1 + (t2 - t1),
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, v2), candidates,
		floatTopUp(v2, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สปริมาตร %.1f L ที่ %.0f °C ถูกทำให้ร้อนขึ้นเป็น %.0f °C ที่ความดันคงที่ ปริมาตรใหม่เป็นเท่าใด?",
			v1, c1, c2),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// P₁V₁/T₁ = P₂V₂/T₂
func (t gasLawTopic) combined(s *Source) Problem {
	const format = "%.2f L"

	p1 := snapTo(1+s.Float64()*3, 0.1)
	v1 := snapTo(1+s.Float64()*9, 0.1)
	c1 := snapTo(20+s.Float64()*60, 1)
	p2 := snapTo(1+s.Float64()*3, 0.1)
	c2 := snapTo(100+s.Float64()*200, 1)
	t1 := c1 + kelvinOffset
	t2 := c2 + kelvinOffset
	v2 := snapTo(p1*v1*t2/(p2*t1), 0.01)

	candidates := floatCandidates(v2, format, []float64{
		p1 * v1 * c2 / (p2 * c1), // stayed in celsius
		p2 * v1 * t2 / (p1 * t1), // inverted the pressure ratio
		p1 * v1 * t1 / (p2 * t2), // inverted the temperature ratio
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, v2), candidates,
		floatTopUp(v2, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สปริมาตร %.1f L ที่ %.1f atm, %.0f °C เปลี่ยนไปอยู่ที่ %.1f atm, %.0f °C ปริมาตรใหม่เป็นเท่าใด?",
			v1, p1, c1, p2, c2),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// V = n × 22.4 L at STP
func (t gasLawTopic) molarVolume(s *Source) Problem {
	const format = "%.2f L"

	moles := snapTo(0.2+s.Float64()*4.8, 0.01)
	volume := snapTo(moles*molarVolumeSTP, 0.01)

	candidates := floatCandidates(volume, format, []float64{
		moles * molarVolumeRTP, // used 24.0 L/mol, the value at 25 °C, not at STP
		moles / molarVolumeSTP, // divided where the rule multiplies
		molarVolumeSTP / moles,
	})

	choices, idx := assembleChoices(s, fmt.Sprintf(format, volume), candidates,
		floatTopUp(volume, format))

	return Problem{
		Question: fmt.Sprintf(
			"แก๊สอุดมคติจำนวน %.2f mol ที่ STP มีปริมาตรเท่าใด? (1 mol ของแก๊สที่ STP = 22.4 L)",
			moles),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}
