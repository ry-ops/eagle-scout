# eagle-scout - Docker Scout MCP Server
# Multi-stage build for minimal production image

# Use latest Go 1.25 patch to fix stdlib vulnerabilities
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /eagle-scout ./cmd/eagle-scout

# Production stage - use docker CLI image as base for scout access
# Pin to 29.2.1 — Docker CLI built with Go 1.25.6+ (fixes CVE-2025-61726, CVE-2025-61728, CVE-2025-61730)
FROM docker:29.2.1-cli

LABEL org.opencontainers.image.title="eagle-scout"
LABEL org.opencontainers.image.description="MCP server for Docker Scout - container security scanning"
LABEL org.opencontainers.image.source="https://github.com/ry-ops/eagle-scout"
LABEL org.opencontainers.image.licenses="MIT"

# Upgrade Alpine packages to fix CVE-2026-25210 (expat)
RUN apk upgrade --no-cache

# Install Docker Scout CLI plugin
COPY --from=docker/scout-cli:1.19.0 /docker-scout /usr/libexec/docker/cli-plugins/docker-scout

# Copy binary from builder
COPY --from=builder /eagle-scout /usr/local/bin/eagle-scout

# Create non-root user
RUN adduser -D -u 1001 scout

USER scout

ENTRYPOINT ["/usr/local/bin/eagle-scout"]
