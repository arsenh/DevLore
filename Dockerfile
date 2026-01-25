# ---------- Build Stage ----------
FROM golang:1.25-alpine AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o devlore cmd/app/main.go


# ---------- Run Stage ----------
FROM alpine:3.23.2

WORKDIR /app
COPY --from=builder /app/devlore .
COPY --from=builder /app/.env .env
COPY --from=builder /app/internal/migrations /app/internal/migrations
CMD ["./devlore"]
