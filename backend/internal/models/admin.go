package models

import "time"

// Shapes served under /api/admin. They are read models for the admin console
// and nothing else reads them, so they aggregate freely across tables — but
// they are still the JSON contract, mirrored in frontend/src/types/index.ts.

// AdminBucket is one labelled tally. Charts in the console are built from
// these, so a new breakdown needs no new type.
type AdminBucket struct {
	Label string  `json:"label"`
	Count int     `json:"count"`
	Value float64 `json:"value"` // rate, average, or share — meaning depends on the series
}

// AdminOverview is the headline panel: one number per thing worth watching.
type AdminOverview struct {
	Users             int     `json:"users"`
	Admins            int     `json:"admins"`
	NewUsers7d        int     `json:"newUsers7d"`
	ActivePlayers7d   int     `json:"activePlayers7d"`
	Games             int     `json:"games"`
	Games24h          int     `json:"games24h"`
	SoloMatches       int     `json:"soloMatches"`
	RankedMatches     int     `json:"rankedMatches"`
	QuestionsAnswered int     `json:"questionsAnswered"`
	TotalPoints       int64   `json:"totalPoints"`
	AverageAccuracy   float64 `json:"averageAccuracy"`
	Findings          int     `json:"findings"`
	Findings24h       int     `json:"findings24h"`
	FindingsDropped   int64   `json:"findingsDropped"`
	RejectingRules    int     `json:"rejectingRules"`
	QueueSize         int     `json:"queueSize"`
	LiveMatches       int     `json:"liveMatches"`
	OpenRooms         int     `json:"openRooms"`
	LiveSoloMatches   int     `json:"liveSoloMatches"`
	StoredQuestions   int     `json:"storedQuestions"`
	Topics            int     `json:"topics"`
	UptimeSeconds     float64 `json:"uptimeSeconds"`
	DatabaseBytes     int64   `json:"databaseBytes"`
	Goroutines        int     `json:"goroutines"`
}

// AdminDailyPoint is one day of activity.
type AdminDailyPoint struct {
	Date          string `json:"date"` // YYYY-MM-DD, UTC
	Games         int    `json:"games"`
	Players       int    `json:"players"`
	Points        int    `json:"points"`
	RankedMatches int    `json:"rankedMatches"`
	Findings      int    `json:"findings"`
}

// AdminTopicStat is how one question topic is performing, measured from ranked
// results because those carry the topic id per question.
type AdminTopicStat struct {
	Topic       string  `json:"topic"`
	Difficulty  string  `json:"difficulty"`
	Attempts    int     `json:"attempts"`
	Correct     int     `json:"correct"`
	Accuracy    float64 `json:"accuracy"`    // percent
	AverageTime float64 `json:"averageTime"` // seconds
}

// AdminAnalytics is every chart on the analytics tab.
type AdminAnalytics struct {
	Daily      []AdminDailyPoint `json:"daily"`
	Difficulty []AdminBucket     `json:"difficulty"` // games and mean accuracy per difficulty
	Topics     []AdminTopicStat  `json:"topics"`
	Hours      []AdminBucket     `json:"hours"`   // games per hour of day, UTC
	Ratings    []AdminBucket     `json:"ratings"` // player count per rating band
	Rules      []AdminBucket     `json:"rules"`   // findings per rule
}

// AdminUserRow is one row of the user table, account plus play record.
type AdminUserRow struct {
	ID            int64      `json:"id"`
	Username      string     `json:"username"`
	IsAdmin       bool       `json:"isAdmin"`
	TotalPoints   int        `json:"totalPoints"`
	Rating        int        `json:"rating"`
	RankedWins    int        `json:"rankedWins"`
	RankedLosses  int        `json:"rankedLosses"`
	HighestRating int        `json:"highestRating"`
	Games         int        `json:"games"`
	HighScore     int        `json:"highScore"`
	TotalAnswered int        `json:"totalAnswered"`
	TotalCorrect  int        `json:"totalCorrect"`
	Accuracy      float64    `json:"accuracy"`
	Findings      int        `json:"findings"`
	CreatedAt     time.Time  `json:"createdAt"`
	LastPlayed    *time.Time `json:"lastPlayed"`
	Online        bool       `json:"online"` // in queue or in a live match right now
}

// AdminUsersResponse is one page of the user table.
type AdminUsersResponse struct {
	Users    []AdminUserRow `json:"users"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

// AdminUserUpdateRequest edits one account. Every editable field is a pointer:
// an omitted field is left alone, which is what lets the console send a single
// toggle without echoing the whole row back.
type AdminUserUpdateRequest struct {
	UserID      int64   `json:"userId" binding:"required"`
	IsAdmin     *bool   `json:"isAdmin"`
	Rating      *int    `json:"rating"`
	TotalPoints *int    `json:"totalPoints"`
	Password    *string `json:"password"`
}

// AdminUserIDRequest names one account, for actions that take no other input.
type AdminUserIDRequest struct {
	UserID int64 `json:"userId" binding:"required"`
}

// AdminFinding is one stored anti-cheat detection, resolved to a username.
type AdminFinding struct {
	ID         int64     `json:"id"`
	At         time.Time `json:"at"`
	Subject    string    `json:"subject"`
	UserID     int64     `json:"userId"`
	Username   string    `json:"username"` // empty for casual play, which has no account
	MatchID    int64     `json:"matchId"`
	Mode       string    `json:"mode"`
	QuestionID string    `json:"questionId"`
	Difficulty string    `json:"difficulty"`
	TimeSpent  float64   `json:"timeSpent"`
	Rule       string    `json:"rule"`
	Detail     string    `json:"detail"`
	Action     string    `json:"action"`
}

// AdminFindingsResponse is one filtered page of findings plus the breakdowns
// of the whole filtered set, so the summary does not shift as pages turn.
type AdminFindingsResponse struct {
	Findings []AdminFinding `json:"findings"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	ByRule   []AdminBucket  `json:"byRule"`
	ByMode   []AdminBucket  `json:"byMode"`
	ByAction []AdminBucket  `json:"byAction"`
	Dropped  int64          `json:"dropped"` // findings the sink discarded under load
}

