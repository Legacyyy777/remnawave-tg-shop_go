package database

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Database представляет подключение к базе данных
type Database struct {
	DB     *gorm.DB
	config *config.Config
	logger logger.Logger
}

// New создает новое подключение к базе данных
func New(cfg *config.Config, log logger.Logger) (*Database, error) {
	// Формируем DSN для PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Настраиваем GORM логгер
	gormLog := gormLogger.Default.LogMode(gormLogger.Silent)
	if cfg.LogLevel == "debug" {
		gormLog = gormLogger.Default.LogMode(gormLogger.Info)
	}

	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLog,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	database := &Database{
		DB:     db,
		config: cfg,
		logger: log,
	}

	// Выполняем миграции
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return database, nil
}

// migrate выполняет миграции базы данных
func (d *Database) migrate() error {
	d.logger.Info("Running database migrations...")

	// Автомиграция моделей
	if err := d.DB.AutoMigrate(
		&models.User{},
		&models.Subscription{},
		&models.Payment{},
		&models.Server{},
		&models.Plan{},
	); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	d.logger.Info("Database migrations completed successfully")
	return nil
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Health проверяет состояние базы данных
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
