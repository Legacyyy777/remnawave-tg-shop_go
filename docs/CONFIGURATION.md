# Руководство по конфигурации

## 📋 Обзор

Этот документ описывает все доступные параметры конфигурации для Remnawave Telegram Shop Bot.

## 🔧 Основные параметры

### Telegram Bot

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `BOT_TOKEN` | Токен бота от @BotFather | ✅ | - |
| `BOT_WEBHOOK_URL` | URL для webhook (для продакшена) | ❌ | - |
| `BOT_WEBHOOK_PORT` | Порт для webhook | ❌ | 8080 |

### База данных

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `DB_HOST` | Хост базы данных | ❌ | localhost |
| `DB_PORT` | Порт базы данных | ❌ | 5432 |
| `DB_USER` | Пользователь базы данных | ❌ | remnawave_bot |
| `DB_PASSWORD` | Пароль базы данных | ✅ | - |
| `DB_NAME` | Имя базы данных | ❌ | remnawave_bot |
| `DB_SSL_MODE` | Режим SSL | ❌ | disable |

### Remnawave API

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `REMNAWAVE_API_URL` | URL API панели | ✅ | - |
| `REMNAWAVE_API_KEY` | Ключ API | ✅ | - |
| `REMNAWAVE_SECRET_KEY` | Секретный ключ (опционально) | ❌ | - |

### Платежные системы

#### Tribute

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `TRIBUTE_WEBHOOK_URL` | URL для webhook Tribute | ❌ | - |

#### ЮKassa

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `YOOKASSA_SHOP_ID` | ID магазина | ❌ | - |
| `YOOKASSA_SECRET_KEY` | Секретный ключ | ❌ | - |
| `YOOKASSA_WEBHOOK_URL` | URL для webhook | ❌ | - |

### Сервер

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `SERVER_PORT` | Порт HTTP сервера | ❌ | 8080 |
| `LOG_LEVEL` | Уровень логирования | ❌ | info |
| `ENVIRONMENT` | Окружение (development/production) | ❌ | development |

### Администратор

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `ADMIN_TELEGRAM_ID` | Telegram ID администратора | ❌ | 0 |
| `MAINTENANCE_MODE` | Режим обслуживания | ❌ | false |
| `MAINTENANCE_AUTO_ENABLE` | Автоматическое включение обслуживания | ❌ | true |

### Безопасность

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `JWT_SECRET` | Секрет для JWT токенов | ❌ | - |
| `ENCRYPTION_KEY` | Ключ шифрования (32 символа) | ✅ | - |

### Мониторинг

| Параметр | Описание | Обязательный | По умолчанию |
|----------|----------|--------------|--------------|
| `HEALTH_CHECK_INTERVAL` | Интервал проверки здоровья | ❌ | 30s |
| `STATS_CLEANUP_INTERVAL` | Интервал очистки статистики | ❌ | 24h |

## 🚀 Примеры конфигурации

### Development

```env
# Telegram Bot
BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
BOT_WEBHOOK_URL=
BOT_WEBHOOK_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=remnawave_bot
DB_PASSWORD=dev_password
DB_NAME=remnawave_bot
DB_SSL_MODE=disable

# Remnawave API
REMNAWAVE_API_URL=https://dev-panel.example.com/api
REMNAWAVE_API_KEY=dev_api_key
REMNAWAVE_SECRET_KEY=

# Payment Systems
TRIBUTE_WEBHOOK_URL=
YOOKASSA_SHOP_ID=
YOOKASSA_SECRET_KEY=
YOOKASSA_WEBHOOK_URL=

# Server
SERVER_PORT=8080
LOG_LEVEL=debug
ENVIRONMENT=development

# Admin
ADMIN_TELEGRAM_ID=123456789
MAINTENANCE_MODE=false
MAINTENANCE_AUTO_ENABLE=false

# Security
JWT_SECRET=dev_jwt_secret
ENCRYPTION_KEY=dev_encryption_key_32_chars

# Monitoring
HEALTH_CHECK_INTERVAL=30s
STATS_CLEANUP_INTERVAL=24h
```

### Production

```env
# Telegram Bot
BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
BOT_WEBHOOK_URL=https://yourdomain.com/webhook
BOT_WEBHOOK_PORT=8080

# Database
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=remnawave_bot
DB_PASSWORD=very_secure_password
DB_NAME=remnawave_bot
DB_SSL_MODE=require

# Remnawave API
REMNAWAVE_API_URL=https://your-panel.com/api
REMNAWAVE_API_KEY=production_api_key
REMNAWAVE_SECRET_KEY=secret_name:secret_value

# Payment Systems
TRIBUTE_WEBHOOK_URL=https://yourdomain.com/tribute-webhook
YOOKASSA_SHOP_ID=123456
YOOKASSA_SECRET_KEY=test_1234567890abcdef
YOOKASSA_WEBHOOK_URL=https://yourdomain.com/yookassa-webhook

# Server
SERVER_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=production

# Admin
ADMIN_TELEGRAM_ID=123456789
MAINTENANCE_MODE=false
MAINTENANCE_AUTO_ENABLE=true

# Security
JWT_SECRET=very_secure_jwt_secret
ENCRYPTION_KEY=very_secure_encryption_key_32

# Monitoring
HEALTH_CHECK_INTERVAL=30s
STATS_CLEANUP_INTERVAL=24h
```

## 🔐 Безопасность

### Генерация безопасных ключей

#### JWT Secret
```bash
# Генерируем случайный JWT secret
openssl rand -base64 32
```

