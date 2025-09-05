# Руководство по тестированию

## 🧪 Обзор

Этот документ описывает процесс тестирования Remnawave Telegram Shop Bot.

## 🚀 Быстрый старт

### Запуск тестов

```bash
# Все тесты
make test

# Тесты с покрытием
make test-coverage

# Бенчмарки
make bench

# Конкретный пакет
go test ./internal/services/...

# Конкретный тест
go test -run TestUserService_CreateOrGetUser ./internal/services/
```

## 📋 Типы тестов

### Unit тесты

Тестируют отдельные функции и методы.

```go
func TestUserService_CreateOrGetUser(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    mockLogger := new(MockLogger)
    service := NewUserService(mockRepo, mockLogger)

    // Act
    user, err := service.CreateOrGetUser(123456789, "user", "John", "Doe", "ru")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, int64(123456789), user.TelegramID)
}
```

### Integration тесты

Тестируют взаимодействие между компонентами.

```go
func TestUserService_Integration(t *testing.T) {
    // Настройка тестовой БД
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Создание сервиса с реальной БД
    userRepo := repositories.NewUserRepository(db)
    service := services.NewUserService(userRepo, mockLogger)

    // Тестирование
    user, err := service.CreateOrGetUser(123456789, "user", "John", "Doe", "ru")
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### E2E тесты

Тестируют полный пользовательский сценарий.

```go
func TestE2E_UserRegistrationAndSubscription(t *testing.T) {
    // Запуск тестового сервера
    server := startTestServer(t)
    defer server.Close()

    // Регистрация пользователя
    user := registerUser(t, server)

    // Пополнение баланса
    addBalance(t, server, user.ID, 1000.0)

    // Покупка подписки
    subscription := buySubscription(t, server, user.ID, 1, 1)

    // Проверка результата
    assert.Equal(t, "active", subscription.Status)
}
```

## 🔧 Настройка тестового окружения

### Тестовая база данных

```go
func setupTestDB(t *testing.T) *gorm.DB {
    // Создаем временную БД
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // Применяем миграции
    err = db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Payment{})
    require.NoError(t, err)

    return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
    sqlDB, err := db.DB()
    require.NoError(t, err)
    sqlDB.Close()
}
```

### Моки

```go
// Mock для UserRepository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
    args := m.Called(telegramID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}
