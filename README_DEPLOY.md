# 🚀 Развертывание Remnawave Telegram Bot

## Быстрый старт

### 1. Подготовка VPS

```bash
# Обновляем систему
sudo apt update && sudo apt upgrade -y

# Устанавливаем Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Устанавливаем Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Перезагружаемся для применения изменений
sudo reboot
```

### 2. Клонирование проекта

```bash
# Клонируем репозиторий
git clone https://github.com/your-username/remnawave-tg-shop.git
cd remnawave-tg-shop
```

### 3. Настройка

```bash
# Копируем конфигурацию
cp env.example .env

# Редактируем настройки
nano .env
```

**Обязательные параметры в .env:**
```env
# Telegram Bot
BOT_TOKEN=your_telegram_bot_token_here

# Remnawave API (если Remnawave в Docker на порту 3000)
REMNAWAVE_API_URL=http://remnawave:3000/api
REMNAWAVE_API_KEY=your_api_key_here

# Security
ENCRYPTION_KEY=your_32_character_encryption_key

# Database
DB_PASSWORD=your_secure_database_password
```

### 4. Запуск

```bash
# Запускаем развертывание
chmod +x deploy.sh
./deploy.sh
```

## Конфигурация

### Если Remnawave в отдельном Docker контейнере

```yaml
# docker-compose.prod.yml
services:
  remnawave:
    image: your-remnawave-image
    container_name: remnawave_panel
    ports:
      - "3000:3000"
    networks:
      - remnawave_network

  bot:
    environment:
      REMNAWAVE_API_URL: http://remnawave:3000/api
```

### Если Remnawave на том же сервере

```yaml
# docker-compose.prod.yml
services:
  bot:
    environment:
      REMNAWAVE_API_URL: http://host.docker.internal:3000/api
    extra_hosts:
      - "host.docker.internal:host-gateway"
```

## Мониторинг

```bash
# Просмотр логов
docker-compose -f docker-compose.prod.yml logs -f bot

# Статус сервисов
docker-compose -f docker-compose.prod.yml ps

# Использование ресурсов
docker stats

# Health check
curl http://localhost:8080/health
```

## Обновление

```bash
# Останавливаем сервисы
docker-compose -f docker-compose.prod.yml down

# Обновляем код
git pull origin main

# Пересобираем и запускаем
docker-compose -f docker-compose.prod.yml build --no-cache
docker-compose -f docker-compose.prod.yml up -d
```

## Устранение неполадок

### Бот не отвечает

```bash
# Проверяем логи
docker-compose -f docker-compose.prod.yml logs bot

# Проверяем конфигурацию
docker-compose -f docker-compose.prod.yml config

# Проверяем токен
curl -X GET "https://api.telegram.org/bot$BOT_TOKEN/getWebhookInfo"
```

### Проблемы с базой данных

```bash
# Проверяем статус БД
docker-compose -f docker-compose.prod.yml exec postgres pg_isready -U remnawave_bot

# Подключаемся к БД
docker-compose -f docker-compose.prod.yml exec postgres psql -U remnawave_bot -d remnawave_bot
```

### Проблемы с Remnawave API

```bash
# Проверяем доступность API
curl http://remnawave:80/api/servers

# Проверяем логи Remnawave
docker-compose -f docker-compose.prod.yml logs remnawave
```

## Безопасность

### Firewall

```bash
# Открываем только необходимые порты
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### SSL (опционально)

```bash
# Устанавливаем Certbot
sudo apt install certbot python3-certbot-nginx -y

# Получаем SSL сертификат
sudo certbot --nginx -d yourdomain.com
```

## Резервное копирование

```bash
# Создаем бэкап БД
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U remnawave_bot remnawave_bot > backup.sql

# Восстанавливаем из бэкапа
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U remnawave_bot -d remnawave_bot < backup.sql
```

---

**Готово! Ваш бот развернут и готов к работе! 🎉**
