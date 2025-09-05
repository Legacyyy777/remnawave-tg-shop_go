# Руководство по развертыванию

## 🚀 Быстрый старт

### 1. Подготовка сервера

#### Минимальные требования
- **CPU**: 1 vCPU
- **RAM**: 512 MB
- **Диск**: 10 GB SSD
- **ОС**: Ubuntu 20.04+ / CentOS 8+ / Debian 11+

#### Рекомендуемые требования
- **CPU**: 2 vCPU
- **RAM**: 1 GB
- **Диск**: 20 GB SSD
- **ОС**: Ubuntu 22.04 LTS

### 2. Установка зависимостей

#### Ubuntu/Debian
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

# Устанавливаем Nginx (опционально)
sudo apt install nginx -y

# Устанавливаем Certbot для SSL
sudo apt install certbot python3-certbot-nginx -y
```

#### CentOS/RHEL
```bash
# Обновляем систему
sudo yum update -y

# Устанавливаем Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker

# Устанавливаем Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 3. Настройка проекта

```bash
# Клонируем репозиторий
git clone https://github.com/your-username/remnawave-tg-shop.git
cd remnawave-tg-shop

# Копируем конфигурацию
cp configs/production.env .env

# Редактируем конфигурацию
nano .env
```

### 4. Настройка переменных окружения

```env
# Обязательные параметры
BOT_TOKEN=your_telegram_bot_token
DB_PASSWORD=your_secure_database_password
REMNAWAVE_API_URL=https://your-panel.com/api
REMNAWAVE_API_KEY=your_api_key
ENCRYPTION_KEY=your_32_character_encryption_key

# Опциональные параметры
BOT_WEBHOOK_URL=https://yourdomain.com/webhook
ADMIN_TELEGRAM_ID=your_telegram_id
```

### 5. Запуск с Docker Compose

```bash
# Запускаем все сервисы
docker-compose up -d

# Проверяем статус
docker-compose ps

# Просматриваем логи
docker-compose logs -f bot
```

## 🔧 Настройка Nginx

### 1. Создание конфигурации

```bash
sudo nano /etc/nginx/sites-available/remnawave-bot
```

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";

    # Telegram webhook
    location /webhook {
        proxy_pass http://localhost:8080/webhook;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check
    location /health {
        proxy_pass http://localhost:8080/health;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Payment webhooks
    location /tribute-webhook {
        proxy_pass http://localhost:8080/tribute-webhook;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /yookassa-webhook {
        proxy_pass http://localhost:8080/yookassa-webhook;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. Активация конфигурации

```bash
# Создаем символическую ссылку
sudo ln -s /etc/nginx/sites-available/remnawave-bot /etc/nginx/sites-enabled/

# Удаляем дефолтную конфигурацию
sudo rm /etc/nginx/sites-enabled/default

# Проверяем конфигурацию
sudo nginx -t

# Перезапускаем Nginx
sudo systemctl restart nginx
```

## 🔒 Настройка SSL

### 1. Получение SSL сертификата

```bash
# Получаем сертификат от Let's Encrypt
sudo certbot --nginx -d yourdomain.com

# Проверяем автообновление
sudo certbot renew --dry-run
```

### 2. Обновление конфигурации Nginx

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;

    # Same location blocks as HTTP
}
```

## 📊 Мониторинг

### 1. Настройка логирования

```bash
# Создаем директорию для логов
sudo mkdir -p /var/log/remnawave-bot

# Настраиваем ротацию логов
sudo nano /etc/logrotate.d/remnawave-bot
```

```
/var/log/remnawave-bot/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 root root
    postrotate
        docker-compose restart bot
    endscript
}
```

### 2. Настройка мониторинга

```bash
# Устанавливаем htop для мониторинга
sudo apt install htop -y

# Устанавливаем iotop для мониторинга диска
sudo apt install iotop -y

# Устанавливаем nethogs для мониторинга сети
sudo apt install nethogs -y
```

### 3. Настройка алертов

```bash
# Создаем скрипт для проверки здоровья
sudo nano /usr/local/bin/health-check.sh
```

```bash
#!/bin/bash
# Health check script

URL="http://localhost:8080/health"
RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $URL)

if [ $RESPONSE -ne 200 ]; then
    echo "Health check failed: HTTP $RESPONSE"
    # Здесь можно добавить отправку уведомлений
    exit 1
fi

echo "Health check passed"
exit 0
```

```bash
# Делаем скрипт исполняемым
sudo chmod +x /usr/local/bin/health-check.sh

# Добавляем в crontab
sudo crontab -e
```

```
# Проверка каждые 5 минут
*/5 * * * * /usr/local/bin/health-check.sh
```

## 🔄 Обновление

### 1. Обновление кода

```bash
# Останавливаем сервисы
docker-compose down

