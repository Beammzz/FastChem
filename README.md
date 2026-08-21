# ⚡ FastChem

A fast-paced chemistry practice game — like [fastmath.io](https://fastmath.io), but for chemistry.

## Stack

- **Frontend:** Next.js 14 (App Router, TypeScript, Tailwind CSS)
- **Backend:** Go (Gin framework)
- **Communication:** REST API (JSON)

## Features

- **Single player** — pick the question count, a difficulty, or your own set of
  topics grouped by curriculum chapter
- **Ranked 1v1** — matchmaking by ELO over a seeded question set, so both
  players answer the same questions in the same order
- **Custom rooms** — code-based lobbies
- Auto-generated questions across 23 topics (see [Question coverage](#question-coverage))
- Per-difficulty timers and scoring, with a combo multiplier for answer streaks
- Casual play needs no account; leaderboards and ranked require one
- **Admin console** at `/admin` — analytics, both leaderboards, user management,
  live match/queue/room state, the anti-cheat manager and a question-bank
  browser (see [Admin console](#admin-console))
- **`fastchemctl`** — a management CLI for the jobs a browser cannot do: grant
  the first admin, retune anti-cheat, back up a live database
  (see [fastchemctl](#fastchemctl))

## Project Structure

```
FastChem/
├── backend/
│   ├── cmd/server/main.go        # Entry point
│   ├── internal/
│   │   ├── handlers/             # HTTP handlers
│   │   ├── models/               # Data models
│   │   └── services/             # Question generation logic
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── page.tsx          # Home page
│   │   │   ├── game/page.tsx     # Game page
│   │   │   ├── layout.tsx
│   │   │   └── globals.css
│   │   ├── components/           # UI components
│   │   ├── hooks/                # Custom React hooks
│   │   ├── lib/                  # API client
│   │   └── types/                # TypeScript types
│   ├── .env.local
│   └── package.json
└── README.md
```

## Setup & Run

### Prerequisites

- Go 1.22+
- Node.js 18+
- npm 9+

### Backend

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

Backend runs on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:3000`.

### Full stack in one command

```bash
./run.sh      # macOS / Linux
```

```powershell
.\run.ps1     # Windows PowerShell
```

Both build the frontend static export and start the backend, which serves the
export and the API together on `http://localhost:8080`. Output goes to `logs/`.

### API Endpoints

Everything below sits under `/api` behind an IP rate limiter (burst 30,
refilling 10/s). Rows marked ✔ need an `Authorization: Bearer <token>` header;
the two WebSocket routes take the same token as a `?token=` query parameter,
since a browser cannot set headers on a WebSocket handshake.

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET  | `/api/question` | — | Generate a question. `?difficulty=easy\|medium\|hard`, or `?categories=<csv>` to draw from chosen topics |
| POST | `/api/answer` | — | Submit `questionId` + `selectedIndex`; the server times the answer, marks it, and returns the score along with the correct index |
| POST | `/api/validate` | — | Legacy self-scoring check — the caller supplies the answer it wants compared. Nothing in the app calls it |
| GET  | `/api/health` | — | Health check |
| POST | `/api/auth/register` | — | Create an account |
| POST | `/api/auth/login` | — | Exchange credentials for a JWT |
| GET  | `/api/auth/me` | ✔ | The current user |
| GET  | `/api/leaderboard` | — | Top scores |
| GET  | `/api/profile/:username` | — | Public profile and score history (`?page=`) |
| POST | `/api/scores` | ✔ | Record a finished single-player run |
| GET  | `/api/scores/me` | ✔ | Your own score history |
| POST | `/api/match/start` | ✔ | Start a server-tracked match |
| POST | `/api/match/answer` | ✔ | Answer inside a match — additionally limited to burst 3, refilling 1/s |
| POST | `/api/match/end` | ✔ | Finish a match and persist the result |
| GET  | `/api/ranked/ws` | `?token=` | WebSocket: matchmaking and live ranked play |
| GET  | `/api/ranked/stats` | ✔ | Your ELO, wins and losses |
| GET  | `/api/ranked/history` | ✔ | Your past ranked matches |
| GET  | `/api/ranked/leaderboard` | — | Ranked ladder |
| GET  | `/api/room/ws` | `?token=` | WebSocket: custom rooms, `?action=create` or `?action=join&code=<code>` |

Admin routes need a token **and** `users.is_admin`; a signed-in non-admin gets
`403`. The flag is read from the database on every request, so demoting an
account takes effect at once rather than when its week-long token expires.
Writes are `POST` because CORS allows only `GET`, `POST` and `OPTIONS`.

| Method | Endpoint | Description |
|---|---|---|
| GET  | `/api/admin/overview` | Headline counts: accounts, games, findings, live state, uptime, DB size |
| GET  | `/api/admin/analytics` | Daily activity (`?days=`, max 90), difficulty split, per-topic accuracy, hour-of-day, rating bands, findings per rule |
| GET  | `/api/admin/leaderboards` | Both boards, read live rather than from the public 30s cache (`?limit=`) |
| GET  | `/api/admin/users` | Paged user table with play stats (`?page=`, `?pageSize=`, `?search=`, `?sort=`, `?order=`) |
| POST | `/api/admin/users/update` | Set `isAdmin`, `rating`, `totalPoints` or `password`; omitted fields are left alone |
| POST | `/api/admin/users/delete` | Delete an account and everything it owns |
| GET  | `/api/admin/anticheat/rules` | Live rule settings plus how often each has fired |
| POST | `/api/admin/anticheat/rules` | Retune a rule — takes effect immediately, no restart and no waiting for the reload ticker |
| GET  | `/api/admin/anticheat/findings` | Paged findings with filters (`?rule=`, `?mode=`, `?action=`, `?userId=`) and breakdowns of the filtered set |
| GET  | `/api/admin/live` | Matchmaking queue, in-progress matches, open rooms — in-memory state no table holds |
| GET  | `/api/admin/matches` | Recent ranked matches, single-player matches and submitted scores (`?limit=`) |
| GET  | `/api/admin/topics` | Every registered question topic |
| GET  | `/api/admin/topics/preview` | Generate a throwaway sample question with its answer (`?category=`) |

Any path that does not start with `/api/` falls through to the Next.js static
export in `frontend/out/`.

### Question Response

```json
{
  "id": "uuid",
  "question": "ธาตุ คาร์บอน (C) มีจำนวนโปรตอนเท่าใด?",
  "choices": ["4", "6", "8", "5"],
  "timeLimit": 30,
  "category": "atomic_structure",
  "difficulty": "easy"
}
```

No `correctIndex`: the server keeps the answer, keyed by question id, and
returns it from `POST /api/answer` together with the marking and the score.
So the client can highlight the right choice afterwards, but cannot read it
out of the network tab beforehand.

## Question coverage

Questions are generated, not stored, from 23 topics spanning บทที่ 2–13 of
สาระเคมี in หลักสูตรแกนกลางการศึกษาขั้นพื้นฐาน พ.ศ. 2551 (ฉบับปรับปรุง พ.ศ. 2560)
— the curriculum Thai ม.4–ม.6 chemistry is taught from:

| Chapter | Topics |
|---|---|
| บทที่ 2 อะตอมและสมบัติของธาตุ | โครงสร้างอะตอม, การจัดเรียงอิเล็กตรอน |
| บทที่ 3 พันธะเคมี | ชนิดพันธะเคมี, รูปร่างโมเลกุล (VSEPR) |
| บทที่ 4 โมลและสูตรเคมี | โมลคอนเซ็ปต์, มวลต่อโมล, ร้อยละโดยมวล |
| บทที่ 5 สารละลาย | ความเข้มข้น, การเจือจาง, เตรียมสารละลาย, จุดเยือกแข็ง |
| บทที่ 6 ปริมาณสัมพันธ์ | ปริมาณสัมพันธ์, สารกำหนดปริมาณและร้อยละผลได้ |
| บทที่ 7 แก๊ส | กฎของแก๊ส, แก๊สอุดมคติ (PV = nRT) |
| บทที่ 8 อัตราการเกิดปฏิกิริยาเคมี | อัตราเฉลี่ยและอันดับของปฏิกิริยา |
| บทที่ 9 สมดุลเคมี | ค่าคงที่สมดุล K |
| บทที่ 10 กรด–เบส | pH, pOH, Ka |
| บทที่ 11 เคมีไฟฟ้า | เลขออกซิเดชัน, E°cell |
| บทที่ 12 เคมีอินทรีย์ | หมู่ฟังก์ชัน |
| บทที่ 13 พอลิเมอร์ | มอนอเมอร์, แบบเติม / แบบควบแน่น |

Single player can filter by any of these; ranked draws 4 easy, 3 medium and
3 hard from the same registry.

## Anti-cheat

The server already withholds the answer, measures answer times itself, and
scores every mode server-side. On top of that, `backend/internal/anticheat`
watches for answer patterns a human cannot produce:

| Rule | Fires on |
|---|---|
| `impossible_speed` | A correct answer returned faster than the question can be read — per-difficulty floor, 1.0s easy to 2.0s hard |
| `fast_streak` | 5 correct answers in a row, each under 3s — fast on every question, including the hard ones |
| `uniform_timing` | 6 answers with under 0.2s of spread — machine cadence rather than a person |

**Everything ships in observe mode.** Rules record findings and change nothing
a player sees. Enforcement is per-rule, stored in the `anticheat_rules` table,
and reloaded every 30 seconds — so escalating is a setting, not a deploy. The
admin console's Anti-cheat tab edits the same rules and applies them
immediately; the equivalent by hand is:

```sql
-- reject flagged answers (they score zero) instead of just logging them
UPDATE anticheat_rules SET action = 'reject' WHERE name = 'impossible_speed';

-- retune a threshold; params merge over the defaults
UPDATE anticheat_rules SET params = '{"window": 8}' WHERE name = 'fast_streak';

-- turn a rule off entirely
UPDATE anticheat_rules SET enabled = 0 WHERE name = 'uniform_timing';
```

Adding a rule means implementing `Detector`, registering it, and giving it a
default row — see `backend/internal/anticheat/AGENTS.md`. Findings go to the
structured log and to the `anticheat_findings` table through the `Sink`
interface; the database sink buffers and writes on its own goroutine, so it
never sits on a player's answer path. Under sustained load it drops rather than
blocks, and the console reports the drop count so the table is never silently
an undercount.

## Admin console

`/admin` is a single page of tabs backed by the `/api/admin/*` routes above:

| Tab | What it shows |
|---|---|
| ภาพรวม | Headline counts, daily activity, difficulty and hour-of-day distributions, rating bands, per-topic accuracy, and process health (uptime, DB size, goroutines) |
| ผู้เล่น | Searchable, sortable user table — points, rating, accuracy, flag count, last played, online now. Edit rating, points and password; grant or remove admin; delete an account |
| Anti-cheat | Every rule with its live thresholds and hit counts, editable in place; the findings log with filters and breakdowns |
| สด | Matchmaking queue, matches in progress with per-player score and connection state, open rooms, and recent finished matches. Refreshes every 5s |
| กระดานผู้นำ | Both leaderboards, read live rather than cached |
| คลังคำถาม | All 23 topics; generate a sample question with its answer |

**Access.** No account is an admin by default. Either set `ADMIN_USERNAMES` to a
comma-separated list and restart — accounts that already exist are promoted at
startup, a name with no account yet is skipped — or grant it with no restart at
all using the CLI below:

```bash
ADMIN_USERNAMES=alice,bob ./fastchem-server   # at startup
fastchemctl user grant alice                  # any time, takes effect at once
```

An admin cannot remove their own admin flag or delete their own account —
otherwise the last admin could lock everyone out.

## fastchemctl

A second binary (`backend/cmd/fastchemctl`) that manages the database from a
shell. It opens SQLite directly rather than calling the admin API, which is the
point: it works before the first admin exists, while the server is down, and
over `docker exec`. The web console is the better tool for everything else.

```bash
cd backend && go build -o fastchemctl ./cmd/fastchemctl
./fastchemctl status
```

| Command | Does |
|---|---|
| `status` | Counts, plus whether any admin exists and whether any rule is rejecting |
| `user list` | Accounts with points, rating, games and flag counts — `-search`, `-limit`, `-sort` |
| `user show <name>` | One account in full |
| `user grant\|revoke <name>` | Admin console access |
| `user passwd <name>` | New password, prompted and not echoed |
| `user delete <name>` | The account and everything it owns — confirmation required, `-yes` to skip |
| `rules list` | Anti-cheat rules, live thresholds and hit counts |
| `rules set <name>` | Retune — `-enabled=false`, `-action observe\|reject`, repeatable `-param key=value` |
| `findings` | Recent detections — `-rule`, `-mode`, `-user`, `-limit` |
| `backup <dest>` | `VACUUM INTO` snapshot, consistent even while the server writes |

The database is `-db`, then `$DB_PATH`, then `fastchem.db` — the same order the
server uses, so both find the same file from the same directory.

**Two changes are safe while the server is running**, because the server
re-reads both: `is_admin` on every admin request, and `anticheat_rules` on a
30-second ticker. So this needs no restart and no deploy:

```bash
fastchemctl user grant alice
fastchemctl rules set impossible_speed -action reject -param min_seconds_easy=1.5
```

Unset flags keep their stored value, so retuning one threshold leaves the rule's
action and its other thresholds alone. Deleting a player who is mid-match is the
one operation that does not reconcile with live state — the server keeps that
match in memory until it ends or is swept.

In Docker the binary is on `PATH`:

```bash
docker compose exec fastchem fastchemctl user grant alice
```

## Extending

- **New question types** — implement `Topic` in `backend/internal/services/topics_*.go`,
  append it to `topicRegistry`, and add its id to `frontend/src/data/categories.ts`.
  See `backend/internal/services/AGENTS.md` for the full contract.
- **New reference data** — add a data file next to `elements.go`; topics read
  data, they never embed it.

## Docker: Build & Run with Persistent DB

- **DB file location**: the backend uses SQLite and defaults to `fastchem.db` in the working directory (see [backend/internal/database/db.go](backend/internal/database/db.go#L18) and [backend/internal/config/config.go](backend/internal/config/config.go#L22)). To persist data across container restarts, mount a host directory and point `DB_PATH` to a file inside that directory.

- **Build the image**:

```bash
docker build -t fastchem:local .
```

- **Run with a host volume mapped to `/data` and persistent DB**:

```bash
mkdir -p ./data
docker run --rm -p 8080:8080 \
  -v "$(pwd)/data:/data" \
  -e DB_PATH=/data/fastchem.db \
  fastchem:local
```

- **Docker Compose**: see `docker-compose.yml` included in the repo for an example service that maps `./data` → container `/data` and sets `DB_PATH=/data/fastchem.db`.

If you want, I can add a small helper script to build & run the image locally.
