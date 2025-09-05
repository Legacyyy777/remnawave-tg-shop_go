package bot

import (
	"fmt"
	"strings"
	"time"

	"remnawave-tg-shop/internal/bot/handlers/callbacks"
	"remnawave-tg-shop/internal/bot/handlers/commands"
	"remnawave-tg-shop/internal/bot/handlers/messages"
	"remnawave-tg-shop/internal/bot/middleware"
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/telebot.v3"
)

// Bot представляет телеграм-бота
type Bot struct {
	api                 *telebot.Bot
	config              *config.Config
	logger              logger.Logger
	userService         services.UserService
	subscriptionService services.SubscriptionService
	paymentService      services.PaymentService

	// Обработчики команд
	startHandler *commands.StartHandler
	helpHandler  *commands.HelpHandler
	adminHandler *commands.AdminHandler

	// Обработчики callback'ов
	balanceHandler   *callbacks.BalanceHandler
	promoCodeHandler *callbacks.PromoCodeHandler

	// Обработчики сообщений
	textHandler *messages.TextHandler

	// Middleware
	authMiddleware *middleware.AuthMiddleware
}

// NewBot создает нового бота
func NewBot(cfg *config.Config, log logger.Logger, userService services.UserService, subscriptionService services.SubscriptionService, paymentService services.PaymentService, promoCodeService services.IPromoCodeService, notificationService services.INotificationService, activityLogService services.IActivityLogService) (*Bot, error) {
	pref := telebot.Settings{
		Token: cfg.BotToken,
		// Используем Long Polling для простоты
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	api, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	// Создаем обработчики
	startHandler := commands.NewStartHandler(cfg, userService, subscriptionService)
	helpHandler := commands.NewHelpHandler(cfg)
	adminHandler := commands.NewAdminHandler(cfg, userService, subscriptionService, paymentService, promoCodeService, notificationService, activityLogService)
	balanceHandler := callbacks.NewBalanceHandler(cfg, userService)
	promoCodeHandler := callbacks.NewPromoCodeHandler(cfg, userService, promoCodeService, activityLogService)
	textHandler := messages.NewTextHandler(cfg)
	authMiddleware := middleware.NewAuthMiddleware(userService, log)

	bot := &Bot{
		api:                 api,
		config:              cfg,
		logger:              log,
		userService:         userService,
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
		startHandler:        startHandler,
		helpHandler:         helpHandler,
		adminHandler:        adminHandler,
		balanceHandler:      balanceHandler,
		promoCodeHandler:    promoCodeHandler,
		textHandler:         textHandler,
		authMiddleware:      authMiddleware,
	}

	// Регистрируем обработчики
	bot.setupHandlers()

	return bot, nil
}

// Start запускает бота
func (b *Bot) Start() error {
	b.logger.Info("Starting Telegram bot with Telebot...")
	b.api.Start()
	return nil
}

// HandleUpdate обрабатывает обновления для webhook
func (b *Bot) HandleUpdate(update interface{}) error {
	// Преобразуем обновление в формат telebot
	if tgbotUpdate, ok := update.(tgbotapi.Update); ok {
		// Обрабатываем обновление напрямую
		return b.processUpdate(tgbotUpdate)
	}
	return nil
}

// processUpdate обрабатывает обновление
func (b *Bot) processUpdate(update tgbotapi.Update) error {
	b.logger.Info("Processing update", "update_id", update.UpdateID)

	// Обрабатываем сообщения
	if update.Message != nil {
		return b.handleMessage(update.Message)
	}

	// Обрабатываем callback queries
	if update.CallbackQuery != nil {
		return b.handleCallbackQuery(update.CallbackQuery)
	}

	return nil
}

// handleMessage обрабатывает сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	b.logger.Info("Handling message", "chat_id", message.Chat.ID, "text", message.Text)

	// Получаем пользователя
	user, err := utils.GetOrCreateUser(message.From, b.userService)
	if err != nil {
		b.logger.Error("Failed to get user", "error", err)
		return err
	}

	// Обрабатываем команды
	if message.IsCommand() {
		command := message.Command()
		args := message.CommandArguments()

		switch command {
		case "start":
			return b.startHandler.Handle(message, user, args)
		case "help":
			return b.helpHandler.Handle(message, user, args)
		case "admin":
			return b.adminHandler.Handle(message, user, args)
		case "promo":
			// Обрабатываем команду промокода
			return b.promoCodeHandler.HandlePromoCodeMessage(message, user)
		default:
			return b.handleUnknownCommand(message, user, args)
		}
	}

	// Обрабатываем обычные сообщения
	return b.textHandler.Handle(message, user)
}

