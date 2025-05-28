# ---------- Build Stage ----------
FROM golang:1.24.3 AS builder

WORKDIR /app

# Install dependencies needed to build with CGo and Kafka
RUN apt-get update && apt-get install -y \
    gcc \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/api ./cmd/api/main.go

# ---------- Final Stage ----------
FROM debian:bookworm-slim

# Install runtime dependencies including CA certificates
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bin/api /bin/api

EXPOSE 8091

ENTRYPOINT ["/bin/api"]
