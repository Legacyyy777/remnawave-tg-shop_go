package bot

import (
	"fmt"
	"strings"
	"time"

	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	"gopkg.in/telebot.v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot представляет телеграм-бота
type Bot struct {
	api                 *telebot.Bot
	config              *config.Config
	logger              logger.Logger
	userService         services.UserService
	subscriptionService services.SubscriptionService
	paymentService      services.PaymentService
}

// NewBot создает нового бота
func NewBot(cfg *config.Config, log logger.Logger, userService services.UserService, subscriptionService services.SubscriptionService, paymentService services.PaymentService) (*Bot, error) {
	pref := telebot.Settings{
		Token: cfg.BotToken,
		// Используем Long Polling для простоты
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	api, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	bot := &Bot{
		api:                 api,
		config:              cfg,
		logger:              log,
		userService:         userService,
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
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
	// Здесь можно добавить логику обработки обновлений
	// Пока что просто логируем
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
	user, err := b.getOrCreateUser(message.From)
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
			return b.handleStartCommandTgBot(message, user, args)
		case "help":
			return b.handleHelpCommandTgBot(message, user, args)
		case "balance":
			return b.handleBalanceCommandTgBot(message, user, args)
		case "subscriptions":
			return b.handleSubscriptionsCommandTgBot(message, user, args)
		case "referrals":
			return b.handleReferralsCommandTgBot(message, user, args)
		case "admin":
			return b.handleAdminCommandTgBot(message, user, args)
		default:
			return b.handleUnknownCommand(message, user, args)
		}
	}
	
	// Обрабатываем обычные сообщения
	return b.handleTextMessage(message, user)
}

// handleCallbackQuery обрабатывает callback queries
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) error {
	b.logger.Info("Handling callback query", "chat_id", query.Message.Chat.ID, "data", query.Data)
	
	// Получаем пользователя
	user, err := b.getOrCreateUser(query.From)
	if err != nil {
		b.logger.Error("Failed to get user", "error", err)
		return err
	}
	
	// Обрабатываем callback query
	return b.handleCallbackQueryData(query, user)
}

// getOrCreateUser получает или создает пользователя
func (b *Bot) getOrCreateUser(from *tgbotapi.User) (*models.User, error) {
	// Используем CreateOrGetUser для получения или создания пользователя
	user, err := b.userService.CreateOrGetUser(
		from.ID,
		from.UserName,
		from.FirstName,
		from.LastName,
		from.LanguageCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}
	
	return user, nil
}

// handleTextMessage обрабатывает обычные текстовые сообщения
func (b *Bot) handleTextMessage(message *tgbotapi.Message, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling text message", "chat_id", message.Chat.ID, "text", message.Text)
	return nil
}

