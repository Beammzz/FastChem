package services

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

// withDeadline fails the test if fn has not returned in time, instead of
// letting a non-terminating generator hang the whole test binary.
func withDeadline(t *testing.T, name string, d time.Duration, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not return within %s (infinite loop regression)", name, d)
	}
}

func assertWellFormed(t *testing.T, choices []string, idx int, want string) {
	t.Helper()
	if len(choices) != 4 {
		t.Fatalf("expected 4 choices for %q, got %d: %v", want, len(choices), choices)
	}
	seen := map[string]bool{}
	for _, c := range choices {
		if c == "" {
			t.Fatalf("empty choice in %v", choices)
		}
		if seen[c] {
			t.Fatalf("duplicate choice %q in %v", c, choices)
		}
		seen[c] = true
	}
	if idx < 0 || idx >= len(choices) {
		t.Fatalf("correctIndex %d out of range for %v", idx, choices)
	}
	if choices[idx] != want {
		t.Fatalf("choices[%d] = %q, want %q (%v)", idx, choices[idx], want, choices)
	}
}

// ─── Determinism ────────────────────────────────────────────────

// The seed on a ranked match row is only worth storing if it reproduces the
// match. Choices used to be collected by ranging over a map, so the same seed
// produced the same questions in a different choice order roughly 58% of the
// time, moving the correct answer's index.
func TestGenerateRankedQuestionsIsDeterministic(t *testing.T) {
	withDeadline(t, "determinism sweep", 60*time.Second, func() {
		for seed := int64(1); seed <= 400; seed++ {
			a := GenerateRankedQuestions(seed)
			b := GenerateRankedQuestions(seed)
			if len(a.Questions) != len(b.Questions) {
				t.Fatalf("seed %d: question count differs", seed)
			}
			for i := range a.Questions {
				qa, qb := a.Questions[i], b.Questions[i]
				if qa.Question != qb.Question {
					t.Fatalf("seed %d q%d: text differs\n A: %s\n B: %s", seed, i, qa.Question, qb.Question)
				}
				if strings.Join(qa.Choices, "|") != strings.Join(qb.Choices, "|") {
					t.Fatalf("seed %d q%d: choice order differs\n A: %v\n B: %v", seed, i, qa.Choices, qb.Choices)
				}
				if qa.CorrectIndex != qb.CorrectIndex {
					t.Fatalf("seed %d q%d: correctIndex %d vs %d", seed, i, qa.CorrectIndex, qb.CorrectIndex)
				}
				if qa.Category != qb.Category || qa.Difficulty != qb.Difficulty {
					t.Fatalf("seed %d q%d: metadata differs", seed, i)
				}
			}
		}
	})
}

// Different seeds must actually produce different matches.
func TestDifferentSeedsDiffer(t *testing.T) {
	first := GenerateRankedQuestions(1).Questions[0].Question
	same := 0
	for seed := int64(2); seed <= 50; seed++ {
		if GenerateRankedQuestions(seed).Questions[0].Question == first {
			same++
		}
	}
	if same > 10 {
		t.Errorf("%d of 49 seeds produced the same first question; seeding looks ineffective", same)
	}
}

// ─── Termination ────────────────────────────────────────────────

// Seed 99 used to spin forever: mass→moles produced 0.03 mol, and fewer than
// three distinct "%.2f" values exist within ±50% of it.
func TestGenerateRankedQuestionsTerminatesOnSeed99(t *testing.T) {
	withDeadline(t, "GenerateRankedQuestions(99)", 10*time.Second, func() {
		set := GenerateRankedQuestions(99)
		if got := len(set.Questions); got != 10 {
			t.Errorf("expected 10 questions, got %d", got)
		}
	})
}

// The distractor ladder must fill four slots for any answer, however small,
// in any of the formats the topics use.
func TestChoiceAssemblyTerminatesForSmallAnswers(t *testing.T) {
	formats := []string{"%.2f mol", "%.2f g", "%.1f mL", "%.2f °C"}
	s := NewSeededSource(1)
	withDeadline(t, "small-answer sweep", 30*time.Second, func() {
		for _, format := range formats {
			for cents := 1; cents <= 80; cents++ {
				for _, sign := range []float64{1, -1} {
					correct := sign * float64(cents) / 100
					want := fmt.Sprintf(format, correct)
					choices, idx := assembleChoices(s, want, nil, floatTopUp(correct, format))
					assertWellFormed(t, choices, idx, want)
				}
			}
		}
	})
}

