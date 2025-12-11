# ============================================
# Development Stage (with Delve + Air)
# ============================================
FROM golang:1.23-alpine AS development

# Install development tools
RUN apk add --no-cache git curl

# Install Air for hot reload (pin to v1.61.1 for Go 1.23 compatibility)
RUN go install github.com/air-verse/air@v1.61.1

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (will be overridden by volume mount in docker-compose)
COPY . .

# Expose API port and Delve debugger port
EXPOSE 8080 2345

# Start with Air for hot reload + debugging
CMD ["air", "-c", ".air.toml"]

# ============================================
# Builder Stage (for production binary)
# ============================================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized production binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/api

# ============================================
# Production Stage (minimal, secure)
# ============================================
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
