package keyboards

import (
	"remnawave-tg-shop/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AdminMenuKeyboard создает клавиатуру админ-панели
type AdminMenuKeyboard struct{}

// NewAdminMenuKeyboard создает новый AdminMenuKeyboard
func NewAdminMenuKeyboard() *AdminMenuKeyboard {
	return &AdminMenuKeyboard{}
}

// CreateMainMenu создает главное меню админ-панели
func (k *AdminMenuKeyboard) CreateMainMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Статистика и пользователи
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "admin:stats"),
		tgbotapi.NewInlineKeyboardButtonData("👥 Пользователи", "admin:users"),
	})

	// Управление пользователями
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔍 Найти пользователя", "admin:find_user"),
		tgbotapi.NewInlineKeyboardButtonData("💰 Управление балансом", "admin:balance"),
	})

	// Промокоды и уведомления
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Промокоды", "admin:promo"),
		tgbotapi.NewInlineKeyboardButtonData("📢 Уведомления", "admin:notify"),
	})

	// Логи и настройки
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📋 Логи", "admin:logs"),
		tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "admin:settings"),
	})

	// Назад в главное меню
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "start"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateUserManagementMenu создает меню управления пользователями
func (k *AdminMenuKeyboard) CreateUserManagementMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Поиск и список
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск по ID", "admin:search_user_id"),
		tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск по username", "admin:search_username"),
	})

	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📋 Список пользователей", "admin:list_users"),
		tgbotapi.NewInlineKeyboardButtonData("📊 Статистика пользователей", "admin:user_stats"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateUserActionsMenu создает меню действий с пользователем
func (k *AdminMenuKeyboard) CreateUserActionsMenu(user *models.User) tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Основная информация
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("ℹ️ Информация", "admin:user_info"),
		tgbotapi.NewInlineKeyboardButtonData("💰 Баланс", "admin:user_balance"),
	})

	// Блокировка/разблокировка
	blockText := "✅ Разблокировать"
	if !user.IsBlocked {
		blockText = "🚫 Заблокировать"
	}
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData(blockText, "admin:toggle_block"),
	})

	// Подписки и платежи
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔒 Подписки", "admin:user_subscriptions"),
		tgbotapi.NewInlineKeyboardButtonData("💳 Платежи", "admin:user_payments"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:users"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreatePromoCodeMenu создает меню управления промокодами
func (k *AdminMenuKeyboard) CreatePromoCodeMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Создание и просмотр
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("➕ Создать промокод", "admin:promo_create"),
		tgbotapi.NewInlineKeyboardButtonData("📋 Список промокодов", "admin:promo_list"),
	})

	// Статистика
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📊 Статистика промокодов", "admin:promo_stats"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateBalanceMenu создает меню управления балансом
func (k *AdminMenuKeyboard) CreateBalanceMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Операции с балансом
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("➕ Пополнить", "admin:balance_add"),
		tgbotapi.NewInlineKeyboardButtonData("➖ Списать", "admin:balance_subtract"),
	})

	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔢 Установить сумму", "admin:balance_set"),
		tgbotapi.NewInlineKeyboardButtonData("📊 История операций", "admin:balance_history"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateNotificationMenu создает меню уведомлений
func (k *AdminMenuKeyboard) CreateNotificationMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Типы уведомлений
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📢 Всем пользователям", "admin:notify_all"),
		tgbotapi.NewInlineKeyboardButtonData("👤 Конкретному пользователю", "admin:notify_user"),
	})

	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📊 Статистика уведомлений", "admin:notify_stats"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateLogsMenu создает меню логов
func (k *AdminMenuKeyboard) CreateLogsMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Типы логов
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📋 Все логи", "admin:logs_all"),
		tgbotapi.NewInlineKeyboardButtonData("👤 Логи пользователя", "admin:logs_user"),
	})

	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск по действию", "admin:logs_search"),
		tgbotapi.NewInlineKeyboardButtonData("📊 Статистика логов", "admin:logs_stats"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}

// CreateSettingsMenu создает меню настроек
func (k *AdminMenuKeyboard) CreateSettingsMenu() tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Настройки бота
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🤖 Настройки бота", "admin:settings_bot"),
		tgbotapi.NewInlineKeyboardButtonData("💳 Платежные системы", "admin:settings_payments"),
	})

	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Настройки промокодов", "admin:settings_promo"),
		tgbotapi.NewInlineKeyboardButtonData("📢 Настройки уведомлений", "admin:settings_notify"),
	})

	// Назад
	keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "admin:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}