func TestSciChoiceAssemblyTerminates(t *testing.T) {
	s := NewSeededSource(2)
	withDeadline(t, "sci sweep", 10*time.Second, func() {
		for _, v := range []float64{6.02e23, 6.02e-5, 1.0, 1.2e10, 9.99e22} {
			want := formatSci(v)
			choices, idx := assembleChoices(s, want, nil, sciTopUp(v))
			assertWellFormed(t, choices, idx, want)
		}
	})
}

// ─── Well-formedness ────────────────────────────────────────────

// Scoring resolves correctness by index, so a malformed problem mis-scores
// silently. Every topic must produce four distinct choices every time.
func TestEveryTopicProducesWellFormedProblems(t *testing.T) {
	withDeadline(t, "topic sweep", 60*time.Second, func() {
		for _, topic := range topicRegistry {
			s := NewSeededSource(int64(len(topic.Category())))
			for i := 0; i < 2000; i++ {
				p := topic.Generate(s)
				if !validProblem(p) {
					t.Fatalf("%s produced an invalid problem: %+v", topic.Category(), p)
				}
				if p.Category != topic.Category() || p.Difficulty != topic.Difficulty() {
					t.Fatalf("%s: metadata mismatch: %+v", topic.Category(), p)
				}
			}
		}
	})
}

func TestGeneratedQuestionsAreWellFormed(t *testing.T) {
	g := NewQuestionGenerator()
	withDeadline(t, "question sweep", 30*time.Second, func() {
		for _, difficulty := range []string{"easy", "medium", "hard", "nonsense"} {
			for i := 0; i < 500; i++ {
				q := g.GenerateQuestion(difficulty)
				assertWellFormed(t, q.Choices, q.CorrectIndex, q.Choices[q.CorrectIndex])
				if q.Question == "" {
					t.Fatalf("%s: empty question text", difficulty)
				}
				if q.TimeLimit <= 0 {
					t.Fatalf("%s: non-positive time limit %d", difficulty, q.TimeLimit)
				}
			}
		}
	})
}

// ─── Mode parity ────────────────────────────────────────────────

// Ranked play used to draw easy questions from only two of the three easy
// topics, so state_of_matter never appeared in a 1v1 match.
func TestRankedCoversEveryTopic(t *testing.T) {
	seen := map[string]bool{}
	withDeadline(t, "coverage sweep", 60*time.Second, func() {
		for seed := int64(1); seed <= 600; seed++ {
			for _, q := range GenerateRankedQuestions(seed).Questions {
				seen[q.Category] = true
			}
		}
	})
	for _, topic := range topicRegistry {
		if !seen[topic.Category()] {
			t.Errorf("category %q never appeared in 600 ranked matches", topic.Category())
		}
	}
}

// Category IDs are load-bearing in the frontend, which switches on them for
// labels and colours.
func TestTopicCategoriesMatchFrontendIDs(t *testing.T) {
	want := []string{
		"atomic_structure", "oxidation_number", "state_of_matter",
		"mole_concept", "dilution", "preparing_solution", "freezing_point",
		"electron_config", "bond_type", "molecular_shape",
		"molar_mass", "percent_composition", "concentration",
		"stoichiometry", "limiting_reagent", "gas_law", "ideal_gas",
		"reaction_rate", "equilibrium", "acid_base", "electrochemistry",
		"functional_group", "polymer",
	}
	for _, category := range want {
		if _, ok := topicForCategory(category); !ok {
			t.Errorf("no topic registered for category %q", category)
		}
	}
	if len(topicRegistry) != len(want) {
		t.Errorf("registry has %d topics, expected %d", len(topicRegistry), len(want))
	}
}