// handleCallbackQuery обрабатывает callback queries
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) error {
	b.logger.Info("Handling callback query", "chat_id", query.Message.Chat.ID, "data", query.Data)

	// Получаем пользователя
	user, err := utils.GetOrCreateUser(query.From, b.userService)
	if err != nil {
		b.logger.Error("Failed to get user", "error", err)
		return err
	}

	// Обрабатываем callback query
	return b.handleCallbackQueryData(query, user)
}

// handleCallbackQueryData обрабатывает данные callback query
func (b *Bot) handleCallbackQueryData(query *tgbotapi.CallbackQuery, user *models.User) error {
	data := query.Data

	// Обрабатываем различные типы callback'ов
	switch {
	case data == "balance":
		return b.balanceHandler.Handle(query, user)
	case data == "buy_subscription":
		return b.handleBuySubscription(query, user)
	case strings.HasPrefix(data, "subscription:"):
		return b.handleSubscriptionSelection(query, user)
	case data == "start":
		return b.handleStartCallback(query, user)
	case strings.HasPrefix(data, "promo_code:"):
		return b.promoCodeHandler.Handle(query, user)
	default:
		b.logger.Info("Unknown callback data", "data", data)
		return nil
	}
}

// handleBuySubscription обрабатывает callback для покупки подписки
func (b *Bot) handleBuySubscription(query *tgbotapi.CallbackQuery, _ *models.User) error {
	text := "🚀 Выберите тарифный план:\n\n"
	text += "📦 Basic (30 дней) - 299₽\n"
	text += "⭐ Premium (90 дней) - 799₽\n"
	text += "💎 Pro (365 дней) - 2499₽\n\n"
	text += "Выберите подходящий тариф:"

	// Создаем клавиатуру с тарифами
	keyboard := b.createSubscriptionKeyboard()

	// Отправляем сообщение
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
}

// createSubscriptionKeyboard создает клавиатуру с тарифами подписки
func (b *Bot) createSubscriptionKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Тарифы
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📦 Basic (30 дней) - 299₽", "subscription:basic"),
	})
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⭐ Premium (90 дней) - 799₽", "subscription:premium"),
	})
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("💎 Pro (365 дней) - 2499₽", "subscription:pro"),
	})

	// Кнопка "Назад"
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "start"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// handleSubscriptionSelection обрабатывает выбор тарифа подписки
func (b *Bot) handleSubscriptionSelection(query *tgbotapi.CallbackQuery, user *models.User) error {
	data := query.Data
	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		return b.handleBuySubscription(query, user)
	}

	plan := parts[1]

	// Определяем параметры тарифа
	var duration int
	var price float64
	var planName string

	switch plan {
	case "basic":
		duration = 30
		price = 299
		planName = "Basic"
	case "premium":
		duration = 90
		price = 799
		planName = "Premium"
	case "pro":
		duration = 365
		price = 2499
		planName = "Pro"
	default:
		return b.handleBuySubscription(query, user)
	}

	// Проверяем баланс пользователя
	if user.Balance < price {
		text := "❌ Недостаточно средств на балансе!\n\n"
		text += fmt.Sprintf("💰 Ваш баланс: %.0f₽\n", user.Balance)
		text += fmt.Sprintf("💳 Стоимость: %.0f₽\n\n", price)
		text += "Пополните баланс для покупки подписки."

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💰 Пополнить баланс", "balance"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "buy_subscription"),
			),
		)

		return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
	}

	// Создаем подписку (конвертируем дни в месяцы)
	durationMonths := duration / 30
	if durationMonths < 1 {
		durationMonths = 1
	}

	err := b.subscriptionService.CreateSubscriptionByPlan(user.ID, planName, durationMonths, int(price))
	if err != nil {
		b.logger.Error("Failed to create subscription", "error", err, "user_id", user.ID, "plan", plan)
		text := "❌ Ошибка при создании подписки. Попробуйте позже."
		return utils.SendMessage(query.Message.Chat.ID, text, b.config.BotToken)
	}

	// Списываем средства с баланса
	err = b.userService.SubtractBalance(user.ID, price)
	if err != nil {
		b.logger.Error("Failed to subtract balance", "error", err, "user_id", user.ID, "amount", price)
		text := "❌ Ошибка при списании средств. Попробуйте позже."
		return utils.SendMessage(query.Message.Chat.ID, text, b.config.BotToken)
	}

	// Отправляем подтверждение
	text := fmt.Sprintf("✅ Подписка %s успешно активирована!\n\n", planName)
	text += fmt.Sprintf("📅 Срок действия: %d дней\n", duration)
	text += fmt.Sprintf("💰 Стоимость: %.0f₽\n", price)
	text += "🔒 Используйте кнопку 'Моя подписка' для получения конфигурации VPN."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Главное меню", "start"),
		),
	)

	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
}

