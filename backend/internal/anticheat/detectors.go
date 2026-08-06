package anticheat

import (
	"fmt"
	"math"
)

// Detector examines one answer and returns a Finding when it looks inhuman.
//
// A detector reads its thresholds from Params rather than hard-coding them, so
// an operator can retune it without a deploy. It must tolerate missing or
// nonsensical params by falling back to its own defaults: a bad row in the
// database is an operator mistake, not a reason to break scoring.
//
// Inspect must not block. It runs on the answer path, sometimes under a match
// lock, so no I/O and no sleeping.
type Detector interface {
	// Name is the stable identifier used as the settings key and in logs.
	// Renaming one orphans its settings row.
	Name() string
	// Inspect returns nil when the answer looks fine. Rule and Action are
	// filled in by the engine; a detector sets only Detail.
	Inspect(sig Signal, h History, p Params) *Finding
}

// builtinDetectors is append-only, like the question topic registry: the order
// fixes the order findings appear in, and each name is a settings key that
// exists in the database.
func builtinDetectors() []Detector {
	return []Detector{
		impossibleSpeed{},
		fastStreak{},
		uniformTiming{},
	}
}

// ─── impossible_speed ───────────────────────────────────────────

// impossibleSpeed catches a correct answer returned faster than the question
// can be read. It is the crudest signal and the hardest to argue with: nobody
// reads a Thai stoichiometry problem and four options in under a second.
//
// Only correct answers count. A fast wrong answer is a misclick, and flagging
// those would bury the real signal in noise.
type impossibleSpeed struct{}

func (impossibleSpeed) Name() string { return "impossible_speed" }

func (impossibleSpeed) Inspect(sig Signal, _ History, p Params) *Finding {
	if !sig.Correct || sig.TimedOut {
		return nil
	}
	// A harder question takes longer to read, so the floor is per-difficulty
	// with a shared fallback.
	floor := p.Float("min_seconds_"+sig.Difficulty, p.Float("min_seconds", 1.0))
	if floor <= 0 || sig.TimeSpent >= floor {
		return nil
	}
	return &Finding{
		Detail: fmt.Sprintf("correct %s answer in %.2fs, under the %.2fs floor",
			sig.Difficulty, sig.TimeSpent, floor),
	}
}

// ─── fast_streak ────────────────────────────────────────────────

// fastStreak catches the player who is not impossibly fast on any single
// question but is fast and correct every time. A run of quick correct answers
// is what an answer key looks like; a strong student still slows down on the
// hard ones.
type fastStreak struct{}

func (fastStreak) Name() string { return "fast_streak" }

func (fastStreak) Inspect(_ Signal, h History, p Params) *Finding {
	window := int(p.Float("window", 5))
	limit := p.Float("max_seconds", 3.0)
	if window < 2 || limit <= 0 || len(h.Recent) < window {
		return nil
	}
	run := h.Recent[len(h.Recent)-window:]
	slowest := 0.0
	for _, a := range run {
		if !a.Correct || a.TimedOut || a.TimeSpent >= limit {
			return nil
		}
		if a.TimeSpent > slowest {
			slowest = a.TimeSpent
		}
	}
	return &Finding{
		Detail: fmt.Sprintf("%d correct answers in a row, slowest %.2fs, all under %.2fs",
			window, slowest, limit),
	}
}

// ─── uniform_timing ─────────────────────────────────────────────

// uniformTiming catches a script rather than a cheat sheet. Human answer times
// scatter — an easy recall question and a two-step calculation do not take the
// same seven-tenths of a second — so a run with almost no spread is machinery,
// whether the answers are right or wrong.
type uniformTiming struct{}

func (uniformTiming) Name() string { return "uniform_timing" }

func (uniformTiming) Inspect(_ Signal, h History, p Params) *Finding {
	window := int(p.Float("window", 6))
	maxSD := p.Float("max_stddev", 0.2)
	if window < 3 || maxSD <= 0 {
		return nil
	}

	// Timed-out answers all land on the difficulty's time limit, so a player
	// who walks away produces a run of near-identical times. That is absence,
	// not automation — skip them.
	times := make([]float64, 0, window)
	for i := len(h.Recent) - 1; i >= 0 && len(times) < window; i-- {
		if h.Recent[i].TimedOut {
			continue
		}
		times = append(times, h.Recent[i].TimeSpent)
	}
	if len(times) < window {
		return nil
	}

	mean, sd := meanStdDev(times)
	if sd >= maxSD {
		return nil
	}
	return &Finding{
		Detail: fmt.Sprintf("%d answers averaging %.2fs with σ=%.3fs, under the %.3fs spread floor",
			len(times), mean, sd, maxSD),
	}
}

// meanStdDev returns the mean and population standard deviation of xs, which
// must not be empty.
func meanStdDev(xs []float64) (mean, sd float64) {
	for _, x := range xs {
		mean += x
	}
	mean /= float64(len(xs))
	for _, x := range xs {
		d := x - mean
		sd += d * d
	}
	return mean, math.Sqrt(sd / float64(len(xs)))
}