// Every topic must be reachable through the difficulty filter as well as the
// category filter — a topic whose Difficulty() is a typo would be registered,
// pass every other test, and never be served.
func TestEveryTopicIsReachableByDifficulty(t *testing.T) {
	seen := map[string]bool{}
	for _, difficulty := range []string{"easy", "medium", "hard"} {
		for _, topic := range topicsForDifficulty(difficulty) {
			if topic.Difficulty() != difficulty {
				t.Errorf("%s appeared under difficulty %q", topic.Category(), difficulty)
			}
			seen[topic.Category()] = true
		}
	}
	for _, topic := range topicRegistry {
		if !seen[topic.Category()] {
			t.Errorf("%s has difficulty %q, which no filter selects",
				topic.Category(), topic.Difficulty())
		}
	}
}

func TestGenerateByCategoryHonoursTheRequest(t *testing.T) {
	g := NewQuestionGenerator()
	for _, category := range []string{"state_of_matter", "dilution", "mole_concept"} {
		for i := 0; i < 100; i++ {
			q := g.GenerateQuestionByCategories([]string{category})
			if q.Category != category {
				t.Fatalf("requested %q, got %q", category, q.Category)
			}
		}
	}
	// Unknown categories fall back rather than failing.
	q := g.GenerateQuestionByCategories([]string{"not_a_category"})
	if q.Question == "" || len(q.Choices) != 4 {
		t.Errorf("fallback question malformed: %+v", q)
	}
}

// ─── Chemistry sanity ───────────────────────────────────────────

