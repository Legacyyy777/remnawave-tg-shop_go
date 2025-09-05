# Makefile for Remnawave Telegram Shop Bot

.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=remnawave-bot
DOCKER_IMAGE=remnawave-bot
DOCKER_TAG=latest

# Build the application
build:
	go build -o $(BINARY_NAME) ./cmd/main.go

# Run the application
run:
	go run ./cmd/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run benchmarks
bench:
	go test -bench=. ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate documentation
docs:
	godoc -http=:6060

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run:
	docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker Compose up
docker-up:
	docker-compose up -d

# Docker Compose down
docker-down:
	docker-compose down

# Docker Compose logs
docker-logs:
	docker-compose logs -f

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  bench          - Run benchmarks"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Download dependencies"
	@echo "  docs           - Generate documentation"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-up      - Start with Docker Compose"
	@echo "  docker-down    - Stop Docker Compose"
	@echo "  docker-logs    - Show Docker Compose logs"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help"
