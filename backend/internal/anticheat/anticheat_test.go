package anticheat

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func newTestEngine() *Engine { return NewEngine(NewStore(nil), DiscardSink{}) }

func answer(subject, difficulty string, seconds float64, correct bool) Signal {
	return Signal{
		Subject:    subject,
		Mode:       ModeMatch,
		Difficulty: difficulty,
		TimeSpent:  seconds,
		TimeLimit:  30,
		Correct:    correct,
	}
}

func fired(v Verdict, rule string) bool {
	for _, f := range v.Findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

// ─── The package is inert until an operator says otherwise ──────

// Turning anticheat on must not change what any player experiences. Every
// shipped rule observes; nothing rejects until someone edits the table.
func TestDefaultRulesOnlyObserve(t *testing.T) {
	for _, r := range defaultRules {
		if r.Action != ActionObserve {
			t.Errorf("rule %q ships with action %q; defaults must observe so enabling the package changes no gameplay", r.Name, r.Action)
		}
	}
}

// A flagged answer under the default settings is still scored normally.
func TestObserveFlagsWithoutRejecting(t *testing.T) {
	e := newTestEngine()
	v := e.Evaluate(answer("user:1", "hard", 0.2, true))
	if !v.Flagged() {
		t.Fatal("a correct hard answer in 0.2s should be flagged")
	}
	if v.Rejected() {
		t.Fatal("observe must not reject")
	}
}

// Every rule name must have a settings row, or an operator cannot tune it.
func TestEveryDetectorHasDefaultSettings(t *testing.T) {
	byName := make(map[string]bool, len(defaultRules))
	for _, r := range defaultRules {
		byName[r.Name] = true
	}
	for _, d := range builtinDetectors() {
		if !byName[d.Name()] {
			t.Errorf("detector %q has no row in defaultRules; it would be untunable", d.Name())
		}
	}
}

// ─── impossible_speed ───────────────────────────────────────────

func TestImpossibleSpeedIgnoresWrongAnswers(t *testing.T) {
	e := newTestEngine()
	// A fast wrong answer is a misclick. Flagging those buries the signal.
	if v := e.Evaluate(answer("user:2", "hard", 0.1, false)); fired(v, "impossible_speed") {
		t.Fatal("a wrong answer should not trip impossible_speed")
	}
}

func TestImpossibleSpeedUsesThePerDifficultyFloor(t *testing.T) {
	e := newTestEngine()
	// 1.7s is over the easy floor (1.0) and under the hard floor (2.0).
	if v := e.Evaluate(answer("user:3", "easy", 1.7, true)); fired(v, "impossible_speed") {
		t.Fatal("1.7s on an easy question is above the 1.0s floor")
	}
	if v := e.Evaluate(answer("user:4", "hard", 1.7, true)); !fired(v, "impossible_speed") {
		t.Fatal("1.7s on a hard question is below the 2.0s floor")
	}
}

func TestImpossibleSpeedIgnoresTimeouts(t *testing.T) {
	e := newTestEngine()
	sig := answer("user:5", "hard", 0.1, true)
	sig.TimedOut = true
	if v := e.Evaluate(sig); fired(v, "impossible_speed") {
		t.Fatal("a timed-out answer should not trip impossible_speed")
	}
}

// ─── fast_streak ────────────────────────────────────────────────

func TestFastStreakNeedsTheWholeWindow(t *testing.T) {
	e := newTestEngine()
	// Four quick correct answers: one short of the default window of five.
	for i := 0; i < 4; i++ {
		if v := e.Evaluate(answer("user:6", "medium", 2.5, true)); fired(v, "fast_streak") {
			t.Fatalf("fast_streak fired after %d answers, window is 5", i+1)
		}
	}
	if v := e.Evaluate(answer("user:6", "medium", 2.5, true)); !fired(v, "fast_streak") {
		t.Fatal("fast_streak should fire on the fifth quick correct answer")
	}
}

func TestFastStreakBreaksOnAWrongAnswer(t *testing.T) {
	e := newTestEngine()
	for i := 0; i < 4; i++ {
		e.Evaluate(answer("user:7", "medium", 2.5, true))
	}
	e.Evaluate(answer("user:7", "medium", 2.5, false))
	if v := e.Evaluate(answer("user:7", "medium", 2.5, true)); fired(v, "fast_streak") {
		t.Fatal("a wrong answer should break the streak")
	}
}

func TestFastStreakIsPerSubject(t *testing.T) {
	e := newTestEngine()
	// Two players answering quickly must not pool into one streak.
	for i := 0; i < 4; i++ {
		e.Evaluate(answer("user:8", "medium", 2.5, true))
		e.Evaluate(answer("user:9", "medium", 2.5, true))
	}
	if v := e.Evaluate(answer("user:8", "medium", 2.5, true)); !fired(v, "fast_streak") {
		t.Fatal("user:8 reached five quick correct answers")
	}
}

// ─── uniform_timing ─────────────────────────────────────────────

// The false positive this rule most easily produces: a player who walks away
// times out repeatedly, and every timeout lands on the same time limit.
func TestUniformTimingIgnoresTimeouts(t *testing.T) {
	e := newTestEngine()
	for i := 0; i < 10; i++ {
		sig := answer("user:10", "easy", 30.0, false)
		sig.TimedOut = true
		if v := e.Evaluate(sig); fired(v, "uniform_timing") {
			t.Fatal("a run of timeouts is absence, not automation")
		}
	}
}

func TestUniformTimingCatchesMachineCadence(t *testing.T) {
	e := newTestEngine()
	var v Verdict
	for i := 0; i < 6; i++ {
		// Well clear of every speed floor, but with almost no spread.
		v = e.Evaluate(answer("user:11", "hard", 5.0+float64(i)*0.01, i%2 == 0))
	}
	if !fired(v, "uniform_timing") {
		t.Fatal("six answers within 50ms of each other should trip uniform_timing")
	}
}

func TestUniformTimingAcceptsHumanScatter(t *testing.T) {
	e := newTestEngine()
	spread := []float64{4.1, 9.7, 6.2, 12.4, 5.5, 8.8}
	var v Verdict
	for _, s := range spread {
		v = e.Evaluate(answer("user:12", "medium", s, true))
	}
	if fired(v, "uniform_timing") {
		t.Fatal("ordinary variation in answer times should not be flagged")
	}
}

// ─── Escalation ─────────────────────────────────────────────────

// Flipping a rule to reject is the whole point of the settings table: it must
// change the verdict with no code change.
func TestRejectActionRejects(t *testing.T) {
	store := NewStore(nil)
	store.rules["impossible_speed"] = Rule{
		Name:    "impossible_speed",
		Enabled: true,
		Action:  ActionReject,
		Params:  Params{"min_seconds": 1.0},
	}
	e := NewEngine(store, DiscardSink{})

	v := e.Evaluate(answer("user:13", "easy", 0.2, true))
	if !v.Rejected() {
		t.Fatal("a rule set to reject should produce a rejecting verdict")
	}
}

// The strongest action wins when several rules fire at once.
func TestVerdictTakesTheStrongestAction(t *testing.T) {
	store := NewStore(nil)
	r := store.Rule("fast_streak")
	r.Action = ActionReject
	store.rules["fast_streak"] = r
	e := NewEngine(store, DiscardSink{})

	var v Verdict
	for i := 0; i < 5; i++ {
		v = e.Evaluate(answer("user:14", "hard", 0.5, true))
	}
	if !fired(v, "impossible_speed") || !fired(v, "fast_streak") {
		t.Fatalf("expected both rules to fire, got %v", v.Rules())
	}
	if !v.Rejected() {
		t.Fatal("one rejecting rule should make the whole verdict reject")
	}
}

func TestDisabledRuleNeverFires(t *testing.T) {
	store := NewStore(nil)
	r := store.Rule("impossible_speed")
	r.Enabled = false
	store.rules["impossible_speed"] = r
	e := NewEngine(store, DiscardSink{})

	if v := e.Evaluate(answer("user:15", "hard", 0.1, true)); fired(v, "impossible_speed") {
		t.Fatal("a disabled rule must not fire")
	}
}

// ─── Robustness ─────────────────────────────────────────────────

// A caller that never called Init must still work; anticheat is advisory.
func TestNilEngineIsANoOp(t *testing.T) {
	var e *Engine
	v := e.Evaluate(answer("user:16", "hard", 0.01, true))
	if v.Flagged() || v.Rejected() {
		t.Fatal("an uninitialised engine must return an empty verdict")
	}
	e.Sweep(time.Minute) // must not panic
}

func TestParamsFallBackOnGarbage(t *testing.T) {
	var missing Params
	if got := missing.Float("window", 5); got != 5 {
		t.Errorf("nil params should return the default, got %v", got)
	}
	p := Params{"window": 0}
	if got := p.Float("window", 5); got != 0 {
		t.Errorf("an explicit 0 is a real value, got %v", got)
	}
	if got := p.Float("absent", 7); got != 7 {
		t.Errorf("a missing key should return the default, got %v", got)
	}
}

// Detectors must not be handed a slice the tracker keeps mutating.
func TestHistoryIsCopiedPerCall(t *testing.T) {
	tr := newTracker()
	h := tr.record(answer("user:17", "easy", 1.0, true))
	tr.record(answer("user:17", "easy", 2.0, true))
	if len(h.Recent) != 1 {
		t.Fatalf("the first history should still hold one answer, got %d", len(h.Recent))
	}
}

func TestSweepDropsIdleSubjects(t *testing.T) {
	tr := newTracker()
	tr.record(answer("user:18", "easy", 1.0, true))
	tr.sweep(time.Hour)
	if len(tr.subjects) != 1 {
		t.Fatal("a fresh subject should survive a sweep")
	}
	tr.sweep(0)
	if len(tr.subjects) != 0 {
		t.Fatal("an idle subject should be dropped")
	}
}

// Answers arrive from many connections at once. Run with -race.
func TestConcurrentEvaluationIsRaceFree(t *testing.T) {
	e := newTestEngine()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			subject := SubjectUser(int64(n % 4))
			for j := 0; j < 50; j++ {
				e.Evaluate(answer(subject, "medium", float64(j%7)+0.5, j%3 != 0))
			}
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 50; j++ {
			e.Sweep(time.Hour)
			e.Settings().Snapshot()
		}
	}()
	wg.Wait()
}

