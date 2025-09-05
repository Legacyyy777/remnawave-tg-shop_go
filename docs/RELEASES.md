# Руководство по релизам

## 🚀 Обзор

Этот документ описывает процесс создания и управления релизами Remnawave Telegram Shop Bot.

## 📋 Типы релизов

### Major Release (X.0.0)

- **Несовместимые изменения** API
- **Удаление функций** - удаление устаревших функций
- **Изменение архитектуры** - значительные изменения в архитектуре
- **Новые зависимости** - добавление новых зависимостей

### Minor Release (X.Y.0)

- **Новые функции** - добавление нового функционала
- **Улучшения** - улучшение существующих функций
- **Новые API endpoints** - добавление новых API
- **Новые конфигурации** - добавление новых параметров

### Patch Release (X.Y.Z)

- **Исправления багов** - исправление ошибок
- **Улучшения безопасности** - исправления уязвимостей
- **Обновления зависимостей** - обновление версий зависимостей
- **Документация** - обновление документации

## 🔄 Процесс релиза

### 1. Планирование

#### Создание milestone

1. Перейдите в **Issues** → **Milestones**
2. Нажмите **New milestone**
3. Заполните информацию:
   - **Title**: `v1.1.0 - New Features`
   - **Description**: Описание релиза
   - **Due date**: Планируемая дата релиза

#### Планирование функций

1. Создайте issues для новых функций
2. Назначьте их на milestone
3. Оцените сложность и время выполнения
4. Назначьте ответственных

### 2. Разработка

#### Создание feature branch

```bash
# Обновляем main branch
git checkout main
git pull origin main

# Создаем feature branch
git checkout -b feature/new-feature

# Разрабатываем функцию
# ... код ...

# Коммитим изменения
git add .
git commit -m "feat: add new feature"

# Пушим в remote
git push origin feature/new-feature
```

#### Code Review

1. Создайте **Pull Request**
2. Назначьте ревьюеров
3. Дождитесь approval
4. Мержите в main

### 3. Тестирование

#### Unit тесты

```bash
# Запускаем unit тесты
make test

# Проверяем покрытие
make test-coverage
```

#### Integration тесты

```bash
# Запускаем integration тесты
go test ./tests/integration/...
```

#### E2E тесты

```bash
# Запускаем E2E тесты
go test ./tests/e2e/...
```

#### Performance тесты

```bash
# Запускаем бенчмарки
make bench

# Профилирование
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

### 4. Подготовка релиза

#### Обновление версии

```bash
# Обновляем версию в go.mod
go mod edit -module=remnawave-tg-shop/v1.1.0

# Обновляем версию в коде
# internal/version/version.go
const Version = "1.1.0"
```

#### Обновление документации

1. **CHANGELOG.md** - добавьте записи о изменениях
2. **README.md** - обновите информацию о версии
3. **API.md** - обновите API документацию
4. **USAGE.md** - обновите руководство пользователя

#### Создание release notes

```markdown
# Release Notes v1.1.0

## 🎉 Новые функции

- ✨ Добавлена поддержка новых платежных систем
- 🔄 Улучшена система уведомлений
- 📊 Добавлена расширенная аналитика

## 🐛 Исправления

- 🔧 Исправлена ошибка в обработке платежей
- 🛡️ Улучшена безопасность API
- 📱 Исправлена проблема с мобильным интерфейсом

## 🔄 Изменения

- ⚠️ Изменен формат API ответов (обратная совместимость сохранена)
- 📝 Обновлена документация
- 🧪 Добавлены новые тесты

## 📦 Обновления

- 🔄 Обновлены зависимости
- 🐳 Обновлены Docker образы
- 📚 Обновлена документация

## 🚀 Развертывание

### Требования

- Go 1.21+
- PostgreSQL 15+
- Docker 20.10+

### Обновление

1. Остановите текущую версию
2. Создайте резервную копию
3. Обновите код
4. Примените миграции
5. Запустите новую версию

### Откат

Если возникли проблемы:

1. Остановите новую версию
2. Восстановите из резервной копии
3. Запустите предыдущую версию
4. Сообщите о проблеме

## 📞 Поддержка

