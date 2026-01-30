# Stage 1: Build binary
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy module files
COPY go.mod go.sum ./
RUN go mod download

# Copy EVERYTHING from backend/ into /app
COPY . .

# Argument passed from docker-compose (auth, contract, etc.)
ARG SERVICE_NAME

# We build from the cmd folder where your main.go lives
RUN go build -o /app/service ./cmd/${SERVICE_NAME}/main.go

# Stage 2: Small runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]