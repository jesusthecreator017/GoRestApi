# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api
RUN go build -o migrate ./cmd/migrate

# Stage 2: Runtime
FROM alpine:3.21

RUN adduser -D -u 1000 appuser

WORKDIR /app

COPY --from=builder /build/main .
COPY --from=builder /build/migrate .
COPY --from=builder /build/cmd/migrate/migrations ./cmd/migrate/migrations
COPY scripts/entrypoint.sh .

RUN chmod +x entrypoint.sh

USER appuser

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
