package callbacks

import (
	"fmt"
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

// PromoCodeHandler обрабатывает callback'и для промокодов
type PromoCodeHandler struct {
	config             *config.Config
	userService        services.UserService
	promoCodeService   services.IPromoCodeService
	activityLogService services.IActivityLogService
}

// NewPromoCodeHandler создает новый PromoCodeHandler
func NewPromoCodeHandler(
	config *config.Config,
	userService services.UserService,
	promoCodeService services.IPromoCodeService,
	activityLogService services.IActivityLogService,
) *PromoCodeHandler {
	return &PromoCodeHandler{
		config:             config,
		userService:        userService,
		promoCodeService:   promoCodeService,
		activityLogService: activityLogService,
	}
}

// Handle обрабатывает callback для промокодов
func (h *PromoCodeHandler) Handle(query *tgbotapi.CallbackQuery, user *models.User) error {
	data := query.Data

	// Парсим данные callback'а
	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		return h.showPromoCodeMenu(query, user)
	}

	action := parts[1]

	switch action {
	case "menu":
		return h.showPromoCodeMenu(query, user)
	case "apply":
		if len(parts) < 3 {
			return h.showPromoCodeInput(query, user)
		}
		code := parts[2]
		return h.applyPromoCode(query, user, code)
	case "input":
		return h.showPromoCodeInput(query, user)
	default:
		return h.showPromoCodeMenu(query, user)
	}
}

// showPromoCodeMenu показывает меню промокодов
func (h *PromoCodeHandler) showPromoCodeMenu(query *tgbotapi.CallbackQuery, user *models.User) error {
	text := "🎟️ *Промокоды*\n\n"
	text += "Введите промокод для получения бонусов!\n\n"
	text += "Доступные типы промокодов:\n"
	text += "• 🎁 Бонусные дни подписки\n"
	text += "• 💰 Скидка на покупку\n"
	text += "• 🎯 Специальные предложения\n\n"
	text += "Нажмите кнопку ниже, чтобы ввести промокод:"

	// Создаем клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести промокод", "promo_code:input"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "start"),
		),
	)

	// Отправляем сообщение
	msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = &keyboard

	bot, err := tgbotapi.NewBotAPI(h.config.BotToken)
	if err != nil {
		return err
	}

	_, err = bot.Send(msg)
	return err
}

// showPromoCodeInput показывает форму ввода промокода
func (h *PromoCodeHandler) showPromoCodeInput(query *tgbotapi.CallbackQuery, user *models.User) error {
	text := "📝 *Ввод промокода*\n\n"
	text += "Отправьте промокод в следующем сообщении.\n\n"
	text += "Пример: `PROMO2024` или `BONUS50`\n\n"
	text += "⚠️ Промокод можно использовать только один раз!"

	// Создаем клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад к промокодам", "promo_code:menu"),
		),
	)

	// Отправляем сообщение
	msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = &keyboard

	bot, err := tgbotapi.NewBotAPI(h.config.BotToken)
	if err != nil {
		return err
	}

	_, err = bot.Send(msg)
	return err
}

// applyPromoCode применяет промокод
func (h *PromoCodeHandler) applyPromoCode(query *tgbotapi.CallbackQuery, user *models.User, code string) error {
	// Логируем попытку применения промокода
	h.activityLogService.LogPromoCode(user.ID, uuid.Nil, code, "", "")

	// Применяем промокод
	promoCode, err := h.promoCodeService.ApplyPromoCode(user.ID, code)
	if err != nil {
		text := "❌ *Ошибка применения промокода*\n\n"
		text += fmt.Sprintf("Причина: %s\n\n", err.Error())
		text += "Проверьте правильность введенного кода и попробуйте снова."

		// Создаем клавиатуру
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Попробовать снова", "promo_code:input"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "start"),
			),
		)

		// Отправляем сообщение
		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = &keyboard

		bot, err := tgbotapi.NewBotAPI(h.config.BotToken)
		if err != nil {
			return err
		}

		_, err = bot.Send(msg)
		return err
	}

	// Промокод успешно применен
	text := "✅ *Промокод успешно применен!*\n\n"
	text += fmt.Sprintf("🎟️ Код: `%s`\n", promoCode.Code)
	text += fmt.Sprintf("📝 Тип: %s\n", promoCode.GetTypeText())
	text += fmt.Sprintf("💎 Значение: %.2f\n", promoCode.Value)

	if promoCode.Description != "" {
		text += fmt.Sprintf("📄 Описание: %s\n", promoCode.Description)
	}

	text += "\n🎉 Бонус добавлен к вашему аккаунту!"

	// Создаем клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎟️ Еще промокод", "promo_code:input"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "start"),
		),
	)

	// Отправляем сообщение
	msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = &keyboard

	bot, err := tgbotapi.NewBotAPI(h.config.BotToken)
	if err != nil {
		return err
	}

	_, err = bot.Send(msg)
	return err
}

// HandlePromoCodeMessage обрабатывает текстовое сообщение с промокодом
func (h *PromoCodeHandler) HandlePromoCodeMessage(message *tgbotapi.Message, user *models.User) error {
	code := strings.TrimSpace(message.Text)

	// Логируем попытку применения промокода
	h.activityLogService.LogPromoCode(user.ID, uuid.Nil, code, "", "")

	// Применяем промокод
	promoCode, err := h.promoCodeService.ApplyPromoCode(user.ID, code)
	if err != nil {
		text := "❌ *Ошибка применения промокода*\n\n"
		text += fmt.Sprintf("Причина: %s\n\n", err.Error())
		text += "Проверьте правильность введенного кода и попробуйте снова.\n\n"
		text += "Для ввода нового промокода используйте команду /promo"

		return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
	}

	// Промокод успешно применен
	text := "✅ *Промокод успешно применен!*\n\n"
	text += fmt.Sprintf("🎟️ Код: `%s`\n", promoCode.Code)
	text += fmt.Sprintf("📝 Тип: %s\n", promoCode.GetTypeText())
	text += fmt.Sprintf("💎 Значение: %.2f\n", promoCode.Value)

	if promoCode.Description != "" {
		text += fmt.Sprintf("📄 Описание: %s\n", promoCode.Description)
	}

	text += "\n🎉 Бонус добавлен к вашему аккаунту!"

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}
