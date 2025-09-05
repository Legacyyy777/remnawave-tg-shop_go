# API Examples

## Remnawave API Client

### Инициализация клиента

```go
package main

import (
    "fmt"
    "log"
    "remnawave-tg-shop/internal/services/remnawave"
)

func main() {
    // Создаем клиент
    client := remnawave.NewClient(
        "https://your-panel.com/api",
        "your_api_key",
        "secret_name:secret_value", // опционально
    )

    // Получаем список серверов
    servers, err := client.GetServers()
    if err != nil {
        log.Fatal(err)
    }

    for _, server := range servers {
        fmt.Printf("Server: %s (ID: %d)\n", server.Name, server.ID)
    }
}
```

### Получение серверов

```go
// Получить все серверы
servers, err := client.GetServers()
if err != nil {
    log.Fatal(err)
}

for _, server := range servers {
    fmt.Printf("Server: %s\n", server.Name)
    fmt.Printf("Description: %s\n", server.Description)
    fmt.Printf("Active: %t\n", server.IsActive)
}
```

### Получение тарифных планов

```go
// Получить планы для сервера
plans, err := client.GetPlans(serverID)
if err != nil {
    log.Fatal(err)
}

for _, plan := range plans {
    fmt.Printf("Plan: %s\n", plan.Name)
    fmt.Printf("Price: %.2f\n", plan.Price)
    fmt.Printf("Duration: %d days\n", plan.Duration)
}
```

### Создание пользователя

```go
// Создать пользователя в Remnawave
user, err := client.CreateUser("username", "email@example.com")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User created: ID=%d, Username=%s\n", user.ID, user.Username)
```

### Создание подписки

```go
// Создать подписку
subscription, err := client.CreateSubscription(userID, serverID, planID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Subscription created: ID=%d\n", subscription.ID)
fmt.Printf("Expires at: %s\n", subscription.ExpiresAt)
```

### Обновление подписки

```go
// Обновить подписку
data := map[string]interface{}{
    "status": "active",
    "expires_at": "2024-12-31T23:59:59Z",
}

subscription, err := client.UpdateSubscription(subscriptionID, data)
if err != nil {
    log.Fatal(err)
}
```

### Удаление подписки

```go
// Удалить подписку
err := client.DeleteSubscription(subscriptionID)
if err != nil {
    log.Fatal(err)
}
```

## Telegram Bot API

### Обработка команд

```go
// В internal/bot/bot.go
func (b *Bot) handleCommand(message *tgbotapi.Message, user *models.User) {
    command := message.Command()
    
    switch command {
    case "start":
        b.handleStartCommand(message, user, message.CommandArguments())
    case "balance":
        b.handleBalanceCommand(message, user)
    case "subscriptions":
        b.handleSubscriptionsCommand(message, user)
    // ... другие команды
    }
}
```

### Отправка сообщений

```go
// Отправить текстовое сообщение
func (b *Bot) sendMessage(chatID int64, text string) {
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = tgbotapi.ModeHTML
    b.api.Send(msg)
}

// Отправить сообщение с клавиатурой
func (b *Bot) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ReplyMarkup = keyboard
    b.api.Send(msg)
}
```

### Обработка callback query

```go
// Обработать нажатие на inline кнопку
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
    data := query.Data
    
    switch {
    case strings.HasPrefix(data, "balance"):
        b.handleBalanceCallback(query, user)
    case strings.HasPrefix(data, "buy_subscription"):
        b.handleBuySubscriptionCallback(query, user)
    // ... другие обработчики
    }
}
```

## Database Operations

### Создание пользователя

