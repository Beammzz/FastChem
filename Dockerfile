# syntax=docker/dockerfile:1.7

############################
# Stage 1 — Frontend Builder
############################
FROM node:20-bullseye-slim AS frontend-builder

WORKDIR /app/frontend

# Install build dependencies only once (better cache)
COPY frontend/package*.json ./
RUN npm ci

# Copy rest of frontend
COPY frontend/ ./

# Build static export (Next.js with output: "export")
RUN npm run build


############################
# Stage 2 — Backend Builder
############################
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app/backend

# Install required packages for cgo + sqlite build
RUN apk add --no-cache \
    build-base \
    sqlite-dev \
    git \
    ca-certificates

# Cache Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Copy frontend static export
COPY --from=frontend-builder /app/frontend/out ./frontend/out

# Enable cgo and build for linux
ENV CGO_ENABLED=1 GOOS=linux

# Build optimized binary
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server


############################
# Stage 3 — Runtime (Minimal)
############################
FROM alpine:3.20

WORKDIR /app

# Copy binary
COPY --from=backend-builder /app/server ./server

# Copy static frontend (if NOT embedding)
COPY --from=backend-builder /app/backend/frontend/out ./frontend/out

# Install runtime dependencies and create nonroot user
RUN apk add --no-cache sqlite-libs ca-certificates \
    && addgroup -S nonroot \
    && adduser -S -G nonroot nonroot

ENV PORT=8080

EXPOSE 8080

USER nonroot

ENTRYPOINT ["/app/server"]