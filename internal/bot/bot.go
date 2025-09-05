package bot

import (
	"fmt"
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
	
	// Обработчики команд
	startHandler *commands.StartHandler
	helpHandler  *commands.HelpHandler
	
	// Обработчики callback'ов
	balanceHandler *callbacks.BalanceHandler
	
	// Обработчики сообщений
	textHandler *messages.TextHandler
	
	// Middleware
	authMiddleware *middleware.AuthMiddleware
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

	// Создаем обработчики
	startHandler := commands.NewStartHandler(cfg, userService, subscriptionService)
	helpHandler := commands.NewHelpHandler(cfg)
	balanceHandler := callbacks.NewBalanceHandler(cfg, userService)
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
		balanceHandler:      balanceHandler,
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
	case data == "start":
		return b.handleStartCallback(query, user)
	default:
		b.logger.Info("Unknown callback data", "data", data)
		return nil
	}
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
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message, user *models.User, args string) error {
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