// ─── Settings, against a real database ──────────────────────────

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(`CREATE TABLE anticheat_rules (
		name TEXT PRIMARY KEY,
		enabled INTEGER NOT NULL DEFAULT 1,
		action TEXT NOT NULL DEFAULT 'observe',
		params TEXT NOT NULL DEFAULT '{}',
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db
}

func TestSeedThenReloadRoundTrips(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db)
	ctx := context.Background()

	if err := store.Seed(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := store.Reload(ctx); err != nil {
		t.Fatalf("reload: %v", err)
	}
	for _, want := range defaultRules {
		got := store.Rule(want.Name)
		if got.Action != want.Action || got.Enabled != want.Enabled {
			t.Errorf("rule %q round-tripped as %+v, want %+v", want.Name, got, want)
		}
		for k, v := range want.Params {
			if got.Params.Float(k, -1) != v {
				t.Errorf("rule %q param %q = %v, want %v", want.Name, k, got.Params.Float(k, -1), v)
			}
		}
	}
}

// Editing a row is how thresholds are retuned without a restart.
func TestReloadPicksUpAnEdit(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db)
	ctx := context.Background()
	if err := store.Seed(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if _, err := db.Exec(`UPDATE anticheat_rules SET action = 'reject', params = '{"min_seconds_hard": 4.0}' WHERE name = 'impossible_speed'`); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := store.Reload(ctx); err != nil {
		t.Fatalf("reload: %v", err)
	}

	got := store.Rule("impossible_speed")
	if got.Action != ActionReject {
		t.Errorf("action = %q, want reject", got.Action)
	}
	if v := got.Params.Float("min_seconds_hard", -1); v != 4.0 {
		t.Errorf("min_seconds_hard = %v, want 4", v)
	}
	// A row that sets one threshold must not blank the others.
	if v := got.Params.Float("min_seconds_easy", -1); v != 1.0 {
		t.Errorf("min_seconds_easy = %v, want the default 1.0 to survive a partial edit", v)
	}
}

// A bad hand-edit is an operator mistake. It must degrade, not escalate and
// not break the answer path.
func TestReloadSurvivesABadRow(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db)
	ctx := context.Background()
	if err := store.Seed(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}

	if _, err := db.Exec(`UPDATE anticheat_rules SET action = 'delete-account', params = 'not json' WHERE name = 'fast_streak'`); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := store.Reload(ctx); err != nil {
		t.Fatalf("reload should tolerate a bad row: %v", err)
	}

	got := store.Rule("fast_streak")
	if got.Action != ActionObserve {
		t.Errorf("an unknown action must fall back to observe, got %q", got.Action)
	}
	if v := got.Params.Float("window", -1); v != 5 {
		t.Errorf("unreadable params must fall back to defaults, window = %v", v)
	}
}

// Deleting a row resets the rule rather than disabling the detector.
func TestDeletedRowFallsBackToTheDefault(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db)
	ctx := context.Background()
	if err := store.Seed(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := db.Exec(`DELETE FROM anticheat_rules WHERE name = 'uniform_timing'`); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if err := store.Reload(ctx); err != nil {
		t.Fatalf("reload: %v", err)
	}

	got := store.Rule("uniform_timing")
	if !got.Enabled || got.Action != ActionObserve {
		t.Errorf("a missing row should leave the compiled default, got %+v", got)
	}
	if v := got.Params.Float("window", -1); v != 6 {
		t.Errorf("window = %v, want the default 6", v)
	}
}

// Seed must not stomp an operator's edits on restart.
func TestSeedIsIdempotent(t *testing.T) {
	db := openTestDB(t)
	store := NewStore(db)
	ctx := context.Background()
	if err := store.Seed(ctx); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := db.Exec(`UPDATE anticheat_rules SET enabled = 0 WHERE name = 'fast_streak'`); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := store.Seed(ctx); err != nil {
		t.Fatalf("reseed: %v", err)
	}
	if err := store.Reload(ctx); err != nil {
		t.Fatalf("reload: %v", err)
	}
	if store.Rule("fast_streak").Enabled {
		t.Fatal("re-seeding on restart overwrote an operator's edit")
	}
}

// A nil database is the "settings unavailable" path: defaults, no errors.
func TestNilDatabaseRunsOnDefaults(t *testing.T) {
	store := NewStore(nil)
	if err := store.Seed(context.Background()); err != nil {
		t.Fatalf("seed with no db: %v", err)
	}
	if err := store.Reload(context.Background()); err != nil {
		t.Fatalf("reload with no db: %v", err)
	}
	if got := store.Rule("impossible_speed"); !got.Enabled {
		t.Fatal("defaults should still be in place")
	}
}