- **GitHub Issues**: [Создать issue](https://github.com/your-username/remnawave-tg-shop/issues)
- **Telegram**: @your_support_bot
- **Email**: support@yourdomain.com
```

### 5. Создание релиза

#### Создание Git tag

```bash
# Создаем tag
git tag -a v1.1.0 -m "Release v1.1.0"

# Пушим tag
git push origin v1.1.0
```

#### Создание GitHub Release

1. Перейдите в **Releases** → **Create a new release**
2. Выберите tag `v1.1.0`
3. Заполните информацию:
   - **Release title**: `v1.1.0 - New Features`
   - **Description**: Вставьте release notes
   - **Attach binaries**: Прикрепите собранные бинарники

#### Создание Docker образа

```bash
# Собираем Docker образ
docker build -t remnawave-bot:v1.1.0 .

# Тегируем для Docker Hub
docker tag remnawave-bot:v1.1.0 your-username/remnawave-bot:v1.1.0
docker tag remnawave-bot:v1.1.0 your-username/remnawave-bot:latest

# Пушим в Docker Hub
docker push your-username/remnawave-bot:v1.1.0
docker push your-username/remnawave-bot:latest
```

### 6. Развертывание

#### Staging

```bash
# Развертываем на staging
docker-compose -f docker-compose.staging.yml up -d

# Проверяем работу
curl http://staging.yourdomain.com/health
```

#### Production

```bash
# Развертываем на production
docker-compose -f docker-compose.prod.yml up -d

# Проверяем работу
curl http://yourdomain.com/health
```

### 7. Мониторинг

#### Проверка работы

- **Health checks** - проверка состояния сервисов
- **Логи** - мониторинг логов на ошибки
- **Метрики** - проверка производительности
- **Пользователи** - мониторинг активности пользователей

#### Откат при проблемах

```bash
# Останавливаем новую версию
docker-compose down

# Запускаем предыдущую версию
docker-compose -f docker-compose.previous.yml up -d
```

## 📅 Календарь релизов

### Регулярные релизы

- **Patch релизы** - каждые 2 недели
- **Minor релизы** - каждые 2 месяца
- **Major релизы** - каждые 6 месяцев

### Специальные релизы

- **Security релизы** - по мере необходимости
- **Hotfix релизы** - для критических исправлений
- **Beta релизы** - для тестирования новых функций

## 🔧 Инструменты

### Автоматизация

#### GitHub Actions

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.21
    - name: Build
      run: go build -o remnawave-bot ./cmd/main.go
    - name: Create Release
      uses: actions/create-release@v1
      with:
        tag_name: ${{ github.ref }}
        release_name: Release ${{ github.ref }}
        body: |
          Release ${{ github.ref }}
        draft: false
        prerelease: false
```

#### Makefile

```makefile
# Makefile
.PHONY: release patch minor major

release: test build
	@echo "Creating release..."

patch:
	@echo "Creating patch release..."
	@$(MAKE) update-version PATCH=1
	@$(MAKE) create-release

minor:
	@echo "Creating minor release..."
	@$(MAKE) update-version MINOR=1
	@$(MAKE) create-release

major:
	@echo "Creating major release..."
	@$(MAKE) update-version MAJOR=1
	@$(MAKE) create-release

update-version:
	@echo "Updating version..."

create-release:
	@echo "Creating release..."
	git add .
	git commit -m "chore: release $(VERSION)"
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
```

### Мониторинг

#### Health checks

```bash
# Проверка здоровья
curl http://localhost:8080/health

# Проверка метрик
curl http://localhost:8080/metrics
```

#### Логи

```bash
# Просмотр логов
docker-compose logs -f bot

# Фильтрация логов
docker-compose logs -f bot | grep ERROR
```

## 📊 Метрики релизов

### Отслеживание

- **Время релиза** - время от планирования до релиза
- **Количество багов** - количество багов после релиза
- **Время исправления** - время исправления критических багов
- **Удовлетворенность** - отзывы пользователей

### Анализ

- **Анализ багов** - анализ причин багов
- **Улучшения процесса** - улучшения процесса релиза
- **Обучение команды** - обучение на основе опыта

## 🚨 Экстренные релизы

### Hotfix процесс

1. **Обнаружение проблемы** - выявление критической проблемы
2. **Создание hotfix branch** - создание ветки для исправления
3. **Быстрое исправление** - исправление проблемы
4. **Тестирование** - минимальное тестирование
5. **Релиз** - быстрый релиз исправления

### Security релизы

1. **Обнаружение уязвимости** - выявление уязвимости безопасности
2. **Оценка риска** - оценка серьезности уязвимости
3. **Исправление** - исправление уязвимости
4. **Тестирование** - тестирование исправления
5. **Релиз** - релиз с уведомлением о безопасности

## 📚 Документация

### Release Notes

- **Что нового** - описание новых функций
- **Исправления** - список исправленных багов
- **Изменения** - список изменений
- **Обновления** - обновления зависимостей

### Migration Guide

- **Breaking changes** - описание несовместимых изменений
- **Миграция данных** - инструкции по миграции
- **Обновление конфигурации** - изменения в конфигурации
- **Обновление API** - изменения в API

### Changelog

- **Ведение changelog** - регулярное обновление
- **Формат** - стандартный формат записей
- **Категории** - группировка изменений
- **Версии** - привязка к версиям

## 🤝 Команда

### Роли

- **Release Manager** - координация релиза
- **Developers** - разработка функций
- **QA Engineers** - тестирование
- **DevOps Engineers** - развертывание
- **Technical Writers** - документация

### Ответственности

- **Планирование** - планирование релиза
- **Разработка** - разработка функций
- **Тестирование** - тестирование релиза
- **Развертывание** - развертывание в продакшен
- **Мониторинг** - мониторинг после релиза

## 📞 Поддержка

### Контакты

- **Release Manager**: release@yourdomain.com
- **Development Team**: dev@yourdomain.com
- **QA Team**: qa@yourdomain.com
- **DevOps Team**: devops@yourdomain.com

### Каналы связи

- **Slack**: #releases
- **Telegram**: @your_release_bot
- **Email**: releases@yourdomain.com
- **GitHub**: [Discussions](https://github.com/your-username/remnawave-tg-shop/discussions)

---

**Удачных релизов! 🚀**
