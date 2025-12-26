# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o lox-bridge ./cmd/bridge

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/lox-bridge .

# Create a default config directory
RUN mkdir configs

# Command to run
ENTRYPOINT ["./lox-bridge"]
