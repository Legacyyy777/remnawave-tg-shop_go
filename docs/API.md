# API Documentation

## 📋 Обзор

Remnawave Telegram Shop Bot предоставляет REST API для интеграции с внешними системами.

## 🔗 Базовый URL

```
http://localhost:8080/api/v1
```

## 🔐 Аутентификация

API использует JWT токены для аутентификации.

### Получение токена

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### Использование токена

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/users
```

## 📊 Endpoints

### Health Check

#### GET /health

Проверка состояния сервиса.

**Response:**
```json
{
  "status": "ok",
  "time": "2024-01-01T12:00:00Z",
  "version": "1.0.0"
}
```

### Пользователи

#### GET /api/v1/users

Получить список пользователей.

**Query Parameters:**
- `limit` (int): Количество записей (по умолчанию: 20)
- `offset` (int): Смещение (по умолчанию: 0)
- `search` (string): Поиск по имени или username

**Response:**
```json
{
  "users": [
    {
      "id": "uuid",
      "telegram_id": 123456789,
      "username": "user123",
      "first_name": "John",
      "last_name": "Doe",
      "language_code": "ru",
      "is_blocked": false,
      "is_admin": false,
      "balance": 100.50,
      "referral_code": "ABC123",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 100,
  "limit": 20,
  "offset": 0
}
```

#### GET /api/v1/users/{id}

Получить пользователя по ID.

**Response:**
```json
{
  "id": "uuid",
  "telegram_id": 123456789,
  "username": "user123",
  "first_name": "John",
  "last_name": "Doe",
  "language_code": "ru",
  "is_blocked": false,
  "is_admin": false,
  "balance": 100.50,
  "referral_code": "ABC123",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### PUT /api/v1/users/{id}

Обновить пользователя.

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "is_blocked": false,
  "balance": 150.00
}
```

**Response:**
```json
{
  "id": "uuid",
  "telegram_id": 123456789,
  "username": "user123",
  "first_name": "John",
  "last_name": "Doe",
  "language_code": "ru",
  "is_blocked": false,
  "is_admin": false,
  "balance": 150.00,
  "referral_code": "ABC123",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Подписки

#### GET /api/v1/subscriptions

Получить список подписок.

**Query Parameters:**
- `user_id` (uuid): Фильтр по пользователю
- `status` (string): Фильтр по статусу
- `limit` (int): Количество записей
- `offset` (int): Смещение

**Response:**
```json
{
  "subscriptions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "server_id": 1,
      "server_name": "Test Server",
      "plan_id": 1,
      "plan_name": "Basic Plan",
      "status": "active",
      "expires_at": "2024-02-01T12:00:00Z",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "total": 50,
  "limit": 20,
  "offset": 0
}
```

#### POST /api/v1/subscriptions

Создать новую подписку.

**Request Body:**
```json
{
  "user_id": "uuid",
  "server_id": 1,
  "plan_id": 1
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "server_id": 1,
  "server_name": "Test Server",
  "plan_id": 1,
  "plan_name": "Basic Plan",
  "status": "active",
  "expires_at": "2024-02-01T12:00:00Z",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### PUT /api/v1/subscriptions/{id}

Обновить подписку.

**Request Body:**
```json
{
  "status": "cancelled"
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "server_id": 1,
  "server_name": "Test Server",
  "plan_id": 1,
  "plan_name": "Basic Plan",
  "status": "cancelled",
  "expires_at": "2024-02-01T12:00:00Z",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Платежи

#### GET /api/v1/payments

Получить список платежей.

**Query Parameters:**
- `user_id` (uuid): Фильтр по пользователю
- `status` (string): Фильтр по статусу
- `payment_method` (string): Фильтр по способу оплаты
- `limit` (int): Количество записей
- `offset` (int): Смещение

**Response:**
```json
{
  "payments": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "amount": 299.00,
      "currency": "RUB",
      "payment_method": "stars",
      "status": "completed",
      "external_id": "payment_123",
      "description": "Пополнение баланса",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "completed_at": "2024-01-01T12:01:00Z"
    }
  ],
  "total": 25,
  "limit": 20,
  "offset": 0
}
```

#### POST /api/v1/payments

Создать новый платеж.

**Request Body:**
```json
{
  "user_id": "uuid",
  "amount": 299.00,
  "payment_method": "stars",
  "description": "Пополнение баланса"
}
```

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "amount": 299.00,
  "currency": "RUB",
  "payment_method": "stars",
  "status": "pending",
  "external_id": "payment_123",
  "description": "Пополнение баланса",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### Серверы

#### GET /api/v1/servers

Получить список серверов.

**Response:**
```json
{
  "servers": [
    {
      "id": 1,
      "name": "Test Server",
      "description": "Test server for development",
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### GET /api/v1/servers/{id}/plans

Получить планы сервера.

**Response:**
```json
{
  "plans": [
    {
      "id": 1,
      "server_id": 1,
      "name": "Basic Plan",
      "description": "Basic VPN plan for 30 days",
      "price": 299.00,
      "duration": 30,
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

### Статистика

#### GET /api/v1/stats

Получить общую статистику.

**Response:**
```json
{
  "users": {
    "total": 1000,
    "active": 800,
    "blocked": 10,
    "new_today": 5
  },
  "subscriptions": {
    "total": 500,
    "active": 400,
    "expired": 80,
    "cancelled": 20
  },
  "payments": {
    "total": 25000.00,
    "today": 500.00,
    "this_month": 5000.00
  },
  "servers": {
    "total": 5,
    "active": 4
  }
}
```

#### GET /api/v1/stats/revenue

Получить статистику по доходам.

**Query Parameters:**
- `period` (string): Период (day, week, month, year)
- `start_date` (string): Начальная дата (ISO 8601)
- `end_date` (string): Конечная дата (ISO 8601)

**Response:**
```json
{
  "period": "month",
  "start_date": "2024-01-01T00:00:00Z",
  "end_date": "2024-01-31T23:59:59Z",
  "revenue": 5000.00,
  "payments_count": 20,
  "average_payment": 250.00,
  "by_method": {
    "stars": 2000.00,
    "tribute": 1500.00,
    "yookassa": 1500.00
  }
}
```

## 🔄 Webhooks

### Telegram Webhook

#### POST /webhook

Обработка обновлений от Telegram.

**Request Body:**
```json
{
  "update_id": 123456789,
  "message": {
    "message_id": 1,
    "from": {
      "id": 123456789,
      "is_bot": false,
      "first_name": "John",
      "last_name": "Doe",
      "username": "johndoe",
      "language_code": "ru"
    },
    "chat": {
      "id": 123456789,
      "first_name": "John",
      "last_name": "Doe",
      "username": "johndoe",
      "type": "private"
    },
    "date": 1640995200,
    "text": "/start"
  }
}
```

### Tribute Webhook

#### POST /tribute-webhook

Обработка платежей от Tribute.

**Request Body:**
```json
{
  "id": "payment_123",
  "status": "completed",
  "amount": 299.00,
  "currency": "RUB",
  "user_id": "uuid"
}
```

### ЮKassa Webhook

#### POST /yookassa-webhook

Обработка платежей от ЮKassa.

**Request Body:**
```json
{
  "type": "payment.succeeded",
  "event": {
    "id": "payment_123",
    "status": "succeeded",
    "amount": {
      "value": "299.00",
      "currency": "RUB"
    },
    "metadata": {
      "user_id": "uuid"
    }
  }
}
```

## 📝 Коды ошибок

### HTTP Status Codes

| Код | Описание |
|-----|----------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 422 | Unprocessable Entity |
| 500 | Internal Server Error |

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "field": "email",
      "reason": "Invalid email format"
    }
  }
}
```

### Error Codes

| Код | Описание |
|-----|----------|
| `VALIDATION_ERROR` | Ошибка валидации |
| `NOT_FOUND` | Ресурс не найден |
| `UNAUTHORIZED` | Не авторизован |
| `FORBIDDEN` | Доступ запрещен |
| `INTERNAL_ERROR` | Внутренняя ошибка сервера |

## 🔧 Примеры использования

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "telegram_id": 123456789,
    "username": "user123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Создание подписки

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": "uuid",
    "server_id": 1,
    "plan_id": 1
  }'
```

### Получение статистики

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/stats
```

## 🧪 Тестирование API

### Postman Collection

```json
{
  "info": {
    "name": "Remnawave Bot API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/health",
          "host": ["{{base_url}}"],
          "path": ["health"]
        }
      }
    }
  ]
}
```

### cURL Examples

```bash
# Health check
curl http://localhost:8080/health

# Get users
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/users

# Create subscription
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"user_id": "uuid", "server_id": 1, "plan_id": 1}'
```

## 📚 Дополнительные ресурсы

- [OpenAPI Specification](openapi.yaml)
- [Postman Collection](postman.json)
- [API Examples](examples/)
- [Error Handling](errors.md)

---

**Удачной интеграции! 🚀**
