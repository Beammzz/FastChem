package services

import (
	"fmt"
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

// Seed 99 used to spin forever in seededFloatDistractors: mass→moles produced
// 0.03 mol, and fewer than three distinct "%.2f" values exist within ±50% of it.
func TestGenerateRankedQuestionsTerminatesOnSeed99(t *testing.T) {
	withDeadline(t, "GenerateRankedQuestions(99)", 10*time.Second, func() {
		set := GenerateRankedQuestions(99)
		if got := len(set.Questions); got != 10 {
			t.Errorf("expected 10 questions, got %d", got)
		}
	})
}

// Sweep the range of small answers that previously had too few distinct
// formatted values in the ±50% band to fill four choices.
func TestFloatDistractorsTerminateForSmallAnswers(t *testing.T) {
	formats := []string{"%.2f mol", "%.2f g", "%.1f mL", "%.2f °C"}
	withDeadline(t, "generateFloatDistractors sweep", 20*time.Second, func() {
		for _, format := range formats {
			for cents := 1; cents <= 60; cents++ {
				correct := float64(cents) / 100
				for _, v := range []float64{correct, -correct} {
					choices, idx := generateFloatDistractors(v, format)
					assertWellFormed(t, choices, idx, fmt.Sprintf(format, v))
				}
			}
		}
	})
}

func TestSciDistractorsTerminate(t *testing.T) {
	withDeadline(t, "generateSciDistractors sweep", 10*time.Second, func() {
		for _, v := range []float64{6.02e23, 6.02e-5, 1.0, 1.2e10, 9.99e22} {
			choices, idx := generateSciDistractors(v)
			assertWellFormed(t, choices, idx, formatSci(v))
		}
	})
}

// Every generated question must offer four distinct choices with the correct
// answer among them — the invariant the game's scoring depends on.
func TestGeneratedQuestionsAreWellFormed(t *testing.T) {
	g := NewQuestionGenerator()
	withDeadline(t, "question sweep", 30*time.Second, func() {
		for _, difficulty := range []string{"easy", "medium", "hard"} {
			for i := 0; i < 500; i++ {
				q := g.GenerateQuestion(difficulty)
				if len(q.Choices) != 4 {
					t.Fatalf("%s: expected 4 choices, got %d (%v)", difficulty, len(q.Choices), q.Choices)
				}
				if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Choices) {
					t.Fatalf("%s: correctIndex %d out of range", difficulty, q.CorrectIndex)
				}
				seen := map[string]bool{}
				for _, c := range q.Choices {
					if seen[c] {
						t.Fatalf("%s: duplicate choice %q in %v", difficulty, c, q.Choices)
					}
					seen[c] = true
				}
				if q.Question == "" {
					t.Fatalf("%s: empty question text", difficulty)
				}
			}
		}
	})
}

// pickFreshIndex mutates package-level slices; Gin serves requests concurrently.
// Run with -race.
func TestConcurrentGenerationIsRaceFree(t *testing.T) {
	g := NewQuestionGenerator()
	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = g.GenerateQuestion("easy")
			}
		}()
	}
	wg.Wait()
}

func assertWellFormed(t *testing.T, choices []string, idx int, want string) {
	t.Helper()
	if len(choices) != 4 {
		t.Fatalf("expected 4 choices for %q, got %d: %v", want, len(choices), choices)
	}
	seen := map[string]bool{}
	for _, c := range choices {
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
