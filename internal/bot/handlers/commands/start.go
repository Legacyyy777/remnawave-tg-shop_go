package commands

import (
	"fmt"
	"remnawave-tg-shop/internal/bot/keyboards"
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	config              *config.Config
	userService         services.UserService
	subscriptionService services.SubscriptionService
	keyboard            *keyboards.MainMenuKeyboard
}

// NewStartHandler создает новый StartHandler
func NewStartHandler(
	config *config.Config,
	userService services.UserService,
	subscriptionService services.SubscriptionService,
) *StartHandler {
	return &StartHandler{
		config:              config,
		userService:         userService,
		subscriptionService: subscriptionService,
		keyboard:            keyboards.NewMainMenuKeyboard(config, subscriptionService),
	}
}

// Handle обрабатывает команду /start
func (h *StartHandler) Handle(message *tgbotapi.Message, user *models.User, args string) error {
	// Обработка реферального кода
	if args != "" {
		referralUser, err := h.userService.GetUserByReferralCode(args)
		if err == nil && referralUser != nil && referralUser.ID != user.ID {
			user.ReferredBy = &referralUser.ID
			h.userService.UpdateUser(user)
			h.userService.AddBalance(referralUser.ID, 50)
		}
	}

	// Формируем приветствие с именем пользователя
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	// Создаем клавиатуру
	keyboard := h.keyboard.Create(user)

	// Отправляем сообщение
	return utils.SendMessageWithTelegoKeyboard(message.Chat.ID, text, keyboard, h.config.BotToken)
}
