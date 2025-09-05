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
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создаем логгер
	logger := logger.New(cfg.LogLevel)

	// Создаем приложение
	application := app.New(cfg, logger)

	// Запускаем приложение
	if err := application.Run(); err != nil {
		logger.Fatal("Ошибка запуска приложения", "error", err)
		os.Exit(1)
	}
}
