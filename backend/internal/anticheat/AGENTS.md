# internal/anticheat

## Purpose

Watches answers that have already been scored and reports the ones a human could not plausibly have produced. Detection only — it decides nothing about chemistry, scoring, or match flow, and it is advisory by default.

## Ownership

- `anticheat.go` — `Signal`, `Verdict`, `Finding`, `Action`, the `Engine`, and the process-wide default (`Init`, `Evaluate`, `Sweep`)
- `detectors.go` — the `Detector` interface and the three built-in rules: `impossible_speed`, `fast_streak`, `uniform_timing`
- `settings.go` — `Rule`, `Params`, the `Store` that loads them from the `anticheat_rules` table, and `Save` for operator edits
- `history.go` — the bounded per-subject answer window the windowed rules read
- `sink.go` — the `Sink` seam findings are written through, plus `SlogSink`, `DBSink`, `DiscardSink`, `MultiSink`

## Local Contracts

- **Everything ships as `observe`.** A rule's default action is to record and change nothing, so enabling this package alters no player's experience. Escalation is an operator's decision, made in the database. A new detector that ships as `reject` breaks that promise.
- **Enforcement is wired even though it is off.** Call sites honour `Verdict.Rejected()` today. That is what makes escalation a config change rather than a code change — do not "simplify" by ignoring the verdict.
- **This package imports nothing from `services` or `handlers`.** That one-way edge is what lets both of them call it. It depends on the standard library and `database/sql` only. Anything a detector needs arrives on the `Signal`.
- **`Evaluate` touches memory only.** Settings are cached and reloaded on a timer, never read from SQLite on the answer path. It is called while `MatchStore` and `ActiveRankedMatch` locks are held, so any I/O or sleep added here becomes a stall in every match. A `Sink` that writes to a database must buffer.
- **Detectors read thresholds from `Params`, never constants.** A rule with a hard-coded number cannot be retuned, which is the whole point of the table. Use `Params.Float(key, default)` so a missing or malformed value degrades to the compiled default.
- **Bad settings degrade, never escalate.** An unknown action becomes `observe`; unreadable params fall back to defaults; a deleted row restores the compiled default rather than disabling the detector. `Reload` logs and continues rather than failing.
- **Detector names are stable identifiers.** The name is the primary key in `anticheat_rules` and appears in logs. Renaming one orphans its settings row and loses its history.
- **Registration is append-only** (`builtinDetectors`, then `Engine.Use`). Order fixes the order findings appear in.
- **Timed-out answers are not evidence.** Every timeout lands on the same difficulty time limit, so a player who walks away produces a run of identical times. `uniform_timing` skips them and `fast_streak` requires correctness. A new rule that reads `TimeSpent` across a window must do the same or it will flag absence as automation.
- **The subject key differs by mode.** Signed-in play keys on `SubjectUser`; casual keys on `SubjectIP`, which groups everyone behind one NAT together. Casual findings are correspondingly weak.
- **`DBSink` drops rather than blocks.** Findings go onto a buffered channel drained by one writer goroutine; a full buffer increments `Dropped()` and discards the event. `Record` runs under a match lock, so waiting on SQLite's single writer there would stall every match. The admin console surfaces the drop count — a sink that silently loses findings would make the table look like quiet rather than saturation.
- **`Save` validates, `Reload` degrades.** They are not symmetric on purpose: `Reload` reads whatever is in the table and falls back on anything it cannot parse, because the answer path runs through it. `Save` takes input from a UI, so an unknown rule name or action is an error it can report instead of a silent downgrade. It reloads on success, so a console edit applies at once rather than on the next tick.

## Call sites

Four, all passing server-measured values — `TimeSpent` is derived from when the server issued the question, never from the client:

- `handlers/question.go` — `POST /api/answer`, `ModeCasual`
- `services/match_store.go` — `RecordAnswer`, `ModeMatch`
- `services/ranked_match.go` — `SubmitAnswer`, `ModeRanked` or `ModeRoom` depending on the sign of the match ID (rooms use a negative synthetic one)

## Settings

Rules live in `anticheat_rules`, seeded at startup and reloaded every 30 seconds by `cmd/server/main.go`. The admin console's Anti-cheat tab edits them through `Store.Save`; the SQL below is the same change by hand, picked up on the next reload:

```sql
UPDATE anticheat_rules SET action = 'reject'          WHERE name = 'impossible_speed';
UPDATE anticheat_rules SET enabled = 0                WHERE name = 'uniform_timing';
UPDATE anticheat_rules SET params = '{"window": 8}'   WHERE name = 'fast_streak';
```

`params` is merged over the compiled defaults, so setting one threshold leaves the rest alone.

## Work Guidance

- A new rule is a `Detector` in `detectors.go`, an entry in `builtinDetectors`, and a row in `defaultRules` — all three, or it is untunable.
- Pick thresholds that under-flag. A rule that cries wolf gets ignored, and observe mode exists to show real distributions before anything is tightened.
- A new action means extending the `Action` enum, `ParseAction`, `severity`, and every call site. Do not return an action no call site handles.
- Another destination for findings is a `Sink` implementation passed to `NewEngine`; no detector or call site changes. `main.go` already fans out to `SlogSink` and `DBSink` through a `MultiSink`.

## Verification

```bash
go build ./...
go vet ./...
go test ./internal/anticheat/ -race
```

End to end, against a running server: answer questions instantly and confirm `anticheat finding` lines appear in the log; then set a rule's action to `reject`, wait for the reload, and confirm fast correct answers come back `"correct": false`.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