// AdminRule is one anti-cheat rule as the console shows it: its live settings
// plus how often it has actually fired.
type AdminRule struct {
	Name        string             `json:"name"`
	Enabled     bool               `json:"enabled"`
	Action      string             `json:"action"`
	Params      map[string]float64 `json:"params"`
	Findings    int                `json:"findings"`
	Findings24h int                `json:"findings24h"`
}

// AdminRuleUpdateRequest retunes one rule. Params replace the stored set and
// are merged over the compiled defaults on reload.
type AdminRuleUpdateRequest struct {
	Name    string             `json:"name" binding:"required"`
	Enabled bool               `json:"enabled"`
	Action  string             `json:"action" binding:"required"`
	Params  map[string]float64 `json:"params"`
}

// AdminQueueEntry is one player waiting for a ranked match.
type AdminQueueEntry struct {
	UserID         int64   `json:"userId"`
	Username       string  `json:"username"`
	Rating         int     `json:"rating"`
	WaitingSeconds float64 `json:"waitingSeconds"`
}

// AdminLivePlayer is one side of a match in progress.
type AdminLivePlayer struct {
	UserID     int64  `json:"userId"`
	Username   string `json:"username"`
	Rating     int    `json:"rating"`
	TotalScore int    `json:"totalScore"`
	Answered   int    `json:"answered"`
	Correct    int    `json:"correct"`
	Combo      int    `json:"combo"`
	Connected  bool   `json:"connected"`
}

// AdminLiveMatch is a ranked or room match currently in memory. A negative
// MatchID is a room match, matching the synthetic ids RoomService assigns.
type AdminLiveMatch struct {
	MatchID   int64           `json:"matchId"`
	Status    string          `json:"status"`
	Question  int             `json:"question"` // last question index sent to both sides
	CreatedAt time.Time       `json:"createdAt"`
	Player1   AdminLivePlayer `json:"player1"`
	Player2   AdminLivePlayer `json:"player2"`
}

// AdminRoom is one custom room.
type AdminRoom struct {
	Code      string    `json:"code"`
	HostID    int64     `json:"hostId"`
	HostName  string    `json:"hostName"`
	GuestID   int64     `json:"guestId"`
	GuestName string    `json:"guestName"`
	MatchID   int64     `json:"matchId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// AdminLive is everything the server is holding in memory right now.
type AdminLive struct {
	Queue           []AdminQueueEntry `json:"queue"`
	Matches         []AdminLiveMatch  `json:"matches"`
	Rooms           []AdminRoom       `json:"rooms"`
	SoloMatches     int               `json:"soloMatches"`
	StoredQuestions int               `json:"storedQuestions"`
}

// AdminRankedMatchRow is one finished 1v1.
type AdminRankedMatchRow struct {
	MatchID      int64      `json:"matchId"`
	Player1      string     `json:"player1"`
	Player2      string     `json:"player2"`
	Player1Score int        `json:"player1Score"`
	Player2Score int        `json:"player2Score"`
	Winner       string     `json:"winner"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	FinishedAt   *time.Time `json:"finishedAt"`
}

// AdminSoloMatchRow is one server-scored single-player session.
type AdminSoloMatchRow struct {
	MatchID    int64      `json:"matchId"`
	Username   string     `json:"username"`
	Difficulty string     `json:"difficulty"`
	TotalScore int        `json:"totalScore"`
	BestCombo  int        `json:"bestCombo"`
	Attempts   int        `json:"attempts"`
	StartedAt  time.Time  `json:"startedAt"`
	EndedAt    *time.Time `json:"endedAt"`
}

// AdminScoreRow is one submitted casual score, with its player named.
type AdminScoreRow struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"userId"`
	Username       string    `json:"username"`
	Score          int       `json:"score"`
	TotalAnswered  int       `json:"totalAnswered"`
	CorrectAnswers int       `json:"correctAnswers"`
	Difficulty     string    `json:"difficulty"`
	TimeSpent      float64   `json:"timeSpent"`
	PlayedAt       time.Time `json:"playedAt"`
}

// AdminMatchesResponse is the recent-activity tab.
type AdminMatchesResponse struct {
	Ranked []AdminRankedMatchRow `json:"ranked"`
	Solo   []AdminSoloMatchRow   `json:"solo"`
	Scores []AdminScoreRow       `json:"scores"`
}

// AdminLeaderboards is both boards in one response, read live rather than from
// the public leaderboard's TTL cache.
type AdminLeaderboards struct {
	Casual []LeaderboardEntry       `json:"casual"`
	Ranked []RankedLeaderboardEntry `json:"ranked"`
}

// AdminTopic is one registered question topic.
type AdminTopic struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
}

// AdminQuestionPreview is a sample question with its answer.
//
// Question hides CorrectIndex because it goes out before the player answers;
// this type is a separate shape on an admin-only route that generates a
// throwaway question — it is never stored, never scored, and never reaches a
// player, so revealing the answer here costs nothing.
type AdminQuestionPreview struct {
	Question     string   `json:"question"`
	Choices      []string `json:"choices"`
	CorrectIndex int      `json:"correctIndex"`
	Category     string   `json:"category"`
	Difficulty   string   `json:"difficulty"`
	TimeLimit    int      `json:"timeLimit"`
}