// handleUnknownCommand обрабатывает неизвестные команды
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message, user *models.User, args string) error {
	// Отправляем сообщение о неизвестной команде
	msg := tgbotapi.NewMessage(message.Chat.ID, "❓ Неизвестная команда. Используйте /help для получения списка команд.")
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// handleCallbackQueryData обрабатывает данные callback query
func (b *Bot) handleCallbackQueryData(query *tgbotapi.CallbackQuery, user *models.User) error {
	data := query.Data
	
	// Обрабатываем различные типы callback'ов
	switch {
	case data == "balance":
		return b.handleBalanceCallbackTgBot(query, user)
	case data == "buy_subscription":
		return b.handleBuySubscriptionCallbackTgBot(query, user)
	case data == "my_subscriptions":
		return b.handleMySubscriptionsCallbackTgBot(query, user)
	case data == "referrals":
		return b.handleReferralsCallbackTgBot(query, user)
	case data == "promo_code":
		return b.handlePromoCodeCallbackTgBot(query, user)
	case data == "language":
		return b.handleLanguageCallbackTgBot(query, user)
	case data == "status":
		return b.handleStatusCallbackTgBot(query, user)
	case data == "support":
		return b.handleSupportCallbackTgBot(query, user)
	case data == "trial":
		return b.handleTrialCallbackTgBot(query, user)
	case data == "start":
		return b.handleStartCallbackTgBot(query, user)
	case strings.HasPrefix(data, "tariff_"):
		return b.handleTariffCallbackTgBot(query, user)
	case strings.HasPrefix(data, "payment_"):
		return b.handlePaymentCallbackTgBot(query, user)
	default:
		b.logger.Info("Unknown callback data", "data", data)
		return nil
	}
}

// setupHandlers настраивает обработчики команд и callback'ов
func (b *Bot) setupHandlers() {
	// Middleware для логирования и аутентификации
	b.api.Use(b.authMiddleware)

	// Команды
	b.api.Handle("/start", b.handleStartCommand)
	b.api.Handle("/help", b.handleHelpCommand)
	b.api.Handle("/balance", b.handleBalanceCommand)
	b.api.Handle("/subscriptions", b.handleSubscriptionsCommand)
	b.api.Handle("/referrals", b.handleReferralsCommand)
	b.api.Handle("/admin", b.handleAdminCommand)

	// Callback queries - используем строки вместо указателей на кнопки
	b.api.Handle("\fbalance", b.handleBalanceCallback)
	b.api.Handle("\fbuy_subscription", b.handleBuySubscriptionCallback)
	b.api.Handle("\fmy_subscriptions", b.handleMySubscriptionsCallback)
	b.api.Handle("\freferrals", b.handleReferralsCallback)
	b.api.Handle("\fpromo_code", b.handlePromoCodeCallback)
	b.api.Handle("\flanguage", b.handleLanguageCallback)
	b.api.Handle("\fstatus", b.handleStatusCallback)
	b.api.Handle("\fsupport", b.handleSupportCallback)
	b.api.Handle("\ftrial", b.handleTrialCallback)
	b.api.Handle("\fstart", b.handleStartCallback)

	// Tariff callbacks
	b.api.Handle("\ftariff_1", func(c telebot.Context) error { return b.handleTariffCallback(c, "tariff_1") })
	b.api.Handle("\ftariff_3", func(c telebot.Context) error { return b.handleTariffCallback(c, "tariff_3") })
	b.api.Handle("\ftariff_6", func(c telebot.Context) error { return b.handleTariffCallback(c, "tariff_6") })
	b.api.Handle("\ftariff_12", func(c telebot.Context) error { return b.handleTariffCallback(c, "tariff_12") })

	// Payment callbacks
	b.api.Handle("\fpayment_stars", func(c telebot.Context) error { return b.handlePaymentCallback(c, "payment_stars") })
	b.api.Handle("\fpayment_tribute", func(c telebot.Context) error { return b.handlePaymentCallback(c, "payment_tribute") })
	b.api.Handle("\fpayment_yookassa", func(c telebot.Context) error { return b.handlePaymentCallback(c, "payment_yookassa") })

	// Text messages
	b.api.Handle(telebot.OnText, b.handleTextMessageTelebot)
}

// authMiddleware - middleware для аутентификации и логирования
func (b *Bot) authMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user := c.Sender()
		if user == nil {
			return c.Send("❌ Ошибка аутентификации")
		}

		b.logger.Info("Request from user", "user_id", user.ID, "username", user.Username)

		// Создаем или получаем пользователя
		dbUser, err := b.userService.CreateOrGetUser(
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName,
			user.LanguageCode,
		)
		if err != nil {
			b.logger.Error("Failed to create/get user", "error", err)
			return c.Send("❌ Ошибка при получении данных пользователя")
		}

		// Проверяем блокировку
		if dbUser.IsBlocked {
			return c.Send("❌ Вы заблокированы и не можете использовать бота.")
		}

		// Сохраняем пользователя в контексте
		c.Set("user", dbUser)

		return next(c)
	}
}

// getUserFromContext получает пользователя из контекста
func (b *Bot) getUserFromContext(c telebot.Context) *models.User {
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}

// handleStartCommand обрабатывает команду /start
func (b *Bot) handleStartCommand(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	// Обработка реферального кода
	args := c.Message().Payload
	if args != "" {
		referralUser, err := b.userService.GetUserByReferralCode(args)
		if err == nil && referralUser != nil && referralUser.ID != user.ID {
			user.ReferredBy = &referralUser.ID
			b.userService.UpdateUser(user)
			b.userService.AddBalance(referralUser.ID, 50)
		}
	}

	// Формируем текст
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	// Создаем главное меню
	keyboard := b.createMainMenuKeyboard(user)

	return c.Send(text, keyboard)
}

