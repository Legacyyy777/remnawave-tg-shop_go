@echo off
REM Remnawave Telegram Shop Bot - Initialization Script for Windows

echo 🚀 Initializing Remnawave Telegram Shop Bot...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go 1.21 or later.
    echo    Download from: https://golang.org/dl/
    pause
    exit /b 1
)

echo ✅ Go is installed

REM Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Docker is not installed. You'll need Docker to run the database.
    echo    Download from: https://www.docker.com/get-started
)

REM Check if Docker Compose is installed
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  Docker Compose is not installed. You'll need it to run the full stack.
    echo    Download from: https://docs.docker.com/compose/install/
)

REM Create .env file if it doesn't exist
if not exist .env (
    echo 📝 Creating .env file from template...
    copy env.example .env
    echo ✅ .env file created. Please edit it with your configuration.
) else (
    echo ✅ .env file already exists
)

REM Download dependencies
echo 📦 Downloading dependencies...
go mod download
go mod tidy

REM Create necessary directories
echo 📁 Creating directories...
if not exist logs mkdir logs
if not exist ssl mkdir ssl
if not exist migrations mkdir migrations

REM Build the application
echo 🔨 Building application...
go build -o remnawave-bot.exe ./cmd/main.go

if %errorlevel% equ 0 (
    echo ✅ Application built successfully
) else (
    echo ❌ Build failed
    pause
    exit /b 1
)

REM Check if PostgreSQL is running
docker --version >nul 2>&1
if %errorlevel% equ 0 (
    echo 🐘 Starting PostgreSQL database...
    docker run -d --name remnawave-postgres -e POSTGRES_DB=remnawave_bot -e POSTGRES_USER=remnawave_bot -e POSTGRES_PASSWORD=remnawave_bot_password -p 5432:5432 postgres:15-alpine
    
    echo ✅ PostgreSQL started
    echo    Database: remnawave_bot
    echo    User: remnawave_bot
    echo    Password: remnawave_bot_password
    echo    Port: 5432
)

echo.
echo 🎉 Initialization complete!
echo.
echo Next steps:
echo 1. Edit .env file with your configuration
echo 2. Run: remnawave-bot.exe
echo    or: go run ./cmd/main.go
echo.
echo For Docker deployment:
echo 1. Edit .env file
echo 2. Run: docker-compose up -d
echo.
echo For development:
echo 1. Install tools: make install-tools
echo 2. Run tests: make test
echo 3. Format code: make fmt
echo.
echo Happy coding! 🚀
pause
