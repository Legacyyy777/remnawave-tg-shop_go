package utils

import (
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mymmrac/telego"
)

// GetOrCreateUser получает или создает пользователя из Telegram User
func GetOrCreateUser(from *tgbotapi.User, userService services.UserService) (*models.User, error) {
	// Используем CreateOrGetUser для получения или создания пользователя
	user, err := userService.CreateOrGetUser(
		from.ID,
		from.UserName,
		from.FirstName,
		from.LastName,
		from.LanguageCode,
	)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// SendMessage отправляет сообщение через tgbotapi
func SendMessage(chatID int64, text string, botToken string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	
	msg := tgbotapi.NewMessage(chatID, text)
	_, err = bot.Send(msg)
	return err
}

// SendMessageWithKeyboard отправляет сообщение с клавиатурой
func SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup, botToken string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}
	
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err = bot.Send(msg)
	return err
}

// SendMessageWithTelegoKeyboard отправляет сообщение с клавиатурой telego
func SendMessageWithTelegoKeyboard(chatID int64, text string, keyboard *telego.InlineKeyboardMarkup, botToken string) error {
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		return err
	}
	
	msg := telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID},
		Text:      text,
		ParseMode: telego.ModeHTML,
	}
	
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	
	_, err = bot.SendMessage(&msg)
	return err
}