// handleStartCommandTgBot обрабатывает команду /start для tgbotapi
func (b *Bot) handleStartCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	// Обработка реферального кода
	if args != "" {
		referralUser, err := b.userService.GetUserByReferralCode(args)
		if err == nil && referralUser != nil && referralUser.ID != user.ID {
			user.ReferredBy = &referralUser.ID
			b.userService.UpdateUser(user)
			b.userService.AddBalance(referralUser.ID, 50)
		}
	}

	// Формируем приветствие с именем пользователя
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	// Формируем кнопку с балансом
	balanceText := fmt.Sprintf("Баланс %.0f₽", user.Balance)

	// Создаем кнопки главного меню
	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	
	// Баланс
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("💰 " + balanceText, "balance"),
	})
	
	// Купить
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🚀 Купить", "buy_subscription"),
	})
	
	// Пробный период (если включен и пользователь еще не использовал)
	if b.config.Trial.Enabled {
		// Проверяем, использовал ли пользователь пробный период
		hasUsedTrial, err := b.subscriptionService.HasUsedTrial(user.ID)
		if err == nil && !hasUsedTrial {
			keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🎁 Пробный период", "trial"),
			})
		}
	}
	
	// Моя подписка - прямая кнопка миниаппа
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonURL("🔒 Моя подписка", b.config.MiniApp.URL),
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard

	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// createMainMenuKeyboard создает главное меню
func (b *Bot) createMainMenuKeyboard(user *models.User) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}

	// Баланс
	balanceText := fmt.Sprintf("💰 Баланс %.0f₽", user.Balance)
	balanceBtn := menu.Data(balanceText, "balance")

	// Основные кнопки
	buyBtn := menu.Data("🚀 Купить", "buy_subscription")

	// Пробный период (если доступен)
	var trialRow []telebot.Btn
	if b.config.Trial.Enabled {
		hasUsedTrial, err := b.subscriptionService.HasUsedTrial(user.ID)
		if err == nil && !hasUsedTrial {
			trialBtn := menu.Data("🎁 Пробный период", "trial")
			trialRow = append(trialRow, trialBtn)
		}
	}

	// Моя подписка - WebApp кнопка
	webAppBtn := menu.WebApp("🔒 Моя подписка", &telebot.WebApp{URL: b.config.MiniApp.URL})

	// Дополнительные кнопки
	referralsBtn := menu.Data("🎁 Рефералы", "referrals")
	promoBtn := menu.Data("🎟️ Промокод", "promo_code")
	langBtn := menu.Data("🌐 Язык", "language")
	statusBtn := menu.Data("📊 Статус", "status")
	supportBtn := menu.Data("💬 Поддержка", "support")

	// Формируем клавиатуру
	rows := []telebot.Row{
		{balanceBtn},
		{buyBtn},
	}

	if len(trialRow) > 0 {
		rows = append(rows, trialRow)
	}

	rows = append(rows, []telebot.Row{
		{webAppBtn},
		{referralsBtn, promoBtn},
		{langBtn, statusBtn},
		{supportBtn},
	}...)

	menu.Inline(rows...)
	return menu
}

// handleBalanceCommand обрабатывает команду /balance
func (b *Bot) handleBalanceCommand(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	text := fmt.Sprintf("💰 Ваш баланс: %.2f ₽\n\n", user.Balance)
	text += "💳 Пополнить баланс:"

	keyboard := b.createPaymentKeyboard()
	return c.Send(text, keyboard)
}

// createPaymentKeyboard создает клавиатуру платежей
func (b *Bot) createPaymentKeyboard() *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}

	var buttons []telebot.Btn

	if b.config.Payments.StarsEnabled {
		buttons = append(buttons, menu.Data("⭐ Telegram Stars", "payment_stars"))
	}
	if b.config.Payments.TributeEnabled {
		buttons = append(buttons, menu.Data("💎 Tribute", "payment_tribute"))
	}
	if b.config.Payments.YooKassaEnabled {
		buttons = append(buttons, menu.Data("💳 ЮKassa", "payment_yookassa"))
	}

	// Группируем по 2 кнопки в ряд
	var rows []telebot.Row
	for i := 0; i < len(buttons); i += 2 {
		if i+1 < len(buttons) {
			rows = append(rows, telebot.Row{buttons[i], buttons[i+1]})
		} else {
			rows = append(rows, telebot.Row{buttons[i]})
		}
	}

	// Добавляем кнопку "Назад"
	backBtn := menu.Data("🔙 Назад", "start")
	rows = append(rows, telebot.Row{backBtn})

	menu.Inline(rows...)
	return menu
}

