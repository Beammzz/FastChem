package services

import (
	"fmt"
	"math"

	"github.com/takumi/fastchem/internal/models"
)

// ─── Atomic structure ────────────────────────────────────────────

type atomicStructureTopic struct{}

func (atomicStructureTopic) Category() string   { return "atomic_structure" }
func (atomicStructureTopic) Difficulty() string { return "easy" }

// Bounds for particle counts: a distractor of -1 protons is not a mistake
// anyone makes, and 0 protons is not an atom.
const (
	minParticleCount = 0
	minProtonCount   = 1
	maxNucleonCount  = 250
)

func (t atomicStructureTopic) Generate(s *Source) Problem {
	elem := elements[s.Pick("elements", len(elements))]
	z := elem.AtomicNumber

	if s.Intn(4) == 3 && len(elem.Isotopes) > 0 {
		return t.neutronCount(s, elem)
	}

	var text string
	switch s.Intn(3) {
	case 0:
		text = fmt.Sprintf("ธาตุ %s (%s) มีเลขอะตอมเท่าใด?", elem.NameTH, elem.Symbol)
	case 1:
		text = fmt.Sprintf("ธาตุ %s (%s) มีจำนวนโปรตอนเท่าใด?", elem.NameTH, elem.Symbol)
	default:
		text = fmt.Sprintf("อะตอมที่เป็นกลางของ %s (%s) มีจำนวนอิเล็กตรอนเท่าใด?", elem.NameTH, elem.Symbol)
	}

	candidates := intCandidates(minProtonCount, maxNucleonCount, []int{
		int(math.Round(elem.MolarMass)), // read the atomic mass as the atomic number
		z + 1,                           // miscounted one place along the periodic table
		z - 1,
	}, formatInt)

	choices, idx := assembleChoices(s, formatInt(z), candidates,
		intTopUp(z, minProtonCount, maxNucleonCount, formatInt))

	return Problem{
		Question:     text,
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t atomicStructureTopic) neutronCount(s *Source, elem models.Element) Problem {
	z := elem.AtomicNumber
	a := elem.Isotopes[s.Intn(len(elem.Isotopes))]
	neutrons := a - z

	candidates := intCandidates(minParticleCount, maxNucleonCount, []int{
		a,     // forgot to subtract the atomic number
		z,     // answered with the proton count
		a + z, // added where the rule subtracts
	}, formatInt)

	choices, idx := assembleChoices(s, formatInt(neutrons), candidates,
		intTopUp(neutrons, minParticleCount, maxNucleonCount, formatInt))

	return Problem{
		Question: fmt.Sprintf("ไอโซโทปของ %s (%s) ที่มีเลขมวล %d มีจำนวนนิวตรอนเท่าใด?",
			elem.NameTH, elem.Symbol, a),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Oxidation number ────────────────────────────────────────────

type oxidationNumberTopic struct{}

func (oxidationNumberTopic) Category() string   { return "oxidation_number" }
func (oxidationNumberTopic) Difficulty() string { return "easy" }

// Oxidation numbers taught at this level span -4 to +7.
const (
	minOxidation = -4
	maxOxidation = 7
)

func (t oxidationNumberTopic) Generate(s *Source) Problem {
	compound := compounds[s.Pick("compounds", len(compounds))]
	targetIdx := s.Intn(len(compound.Elements))
	target := compound.Elements[targetIdx]
	ox := target.OxidationNumber

	// Reversing the sign is far and away the most common error, so it leads.
	values := []int{-ox}
	for i, e := range compound.Elements {
		if i != targetIdx {
			values = append(values, e.OxidationNumber) // read the wrong element
		}
	}
	values = append(values,
		ox*target.Count, // multiplied by the subscript
		0,               // assumed the element is in its elemental state
	)

	candidates := intCandidates(minOxidation, maxOxidation, values, formatOxidation)
	choices, idx := assembleChoices(s, formatOxidation(ox), candidates,
		intTopUp(ox, minOxidation, maxOxidation, formatOxidation))

	return Problem{
		Question: fmt.Sprintf("เลขออกซิเดชันของ %s ใน %s คือเท่าใด?",
			target.Symbol, compound.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── State of matter ─────────────────────────────────────────────

type stateOfMatterTopic struct{}

func (stateOfMatterTopic) Category() string   { return "state_of_matter" }
func (stateOfMatterTopic) Difficulty() string { return "easy" }

// stateLabels maps internal state codes to Thai display labels.
var stateLabels = map[string]string{
	"s":  "ของแข็ง (s)",
	"l":  "ของเหลว (l)",
	"g":  "แก๊ส (g)",
	"aq": "สารละลายในน้ำ (aq)",
}

// stateOrder fixes the iteration order of the states. Ranging over stateLabels
// would reintroduce map-order nondeterminism into the seeded path.
var stateOrder = []string{"s", "l", "g", "aq"}

func (t stateOfMatterTopic) Generate(s *Source) Problem {
	compound := stateCompounds[s.Pick("stateCompounds", len(stateCompounds))]

	// The three states this compound is not in are the only sensible
	// distractors, and there are exactly three of them.
	candidates := make([]string, 0, 3)
	for _, code := range stateOrder {
		if code != compound.State {
			candidates = append(candidates, stateLabels[code])
		}
	}

	choices, idx := assembleChoices(s, stateLabels[compound.State], candidates,
		func(int) string { return "" })

	return Problem{
		Question: fmt.Sprintf("%s (%s) มีสถานะที่ภาวะมาตรฐาน (25°C, 1 atm) เป็นอย่างไร?",
			compound.NameTH, compound.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}
