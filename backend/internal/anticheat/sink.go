package anticheat

import (
	"log/slog"
	"time"
)

// Event is one finding, ready to be recorded. It carries the measurements that
// produced it so a reviewer can judge the finding without rerunning the rule.
type Event struct {
	At         time.Time
	Subject    string
	UserID     int64
	MatchID    int64
	Mode       Mode
	QuestionID string
	Difficulty string
	TimeSpent  float64
	Rule       string
	Detail     string
	Action     Action
}

// Sink receives findings. It is the seam for persistence: the engine knows
// nothing about where findings go, so adding a cheat_flags table later means
// writing one Sink and passing it to NewEngine, with no change to any
// detector or call site.
//
// Record runs on the answer path. An implementation that talks to a database
// should buffer and write asynchronously rather than block a player's answer.
type Sink interface {
	Record(Event)
}

// SlogSink writes findings to the structured log. Findings are warnings, not
// errors: the server handled the answer fine, a human should look at it.
type SlogSink struct{}

func (SlogSink) Record(e Event) {
	slog.Warn("anticheat finding",
		"rule", e.Rule,
		"action", string(e.Action),
		"detail", e.Detail,
		"subject", e.Subject,
		"user_id", e.UserID,
		"match_id", e.MatchID,
		"mode", string(e.Mode),
		"question_id", e.QuestionID,
		"difficulty", e.Difficulty,
		"time_spent", e.TimeSpent,
	)
}

// DiscardSink throws findings away. Used when no sink is configured and in
// tests that only care about the verdict.
type DiscardSink struct{}

func (DiscardSink) Record(Event) {}

// MultiSink fans one finding out to several sinks — log it and store it, say.
type MultiSink []Sink

func (m MultiSink) Record(e Event) {
	for _, s := range m {
		if s != nil {
			s.Record(e)
		}
	}
}
