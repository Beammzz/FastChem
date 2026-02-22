package services

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/takumi/fastchem/internal/database"
	"github.com/takumi/fastchem/internal/models"
)

// RankedPlayerState holds per-player state within an active ranked match.
type RankedPlayerState struct {
	UserID         int64
	Username       string
	Rating         int
	Combo          int
	BestCombo      int
	TotalScore     int
	Answered       int
	Correct        int
	HardScore      int
	TotalTime      float64 // total time spent (seconds)
	Connected      bool
	DisconnectedAt *time.Time
	QuestionStart  time.Time // when the current question was issued
	Results        []models.RankedQuestionResult
}

// ActiveRankedMatch holds the full in-memory state of a ranked match.
type ActiveRankedMatch struct {
	mu             sync.Mutex
	MatchID        int64
	Seed           int64
	Player1        *RankedPlayerState
	Player2        *RankedPlayerState
	Questions      *RankedQuestionSet
	Status         string // "waiting", "active", "finished", "abandoned"
	CurrentP1Index int    // current question index for player 1
	CurrentP2Index int    // current question index for player 2
	CreatedAt      time.Time
	FinishedAt     *time.Time

	// Channels to notify WebSocket handlers
	P1Send chan models.WSMessage
	P2Send chan models.WSMessage
}

// RankedMatchService manages active ranked matches.
type RankedMatchService struct {
	mu      sync.RWMutex
	matches map[int64]*ActiveRankedMatch // matchID → match
	rating  *RatingService
}

// NewRankedMatchService creates a new RankedMatchService.
func NewRankedMatchService() *RankedMatchService {
	return &RankedMatchService{
		matches: make(map[int64]*ActiveRankedMatch),
		rating:  NewRatingService(),
	}
}

// CreateMatch creates a new ranked match in the database and in-memory store.
func (rms *RankedMatchService) CreateMatch(
	p1ID int64, p1Username string, p1Rating int,
	p2ID int64, p2Username string, p2Rating int,
	seed int64,
) (*ActiveRankedMatch, error) {
	// Insert into database
	result, err := database.DB.Exec(
		`INSERT INTO ranked_matches (player1_id, player2_id, seed, player1_score, player2_score, status)
		 VALUES (?, ?, ?, 0, 0, 'active')`,
		p1ID, p2ID, seed,
	)
	if err != nil {
		return nil, err
	}
	matchID, _ := result.LastInsertId()

	// Generate questions deterministically
	questions := GenerateRankedQuestions(seed)

	match := &ActiveRankedMatch{
		MatchID: matchID,
		Seed:    seed,
		Player1: &RankedPlayerState{
			UserID:    p1ID,
			Username:  p1Username,
			Rating:    p1Rating,
			Connected: true,
			Results:   make([]models.RankedQuestionResult, 0, models.RankedQuestionsPerMatch),
		},
		Player2: &RankedPlayerState{
			UserID:    p2ID,
			Username:  p2Username,
			Rating:    p2Rating,
			Connected: true,
			Results:   make([]models.RankedQuestionResult, 0, models.RankedQuestionsPerMatch),
		},
		Questions:      questions,
		Status:         "active",
		CurrentP1Index: 0,
		CurrentP2Index: 0,
		CreatedAt:      time.Now(),
		P1Send:         make(chan models.WSMessage, 32),
		P2Send:         make(chan models.WSMessage, 32),
	}

	rms.mu.Lock()
	rms.matches[matchID] = match
	rms.mu.Unlock()

	return match, nil
}

// GetMatch retrieves an active ranked match by ID.
func (rms *RankedMatchService) GetMatch(matchID int64) (*ActiveRankedMatch, bool) {
	rms.mu.RLock()
	defer rms.mu.RUnlock()
	m, ok := rms.matches[matchID]
	return m, ok
}

// GetMatchByUserID finds the active ranked match for a given user.
func (rms *RankedMatchService) GetMatchByUserID(userID int64) (*ActiveRankedMatch, bool) {
	rms.mu.RLock()
	defer rms.mu.RUnlock()
	for _, m := range rms.matches {
		if m.Player1.UserID == userID || m.Player2.UserID == userID {
			return m, true
		}
	}
	return nil, false
}

