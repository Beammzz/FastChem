# internal/services

## Purpose

All game logic: question generation, chemistry reference data, scoring, ELO rating, matchmaking, ranked/room match state, and the in-memory stores.

## Ownership

Question generation:
- `generator.go` — `QuestionGenerator`, the `QuestionGen` interface, `EasyGenerator`, distractor helpers, recent-question tracking
- `medium_generator.go` — mole conversions (mass ↔ moles, moles ↔ particles), `avogadro`
- `hard_generator.go` — solution chemistry: dilution, preparing from solid, freezing-point depression
- `ranked_questions.go` — seeded, deterministic question sets for 1v1

Reference data:
- `elements.go` — element list
- `compounds.go` — compound and state-of-matter data

Rules:
- `scoring.go` — `DifficultyConfigs`, `ComboMultiplier`, `CalculateScore`, `CalculateScoreWithCombo`
- `rating.go` — ELO expected score and rating update, `RankedComboMultiplier`, `CalculateRankedScore`

Match state:
- `question_store.go` — `GlobalQuestionStore`, correct answers held server-side by question ID
- `match_store.go` — `GlobalMatchStore`, active single-player sessions
- `ranked_match.go` — `RankedMatchService`, live 1v1 state, disconnects, finalization, persistence
- `matchmaking.go` — `MatchmakingQueue`, rating-based pairing loop
- `room.go` — `RoomService`, code-based custom rooms

## Local Contracts

- **The client is never trusted with the answer.** `GlobalQuestionStore` holds `correctIndex` keyed by question ID; `/api/answer` and match endpoints resolve correctness from the store. Stored questions expire — `main.go` sweeps entries older than 10 minutes — so a stale ID is a normal, expected miss, not an error to paper over.
- **Scoring is authoritative here.** `CalculateScore` returns base + speed bonus before combo; `CalculateScoreWithCombo` applies it. Per-difficulty parameters live in one map:

  | difficulty | time limit | base | speed multiplier |
  |---|---|---|---|
  | easy | 30s | 5 | 1 |
  | medium | 60s | 10 | 2 |
  | hard | 120s | 20 | 3 |

  Unknown difficulties fall back to easy via `GetDifficultyConfig`.
- **Two combo curves exist and are not interchangeable.** Single player uses `ComboMultiplier` (3→1.2, 6→1.5, 10+→2.0); ranked uses `RankedComboMultiplier` (2→1.1, 3→1.2, 5+→1.5). Wrong answers reset the streak in both.
- **Ranked questions must be reproducible from the seed.** `GenerateRankedQuestions(seed)` builds the whole set from a seeded `*rand.Rand` so both players get identical questions in identical order. Never use package-level `rand`, `time.Now()`, or map iteration order inside the seeded path — the `seeded*` helpers exist for exactly this reason.
- **ELO updates go through `rating.go`** using `models.EloKFactor`; ratings floor at 0.
- **Distractor generation must terminate for any answer value.** The random phase is capped at `distractorAttempts`, then `fillFloatLadder` / `fillSciLadder` top the set up using steps of exactly one unit in the last printed decimal, which always yields distinct strings. A `for len(choices) < 4` loop with no bound spins forever whenever the answer is small enough that too few distinct formatted values exist nearby — never write one.
- **`pickFreshIndex` holds `recentMu`.** The `recent*` slices are package-level and Gin serves requests concurrently.
- **Every store is mutex-guarded and sweeper-backed.** `MatchStore`, `QuestionStore`, `RankedMatchService`, and `RoomService` each expose a `CleanupOlderThan` / `CleanupStale*` method driven by `main.go`. A new store must follow the same shape or it grows without bound.
- **Send on match channels only via `SafeSend` / `SafeSendTimeout`.** They recover from sends on closed channels, which is what keeps a disconnecting player from panicking the opponent's goroutine.
- **`ActiveRankedMatch` state is guarded by its own mutex,** reachable through `Mu()`. Read-modify-write of scores, progress, or completion must hold it; `finalizeMatchLocked` assumes the caller already does.

## Work Guidance

- Add a difficulty tier by implementing `Generate() models.Question` and registering it in `QuestionGenerator`; do not branch on difficulty strings inside handlers.
- Generated distractors must be plausible, unique, and never duplicate the correct answer — reuse `generateDistractors`, `generateOxidationDistractors`, `generateFloatDistractors`, or `generateSciDistractors` rather than writing new ad-hoc ones.
- Keep chemistry facts in `elements.go` / `compounds.go`; generators read data, they do not embed it.
- `pickFreshIndex` with the `recentWindow` of 8 keeps questions from repeating back-to-back — route new random picks through it.
- Persistence belongs at match finalization (`persistMatchResult`), not on every answer.

## Verification

```bash
go build ./...
go vet ./...
go test ./internal/services/ -race
```

`generator_test.go` covers the invariants that scoring depends on: every question offers four distinct choices with the correct answer among them, distractor generation terminates for small answers, and concurrent generation is race-free. Tests wrap generation in `withDeadline` so a non-terminating generator fails instead of hanging the suite.

Chemistry changes still need a manual sanity check of the generated question text and the correct answer.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
