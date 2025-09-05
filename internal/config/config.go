package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	// Telegram Bot
	BotToken      string
	BotWebhookURL string
	BotWebhookPort int

	// Database
	Database DatabaseConfig

	// Remnawave API
	Remnawave RemnawaveConfig

	// Payment Systems
	Payments PaymentConfig

	// Server
	Server ServerConfig

	// Admin
	Admin AdminConfig

	// Security
	Security SecurityConfig

	// Monitoring
	Monitoring MonitoringConfig

	// Environment
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RemnawaveConfig struct {
	APIURL     string
	APIKey     string
	SecretKey  string
}

type PaymentConfig struct {
	Tribute TributeConfig
	YooKassa YooKassaConfig
}

type TributeConfig struct {
	WebhookURL string
}

type YooKassaConfig struct {
	ShopID     string
	SecretKey  string
	WebhookURL string
}

type ServerConfig struct {
	Port int
}

type AdminConfig struct {
	TelegramID           int64
	MaintenanceMode      bool
	MaintenanceAutoEnable bool
}

type SecurityConfig struct {
	JWTSecret     string
	EncryptionKey string
}

type MonitoringConfig struct {
	HealthCheckInterval   time.Duration
	StatsCleanupInterval  time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загружаем .env файл если он существует
	if err := godotenv.Load(); err != nil {
		// Игнорируем ошибку если файл не найден
	}

	cfg := &Config{}

	// Telegram Bot
	cfg.BotToken = getEnv("BOT_TOKEN", "")
	cfg.BotWebhookURL = getEnv("BOT_WEBHOOK_URL", "")
	cfg.BotWebhookPort = getEnvAsInt("BOT_WEBHOOK_PORT", 8080)

	// Database
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "remnawave_bot")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Name = getEnv("DB_NAME", "remnawave_bot")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")

	// Remnawave
	cfg.Remnawave.APIURL = getEnv("REMNAWAVE_API_URL", "")
	cfg.Remnawave.APIKey = getEnv("REMNAWAVE_API_KEY", "")
	cfg.Remnawave.SecretKey = getEnv("REMNAWAVE_SECRET_KEY", "")

	// Payments
	cfg.Payments.Tribute.WebhookURL = getEnv("TRIBUTE_WEBHOOK_URL", "")
	cfg.Payments.YooKassa.ShopID = getEnv("YOOKASSA_SHOP_ID", "")
	cfg.Payments.YooKassa.SecretKey = getEnv("YOOKASSA_SECRET_KEY", "")
	cfg.Payments.YooKassa.WebhookURL = getEnv("YOOKASSA_WEBHOOK_URL", "")

	// Server
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)

	// Admin
	cfg.Admin.TelegramID = getEnvAsInt64("ADMIN_TELEGRAM_ID", 0)
	cfg.Admin.MaintenanceMode = getEnvAsBool("MAINTENANCE_MODE", false)
	cfg.Admin.MaintenanceAutoEnable = getEnvAsBool("MAINTENANCE_AUTO_ENABLE", true)

	// Security
	cfg.Security.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.Security.EncryptionKey = getEnv("ENCRYPTION_KEY", "")

	// Monitoring
	cfg.Monitoring.HealthCheckInterval = getEnvAsDuration("HEALTH_CHECK_INTERVAL", "30s")
	cfg.Monitoring.StatsCleanupInterval = getEnvAsDuration("STATS_CLEANUP_INTERVAL", "24h")

	// Environment
	cfg.Environment = getEnv("ENVIRONMENT", "development")
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")

	// Валидация обязательных параметров
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return fmt.Errorf("BOT_TOKEN is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Remnawave.APIURL == "" {
		return fmt.Errorf("REMNAWAVE_API_URL is required")
	}
	if c.Remnawave.APIKey == "" {
		return fmt.Errorf("REMNAWAVE_API_KEY is required")
	}
	if c.Security.EncryptionKey == "" || len(c.Security.EncryptionKey) != 32 {
		return fmt.Errorf("ENCRYPTION_KEY must be 32 characters long")
	}
	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
