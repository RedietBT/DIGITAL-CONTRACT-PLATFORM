# Stage 1: Build binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE_NAME
RUN go build -o /app/service ./cmd/${SERVICE_NAME}/main.go

# Stage 2: Small runtime
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]