// handleHelpCommand обрабатывает команду /help
func (b *Bot) handleHelpCommand(c telebot.Context) error {
	text := "❓ Помощь по использованию бота\n\n"
	text += "📋 Основные команды:\n"
	text += "/start - 🏠 Главное меню\n"
	text += "/balance - 💰 Ваш баланс\n"
	text += "/subscriptions - 📱 Мои подписки\n"
	text += "/referrals - 👥 Рефералы\n"
	text += "/help - ❓ Эта справка\n\n"
	text += "💡 Как пользоваться:\n"
	text += "1. Пополните баланс через Stars, Tribute или ЮKassa\n"
	text += "2. Выберите сервер и тарифный план\n"
	text += "3. Оплатите подписку\n"
	text += "4. Получите конфигурацию VPN\n\n"
	text += "🆘 Если у вас есть вопросы, обратитесь к администратору."

	return c.Send(text)
}

// handleSubscriptionsCommand обрабатывает команду /subscriptions
func (b *Bot) handleSubscriptionsCommand(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	subscriptions, err := b.subscriptionService.GetUserSubscriptions(user.ID)
	if err != nil {
		b.logger.Error("Failed to get user subscriptions", "error", err)
		return c.Send("❌ Ошибка при получении подписок.")
	}

	if len(subscriptions) == 0 {
		text := "📱 У вас пока нет активных подписок.\n\n"
		text += "🛒 Купить подписку:"

		menu := &telebot.ReplyMarkup{}
		buyBtn := menu.Data("🛒 Купить подписку", "buy_subscription")
		menu.Inline(menu.Row(buyBtn))

		return c.Send(text, menu)
	}

	text := "📱 Ваши подписки:\n\n"
	for i, sub := range subscriptions {
		status := "🟢 Активна"
		if !sub.IsActive() {
			status = "🔴 " + sub.GetStatusText()
		}

		text += fmt.Sprintf("%d. %s - %s\n", i+1, sub.ServerName, sub.PlanName)
		text += fmt.Sprintf("   Статус: %s\n", status)
		text += fmt.Sprintf("   Истекает: %s\n", sub.ExpiresAt.Format("02.01.2006 15:04"))
		if sub.IsActive() {
			text += fmt.Sprintf("   Осталось дней: %d\n", sub.GetDaysLeft())
		}
		text += "\n"
	}

	menu := &telebot.ReplyMarkup{}
	buyBtn := menu.Data("🛒 Купить подписку", "buy_subscription")
	menu.Inline(menu.Row(buyBtn))

	return c.Send(text, menu)
}

// handleReferralsCommand обрабатывает команду /referrals
func (b *Bot) handleReferralsCommand(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	referrals, err := b.userService.GetReferrals(user.ID)
	if err != nil {
		b.logger.Error("Failed to get referrals", "error", err)
		return c.Send("❌ Ошибка при получении рефералов.")
	}

	text := "👥 Реферальная программа\n\n"
	text += fmt.Sprintf("🔗 Ваша реферальная ссылка:\n")
	text += fmt.Sprintf("https://t.me/%s?start=%s\n\n", b.api.Me.Username, user.ReferralCode)
	text += fmt.Sprintf("👥 Количество рефералов: %d\n", len(referrals))
	text += "💰 За каждого реферала вы получаете 50 ₽ бонуса\n\n"

	if len(referrals) > 0 {
		text += "📋 Ваши рефералы:\n"
		for i, ref := range referrals {
			text += fmt.Sprintf("%d. %s\n", i+1, ref.GetDisplayName())
		}
	}

	return c.Send(text)
}

// handleAdminCommand обрабатывает команду /admin
func (b *Bot) handleAdminCommand(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	if !b.userService.IsAdmin(user.TelegramID) {
		return c.Send("❌ У вас нет прав администратора.")
	}

	text := "⚙️ Админ панель\n\n"
	text += "📊 Статистика:\n"
	text += "👥 Пользователи: [загрузка...]\n"
	text += "📱 Подписки: [загрузка...]\n"
	text += "💰 Доход: [загрузка...]\n\n"
	text += "Выберите действие:"

	menu := &telebot.ReplyMarkup{}
	usersBtn := menu.Data("👥 Пользователи", "admin_users")
	subsBtn := menu.Data("📱 Подписки", "admin_subscriptions")
	paymentsBtn := menu.Data("💰 Платежи", "admin_payments")
	statsBtn := menu.Data("📊 Статистика", "admin_stats")

	menu.Inline(
		menu.Row(usersBtn, subsBtn),
		menu.Row(paymentsBtn, statsBtn),
	)

	return c.Send(text, menu)
}

// handleTextMessageTelebot обрабатывает текстовые сообщения
func (b *Bot) handleTextMessageTelebot(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Send("❌ Ошибка получения данных пользователя")
	}

	switch c.Text() {
	case "🔙 Главное меню":
		return b.handleStartCommand(c)
	default:
		// Обработка промокодов и других сообщений
		return nil
	}
}

