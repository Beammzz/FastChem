# internal/services

## Purpose

All game logic: question generation, chemistry reference data, scoring, ELO rating, matchmaking, ranked/room match state, and the in-memory stores.

## Ownership

Question generation:
- `problem.go` — `Problem`, the `Topic` interface, `Source`, `topicRegistry`, `validProblem`, `generateValid`
- `distractors.go` — `assembleChoices` and the misconception/top-up machinery every topic shares
- `topics_easy.go` — atomic structure, oxidation number, state of matter
- `topics_medium.go` — mole concept (mass↔moles, moles→particles), `avogadro`
- `topics_hard.go` — dilution, preparing from solid, freezing-point depression, `solutes`, `kfWater`, `kbWater`
- `generator.go` — `QuestionGenerator`, the casual-play entry point
- `ranked_questions.go` — `GenerateRankedQuestions`, the seeded set for 1v1

Reference data:
- `elements.go` — elements 1–30 with real isotope lists
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

### Question generation

- **One implementation serves every mode.** A `Topic` takes a `*Source`; casual play passes a clock-seeded Source held by `QuestionGenerator`, ranked and room matches pass a match-seeded one. Never fork a parallel "seeded" copy of a generator — that is what let ranked silently lose the `state_of_matter` topic and its anti-repeat.
- **Seeded generation must be reproducible.** `GenerateRankedQuestions(seed)` has to yield the same questions, the same order, and the same choice order every time, because `ranked_matches.seed` is the only record of what was asked. Nothing in the generation path may read the global `rand`, the clock, or **Go map iteration order** — ranging over a map to collect choices is what broke this before, moving the correct answer's index between runs.
- **`topicRegistry` order is part of the seeded contract.** Selection indexes into slices derived from it, so inserting a topic anywhere but the end changes what every existing seed generates. Append.
- **Distractors come from misconceptions, not noise.** Each topic authors its wrong answers as named student errors — sign reversed, operation inverted, unit conversion dropped, wrong constant — in priority order, with the mistake in a trailing comment. Perturbing the answer by a random factor makes the correct value the only plausible number in the list and teaches nothing when picked.
- **Choice assembly must terminate for any answer.** `assembleChoices` takes candidates plus a `topUp` function and is bounded by `topUpLimit`. `floatTopUp` falls back to steps of exactly one unit in the last printed decimal and `sciTopUp` to 5% mantissa shifts, both of which always yield distinct strings. A `for len(choices) < 4` loop with no bound spins forever whenever the answer is small enough that too few distinct formatted values exist nearby.
- **Every problem is validated before it is served.** `validProblem` requires four distinct non-empty choices with `CorrectIndex` in range; scoring resolves correctness by index, so a malformed problem mis-scores silently rather than failing loudly. `generateValid` retries a bounded number of times.
- **Category IDs are a frontend contract.** `atomic_structure`, `oxidation_number`, `state_of_matter`, `mole_concept`, `dilution`, `preparing_solution`, `freezing_point` — `QuestionCard.tsx` and `RankedMatchUI.tsx` switch on these for labels and colours, and the single-player page offers them as filters. Renaming one is a two-sided change.
- **`Source` is the only randomness.** It is mutex-guarded because the casual generator holds one for the life of the process while Gin serves requests in parallel. `Source.Pick` also supplies the anti-repeat window, so consecutive questions do not reuse the same element or compound.
- **Chemistry data stays in `elements.go` / `compounds.go` / `solutes`.** Topics read data, they do not embed it. Neutron questions draw from the real `Isotopes` list rather than offsetting the atomic number, which used to invent nuclides like hydrogen-7.

### Scoring and matches

- **The client is never trusted with the answer.** `GlobalQuestionStore` holds `correctIndex` keyed by question ID; `/api/answer` and match endpoints resolve correctness from the store. Stored questions expire — `main.go` sweeps entries older than 10 minutes — so a stale ID is a normal, expected miss, not an error to paper over.
- **Scoring is authoritative here.** `CalculateScore` returns base + speed bonus before combo; `CalculateScoreWithCombo` applies it. Per-difficulty parameters live in one map:

  | difficulty | time limit | base | speed multiplier |
  |---|---|---|---|
  | easy | 30s | 5 | 1 |
  | medium | 60s | 10 | 2 |
  | hard | 120s | 20 | 3 |

  Unknown difficulties fall back to easy via `GetDifficultyConfig`.
- **Two combo curves exist and are not interchangeable.** Single player uses `ComboMultiplier` (3→1.2, 6→1.5, 10+→2.0); ranked uses `RankedComboMultiplier` (2→1.1, 3→1.2, 5+→1.5). Wrong answers reset the streak in both.
- **ELO updates go through `rating.go`** using `models.EloKFactor`; ratings floor at 0.
- **Every store is mutex-guarded and sweeper-backed.** `MatchStore`, `QuestionStore`, `RankedMatchService`, and `RoomService` each expose a `CleanupOlderThan` / `CleanupStale*` method driven by `main.go`. A new store must follow the same shape or it grows without bound.
- **Send on match channels only via `SafeSend` / `SafeSendTimeout`.** They recover from sends on closed channels, which is what keeps a disconnecting player from panicking the opponent's goroutine.
- **`ActiveRankedMatch` state is guarded by its own mutex,** reachable through `Mu()`. Read-modify-write of scores, progress, or completion must hold it; `finalizeMatchLocked` assumes the caller already does.

## Work Guidance

Adding a topic:

1. Implement `Topic` in the `topics_*.go` file for its difficulty — `Category()`, `Difficulty()`, `Generate(*Source)`.
2. Draw data with `s.Pick(pool, n)` so the anti-repeat window applies, and all randomness from the `Source`.
3. Author distractors as misconception values with the mistake named in a comment; pass them through `floatCandidates` / `intCandidates` / `sciCandidates` and finish with `assembleChoices`.
4. Append to `topicRegistry` — never insert.
5. Add the category to the label and colour switches in `QuestionCard.tsx` and `RankedMatchUI.tsx`, and to `ALL_CATEGORIES` in the single-player page if it should be selectable.

Other rules:

- Constrain generated values so answers stay in a range that has room for distinct distractors — `minMoles` exists for exactly this reason.
- Keep `withinScale` in mind when authoring candidates: a distractor more than `scaleBand` away from the answer is dropped, because it would be eliminated on sight. The factor-of-1000 unit slip sits deliberately inside that band.
- Persistence belongs at match finalization (`persistMatchResult`), not on every answer.

## Verification

```bash
go build ./...
go vet ./...
go test ./internal/services/ -race
```

`generator_test.go` locks down the properties this package's correctness rests on: seeded generation is byte-identical across runs, every topic yields four distinct choices with the answer among them, generation terminates for small answers, ranked covers every registered topic, neutron questions use real isotopes, and concurrent generation is race-free. Generation is wrapped in `withDeadline` so a non-terminating generator fails the test instead of hanging the suite.

Chemistry changes still need a manual sanity check of the generated question text and the correct answer.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
