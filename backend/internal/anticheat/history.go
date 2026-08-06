package anticheat

import (
	"sync"
	"time"
)

// historyWindow bounds how many answers are kept per subject. It must exceed
// the largest rule window; a rule asking for more than this can never fire.
const historyWindow = 16

// Answer is one past answer, kept only for the fields a windowed rule reads.
type Answer struct {
	TimeSpent  float64
	Correct    bool
	TimedOut   bool
	Difficulty string
	At         time.Time
}

// History is a subject's recent answers, oldest first, ending with the answer
// currently under inspection. The slice is a copy — detectors run outside the
// tracker's lock and must not see it mutate underneath them.
type History struct {
	Recent []Answer
}

type subjectHistory struct {
	recent  []Answer
	updated time.Time
}

// tracker keeps a bounded, expiring window of answers per subject.
type tracker struct {
	mu       sync.Mutex
	subjects map[string]*subjectHistory
}

func newTracker() *tracker {
	return &tracker{subjects: make(map[string]*subjectHistory)}
}

// record appends sig to its subject's history and returns the window,
// including sig itself as the last entry.
func (t *tracker) record(sig Signal) History {
	t.mu.Lock()
	defer t.mu.Unlock()

	h, ok := t.subjects[sig.Subject]
	if !ok {
		h = &subjectHistory{recent: make([]Answer, 0, historyWindow)}
		t.subjects[sig.Subject] = h
	}

	h.recent = append(h.recent, Answer{
		TimeSpent:  sig.TimeSpent,
		Correct:    sig.Correct,
		TimedOut:   sig.TimedOut,
		Difficulty: sig.Difficulty,
		At:         time.Now(),
	})
	if len(h.recent) > historyWindow {
		h.recent = h.recent[len(h.recent)-historyWindow:]
	}
	h.updated = time.Now()

	return History{Recent: append([]Answer(nil), h.recent...)}
}

// sweep drops subjects idle for longer than maxIdle, so a long-running server
// does not accumulate a row per IP that ever answered a question.
func (t *tracker) sweep(maxIdle time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := time.Now().Add(-maxIdle)
	for k, h := range t.subjects {
		if h.updated.Before(cutoff) {
			delete(t.subjects, k)
		}
	}
}
