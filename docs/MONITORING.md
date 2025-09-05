# Руководство по мониторингу

## 📊 Обзор

Этот документ описывает систему мониторинга Remnawave Telegram Shop Bot.

## 🎯 Цели мониторинга

### Доступность

- **Uptime** - время работы сервиса
- **Health checks** - проверка состояния компонентов
- **Response time** - время отклика API
- **Error rate** - частота ошибок

### Производительность

- **Throughput** - количество обработанных запросов
- **Latency** - задержка обработки запросов
- **Resource usage** - использование ресурсов
- **Queue length** - длина очередей

### Безопасность

- **Failed logins** - неудачные попытки входа
- **Suspicious activity** - подозрительная активность
- **Rate limiting** - срабатывание ограничений
- **Security events** - события безопасности

## 🔧 Инструменты мониторинга

### Логирование

#### Структурированные логи

```go
// internal/logger/logger.go
package logger

import (
    "github.com/sirupsen/logrus"
)

type Logger struct {
    *logrus.Logger
}

func New(level string) *Logger {
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })
    log.SetLevel(getLogLevel(level))
    return &Logger{log}
}

func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
    l.WithFields(fields).Info(msg)
}

func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
    l.WithFields(fields).Error(msg)
}
```

#### Использование в коде

```go
// В сервисах
func (s *userService) CreateUser(user *models.User) error {
    s.logger.InfoWithFields("Creating user", map[string]interface{}{
        "user_id": user.ID,
        "telegram_id": user.TelegramID,
        "username": user.Username,
    })
    
    if err := s.userRepo.Create(user); err != nil {
        s.logger.ErrorWithFields("Failed to create user", map[string]interface{}{
            "user_id": user.ID,
            "error": err.Error(),
        })
        return err
    }
    
    s.logger.InfoWithFields("User created successfully", map[string]interface{}{
        "user_id": user.ID,
    })
    
    return nil
}
```

### Метрики

#### Prometheus метрики

```go
// internal/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP метрики
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
    
    // Бизнес метрики
    usersTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "users_total",
            Help: "Total number of users",
        },
    )
    
    subscriptionsTotal = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "subscriptions_total",
            Help: "Total number of subscriptions",
        },
    )
    
    paymentsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "payments_total",
            Help: "Total number of payments",
        },
        []string{"method", "status"},
    )
    
    // Системные метрики
    databaseConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Number of active database connections",
        },
    )
    
    memoryUsage = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "memory_usage_bytes",
            Help: "Memory usage in bytes",
        },
    )
)
```

#### Использование метрик

```go
// В HTTP обработчиках
func (h *userHandler) CreateUser(c *gin.Context) {
    start := time.Now()
    
    // ... обработка запроса ...
    
    // Записываем метрики
    httpRequestsTotal.WithLabelValues(
        c.Request.Method,
        c.Request.URL.Path,
        strconv.Itoa(c.Writer.Status()),
    ).Inc()
    
    httpRequestDuration.WithLabelValues(
        c.Request.Method,
        c.Request.URL.Path,
    ).Observe(time.Since(start).Seconds())
}

// В бизнес-логике
func (s *userService) CreateUser(user *models.User) error {
    // ... создание пользователя ...
    
    // Увеличиваем счетчик пользователей
    usersTotal.Inc()
    
    return nil
}
```

### Health Checks

#### Health check endpoint

```go
// internal/handlers/health.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type HealthHandler struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewHealthHandler(db *gorm.DB, log logger.Logger) *HealthHandler {
    return &HealthHandler{db: db, logger: log}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
    status := "ok"
    checks := make(map[string]string)
    
    // Проверка базы данных
    if err := h.db.DB().Ping(); err != nil {
        status = "error"
        checks["database"] = "error"
        h.logger.Error("Database health check failed", "error", err)
    } else {
        checks["database"] = "ok"
    }
    
    // Проверка Redis (если используется)
    // if err := h.redis.Ping().Err(); err != nil {
    //     status = "error"
    //     checks["redis"] = "error"
    // } else {
    //     checks["redis"] = "ok"
    // }
    
    // Проверка Remnawave API
    // if err := h.remnawaveClient.HealthCheck(); err != nil {
    //     status = "error"
    //     checks["remnawave"] = "error"
    // } else {
    //     checks["remnawave"] = "ok"
    // }
    
    response := gin.H{
        "status": status,
        "timestamp": time.Now().UTC(),
        "version": "1.0.0",
        "checks": checks,
    }
    
    if status == "error" {
        c.JSON(http.StatusServiceUnavailable, response)
    } else {
        c.JSON(http.StatusOK, response)
    }
}
```

