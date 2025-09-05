package callbacks

import (
	"fmt"
	"remnawave-tg-shop/internal/bot/keyboards"
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BalanceHandler обрабатывает callback для баланса
type BalanceHandler struct {
	config      *config.Config
	userService services.UserService
}

// NewBalanceHandler создает новый BalanceHandler
func NewBalanceHandler(config *config.Config, userService services.UserService) *BalanceHandler {
	return &BalanceHandler{
		config:      config,
		userService: userService,
	}
}

// Handle обрабатывает callback для баланса
func (h *BalanceHandler) Handle(query *tgbotapi.CallbackQuery, user *models.User) error {
	text := fmt.Sprintf("💰 Ваш баланс: %.0f₽\n\n", user.Balance)
	text += "Выберите способ пополнения:"

	// Создаем клавиатуру с методами оплаты
	keyboard := h.createPaymentKeyboard()

	// Отправляем сообщение
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, h.config.BotToken)
}

// createPaymentKeyboard создает клавиатуру с методами оплаты
func (h *BalanceHandler) createPaymentKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Кнопка "Назад"
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "start"),
	})

	// Методы оплаты (если включены)
	if h.config.Payments.StarsEnabled {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("⭐ Telegram Stars", "payment_stars"),
		})
	}

	if h.config.Payments.TributeEnabled {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("💎 Tribute", "payment_tribute"),
		})
	}

	if h.config.Payments.YooKassaEnabled {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("💳 ЮKassa", "payment_yookassa"),
		})
	}

	if h.config.Payments.CryptoPayEnabled {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("₿ CryptoPay", "payment_cryptopay"),
		})
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}
