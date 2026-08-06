package anticheat

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"math"
	"sort"
	"sync"
)

// Params holds a rule's thresholds. Values are float64 because every threshold
// so far is a count, a duration in seconds, or a spread in seconds; a count is
// read back with int(p.Float(...)).
type Params map[string]float64

// Float returns the named threshold, falling back to def when it is missing or
// not a usable number. A detector calls this rather than indexing the map, so
// a typo or a NaN in the database degrades to the compiled default instead of
// silently disabling a rule or flagging everyone.
func (p Params) Float(key string, def float64) float64 {
	if p == nil {
		return def
	}
	v, ok := p[key]
	if !ok || math.IsNaN(v) || math.IsInf(v, 0) {
		return def
	}
	return v
}

// Clone returns a copy, so a caller holding a Rule cannot mutate the store.
func (p Params) Clone() Params {
	if p == nil {
		return nil
	}
	out := make(Params, len(p))
	for k, v := range p {
		out[k] = v
	}
	return out
}

// Rule is one detector's configuration.
type Rule struct {
	Name    string
	Enabled bool
	Action  Action
	Params  Params
}

// defaultRules is the compiled fallback, and what Seed writes into a fresh
// database. Every rule ships enabled and observing: turning the package on
// must not change what any player experiences until an operator says so.
//
// The thresholds are deliberately conservative. A first cut should under-flag;
// a rule that cries wolf gets ignored, and the whole point of starting in
// observe mode is to see real distributions before tightening anything.
var defaultRules = []Rule{
	{
		Name:    "impossible_speed",
		Enabled: true,
		Action:  ActionObserve,
		Params: Params{
			"min_seconds":        1.0,
			"min_seconds_easy":   1.0,
			"min_seconds_medium": 1.5,
			"min_seconds_hard":   2.0,
		},
	},
	{
		Name:    "fast_streak",
		Enabled: true,
		Action:  ActionObserve,
		Params: Params{
			"window":      5,
			"max_seconds": 3.0,
		},
	},
	{
		Name:    "uniform_timing",
		Enabled: true,
		Action:  ActionObserve,
		Params: Params{
			"window":     6,
			"max_stddev": 0.2,
		},
	},
}

// unknownRule is what an unregistered detector gets: it runs, on its own
// defaults, and observes. A detector added in code before its settings row
// exists should work, not vanish.
var unknownRule = Rule{Enabled: true, Action: ActionObserve}

// Store holds the current rule set. Reads happen on every answer, writes only
// on Reload, so it is a plain RWMutex over a map.
type Store struct {
	mu    sync.RWMutex
	db    *sql.DB
	rules map[string]Rule
}

// NewStore returns a store preloaded with the compiled defaults. A nil db is
// valid and means "defaults only" — used in tests and if the settings table is
// unavailable.
func NewStore(db *sql.DB) *Store {
	s := &Store{db: db, rules: make(map[string]Rule, len(defaultRules))}
	for _, r := range defaultRules {
		s.rules[r.Name] = Rule{Name: r.Name, Enabled: r.Enabled, Action: r.Action, Params: r.Params.Clone()}
	}
	return s
}

// Rule returns the configuration for a detector.
func (s *Store) Rule(name string) Rule {
	if s == nil {
		r := unknownRule
		r.Name = name
		return r
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[name]
	if !ok {
		r = unknownRule
		r.Name = name
		return r
	}
	r.Params = r.Params.Clone()
	return r
}

// Snapshot returns every rule, name-sorted, for logging and diagnostics.
func (s *Store) Snapshot() []Rule {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	out := make([]Rule, 0, len(s.rules))
	for _, r := range s.rules {
		r.Params = r.Params.Clone()
		out = append(out, r)
	}
	s.mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Seed writes the compiled defaults into the settings table, leaving any row
// that already exists alone. It runs at startup so an operator can discover
// and edit the rules with plain SQL:
//
//	UPDATE anticheat_rules SET action = 'reject' WHERE name = 'impossible_speed';
//
// The change takes effect on the next Reload — no restart.
func (s *Store) Seed(ctx context.Context) error {
	if s == nil || s.db == nil {
		return nil
	}
	for _, r := range defaultRules {
		params, err := json.Marshal(r.Params)
		if err != nil {
			return err
		}
		enabled := 0
		if r.Enabled {
			enabled = 1
		}
		if _, err := s.db.ExecContext(ctx,
			`INSERT OR IGNORE INTO anticheat_rules (name, enabled, action, params) VALUES (?, ?, ?, ?)`,
			r.Name, enabled, string(r.Action), string(params),
		); err != nil {
			return err
		}
	}
	return nil
}

// Reload replaces the cached rule set from the database. Rules with no row
// keep their compiled defaults, so deleting a row resets it rather than
// disabling the detector.
//
// A malformed row is logged and skipped, never fatal: the answer path runs
// through here, and a bad hand-edit must not take scoring down with it.
func (s *Store) Reload(ctx context.Context) error {
	if s == nil || s.db == nil {
		return nil
	}

	rows, err := s.db.QueryContext(ctx, `SELECT name, enabled, action, params FROM anticheat_rules`)
	if err != nil {
		return err
	}
	defer rows.Close()

	next := make(map[string]Rule, len(defaultRules))
	for _, r := range defaultRules {
		next[r.Name] = Rule{Name: r.Name, Enabled: r.Enabled, Action: r.Action, Params: r.Params.Clone()}
	}

	for rows.Next() {
		var (
			name    string
			enabled int
			action  string
			params  string
		)
		if err := rows.Scan(&name, &enabled, &action, &params); err != nil {
			return err
		}

		rule := next[name]
		rule.Name = name
		rule.Enabled = enabled != 0

		parsed, ok := ParseAction(action)
		if !ok {
			slog.Warn("anticheat: unknown action, observing instead", "rule", name, "action", action)
		}
		rule.Action = parsed

		var p Params
		if err := json.Unmarshal([]byte(params), &p); err != nil {
			slog.Warn("anticheat: unreadable params, using defaults", "rule", name, "error", err)
		} else if p != nil {
			// Merge over the defaults so a row that sets one threshold does not
			// blank the rest.
			merged := rule.Params.Clone()
			if merged == nil {
				merged = make(Params, len(p))
			}
			for k, v := range p {
				merged[k] = v
			}
			rule.Params = merged
		}

		next[name] = rule
	}
	if err := rows.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.rules = next
	s.mu.Unlock()
	return nil
}
