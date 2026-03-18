# Stage 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git and build essentials
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

# FIX 1: Grant execution permissions to the binary
RUN chmod +x ./service

# FIX 2: Copy migrations so GORM can find them at runtime
RUN mkdir -p ./internal/migrations
COPY --from=builder /app/internal/migrations ./internal/migrations

# FIX 3: Copy Swagger documentation for the specific service
# This maps the generated docs to the path expected by your gin-swagger handler
ARG SERVICE_NAME
RUN mkdir -p ./cmd/${SERVICE_NAME}/docs
COPY --from=builder /app/cmd/${SERVICE_NAME}/docs ./cmd/${SERVICE_NAME}/docs

# The port is handled by the docker-compose mapping, 
# but 8080 is a good default for documentation.
EXPOSE 8080

# Run the service
CMD ["./service"]