package commands

import (
	"fmt"
	"remnawave-tg-shop/internal/bot/keyboards"
	"remnawave-tg-shop/internal/bot/utils"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AdminHandler обрабатывает админские команды
type AdminHandler struct {
	config              *config.Config
	userService         services.UserService
	subscriptionService services.SubscriptionService
	paymentService      services.PaymentService
	promoCodeService    services.IPromoCodeService
	notificationService services.INotificationService
	activityLogService  services.IActivityLogService
	adminKeyboard       *keyboards.AdminMenuKeyboard
}

// NewAdminHandler создает новый AdminHandler
func NewAdminHandler(
	config *config.Config,
	userService services.UserService,
	subscriptionService services.SubscriptionService,
	paymentService services.PaymentService,
	promoCodeService services.IPromoCodeService,
	notificationService services.INotificationService,
	activityLogService services.IActivityLogService,
) *AdminHandler {
	return &AdminHandler{
		config:              config,
		userService:         userService,
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
		promoCodeService:    promoCodeService,
		notificationService: notificationService,
		activityLogService:  activityLogService,
		adminKeyboard:       keyboards.NewAdminMenuKeyboard(),
	}
}

// Handle обрабатывает админские команды
func (h *AdminHandler) Handle(message *tgbotapi.Message, user *models.User, args string) error {
	// Проверяем, является ли пользователь админом
	if !h.userService.IsAdmin(user.TelegramID) {
		return utils.SendMessage(message.Chat.ID, "❌ У вас нет прав администратора", h.config.BotToken)
	}

	// Логируем команду
	h.activityLogService.LogCommand(user.ID, "admin", args, "", "")

	// Парсим команду
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return h.showAdminMenu(message, user)
	}

	command := parts[0]
	commandArgs := strings.Join(parts[1:], " ")

	switch command {
	case "stats":
		return h.showStats(message, user)
	case "users":
		return h.showUsers(message, user, commandArgs)
	case "user":
		return h.showUser(message, user, commandArgs)
	case "block":
		return h.blockUser(message, user, commandArgs)
	case "unblock":
		return h.unblockUser(message, user, commandArgs)
	case "balance":
		return h.manageBalance(message, user, commandArgs)
	case "promo":
		return h.managePromoCodes(message, user, commandArgs)
	case "notify":
		return h.sendNotification(message, user, commandArgs)
	case "logs":
		return h.showLogs(message, user, commandArgs)
	case "help":
		return h.showAdminHelp(message, user)
	default:
		return h.showAdminMenu(message, user)
	}
}

// showAdminMenu показывает главное меню админ-панели
func (h *AdminHandler) showAdminMenu(message *tgbotapi.Message, _ *models.User) error {
	text := "🔧 *Админ-панель*\n\n"
	text += "Выберите раздел для управления ботом:"

	keyboard := h.adminKeyboard.CreateMainMenu()
	return utils.SendMessageWithKeyboard(message.Chat.ID, text, keyboard, h.config.BotToken)
}

