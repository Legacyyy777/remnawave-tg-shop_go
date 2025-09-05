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
	case data == "payment_tribute":
		return b.handleTributePayment(query, user)
	case data == "payment_stars":
		return b.handleStarsPayment(query, user)
	case data == "payment_yookassa":
		return b.handleYooKassaPayment(query, user)
	case data == "payment_cryptopay":
		return b.handleCryptoPayPayment(query, user)
	case data == "start":
		return b.handleStartCallback(query, user)
	case data == "support":
		return b.handleSupport(query, user)
	case data == "language":
		return b.handleLanguage(query, user)
	case data == "status":
		return b.handleStatus(query, user)
	case data == "referrals":
		return b.handleReferrals(query, user)
	case data == "trial":
		return b.handleTrial(query, user)
	case strings.HasPrefix(data, "admin:"):
		return b.handleAdminCallback(query, user)
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

// createMainMenuKeyboard создает главное меню
func (b *Bot) createMainMenuKeyboard(user *models.User) tgbotapi.InlineKeyboardMarkup {
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

	// Рефералы и Промокод
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🎁 Рефералы", "referrals"),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Промокод", "promo_code:menu"),
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

// handleTributePayment обрабатывает платеж через Tribute
func (b *Bot) handleTributePayment(query *tgbotapi.CallbackQuery, _ *models.User) error {
	text := "💎 *Пополнение через Tribute*\n\n"
	text += "Для пополнения баланса перейдите по ссылке:\n\n"
	text += "🔗 " + b.config.Payments.Tribute.AppURL + "\n\n"
	text += "После успешного платежа средства будут автоматически зачислены на ваш баланс."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("💎 Перейти к оплате", b.config.Payments.Tribute.AppURL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "balance"),
		),
	)

	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
}

// handleStarsPayment обрабатывает платеж через Telegram Stars
func (b *Bot) handleStarsPayment(query *tgbotapi.CallbackQuery, _ *models.User) error {
	text := "⭐ *Пополнение через Telegram Stars*\n\n"
	text += "Функция пополнения через Telegram Stars временно недоступна.\n"
	text += "Используйте другие способы оплаты."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "balance"),
		),
	)

	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
}

// handleYooKassaPayment обрабатывает платеж через ЮKassa
func (b *Bot) handleYooKassaPayment(query *tgbotapi.CallbackQuery, _ *models.User) error {
	text := "💳 *Пополнение через ЮKassa*\n\n"
	text += "Функция пополнения через ЮKassa временно недоступна.\n"
	text += "Используйте другие способы оплаты."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "balance"),
		),
	)

	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, text, keyboard, b.config.BotToken)
}

// handleCryptoPayPayment обрабатывает платеж через CryptoPay
func (b *Bot) handleCryptoPayPayment(query *tgbotapi.CallbackQuery, _ *models.User) error {
	text := "₿ *Пополнение через CryptoPay*\n\n"
	text += "Функция пополнения через CryptoPay временно недоступна.\n"
	text += "Используйте другие способы оплаты."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "balance"),
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

