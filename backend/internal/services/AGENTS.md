# internal/services

## Purpose

All game logic: question generation, chemistry reference data, scoring, ELO rating, matchmaking, ranked/room match state, and the in-memory stores.

## Ownership

Question generation — 23 topics covering สาระเคมี บทที่ 2–13 of หลักสูตรแกนกลาง
พ.ศ. 2551 (ฉบับปรับปรุง พ.ศ. 2560), the curriculum ม.4–ม.6 is taught from:

- `problem.go` — `Problem`, the `Topic` interface, `Source`, `topicRegistry`, `validProblem`, `generateValid`
- `distractors.go` — `assembleChoices` and the misconception/top-up machinery every topic shares
- `topics_easy.go` — atomic structure, oxidation number, state of matter, electron configuration, bond type, molecular shape (VSEPR), functional group, polymer
- `topics_medium.go` — mole concept, molar mass, percent composition, concentration units, stoichiometry, gas laws; `avogadro`, `minMoles`, `molarVolumeSTP`, `kelvinOffset`
- `topics_hard.go` — dilution, preparing from solid, freezing-point depression, limiting reagent and percent yield, ideal gas, reaction rate, equilibrium, acid–base, electrochemistry; `solutes`, `kfWater`, `kbWater`, `gasConstantAtm`, `pKw`
- `generator.go` — `QuestionGenerator`, the casual-play entry point
- `ranked_questions.go` — `GenerateRankedQuestions`, the seeded set for 1v1

Reference data — topics read these, they do not embed data:
- `elements.go` — elements 1–30 with real isotope lists
- `compounds.go` — compound and state-of-matter data
- `periodic.go` — `shellElements` (Z = 1–20 shell/subshell configurations, group, period) and `representativeGroups` (บทที่ 2)
- `bonding.go` — `bondClasses` / `bondSubstances` and `molecularShapes` / `shapeMolecules` (บทที่ 3)
- `formulas.go` — `atomicMasses` and `formulaCompounds`, whose molar mass is derived from their parts (บทที่ 4)
- `reactions.go` — `reactions` and `twoReactantReactions`, plus `equilibriumReactions` (บทที่ 6, 9)
- `organic.go` — `organicClasses` / `functionalGroups` / `organicCompounds` and `polymers` (บทที่ 12–13)
- `electrolytes.go` — `weakAcids`, `strongAcids`, `standardPotentials`, `kWater` (บทที่ 10–11)

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
- **A distractor a student can rule out without doing the chemistry is a wasted slot.** `floatCandidates` drops any value that renders as `0.00` beside a non-zero answer, and `percentCandidates` drops anything above 100% for an answer that is a share of a whole — so the inverted-ratio error cannot be taught on a percentage question and its slot goes to an error that produces a believable figure.
- **Choice assembly must terminate for any answer.** `assembleChoices` takes candidates plus a `topUp` function and is bounded by `topUpLimit`. `floatTopUp` falls back to steps of exactly one unit in the last printed decimal and `sciTopUp` to 5% mantissa shifts, both of which always yield distinct strings. A `for len(choices) < 4` loop with no bound spins forever whenever the answer is small enough that too few distinct formatted values exist nearby.
- **Closed-vocabulary topics answer from a fixed pool and have no ladder.** Bond class, molecular shape, organic class, functional group, group number, polymer name: the only honest wrong answers are the other real entries, so these use `textCandidates` with `noTopUp`. Each pool must hold at least four distinct entries — with fewer, every problem the topic makes fails `validProblem`, and `generateValid` serves the malformed one after its retries. `TestClosedVocabularyPoolsCanFillFourChoices` guards this.
- **Every problem is validated before it is served.** `validProblem` requires four distinct non-empty choices with `CorrectIndex` in range; scoring resolves correctness by index, so a malformed problem mis-scores silently rather than failing loudly. `generateValid` retries a bounded number of times.
- **Category IDs are a frontend contract.** All 23 are listed in `frontend/src/data/categories.ts`, which supplies the Thai label, chip colour, and chapter for each; `QuestionCard.tsx` and `RankedMatchUI.tsx` read it, and the single-player page offers them as filters grouped by chapter. Adding or renaming an id is a two-sided change — `TestTopicCategoriesMatchFrontendIDs` pins the list on this side.
- **`Source` is the only randomness.** It is mutex-guarded because the casual generator holds one for the life of the process while Gin serves requests in parallel. `Source.Pick` also supplies the anti-repeat window, so consecutive questions do not reuse the same element or compound.
- **Chemistry data stays in the reference files listed under Ownership.** Topics read data, they do not embed it. Neutron questions draw from the real `Isotopes` list rather than offsetting the atomic number, which used to invent nuclides like hydrogen-7.
- **Derive what can be derived.** `formulaCompound.MolarMass()` sums its parts instead of storing a number, because a hand-typed molar mass that drifts from the subscripts makes the answer and every distractor wrong at once. Where a value cannot be derived — `reactionSpecies.MolarMass`, `standardPotentials`, `weakAcids.Ka` — a test checks it: `TestReactionsConserveMass` catches a wrong coefficient or molar mass, `TestFormulaMolarMassesAreCorrect` spot-checks against the book.
- **`shellElements` stops at calcium on purpose.** Past Z = 20 the 4s/3d overlap makes the 2-8-8-2 picture wrong, and a question built on it would teach the error rather than test it.

### Scoring and matches

- **The client is never trusted with the answer, and never told it in advance.** `GlobalQuestionStore` holds `correctIndex` keyed by question ID; `/api/answer` and match endpoints resolve correctness from the store, and the question payload itself omits the answer (`json:"-"` on both `models.Question.CorrectIndex` and `RankedQuestion.CorrectIndex`). It is revealed in the marking response, once the player has committed. Stored questions expire — `main.go` sweeps entries older than 10 minutes — so a stale ID is a normal, expected miss, not an error to paper over.
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
3. Author distractors as misconception values with the mistake named in a comment; pass them through `floatCandidates` / `percentCandidates` / `intCandidates` / `sciCandidates` / `textCandidates` and finish with `assembleChoices`.
4. Append to `topicRegistry` — never insert.
5. Add the category id to `frontend/src/data/categories.ts` with its label, chapter, and chip colour, and to the `want` list in `TestTopicCategoriesMatchFrontendIDs`.

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

`generator_test.go` locks down the properties this package's correctness rests on: seeded generation is byte-identical across runs, every topic yields four distinct choices with the answer among them, generation terminates for small answers, ranked and both filters reach every registered topic, and concurrent generation is race-free. Generation is wrapped in `withDeadline` so a non-terminating generator fails the test instead of hanging the suite.

It also checks the chemistry the data files assert: reactions conserve mass, molar masses match the book, every formula symbol has an atomic mass (a missing key reads as 0.0 rather than failing), shell configurations sum to the atomic number, neutron questions use real isotopes, cell potentials are non-negative, pH answers stay on the 0–14 scale, and the closed-vocabulary pools can fill four choices.

Chemistry changes still need a manual sanity check of the generated question text and the correct answer — a topic can be well-formed, deterministic, and still ask something false.

## Child DOX Index

- No child AGENTS.md files. This folder is flat.
