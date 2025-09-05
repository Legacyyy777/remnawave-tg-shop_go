package middleware

import (
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	"gopkg.in/telebot.v3"
)

// AuthMiddleware создает middleware для аутентификации
type AuthMiddleware struct {
	userService services.UserService
	logger      logger.Logger
}

// NewAuthMiddleware создает новый AuthMiddleware
func NewAuthMiddleware(userService services.UserService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		logger:      logger,
	}
}

// Handle обрабатывает middleware для аутентификации
func (m *AuthMiddleware) Handle(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Логируем входящее сообщение
		if c.Message() != nil {
			m.logger.Info("Received message", 
				"chat_id", c.Message().Chat.ID, 
				"user_id", c.Message().From.ID,
				"text", c.Message().Text)
		}

		// Получаем пользователя из контекста
		user := m.getUserFromContext(c)
		if user == nil {
			return c.Send("❌ Ошибка получения данных пользователя")
		}

		// Сохраняем пользователя в контексте
		c.Set("user", user)

		return next(c)
	}
}

// getUserFromContext извлекает пользователя из контекста
func (m *AuthMiddleware) getUserFromContext(c telebot.Context) *models.User {
	if c.Message() == nil {
		return nil
	}

	from := c.Message().From
	if from == nil {
		return nil
	}

	// Получаем или создаем пользователя
	user, err := m.userService.CreateOrGetUser(
		from.ID,
		from.UserName,
		from.FirstName,
		from.LastName,
		from.LanguageCode,
	)
	if err != nil {
		m.logger.Error("Failed to get or create user", "error", err)
		return nil
	}

	return user
}
