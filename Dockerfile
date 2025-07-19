# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o main .

# Run stage (slim image)
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main", "serve-http"]