# Руководство по разработке

## 🚀 Начало работы

### Требования

- **Go 1.21+**
- **Docker & Docker Compose**
- **PostgreSQL 15+**
- **Git**

### Установка

```bash
# Клонируем репозиторий
git clone https://github.com/your-username/remnawave-tg-shop.git
cd remnawave-tg-shop

# Устанавливаем зависимости
go mod download

# Копируем конфигурацию
cp configs/development.env .env

# Запускаем базу данных
docker-compose up -d postgres

# Запускаем приложение
go run cmd/main.go
```

## 🏗️ Архитектура

### Структура проекта

```
remnawave-tg-shop/
├── cmd/                    # Точка входа
│   └── main.go
├── internal/               # Внутренние пакеты
│   ├── app/               # Основное приложение
│   ├── bot/               # Telegram бот
│   ├── config/            # Конфигурация
│   ├── database/          # База данных
│   ├── handlers/          # HTTP обработчики
│   ├── logger/            # Логирование
│   ├── models/            # Модели данных
│   ├── repositories/      # Репозитории
│   └── services/          # Бизнес-логика
├── migrations/            # Миграции БД
├── docs/                  # Документация
├── examples/              # Примеры
├── scripts/               # Скрипты
└── tests/                 # Тесты
```

### Принципы архитектуры

1. **Clean Architecture** - разделение на слои
2. **SOLID принципы** - гибкость и расширяемость
3. **Dependency Injection** - слабая связанность
4. **Interface Segregation** - четкие интерфейсы
5. **Single Responsibility** - одна ответственность

## 🔧 Разработка

### Создание новой функции

#### 1. Создание модели

```go
// internal/models/feature.go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Feature struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name        string    `gorm:"size:255;not null" json:"name"`
    Description string    `gorm:"type:text" json:"description"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

#### 2. Создание репозитория

```go
// internal/repositories/feature_repository.go
package repositories

import (
    "fmt"
    "remnawave-tg-shop/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type FeatureRepository interface {
    Create(feature *models.Feature) error
    GetByID(id uuid.UUID) (*models.Feature, error)
    List(limit, offset int) ([]models.Feature, error)
    Update(feature *models.Feature) error
    Delete(id uuid.UUID) error
}

type featureRepository struct {
    db *gorm.DB
}

func NewFeatureRepository(db *gorm.DB) FeatureRepository {
    return &featureRepository{db: db}
}

func (r *featureRepository) Create(feature *models.Feature) error {
    if err := r.db.Create(feature).Error; err != nil {
        return fmt.Errorf("failed to create feature: %w", err)
    }
    return nil
}

// ... остальные методы
```

#### 3. Создание сервиса

```go
// internal/services/feature_service.go
package services

import (
    "fmt"
    "remnawave-tg-shop/internal/logger"
    "remnawave-tg-shop/internal/models"
    "remnawave-tg-shop/internal/repositories"
    "github.com/google/uuid"
)

type FeatureService interface {
    CreateFeature(name, description string) (*models.Feature, error)
    GetFeature(id uuid.UUID) (*models.Feature, error)
    ListFeatures(limit, offset int) ([]models.Feature, error)
    UpdateFeature(id uuid.UUID, name, description string) error
    DeleteFeature(id uuid.UUID) error
}

type featureService struct {
    featureRepo repositories.FeatureRepository
    logger      logger.Logger
}

func NewFeatureService(featureRepo repositories.FeatureRepository, log logger.Logger) FeatureService {
    return &featureService{
        featureRepo: featureRepo,
        logger:      log,
    }
}

func (s *featureService) CreateFeature(name, description string) (*models.Feature, error) {
    feature := &models.Feature{
        Name:        name,
        Description: description,
        IsActive:    true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.featureRepo.Create(feature); err != nil {
        return nil, fmt.Errorf("failed to create feature: %w", err)
    }

    s.logger.Info("Feature created", "id", feature.ID, "name", name)
    return feature, nil
}

// ... остальные методы
```

#### 4. Создание HTTP обработчика

```go
// internal/handlers/feature_handler.go
package handlers

import (
    "net/http"
    "strconv"
    "remnawave-tg-shop/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type FeatureHandler struct {
    featureService services.FeatureService
}

func NewFeatureHandler(featureService services.FeatureService) *FeatureHandler {
    return &FeatureHandler{featureService: featureService}
}

func (h *FeatureHandler) CreateFeature(c *gin.Context) {
    var request struct {
        Name        string `json:"name" binding:"required"`
        Description string `json:"description"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    feature, err := h.featureService.CreateFeature(request.Name, request.Description)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, feature)
}

