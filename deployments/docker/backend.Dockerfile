# =============================================================================
# DeltaGov Backend Dockerfile
# Multi-stage build optimized for Google Cloud Run
# =============================================================================
# Build from repo root:
#   docker build -f deployments/docker/backend.Dockerfile -t deltagov-api .
#
# Run API server:
#   docker run -p 8080:8080 -e PORT=8080 -e DATABASE_URL=... deltagov-api
#
# Run Ingestor (single pass for Cloud Run Jobs):
#   docker run --entrypoint /app/ingestor -e DATABASE_URL=... deltagov-api --single-run
#
# Run Ingestor (continuous polling):
#   docker run --entrypoint /app/ingestor -e DATABASE_URL=... deltagov-api
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Build
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

# Install git for fetching dependencies and ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go module files first for better layer caching
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY backend/ .

# Build both binaries with optimizations
# CGO_ENABLED=0 for static binary (required for distroless)
# -ldflags="-w -s" strips debug info for smaller binary
# -trimpath removes file system paths from binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /build/bin/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /build/bin/ingestor ./cmd/ingestor

# -----------------------------------------------------------------------------
# Stage 2: Runtime (Distroless for minimal attack surface)
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

# Labels for container registry
LABEL org.opencontainers.image.source="https://github.com/drewjst/deltagov"
LABEL org.opencontainers.image.description="DeltaGov API Server"
LABEL org.opencontainers.image.licenses="AGPL-3.0"

# Copy timezone data for proper time handling
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates for HTTPS requests to external APIs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binaries
COPY --from=builder /build/bin/api /app/api
COPY --from=builder /build/bin/ingestor /app/ingestor

WORKDIR /app

# Environment variables
# PORT: HTTP port for the API server (Cloud Run sets this automatically)
ENV PORT=8080

# Expose the default port
EXPOSE 8080

# Default: Run the API server
# Override with --entrypoint /app/ingestor for ingestion jobs
ENTRYPOINT ["/app/api"]
