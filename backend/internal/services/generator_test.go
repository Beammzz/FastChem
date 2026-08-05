package services

import (
	"fmt"
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
