package keyboards

import (
	"fmt"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	"github.com/mymmrac/telego"
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
func (k *MainMenuKeyboard) Create(user *models.User) *telego.InlineKeyboardMarkup {
	var keyboardRows [][]telego.InlineKeyboardButton
	
	// Баланс
	balanceText := fmt.Sprintf("💰 Баланс %.0f₽", user.Balance)
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text:         balanceText,
			CallbackData: "balance",
		},
	})
	
	// Купить
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text:         "🚀 Купить",
			CallbackData: "buy_subscription",
		},
	})
	
	// Пробный период (если включен и пользователь еще не использовал)
	if k.config.Trial.Enabled {
		hasUsedTrial, err := k.subscriptionService.HasUsedTrial(user.ID)
		if err == nil && !hasUsedTrial {
			keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
				{
					Text:         "🎁 Пробный период",
					CallbackData: "trial",
				},
			})
		}
	}
	
	// Моя подписка - Mini App кнопка
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text: "🔒 Моя подписка",
			WebApp: &telego.WebAppInfo{
				URL: k.config.MiniApp.URL,
			},
		},
	})
	
	// Рефералы и Промокод
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text:         "🎁 Рефералы",
			CallbackData: "referrals",
		},
		{
			Text:         "🎟️ Промокод",
			CallbackData: "promo_code",
		},
	})
	
	// Язык и Статус
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text:         "🌐 Язык",
			CallbackData: "language",
		},
		{
			Text:         "📊 Статус",
			CallbackData: "status",
		},
	})
	
	// Поддержка
	keyboardRows = append(keyboardRows, []telego.InlineKeyboardButton{
		{
			Text:         "🆘 Поддержка",
			CallbackData: "support",
		},
	})

	return &telego.InlineKeyboardMarkup{
		InlineKeyboard: keyboardRows,
	}
}
