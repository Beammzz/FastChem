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
FROM golang:1.24-bullseye AS backend-builder

WORKDIR /app/backend

# Install required packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Cache Go modules
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Copy frontend static export
COPY --from=frontend-builder /app/frontend/out ./frontend/out

# Multi-arch support
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

# Build optimized binary
RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server


############################
# Stage 3 — Runtime (Minimal)
############################
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary
COPY --from=backend-builder /app/server ./server

# Copy static frontend (if NOT embedding)
COPY --from=backend-builder /app/backend/frontend/out ./frontend/out

# Copy system CA certificates
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PORT=8080

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]