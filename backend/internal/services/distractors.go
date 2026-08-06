package services

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Distractors are built from misconceptions — sign flips, inverted operations,
// dropped unit conversions, the wrong constant — rather than from random
// perturbation of the answer.
//
// Two things follow. A wrong answer becomes diagnostic: picking it means you
// made a specific, nameable mistake. And the correct value stops standing out,
// because every option is the result of a real calculation rather than noise
// scattered around the truth.
//
// Each topic authors its candidate list in priority order. Candidates that
// duplicate the answer, fail to render, or fall wildly out of scale are
// dropped, and the top-up ladder fills any remaining slots.

// topUpLimit bounds the fallback search. The ladders below reach three
// distinct values within a handful of steps; this only guards against a
// topUp function that never makes progress.
const topUpLimit = 64

// scaleBand is the widest ratio a distractor may sit at relative to the
// answer. It is deliberately loose enough to admit the classic factor-of-1000
// unit slip, which is worth keeping despite being easy to eliminate.
const scaleBand = 2000.0

// assembleChoices returns four distinct choices — three distractors plus the
// correct answer — shuffled, with the index of the correct one.
//
// candidates are misconception-derived distractors in priority order. topUp
// supplies more when fewer than three survive, and must eventually return
// values distinct from everything already chosen.
func assembleChoices(s *Source, correctStr string, candidates []string, topUp func(k int) string) ([]string, int) {
	out := make([]string, 0, 4)
	seen := map[string]bool{correctStr: true}

	add := func(c string) {
		if c == "" || seen[c] {
			return
		}
		seen[c] = true
		out = append(out, c)
	}

	for _, c := range candidates {
		if len(out) == 3 {
			break
		}
		add(c)
	}
	for k := 1; len(out) < 3 && k <= topUpLimit; k++ {
		add(topUp(k))
	}

	out = append(out, correctStr)
	s.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })

	idx := 0
	for i, c := range out {
		if c == correctStr {
			idx = i
			break
		}
	}
	return out, idx
}

// ─── Float answers ───────────────────────────────────────────────

// formatPrecisionStep returns one unit in the last decimal place of a
// printf-style float format, e.g. "%.2f mol" → 0.01.
func formatPrecisionStep(format string) float64 {
	i := strings.Index(format, "%.")
	if i < 0 {
		return 0.01
	}
	prec := 0
	for j := i + 2; j < len(format) && format[j] >= '0' && format[j] <= '9'; j++ {
		prec = prec*10 + int(format[j]-'0')
	}
	return math.Pow(10, -float64(prec))
}

// snapTo rounds v onto the grid the answer will be printed on, so a value and
// its rendered form never disagree.
func snapTo(v, step float64) float64 {
	if step <= 0 {
		return v
	}
	return math.Round(v/step) * step
}

// withinScale drops distractors that are orders of magnitude away from the
// answer. Past that distance a wrong option is eliminated on sight and wastes
// a slot that a plausible near-miss could fill.
func withinScale(correct, v float64) bool {
	a, b := math.Abs(correct), math.Abs(v)
	if a == 0 || b == 0 {
		return a == b
	}
	r := a / b
	if r < 1 {
		r = 1 / r
	}
	return r <= scaleBand
}

// floatCandidates renders misconception values, dropping non-finite results
// and anything out of scale with the answer.
func floatCandidates(correct float64, format string, values []float64) []string {
	step := formatPrecisionStep(format)
	out := make([]string, 0, len(values))
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) || !withinScale(correct, v) {
			continue
		}
		out = append(out, fmt.Sprintf(format, snapTo(v, step)))
	}
	return out
}

// floatTopUp yields fallback distractors in two tiers. First proportional
// near-misses — 10%, 20%, 30% either side — which read as arithmetic slips and
// stay well separated. Then a ladder one format-step at a time, which is
// guaranteed to produce distinct strings however small the answer is; this is
// what makes generation terminate for values like 0.02 mol.
func floatTopUp(correct float64, format string) func(int) string {
	step := formatPrecisionStep(format)
	return func(k int) string {
		if k <= 6 {
			mag := 0.1 * float64((k+1)/2)
			if k%2 == 0 {
				mag = -mag
			}
			v := snapTo(correct*(1+mag), step)
			if math.Abs(v-correct) < step/2 {
				return "" // indistinguishable once printed
			}
			return fmt.Sprintf(format, v)
		}
		dir := 1.0
		if correct < 0 {
			dir = -1.0
		}
		return fmt.Sprintf(format, snapTo(correct+dir*float64(k-6)*step, step))
	}
}

// ─── Integer answers ─────────────────────────────────────────────

func formatInt(n int) string { return strconv.Itoa(n) }

// formatOxidation renders an oxidation number the way chemistry writes it.
func formatOxidation(n int) string {
	if n > 0 {
		return "+" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}

// intCandidates renders misconception values that fall inside [minVal, maxVal].
func intCandidates(minVal, maxVal int, values []int, format func(int) string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if v < minVal || v > maxVal {
			continue
		}
		out = append(out, format(v))
	}
	return out
}

// intTopUp walks outward from the answer, alternating above and below, and
// skips anything outside [minVal, maxVal]. As long as the band holds at least
// four values this always finds three distinct distractors.
func intTopUp(correct, minVal, maxVal int, format func(int) string) func(int) string {
	return func(k int) string {
		d := (k + 1) / 2
		if k%2 == 0 {
			d = -d
		}
		v := correct + d
		if v < minVal || v > maxVal {
			return ""
		}
		return format(v)
	}
}

// ─── Scientific-notation answers ─────────────────────────────────

// formatSci formats a number in "a.bb × 10^n" style.
func formatSci(v float64) string {
	if v == 0 {
		return "0"
	}
	exp := math.Floor(math.Log10(math.Abs(v)))
	mantissa := v / math.Pow(10, exp)
	mantissa = math.Round(mantissa*100) / 100
	return fmt.Sprintf("%.2f × 10^%.0f", mantissa, exp)
}

func sciCandidates(correct float64, values []float64) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) || v == 0 || !withinScale(correct, v) {
			continue
		}
		out = append(out, formatSci(v))
	}
	return out
}

// sciTopUp shifts the mantissa in 5% steps. The mantissa always lies in
// [1, 10), so every step changes the second decimal and the values stay
// distinct.
func sciTopUp(correct float64) func(int) string {
	return func(k int) string {
		f := 1 + 0.05*float64((k+1)/2)
		if k%2 == 0 {
			f = 1 - 0.05*float64((k+1)/2)
		}
		if f <= 0 {
			return ""
		}
		return formatSci(correct * f)
	}
}