// RemoveMatch removes a match from the in-memory store.
func (rms *RankedMatchService) RemoveMatch(matchID int64) {
	rms.mu.Lock()
	defer rms.mu.Unlock()
	if m, ok := rms.matches[matchID]; ok {
		close(m.P1Send)
		close(m.P2Send)
		delete(rms.matches, matchID)
	}
}

// SubmitAnswer processes an answer from a player in a ranked match.
// Returns the answer result and whether the match is complete for this player.
func (rms *RankedMatchService) SubmitAnswer(
	matchID int64, userID int64, questionIndex int, selectedIndex int,
) (*models.AnswerResultPayload, bool, error) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return nil, false, nil
	}

	match.mu.Lock()
	defer match.mu.Unlock()

	if match.Status != "active" {
		return nil, false, nil
	}

	// Determine which player
	var player *RankedPlayerState
	var opponent *RankedPlayerState
	var currentIndex *int

	if userID == match.Player1.UserID {
		player = match.Player1
		opponent = match.Player2
		currentIndex = &match.CurrentP1Index
	} else if userID == match.Player2.UserID {
		player = match.Player2
		opponent = match.Player1
		currentIndex = &match.CurrentP2Index
	} else {
		return nil, false, nil
	}

	// Validate question index
	if questionIndex != *currentIndex {
		return nil, false, nil
	}
	if questionIndex >= len(match.Questions.Questions) {
		return nil, false, nil
	}

	q := match.Questions.Questions[questionIndex]
	timeSpent := time.Since(player.QuestionStart).Seconds()

	cfg := GetDifficultyConfig(q.Difficulty)
	timedOut := timeSpent >= float64(cfg.TimeLimit)
	correct := selectedIndex == q.CorrectIndex && !timedOut

	// Update combo
	if correct {
		player.Combo++
		if player.Combo > player.BestCombo {
			player.BestCombo = player.Combo
		}
	} else {
		player.Combo = 0
	}

	// Calculate score
	finalScore, speedBonus, comboMult := CalculateRankedScore(q.Difficulty, timeSpent, correct, player.Combo)

	player.TotalScore += finalScore
	player.Answered++
	player.TotalTime += timeSpent
	if correct {
		player.Correct++
	}
	if q.Difficulty == "hard" && correct {
		player.HardScore += finalScore
	}

	// Record result
	qResult := models.RankedQuestionResult{
		RankedMatchID: matchID,
		UserID:        userID,
		QuestionIndex: questionIndex,
		Difficulty:    q.Difficulty,
		Topic:         q.Category,
		Correct:       correct,
		TimeSpent:     timeSpent,
		ScoreEarned:   finalScore,
	}
	player.Results = append(player.Results, qResult)

	*currentIndex++

	resp := &models.AnswerResultPayload{
		Index:            questionIndex,
		Correct:          correct,
		CorrectIndex:     q.CorrectIndex,
		ScoreEarned:      finalScore,
		SpeedBonus:       speedBonus,
		ComboMultiplier:  comboMult,
		Combo:            player.Combo,
		TotalScore:       player.TotalScore,
		TimedOut:         timedOut,
		TimeSpent:        timeSpent,
		OpponentAnswered: opponent.Answered,
		OpponentScore:    opponent.TotalScore,
	}

	playerDone := player.Answered >= models.RankedQuestionsPerMatch
	bothDone := match.Player1.Answered >= models.RankedQuestionsPerMatch &&
		match.Player2.Answered >= models.RankedQuestionsPerMatch

	if bothDone {
		rms.finalizeMatchLocked(match)
	}

	return resp, playerDone, nil
}

// StartQuestion marks the beginning of a question for timing purposes.
func (rms *RankedMatchService) StartQuestion(matchID int64, userID int64) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	if userID == match.Player1.UserID {
		match.Player1.QuestionStart = time.Now()
	} else if userID == match.Player2.UserID {
		match.Player2.QuestionStart = time.Now()
	}
}

