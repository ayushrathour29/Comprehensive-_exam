# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git (for go mod) and build tools
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o jobqueue ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary and .env
COPY --from=builder /app/jobqueue .
COPY .env .env

# Expose port
EXPOSE 8080

CMD ["./jobqueue"]