// Callback handlers

func (b *Bot) handleBalanceCallback(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	text := fmt.Sprintf("💰 Ваш баланс: %.2f ₽\n\n", user.Balance)
	text += "💳 Пополнить баланс:"

	keyboard := b.createPaymentKeyboard()
	
	return c.Edit(text, keyboard)
}

func (b *Bot) handleBuySubscriptionCallback(c telebot.Context) error {
	text := "🛒 Выберите тарифный план:\n\n"

	menu := &telebot.ReplyMarkup{}
	var buttons []telebot.Btn

	if b.config.Payments.Price1Month > 0 {
		text += fmt.Sprintf("1️⃣ 1 месяц - %d₽\n", b.config.Payments.Price1Month)
		buttons = append(buttons, menu.Data("1️⃣ 1 месяц", "tariff_1"))
	}
	if b.config.Payments.Price3Months > 0 {
		text += fmt.Sprintf("3️⃣ 3 месяца - %d₽\n", b.config.Payments.Price3Months)
		buttons = append(buttons, menu.Data("3️⃣ 3 месяца", "tariff_3"))
	}
	if b.config.Payments.Price6Months > 0 {
		text += fmt.Sprintf("6️⃣ 6 месяцев - %d₽\n", b.config.Payments.Price6Months)
		buttons = append(buttons, menu.Data("6️⃣ 6 месяцев", "tariff_6"))
	}
	if b.config.Payments.Price12Months > 0 {
		text += fmt.Sprintf("1️⃣2️⃣ 12 месяцев - %d₽\n", b.config.Payments.Price12Months)
		buttons = append(buttons, menu.Data("1️⃣2️⃣ 12 месяцев", "tariff_12"))
	}

	text += "\nВыберите тарифный план:"

	// Группируем кнопки по 2
	var rows []telebot.Row
	for i := 0; i < len(buttons); i += 2 {
		if i+1 < len(buttons) {
			rows = append(rows, telebot.Row{buttons[i], buttons[i+1]})
		} else {
			rows = append(rows, telebot.Row{buttons[i]})
		}
	}

	backBtn := menu.Data("🔙 Назад", "start")
	rows = append(rows, telebot.Row{backBtn})

	menu.Inline(rows...)
	return c.Edit(text, menu)
}

func (b *Bot) handleMySubscriptionsCallback(c telebot.Context) error {
	text := "📱 Управление подписками\n\n"
	text += "Нажмите на кнопку ниже, чтобы открыть мини-приложение для управления вашими подписками."

	// Создаем Reply клавиатуру с WebApp кнопкой
	menu := &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
	webAppBtn := menu.WebApp("🔒 Открыть мини-приложение", &telebot.WebApp{URL: b.config.MiniApp.URL})
	mainMenuBtn := menu.Text("🔙 Главное меню")

	menu.Reply(
		menu.Row(webAppBtn),
		menu.Row(mainMenuBtn),
	)

	return c.Send(text, menu)
}

func (b *Bot) handleReferralsCallback(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	referrals, err := b.userService.GetReferrals(user.ID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при получении рефералов"})
	}

	text := "👥 Реферальная программа\n\n"
	text += fmt.Sprintf("🔗 Ваша реферальная ссылка:\n")
	text += fmt.Sprintf("https://t.me/%s?start=%s\n\n", b.api.Me.Username, user.ReferralCode)
	text += fmt.Sprintf("👥 Количество рефералов: %d\n", len(referrals))
	text += "💰 За каждого реферала вы получаете 50 ₽ бонуса\n\n"

	if len(referrals) > 0 {
		text += "📋 Ваши рефералы:\n"
		for i, ref := range referrals {
			text += fmt.Sprintf("%d. %s\n", i+1, ref.GetDisplayName())
		}
	}

	menu := &telebot.ReplyMarkup{}
	backBtn := menu.Data("🔙 Назад", "start")
	menu.Inline(menu.Row(backBtn))

	return c.Edit(text, menu)
}

func (b *Bot) handlePromoCodeCallback(c telebot.Context) error {
	text := "🎟️ Промокоды\n\n"
	text += "Введите промокод для получения скидки или бонуса.\n\n"
	text += "💡 Промокоды можно получить:\n"
	text += "• От друзей\n"
	text += "• В рекламных акциях\n"
	text += "• За участие в конкурсах\n\n"
	text += "Просто отправьте промокод в чат."

	menu := &telebot.ReplyMarkup{}
	backBtn := menu.Data("🔙 Назад", "start")
	menu.Inline(menu.Row(backBtn))

	return c.Edit(text, menu)
}