// handleStartCallback обрабатывает callback для главного меню
func (b *Bot) handleStartCallback(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Создаем сообщение как для команды /start
	message := &tgbotapi.Message{
		Chat: query.Message.Chat,
		From: query.From,
	}

	return b.startHandler.Handle(message, user, "")
}

// handleUnknownCommand обрабатывает неизвестные команды
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message, _ *models.User, _ string) error {
	text := "❓ Неизвестная команда. Используйте /help для получения списка команд."
	return utils.SendMessage(message.Chat.ID, text, b.config.BotToken)
}

// setupHandlers настраивает обработчики команд и callback'ов
func (b *Bot) setupHandlers() {
	// Middleware для логирования и аутентификации
	b.api.Use(b.authMiddleware.Handle)

	// Команды
	b.api.Handle("/start", b.handleStartCommand)
	b.api.Handle("/help", b.handleHelpCommand)

	// Callback queries
	b.api.Handle("\fbalance", b.handleBalanceCallback)
	b.api.Handle("\fstart", b.handleStartCallbackTelebot)

	// Text messages
	b.api.Handle(telebot.OnText, b.handleTextMessage)
}

// Обработчики для telebot (для совместимости)
func (b *Bot) handleStartCommand(c telebot.Context) error {
	user := c.Get("user").(*models.User)
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: c.Message().Chat.ID},
		From: &tgbotapi.User{ID: c.Message().Sender.ID},
	}
	return b.startHandler.Handle(message, user, "")
}

func (b *Bot) handleHelpCommand(c telebot.Context) error {
	user := c.Get("user").(*models.User)
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: c.Message().Chat.ID},
		From: &tgbotapi.User{ID: c.Message().Sender.ID},
	}
	return b.helpHandler.Handle(message, user, "")
}

func (b *Bot) handleBalanceCallback(c telebot.Context) error {
	user := c.Get("user").(*models.User)
	query := &tgbotapi.CallbackQuery{
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: c.Message().Chat.ID}},
		From:    &tgbotapi.User{ID: c.Message().Sender.ID},
		Data:    c.Callback().Data,
	}
	return b.balanceHandler.Handle(query, user)
}

func (b *Bot) handleStartCallbackTelebot(c telebot.Context) error {
	user := c.Get("user").(*models.User)
	query := &tgbotapi.CallbackQuery{
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: c.Message().Chat.ID}},
		From:    &tgbotapi.User{ID: c.Message().Sender.ID},
		Data:    c.Callback().Data,
	}
	return b.handleStartCallback(query, user)
}

func (b *Bot) handleTextMessage(c telebot.Context) error {
	user := c.Get("user").(*models.User)
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: c.Message().Chat.ID},
		From: &tgbotapi.User{ID: c.Message().Sender.ID},
		Text: c.Message().Text,
	}
	return b.textHandler.Handle(message, user)
}