# Обновляем код
git pull origin main

# Пересобираем образы
docker-compose build --no-cache

# Запускаем сервисы
docker-compose up -d
```

### 2. Обновление базы данных

```bash
# Применяем миграции
docker-compose exec bot ./migrate up

# Или через SQL
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -f /migrations/001_initial_schema.sql
```

## 🛠️ Резервное копирование

### 1. Настройка бэкапов

```bash
# Создаем директорию для бэкапов
sudo mkdir -p /backup/remnawave-bot

# Создаем скрипт бэкапа
sudo nano /usr/local/bin/backup.sh
```

```bash
#!/bin/bash
# Backup script

BACKUP_DIR="/backup/remnawave-bot"
DATE=$(date +%Y%m%d_%H%M%S)

# Создаем бэкап базы данных
docker-compose exec -T postgres pg_dump -U remnawave_bot remnawave_bot > $BACKUP_DIR/db_$DATE.sql

# Создаем бэкап конфигурации
cp .env $BACKUP_DIR/env_$DATE

# Удаляем старые бэкапы (старше 30 дней)
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "env_*" -mtime +30 -delete

echo "Backup completed: $DATE"
```

```bash
# Делаем скрипт исполняемым
sudo chmod +x /usr/local/bin/backup.sh

# Добавляем в crontab (ежедневно в 2:00)
sudo crontab -e
```

```
# Ежедневный бэкап в 2:00
0 2 * * * /usr/local/bin/backup.sh
```

### 2. Восстановление из бэкапа

```bash
# Останавливаем сервисы
docker-compose down

# Восстанавливаем базу данных
docker-compose exec -T postgres psql -U remnawave_bot -d remnawave_bot < /backup/remnawave-bot/db_20240101_020000.sql

# Восстанавливаем конфигурацию
cp /backup/remnawave-bot/env_20240101_020000 .env

# Запускаем сервисы
docker-compose up -d
```

## 🚨 Устранение неполадок

### 1. Проверка логов

```bash
# Логи бота
docker-compose logs bot

# Логи базы данных
docker-compose logs postgres

# Логи Nginx
sudo tail -f /var/log/nginx/error.log
```

### 2. Проверка статуса сервисов

```bash
# Статус Docker контейнеров
docker-compose ps

# Статус Nginx
sudo systemctl status nginx

# Статус базы данных
docker-compose exec postgres pg_isready -U remnawave_bot
```

### 3. Частые проблемы

#### Бот не отвечает
```bash
# Проверяем токен
echo $BOT_TOKEN

# Проверяем webhook
curl -X GET "https://api.telegram.org/bot$BOT_TOKEN/getWebhookInfo"
```

#### Ошибки базы данных
```bash
# Проверяем подключение
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "SELECT 1;"

# Проверяем таблицы
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "\dt"
```

#### Проблемы с Nginx
```bash
# Проверяем конфигурацию
sudo nginx -t

# Перезапускаем Nginx
sudo systemctl restart nginx
```

## 📈 Оптимизация производительности

### 1. Настройка PostgreSQL

```bash
# Редактируем конфигурацию PostgreSQL
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "ALTER SYSTEM SET shared_buffers = '256MB';"
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "ALTER SYSTEM SET effective_cache_size = '1GB';"
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "ALTER SYSTEM SET maintenance_work_mem = '64MB';"
docker-compose exec postgres psql -U remnawave_bot -d remnawave_bot -c "SELECT pg_reload_conf();"
```

### 2. Настройка Docker

```bash
# Ограничиваем ресурсы контейнеров
# В docker-compose.yml добавьте:
# deploy:
#   resources:
#     limits:
#       memory: 512M
#       cpus: '0.5'
```

### 3. Настройка Nginx

```nginx
# В nginx.conf добавьте:
worker_processes auto;
worker_connections 1024;

# Кэширование
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 🔐 Безопасность

### 1. Настройка файрвола

```bash
# Устанавливаем UFW
sudo apt install ufw -y

# Настраиваем правила
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443
sudo ufw enable
```

### 2. Настройка SSH

```bash
# Отключаем root login
sudo nano /etc/ssh/sshd_config
```

```
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
```

```bash
# Перезапускаем SSH
sudo systemctl restart ssh
```

### 3. Регулярные обновления

```bash
# Автоматические обновления безопасности
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure -plow unattended-upgrades
```

## 📞 Поддержка

Если у вас возникли проблемы:

1. Проверьте логи: `docker-compose logs -f`
2. Проверьте статус: `docker-compose ps`
3. Проверьте конфигурацию: `docker-compose config`
4. Создайте issue в GitHub
5. Обратитесь в Telegram: @your_support_bot

---

**Удачного развертывания! 🚀**
