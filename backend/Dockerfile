# Stage 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
# Build the binary
RUN go build -o /app/service ./cmd/${SERVICE_NAME}/main.go

# Stage 2: Small runtime
FROM alpine:latest
WORKDIR /root/

# Copy the binary
COPY --from=builder /app/service .

# --- CRITICAL FIXES START HERE ---
# 1. Grant execution permission (This fixes the "Permission Denied" error)
RUN chmod +x ./service

# 2. Copy migrations
RUN mkdir -p ./internal/migrations
COPY --from=builder /app/internal/migrations ./internal/migrations

# 3. Copy Swagger docs (This fixes the "Fetch Error" in Swagger UI)
ARG SERVICE_NAME
RUN mkdir -p ./cmd/${SERVICE_NAME}/docs
COPY --from=builder /app/cmd/${SERVICE_NAME}/docs ./cmd/${SERVICE_NAME}/docs
# --- CRITICAL FIXES END HERE ---

EXPOSE 8080
CMD ["./service"]