func (b *Bot) handleLanguageCallback(c telebot.Context) error {
	text := "🌐 Выбор языка\n\n"
	text += "Выберите предпочитаемый язык интерфейса:\n\n"
	text += "🇷🇺 Русский (текущий)\n"
	text += "🇺🇸 English\n"
	text += "🇩🇪 Deutsch\n"
	text += "🇫🇷 Français\n\n"
	text += "Смена языка будет доступна в следующих версиях."

	menu := &telebot.ReplyMarkup{}
	ruBtn := menu.Data("🇷🇺 Русский", "lang_ru")
	enBtn := menu.Data("🇺🇸 English", "lang_en")
	deBtn := menu.Data("🇩🇪 Deutsch", "lang_de")
	frBtn := menu.Data("🇫🇷 Français", "lang_fr")
	backBtn := menu.Data("🔙 Назад", "start")

	menu.Inline(
		menu.Row(ruBtn, enBtn),
		menu.Row(deBtn, frBtn),
		menu.Row(backBtn),
	)

	return c.Edit(text, menu)
}

func (b *Bot) handleStatusCallback(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	text := "📊 Статус аккаунта\n\n"
	text += fmt.Sprintf("👤 Пользователь: %s\n", user.GetDisplayName())
	text += fmt.Sprintf("💰 Баланс: %.2f ₽\n", user.Balance)
	text += fmt.Sprintf("📅 Регистрация: %s\n", user.CreatedAt.Format("02.01.2006"))
	text += fmt.Sprintf("🔗 Реферальный код: %s\n", user.ReferralCode)

	if user.ReferredBy != nil {
		text += "🎁 Получен по реферальной ссылке\n"
	}

	text += "\n📱 Активные подписки:\n"
	subscriptions, err := b.subscriptionService.GetUserSubscriptions(user.ID)
	if err == nil && len(subscriptions) > 0 {
		for _, sub := range subscriptions {
			if sub.IsActive() {
				text += fmt.Sprintf("• %s - %s (до %s)\n", sub.ServerName, sub.PlanName, sub.ExpiresAt.Format("02.01.2006"))
			}
		}
	} else {
		text += "Нет активных подписок\n"
	}

	menu := &telebot.ReplyMarkup{}
	backBtn := menu.Data("🔙 Назад", "start")
	menu.Inline(menu.Row(backBtn))

	return c.Edit(text, menu)
}

func (b *Bot) handleSupportCallback(c telebot.Context) error {
	text := "💬 Поддержка\n\n"
	text += "Если у вас возникли вопросы или проблемы, обратитесь к нашей службе поддержки:\n\n"
	text += "📧 Email: support@remnawave.com\n"
	text += "💬 Telegram: @remnawave_support\n"
	text += "🕐 Время работы: 24/7\n\n"
	text += "📋 Часто задаваемые вопросы:\n"
	text += "• Как пополнить баланс?\n"
	text += "• Как получить конфигурацию VPN?\n"
	text += "• Как отменить подписку?\n"
	text += "• Проблемы с подключением\n\n"
	text += "Мы ответим в течение 15 минут!"

	menu := &telebot.ReplyMarkup{}
	backBtn := menu.Data("🔙 Назад", "start")
	menu.Inline(menu.Row(backBtn))

	return c.Edit(text, menu)
}

func (b *Bot) handleTrialCallback(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	// Проверяем, использовал ли пользователь пробный период
	hasUsedTrial, err := b.subscriptionService.HasUsedTrial(user.ID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при проверке пробного периода"})
	}

	if hasUsedTrial {
		text := "🎁 Пробный период\n\n"
		text += "❌ Вы уже использовали пробный период.\n"
		text += "🛒 Купите подписку для продолжения использования VPN."

		menu := &telebot.ReplyMarkup{}
		buyBtn := menu.Data("🛒 Купить подписку", "buy_subscription")
		backBtn := menu.Data("🔙 Назад", "start")

		menu.Inline(
			menu.Row(buyBtn),
			menu.Row(backBtn),
		)

		return c.Edit(text, menu)
	}

	// Создаем пробную подписку
	err = b.subscriptionService.CreateTrialSubscription(user.ID, b.config.Trial.DurationDays, b.config.Trial.TrafficLimitGB, b.config.Trial.TrafficStrategy)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при создании пробной подписки"})
	}

	text := "🎁 Пробный период активирован!\n\n"
	text += fmt.Sprintf("⏰ Длительность: %d дней\n", b.config.Trial.DurationDays)
	if b.config.Trial.TrafficLimitGB > 0 {
		text += fmt.Sprintf("📊 Лимит трафика: %d ГБ\n", b.config.Trial.TrafficLimitGB)
	} else {
		text += "📊 Лимит трафика: безлимитный\n"
	}
	text += "\n🔗 Конфигурация VPN будет отправлена в течение 5 минут.\n"
	text += "📱 Используйте кнопку 'Моя подписка' для управления."

	menu := &telebot.ReplyMarkup{}
	subBtn := menu.Data("🔒 Моя подписка", "my_subscriptions")
	backBtn := menu.Data("🔙 Назад", "start")

	menu.Inline(
		menu.Row(subBtn),
		menu.Row(backBtn),
	)

	return c.Edit(text, menu)
}

