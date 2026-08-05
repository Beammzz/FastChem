package services

import (
	"log/slog"
	"math/rand"
	"sync"
)

// Problem is one generated question in transport-neutral form. Casual play maps
// it onto models.Question; ranked and room matches map it onto RankedQuestion.
type Problem struct {
	Question     string
	Choices      []string
	CorrectIndex int
	Category     string
	Difficulty   string
}

// Topic generates one kind of chemistry problem. A single implementation
// serves casual and ranked play alike — the Source decides whether generation
// is reproducible, so the two modes cannot drift apart.
type Topic interface {
	Category() string
	Difficulty() string
	Generate(s *Source) Problem
}

// recentWindow is how many previous picks to remember per data pool.
const recentWindow = 8

// Source supplies randomness and short-term history to a Topic.
//
// Seeding it makes generation reproducible, which is what lets a ranked match
// be rebuilt from the seed stored on its row. Reproducibility only holds if
// nothing in the generation path depends on wall-clock time, the global rand,
// or Go map iteration order — see the contract in AGENTS.md.
//
// Every method is safe for concurrent use, because the casual generator holds
// one Source for the lifetime of the process and Gin serves requests in
// parallel.
type Source struct {
	mu     sync.Mutex
	rng    *rand.Rand
	recent map[string][]int
}

// NewSeededSource returns a Source that replays identically for a given seed.
func NewSeededSource(seed int64) *Source {
	return &Source{
		rng:    rand.New(rand.NewSource(seed)),
		recent: make(map[string][]int),
	}
}

func (s *Source) Intn(n int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rng.Intn(n)
}

func (s *Source) Float64() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rng.Float64()
}

func (s *Source) Shuffle(n int, swap func(i, j int)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rng.Shuffle(n, swap)
}

// Pick returns an index in [0, n) that is not among the last recentWindow
// picks for key, so consecutive questions do not reuse the same element or
// compound. Ranked matches get this too, which they did not before.
func (s *Source) Pick(key string, n int) int {
	if n <= 1 {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	hist := s.recent[key]
	for {
		idx := s.rng.Intn(n)
		dup := false
		for _, r := range hist {
			if r == idx {
				dup = true
				break
			}
		}
		if !dup {
			hist = append(hist, idx)
			if len(hist) > recentWindow {
				hist = hist[1:]
			}
			s.recent[key] = hist
			return idx
		}
		// Pool smaller than the window: forget history rather than spin.
		if n <= len(hist) {
			hist = nil
		}
	}
}

// topicRegistry lists every topic in a fixed order.
//
// The order is part of the seeded contract: selection indexes into slices
// derived from this one, so inserting a topic anywhere but the end changes
// what every existing seed generates.
var topicRegistry = []Topic{
	atomicStructureTopic{},
	oxidationNumberTopic{},
	stateOfMatterTopic{},
	moleConceptTopic{},
	dilutionTopic{},
	preparingSolutionTopic{},
	freezingPointTopic{},
}

// topicsForDifficulty returns the topics at a difficulty, in registry order.
// Unknown difficulties fall back to easy, matching GetDifficultyConfig.
func topicsForDifficulty(difficulty string) []Topic {
	out := make([]Topic, 0, len(topicRegistry))
	for _, t := range topicRegistry {
		if t.Difficulty() == difficulty {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		for _, t := range topicRegistry {
			if t.Difficulty() == "easy" {
				out = append(out, t)
			}
		}
	}
	return out
}

// topicForCategory looks a topic up by the category ID the frontend sends.
func topicForCategory(category string) (Topic, bool) {
	for _, t := range topicRegistry {
		if t.Category() == category {
			return t, true
		}
	}
	return nil, false
}

// validProblem reports whether p is safe to serve. Scoring resolves
// correctness by index, so a malformed problem would mis-score silently
// rather than fail loudly.
func validProblem(p Problem) bool {
	if p.Question == "" || len(p.Choices) != 4 {
		return false
	}
	if p.CorrectIndex < 0 || p.CorrectIndex >= len(p.Choices) {
		return false
	}
	seen := make(map[string]bool, len(p.Choices))
	for _, c := range p.Choices {
		if c == "" || seen[c] {
			return false
		}
		seen[c] = true
	}
	return true
}

// generateValid draws from a topic until the result passes validation. The
// retry is bounded: an unbounded "keep trying until it looks right" loop is
// exactly the shape that used to hang question generation.
func generateValid(s *Source, t Topic) Problem {
	var p Problem
	for i := 0; i < 5; i++ {
		p = t.Generate(s)
		if validProblem(p) {
			return p
		}
		slog.Warn("generated malformed problem", "category", t.Category(), "attempt", i+1)
	}
	return p
}
