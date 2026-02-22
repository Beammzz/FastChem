package services

import (
	"sync"
	"time"
)

// StoredQuestion holds server-side data for an issued question (anti-cheat).
type StoredQuestion struct {
	CorrectIndex int
	Difficulty   string
	IssuedAt     time.Time
}

// QuestionStore is a thread-safe in-memory store for issued questions.
// It maps questionID → StoredQuestion so the server can validate answers
// and calculate time_spent independently of the client.
type QuestionStore struct {
	mu    sync.RWMutex
	store map[string]StoredQuestion
}

// GlobalQuestionStore is the singleton question store instance.
var GlobalQuestionStore = &QuestionStore{
	store: make(map[string]StoredQuestion),
}

// Store records a newly issued question.
func (qs *QuestionStore) Store(id string, correctIndex int, difficulty string) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	qs.store[id] = StoredQuestion{
		CorrectIndex: correctIndex,
		Difficulty:   difficulty,
		IssuedAt:     time.Now(),
	}
}

// Get retrieves a stored question by ID.
func (qs *QuestionStore) Get(id string) (StoredQuestion, bool) {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	q, ok := qs.store[id]
	return q, ok
}

// Delete removes a question from the store (after it's been answered).
func (qs *QuestionStore) Delete(id string) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	delete(qs.store, id)
}

// CleanupOlderThan removes all questions older than the given duration.
// Call periodically to prevent memory leaks from unanswered questions.
func (qs *QuestionStore) CleanupOlderThan(maxAge time.Duration) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for id, q := range qs.store {
		if q.IssuedAt.Before(cutoff) {
			delete(qs.store, id)
		}
	}
}
