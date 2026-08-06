package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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

	choices, idx := assembleChoices(s, stateLabels[compound.State], candidates, noTopUp)

	return Problem{
		Question: fmt.Sprintf("%s (%s) มีสถานะที่ภาวะมาตรฐาน (25°C, 1 atm) เป็นอย่างไร?",
			compound.NameTH, compound.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Electron configuration (บทที่ 2) ────────────────────────────

type electronConfigTopic struct{}

func (electronConfigTopic) Category() string   { return "electron_config" }
func (electronConfigTopic) Difficulty() string { return "easy" }

// Valence electrons run 1–8; periods run 1–7 on the printed table, which is
// the range a student would consider even though the pool stops at calcium.
const (
	minValence = 1
	maxValence = 8
	minPeriod  = 1
	maxPeriod  = 7
)

// shellString renders a shell configuration the way the book writes it.
func shellString(shells []int) string {
	parts := make([]string, len(shells))
	for i, n := range shells {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ", ")
}

// greedyShells fills 2, 8, 18, 32 without the outer-shell rule that caps the
// third level at 8 until the fourth has opened. It matches the correct answer
// for most elements and diverges exactly where students go wrong — potassium
// and calcium, which it puts at "2, 8, 9" and "2, 8, 10".
func greedyShells(z int) []int {
	capacities := []int{2, 8, 18, 32}
	out := make([]int, 0, len(capacities))
	for _, capacity := range capacities {
		if z <= 0 {
			break
		}
		n := z
		if n > capacity {
			n = capacity
		}
		out = append(out, n)
		z -= n
	}
	return out
}

func reversedShells(shells []int) []int {
	out := make([]int, len(shells))
	for i, n := range shells {
		out[len(shells)-1-i] = n
	}
	return out
}

func (t electronConfigTopic) Generate(s *Source) Problem {
	idx := s.Pick("shellElements", len(shellElements))
	elem := shellElements[idx]

	switch s.Intn(4) {
	case 0:
		return t.shellConfig(s, idx, elem)
	case 1:
		return t.valenceElectrons(s, elem)
	case 2:
		return t.group(s, elem)
	default:
		return t.period(s, elem)
	}
}

func (t electronConfigTopic) shellConfig(s *Source, idx int, elem shellElement) Problem {
	correct := shellString(elem.Shells)

	candidates := make([]string, 0, 4)
	if idx > 0 {
		candidates = append(candidates, shellString(shellElements[idx-1].Shells)) // off by one element
	}
	if idx+1 < len(shellElements) {
		candidates = append(candidates, shellString(shellElements[idx+1].Shells))
	}
	candidates = append(candidates,
		shellString(greedyShells(elem.Number)),   // filled 2-8-18 without capping the third shell
		shellString(reversedShells(elem.Shells)), // wrote the levels outermost first
	)

	// Any other element's configuration is a real configuration, so the ladder
	// never has to invent one.
	topUp := func(k int) string {
		if k-1 >= len(shellElements) {
			return ""
		}
		return shellString(shellElements[k-1].Shells)
	}

	choices, correctIdx := assembleChoices(s, correct, candidates, topUp)

	return Problem{
		Question: fmt.Sprintf("การจัดเรียงอิเล็กตรอนในระดับพลังงานหลักของ %s (%s, Z=%d) เป็นอย่างไร?",
			elem.NameTH, elem.Symbol, elem.Number),
		Choices:      choices,
		CorrectIndex: correctIdx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t electronConfigTopic) valenceElectrons(s *Source, elem shellElement) Problem {
	valence := elem.Shells[len(elem.Shells)-1]

	candidates := intCandidates(minValence, maxValence, []int{
		elem.Number,      // answered with the atomic number
		8 - valence,      // counted the electrons still needed for an octet
		len(elem.Shells), // answered with the number of shells
	}, formatInt)

	choices, idx := assembleChoices(s, formatInt(valence), candidates,
		intTopUp(valence, minValence, maxValence, formatInt))

	return Problem{
		Question: fmt.Sprintf("ธาตุ %s (%s, Z=%d) มีเวเลนซ์อิเล็กตรอนกี่ตัว?",
			elem.NameTH, elem.Symbol, elem.Number),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t electronConfigTopic) group(s *Source, elem shellElement) Problem {
	candidates := textCandidates(s, elem.Group, representativeGroups)
	choices, idx := assembleChoices(s, elem.Group, candidates, noTopUp)

	return Problem{
		Question: fmt.Sprintf("ธาตุ %s (%s) ที่มีการจัดเรียงอิเล็กตรอน %s อยู่หมู่ใดในตารางธาตุ?",
			elem.NameTH, elem.Symbol, shellString(elem.Shells)),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

func (t electronConfigTopic) period(s *Source, elem shellElement) Problem {
	valence := elem.Shells[len(elem.Shells)-1]

	candidates := intCandidates(minPeriod, maxPeriod, []int{
		valence,         // read the group off the outer shell and called it the period
		elem.Period + 1, // counted the shells starting from zero
		len(elem.Shells) + 1,
	}, formatInt)

	choices, idx := assembleChoices(s, formatInt(elem.Period), candidates,
		intTopUp(elem.Period, minPeriod, maxPeriod, formatInt))

	return Problem{
		Question: fmt.Sprintf("ธาตุ %s (%s) ที่มีการจัดเรียงอิเล็กตรอน %s อยู่คาบใดในตารางธาตุ?",
			elem.NameTH, elem.Symbol, shellString(elem.Shells)),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Bond type (บทที่ 3) ──────────────────────────────────────────

type bondTypeTopic struct{}

func (bondTypeTopic) Category() string   { return "bond_type" }
func (bondTypeTopic) Difficulty() string { return "easy" }

func (t bondTypeTopic) Generate(s *Source) Problem {
	sub := bondSubstances[s.Pick("bondSubstances", len(bondSubstances))]

	// The three classes the substance is not in are the only sensible wrong
	// answers, and there are exactly three of them.
	candidates := textCandidates(s, sub.Class, bondClasses)
	choices, idx := assembleChoices(s, sub.Class, candidates, noTopUp)

	return Problem{
		Question: fmt.Sprintf("%s (%s) ยึดเหนี่ยวกันด้วยพันธะชนิดใด?",
			sub.NameTH, sub.Formula),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Molecular shape, VSEPR (บทที่ 3) ─────────────────────────────

type molecularShapeTopic struct{}

func (molecularShapeTopic) Category() string   { return "molecular_shape" }
func (molecularShapeTopic) Difficulty() string { return "easy" }

func (t molecularShapeTopic) Generate(s *Source) Problem {
	mol := shapeMolecules[s.Pick("shapeMolecules", len(shapeMolecules))]

	candidates := textCandidates(s, mol.Shape, molecularShapes)
	choices, idx := assembleChoices(s, mol.Shape, candidates, noTopUp)

	// Quoting the electron-pair counts keeps this derivable from the VSEPR
	// rule rather than recallable only by having seen the molecule before.
	return Problem{
		Question: fmt.Sprintf(
			"โมเลกุล %s (%s) มีอะตอมกลางคือ %s ล้อมรอบด้วยพันธะ %d กลุ่ม และอิเล็กตรอนคู่โดดเดี่ยว %d คู่ รูปร่างโมเลกุลเป็นแบบใด?",
			mol.NameTH, mol.Formula, mol.Central, mol.Bonded, mol.Lone),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Functional groups (บทที่ 12) ─────────────────────────────────

type functionalGroupTopic struct{}

func (functionalGroupTopic) Category() string   { return "functional_group" }
func (functionalGroupTopic) Difficulty() string { return "easy" }

func (t functionalGroupTopic) Generate(s *Source) Problem {
	c := organicCompounds[s.Pick("organicCompounds", len(organicCompounds))]

	var text, correct string
	var pool []string
	if s.Intn(2) == 0 {
		text = fmt.Sprintf("สาร %s (%s) จัดเป็นสารประกอบอินทรีย์ประเภทใด?", c.NameTH, c.Formula)
		correct, pool = c.Class, organicClasses
	} else {
		text = fmt.Sprintf("หมู่ฟังก์ชันของ %s (%s) คือข้อใด?", c.NameTH, c.Formula)
		correct, pool = c.Group, functionalGroups
	}

	candidates := textCandidates(s, correct, pool)
	choices, idx := assembleChoices(s, correct, candidates, noTopUp)

	return Problem{
		Question:     text,
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// ─── Polymers (บทที่ 13) ──────────────────────────────────────────

type polymerTopic struct{}

func (polymerTopic) Category() string   { return "polymer" }
func (polymerTopic) Difficulty() string { return "easy" }

func (t polymerTopic) Generate(s *Source) Problem {
	p := polymers[s.Pick("polymers", len(polymers))]
	if s.Intn(2) == 0 {
		return t.monomer(s, p)
	}
	return t.route(s, p)
}

func (t polymerTopic) monomer(s *Source, p polymer) Problem {
	pool := make([]string, 0, len(polymers))
	for _, other := range polymers {
		pool = append(pool, other.Monomer)
	}

	candidates := textCandidates(s, p.Monomer, pool)
	choices, idx := assembleChoices(s, p.Monomer, candidates, noTopUp)

	return Problem{
		Question:     fmt.Sprintf("%s เกิดจากมอนอเมอร์ใด?", p.NameTH),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}

// route asks which polymer comes from a given polymerisation, so the three
// wrong answers are real polymers made the other way — the exact confusion the
// question is about.
func (t polymerTopic) route(s *Source, p polymer) Problem {
	pool := make([]string, 0, len(polymers))
	for _, other := range polymers {
		if other.Route != p.Route {
			pool = append(pool, other.NameTH)
		}
	}

	candidates := textCandidates(s, p.NameTH, pool)
	choices, idx := assembleChoices(s, p.NameTH, candidates, noTopUp)

	return Problem{
		Question:     fmt.Sprintf("ข้อใดเป็นพอลิเมอร์ที่เกิดจากปฏิกิริยา%s?", p.Route),
		Choices:      choices,
		CorrectIndex: idx,
		Category:     t.Category(),
		Difficulty:   t.Difficulty(),
	}
}
