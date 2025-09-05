package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	// Telegram Bot
	BotToken       string
	BotWebhookURL  string
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

	// Trial Settings
	Trial TrialConfig

	// Mini App
	MiniApp MiniAppConfig

	// Referral System
	Referral ReferralConfig

	// Notifications
	Notifications NotificationConfig

	// Promo Codes
	PromoCodes PromoCodeConfig

	// Localization
	Localization LocalizationConfig

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
	APIURL    string
	APIKey    string
	SecretKey string
}

type PaymentConfig struct {
	Tribute  TributeConfig
	YooKassa YooKassaConfig

	// Payment Method Toggles
	StarsEnabled     bool
	TributeEnabled   bool
	YooKassaEnabled  bool
	CryptoPayEnabled bool

	// Subscription Prices
	Price1Month   int
	Price3Months  int
	Price6Months  int
	Price12Months int
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
	TelegramIDs           []int64
	MaintenanceMode       bool
	MaintenanceAutoEnable bool
}

type SecurityConfig struct {
	JWTSecret     string
	EncryptionKey string
}

type MonitoringConfig struct {
	HealthCheckInterval  time.Duration
	StatsCleanupInterval time.Duration
}

type TrialConfig struct {
	Enabled         bool
	DurationDays    int
	TrafficLimitGB  int
	TrafficStrategy string
}

type MiniAppConfig struct {
	URL string
}

// ReferralConfig настройки реферальной системы
type ReferralConfig struct {
	Enabled       bool
	BonusDays     int
	ReferrerBonus int
	ReferredBonus int
}

// NotificationConfig настройки уведомлений
type NotificationConfig struct {
	Enabled            bool
	ExpiringDaysBefore int
	CheckInterval      time.Duration
	MaxRetries         int
}

// PromoCodeConfig настройки промокодов
type PromoCodeConfig struct {
	Enabled       bool
	MaxCodeLength int
	MinCodeLength int
}

