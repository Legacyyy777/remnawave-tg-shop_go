# Руководство по контрибуции

## 🤝 Спасибо за интерес к проекту!

Мы приветствуем вклад в развитие Remnawave Telegram Shop Bot. Это руководство поможет вам начать работу.

## 🚀 Быстрый старт

### 1. Fork репозитория

Нажмите кнопку "Fork" в правом верхнем углу страницы репозитория.

### 2. Клонируйте ваш fork

```bash
git clone https://github.com/your-username/remnawave-tg-shop.git
cd remnawave-tg-shop
```

### 3. Добавьте upstream remote

```bash
git remote add upstream https://github.com/original-username/remnawave-tg-shop.git
```

### 4. Создайте feature branch

```bash
git checkout -b feature/your-feature-name
```

## 🔧 Настройка окружения

### Требования

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Git

### Установка

```bash
# Устанавливаем зависимости
go mod download

# Копируем конфигурацию
cp configs/development.env .env

# Запускаем базу данных
docker-compose up -d postgres

# Запускаем приложение
go run cmd/main.go
```

### Установка инструментов разработки

```bash
# Устанавливаем линтеры
make install-tools

# Запускаем тесты
make test

# Проверяем код
make lint
```

## 📋 Процесс контрибуции

### 1. Планирование

Перед началом работы:

- [ ] Проверьте существующие issues
- [ ] Создайте issue для обсуждения (если нужно)
- [ ] Убедитесь, что ваша идея соответствует целям проекта

### 2. Разработка

#### Создание feature branch

```bash
# Обновляем main branch
git checkout main
git pull upstream main

# Создаем новую ветку
git checkout -b feature/your-feature-name
```

#### Следование стандартам кода

- Используйте `gofmt` для форматирования
- Следуйте Go conventions
- Пишите тесты для нового кода
- Используйте осмысленные имена переменных и функций

#### Коммиты

```bash
# Добавляем изменения
git add .

# Коммитим с осмысленным сообщением
git commit -m "feat: add user authentication"

# Пушим в ваш fork
git push origin feature/your-feature-name
```

#### Conventional Commits

Используйте следующий формат для commit сообщений:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Типы:
- `feat`: новая функция
- `fix`: исправление бага
- `docs`: изменения в документации
- `style`: форматирование, отсутствующие точки с запятой и т.д.
- `refactor`: рефакторинг кода
- `test`: добавление тестов
- `chore`: изменения в build процессе, зависимостях и т.д.

Примеры:
```
feat(auth): add JWT token validation
fix(api): handle null values in user response
docs: update API documentation
test(user): add unit tests for user service
```

### 3. Тестирование

#### Unit тесты

```bash
# Запускаем unit тесты
go test ./internal/services/...

# Запускаем тесты с покрытием
go test -cover ./internal/services/...
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

### 4. Code Review

#### Создание Pull Request

1. Перейдите на страницу вашего fork
2. Нажмите "New Pull Request"
3. Заполните описание:
   - Что изменилось
   - Почему это нужно
   - Как тестировалось
   - Скриншоты (если применимо)

#### Шаблон Pull Request

```markdown
## Описание
Краткое описание изменений

## Тип изменений
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Как тестировалось
Опишите, как вы тестировали изменения

## Чеклист
- [ ] Код соответствует стандартам проекта
- [ ] Добавлены тесты для нового функционала
- [ ] Обновлена документация
- [ ] Все тесты проходят
- [ ] Код проверен линтерами

## Скриншоты (если применимо)
Добавьте скриншоты для UI изменений

## Дополнительная информация
Любая дополнительная информация для ревьюеров
```

### 5. Обновление

#### Синхронизация с upstream

```bash
# Переключаемся на main
git checkout main

# Получаем изменения из upstream
git pull upstream main

# Пушим в ваш fork
git push origin main

# Переключаемся на feature branch
git checkout feature/your-feature-name

# Мержим изменения из main
git merge main
```

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
make test

# Конкретный пакет
go test ./internal/services/...

# С покрытием
make test-coverage

# Бенчмарки
make bench
```

### Написание тестов

#### Unit тесты

```go
func TestFunction(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"

    // Act
    result := function(input)

    // Assert
    assert.Equal(t, expected, result)
}
```

#### Integration тесты

```go
func TestIntegration(t *testing.T) {
    // Настройка
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    // Тестирование
    // ...
}
```

### Покрытие кода

```bash
# Генерируем отчет о покрытии
go test -coverprofile=coverage.out ./...

# Просматриваем в браузере
go tool cover -html=coverage.out

# Текстовый отчет
go tool cover -func=coverage.out
```

## 📝 Документация

### Обновление документации

При внесении изменений обновляйте соответствующую документацию:

- **README.md** - общая информация о проекте
- **docs/API.md** - API документация
- **docs/USAGE.md** - руководство пользователя
- **docs/DEPLOYMENT.md** - руководство по развертыванию
- **docs/CONFIGURATION.md** - руководство по конфигурации
- **docs/DEVELOPMENT.md** - руководство по разработке
- **docs/TESTING.md** - руководство по тестированию

### Комментарии в коде

