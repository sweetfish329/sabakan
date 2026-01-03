# Production Containerfile for Sabakan
# Multi-stage build for minimal production image.
#
# Build:
#   podman build -f Containerfile -t sabakan:latest .
#
# Run:
#   podman run -d --name sabakan \
#     -v /run/podman/podman.sock:/run/podman/podman.sock:Z \
#     -v sabakan-data:/data:Z \
#     -p 1323:1323 \
#     sabakan:latest

# ============================================
# Stage 1: Build Backend (Go)
# ============================================
FROM golang:1.25-alpine AS backend-builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source and build
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o sabakan ./cmd/sabakan

# ============================================
# Stage 2: Build Frontend (Angular)
# ============================================
FROM oven/bun:1-alpine AS frontend-builder

WORKDIR /build

# Copy package files and install dependencies (skip prepare scripts for CI)
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile --ignore-scripts

# Copy source and build
COPY frontend/ ./
RUN bun run build

# ============================================
# Stage 3: Production Runtime
# ============================================
FROM alpine:3.21

LABEL maintainer="Sabakan Development Team"
LABEL description="Sabakan - Game Container Management System"
LABEL version="1.0.0"

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S sabakan && adduser -S sabakan -G sabakan

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /build/sabakan ./

# Copy frontend static files
COPY --from=frontend-builder /build/dist/frontend/browser ./public

# Copy default configuration
COPY backend/config.example.toml ./config.example.toml

# Create data directory
RUN mkdir -p /data && chown -R sabakan:sabakan /data /app

# Switch to non-root user
USER sabakan

# Environment variables
ENV CONFIG_PATH=/app/config.toml
ENV DATABASE_PATH=/data/sabakan.db
ENV PODMAN_SOCKET=unix:///run/podman/podman.sock

# Expose port
EXPOSE 1323

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1323/health || exit 1

# Entrypoint
ENTRYPOINT ["./sabakan"]
