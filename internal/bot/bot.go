package bot

import (
	"fmt"
	"strings"
	"time"

	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot представляет телеграм-бота
type Bot struct {
	api            *tgbotapi.BotAPI
	config         *config.Config
	logger         logger.Logger
	userService    services.UserService
	subscriptionService services.SubscriptionService
	paymentService services.PaymentService
}

// NewBot создает нового бота
func NewBot(cfg *config.Config, log logger.Logger, userService services.UserService, subscriptionService services.SubscriptionService, paymentService services.PaymentService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = cfg.LogLevel == "debug"

	bot := &Bot{
		api:                api,
		config:             cfg,
		logger:             log,
		userService:        userService,
		subscriptionService: subscriptionService,
		paymentService:     paymentService,
	}

	return bot, nil
}

// Start запускает бота
func (b *Bot) Start() error {
	b.logger.Info("Starting Telegram bot...")

	// Настраиваем webhook если указан URL
	if b.config.BotWebhookURL != "" {
		return b.startWebhook()
	}

	// Запускаем в режиме polling
	return b.startPolling()
}

// startWebhook запускает бота в режиме webhook
func (b *Bot) startWebhook() error {
	webhook, err := tgbotapi.NewWebhook(b.config.BotWebhookURL)
	if err != nil {
		return fmt.Errorf("failed to create webhook: %w", err)
	}

	_, err = b.api.Request(webhook)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	info, err := b.api.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	if info.LastErrorDate != 0 {
		b.logger.Error("Webhook error", "message", info.LastErrorMessage)
	}

	b.logger.Info("Bot started in webhook mode", "url", b.config.BotWebhookURL)
	return nil
}

// startPolling запускает бота в режиме polling
func (b *Bot) startPolling() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	b.logger.Info("Bot started in polling mode")

	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

// HandleUpdate обрабатывает обновление от Telegram (публичный метод)
func (b *Bot) HandleUpdate(update tgbotapi.Update) error {
	return b.handleUpdate(update)
}

// handleUpdate обрабатывает обновление от Telegram
func (b *Bot) handleUpdate(update tgbotapi.Update) error {
	if update.Message != nil {
		b.handleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.handleCallbackQuery(update.CallbackQuery)
	}
	return nil
}

// handleMessage обрабатывает сообщения
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	b.logger.Info("Received message", "chat_id", message.Chat.ID, "text", message.Text, "from", message.From.UserName)
	
	// Игнорируем старые сообщения
	if message.Date < int(time.Now().Unix()-300) {
		b.logger.Info("Ignoring old message", "date", message.Date, "now", time.Now().Unix())
		return
	}

	// Создаем или получаем пользователя
	user, err := b.userService.CreateOrGetUser(
		message.From.ID,
		message.From.UserName,
		message.From.FirstName,
		message.From.LastName,
		message.From.LanguageCode,
	)
	if err != nil {
		b.logger.Error("Failed to create/get user", "error", err)
		return
	}

	// Проверяем, не заблокирован ли пользователь
	if user.IsBlocked {
		b.sendMessage(message.Chat.ID, "❌ Вы заблокированы и не можете использовать бота.")
		return
	}

	// Обрабатываем команды
	if message.IsCommand() {
		b.handleCommand(message, user)
		return
	}

	// Обрабатываем обычные сообщения
	b.handleTextMessage(message, user)
}

// handleCommand обрабатывает команды
func (b *Bot) handleCommand(message *tgbotapi.Message, user *models.User) {
	command := message.Command()
	args := message.CommandArguments()

	switch command {
	case "start":
		b.handleStartCommand(message, user, args)
	case "help":
		b.handleHelpCommand(message, user)
	case "balance":
		b.handleBalanceCommand(message, user)
	case "subscriptions":
		b.handleSubscriptionsCommand(message, user)
	case "referrals":
		b.handleReferralsCommand(message, user)
	case "admin":
		b.handleAdminCommand(message, user)
	default:
		b.sendMessage(message.Chat.ID, "❓ Неизвестная команда. Используйте /help для получения списка команд.")
	}
}