```go
// В internal/services/user_service.go
func (s *userService) CreateOrGetUser(telegramID int64, username, firstName, lastName, languageCode string) (*models.User, error) {
    // Проверяем, существует ли пользователь
    user, err := s.userRepo.GetByTelegramID(telegramID)
    if err != nil {
        return nil, err
    }
    
    if user != nil {
        // Обновляем данные существующего пользователя
        user.Username = username
        user.FirstName = firstName
        user.LastName = lastName
        user.LanguageCode = languageCode
        user.UpdatedAt = time.Now()
        
        return user, s.userRepo.Update(user)
    }
    
    // Создаем нового пользователя
    user = &models.User{
        TelegramID:   telegramID,
        Username:     username,
        FirstName:    firstName,
        LastName:     lastName,
        LanguageCode: languageCode,
        IsBlocked:    false,
        IsAdmin:      false,
        Balance:      0,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    return user, s.userRepo.Create(user)
}
```

### Создание подписки

```go
// В internal/services/subscription_service.go
func (s *subscriptionService) CreateSubscription(userID uuid.UUID, serverID, planID int) (*models.Subscription, error) {
    // Получаем информацию о сервере и плане
    servers, err := s.remnawaveClient.GetServers()
    if err != nil {
        return nil, err
    }
    
    var serverName string
    for _, server := range servers {
        if server.ID == serverID {
            serverName = server.Name
            break
        }
    }
    
    plans, err := s.remnawaveClient.GetPlans(serverID)
    if err != nil {
        return nil, err
    }
    
    var planName string
    var planDuration int
    for _, plan := range plans {
        if plan.ID == planID {
            planName = plan.Name
            planDuration = plan.Duration
            break
        }
    }
    
    // Создаем подписку в Remnawave
    remnawaveSub, err := s.remnawaveClient.CreateSubscription(0, serverID, planID)
    if err != nil {
        return nil, err
    }
    
    // Создаем подписку в нашей БД
    subscription := &models.Subscription{
        UserID:     userID,
        ServerID:   serverID,
        ServerName: serverName,
        PlanID:     planID,
        PlanName:   planName,
        Status:     "active",
        ExpiresAt:  time.Now().AddDate(0, 0, planDuration),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    return subscription, s.subscriptionRepo.Create(subscription)
}
```

## Payment Processing

### Создание платежа

```go
// В internal/services/payment_service.go
func (s *paymentService) CreatePayment(userID uuid.UUID, amount float64, method, description string) (*models.Payment, error) {
    payment := &models.Payment{
        UserID:        userID,
        Amount:        amount,
        Currency:      "RUB",
        PaymentMethod: method,
        Status:        "pending",
        Description:   description,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    
    return payment, s.paymentRepo.Create(payment)
}
```

### Обработка завершенного платежа

```go
func (s *paymentService) UpdatePaymentStatus(id uuid.UUID, status string) error {
    payment, err := s.paymentRepo.GetByID(id)
    if err != nil {
        return err
    }
    
    payment.Status = status
    payment.UpdatedAt = time.Now()
    
    if status == "completed" {
        now := time.Now()
        payment.CompletedAt = &now
        
        // Добавляем средства на баланс
        if err := s.userService.AddBalance(payment.UserID, payment.Amount); err != nil {
            return err
        }
    }
    
    return s.paymentRepo.Update(payment)
}
```

## Error Handling

### Обработка ошибок API

```go
func (c *Client) makeRequest(method, endpoint string, data interface{}, result interface{}) error {
    // ... выполнение запроса
    
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(responseBody))
    }
    
    // ... парсинг ответа
}
```

### Обработка ошибок базы данных

```go
func (r *userRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
    var user models.User
    if err := r.db.First(&user, "telegram_id = ?", telegramID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by Telegram ID: %w", err)
    }
    return &user, nil
}
```

## Configuration

### Загрузка конфигурации

```go
// В internal/config/config.go
func Load() (*Config, error) {
    cfg := &Config{}
    
    // Загружаем переменные окружения
    cfg.BotToken = getEnv("BOT_TOKEN", "")
    cfg.Database.Host = getEnv("DB_HOST", "localhost")
    cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
    // ... другие параметры
    
    // Валидируем конфигурацию
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return cfg, nil
}
```

### Валидация конфигурации

```go
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
    return nil
}
```