// handleSupport обрабатывает callback для поддержки
func (b *Bot) handleSupport(query *tgbotapi.CallbackQuery, user *models.User) error {
	message := "🆘 **Поддержка**\n\n" +
		"Если у вас возникли вопросы или проблемы, обратитесь к администратору:\n\n" +
		"• Напишите в личные сообщения администратору\n" +
		"• Опишите вашу проблему подробно\n" +
		"• Укажите ваш Telegram ID: `" + fmt.Sprintf("%d", user.TelegramID) + "`\n\n" +
		"Мы постараемся ответить как можно скорее! 🚀"

	keyboard := b.createMainMenuKeyboard(user)
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleLanguage обрабатывает callback для смены языка
func (b *Bot) handleLanguage(query *tgbotapi.CallbackQuery, user *models.User) error {
	message := "🌐 **Выбор языка**\n\n" +
		"В данный момент доступен только русский язык.\n" +
		"В будущих версиях будут добавлены другие языки."

	keyboard := b.createMainMenuKeyboard(user)
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleStatus обрабатывает callback для статуса
func (b *Bot) handleStatus(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Получаем активные подписки пользователя
	subscriptions, err := b.subscriptionService.GetActiveSubscriptions(user.ID)
	if err != nil {
		b.logger.Error("Failed to get user subscriptions", "error", err)
		subscriptions = []models.Subscription{}
	}

	message := "📊 **Ваш статус**\n\n"
	message += fmt.Sprintf("💰 Баланс: %.0f₽\n", user.Balance)
	message += fmt.Sprintf("👤 Telegram ID: `%d`\n", user.TelegramID)
	message += fmt.Sprintf("📅 Регистрация: %s\n\n", user.CreatedAt.Format("02.01.2006"))

	if len(subscriptions) > 0 {
		message += "🔒 **Активные подписки:**\n"
		for _, sub := range subscriptions {
			message += fmt.Sprintf("• %s (%s) - до %s\n",
				sub.ServerName,
				sub.PlanName,
				sub.ExpiresAt.Format("02.01.2006 15:04"))
		}
	} else {
		message += "❌ **Нет активных подписок**\n"
		message += "Используйте кнопку \"🚀 Купить\" для приобретения подписки."
	}

	keyboard := b.createMainMenuKeyboard(user)
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleReferrals обрабатывает callback для рефералов
func (b *Bot) handleReferrals(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Получаем рефералов пользователя
	referrals, err := b.userService.GetReferrals(user.ID)
	if err != nil {
		b.logger.Error("Failed to get referrals", "error", err)
		referrals = []models.User{}
	}

	message := "🎁 **Реферальная программа**\n\n"
	message += fmt.Sprintf("Ваш реферальный код: `%s`\n\n", user.ReferralCode)
	message += "Приглашайте друзей и получайте бонусы!\n\n"
	message += fmt.Sprintf("👥 Приглашено пользователей: %d\n", len(referrals))

	if len(referrals) > 0 {
		message += "\n**Ваши рефералы:**\n"
		for i, ref := range referrals {
			if i >= 10 { // Показываем только первых 10
				message += fmt.Sprintf("... и еще %d пользователей\n", len(referrals)-10)
				break
			}
			username := "Без имени"
			if ref.Username != "" {
				username = "@" + ref.Username
			} else if ref.FirstName != "" {
				username = ref.FirstName
			}
			message += fmt.Sprintf("• %s (ID: %d)\n", username, ref.TelegramID)
		}
	}

	keyboard := b.createMainMenuKeyboard(user)
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleTrial обрабатывает callback для пробного периода
func (b *Bot) handleTrial(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Проверяем, использовал ли пользователь пробный период
	hasUsedTrial, err := b.subscriptionService.HasUsedTrial(user.ID)
	if err != nil {
		b.logger.Error("Failed to check trial usage", "error", err)
		message := "❌ Произошла ошибка при проверке пробного периода."
		keyboard := b.createMainMenuKeyboard(user)
		return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
	}

	if hasUsedTrial {
		message := "🎁 **Пробный период**\n\n" +
			"Вы уже использовали пробный период.\n" +
			"Используйте кнопку \"🚀 Купить\" для приобретения подписки."
		keyboard := b.createMainMenuKeyboard(user)
		return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
	}

	// Здесь должна быть логика активации пробного периода
	message := "🎁 **Пробный период**\n\n" +
		"Функция пробного периода будет реализована в следующих версиях.\n" +
		"Используйте кнопку \"🚀 Купить\" для приобретения подписки."

	keyboard := b.createMainMenuKeyboard(user)
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminCallback обрабатывает callback'ы админ-панели
func (b *Bot) handleAdminCallback(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Проверяем, является ли пользователь админом
	if !b.userService.IsAdmin(user.TelegramID) {
		message := "❌ У вас нет прав администратора"
		keyboard := b.createMainMenuKeyboard(user)
		return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
	}

	data := query.Data
	action := strings.TrimPrefix(data, "admin:")

	// Создаем сообщение для админ-обработчика
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: query.Message.Chat.ID},
		From: &tgbotapi.User{ID: query.From.ID},
	}

	// Обрабатываем различные действия админ-панели
	switch action {
	case "main":
		return b.adminHandler.Handle(message, user, "")
	case "stats":
		return b.adminHandler.Handle(message, user, "stats")
	case "users":
		return b.adminHandler.Handle(message, user, "users")
	case "find_user":
		return b.handleAdminFindUser(query, user)
	case "balance":
		return b.handleAdminBalance(query, user)
	case "promo":
		return b.handleAdminPromo(query, user)
	case "notify":
		return b.handleAdminNotify(query, user)
	case "logs":
		return b.handleAdminLogs(query, user)
	case "settings":
		return b.handleAdminSettings(query, user)
	default:
		message := "❌ Неизвестное действие админ-панели"
		keyboard := b.adminHandler.GetAdminKeyboard().CreateMainMenu()
		return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
	}
}

// handleAdminFindUser обрабатывает поиск пользователя
func (b *Bot) handleAdminFindUser(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "🔍 *Поиск пользователя*\n\n"
	message += "Выберите способ поиска:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreateUserManagementMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminBalance обрабатывает управление балансом
func (b *Bot) handleAdminBalance(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "💰 *Управление балансом*\n\n"
	message += "Выберите операцию:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreateBalanceMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminPromo обрабатывает управление промокодами
func (b *Bot) handleAdminPromo(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "🎟️ *Управление промокодами*\n\n"
	message += "Выберите действие:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreatePromoCodeMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminNotify обрабатывает уведомления
func (b *Bot) handleAdminNotify(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "📢 *Уведомления*\n\n"
	message += "Выберите тип уведомления:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreateNotificationMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminLogs обрабатывает логи
func (b *Bot) handleAdminLogs(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "📋 *Логи активности*\n\n"
	message += "Выберите тип логов:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreateLogsMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
}

// handleAdminSettings обрабатывает настройки
func (b *Bot) handleAdminSettings(query *tgbotapi.CallbackQuery, _ *models.User) error {
	message := "⚙️ *Настройки бота*\n\n"
	message += "Выберите раздел настроек:"

	keyboard := b.adminHandler.GetAdminKeyboard().CreateSettingsMenu()
	return utils.SendMessageWithKeyboard(query.Message.Chat.ID, message, keyboard, b.config.BotToken)
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