// handleStartCommand обрабатывает команду /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message, user *models.User, args string) {
	b.logger.Info("Handling start command", "chat_id", message.Chat.ID, "user_id", user.ID)
	
	// Формируем приветствие с именем пользователя
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	// Обработка реферального кода
	if args != "" {
		referralUser, err := b.userService.GetUserByReferralCode(args)
		if err == nil && referralUser != nil && referralUser.ID != user.ID {
			// Добавляем реферала
			user.ReferredBy = &referralUser.ID
			b.userService.UpdateUser(user)
			
			// Начисляем бонус рефереру
			b.userService.AddBalance(referralUser.ID, 50) // 50 рублей бонус
			
			text += "\n\n🎁 Вы получили бонус за переход по реферальной ссылке!\n"
			text += "💰 На ваш баланс начислено 50 рублей."
		}
	}

	// Формируем кнопку с балансом
	balanceText := fmt.Sprintf("Баланс %.0f₽", user.Balance)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 " + balanceText, "balance"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚀 Купить", "buy_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔒 Моя подписка", "my_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Рефералы", "referrals"),
			tgbotapi.NewInlineKeyboardButtonData("🎟️ Промокод", "promo_code"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "language"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Статус", "status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "support"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// handleHelpCommand обрабатывает команду /help
func (b *Bot) handleHelpCommand(message *tgbotapi.Message, user *models.User) {
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

	b.sendMessage(message.Chat.ID, text)
}

// handleBalanceCommand обрабатывает команду /balance
func (b *Bot) handleBalanceCommand(message *tgbotapi.Message, user *models.User) {
	text := fmt.Sprintf("💰 Ваш баланс: %.2f ₽\n\n", user.Balance)
	text += "💳 Пополнить баланс:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐ Telegram Stars", "payment_stars"),
			tgbotapi.NewInlineKeyboardButtonData("💎 Tribute", "payment_tribute"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 ЮKassa", "payment_yookassa"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// handleSubscriptionsCommand обрабатывает команду /subscriptions
func (b *Bot) handleSubscriptionsCommand(message *tgbotapi.Message, user *models.User) {
	subscriptions, err := b.subscriptionService.GetUserSubscriptions(user.ID)
	if err != nil {
		b.logger.Error("Failed to get user subscriptions", "error", err)
		b.sendMessage(message.Chat.ID, "❌ Ошибка при получении подписок.")
		return
	}

	if len(subscriptions) == 0 {
		text := "📱 У вас пока нет активных подписок.\n\n"
		text += "🛒 Купить подписку:"

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🛒 Купить подписку", "buy_subscription"),
			},
		)

		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ReplyMarkup = keyboard
		b.api.Send(msg)
		return
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🛒 Купить подписку", "buy_subscription"),
		},
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// handleReferralsCommand обрабатывает команду /referrals
func (b *Bot) handleReferralsCommand(message *tgbotapi.Message, user *models.User) {
	referrals, err := b.userService.GetReferrals(user.ID)
	if err != nil {
		b.logger.Error("Failed to get referrals", "error", err)
		b.sendMessage(message.Chat.ID, "❌ Ошибка при получении рефералов.")
		return
	}

	text := "👥 Реферальная программа\n\n"
	text += fmt.Sprintf("🔗 Ваша реферальная ссылка:\n")
	text += fmt.Sprintf("https://t.me/%s?start=%s\n\n", b.api.Self.UserName, user.ReferralCode)
	text += fmt.Sprintf("👥 Количество рефералов: %d\n", len(referrals))
	text += "💰 За каждого реферала вы получаете 50 ₽ бонуса\n\n"

	if len(referrals) > 0 {
		text += "📋 Ваши рефералы:\n"
		for i, ref := range referrals {
			text += fmt.Sprintf("%d. %s\n", i+1, ref.GetDisplayName())
		}
	}

	b.sendMessage(message.Chat.ID, text)
}

// handleAdminCommand обрабатывает команду /admin
func (b *Bot) handleAdminCommand(message *tgbotapi.Message, user *models.User) {
	b.logger.Info("Admin command received", 
		"user_telegram_id", user.TelegramID, 
		"user_id", user.ID,
		"username", user.Username)
	
	// Добавляем отладку конфигурации
	b.logger.Info("Config debug", 
		"admin_telegram_id", b.config.Admin.TelegramID,
		"admin_telegram_id_zero", b.config.Admin.TelegramID == 0)
		
	if !b.userService.IsAdmin(user.TelegramID) {
		b.sendMessage(message.Chat.ID, "❌ У вас нет прав администратора.")
		return
	}

	text := "⚙️ Админ панель\n\n"
	text += "📊 Статистика:\n"
	text += "👥 Пользователи: [загрузка...]\n"
	text += "📱 Подписки: [загрузка...]\n"
	text += "💰 Доход: [загрузка...]\n\n"
	text += "Выберите действие:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", "admin_users"),
			tgbotapi.NewInlineKeyboardButtonData("📱 Подписки", "admin_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Платежи", "admin_payments"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "admin_stats"),
		),
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// handleTextMessage обрабатывает обычные текстовые сообщения
func (b *Bot) handleTextMessage(message *tgbotapi.Message, user *models.User) {
	// Здесь можно добавить обработку обычных сообщений
	// Например, ответы на вопросы, поиск и т.д.
}