```go
// Package services provides business logic for the application.
package services

// UserService defines the interface for user-related operations.
type UserService interface {
    // CreateOrGetUser creates a new user or returns existing one.
    // It takes telegram ID, username, first name, last name and language code.
    CreateOrGetUser(telegramID int64, username, firstName, lastName, languageCode string) (*models.User, error)
}
```

## 🔍 Code Review

### Как проводить code review

1. **Проверьте функциональность**
   - Код делает то, что должен?
   - Есть ли edge cases?

2. **Проверьте качество кода**
   - Следует ли код стандартам?
   - Есть ли дублирование?
   - Легко ли читается?

3. **Проверьте тесты**
   - Есть ли тесты для нового кода?
   - Покрывают ли тесты edge cases?

4. **Проверьте безопасность**
   - Есть ли уязвимости?
   - Правильно ли обрабатываются ошибки?

### Как получить code review

1. **Будьте вежливы**
   - Используйте "пожалуйста" и "спасибо"
   - Не принимайте критику лично

2. **Будьте открыты к обратной связи**
   - Слушайте предложения
   - Задавайте вопросы

3. **Отвечайте на комментарии**
   - Объясняйте свои решения
   - Вносите изменения по запросу

## 🐛 Сообщение о багах

### Создание issue

1. Проверьте существующие issues
2. Используйте шаблон для багов
3. Предоставьте подробную информацию

#### Шаблон для багов

```markdown
## Описание бага
Краткое описание проблемы

## Шаги для воспроизведения
1. Перейдите к '...'
2. Нажмите на '...'
3. Прокрутите вниз до '...'
4. Увидите ошибку

## Ожидаемое поведение
Что должно происходить

## Фактическое поведение
Что происходит на самом деле

## Скриншоты
Если применимо, добавьте скриншоты

## Окружение
- OS: [e.g. Windows 10, macOS 12.0, Ubuntu 20.04]
- Go version: [e.g. 1.21.0]
- Docker version: [e.g. 20.10.0]
- Browser: [e.g. Chrome 91, Firefox 89]

## Дополнительная информация
Любая дополнительная информация о проблеме
```

## 💡 Предложение новых функций

### Создание feature request

1. Проверьте существующие feature requests
2. Используйте шаблон для предложений
3. Обоснуйте необходимость функции

#### Шаблон для feature requests

```markdown
## Описание функции
Краткое описание предлагаемой функции

## Проблема, которую решает
Какая проблема решается этой функцией?

## Предлагаемое решение
Как вы видите реализацию?

## Альтернативы
Какие альтернативы вы рассматривали?

## Дополнительная информация
Любая дополнительная информация о предложении
```

## 🏷️ Labels

Мы используем следующие labels:

- **bug** - что-то не работает
- **enhancement** - новая функция или улучшение
- **documentation** - улучшения документации
- **good first issue** - хорошая задача для новичков
- **help wanted** - нужна помощь
- **question** - вопрос
- **wontfix** - не будет исправлено

## 📋 Pull Request процесс

### 1. Создание PR

1. Убедитесь, что ваша ветка актуальна
2. Запустите тесты
3. Создайте PR с описанием

### 2. Обсуждение

1. Отвечайте на комментарии
2. Вносите изменения по запросу
3. Обновляйте описание при необходимости

### 3. Мерж

1. Дождитесь approval от maintainers
2. Убедитесь, что все проверки пройдены
3. Maintainer замержит PR

## 🎯 Типы контрибуций

### Код

- Исправление багов
- Добавление новых функций
- Рефакторинг
- Оптимизация производительности

### Документация

- Обновление README
- API документация
- Руководства пользователя
- Примеры использования

### Тестирование

- Unit тесты
- Integration тесты
- E2E тесты
- Performance тесты

### Дизайн

- UI/UX улучшения
- Логотипы
- Иконки
- Диаграммы

## 🏆 Признание

### Contributors

Все контрибьюторы будут добавлены в:

- **CONTRIBUTORS.md** - список контрибьюторов
- **GitHub Contributors** - автоматически
- **Release Notes** - в описании релизов

### Особые достижения

- **Major Contributor** - значительный вклад в проект
- **Bug Hunter** - нашел и исправил много багов
- **Documentation Hero** - улучшил документацию
- **Test Champion** - добавил много тестов

## 📞 Получение помощи

### Где получить помощь

- **GitHub Issues** - для вопросов и обсуждений
- **GitHub Discussions** - для общих вопросов
- **Telegram** - @your_support_bot
- **Email** - support@yourdomain.com

### Как задать вопрос

1. Проверьте существующие issues и discussions
2. Используйте поиск
3. Создайте новый issue с подробным описанием
4. Будьте вежливы и терпеливы

## 📜 Лицензия

Проект распространяется под лицензией MIT. Внося вклад, вы соглашаетесь с тем, что ваш код будет лицензирован под той же лицензией.

## 🙏 Спасибо!

Спасибо за ваш интерес к проекту! Ваш вклад помогает сделать Remnawave Telegram Shop Bot лучше для всех.

---

**Удачного контрибьютинга! 🚀**
