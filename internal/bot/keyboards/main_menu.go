package keyboards

import (
	"fmt"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MainMenuKeyboard создает главное меню
type MainMenuKeyboard struct {
	config              *config.Config
	subscriptionService services.SubscriptionService
}

// NewMainMenuKeyboard создает новый MainMenuKeyboard
func NewMainMenuKeyboard(config *config.Config, subscriptionService services.SubscriptionService) *MainMenuKeyboard {
	return &MainMenuKeyboard{
		config:              config,
		subscriptionService: subscriptionService,
	}
}

// Create создает клавиатуру главного меню
func (k *MainMenuKeyboard) Create(user *models.User) tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	
	// Баланс
	balanceText := fmt.Sprintf("💰 Баланс %.0f₽", user.Balance)
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(balanceText, "balance"),
	})
	
	// Купить
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🚀 Купить", "buy_subscription"),
	})
	
	// Пробный период (если включен и пользователь еще не использовал)
	if k.config.Trial.Enabled {
		hasUsedTrial, err := k.subscriptionService.HasUsedTrial(user.ID)
		if err == nil && !hasUsedTrial {
			keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🎁 Пробный период", "trial"),
			})
		}
	}
	
	// Моя подписка - Mini App кнопка
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		{
			Text: "🔒 Моя подписка",
			WebApp: &tgbotapi.WebAppInfo{
				URL: k.config.MiniApp.URL,
			},
		},
	})
	
	// Рефералы и Промокод
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🎁 Рефералы", "referrals"),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Промокод", "promo_code"),
	})
	
	// Язык и Статус
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "language"),
		tgbotapi.NewInlineKeyboardButtonData("📊 Статус", "status"),
	})
	
	// Поддержка
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🆘 Поддержка", "support"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}