// handleCallbackQuery обрабатывает нажатия на inline кнопки
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	data := query.Data
	userID := query.From.ID

	// Получаем пользователя
	user, err := b.userService.GetUser(int64(userID))
	if err != nil || user == nil {
		b.answerCallbackQuery(query.ID, "❌ Ошибка при получении данных пользователя.")
		return
	}

	// Проверяем, не заблокирован ли пользователь
	if user.IsBlocked {
		b.answerCallbackQuery(query.ID, "❌ Вы заблокированы.")
		return
	}

	// Обрабатываем callback данные
	switch {
	case data == "balance":
		b.handleBalanceCallback(query, user)
	case data == "buy_subscription":
		b.handleBuySubscriptionCallback(query, user)
	case data == "my_subscriptions":
		b.handleMySubscriptionsCallback(query, user)
	case data == "referrals":
		b.handleReferralsCallback(query, user)
	case data == "promo_code":
		b.handlePromoCodeCallback(query, user)
	case data == "language":
		b.handleLanguageCallback(query, user)
	case data == "status":
		b.handleStatusCallback(query, user)
	case data == "support":
		b.handleSupportCallback(query, user)
	case data == "start":
		b.handleStartCallback(query, user)
	case strings.HasPrefix(data, "payment_"):
		b.handlePaymentCallback(query, user, data)
	case strings.HasPrefix(data, "admin_"):
		b.handleAdminCallback(query, user, data)
	default:
		b.answerCallbackQuery(query.ID, "❓ Неизвестное действие.")
	}
}

// sendMessage отправляет сообщение пользователю
func (b *Bot) sendMessage(chatID int64, text string) {
	b.logger.Info("Sending message", "chat_id", chatID, "text", text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Error("Failed to send message", "error", err, "chat_id", chatID)
	} else {
		b.logger.Info("Message sent successfully", "chat_id", chatID)
	}
}

// answerCallbackQuery отвечает на callback query
func (b *Bot) answerCallbackQuery(callbackQueryID string, text string) {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	_, err := b.api.Request(callback)
	if err != nil {
		b.logger.Error("Failed to answer callback query", "error", err)
	}
}

// editMessage редактирует сообщение
func (b *Bot) editMessage(chatID int64, messageID int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, err := b.api.Send(msg)
	if err != nil {
		b.logger.Error("Failed to edit message", "error", err)
	}
}

// handleBalanceCallback обрабатывает callback для баланса
func (b *Bot) handleBalanceCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	text := fmt.Sprintf("💰 Ваш баланс: %.2f ₽\n\n", user.Balance)
	text += "💳 Пополнить баланс:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⭐ Telegram Stars", "payment_stars"),
			tgbotapi.NewInlineKeyboardButtonData("💎 Tribute", "payment_tribute"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 ЮKassa", "payment_yookassa"),
		),
	)

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
	b.answerCallbackQuery(query.ID, "💰 Баланс обновлен")
}

// handleBuySubscriptionCallback обрабатывает callback для покупки подписки
func (b *Bot) handleBuySubscriptionCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	text := "🛒 Выберите сервер:\n\n"
	text += "1. 🇺🇸 США - 299₽/месяц\n"
	text += "2. 🇩🇪 Германия - 399₽/месяц\n"
	text += "3. 🇯🇵 Япония - 499₽/месяц\n\n"
	text += "Выберите сервер для покупки подписки:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇺🇸 США", "server_1"),
			tgbotapi.NewInlineKeyboardButtonData("🇩🇪 Германия", "server_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇯🇵 Япония", "server_3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "start"),
		),
	)

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
	b.answerCallbackQuery(query.ID, "🛒 Выберите сервер")
}

// handleMySubscriptionsCallback обрабатывает callback для моих подписок
func (b *Bot) handleMySubscriptionsCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	subscriptions, err := b.subscriptionService.GetUserSubscriptions(user.ID)
	if err != nil {
		b.answerCallbackQuery(query.ID, "❌ Ошибка при получении подписок")
		return
	}

	if len(subscriptions) == 0 {
		text := "📱 У вас пока нет активных подписок.\n\n"
		text += "🛒 Купить подписку:"

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("🛒 Купить подписку", "buy_subscription"),
			},
		)

		b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
		b.answerCallbackQuery(query.ID, "📱 Подписки не найдены")
		return
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🛒 Купить подписку", "buy_subscription"),
		},
	)

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
	b.answerCallbackQuery(query.ID, "📱 Подписки загружены")
}

// handleReferralsCallback обрабатывает callback для рефералов
func (b *Bot) handleReferralsCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	referrals, err := b.userService.GetReferrals(user.ID)
	if err != nil {
		b.answerCallbackQuery(query.ID, "❌ Ошибка при получении рефералов")
		return
	}

	text := "👥 Реферальная программа\n\n"
	text += fmt.Sprintf("🔗 Ваша реферальная ссылка:\n")
	text += fmt.Sprintf("https://t.me/%s?start=%s\n\n", b.api.Self.UserName, user.ReferralCode)
	text += fmt.Sprintf("👥 Количество рефералов: %d\n", len(referrals))
	text += "💰 За каждого реферала вы получаете 50 ₽ бонуса\n\n"

	if len(referrals) > 0 {
		text += "📋 Ваши рефералы:\n"
		for i, ref := range referrals {
			text += fmt.Sprintf("%d. %s\n", i+1, ref.GetDisplayName())
		}
	}

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, nil)
	b.answerCallbackQuery(query.ID, "👥 Рефералы загружены")
}

