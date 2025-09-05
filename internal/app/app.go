package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"remnawave-tg-shop/internal/bot"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/database"
	"remnawave-tg-shop/internal/logger"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"remnawave-tg-shop/internal/repositories"
	"remnawave-tg-shop/internal/services"
	"remnawave-tg-shop/internal/services/remnawave"
)

// App представляет основное приложение
type App struct {
	config *config.Config
	logger logger.Logger
	db     *database.Database
	bot    *bot.Bot
	server *http.Server
}

// New создает новое приложение
func New(cfg *config.Config, log logger.Logger) *App {
	return &App{
		config: cfg,
		logger: log,
	}
}

// Run запускает приложение
func (a *App) Run() error {
	// Инициализируем базу данных
	db, err := database.New(a.config, a.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.db = db

	// Создаем репозитории
	userRepo := repositories.NewUserRepository(db.DB)
	subscriptionRepo := repositories.NewSubscriptionRepository(db.DB)
	paymentRepo := repositories.NewPaymentRepository(db.DB)

	// Создаем клиент Remnawave
	remnawaveClient := remnawave.NewClient(
		a.config.Remnawave.APIURL,
		a.config.Remnawave.APIKey,
		a.config.Remnawave.SecretKey,
	)

	// Создаем сервисы
	userService := services.NewUserService(userRepo, remnawaveClient, a.logger, a.config)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, remnawaveClient, a.logger)
	paymentService := services.NewPaymentService(paymentRepo, userService, a.logger)

	// Создаем бота
	telegramBot, err := bot.NewBot(a.config, a.logger, userService, subscriptionService, paymentService)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}
	a.bot = telegramBot

	// Настраиваем HTTP сервер для дополнительных endpoints
	if err := a.setupHTTPServer(); err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	// Запускаем HTTP сервер
	go func() {
		a.logger.Info("Starting HTTP server", "port", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server failed", "error", err)
		}
	}()

	// Запускаем бота в отдельной горутине
	go func() {
		if err := a.bot.Start(); err != nil {
			a.logger.Error("Failed to start bot", "error", err)
		}
	}()

	a.logger.Info("Application started successfully")

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down application...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Error("HTTP server shutdown failed", "error", err)
		}
	}

	// Закрываем базу данных
	if err := a.db.Close(); err != nil {
		a.logger.Error("Database close failed", "error", err)
	}

	a.logger.Info("Application stopped")
	return nil
}

// setupHTTPServer настраивает HTTP сервер для webhook'ов
func (a *App) setupHTTPServer() error {
	// Настраиваем Gin
	if a.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// Webhook endpoints
	router.POST("/webhook", a.handleTelegramWebhook)
	router.POST("/tribute-webhook", a.handleTributeWebhook)
	router.POST("/yookassa-webhook", a.handleYooKassaWebhook)

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
		Handler: router,
	}

	return nil
}

// handleTelegramWebhook обрабатывает webhook от Telegram
func (a *App) handleTelegramWebhook(c *gin.Context) {
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		a.logger.Error("Failed to parse Telegram webhook", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Передаем обновление в бот
	if err := a.bot.HandleUpdate(update); err != nil {
		a.logger.Error("Failed to handle Telegram update", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleTributeWebhook обрабатывает webhook от Tribute
func (a *App) handleTributeWebhook(c *gin.Context) {
	a.logger.Info("Received Tribute webhook")
	
	// Здесь должна быть обработка webhook'а от Tribute
	// Пример структуры данных от Tribute:
	/*
	var tributeData struct {
		TransactionID string  `json:"transaction_id"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		Status        string  `json:"status"`
		UserID        string  `json:"user_id"`
	}
	
	if err := c.ShouldBindJSON(&tributeData); err != nil {
		a.logger.Error("Failed to parse Tribute webhook", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	
	// Обработка платежа
	// a.paymentService.ProcessTributeWebhook(tributeData)
	*/
	
	// Пока что просто возвращаем OK
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// handleYooKassaWebhook обрабатывает webhook от ЮKassa
func (a *App) handleYooKassaWebhook(c *gin.Context) {
	a.logger.Info("Received YooKassa webhook")
	
	// Здесь должна быть обработка webhook'а от ЮKassa
	// Пример структуры данных от ЮKassa:
	/*
	var yookassaData struct {
		Type   string `json:"type"`
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Amount struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			Metadata map[string]string `json:"metadata"`
		} `json:"object"`
	}
	
	if err := c.ShouldBindJSON(&yookassaData); err != nil {
		a.logger.Error("Failed to parse YooKassa webhook", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	
	// Обработка платежа
	// a.paymentService.ProcessYooKassaWebhook(yookassaData)
	*/
	
	// Пока что просто возвращаем OK
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}