# ── build stage ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gateway ./cmd/gateway

# ── runtime stage ────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# non-root user
RUN adduser -D -u 1000 appuser
USER appuser

COPY --from=builder /gateway /usr/local/bin/gateway

EXPOSE 8080

ENTRYPOINT ["gateway"]
