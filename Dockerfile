# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates curl

# Create a non-root user and app directory
RUN adduser -D -s /bin/sh algafood && \
    mkdir -p /app && \
    chown -R algafood:algafood /app

# Set working directory to /app instead of /root
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy Docker configuration to /app (owned by algafood)
COPY --from=builder /app/config.docker.yaml ./config.yaml

# Change ownership of all files in /app
RUN chown -R algafood:algafood /app

# Switch to non-root user
USER algafood

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Command to run
CMD ["./main"]



