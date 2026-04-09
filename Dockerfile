# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pinta-api ./cmd/main.go

# ── Stage 2: Run ─────────────────────────────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/pinta-api .

EXPOSE 8080

CMD ["./pinta-api"]
