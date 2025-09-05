# Bot Structure

## 📁 Directory Structure

```
internal/bot/
├── handlers/                 # Обработчики событий
│   ├── commands/            # Обработчики команд
│   │   ├── start.go         # /start команда
│   │   ├── help.go          # /help команда
│   │   ├── balance.go       # /balance команда
│   │   ├── subscriptions.go # /subscriptions команда
│   │   ├── referrals.go     # /referrals команда
│   │   └── admin.go         # /admin команда
│   ├── callbacks/           # Обработчики callback'ов
│   │   ├── balance.go       # Обработка кнопки баланса
│   │   ├── buy_subscription.go # Обработка покупки подписки
│   │   ├── my_subscriptions.go # Обработка моих подписок
│   │   ├── referrals.go     # Обработка рефералов
│   │   ├── promo_code.go    # Обработка промокодов
│   │   ├── language.go      # Обработка смены языка
│   │   ├── status.go        # Обработка статуса
│   │   ├── support.go       # Обработка поддержки
│   │   ├── trial.go         # Обработка пробного периода
│   │   ├── tariff.go        # Обработка тарифов
│   │   └── payment.go       # Обработка платежей
│   └── messages/            # Обработчики сообщений
│       ├── text.go          # Текстовые сообщения
│       ├── promo_code.go    # Обработка промокодов
│       └── search.go        # Поиск
├── middleware/              # Промежуточное ПО
│   ├── auth.go             # Аутентификация и логирование
│   ├── rate_limit.go       # Ограничение частоты запросов
│   └── admin.go            # Проверка прав администратора
├── keyboards/              # Клавиатуры и UI компоненты
│   ├── main_menu.go        # Главное меню
│   ├── balance.go          # Меню баланса
│   ├── buy_subscription.go # Меню покупки подписки
│   ├── referrals.go        # Меню рефералов
│   ├── admin.go            # Админ меню
│   └── common.go           # Общие элементы UI
├── utils/                  # Утилиты
│   ├── user.go             # Работа с пользователями
│   ├── message.go          # Отправка сообщений
│   ├── validation.go       # Валидация данных
│   └── formatting.go       # Форматирование текста
├── bot.go                  # Старый файл (можно удалить)
└── bot_new.go              # Новый основной файл бота
```

## 🎯 Design Principles

### 1. **Separation of Concerns**
- Каждый компонент отвечает за свою область
- Обработчики команд отделены от callback'ов
- UI компоненты вынесены в отдельный пакет

### 2. **Single Responsibility**
- Каждый файл содержит одну логическую единицу
- Один обработчик = один файл
- Одна клавиатура = один файл

### 3. **Dependency Injection**
- Все зависимости передаются через конструкторы
- Легко тестировать и мокать
- Слабая связанность компонентов

### 4. **Interface Segregation**
- Каждый обработчик реализует свой интерфейс
- Можно легко заменить реализацию
- Гибкость в настройке

## 🚀 How to Add New Features

### Adding a New Command

1. Create file `handlers/commands/new_command.go`:
```go
package commands

type NewCommandHandler struct {
    config *config.Config
    // other dependencies
}

func NewNewCommandHandler(config *config.Config) *NewCommandHandler {
    return &NewCommandHandler{config: config}
}

func (h *NewCommandHandler) Handle(message *tgbotapi.Message, user *models.User, args string) error {
    // Implementation
}
```

2. Add to `bot_new.go`:
```go
// In NewBot constructor
newCommandHandler := commands.NewNewCommandHandler(cfg)

// In handleMessage switch
case "new_command":
    return b.newCommandHandler.Handle(message, user, args)
```

### Adding a New Callback

1. Create file `handlers/callbacks/new_callback.go`:
```go
package callbacks

type NewCallbackHandler struct {
    config *config.Config
    // other dependencies
}

func (h *NewCallbackHandler) Handle(query *tgbotapi.CallbackQuery, user *models.User) error {
    // Implementation
}
```

2. Add to `handleCallbackQueryData`:
```go
case data == "new_callback":
    return b.newCallbackHandler.Handle(query, user)
```

### Adding a New Keyboard

1. Create file `keyboards/new_keyboard.go`:
```go
package keyboards

type NewKeyboard struct {
    config *config.Config
}

func (k *NewKeyboard) Create(data interface{}) tgbotapi.InlineKeyboardMarkup {
    // Implementation
}
```

2. Use in handlers:
```go
keyboard := keyboards.NewNewKeyboard(config)
markup := keyboard.Create(data)
```

## 🧪 Testing

Each component can be tested independently:

```go
func TestStartHandler(t *testing.T) {
    // Setup
    config := &config.Config{...}
    userService := mockUserService{}
    handler := commands.NewStartHandler(config, userService)
    
    // Test
    err := handler.Handle(message, user, "")
    
    // Assert
    assert.NoError(t, err)
}
```

## 📈 Benefits

1. **Maintainability** - Easy to find and modify code
2. **Scalability** - Easy to add new features
3. **Testability** - Each component can be tested in isolation
4. **Readability** - Clear structure and separation
5. **Reusability** - Components can be reused across different handlers
6. **Team Development** - Multiple developers can work on different components
