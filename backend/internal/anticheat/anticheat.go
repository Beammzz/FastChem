// Package anticheat inspects answers that have already been scored and reports
// the ones a human could not plausibly have produced.
//
// It is deliberately advisory. Every rule ships with action "observe", so
// enabling the package changes no gameplay — it only produces findings. An
// operator escalates a rule to "reject" by editing its row in
// anticheat_rules, and the call sites honour that without a code change.
//
// The package holds no opinion about chemistry, scoring, or matches: callers
// hand it a finished Signal and act on the Verdict. It imports nothing from
// services or handlers, which is what lets both of them call it.
package anticheat

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Mode is the kind of play an answer came from. Detectors may treat modes
// differently; the sink records it so findings can be filtered by stakes.
type Mode string

const (
	ModeCasual Mode = "casual" // /api/answer, no account
	ModeMatch  Mode = "match"  // /api/match/answer, single player, scored server-side
	ModeRanked Mode = "ranked" // ranked 1v1, ELO at stake
	ModeRoom   Mode = "room"   // custom room, no ELO
)

// Action is what the server does about a flagged answer. Actions are ordered:
// a verdict takes the strongest action any of its findings asked for.
type Action string

const (
	// ActionObserve records the finding and changes nothing. The default, and
	// the only action any rule ships with.
	ActionObserve Action = "observe"
	// ActionReject marks the answer wrong and awards no points. The match
	// continues; the player is not told why, because telling a cheater which
	// rule caught them is telling them how to stay under it.
	ActionReject Action = "reject"
)

// severity orders actions so a verdict can pick the strongest one.
func (a Action) severity() int {
	switch a {
	case ActionReject:
		return 2
	default:
		return 1
	}
}

// ParseAction maps a stored string onto a known action, falling back to
// observe. An unrecognised value in the database must never escalate.
func ParseAction(s string) (Action, bool) {
	switch Action(s) {
	case ActionObserve:
		return ActionObserve, true
	case ActionReject:
		return ActionReject, true
	default:
		return ActionObserve, false
	}
}

// Signal is one answer, as the server measured it. Every field is server-side
// truth: TimeSpent comes from when the server issued the question, never from
// the client, and Correct has already been resolved against the stored answer.
type Signal struct {
	// Subject identifies who is answering, for history that spans questions.
	// Use SubjectUser for anyone signed in and SubjectIP for casual play.
	Subject    string
	UserID     int64 // 0 when anonymous
	MatchID    int64 // 0 outside a match
	Mode       Mode
	QuestionID string
	Difficulty string
	TimeSpent  float64 // seconds, measured server-side
	TimeLimit  int     // seconds allowed for this difficulty
	Correct    bool
	TimedOut   bool
}

// SubjectUser keys history by account.
func SubjectUser(userID int64) string { return fmt.Sprintf("user:%d", userID) }

// SubjectIP keys history for casual play, which has no account. It is a weak
// key — a shared NAT groups strangers together — which is one reason casual
// findings are worth less than ranked ones.
func SubjectIP(ip string) string { return "ip:" + ip }

// Finding is one rule's complaint about one answer.
type Finding struct {
	Rule   string // the detector's name, and its settings key
	Detail string // human-readable, with the numbers that tripped it
	Action Action // what the rule was configured to do
}

// Verdict is the result of running every enabled rule over one answer.
type Verdict struct {
	Findings []Finding
	Action   Action // the strongest action among Findings
}

// Flagged reports whether any rule complained, whatever the action.
func (v Verdict) Flagged() bool { return len(v.Findings) > 0 }

// Rejected reports whether the answer should score nothing.
func (v Verdict) Rejected() bool { return v.Action == ActionReject }

// Rules lists the rules that fired, for logging.
func (v Verdict) Rules() []string {
	out := make([]string, 0, len(v.Findings))
	for _, f := range v.Findings {
		out = append(out, f.Rule)
	}
	return out
}

// Engine runs the detectors against a signal using the current settings.
// It is safe for concurrent use; evaluation touches memory only, so it is
// cheap enough to call while holding a match lock.
type Engine struct {
	settings  *Store
	sink      Sink
	history   *tracker
	detectors []Detector
}

// NewEngine builds an engine with the built-in detectors registered. A nil
// store means "compiled defaults", and a nil sink discards findings.
func NewEngine(settings *Store, sink Sink) *Engine {
	if settings == nil {
		settings = NewStore(nil)
	}
	if sink == nil {
		sink = DiscardSink{}
	}
	return &Engine{
		settings:  settings,
		sink:      sink,
		history:   newTracker(),
		detectors: builtinDetectors(),
	}
}

// Use appends a detector. Registration is append-only and order is stable, so
// findings arrive in a predictable order. A detector with no settings row runs
// on its own defaults.
func (e *Engine) Use(d Detector) { e.detectors = append(e.detectors, d) }

// Evaluate runs every enabled rule over one answer and records what fires.
// A zero Verdict — no findings, observe — is the normal result.
func (e *Engine) Evaluate(sig Signal) Verdict {
	v := Verdict{Action: ActionObserve}
	if e == nil {
		return v
	}

	// The answer under inspection is the last entry of its own history, so a
	// windowed rule can look at the run it completes without special-casing it.
	hist := e.history.record(sig)

	for _, d := range e.detectors {
		rule := e.settings.Rule(d.Name())
		if !rule.Enabled {
			continue
		}
		f := d.Inspect(sig, hist, rule.Params)
		if f == nil {
			continue
		}
		f.Rule = d.Name()
		f.Action = rule.Action
		v.Findings = append(v.Findings, *f)
		if f.Action.severity() > v.Action.severity() {
			v.Action = f.Action
		}
	}

	now := time.Now()
	for _, f := range v.Findings {
		e.sink.Record(Event{
			At:         now,
			Subject:    sig.Subject,
			UserID:     sig.UserID,
			MatchID:    sig.MatchID,
			Mode:       sig.Mode,
			QuestionID: sig.QuestionID,
			Difficulty: sig.Difficulty,
			TimeSpent:  sig.TimeSpent,
			Rule:       f.Rule,
			Detail:     f.Detail,
			Action:     f.Action,
		})
	}
	return v
}

// Settings exposes the store so callers can reload it on a schedule.
func (e *Engine) Settings() *Store {
	if e == nil {
		return nil
	}
	return e.settings
}

// Sweep drops history for subjects that have not answered in maxIdle, the way
// the question and match stores drop their own stale entries.
func (e *Engine) Sweep(maxIdle time.Duration) {
	if e == nil {
		return
	}
	e.history.sweep(maxIdle)
}

// ─── Process-wide default ───────────────────────────────────────
//
// services and handlers both submit answers, and neither owns the engine, so
// it lives here as a package default alongside GlobalQuestionStore and
// GlobalMatchStore. Until Init is called every Evaluate is a no-op, which
// keeps tests and any caller that skips setup working.

var defaultEngine atomic.Pointer[Engine]

// Init installs the process-wide engine. Safe to call once at startup.
func Init(e *Engine) { defaultEngine.Store(e) }

// Default returns the process-wide engine, or nil when Init has not run.
func Default() *Engine { return defaultEngine.Load() }

// Evaluate runs the process-wide engine. It returns an empty verdict when no
// engine is installed, so a caller never has to nil-check.
func Evaluate(sig Signal) Verdict { return Default().Evaluate(sig) }

// Sweep prunes the process-wide engine's history.
func Sweep(maxIdle time.Duration) { Default().Sweep(maxIdle) }
