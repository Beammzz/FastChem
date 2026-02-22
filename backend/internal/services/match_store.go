package services

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/takumi/fastchem/internal/models"
)

// ActiveMatch holds server-side state for an in-progress match.
type ActiveMatch struct {
	MatchID    int64
	UserID     int64
	Difficulty string
	Combo      int
	BestCombo  int
	TotalScore int
	Answered   int
	Correct    int
	StartedAt  time.Time
	Attempts   []models.QuestionAttempt
}

// MatchStore is a thread-safe in-memory store for active matches.
type MatchStore struct {
	mu    sync.RWMutex
	store map[int64]*ActiveMatch // matchID → ActiveMatch
}

// GlobalMatchStore is the singleton match store instance.
var GlobalMatchStore = &MatchStore{
	store: make(map[int64]*ActiveMatch),
}

// Create registers a new active match.
func (ms *MatchStore) Create(matchID, userID int64, difficulty string) *ActiveMatch {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	m := &ActiveMatch{
		MatchID:    matchID,
		UserID:     userID,
		Difficulty: difficulty,
		StartedAt:  time.Now(),
		Attempts:   make([]models.QuestionAttempt, 0, models.DefaultQuestionsPerMatch),
	}
	ms.store[matchID] = m
	return m
}

// Get retrieves an active match by ID.
func (ms *MatchStore) Get(matchID int64) (*ActiveMatch, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	m, ok := ms.store[matchID]
	return m, ok
}

// Delete removes a match from the active store.
func (ms *MatchStore) Delete(matchID int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.store, matchID)
}

// CleanupOlderThan removes stale matches (e.g. abandoned).
func (ms *MatchStore) CleanupOlderThan(maxAge time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for id, m := range ms.store {
		if m.StartedAt.Before(cutoff) {
			delete(ms.store, id)
		}
	}
}

// RecordAnswer processes an answer within a match, updates combo/score atomically.
// Returns the MatchAnswerResponse and whether the match is complete.
func (ms *MatchStore) RecordAnswer(
	matchID int64,
	question StoredQuestion,
	questionID string,
	selectedIndex int,
	timeSpent float64,
) (*models.MatchAnswerResponse, bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m, ok := ms.store[matchID]
	if !ok {
		return nil, false, nil
	}

	cfg := GetDifficultyConfig(m.Difficulty)
	timedOut := timeSpent >= float64(cfg.TimeLimit)
	correct := selectedIndex == question.CorrectIndex && !timedOut

	// Update combo
	if correct {
		m.Combo++
		if m.Combo > m.BestCombo {
			m.BestCombo = m.Combo
		}
	} else {
		m.Combo = 0
	}

	// Calculate score with combo
	finalScore, _, speedBonus, comboMult := CalculateScoreWithCombo(
		m.Difficulty, timeSpent, correct, m.Combo,
	)

	m.TotalScore += finalScore
	m.Answered++
	if correct {
		m.Correct++
	}

	// Snapshot the question for permanent storage
	snapshot, _ := json.Marshal(map[string]interface{}{
		"questionId":   questionID,
		"correctIndex": question.CorrectIndex,
		"difficulty":   question.Difficulty,
	})

	attempt := models.QuestionAttempt{
		MatchID:          matchID,
		QuestionSnapshot: string(snapshot),
		Correct:          correct,
		TimedOut:         timedOut,
		TimeSpentMs:      int(timeSpent * 1000),
		ScoreAwarded:     finalScore,
		ComboAtAnswer:    m.Combo,
		ComboMultiplier:  comboMult,
		CreatedAt:        time.Now(),
	}
	m.Attempts = append(m.Attempts, attempt)

	resp := &models.MatchAnswerResponse{
		Correct:         correct,
		CorrectIndex:    question.CorrectIndex,
		ScoreEarned:     finalScore,
		SpeedBonus:      speedBonus,
		ComboMultiplier: comboMult,
		Combo:           m.Combo,
		TimeSpent:       timeSpent,
		TimedOut:        timedOut,
		TotalScore:      m.TotalScore,
		QuestionNumber:  m.Answered,
	}

	matchDone := m.Answered >= models.DefaultQuestionsPerMatch
	return resp, matchDone, nil
}
