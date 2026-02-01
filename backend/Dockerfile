# Stage 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY . .

# Build the specific service binary
ARG SERVICE_NAME
RUN go build -o /app/service ./cmd/${SERVICE_NAME}/main.go

# Stage 2: Small runtime
FROM alpine:latest
WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/service .

# FIX: Create the internal directory structure and copy the migrations
# This ensures that "file://internal/migrations" works in your Go code
RUN mkdir -p ./internal/migrations
COPY --from=builder /app/internal/migrations ./internal/migrations

EXPOSE 8080
CMD ["./service"]