#!/bin/bash

# Скрипт развертывания Remnawave Telegram Bot на VPS

set -e

echo "🚀 Развертывание Remnawave Telegram Bot..."

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен. Установите Docker и попробуйте снова."
    exit 1
fi

# Проверяем наличие Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен. Установите Docker Compose и попробуйте снова."
    exit 1
fi

# Создаем .env файл если его нет
if [ ! -f .env ]; then
    echo "📝 Создаем .env файл..."
    cp env.example .env
    echo "⚠️  Отредактируйте .env файл с вашими настройками перед запуском!"
    echo "   nano .env"
    exit 1
fi

# Проверяем обязательные переменные
source .env

if [ -z "$BOT_TOKEN" ]; then
    echo "❌ BOT_TOKEN не установлен в .env файле"
    exit 1
fi

if [ -z "$REMNAWAVE_API_URL" ]; then
    echo "❌ REMNAWAVE_API_URL не установлен в .env файле"
    exit 1
fi

if [ -z "$REMNAWAVE_API_KEY" ]; then
    echo "❌ REMNAWAVE_API_KEY не установлен в .env файле"
    exit 1
fi

if [ -z "$ENCRYPTION_KEY" ]; then
    echo "❌ ENCRYPTION_KEY не установлен в .env файле"
    exit 1
fi

# Останавливаем существующие контейнеры
echo "🛑 Останавливаем существующие контейнеры..."
docker-compose -f docker-compose.prod.yml down || true

# Собираем образ
echo "🔨 Собираем Docker образ..."
docker-compose -f docker-compose.prod.yml build --no-cache

# Запускаем сервисы
echo "🚀 Запускаем сервисы..."
docker-compose -f docker-compose.prod.yml up -d

# Ждем запуска базы данных
echo "⏳ Ждем запуска базы данных..."
sleep 10

# Проверяем статус
echo "📊 Проверяем статус сервисов..."
docker-compose -f docker-compose.prod.yml ps

# Проверяем логи
echo "📋 Последние логи бота:"
docker-compose -f docker-compose.prod.yml logs --tail=20 bot

# Проверяем health check
echo "🏥 Проверяем health check..."
sleep 5
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Бот успешно запущен и отвечает на health check!"
else
    echo "❌ Бот не отвечает на health check. Проверьте логи:"
    echo "   docker-compose -f docker-compose.prod.yml logs bot"
fi

echo ""
echo "🎉 Развертывание завершено!"
echo ""
echo "📋 Полезные команды:"
echo "   Просмотр логов:     docker-compose -f docker-compose.prod.yml logs -f bot"
echo "   Остановка:          docker-compose -f docker-compose.prod.yml down"
echo "   Перезапуск:         docker-compose -f docker-compose.prod.yml restart bot"
echo "   Статус:             docker-compose -f docker-compose.prod.yml ps"
echo ""
echo "🌐 Health check: http://localhost:8080/health"
echo "📊 Мониторинг: docker stats"
