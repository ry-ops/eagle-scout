# scout-mcp - Docker Scout MCP Server
# Multi-stage build for minimal production image

FROM golang:1.22-alpine AS builder

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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /scout-mcp ./cmd/scout-mcp

# Production stage - use docker CLI image as base for scout access
FROM docker:cli

LABEL org.opencontainers.image.title="scout-mcp"
LABEL org.opencontainers.image.description="MCP server for Docker Scout - container security scanning"
LABEL org.opencontainers.image.source="https://github.com/ry-ops/scout-mcp"
LABEL org.opencontainers.image.licenses="MIT"

# Copy binary from builder
COPY --from=builder /scout-mcp /usr/local/bin/scout-mcp

# Create non-root user
RUN adduser -D -u 1001 scout

USER scout

ENTRYPOINT ["/usr/local/bin/scout-mcp"]