func (b *Bot) handleStartCallback(c telebot.Context) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	// Формируем приветствие
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	keyboard := b.createMainMenuKeyboard(user)
	return c.Edit(text, keyboard)
}

func (b *Bot) handleTariffCallback(c telebot.Context, data string) error {
	user := b.getUserFromContext(c)
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка получения данных пользователя"})
	}

	var price int
	var duration int
	var planName string

	switch data {
	case "tariff_1":
		price = b.config.Payments.Price1Month
		duration = 1
		planName = "1 месяц"
	case "tariff_3":
		price = b.config.Payments.Price3Months
		duration = 3
		planName = "3 месяца"
	case "tariff_6":
		price = b.config.Payments.Price6Months
		duration = 6
		planName = "6 месяцев"
	case "tariff_12":
		price = b.config.Payments.Price12Months
		duration = 12
		planName = "12 месяцев"
	default:
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Неизвестный тариф"})
	}

	// Проверяем баланс
	if user.Balance < float64(price) {
		text := fmt.Sprintf("💰 Недостаточно средств\n\n")
		text += fmt.Sprintf("💳 Стоимость: %d₽\n", price)
		text += fmt.Sprintf("💰 Ваш баланс: %.2f₽\n", user.Balance)
		text += fmt.Sprintf("❌ Не хватает: %.2f₽\n\n", float64(price)-user.Balance)
		text += "Пополните баланс для покупки подписки."

		menu := &telebot.ReplyMarkup{}
		balanceBtn := menu.Data("💰 Пополнить баланс", "balance")
		backBtn := menu.Data("🔙 Назад", "buy_subscription")

		menu.Inline(
			menu.Row(balanceBtn),
			menu.Row(backBtn),
		)

		return c.Edit(text, menu)
	}

	// Создаем подписку
	err := b.subscriptionService.CreateSubscriptionByPlan(user.ID, planName, duration, price)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при создании подписки"})
	}

	// Списываем средства
	err = b.userService.DeductBalance(user.ID, float64(price))
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "❌ Ошибка при списании средств"})
	}

	text := "✅ Подписка успешно создана!\n\n"
	text += fmt.Sprintf("📋 План: %s\n", planName)
	text += fmt.Sprintf("💰 Стоимость: %d₽\n", price)
	text += fmt.Sprintf("💰 Остаток на балансе: %.2f₽\n\n", user.Balance-float64(price))
	text += "🔗 Конфигурация VPN будет отправлена в течение 5 минут.\n"
	text += "📱 Используйте кнопку 'Моя подписка' для управления."

	menu := &telebot.ReplyMarkup{}
	subBtn := menu.Data("🔒 Моя подписка", "my_subscriptions")
	backBtn := menu.Data("🔙 Назад", "start")

	menu.Inline(
		menu.Row(subBtn),
		menu.Row(backBtn),
	)

	return c.Edit(text, menu)
}

func (b *Bot) handlePaymentCallback(c telebot.Context, method string) error {
	switch method {
	case "payment_stars":
		return c.Respond(&telebot.CallbackResponse{Text: "⭐ Платеж через Stars пока не реализован"})
	case "payment_tribute":
		return c.Respond(&telebot.CallbackResponse{Text: "💎 Платеж через Tribute пока не реализован"})
	case "payment_yookassa":
		return c.Respond(&telebot.CallbackResponse{Text: "💳 Платеж через ЮKassa пока не реализован"})
	default:
		return c.Respond(&telebot.CallbackResponse{Text: "❓ Неизвестный способ оплаты"})
	}
}

// ===== Функции для tgbotapi =====

// handleBalanceCallbackTgBot обрабатывает callback для баланса
func (b *Bot) handleBalanceCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling balance callback", "user_id", user.ID)
	return nil
}