// LocalizationConfig настройки локализации
type LocalizationConfig struct {
	DefaultLanguage    string
	SupportedLanguages []string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Отладка загрузки конфигурации
	fmt.Printf("CONFIG DEBUG: Starting config load...\n")

	// Загружаем .env файл если он существует
	if err := godotenv.Load(); err != nil {
		fmt.Printf("CONFIG DEBUG: .env file not loaded: %v\n", err)
		// Игнорируем ошибку если файл не найден
	} else {
		fmt.Printf("CONFIG DEBUG: .env file loaded successfully\n")
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

	// Payment Method Toggles
	cfg.Payments.StarsEnabled = getEnvAsBool("STARS_ENABLED", true)
	cfg.Payments.TributeEnabled = getEnvAsBool("TRIBUTE_ENABLED", true)
	cfg.Payments.YooKassaEnabled = getEnvAsBool("YOOKASSA_ENABLED", true)
	cfg.Payments.CryptoPayEnabled = getEnvAsBool("CRYPTOPAY_ENABLED", false)

	// Subscription Prices
	cfg.Payments.Price1Month = getEnvAsInt("RUB_PRICE_1_MONTH", 150)
	cfg.Payments.Price3Months = getEnvAsInt("RUB_PRICE_3_MONTHS", 300)
	cfg.Payments.Price6Months = getEnvAsInt("RUB_PRICE_6_MONTHS", 500)
	cfg.Payments.Price12Months = getEnvAsInt("RUB_PRICE_12_MONTHS", 900)

	// Server
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)

	// Admin
	cfg.Admin.TelegramIDs = getEnvAsInt64Slice("ADMIN_TELEGRAM_IDS", []int64{})
	cfg.Admin.MaintenanceMode = getEnvAsBool("MAINTENANCE_MODE", false)
	cfg.Admin.MaintenanceAutoEnable = getEnvAsBool("MAINTENANCE_AUTO_ENABLE", true)

	// Отладочная информация для админа
	fmt.Printf("DEBUG: ADMIN_TELEGRAM_IDS loaded: %v\n", cfg.Admin.TelegramIDs)
	fmt.Printf("DEBUG: Environment variables check:\n")
	fmt.Printf("  ADMIN_TELEGRAM_IDS env var: '%s'\n", os.Getenv("ADMIN_TELEGRAM_IDS"))

	// Security
	cfg.Security.JWTSecret = getEnv("JWT_SECRET", "")
	cfg.Security.EncryptionKey = getEnv("ENCRYPTION_KEY", "")

	// Monitoring
	cfg.Monitoring.HealthCheckInterval = getEnvAsDuration("HEALTH_CHECK_INTERVAL", "30s")
	cfg.Monitoring.StatsCleanupInterval = getEnvAsDuration("STATS_CLEANUP_INTERVAL", "24h")

	// Trial Settings
	cfg.Trial.Enabled = getEnvAsBool("TRIAL_ENABLED", true)
	cfg.Trial.DurationDays = getEnvAsInt("TRIAL_DURATION_DAYS", 5)
	cfg.Trial.TrafficLimitGB = getEnvAsInt("TRIAL_TRAFFIC_LIMIT_GB", 0)
	cfg.Trial.TrafficStrategy = getEnv("TRIAL_TRAFFIC_STRATEGY", "NO_RESET")

	// Mini App
	cfg.MiniApp.URL = getEnv("SUBSCRIPTION_MINI_APP_URL", "")
	fmt.Printf("DEBUG: SUBSCRIPTION_MINI_APP_URL loaded: '%s'\n", cfg.MiniApp.URL)
	fmt.Printf("DEBUG: MiniApp.URL length: %d\n", len(cfg.MiniApp.URL))

	// Referral System
	cfg.Referral.Enabled = getEnvAsBool("REFERRAL_ENABLED", true)
	cfg.Referral.BonusDays = getEnvAsInt("REFERRAL_BONUS_DAYS", 7)
	cfg.Referral.ReferrerBonus = getEnvAsInt("REFERRAL_REFERRER_BONUS", 50)
	cfg.Referral.ReferredBonus = getEnvAsInt("REFERRAL_REFERRED_BONUS", 30)

	// Notifications
	cfg.Notifications.Enabled = getEnvAsBool("NOTIFICATIONS_ENABLED", true)
	cfg.Notifications.ExpiringDaysBefore = getEnvAsInt("NOTIFICATIONS_EXPIRING_DAYS_BEFORE", 3)
	cfg.Notifications.CheckInterval = getEnvAsDuration("NOTIFICATIONS_CHECK_INTERVAL", "1h")
	cfg.Notifications.MaxRetries = getEnvAsInt("NOTIFICATIONS_MAX_RETRIES", 3)

	// Promo Codes
	cfg.PromoCodes.Enabled = getEnvAsBool("PROMO_CODES_ENABLED", true)
	cfg.PromoCodes.MaxCodeLength = getEnvAsInt("PROMO_CODES_MAX_LENGTH", 20)
	cfg.PromoCodes.MinCodeLength = getEnvAsInt("PROMO_CODES_MIN_LENGTH", 3)

	// Localization
	cfg.Localization.DefaultLanguage = getEnv("DEFAULT_LANGUAGE", "ru")
	cfg.Localization.SupportedLanguages = getEnvAsStringSlice("SUPPORTED_LANGUAGES", []string{"ru", "en"})

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
	if c.MiniApp.URL == "" {
		return fmt.Errorf("SUBSCRIPTION_MINI_APP_URL is required")
	}
	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	fmt.Printf("DEBUG getEnv: key='%s', value='%s', empty=%t\n", key, value, value == "")
	if value != "" {
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
	value := os.Getenv(key)
	fmt.Printf("DEBUG getEnvAsInt64: key='%s', value='%s', empty=%t\n", key, value, value == "")

	if value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			fmt.Printf("DEBUG getEnvAsInt64: parsed successfully: %d\n", intValue)
			return intValue
		} else {
			fmt.Printf("DEBUG getEnvAsInt64: parse error: %v\n", err)
		}
	}
	fmt.Printf("DEBUG getEnvAsInt64: using default value: %d\n", defaultValue)
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

func getEnvAsInt64Slice(key string, defaultValue []int64) []int64 {
	if value := os.Getenv(key); value != "" {
		// Парсим строку вида "123,456,789"
		parts := strings.Split(value, ",")
		result := make([]int64, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if id, err := strconv.ParseInt(part, 10, 64); err == nil {
				result = append(result, id)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Парсим строку вида "ru,en,de"
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