// handlePaymentCallback обрабатывает callback для платежей
func (b *Bot) handlePaymentCallback(query *tgbotapi.CallbackQuery, user *models.User, data string) {
	switch data {
	case "payment_stars":
		b.answerCallbackQuery(query.ID, "⭐ Платеж через Stars пока не реализован")
	case "payment_tribute":
		b.answerCallbackQuery(query.ID, "💎 Платеж через Tribute пока не реализован")
	case "payment_yookassa":
		b.answerCallbackQuery(query.ID, "💳 Платеж через ЮKassa пока не реализован")
	default:
		b.answerCallbackQuery(query.ID, "❓ Неизвестный способ оплаты")
	}
}

// handleAdminCallback обрабатывает callback для админ-панели
func (b *Bot) handleAdminCallback(query *tgbotapi.CallbackQuery, user *models.User, data string) {
	switch data {
	case "admin_users":
		b.answerCallbackQuery(query.ID, "👥 Управление пользователями пока не реализовано")
	case "admin_subscriptions":
		b.answerCallbackQuery(query.ID, "📱 Управление подписками пока не реализовано")
	case "admin_payments":
		b.answerCallbackQuery(query.ID, "💰 Управление платежами пока не реализовано")
	case "admin_stats":
		b.answerCallbackQuery(query.ID, "📊 Статистика пока не реализована")
	default:
		b.answerCallbackQuery(query.ID, "❓ Неизвестное действие админ-панели")
	}
}

// handlePromoCodeCallback обрабатывает callback для промокода
func (b *Bot) handlePromoCodeCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	text := "🎟️ Промокоды\n\n"
	text += "Введите промокод для получения скидки или бонуса.\n\n"
	text += "💡 Промокоды можно получить:\n"
	text += "• От друзей\n"
	text += "• В рекламных акциях\n"
	text += "• За участие в конкурсах\n\n"
	text += "Просто отправьте промокод в чат."

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, nil)
	b.answerCallbackQuery(query.ID, "🎟️ Введите промокод")
}

// handleLanguageCallback обрабатывает callback для смены языка
func (b *Bot) handleLanguageCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	text := "🌐 Выбор языка\n\n"
	text += "Выберите предпочитаемый язык интерфейса:\n\n"
	text += "🇷🇺 Русский (текущий)\n"
	text += "🇺🇸 English\n"
	text += "🇩🇪 Deutsch\n"
	text += "🇫🇷 Français\n\n"
	text += "Смена языка будет доступна в следующих версиях."

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang_ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇺🇸 English", "lang_en"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇩🇪 Deutsch", "lang_de"),
			tgbotapi.NewInlineKeyboardButtonData("🇫🇷 Français", "lang_fr"),
		),
	)

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
	b.answerCallbackQuery(query.ID, "🌐 Выберите язык")
}

// handleStatusCallback обрабатывает callback для статуса
func (b *Bot) handleStatusCallback(query *tgbotapi.CallbackQuery, user *models.User) {
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

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, nil)
	b.answerCallbackQuery(query.ID, "📊 Статус загружен")
}

// handleSupportCallback обрабатывает callback для поддержки
func (b *Bot) handleSupportCallback(query *tgbotapi.CallbackQuery, user *models.User) {
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

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, nil)
	b.answerCallbackQuery(query.ID, "💬 Поддержка готова помочь")
}

// handleStartCallback обрабатывает callback для возврата в главное меню
func (b *Bot) handleStartCallback(query *tgbotapi.CallbackQuery, user *models.User) {
	// Формируем приветствие с именем пользователя
	username := user.GetDisplayName()
	text := fmt.Sprintf("Привет, %s👋\n\n", username)
	text += "Что бы вы хотели сделать?"

	// Формируем кнопку с балансом
	balanceText := fmt.Sprintf("Баланс %.0f₽", user.Balance)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 " + balanceText, "balance"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚀 Купить", "buy_subscription"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔒 Моя подписка", "my_subscriptions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Рефералы", "referrals"),
			tgbotapi.NewInlineKeyboardButtonData("🎟️ Промокод", "promo_code"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌐 Язык", "language"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Статус", "status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Поддержка", "support"),
		),
	)

	b.editMessage(query.Message.Chat.ID, query.Message.MessageID, text, &keyboard)
	b.answerCallbackQuery(query.ID, "🏠 Главное меню")
}
