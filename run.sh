#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

mkdir -p "$ROOT_DIR/logs"

# Choose frontend build command (prefer pnpm if available)
if command -v pnpm >/dev/null 2>&1; then
  PKG_CMD="pnpm"
elif command -v npm >/dev/null 2>&1; then
  PKG_CMD="npm"
else
  echo "No package manager found (pnpm or npm)." >&2
  exit 1
fi

echo "=== FastChem: Building & starting ==="

# 1. Build the frontend (static export → frontend/out/)
echo "Building frontend..."
(cd "$ROOT_DIR/frontend" && $PKG_CMD run build) 2>&1 | tee "$ROOT_DIR/logs/frontend-build.log"
echo "Frontend build done → frontend/out/"

# 2. Start backend (serves API + static frontend on :8080)
echo "Starting backend on :8080 (logs -> logs/backend.log)"
export FRONTEND_DIR="$ROOT_DIR/frontend/out"
(
  cd "$ROOT_DIR/backend" &&
  go run ./cmd/server
) > "$ROOT_DIR/logs/backend.log" 2>&1 &
BACK_PID=$!

echo "Backend PID: $BACK_PID"
echo ""
echo "===> Open http://localhost:8080 to use FastChem <==="
echo ""

shutdown() {
  echo "Stopping backend..."
  kill "$BACK_PID" 2>/dev/null || true
  wait "$BACK_PID" 2>/dev/null || true
  echo "Stopped. Logs are in $ROOT_DIR/logs/"
}

trap shutdown INT TERM EXIT

wait "$BACK_PID"

wait
