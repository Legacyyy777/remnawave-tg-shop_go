package main

import (
	"log"
	"os"

	"remnawave-tg-shop/internal/app"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем логгер
	logger := logger.New(cfg.LogLevel)

	// Создаем и запускаем приложение
	application := app.New(cfg, logger)
	
	if err := application.Run(); err != nil {
		logger.Fatal("Application failed to start", "error", err)
		os.Exit(1)
	}
}