// GetCurrentQuestion returns the current question for a player.
func (rms *RankedMatchService) GetCurrentQuestion(matchID int64, userID int64) (*RankedQuestion, int, bool) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return nil, 0, false
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	var idx int
	if userID == match.Player1.UserID {
		idx = match.CurrentP1Index
	} else if userID == match.Player2.UserID {
		idx = match.CurrentP2Index
	} else {
		return nil, 0, false
	}

	if idx >= len(match.Questions.Questions) {
		return nil, idx, false
	}

	return &match.Questions.Questions[idx], idx, true
}

// HandleDisconnect marks a player as disconnected.
func (rms *RankedMatchService) HandleDisconnect(matchID int64, userID int64) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	now := time.Now()
	if userID == match.Player1.UserID {
		match.Player1.Connected = false
		match.Player1.DisconnectedAt = &now
	} else if userID == match.Player2.UserID {
		match.Player2.Connected = false
		match.Player2.DisconnectedAt = &now
	}
}

// HandleReconnect marks a player as reconnected.
func (rms *RankedMatchService) HandleReconnect(matchID int64, userID int64) bool {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return false
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	if userID == match.Player1.UserID {
		if match.Player1.DisconnectedAt != nil && time.Since(*match.Player1.DisconnectedAt) > 10*time.Second {
			return false // too late
		}
		match.Player1.Connected = true
		match.Player1.DisconnectedAt = nil
		return true
	} else if userID == match.Player2.UserID {
		if match.Player2.DisconnectedAt != nil && time.Since(*match.Player2.DisconnectedAt) > 10*time.Second {
			return false
		}
		match.Player2.Connected = true
		match.Player2.DisconnectedAt = nil
		return true
	}
	return false
}

// CheckDisconnectTimeout checks if a disconnected player has exceeded the 10-second window.
// If so, the opponent wins by default.
func (rms *RankedMatchService) CheckDisconnectTimeout(matchID int64) (timedOut bool, loserID int64) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return false, 0
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	if match.Status != "active" {
		return false, 0
	}

	if match.Player1.DisconnectedAt != nil && time.Since(*match.Player1.DisconnectedAt) > 10*time.Second {
		return true, match.Player1.UserID
	}
	if match.Player2.DisconnectedAt != nil && time.Since(*match.Player2.DisconnectedAt) > 10*time.Second {
		return true, match.Player2.UserID
	}
	return false, 0
}

// ForceWin ends a match with one player winning due to opponent disconnect.
func (rms *RankedMatchService) ForceWin(matchID int64, winnerID int64) {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	match.Status = "abandoned"
	now := time.Now()
	match.FinishedAt = &now

	rms.persistMatchResult(match, winnerID)
}

// FinalizeIfBothDone checks if both players have finished and finalizes the match.
func (rms *RankedMatchService) FinalizeIfBothDone(matchID int64) bool {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return false
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	if match.Status != "active" {
		return match.Status == "finished" || match.Status == "abandoned"
	}

	if match.Player1.Answered >= models.RankedQuestionsPerMatch &&
		match.Player2.Answered >= models.RankedQuestionsPerMatch {
		rms.finalizeMatchLocked(match)
		return true
	}
	return false
}

// finalizeMatchLocked finalizes a match (must be called with match.mu held).
func (rms *RankedMatchService) finalizeMatchLocked(match *ActiveRankedMatch) {
	if match.Status != "active" {
		return
	}

	match.Status = "finished"
	now := time.Now()
	match.FinishedAt = &now

	winnerID := rms.determineWinner(match)
	rms.persistMatchResult(match, winnerID)
}

