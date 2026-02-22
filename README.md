# ⚡ FastChem

A fast-paced chemistry practice game — like [fastmath.io](https://fastmath.io), but for chemistry.

## Stack

- **Frontend:** Next.js 14 (App Router, TypeScript, Tailwind CSS)
- **Backend:** Go (Gin framework)
- **Communication:** REST API (JSON)

## Features (MVP)

- Single Player Mode (Easy difficulty)
- Auto-generated questions:
  - **Atomic Structure** — protons, neutrons, electrons
  - **Oxidation Numbers** — common compounds
- Configurable timer (30s – 180s)
- Score tracking (+10 per correct answer)
- No authentication required

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
  "question": "How many protons does Carbon have?",
  "choices": ["4", "6", "8", "5"],
  "correctIndex": 1,
  "timeLimit": 15,
  "category": "atomic_structure",
  "difficulty": "easy"
}
```

## Extending

The codebase is designed for extension:

- **Medium/Hard modes** — add new generator methods in `services/generator.go`
- **New question types** — add new compound/element data files
- **PVP mode** — WebSocket support can be added to the Gin server
- **Authentication** — add middleware to Gin router

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
