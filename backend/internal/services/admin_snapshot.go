package services

import (
	"time"

	"github.com/takumi/fastchem/internal/models"
)

// Read-only views of the in-memory stores, for the admin console's live tab.
//
// Every function here copies under the owning lock and returns plain data. No
// caller gets a pointer into live match state, so a slow HTTP response cannot
// hold a match lock or read a score mid-update.

// Len reports how many questions are awaiting an answer.
func (qs *QuestionStore) Len() int {
	qs.mu.RLock()
	defer qs.mu.RUnlock()
	return len(qs.store)
}

// Len reports how many single-player sessions are in progress.
func (ms *MatchStore) Len() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.store)
}

// Snapshot lists everyone waiting for a ranked opponent, longest wait first.
func (mq *MatchmakingQueue) Snapshot() []models.AdminQueueEntry {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	now := time.Now()
	out := make([]models.AdminQueueEntry, 0, len(mq.queue))
	for _, e := range mq.queue {
		out = append(out, models.AdminQueueEntry{
			UserID:         e.UserID,
			Username:       e.Username,
			Rating:         e.Rating,
			WaitingSeconds: now.Sub(e.JoinedAt).Seconds(),
		})
	}
	return out
}

// Snapshot lists every match currently in memory — ranked and room alike,
// since both live in this service.
func (rms *RankedMatchService) Snapshot() []models.AdminLiveMatch {
	rms.mu.RLock()
	matches := make([]*ActiveRankedMatch, 0, len(rms.matches))
	for _, m := range rms.matches {
		matches = append(matches, m)
	}
	rms.mu.RUnlock()

	out := make([]models.AdminLiveMatch, 0, len(matches))
	for _, m := range matches {
		m.mu.Lock()
		out = append(out, models.AdminLiveMatch{
			MatchID:   m.MatchID,
			Status:    m.Status,
			Question:  m.LastSyncedIndex,
			CreatedAt: m.CreatedAt,
			Player1:   livePlayer(m.Player1),
			Player2:   livePlayer(m.Player2),
		})
		m.mu.Unlock()
	}
	return out
}

// ActiveUserIDs returns the set of players in a live match, for marking rows
// in the admin user table as online.
func (rms *RankedMatchService) ActiveUserIDs() map[int64]bool {
	rms.mu.RLock()
	defer rms.mu.RUnlock()
	out := make(map[int64]bool, len(rms.byUser))
	for userID := range rms.byUser {
		out[userID] = true
	}
	return out
}

func livePlayer(p *RankedPlayerState) models.AdminLivePlayer {
	if p == nil {
		return models.AdminLivePlayer{}
	}
	return models.AdminLivePlayer{
		UserID:     p.UserID,
		Username:   p.Username,
		Rating:     p.Rating,
		TotalScore: p.TotalScore,
		Answered:   p.Answered,
		Correct:    p.Correct,
		Combo:      p.Combo,
		Connected:  p.Connected,
	}
}

// Snapshot lists every open custom room.
func (rs *RoomService) Snapshot() []models.AdminRoom {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	out := make([]models.AdminRoom, 0, len(rs.rooms))
	for _, r := range rs.rooms {
		out = append(out, models.AdminRoom{
			Code:      r.Code,
			HostID:    r.HostID,
			HostName:  r.HostName,
			GuestID:   r.GuestID,
			GuestName: r.GuestName,
			MatchID:   r.MatchID,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
		})
	}
	return out
}

// TopicCatalog lists every registered question topic in registry order. The
// admin console uses it to browse and preview what the generator can ask.
func TopicCatalog() []models.AdminTopic {
	out := make([]models.AdminTopic, 0, len(topicRegistry))
	for _, t := range topicRegistry {
		out = append(out, models.AdminTopic{Category: t.Category(), Difficulty: t.Difficulty()})
	}
	return out
}

// PreviewQuestion generates a throwaway question for one topic, answer
// included. It is not stored in GlobalQuestionStore and cannot be scored, so
// nothing a player sees depends on it.
func (g *QuestionGenerator) PreviewQuestion(category string) (models.AdminQuestionPreview, bool) {
	topic, ok := topicForCategory(category)
	if !ok {
		return models.AdminQuestionPreview{}, false
	}
	p := generateValid(g.source, topic)
	return models.AdminQuestionPreview{
		Question:     p.Question,
		Choices:      p.Choices,
		CorrectIndex: p.CorrectIndex,
		Category:     p.Category,
		Difficulty:   p.Difficulty,
		TimeLimit:    GetDifficultyConfig(p.Difficulty).TimeLimit,
	}, true
}