// determineWinner decides the winner using the tiebreaker rules.
func (rms *RankedMatchService) determineWinner(match *ActiveRankedMatch) int64 {
	p1 := match.Player1
	p2 := match.Player2

	// Primary: highest total score
	if p1.TotalScore > p2.TotalScore {
		return p1.UserID
	}
	if p2.TotalScore > p1.TotalScore {
		return p2.UserID
	}

	// Tie-breaker 1: higher hard-question total
	if p1.HardScore > p2.HardScore {
		return p1.UserID
	}
	if p2.HardScore > p1.HardScore {
		return p2.UserID
	}

	// Tie-breaker 2: lower total response time
	if p1.TotalTime < p2.TotalTime {
		return p1.UserID
	}
	if p2.TotalTime < p1.TotalTime {
		return p2.UserID
	}

	// Ultimate fallback: player 1 wins (should be extremely rare)
	return p1.UserID
}

// persistMatchResult saves the match result and updates ratings in the database.
func (rms *RankedMatchService) persistMatchResult(match *ActiveRankedMatch, winnerID int64) {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("ranked: failed to start tx: %v", err)
		return
	}
	defer tx.Rollback()

	// Update ranked_matches
	_, err = tx.Exec(
		`UPDATE ranked_matches 
		 SET player1_score = ?, player2_score = ?, winner_id = ?, status = ?, finished_at = ?
		 WHERE id = ?`,
		match.Player1.TotalScore, match.Player2.TotalScore,
		winnerID, match.Status, match.FinishedAt, match.MatchID,
	)
	if err != nil {
		log.Printf("ranked: failed to update match: %v", err)
		return
	}

	// Insert question results for both players
	for _, results := range [][]models.RankedQuestionResult{match.Player1.Results, match.Player2.Results} {
		for _, r := range results {
			_, err = tx.Exec(
				`INSERT INTO ranked_question_results 
				 (ranked_match_id, user_id, question_index, difficulty, topic, correct, time_spent, score_earned)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				r.RankedMatchID, r.UserID, r.QuestionIndex, r.Difficulty,
				r.Topic, r.Correct, r.TimeSpent, r.ScoreEarned,
			)
			if err != nil {
				log.Printf("ranked: failed to insert question result: %v", err)
				return
			}
		}
	}

	// Calculate and update ELO ratings
	p1Won := winnerID == match.Player1.UserID
	newP1Rating := rms.rating.CalculateNewRating(match.Player1.Rating, match.Player2.Rating, p1Won)
	newP2Rating := rms.rating.CalculateNewRating(match.Player2.Rating, match.Player1.Rating, !p1Won)

	// Update player 1
	if p1Won {
		_, err = tx.Exec(
			`UPDATE users SET rating = ?, ranked_wins = ranked_wins + 1, 
			 highest_rating = MAX(highest_rating, ?) WHERE id = ?`,
			newP1Rating, newP1Rating, match.Player1.UserID,
		)
	} else {
		_, err = tx.Exec(
			`UPDATE users SET rating = ?, ranked_losses = ranked_losses + 1 WHERE id = ?`,
			newP1Rating, match.Player1.UserID,
		)
	}
	if err != nil {
		log.Printf("ranked: failed to update player 1 rating: %v", err)
		return
	}

	// Update player 2
	if !p1Won {
		_, err = tx.Exec(
			`UPDATE users SET rating = ?, ranked_wins = ranked_wins + 1,
			 highest_rating = MAX(highest_rating, ?) WHERE id = ?`,
			newP2Rating, newP2Rating, match.Player2.UserID,
		)
	} else {
		_, err = tx.Exec(
			`UPDATE users SET rating = ?, ranked_losses = ranked_losses + 1 WHERE id = ?`,
			newP2Rating, match.Player2.UserID,
		)
	}
	if err != nil {
		log.Printf("ranked: failed to update player 2 rating: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("ranked: failed to commit: %v", err)
		return
	}

	// Update in-memory ratings for the match end payload
	match.Player1.Rating = newP1Rating
	match.Player2.Rating = newP2Rating
}

// GetMatchEndPayload creates a MatchEndPayload for a specific player.
func (rms *RankedMatchService) GetMatchEndPayload(matchID int64, userID int64) *models.MatchEndPayload {
	match, ok := rms.GetMatch(matchID)
	if !ok {
		return nil
	}
	match.mu.Lock()
	defer match.mu.Unlock()

	var winnerUsername string
	var winnerID int64

	// Determine winner
	if match.Player1.TotalScore >= match.Player2.TotalScore {
		winnerID = rms.determineWinner(match)
	} else {
		winnerID = rms.determineWinner(match)
	}
	if winnerID == match.Player1.UserID {
		winnerUsername = match.Player1.Username
	} else {
		winnerUsername = match.Player2.Username
	}

	// Calculate rating change for the requesting player
	var ratingChange, newRating int
	if userID == match.Player1.UserID {
		newRating = match.Player1.Rating
		// original rating was stored before update; we recalculate
		ratingChange = match.Player1.Rating - getRatingFromDB(match.Player1.UserID)
		if ratingChange == 0 {
			// Already updated in DB, use the stored value
			if winnerID == match.Player1.UserID {
				ratingChange = int(float64(models.EloKFactor) * (1.0 - rms.rating.ExpectedScore(match.Player1.Rating, match.Player2.Rating)))
			} else {
				ratingChange = -int(float64(models.EloKFactor) * rms.rating.ExpectedScore(match.Player1.Rating, match.Player2.Rating))
			}
		}
	} else {
		newRating = match.Player2.Rating
		ratingChange = match.Player2.Rating - getRatingFromDB(match.Player2.UserID)
		if ratingChange == 0 {
			if winnerID == match.Player2.UserID {
				ratingChange = int(float64(models.EloKFactor) * (1.0 - rms.rating.ExpectedScore(match.Player2.Rating, match.Player1.Rating)))
			} else {
				ratingChange = -int(float64(models.EloKFactor) * rms.rating.ExpectedScore(match.Player2.Rating, match.Player1.Rating))
			}
		}
	}

	return &models.MatchEndPayload{
		MatchID:  matchID,
		Winner:   winnerUsername,
		WinnerID: winnerID,
		Player1: models.MatchEndPlayerSummary{
			UserID:         match.Player1.UserID,
			Username:       match.Player1.Username,
			TotalScore:     match.Player1.TotalScore,
			CorrectAnswers: match.Player1.Correct,
			HardScore:      match.Player1.HardScore,
			TotalTime:      match.Player1.TotalTime,
			Combo:          match.Player1.BestCombo,
		},
		Player2: models.MatchEndPlayerSummary{
			UserID:         match.Player2.UserID,
			Username:       match.Player2.Username,
			TotalScore:     match.Player2.TotalScore,
			CorrectAnswers: match.Player2.Correct,
			HardScore:      match.Player2.HardScore,
			TotalTime:      match.Player2.TotalTime,
			Combo:          match.Player2.BestCombo,
		},
		RatingChange: ratingChange,
		NewRating:    newRating,
	}
}

// getRatingFromDB fetches the current rating from the database.
func getRatingFromDB(userID int64) int {
	var rating int
	err := database.DB.QueryRow("SELECT rating FROM users WHERE id = ?", userID).Scan(&rating)
	if err != nil {
		return models.DefaultRating
	}
	return rating
}

// CleanupStaleMatches removes matches that have been active too long.
func (rms *RankedMatchService) CleanupStaleMatches(maxAge time.Duration) {
	rms.mu.Lock()
	defer rms.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	for id, m := range rms.matches {
		if m.CreatedAt.Before(cutoff) && m.Status == "active" {
			m.Status = "abandoned"
			now := time.Now()
			m.FinishedAt = &now
			database.DB.Exec(
				"UPDATE ranked_matches SET status = 'abandoned', finished_at = ? WHERE id = ?",
				now, id,
			)
			close(m.P1Send)
			close(m.P2Send)
			delete(rms.matches, id)
		}
	}
}

// GetUserRankedStats retrieves the ranked stats for a user.
func GetUserRankedStats(userID int64) (rating, wins, losses, highestRating int, err error) {
	err = database.DB.QueryRow(
		"SELECT rating, ranked_wins, ranked_losses, highest_rating FROM users WHERE id = ?",
		userID,
	).Scan(&rating, &wins, &losses, &highestRating)
	if err == sql.ErrNoRows {
		return models.DefaultRating, 0, 0, models.DefaultRating, nil
	}
	return
}