// Neutron questions used to offset the atomic number by a random 0–10,
// inventing nuclides like hydrogen-7.
func TestNeutronQuestionsUseRealIsotopes(t *testing.T) {
	real := map[string]map[int]bool{}
	for _, e := range elements {
		real[e.Symbol] = map[int]bool{}
		for _, a := range e.Isotopes {
			real[e.Symbol][a] = true
		}
	}

	topic := atomicStructureTopic{}
	s := NewSeededSource(11)
	checked := 0
	for i := 0; i < 20000; i++ {
		p := topic.Generate(s)
		var symbol string
		var massNumber int
		if n, _ := fmt.Sscanf(p.Question, "ไอโซโทปของ %s", &symbol); n == 0 {
			continue
		}
		if !strings.Contains(p.Question, "เลขมวล") {
			continue
		}
		// "…(Sym) ที่มีเลขมวล N มี…"
		var sym string
		if _, err := fmt.Sscanf(
			p.Question[strings.Index(p.Question, "("):],
			"(%s", &sym,
		); err != nil {
			continue
		}
		sym = strings.TrimSuffix(sym, ")")
		after := p.Question[strings.Index(p.Question, "เลขมวล")+len("เลขมวล"):]
		if _, err := fmt.Sscanf(strings.TrimSpace(after), "%d", &massNumber); err != nil {
			continue
		}
		if isotopes, ok := real[sym]; ok {
			if !isotopes[massNumber] {
				t.Fatalf("invented isotope %s-%d in %q", sym, massNumber, p.Question)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no neutron questions were generated, test proved nothing")
	}
	t.Logf("verified %d isotope questions against the real isotope table", checked)
}

// Oxidation distractors must stay inside the range chemistry actually uses.
func TestOxidationChoicesStayInRange(t *testing.T) {
	topic := oxidationNumberTopic{}
	s := NewSeededSource(5)
	for i := 0; i < 5000; i++ {
		p := topic.Generate(s)
		for _, c := range p.Choices {
			var v int
			if _, err := fmt.Sscanf(strings.TrimPrefix(c, "+"), "%d", &v); err != nil {
				t.Fatalf("unparseable oxidation choice %q in %v", c, p.Choices)
			}
			if v < minOxidation || v > maxOxidation {
				t.Errorf("oxidation choice %d outside [%d,%d] in %v", v, minOxidation, maxOxidation, p.Choices)
			}
		}
	}
}

// Counts of particles cannot be negative.
func TestParticleCountsAreNonNegative(t *testing.T) {
	topic := atomicStructureTopic{}
	s := NewSeededSource(7)
	for i := 0; i < 5000; i++ {
		p := topic.Generate(s)
		for _, c := range p.Choices {
			var v int
			if _, err := fmt.Sscanf(c, "%d", &v); err != nil {
				t.Fatalf("unparseable count %q in %v", c, p.Choices)
			}
			if v < 0 {
				t.Errorf("negative particle count %d in %v", v, p.Choices)
			}
		}
	}
}

// A missing key in atomicMasses reads as 0.0 rather than failing, which would
// quietly under-count a molar mass and take the answer and every distractor
// down with it.
func TestEveryFormulaSymbolHasAnAtomicMass(t *testing.T) {
	for _, c := range formulaCompounds {
		for _, p := range c.Parts {
			if _, ok := atomicMasses[p.Symbol]; !ok {
				t.Errorf("%s: no atomic mass for %q", c.Formula, p.Symbol)
			}
			if p.Count < 1 {
				t.Errorf("%s: subscript %d on %s", c.Formula, p.Count, p.Symbol)
			}
		}
	}
}

// Spot-check the derived molar masses against the values the IPST book prints.
func TestFormulaMolarMassesAreCorrect(t *testing.T) {
	want := map[string]float64{
		"H₂O":       18.02,
		"CO₂":       44.01,
		"NaCl":      58.44,
		"CaCO₃":     100.09,
		"H₂SO₄":     98.09,
		"NH₃":       17.03,
		"C₆H₁₂O₆":   180.16,
		"KMnO₄":     158.04,
		"K₂Cr₂O₇":   294.20,
		"(NH₄)₂SO₄": 132.15,
	}
	found := 0
	for _, c := range formulaCompounds {
		expected, ok := want[c.Formula]
		if !ok {
			continue
		}
		found++
		if math.Abs(c.MolarMass()-expected) > 0.02 {
			t.Errorf("%s: molar mass %.3f, want %.2f", c.Formula, c.MolarMass(), expected)
		}
	}
	if found != len(want) {
		t.Errorf("checked %d of %d reference compounds; the table was renamed or trimmed", found, len(want))
	}
}

// A stoichiometry question is only meaningful if the equation it quotes is
// balanced. Mass balance catches both a wrong coefficient and a wrong molar
// mass, which is exactly what would make the "correct" answer wrong.
func TestReactionsConserveMass(t *testing.T) {
	side := func(species []reactionSpecies) float64 {
		total := 0.0
		for _, sp := range species {
			total += float64(sp.Coef) * sp.MolarMass
		}
		return total
	}
	for _, r := range reactions {
		in, out := side(r.Reactants), side(r.Products)
		// The tolerance absorbs rounding in the tabulated molar masses only;
		// a missing or wrong coefficient moves the total by grams, not
		// hundredths.
		if math.Abs(in-out) > 0.05 {
			t.Errorf("%s: %.2f g on the left, %.2f g on the right", r.Equation, in, out)
		}
		if len(r.Reactants) == 0 || len(r.Products) == 0 {
			t.Errorf("%s: one side is empty", r.Equation)
		}
	}
}

// Closed-vocabulary topics have no ladder to fall back on: if a pool holds
// fewer than four distinct entries, the topic cannot fill four choices and
// every problem it makes fails validation.
func TestClosedVocabularyPoolsCanFillFourChoices(t *testing.T) {
	pools := map[string][]string{
		"bondClasses":          bondClasses,
		"molecularShapes":      molecularShapes,
		"organicClasses":       organicClasses,
		"functionalGroups":     functionalGroups,
		"representativeGroups": representativeGroups,
	}
	for name, pool := range pools {
		distinct := map[string]bool{}
		for _, v := range pool {
			if v == "" {
				t.Errorf("%s contains an empty entry", name)
			}
			distinct[v] = true
		}
		if len(distinct) < 4 {
			t.Errorf("%s has %d distinct entries, need at least 4", name, len(distinct))
		}
	}

	// The polymer route question draws its distractors from the other route,
	// so each route needs three polymers of its own plus the answer.
	routes := map[string]int{}
	monomers := map[string]string{}
	for _, p := range polymers {
		routes[p.Route]++
		if prev, ok := monomers[p.Monomer]; ok {
			t.Errorf("monomer %q maps to both %s and %s; the question has two right answers",
				p.Monomer, prev, p.NameTH)
		}
		monomers[p.Monomer] = p.NameTH
	}
	for route, n := range routes {
		if n < 4 {
			t.Errorf("route %q has %d polymers, need at least 4", route, n)
		}
	}

	// Every substance and molecule must classify into the vocabulary its
	// question answers from.
	known := func(pool []string, v string) bool {
		for _, entry := range pool {
			if entry == v {
				return true
			}
		}
		return false
	}
	for _, sub := range bondSubstances {
		if !known(bondClasses, sub.Class) {
			t.Errorf("%s has class %q, which is not in bondClasses", sub.Formula, sub.Class)
		}
	}
	for _, mol := range shapeMolecules {
		if !known(molecularShapes, mol.Shape) {
			t.Errorf("%s has shape %q, which is not in molecularShapes", mol.Formula, mol.Shape)
		}
	}
	for _, c := range organicCompounds {
		if !known(organicClasses, c.Class) {
			t.Errorf("%s has class %q, which is not in organicClasses", c.Formula, c.Class)
		}
		if !known(functionalGroups, c.Group) {
			t.Errorf("%s has group %q, which is not in functionalGroups", c.Formula, c.Group)
		}
	}
}

// Shell configurations must add up to the atomic number, and no level may hold
// more electrons than it can.
func TestShellConfigurationsAreConsistent(t *testing.T) {
	capacities := []int{2, 8, 18, 32}
	for _, e := range shellElements {
		total := 0
		for i, n := range e.Shells {
			total += n
			if i < len(capacities) && n > capacities[i] {
				t.Errorf("%s: level %d holds %d electrons, capacity %d", e.Symbol, i+1, n, capacities[i])
			}
		}
		if total != e.Number {
			t.Errorf("%s: shells sum to %d, atomic number is %d", e.Symbol, total, e.Number)
		}
		if e.Period != len(e.Shells) {
			t.Errorf("%s: period %d but %d occupied levels", e.Symbol, e.Period, len(e.Shells))
		}
	}
}

// A galvanic cell runs the higher reduction potential as the cathode, so a
// standard cell potential is never negative.
func TestCellPotentialsAreNonNegative(t *testing.T) {
	topic := electrochemistryTopic{}
	s := NewSeededSource(13)
	checked := 0
	for i := 0; i < 5000; i++ {
		p := topic.Generate(s)
		if !strings.Contains(p.Question, "E°cell") {
			continue
		}
		var v float64
		if _, err := fmt.Sscanf(p.Choices[p.CorrectIndex], "%f", &v); err != nil {
			t.Fatalf("unparseable cell potential %q", p.Choices[p.CorrectIndex])
		}
		if v < 0 {
			t.Errorf("negative E°cell %.2f in %q", v, p.Question)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no cell-potential questions were generated, test proved nothing")
	}
}

// pH questions must land in the range the scale actually covers; a pH of 20 or
// −3 means a sign or a logarithm went the wrong way.
func TestAcidBaseAnswersStayOnTheScale(t *testing.T) {
	topic := acidBaseTopic{}
	s := NewSeededSource(17)
	for i := 0; i < 5000; i++ {
		p := topic.Generate(s)
		var v float64
		if _, err := fmt.Sscanf(p.Choices[p.CorrectIndex], "%f", &v); err != nil {
			t.Fatalf("unparseable pH %q", p.Choices[p.CorrectIndex])
		}
		if v < 0 || v > 14 {
			t.Errorf("pH %.2f outside [0,14] in %q", v, p.Question)
		}
	}
}

// ─── Concurrency ────────────────────────────────────────────────

// The casual generator holds one Source for the life of the process and Gin
// serves requests in parallel. Run with -race.
func TestConcurrentGenerationIsRaceFree(t *testing.T) {
	g := NewQuestionGenerator()
	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = g.GenerateQuestion("easy")
				_ = g.GenerateQuestionByCategories([]string{"oxidation_number", "state_of_matter"})
			}
		}()
	}
	wg.Wait()
}

// ─── The answer stays server-side ───────────────────────────────

// Ranked play sends models.QuestionStartPayload, not this struct, but a
// serializable answer field is one careless handler away from going out.
func TestRankedQuestionJSONOmitsTheAnswer(t *testing.T) {
	set := GenerateRankedQuestions(7)
	encoded, err := json.Marshal(set.Questions[0])
	if err != nil {
		t.Fatalf("marshal ranked question: %v", err)
	}
	if strings.Contains(string(encoded), "correctIndex") {
		t.Fatalf("ranked question JSON leaks the answer: %s", encoded)
	}
}