func (h *FeatureHandler) GetFeature(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    feature, err := h.featureService.GetFeature(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Feature not found"})
        return
    }

    c.JSON(http.StatusOK, feature)
}

// ... остальные обработчики
```

#### 5. Регистрация в приложении

```go
// internal/app/app.go
func (a *App) setupRoutes() {
    // ... существующие маршруты

    // Регистрируем новый сервис
    featureRepo := repositories.NewFeatureRepository(a.db.DB)
    featureService := services.NewFeatureService(featureRepo, a.logger)
    featureHandler := handlers.NewFeatureHandler(featureService)

    // Добавляем маршруты
    api := a.router.Group("/api/v1")
    {
        features := api.Group("/features")
        {
            features.POST("", featureHandler.CreateFeature)
            features.GET("/:id", featureHandler.GetFeature)
            features.GET("", featureHandler.ListFeatures)
            features.PUT("/:id", featureHandler.UpdateFeature)
            features.DELETE("/:id", featureHandler.DeleteFeature)
        }
    }
}
```

### Создание тестов

#### Unit тесты

```go
// internal/services/feature_service_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestFeatureService_CreateFeature(t *testing.T) {
    // Arrange
    mockRepo := new(MockFeatureRepository)
    mockLogger := new(MockLogger)
    service := NewFeatureService(mockRepo, mockLogger)

    name := "Test Feature"
    description := "Test Description"

    mockRepo.On("Create", mock.AnythingOfType("*models.Feature")).Return(nil)
    mockLogger.On("Info", "Feature created", "id", mock.Anything, "name", name).Return()

    // Act
    feature, err := service.CreateFeature(name, description)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, feature)
    assert.Equal(t, name, feature.Name)
    assert.Equal(t, description, feature.Description)
    assert.True(t, feature.IsActive)

    mockRepo.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}
```

#### Integration тесты

```go
// internal/services/feature_service_integration_test.go
package services

import (
    "testing"
    "remnawave-tg-shop/internal/database"
    "remnawave-tg-shop/internal/repositories"
    "github.com/stretchr/testify/assert"
)

func TestFeatureService_Integration(t *testing.T) {
    // Настройка тестовой БД
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Создание сервиса
    featureRepo := repositories.NewFeatureRepository(db)
    service := NewFeatureService(featureRepo, mockLogger)

    // Тестирование
    feature, err := service.CreateFeature("Test Feature", "Test Description")
    assert.NoError(t, err)
    assert.NotNil(t, feature)

    // Проверяем, что данные сохранены в БД
    retrieved, err := service.GetFeature(feature.ID)
    assert.NoError(t, err)
    assert.Equal(t, feature.Name, retrieved.Name)
}
```

### Создание миграций

```sql
-- migrations/002_add_features_table.sql
CREATE TABLE IF NOT EXISTS features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_features_name ON features(name);
CREATE INDEX IF NOT EXISTS idx_features_is_active ON features(is_active);
CREATE INDEX IF NOT EXISTS idx_features_deleted_at ON features(deleted_at);

CREATE TRIGGER update_features_updated_at BEFORE UPDATE ON features
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
make test

# Конкретный пакет
go test ./internal/services/...

# С покрытием
make test-coverage

# Бенчмарки
make bench
```

### Написание тестов

#### 1. Unit тесты

```go
func TestFunction(t *testing.T) {
    // Arrange - настройка
    input := "test input"
    expected := "expected output"

    // Act - выполнение
    result := function(input)

    // Assert - проверка
    assert.Equal(t, expected, result)
}
```

#### 2. Integration тесты

```go
func TestIntegration(t *testing.T) {
    // Настройка
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Тестирование
    // ...
}
```

#### 3. E2E тесты

```go
func TestE2E(t *testing.T) {
    // Запуск сервера
    server := startTestServer(t)
    defer server.Close()

    // Выполнение запросов
    // ...
}
```

## 📝 Логирование

### Использование логгера

```go
// В сервисах
func (s *service) SomeMethod() {
    s.logger.Info("Method called", "param", value)
    s.logger.Error("Error occurred", "error", err)
    s.logger.Debug("Debug info", "data", data)
}

