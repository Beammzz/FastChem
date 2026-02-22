package services

import (
	"log/slog"
	"sort"
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
// The queue is kept sorted by rating for efficient O(n) sliding-window pairing.
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

// Enqueue adds a player to the queue in sorted-by-rating order.
// Returns the notification channel.
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

	// Sorted insert by rating (ascending)
	i := sort.Search(len(mq.queue), func(j int) bool {
		return mq.queue[j].Rating >= rating
	})
	mq.queue = append(mq.queue, nil)
	copy(mq.queue[i+1:], mq.queue[i:])
	mq.queue[i] = entry

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

// tryMatch attempts to match two players from the sorted queue.
// Uses a sliding-window approach over the rating-sorted queue: for each
// player, only the adjacent neighbour is checked. This gives O(n) per tick
// instead of the previous O(n²) brute-force.
func (mq *MatchmakingQueue) tryMatch() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	if len(mq.queue) < 2 {
		return
	}

	now := time.Now()

	// Sliding window: scan adjacent pairs and pick the best eligible match.
	bestI := -1
	bestDiff := int(^uint(0) >> 1)

	for i := 0; i < len(mq.queue)-1; i++ {
		p1 := mq.queue[i]
		p2 := mq.queue[i+1]

		ratingDiff := p2.Rating - p1.Rating // always >= 0 (sorted)

		// Rating window: starts at 100, expands by 50 per second waited
		p1Wait := now.Sub(p1.JoinedAt).Seconds()
		p2Wait := now.Sub(p2.JoinedAt).Seconds()
		maxWait := p1Wait
		if p2Wait > maxWait {
			maxWait = p2Wait
		}

		ratingWindow := 100 + int(maxWait*50)

		if ratingDiff <= ratingWindow && ratingDiff < bestDiff {
			bestI = i
			bestDiff = ratingDiff
		}
	}

	if bestI == -1 {
		return
	}

	p1 := mq.queue[bestI]
	p2 := mq.queue[bestI+1]

	// Remove both from queue (higher index first to keep indices valid)
	mq.queue = append(mq.queue[:bestI+1], mq.queue[bestI+2:]...)
	mq.queue = append(mq.queue[:bestI], mq.queue[bestI+1:]...)

	// Create match (outside lock to avoid deadlock)
	go func() {
		seed := time.Now().UnixNano()
		match, err := mq.matchService.CreateMatch(
			p1.UserID, p1.Username, p1.Rating,
			p2.UserID, p2.Username, p2.Rating,
			seed,
		)
		if err != nil {
			slog.Error("matchmaking: failed to create match", "error", err)
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
