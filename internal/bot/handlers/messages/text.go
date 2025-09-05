package messages

import (
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TextHandler обрабатывает текстовые сообщения
type TextHandler struct {
	config *config.Config
}

// NewTextHandler создает новый TextHandler
func NewTextHandler(config *config.Config) *TextHandler {
	return &TextHandler{
		config: config,
	}
}

// Handle обрабатывает текстовые сообщения
func (h *TextHandler) Handle(message *tgbotapi.Message, user *models.User) error {
	// Пока что просто логируем
	// В будущем здесь можно добавить обработку промокодов, поиск и т.д.
	return nil
}
