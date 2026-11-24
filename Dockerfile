# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT} -w -s" \
    -o /out/webshell-detect ./cmd/webshell-detect

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && addgroup -g 1001 -S webshell \
    && adduser -u 1001 -S webshell -G webshell

WORKDIR /app

COPY --from=builder /out/webshell-detect /app/webshell-detect
COPY config /app/config
COPY LICENSE /app/LICENSE

RUN chown -R webshell:webshell /app
USER webshell

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["./webshell-detect", "-version"]

ENTRYPOINT ["/app/webshell-detect"]
CMD ["-help"]

LABEL org.opencontainers.image.title="webshell-detect" \
      org.opencontainers.image.description="Static PHP webshell detection using taint and pointer analysis" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/0xundef/php-webshell-detect"
