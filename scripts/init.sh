#!/bin/bash

# Remnawave Telegram Shop Bot - Initialization Script

set -e

echo "🚀 Initializing Remnawave Telegram Shop Bot..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "   Download from: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is not supported. Please install Go 1.21 or later."
    exit 1
fi

echo "✅ Go version $GO_VERSION detected"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. You'll need Docker to run the database."
    echo "   Download from: https://www.docker.com/get-started"
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "⚠️  Docker Compose is not installed. You'll need it to run the full stack."
    echo "   Download from: https://docs.docker.com/compose/install/"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please edit it with your configuration."
else
    echo "✅ .env file already exists"
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p logs
mkdir -p ssl
mkdir -p migrations

# Build the application
echo "🔨 Building application..."
go build -o remnawave-bot ./cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Application built successfully"
else
    echo "❌ Build failed"
    exit 1
fi

# Check if PostgreSQL is running
if command -v docker &> /dev/null; then
    echo "🐘 Starting PostgreSQL database..."
    docker run -d --name remnawave-postgres \
        -e POSTGRES_DB=remnawave_bot \
        -e POSTGRES_USER=remnawave_bot \
        -e POSTGRES_PASSWORD=remnawave_bot_password \
        -p 5432:5432 postgres:15-alpine
    
    echo "✅ PostgreSQL started"
    echo "   Database: remnawave_bot"
    echo "   User: remnawave_bot"
    echo "   Password: remnawave_bot_password"
    echo "   Port: 5432"
fi

echo ""
echo "🎉 Initialization complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file with your configuration"
echo "2. Run: ./remnawave-bot"
echo "   or: go run ./cmd/main.go"
echo ""
echo "For Docker deployment:"
echo "1. Edit .env file"
echo "2. Run: docker-compose up -d"
echo ""
echo "For development:"
echo "1. Install tools: make install-tools"
echo "2. Run tests: make test"
echo "3. Format code: make fmt"
echo ""
echo "Happy coding! 🚀"