// В обработчиках
func (h *handler) SomeHandler(c *gin.Context) {
    h.logger.Info("Request received", "path", c.Request.URL.Path)
    // ...
}
```

### Уровни логирования

- **Debug** - отладочная информация
- **Info** - общая информация
- **Warn** - предупреждения
- **Error** - ошибки
- **Fatal** - критические ошибки

## 🔧 Конфигурация

### Добавление новых параметров

#### 1. Обновление структуры конфигурации

```go
// internal/config/config.go
type Config struct {
    // ... существующие поля
    
    // Новые параметры
    NewFeature NewFeatureConfig
}

type NewFeatureConfig struct {
    Enabled bool
    APIKey  string
    URL     string
}
```

#### 2. Загрузка параметров

```go
func Load() (*Config, error) {
    // ... существующий код
    
    // Новые параметры
    cfg.NewFeature.Enabled = getEnvAsBool("NEW_FEATURE_ENABLED", false)
    cfg.NewFeature.APIKey = getEnv("NEW_FEATURE_API_KEY", "")
    cfg.NewFeature.URL = getEnv("NEW_FEATURE_URL", "")
    
    // ... остальной код
}
```

#### 3. Валидация

```go
func (c *Config) Validate() error {
    // ... существующие проверки
    
    if c.NewFeature.Enabled && c.NewFeature.APIKey == "" {
        return fmt.Errorf("NEW_FEATURE_API_KEY is required when NEW_FEATURE_ENABLED is true")
    }
    
    return nil
}
```

## 🚀 Развертывание

### Локальная разработка

```bash
# Запуск базы данных
docker-compose up -d postgres

# Запуск приложения
go run cmd/main.go

# Запуск всех сервисов
docker-compose up -d
```

### Продакшен

```bash
# Сборка образа
docker build -t remnawave-bot .

# Запуск контейнера
docker run -d --name remnawave-bot \
  --env-file .env \
  -p 8080:8080 \
  remnawave-bot
```

## 📚 Документация

### Обновление документации

1. **API документация** - обновляйте `docs/API.md`
2. **Руководство пользователя** - обновляйте `docs/USAGE.md`
3. **Руководство по развертыванию** - обновляйте `docs/DEPLOYMENT.md`
4. **Комментарии в коде** - используйте Go doc

### Комментарии в коде

```go
// Package services provides business logic for the application.
package services

// UserService defines the interface for user-related operations.
type UserService interface {
    // CreateOrGetUser creates a new user or returns existing one.
    // It takes telegram ID, username, first name, last name and language code.
    CreateOrGetUser(telegramID int64, username, firstName, lastName, languageCode string) (*models.User, error)
}
```

## 🔍 Отладка

### Логирование

```go
// Включение debug режима
log.SetLevel(logrus.DebugLevel)

// Логирование с контекстом
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "create_subscription",
}).Info("Creating subscription")
```

### Профилирование

```bash
# CPU профилирование
go test -cpuprofile=cpu.prof ./...

# Memory профилирование
go test -memprofile=mem.prof ./...

# Анализ профиля
go tool pprof cpu.prof
```

## 🚨 Обработка ошибок

### Создание кастомных ошибок

```go
// internal/errors/errors.go
package errors

import "fmt"

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
    return &ValidationError{Field: field, Message: message}
}
```

### Обработка ошибок в сервисах

```go
func (s *service) SomeMethod() error {
    if err := s.validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if err := s.process(); err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }
    
    return nil
}
```

## 📊 Мониторинг

### Метрики

```go
// internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### Health checks

```go
// internal/health/health.go
package health

type Checker interface {
    Check() error
}

type HealthChecker struct {
    checkers map[string]Checker
}

func (h *HealthChecker) Check() map[string]error {
    results := make(map[string]error)
    for name, checker := range h.checkers {
        results[name] = checker.Check()
    }
    return results
}
```

## 🔄 CI/CD

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - name: Run tests
      run: go test ./...
    - name: Run linter
      run: golangci-lint run
    - name: Build
      run: go build ./cmd/main.go
```

## 📝 Лучшие практики

### 1. Код

- Используйте `gofmt` для форматирования
- Следуйте Go conventions
- Пишите тесты для нового кода
- Используйте интерфейсы для тестируемости

### 2. Git

- Используйте осмысленные commit сообщения
- Создавайте feature branches
- Делайте code review
- Используйте conventional commits

### 3. Документация

- Обновляйте README при изменениях
- Документируйте API
- Пишите комментарии к коду
- Ведите changelog

---

**Удачной разработки! 🚀**