// showStats показывает статистику
func (h *AdminHandler) showStats(message *tgbotapi.Message, _ *models.User) error {
	// Получаем статистику (здесь нужно будет реализовать методы в сервисах)
	text := "📊 *Статистика бота*\n\n"
	text += "👥 Пользователи: 0\n"
	text += "🔒 Активные подписки: 0\n"
	text += "💰 Общая выручка: 0₽\n"
	text += "📈 Выручка сегодня: 0₽\n"
	text += "🎟️ Промокоды: 0\n"
	text += "📢 Уведомления: 0"

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// showUsers показывает список пользователей
func (h *AdminHandler) showUsers(message *tgbotapi.Message, _ *models.User, searchQuery string) error {
	text := "👥 *Список пользователей*\n\n"

	if searchQuery != "" {
		text += fmt.Sprintf("Поиск: %s\n\n", searchQuery)
		// Здесь будет поиск пользователей
		text += "Результаты поиска будут здесь..."
	} else {
		text += "Последние 10 пользователей:\n\n"
		// Здесь будет список последних пользователей
		text += "Список пользователей будет здесь..."
	}

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// showUser показывает информацию о пользователе
func (h *AdminHandler) showUser(message *tgbotapi.Message, _ *models.User, userIDStr string) error {
	if userIDStr == "" {
		return utils.SendMessage(message.Chat.ID, "❌ Укажите ID пользователя", h.config.BotToken)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Неверный формат ID пользователя", h.config.BotToken)
	}

	targetUser, err := h.userService.GetUser(userID)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Пользователь не найден", h.config.BotToken)
	}

	text := "👤 *Информация о пользователе*\n\n"
	text += fmt.Sprintf("🆔 ID: %d\n", targetUser.TelegramID)
	text += fmt.Sprintf("👤 Имя: %s\n", targetUser.GetFullName())
	text += fmt.Sprintf("📱 Username: @%s\n", targetUser.Username)
	text += fmt.Sprintf("🌐 Язык: %s\n", targetUser.LanguageCode)
	text += fmt.Sprintf("💰 Баланс: %.2f₽\n", targetUser.Balance)
	text += fmt.Sprintf("🔗 Реферальный код: %s\n", targetUser.ReferralCode)
	text += fmt.Sprintf("🚫 Заблокирован: %t\n", targetUser.IsBlocked)
	text += fmt.Sprintf("👑 Админ: %t\n", targetUser.IsAdmin)
	text += fmt.Sprintf("📅 Регистрация: %s\n", targetUser.CreatedAt.Format("02.01.2006 15:04"))

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// blockUser блокирует пользователя
func (h *AdminHandler) blockUser(message *tgbotapi.Message, user *models.User, userIDStr string) error {
	if userIDStr == "" {
		return utils.SendMessage(message.Chat.ID, "❌ Укажите ID пользователя", h.config.BotToken)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Неверный формат ID пользователя", h.config.BotToken)
	}

	if err := h.userService.BlockUser(userID); err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Ошибка при блокировке пользователя", h.config.BotToken)
	}

	// Логируем действие
	h.activityLogService.LogActivity(user.ID, "admin_action", map[string]interface{}{
		"action":         "block_user",
		"target_user_id": userID,
	}, "", "")

	return utils.SendMessage(message.Chat.ID, "✅ Пользователь заблокирован", h.config.BotToken)
}

// unblockUser разблокирует пользователя
func (h *AdminHandler) unblockUser(message *tgbotapi.Message, user *models.User, userIDStr string) error {
	if userIDStr == "" {
		return utils.SendMessage(message.Chat.ID, "❌ Укажите ID пользователя", h.config.BotToken)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Неверный формат ID пользователя", h.config.BotToken)
	}

	if err := h.userService.UnblockUser(userID); err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Ошибка при разблокировке пользователя", h.config.BotToken)
	}

	// Логируем действие
	h.activityLogService.LogActivity(user.ID, "admin_action", map[string]interface{}{
		"action":         "unblock_user",
		"target_user_id": userID,
	}, "", "")

	return utils.SendMessage(message.Chat.ID, "✅ Пользователь разблокирован", h.config.BotToken)
}

