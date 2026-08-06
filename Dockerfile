# ─────────────────────────────────────────────
# Stage 1 – Build the Next.js static export
# ─────────────────────────────────────────────
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* frontend/yarn.lock* frontend/pnpm-lock.yaml* ./

RUN \
  if [ -f pnpm-lock.yaml ]; then \
    npm install -g pnpm && pnpm install --frozen-lockfile; \
  elif [ -f yarn.lock ]; then \
    yarn install --frozen-lockfile; \
  else \
    npm ci; \
  fi

COPY frontend/ .

RUN \
  if [ -f pnpm-lock.yaml ]; then \
    pnpm run build; \
  elif [ -f yarn.lock ]; then \
    yarn build; \
  else \
    npm run build; \
  fi

# ─────────────────────────────────────────────
# Stage 2 – Build the Go backend (pure Go, no CGO)
# ─────────────────────────────────────────────
FROM golang:1.24-bookworm AS backend-builder

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /fastchem-server ./cmd/server

# ─────────────────────────────────────────────
# Stage 3 – Minimal runtime image
# ─────────────────────────────────────────────
FROM debian:bookworm-slim AS runtime

# ca-certificates for any outbound TLS
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Non-root user
RUN useradd -m -u 1001 fastchem

WORKDIR /app

# Copy compiled server and frontend static export
COPY --from=backend-builder /fastchem-server ./fastchem-server
COPY --from=frontend-builder /app/frontend/out ./frontend/out

# Persistent data directory for the SQLite database
RUN mkdir -p /data && chown fastchem:fastchem /data

USER fastchem

EXPOSE 8080

ENV PORT=8080 \
    GIN_MODE=release \
    DB_PATH=/data/fastchem.db \
    FRONTEND_DIR=/app/frontend/out

CMD ["/app/fastchem-server"]