#### Liveness probe

```go
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
    // Простая проверка, что приложение работает
    c.JSON(http.StatusOK, gin.H{
        "status": "alive",
        "timestamp": time.Now().UTC(),
    })
}
```

#### Readiness probe

```go
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
    // Проверка готовности к обработке запросов
    if err := h.db.DB().Ping(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "not ready",
            "reason": "database unavailable",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "ready",
        "timestamp": time.Now().UTC(),
    })
}
```

### Алерты

#### Настройка алертов

```yaml
# alerts.yml
groups:
- name: remnawave-bot
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }} errors per second"
  
  - alert: HighResponseTime
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High response time detected"
      description: "95th percentile response time is {{ $value }} seconds"
  
  - alert: DatabaseDown
    expr: up{job="database"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Database is down"
      description: "Database has been down for more than 1 minute"
  
  - alert: HighMemoryUsage
    expr: memory_usage_bytes / 1024 / 1024 / 1024 > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage detected"
      description: "Memory usage is {{ $value }} GB"
```

#### Уведомления

```yaml
# notification.yml
route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
- name: 'web.hook'
  webhook_configs:
  - url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
    send_resolved: true
    title: 'Remnawave Bot Alert'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'

- name: 'telegram'
  webhook_configs:
  - url: 'https://api.telegram.org/bot<YOUR_BOT_TOKEN>/sendMessage'
    send_resolved: true
    title: 'Remnawave Bot Alert'
    text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'
```

## 📊 Дашборды

### Grafana дашборд

```json
{
  "dashboard": {
    "title": "Remnawave Bot Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
            "legendFormat": "5xx errors"
          }
        ]
      },
      {
        "title": "Active Users",
        "type": "singlestat",
        "targets": [
          {
            "expr": "users_total",
            "legendFormat": "Total Users"
          }
        ]
      }
    ]
  }
}
```

### Системные метрики

```json
{
  "title": "System Metrics",
  "panels": [
    {
      "title": "CPU Usage",
      "type": "graph",
      "targets": [
        {
          "expr": "rate(process_cpu_seconds_total[5m]) * 100",
          "legendFormat": "CPU Usage %"
        }
      ]
    },
    {
      "title": "Memory Usage",
      "type": "graph",
      "targets": [
        {
          "expr": "memory_usage_bytes / 1024 / 1024",
          "legendFormat": "Memory Usage MB"
        }
      ]
    },
    {
      "title": "Database Connections",
      "type": "graph",
      "targets": [
        {
          "expr": "database_connections_active",
          "legendFormat": "Active Connections"
        }
      ]
    }
  ]
}
```

## 🔍 Логирование

### Структура логов

```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "level": "info",
  "message": "User created",
  "service": "remnawave-bot",
  "version": "1.0.0",
  "fields": {
    "user_id": "uuid",
    "telegram_id": 123456789,
    "username": "user123",
    "action": "create_user",
    "duration_ms": 150
  }
}
```

### Уровни логирования

- **DEBUG** - отладочная информация
- **INFO** - общая информация
- **WARN** - предупреждения
- **ERROR** - ошибки
- **FATAL** - критические ошибки

### Фильтрация логов

```bash
# Фильтрация по уровню
docker-compose logs -f bot | grep "ERROR"

# Фильтрация по полю
docker-compose logs -f bot | grep "user_id"

# Фильтрация по времени
docker-compose logs -f bot --since="2024-01-01T12:00:00Z"
```

## 📈 Метрики производительности

### HTTP метрики

```go
// Middleware для метрик
func MetricsMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        // Записываем метрики
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
        ).Observe(duration.Seconds())
    })
}
```

### Бизнес метрики

