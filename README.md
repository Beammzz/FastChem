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

### API Endpoints

| Method | Endpoint         | Description              |
|--------|------------------|--------------------------|
| GET    | `/api/question`  | Get a random question    |
| POST   | `/api/validate`  | Validate an answer       |
| GET    | `/api/health`    | Health check             |

### Question Response

```json
{
  "id": "uuid",
  "question": "ธาตุ คาร์บอน (C) มีจำนวนโปรตอนเท่าใด?",
  "choices": ["4", "6", "8", "5"],
  "correctIndex": 1,
  "timeLimit": 30,
  "category": "atomic_structure",
  "difficulty": "easy"
}
```

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
