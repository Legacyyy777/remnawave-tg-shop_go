# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Download dependencies
ENV GOSUMDB=off
RUN go mod download

# Copy source code
COPY . .

# Generate go.sum and clean module cache
RUN go mod tidy && go clean -modcache

# Create empty .env file for final stage
RUN echo "# Environment variables will be set via docker-compose" > /tmp/.env

# Build the application
ENV GOSUMDB=off
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy environment file
COPY --from=builder /tmp/.env .env

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
