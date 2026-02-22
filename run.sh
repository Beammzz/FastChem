#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

mkdir -p "$ROOT_DIR/logs"

# Choose frontend dev command (prefer pnpm if available)
if command -v pnpm >/dev/null 2>&1; then
  FRONTEND_CMD="pnpm dev"
elif command -v npm >/dev/null 2>&1; then
  FRONTEND_CMD="npm run dev"
else
  echo "No frontend package manager found (pnpm or npm)." >&2
  exit 1
fi

echo "Starting FastChem stack..."

# Start backend (Go)
echo "Starting backend (logs -> logs/backend.log)"
(
  cd "$ROOT_DIR/backend" && 
  go run ./cmd/server
) > "$ROOT_DIR/logs/backend.log" 2>&1 &
BACK_PID=$!

# Start frontend (Next.js)
echo "Starting frontend (logs -> logs/frontend.log)"
(
  cd "$ROOT_DIR/frontend" && 
  $FRONTEND_CMD
) > "$ROOT_DIR/logs/frontend.log" 2>&1 &
FRONT_PID=$!

echo "Backend PID: $BACK_PID"
echo "Frontend PID: $FRONT_PID"

shutdown() {
  echo "Stopping services..."
  kill "$BACK_PID" 2>/dev/null || true
  kill "$FRONT_PID" 2>/dev/null || true
  wait "$BACK_PID" 2>/dev/null || true
  wait "$FRONT_PID" 2>/dev/null || true
  echo "Stopped. Logs are in $ROOT_DIR/logs/"
}

trap shutdown INT TERM EXIT

wait
