package anticheat

import (
	"database/sql"
	"log/slog"
	"sync/atomic"
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

// DBSink persists findings to the anticheat_findings table, which is what the
// admin console reads.
//
// Record is called on the answer path with a match lock held, so it never
// touches SQLite inline: events go onto a buffered channel that one writer
// goroutine drains. When the buffer is full the event is dropped and counted —
// losing a finding is strictly better than stalling every match behind the
// single SQLite writer.
type DBSink struct {
	db      *sql.DB
	events  chan Event
	dropped atomic.Int64
}

// NewDBSink starts the writer goroutine. buffer is the number of findings that
// may be in flight before drops begin; 0 selects a sane default.
func NewDBSink(db *sql.DB, buffer int) *DBSink {
	if buffer <= 0 {
		buffer = 256
	}
	s := &DBSink{db: db, events: make(chan Event, buffer)}
	go s.write()
	return s
}

func (s *DBSink) Record(e Event) {
	select {
	case s.events <- e:
	default:
		s.dropped.Add(1)
	}
}

// Dropped reports how many findings were discarded because the buffer was
// full. The admin console surfaces it: a non-zero value means the findings
// table is an undercount.
func (s *DBSink) Dropped() int64 { return s.dropped.Load() }

func (s *DBSink) write() {
	for e := range s.events {
		_, err := s.db.Exec(
			`INSERT INTO anticheat_findings
			 (at, subject, user_id, match_id, mode, question_id, difficulty, time_spent, rule, detail, action)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			// Written as UTC text in SQLite's own format, matching what
			// CURRENT_TIMESTAMP produces elsewhere in the schema. A bound
			// time.Time would carry a zone offset, and the string comparisons
			// behind date(at) and at >= datetime('now', …) would then be wrong
			// for every row written outside UTC.
			e.At.UTC().Format(time.DateTime),
			e.Subject, e.UserID, e.MatchID, string(e.Mode), e.QuestionID, e.Difficulty, e.TimeSpent,
			e.Rule, e.Detail, string(e.Action),
		)
		if err != nil {
			slog.Error("anticheat: storing finding failed", "rule", e.Rule, "error", err)
		}
	}
}

// MultiSink fans one finding out to several sinks — log it and store it, say.
type MultiSink []Sink

func (m MultiSink) Record(e Event) {
	for _, s := range m {
		if s != nil {
			s.Record(e)
		}
	}
}