```go
// Метрики для пользователей
func (s *userService) CreateUser(user *models.User) error {
    start := time.Now()
    
    if err := s.userRepo.Create(user); err != nil {
        userCreationErrors.Inc()
        return err
    }
    
    userCreationDuration.Observe(time.Since(start).Seconds())
    usersTotal.Inc()
    
    return nil
}

// Метрики для подписок
func (s *subscriptionService) CreateSubscription(subscription *models.Subscription) error {
    start := time.Now()
    
    if err := s.subscriptionRepo.Create(subscription); err != nil {
        subscriptionCreationErrors.Inc()
        return err
    }
    
    subscriptionCreationDuration.Observe(time.Since(start).Seconds())
    subscriptionsTotal.Inc()
    
    return nil
}
```

## 🚨 Алерты

### Критические алерты

- **Сервис недоступен** - HTTP 5xx ошибки
- **База данных недоступна** - ошибки подключения к БД
- **Высокое использование памяти** - > 90% памяти
- **Высокое использование CPU** - > 90% CPU

### Предупреждения

- **Высокий response time** - > 1 секунды
- **Высокая частота ошибок** - > 5% ошибок
- **Много неудачных попыток входа** - > 10 в минуту
- **Высокое использование диска** - > 80% диска

### Информационные

- **Новый пользователь** - регистрация пользователя
- **Новая подписка** - создание подписки
- **Успешный платеж** - завершение платежа
- **Обновление конфигурации** - изменение настроек

## 🔧 Настройка мониторинга

### Docker Compose

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./alerts.yml:/etc/prometheus/alerts.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
      - '--web.enable-admin-api'

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana/datasources:/etc/grafana/provisioning/datasources

  alertmanager:
    image: prom/alertmanager:latest
    container_name: alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
      - ./notification.yml:/etc/alertmanager/notification.yml

volumes:
  grafana-storage:
```

### Prometheus конфигурация

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alerts.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  - job_name: 'remnawave-bot'
    static_configs:
      - targets: ['bot:8080']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']
    scrape_interval: 30s

  - job_name: 'redis'
    static_configs:
      - targets: ['redis:6379']
    scrape_interval: 30s
```

## 📊 Отчеты

### Ежедневные отчеты

- **Uptime** - время работы сервиса
- **Количество пользователей** - новые регистрации
- **Количество подписок** - новые подписки
- **Количество платежей** - успешные платежи
- **Ошибки** - количество и типы ошибок

### Еженедельные отчеты

- **Производительность** - средний response time
- **Использование ресурсов** - CPU, память, диск
- **Безопасность** - подозрительная активность
- **Пользовательская активность** - активность пользователей

### Ежемесячные отчеты

- **Тренды** - изменения в метриках
- **Рекомендации** - предложения по улучшению
- **Планирование** - планирование ресурсов
- **Анализ** - анализ производительности

## 🛠️ Устранение неполадок

### Частые проблемы

#### Высокий response time

1. Проверьте метрики CPU и памяти
2. Проверьте количество активных соединений к БД
3. Проверьте логи на ошибки
4. Проверьте внешние зависимости

#### Высокая частота ошибок

1. Проверьте логи на детали ошибок
2. Проверьте конфигурацию
3. Проверьте внешние зависимости
4. Проверьте ресурсы системы

#### Проблемы с базой данных

1. Проверьте соединения к БД
2. Проверьте логи БД
3. Проверьте использование диска
4. Проверьте индексы

### Диагностика

```bash
# Проверка здоровья
curl http://localhost:8080/health

# Проверка метрик
curl http://localhost:8080/metrics

# Проверка логов
docker-compose logs -f bot

# Проверка ресурсов
docker stats
```

## 📚 Лучшие практики

### Логирование

- Используйте структурированные логи
- Включайте контекстную информацию
- Не логируйте чувствительные данные
- Используйте соответствующие уровни

### Метрики

- Измеряйте важные бизнес-метрики
- Используйте гистограммы для времени
- Используйте счетчики для событий
- Группируйте метрики по лейблам

### Алерты

- Настройте разумные пороги
- Избегайте ложных срабатываний
- Группируйте связанные алерты
- Настройте эскалацию

### Дашборды

- Показывайте важные метрики
- Используйте понятные названия
- Группируйте связанные метрики
- Обновляйте дашборды регулярно

---

**Удачного мониторинга! 📊**