```

### Тестовые данные

```go
func createTestUser() *models.User {
    return &models.User{
        ID:           uuid.New(),
        TelegramID:   123456789,
        Username:     "testuser",
        FirstName:    "Test",
        LastName:     "User",
        LanguageCode: "ru",
        IsBlocked:    false,
        IsAdmin:      false,
        Balance:      0,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
}

func createTestSubscription() *models.Subscription {
    return &models.Subscription{
        ID:         uuid.New(),
        UserID:     uuid.New(),
        ServerID:   1,
        ServerName: "Test Server",
        PlanID:     1,
        PlanName:   "Test Plan",
        Status:     "active",
        ExpiresAt:  time.Now().AddDate(0, 0, 30),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
}
```

## 📊 Покрытие кода

### Генерация отчета

```bash
# Генерируем отчет о покрытии
go test -coverprofile=coverage.out ./...

# Просматриваем в браузере
go tool cover -html=coverage.out

# Текстовый отчет
go tool cover -func=coverage.out
```

### Настройка покрытия

```go
// В .github/workflows/test.yml
- name: Run tests with coverage
  run: |
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
```

## 🚀 Performance тесты

### Бенчмарки

```go
func BenchmarkUserService_CreateOrGetUser(b *testing.B) {
    mockRepo := new(MockUserRepository)
    mockLogger := new(MockLogger)
    service := NewUserService(mockRepo, mockLogger)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.CreateOrGetUser(int64(i), "user", "John", "Doe", "ru")
    }
}

func BenchmarkDatabase_Query(b *testing.B) {
    db := setupTestDB(b)
    defer cleanupTestDB(b, db)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        var user models.User
        db.First(&user, "telegram_id = ?", 123456789)
    }
}
```

### Load тесты

```go
func TestLoad_ConcurrentUsers(t *testing.T) {
    const numUsers = 1000
    const numGoroutines = 10

    var wg sync.WaitGroup
    errors := make(chan error, numUsers)

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < numUsers/numGoroutines; j++ {
                user, err := service.CreateOrGetUser(int64(j), "user", "John", "Doe", "ru")
                if err != nil {
                    errors <- err
                }
                if user == nil {
                    errors <- fmt.Errorf("user is nil")
                }
            }
        }()
    }

    wg.Wait()
    close(errors)

    for err := range errors {
        t.Error(err)
    }
}
```

## 🔍 Тестирование API

### HTTP тесты

```go
func TestAPI_GetUsers(t *testing.T) {
    // Настройка тестового сервера
    server := httptest.NewServer(setupTestRouter())
    defer server.Close()

    // Создание тестовых данных
    createTestUsers(t, 5)

    // Выполнение запроса
    resp, err := http.Get(server.URL + "/api/v1/users")
    require.NoError(t, err)
    defer resp.Body.Close()

    // Проверка ответа
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var response struct {
        Users []models.User `json:"users"`
        Total int           `json:"total"`
    }
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(t, err)
    assert.Len(t, response.Users, 5)
    assert.Equal(t, 5, response.Total)
}
```

### Webhook тесты

```go
func TestWebhook_Telegram(t *testing.T) {
    server := httptest.NewServer(setupTestRouter())
    defer server.Close()

    // Тестовое сообщение от Telegram
    update := tgbotapi.Update{
        UpdateID: 123456789,
        Message: &tgbotapi.Message{
            MessageID: 1,
            From: &tgbotapi.User{
                ID:           123456789,
                IsBot:        false,
                FirstName:    "John",
                LastName:     "Doe",
                UserName:     "johndoe",
                LanguageCode: "ru",
            },
            Chat: &tgbotapi.Chat{
                ID:       123456789,
                Type:     "private",
                FirstName: "John",
                LastName:  "Doe",
                UserName:  "johndoe",
            },
            Date: 1640995200,
            Text: "/start",
        },
    }

    // Отправка webhook
    body, _ := json.Marshal(update)
    resp, err := http.Post(server.URL+"/webhook", "application/json", bytes.NewBuffer(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

## 🧪 Тестирование базы данных

### Тесты миграций

```go
func TestMigrations(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Применяем миграции
    err := db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Payment{})
    require.NoError(t, err)

    // Проверяем, что таблицы созданы
    var tables []string
    db.Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables)
    
    assert.Contains(t, tables, "users")
    assert.Contains(t, tables, "subscriptions")
    assert.Contains(t, tables, "payments")
}
```

### Тесты транзакций

```go
func TestTransaction_CreateUserAndSubscription(t *testing.T) {
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Начинаем транзакцию
    tx := db.Begin()
    defer tx.Rollback()

    // Создаем пользователя
    user := &models.User{
        TelegramID: 123456789,
        Username:   "testuser",
        FirstName:  "Test",
        LastName:   "User",
    }
    err := tx.Create(user).Error
    require.NoError(t, err)

    // Создаем подписку
    subscription := &models.Subscription{
        UserID:     user.ID,
        ServerID:   1,
        ServerName: "Test Server",
        PlanID:     1,
        PlanName:   "Test Plan",
        Status:     "active",
        ExpiresAt:  time.Now().AddDate(0, 0, 30),
    }
    err = tx.Create(subscription).Error
    require.NoError(t, err)

    // Подтверждаем транзакцию
    err = tx.Commit().Error
    require.NoError(t, err)

    // Проверяем, что данные сохранены
    var count int64
    db.Model(&models.User{}).Count(&count)
    assert.Equal(t, int64(1), count)

    db.Model(&models.Subscription{}).Count(&count)
    assert.Equal(t, int64(1), count)
}
```

## 🔐 Тестирование безопасности

### Тесты аутентификации

```go
func TestAuth_ValidToken(t *testing.T) {
    // Создаем валидный токен
    token := createTestToken(t, "admin")

    // Выполняем запрос с токеном
    req, _ := http.NewRequest("GET", "/api/v1/users", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
}

func TestAuth_InvalidToken(t *testing.T) {
    // Выполняем запрос без токена
    req, _ := http.NewRequest("GET", "/api/v1/users", nil)

    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
```

### Тесты валидации

```go
func TestValidation_InvalidUserData(t *testing.T) {
    server := httptest.NewServer(setupTestRouter())
    defer server.Close()

    // Отправляем невалидные данные
    userData := map[string]interface{}{
        "telegram_id": "invalid",
        "username":    "",
        "first_name":  "",
    }

    body, _ := json.Marshal(userData)
    resp, err := http.Post(server.URL+"/api/v1/users", "application/json", bytes.NewBuffer(body))
    require.NoError(t, err)
    defer resp.Body.Close()

    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Contains(t, response, "error")
}
```

## 📱 Тестирование Telegram бота

### Мокирование Telegram API

```go
func TestBot_HandleMessage(t *testing.T) {
    // Создаем мок для Telegram API
    mockAPI := &MockTelegramAPI{}
    bot := &Bot{api: mockAPI}

    // Настраиваем ожидания
    mockAPI.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).Return(tgbotapi.Message{}, nil)

    // Создаем тестовое сообщение
    message := &tgbotapi.Message{
        MessageID: 1,
        From: &tgbotapi.User{
            ID:        123456789,
            FirstName: "John",
            LastName:  "Doe",
            UserName:  "johndoe",
        },
        Chat: &tgbotapi.Chat{
            ID:   123456789,
            Type: "private",
        },
        Date: 1640995200,
        Text: "/start",
    }

    // Обрабатываем сообщение
    bot.handleMessage(message)

    // Проверяем, что Send был вызван
    mockAPI.AssertExpectations(t)
}
```

## 🚀 CI/CD тесты

### GitHub Actions

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Run tests with coverage
      run: go test -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

## 📊 Метрики тестирования

### Отчеты о покрытии

```bash
# Генерируем отчет
go test -coverprofile=coverage.out ./...

# Просматриваем в браузере
go tool cover -html=coverage.out -o coverage.html

# Текстовый отчет
go tool cover -func=coverage.out | grep total
```

### Производительность

```bash
# Бенчмарки
go test -bench=. ./...

# Профилирование
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...

# Анализ профиля
go tool pprof cpu.prof
go tool pprof mem.prof
```

## 🐛 Отладка тестов

### Включение verbose режима

```bash
# Подробный вывод
go test -v ./...

# Очень подробный вывод
go test -v -count=1 ./...
```

### Логирование в тестах

```go
func TestWithLogging(t *testing.T) {
    // Включаем логирование
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    // Тест
    user, err := service.CreateOrGetUser(123456789, "user", "John", "Doe", "ru")
    if err != nil {
        t.Logf("Error: %v", err)
    }
    t.Logf("User: %+v", user)
}
```

## 📚 Лучшие практики

### 1. Именование тестов

```go
// Хорошо
func TestUserService_CreateOrGetUser_NewUser(t *testing.T) {}
func TestUserService_CreateOrGetUser_ExistingUser(t *testing.T) {}

// Плохо
func Test1(t *testing.T) {}
func TestUser(t *testing.T) {}
```

### 2. Структура тестов

```go
func TestUserService_CreateOrGetUser(t *testing.T) {
    // Arrange - настройка
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo, mockLogger)

    // Act - выполнение
    user, err := service.CreateOrGetUser(123456789, "user", "John", "Doe", "ru")

    // Assert - проверка
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 3. Очистка ресурсов

```go
func TestWithCleanup(t *testing.T) {
    // Настройка
    db := setupTestDB(t)
    
    // Очистка
    defer cleanupTestDB(t, db)
    
    // Тест
    // ...
}
```

### 4. Параллельные тесты

```go
func TestParallel(t *testing.T) {
    t.Parallel()
    
    // Тест
    // ...
}
```

---

**Удачного тестирования! 🧪**
