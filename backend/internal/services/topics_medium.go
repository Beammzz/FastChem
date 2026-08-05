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
