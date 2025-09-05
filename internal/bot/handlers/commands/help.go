package commands

import (
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpHandler обрабатывает команду /help
type HelpHandler struct {
	config *config.Config
}

// NewHelpHandler создает новый HelpHandler
func NewHelpHandler(config *config.Config) *HelpHandler {
	return &HelpHandler{
		config: config,
	}
}

// Handle обрабатывает команду /help
func (h *HelpHandler) Handle(message *tgbotapi.Message, user *models.User, args string) error {
	text := "🤖 Доступные команды:\n\n"
	text += "/start - Главное меню\n"
	text += "/help - Список команд\n"
	text += "/balance - Баланс\n"
	text += "/subscriptions - Мои подписки\n"
	text += "/referrals - Рефералы\n"
	text += "/admin - Админ панель\n\n"
	text += "Используйте кнопки в меню для навигации."

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}