#### Encryption Key
```bash
# Генерируем 32-символьный ключ шифрования
openssl rand -hex 16
```

#### Database Password
```bash
# Генерируем безопасный пароль для БД
openssl rand -base64 32
```

### Рекомендации по безопасности

1. **Никогда не коммитьте .env файлы**
2. **Используйте разные ключи для разных окружений**
3. **Регулярно обновляйте пароли и ключи**
4. **Ограничьте доступ к конфигурационным файлам**
5. **Используйте HTTPS в продакшене**

## 🌍 Переменные окружения

### Установка переменных

#### Linux/macOS
```bash
# Временная установка
export BOT_TOKEN="your_token_here"

# Постоянная установка в ~/.bashrc
echo 'export BOT_TOKEN="your_token_here"' >> ~/.bashrc
source ~/.bashrc
```

#### Windows
```cmd
# Временная установка
set BOT_TOKEN=your_token_here

# Постоянная установка через GUI
# System Properties → Environment Variables
```

#### Docker
```yaml
# В docker-compose.yml
environment:
  - BOT_TOKEN=your_token_here
  - DB_PASSWORD=your_password_here
```

### Приоритет загрузки

1. Переменные окружения
2. .env файл
3. Значения по умолчанию

## 📊 Мониторинг конфигурации

### Проверка конфигурации

```bash
# Проверяем загруженную конфигурацию
curl http://localhost:8080/health

# Проверяем переменные окружения
env | grep -E "(BOT_|DB_|REMNAWAVE_)"
```

### Валидация параметров

```go
// В internal/config/config.go
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
    if c.Security.EncryptionKey == "" || len(c.Security.EncryptionKey) != 32 {
        return fmt.Errorf("ENCRYPTION_KEY must be 32 characters long")
    }
    return nil
}
```

## 🔄 Обновление конфигурации

### Hot Reload

```bash
# Перезапускаем бота с новой конфигурацией
docker-compose restart bot

# Или перезагружаем конфигурацию
docker-compose exec bot kill -HUP 1
```

### Проверка изменений

```bash
# Проверяем статус после обновления
docker-compose ps

# Проверяем логи
docker-compose logs -f bot
```

## 🛠️ Отладка конфигурации

### Включение debug режима

```env
LOG_LEVEL=debug
ENVIRONMENT=development
```

### Проверка подключений

```bash
# Проверяем подключение к БД
docker-compose exec bot ./migrate status

# Проверяем подключение к Remnawave API
curl -H "Authorization: Bearer $REMNAWAVE_API_KEY" $REMNAWAVE_API_URL/servers
```

### Логирование конфигурации

```go
// В internal/config/config.go
func (c *Config) LogConfig(logger logger.Logger) {
    logger.Info("Configuration loaded",
        "bot_token", maskToken(c.BotToken),
        "db_host", c.Database.Host,
        "db_port", c.Database.Port,
        "remnawave_url", c.Remnawave.APIURL,
        "environment", c.Environment,
    )
}

func maskToken(token string) string {
    if len(token) < 8 {
        return "***"
    }
    return token[:4] + "***" + token[len(token)-4:]
}
```

## 📝 Шаблоны конфигурации

### Минимальная конфигурация

```env
BOT_TOKEN=your_token
DB_PASSWORD=your_password
REMNAWAVE_API_URL=https://your-panel.com/api
REMNAWAVE_API_KEY=your_api_key
ENCRYPTION_KEY=your_32_character_key
```

### Полная конфигурация

```env
# Telegram Bot
BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
BOT_WEBHOOK_URL=https://yourdomain.com/webhook
BOT_WEBHOOK_PORT=8080

# Database
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=remnawave_bot
DB_PASSWORD=very_secure_password
DB_NAME=remnawave_bot
DB_SSL_MODE=require

# Remnawave API
REMNAWAVE_API_URL=https://your-panel.com/api
REMNAWAVE_API_KEY=production_api_key
REMNAWAVE_SECRET_KEY=secret_name:secret_value

# Payment Systems
TRIBUTE_WEBHOOK_URL=https://yourdomain.com/tribute-webhook
YOOKASSA_SHOP_ID=123456
YOOKASSA_SECRET_KEY=test_1234567890abcdef
YOOKASSA_WEBHOOK_URL=https://yourdomain.com/yookassa-webhook

# Server
SERVER_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=production

# Admin
ADMIN_TELEGRAM_ID=123456789
MAINTENANCE_MODE=false
MAINTENANCE_AUTO_ENABLE=true

# Security
JWT_SECRET=very_secure_jwt_secret
ENCRYPTION_KEY=very_secure_encryption_key_32

# Monitoring
HEALTH_CHECK_INTERVAL=30s
STATS_CLEANUP_INTERVAL=24h
```

## ❓ Частые проблемы

### Неверный токен бота
```
Error: failed to create bot API: 401 Unauthorized
```
**Решение**: Проверьте правильность токена в `BOT_TOKEN`

### Ошибка подключения к БД
```
Error: failed to connect to database: connection refused
```
**Решение**: Проверьте параметры подключения к БД

### Неверный API ключ Remnawave
```
Error: API error: Invalid API key
```
**Решение**: Проверьте правильность `REMNAWAVE_API_KEY`

### Неверный ключ шифрования
```
Error: ENCRYPTION_KEY must be 32 characters long
```
**Решение**: Сгенерируйте 32-символьный ключ

---

**Удачной конфигурации! 🚀**
