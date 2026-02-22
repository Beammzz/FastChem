package services

import (
	"log"
	"sync"
	"time"
)

// QueueEntry represents a player waiting in the matchmaking queue.
type QueueEntry struct {
	UserID   int64
	Username string
	Rating   int
	JoinedAt time.Time
	// Notify channel — match service sends the matchID when matched.
	Matched chan int64
}

// MatchmakingQueue manages the ranked matchmaking queue.
type MatchmakingQueue struct {
	mu           sync.Mutex
	queue        []*QueueEntry
	matchService *RankedMatchService
}

// NewMatchmakingQueue creates a new matchmaking queue.
func NewMatchmakingQueue(matchService *RankedMatchService) *MatchmakingQueue {
	mq := &MatchmakingQueue{
		queue:        make([]*QueueEntry, 0),
		matchService: matchService,
	}
	// Start the matchmaking loop
	go mq.matchmakingLoop()
	return mq
}

// Enqueue adds a player to the queue. Returns the notification channel.
func (mq *MatchmakingQueue) Enqueue(userID int64, username string, rating int) chan int64 {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	// Check if already in queue
	for _, e := range mq.queue {
		if e.UserID == userID {
			return e.Matched
		}
	}

	// Check if already in an active match
	if _, ok := mq.matchService.GetMatchByUserID(userID); ok {
		return nil // already in a match
	}

	entry := &QueueEntry{
		UserID:   userID,
		Username: username,
		Rating:   rating,
		JoinedAt: time.Now(),
		Matched:  make(chan int64, 1),
	}
	mq.queue = append(mq.queue, entry)
	return entry.Matched
}

// Dequeue removes a player from the queue (e.g., on cancel or disconnect).
func (mq *MatchmakingQueue) Dequeue(userID int64) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for i, e := range mq.queue {
		if e.UserID == userID {
			mq.queue = append(mq.queue[:i], mq.queue[i+1:]...)
			close(e.Matched)
			return
		}
	}
}

// IsInQueue checks if a user is currently in the queue.
func (mq *MatchmakingQueue) IsInQueue(userID int64) bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for _, e := range mq.queue {
		if e.UserID == userID {
			return true
		}
	}
	return false
}

// QueueSize returns the current number of players in the queue.
func (mq *MatchmakingQueue) QueueSize() int {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.queue)
}

// matchmakingLoop runs continuously to match players.
func (mq *MatchmakingQueue) matchmakingLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		mq.tryMatch()
	}
}

// tryMatch attempts to match two players from the queue.
// Simple approach: match closest-rated players if within rating window.
// Rating window expands with wait time.
func (mq *MatchmakingQueue) tryMatch() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.queue) < 2 {
		return
	}

	now := time.Now()

	// Try to find a pair with closest ratings
	bestI, bestJ := -1, -1
	bestDiff := int(^uint(0) >> 1) // max int

	for i := 0; i < len(mq.queue); i++ {
		for j := i + 1; j < len(mq.queue); j++ {
			p1 := mq.queue[i]
			p2 := mq.queue[j]

			ratingDiff := abs(p1.Rating - p2.Rating)

			// Rating window: starts at 100, expands by 50 per second waited
			p1Wait := now.Sub(p1.JoinedAt).Seconds()
			p2Wait := now.Sub(p2.JoinedAt).Seconds()
			maxWait := p1Wait
			if p2Wait > maxWait {
				maxWait = p2Wait
			}

			ratingWindow := 100 + int(maxWait*50)

			if ratingDiff <= ratingWindow && ratingDiff < bestDiff {
				bestI, bestJ = i, j
				bestDiff = ratingDiff
			}
		}
	}

	if bestI == -1 || bestJ == -1 {
		return
	}

	p1 := mq.queue[bestI]
	p2 := mq.queue[bestJ]

	// Remove both from queue (remove higher index first)
	if bestI > bestJ {
		mq.queue = append(mq.queue[:bestI], mq.queue[bestI+1:]...)
		mq.queue = append(mq.queue[:bestJ], mq.queue[bestJ+1:]...)
	} else {
		mq.queue = append(mq.queue[:bestJ], mq.queue[bestJ+1:]...)
		mq.queue = append(mq.queue[:bestI], mq.queue[bestI+1:]...)
	}

	// Create match (outside lock to avoid deadlock)
	go func() {
		seed := time.Now().UnixNano()
		match, err := mq.matchService.CreateMatch(
			p1.UserID, p1.Username, p1.Rating,
			p2.UserID, p2.Username, p2.Rating,
			seed,
		)
		if err != nil {
			log.Printf("matchmaking: failed to create match: %v", err)
			return
		}

		// Notify both players
		select {
		case p1.Matched <- match.MatchID:
		default:
		}
		select {
		case p2.Matched <- match.MatchID:
		default:
		}
	}()
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