// manageBalance управляет балансом пользователя
func (h *AdminHandler) manageBalance(message *tgbotapi.Message, user *models.User, args string) error {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		return utils.SendMessage(message.Chat.ID, "❌ Использование: /admin balance <id> <сумма>", h.config.BotToken)
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Неверный формат ID пользователя", h.config.BotToken)
	}

	amount, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Неверный формат суммы", h.config.BotToken)
	}

	targetUser, err := h.userService.GetUser(userID)
	if err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Пользователь не найден", h.config.BotToken)
	}

	var text string
	if amount > 0 {
		if err := h.userService.AddBalance(targetUser.ID, amount); err != nil {
			return utils.SendMessage(message.Chat.ID, "❌ Ошибка при пополнении баланса", h.config.BotToken)
		}
		text = fmt.Sprintf("✅ Баланс пользователя пополнен на %.2f₽", amount)
	} else {
		amount = -amount // Делаем положительным для вычитания
		if err := h.userService.SubtractBalance(targetUser.ID, amount); err != nil {
			return utils.SendMessage(message.Chat.ID, "❌ Ошибка при списании с баланса", h.config.BotToken)
		}
		text = fmt.Sprintf("✅ С баланса пользователя списано %.2f₽", amount)
	}

	// Логируем действие
	h.activityLogService.LogActivity(user.ID, "admin_action", map[string]interface{}{
		"action":         "manage_balance",
		"target_user_id": userID,
		"amount":         amount,
	}, "", "")

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// managePromoCodes управляет промокодами
func (h *AdminHandler) managePromoCodes(message *tgbotapi.Message, _ *models.User, _ string) error {
	text := "🎟️ *Управление промокодами*\n\n"
	text += "Доступные команды:\n"
	text += "• `/admin promo create <код> <тип> <значение> <макс_использований>` - Создать промокод\n"
	text += "• `/admin promo list` - Список промокодов\n"
	text += "• `/admin promo delete <id>` - Удалить промокод\n\n"
	text += "Типы промокодов:\n"
	text += "• `bonus_days` - Бонусные дни\n"
	text += "• `discount_percent` - Скидка в процентах\n"
	text += "• `discount_amount` - Скидка в рублях"

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// sendNotification отправляет уведомление
func (h *AdminHandler) sendNotification(message *tgbotapi.Message, user *models.User, notificationText string) error {
	if notificationText == "" {
		return utils.SendMessage(message.Chat.ID, "❌ Укажите текст уведомления", h.config.BotToken)
	}

	// Отправляем уведомление всем пользователям
	if err := h.notificationService.SendBulkNotification("admin_message", "Сообщение от администратора", notificationText, h.config.BotToken); err != nil {
		return utils.SendMessage(message.Chat.ID, "❌ Ошибка при отправке уведомления", h.config.BotToken)
	}

	// Логируем действие
	h.activityLogService.LogActivity(user.ID, "admin_action", map[string]interface{}{
		"action":  "send_notification",
		"message": notificationText,
	}, "", "")

	return utils.SendMessage(message.Chat.ID, "✅ Уведомление отправлено всем пользователям", h.config.BotToken)
}

// showLogs показывает логи активности
func (h *AdminHandler) showLogs(message *tgbotapi.Message, _ *models.User, userIDStr string) error {
	text := "📋 *Логи активности*\n\n"

	if userIDStr != "" {
		text += fmt.Sprintf("Логи пользователя %s:\n\n", userIDStr)
		// Здесь будет получение логов конкретного пользователя
		text += "Логи пользователя будут здесь..."
	} else {
		text += "Последние 10 записей:\n\n"
		// Здесь будет получение последних логов
		text += "Последние логи будут здесь..."
	}

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// showAdminHelp показывает справку по админским командам
func (h *AdminHandler) showAdminHelp(message *tgbotapi.Message, _ *models.User) error {
	text := "❓ *Справка по админским командам*\n\n"
	text += "🔧 *Основные команды:*\n"
	text += "`/admin` - Главное меню админ-панели\n"
	text += "`/admin stats` - Статистика бота\n\n"
	text += "👥 *Управление пользователями:*\n"
	text += "`/admin users` - Список всех пользователей\n"
	text += "`/admin users <поиск>` - Поиск пользователей\n"
	text += "`/admin user <id>` - Информация о пользователе\n"
	text += "`/admin block <id>` - Заблокировать пользователя\n"
	text += "`/admin unblock <id>` - Разблокировать пользователя\n\n"
	text += "💰 *Управление балансом:*\n"
	text += "`/admin balance <id> <сумма>` - Изменить баланс\n"
	text += "Положительная сумма - пополнение, отрицательная - списание\n\n"
	text += "🎟️ *Промокоды:*\n"
	text += "`/admin promo` - Управление промокодами\n\n"
	text += "📢 *Уведомления:*\n"
	text += "`/admin notify <сообщение>` - Отправить всем\n\n"
	text += "📋 *Логи:*\n"
	text += "`/admin logs` - Все логи\n"
	text += "`/admin logs <id>` - Логи пользователя"

	return utils.SendMessage(message.Chat.ID, text, h.config.BotToken)
}

// GetAdminKeyboard возвращает клавиатуру админ-панели
func (h *AdminHandler) GetAdminKeyboard() *keyboards.AdminMenuKeyboard {
	return h.adminKeyboard
}
