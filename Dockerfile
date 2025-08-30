# Build stage
FROM golang:1.21-alpine AS builder

# Install git for version information
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version injection
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
ARG GO_VERSION

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -X main.GoVersion=${GO_VERSION} -w -s" \
    -a -installsuffix cgo \
    -o webshell-detect \
    ./cmd/webshell-detect

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S webshell && \
    adduser -u 1001 -S webshell -G webshell

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/webshell-detect .

# Copy configuration files
COPY --from=builder /app/config ./config

# Change ownership to non-root user
RUN chown -R webshell:webshell /app

# Switch to non-root user
USER webshell

# Expose port (if needed for future web interface)
# EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./webshell-detect --version || exit 1

# Set entrypoint
ENTRYPOINT ["./webshell-detect"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="webshell-detect" \
      org.opencontainers.image.description="A tool for detecting PHP webshells using static analysis" \
      org.opencontainers.image.vendor="Security Team" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/your-org/webshell_detect" \
      org.opencontainers.image.documentation="https://github.com/your-org/webshell_detect/blob/main/README.md"