// handleBuySubscriptionCallbackTgBot обрабатывает callback для покупки подписки
func (b *Bot) handleBuySubscriptionCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling buy subscription callback", "user_id", user.ID)
	return nil
}

// handleMySubscriptionsCallbackTgBot обрабатывает callback для моих подписок
func (b *Bot) handleMySubscriptionsCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling my subscriptions callback", "user_id", user.ID)
	return nil
}

// handleReferralsCallbackTgBot обрабатывает callback для рефералов
func (b *Bot) handleReferralsCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling referrals callback", "user_id", user.ID)
	return nil
}

// handlePromoCodeCallbackTgBot обрабатывает callback для промокода
func (b *Bot) handlePromoCodeCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling promo code callback", "user_id", user.ID)
	return nil
}

// handleLanguageCallbackTgBot обрабатывает callback для языка
func (b *Bot) handleLanguageCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling language callback", "user_id", user.ID)
	return nil
}

// handleStatusCallbackTgBot обрабатывает callback для статуса
func (b *Bot) handleStatusCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling status callback", "user_id", user.ID)
	return nil
}

// handleSupportCallbackTgBot обрабатывает callback для поддержки
func (b *Bot) handleSupportCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling support callback", "user_id", user.ID)
	return nil
}

// handleTrialCallbackTgBot обрабатывает callback для пробного периода
func (b *Bot) handleTrialCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling trial callback", "user_id", user.ID)
	return nil
}

// handleStartCallbackTgBot обрабатывает callback для главного меню
func (b *Bot) handleStartCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling start callback", "user_id", user.ID)
	return nil
}

// handleTariffCallbackTgBot обрабатывает callback для тарифов
func (b *Bot) handleTariffCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling tariff callback", "user_id", user.ID, "data", query.Data)
	return nil
}

// handlePaymentCallbackTgBot обрабатывает callback для платежей
func (b *Bot) handlePaymentCallbackTgBot(query *tgbotapi.CallbackQuery, user *models.User) error {
	// Пока что просто логируем
	b.logger.Info("Handling payment callback", "user_id", user.ID, "data", query.Data)
	return nil
}

// ===== Функции команд для tgbotapi =====

// handleHelpCommandTgBot обрабатывает команду /help
func (b *Bot) handleHelpCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	text := "🤖 Доступные команды:\n\n"
	text += "/start - Главное меню\n"
	text += "/help - Список команд\n"
	text += "/balance - Баланс\n"
	text += "/subscriptions - Мои подписки\n"
	text += "/referrals - Рефералы\n"
	text += "/admin - Админ панель\n\n"
	text += "Используйте кнопки в меню для навигации."

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// handleBalanceCommandTgBot обрабатывает команду /balance
func (b *Bot) handleBalanceCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	text := fmt.Sprintf("💰 Ваш баланс: %.0f₽", user.Balance)
	
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// handleSubscriptionsCommandTgBot обрабатывает команду /subscriptions
func (b *Bot) handleSubscriptionsCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	text := "🔒 Ваши подписки:\n\n"
	text += "Для просмотра и управления подписками используйте кнопку 'Моя подписка' в главном меню."
	
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// handleReferralsCommandTgBot обрабатывает команду /referrals
func (b *Bot) handleReferralsCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	text := "🎁 Реферальная программа:\n\n"
	text += "Приглашайте друзей и получайте бонусы!\n"
	text += "Ваша реферальная ссылка: https://t.me/" + b.config.BotToken + "?start=" + user.ReferralCode
	
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}

// handleAdminCommandTgBot обрабатывает команду /admin
func (b *Bot) handleAdminCommandTgBot(message *tgbotapi.Message, user *models.User, args string) error {
	// Проверяем, является ли пользователь админом
	if !b.userService.IsAdmin(user.TelegramID) {
		text := "❌ У вас нет прав администратора."
		
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		
		// Отправляем сообщение через tgbotapi
		bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
		if err != nil {
			return fmt.Errorf("failed to create tgbotapi bot: %w", err)
		}
		
		_, err = bot.Send(msg)
		return err
	}
	
	text := "👑 Админ панель:\n\n"
	text += "Добро пожаловать в админ панель!\n"
	text += "Здесь будут доступны административные функции."
	
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	
	// Отправляем сообщение через tgbotapi
	bot, err := tgbotapi.NewBotAPI(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create tgbotapi bot: %w", err)
	}
	
	_, err = bot.Send(msg)
	return err
}