package services

import (
	"fmt"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type NotificationService struct {
	repo             repositories.NotificationRepository
	userRepo         repositories.UserRepository
	subscriptionRepo repositories.SubscriptionRepository
	config           *config.Config
}

func NewNotificationService(
	repo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	config *config.Config,
) *NotificationService {
	return &NotificationService{
		repo:             repo,
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		config:           config,
	}
}

// CreateNotification создает новое уведомление
func (s *NotificationService) CreateNotification(userID *uuid.UUID, notificationType, title, message string) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		IsRead:  false,
		IsSent:  false,
	}

	return notification, s.repo.Create(notification)
}

// SendNotification отправляет уведомление пользователю
func (s *NotificationService) SendNotification(notificationID uuid.UUID, botToken string) error {
	notification, err := s.repo.GetByID(notificationID)
	if err != nil {
		return fmt.Errorf("уведомление не найдено: %v", err)
	}

	if notification.IsSent {
		return fmt.Errorf("уведомление уже отправлено")
	}

	// Если уведомление для конкретного пользователя
	if notification.UserID != nil {
		user, err := s.userRepo.GetByID(*notification.UserID)
		if err != nil {
			return fmt.Errorf("пользователь не найден: %v", err)
		}

		// Отправляем сообщение пользователю
		message := fmt.Sprintf("🔔 *%s*\n\n%s", notification.Title, notification.Message)
		if err := sendMessage(user.TelegramID, message, botToken); err != nil {
			return fmt.Errorf("ошибка отправки сообщения: %v", err)
		}
	}

	// Помечаем как отправленное
	return s.repo.MarkAsSent(notificationID)
}

// SendBulkNotification отправляет уведомление всем пользователям
func (s *NotificationService) SendBulkNotification(notificationType, title, message string, botToken string) error {
	// Получаем всех пользователей
	users, err := s.userRepo.GetAll(0, 0) // 0, 0 = без лимитов
	if err != nil {
		return fmt.Errorf("ошибка получения пользователей: %v", err)
	}

	// Создаем уведомления для всех пользователей
	notifications := make([]models.Notification, 0, len(users))
	for _, user := range users {
		notifications = append(notifications, models.Notification{
			UserID:  &user.ID,
			Type:    notificationType,
			Title:   title,
			Message: message,
			IsRead:  false,
			IsSent:  false,
		})
	}

	// Создаем уведомления в базе
	if err := s.repo.CreateBulk(notifications); err != nil {
		return fmt.Errorf("ошибка создания уведомлений: %v", err)
	}

	// Отправляем уведомления
	for _, user := range users {
		messageText := fmt.Sprintf("🔔 *%s*\n\n%s", title, message)
		if err := sendMessage(user.TelegramID, messageText, botToken); err != nil {
			// Логируем ошибку, но продолжаем отправку остальным
			fmt.Printf("Ошибка отправки уведомления пользователю %d: %v\n", user.TelegramID, err)
		}
	}

	return nil
}

// SendToUsersWithActiveSubscriptions отправляет уведомление пользователям с активными подписками
func (s *NotificationService) SendToUsersWithActiveSubscriptions(notificationType, title, message string, botToken string) error {
	// Получаем пользователей с активными подписками
	users, err := s.subscriptionRepo.GetUsersWithActiveSubscriptions()
	if err != nil {
		return fmt.Errorf("ошибка получения пользователей с активными подписками: %v", err)
	}

	// Создаем уведомления
	notifications := make([]models.Notification, 0, len(users))
	for _, user := range users {
		notifications = append(notifications, models.Notification{
			UserID:  &user.ID,
			Type:    notificationType,
			Title:   title,
			Message: message,
			IsRead:  false,
			IsSent:  false,
		})
	}

	// Создаем уведомления в базе
	if err := s.repo.CreateBulk(notifications); err != nil {
		return fmt.Errorf("ошибка создания уведомлений: %v", err)
	}

	// Отправляем уведомления
	for _, user := range users {
		messageText := fmt.Sprintf("🔔 *%s*\n\n%s", title, message)
		if err := sendMessage(user.TelegramID, messageText, botToken); err != nil {
			fmt.Printf("Ошибка отправки уведомления пользователю %d: %v\n", user.TelegramID, err)
		}
	}

	return nil
}

// SendToUsersWithExpiredSubscriptions отправляет уведомление пользователям с истекшими подписками
func (s *NotificationService) SendToUsersWithExpiredSubscriptions(notificationType, title, message string, botToken string) error {
	// Получаем пользователей с истекшими подписками
	users, err := s.subscriptionRepo.GetUsersWithExpiredSubscriptions()
	if err != nil {
		return fmt.Errorf("ошибка получения пользователей с истекшими подписками: %v", err)
	}

	// Создаем уведомления
	notifications := make([]models.Notification, 0, len(users))
	for _, user := range users {
		notifications = append(notifications, models.Notification{
			UserID:  &user.ID,
			Type:    notificationType,
			Title:   title,
			Message: message,
			IsRead:  false,
			IsSent:  false,
		})
	}

	// Создаем уведомления в базе
	if err := s.repo.CreateBulk(notifications); err != nil {
		return fmt.Errorf("ошибка создания уведомлений: %v", err)
	}

	// Отправляем уведомления
	for _, user := range users {
		messageText := fmt.Sprintf("🔔 *%s*\n\n%s", title, message)
		if err := sendMessage(user.TelegramID, messageText, botToken); err != nil {
			fmt.Printf("Ошибка отправки уведомления пользователю %d: %v\n", user.TelegramID, err)
		}
	}

	return nil
}

// CheckExpiringSubscriptions проверяет истекающие подписки и отправляет уведомления
func (s *NotificationService) CheckExpiringSubscriptions(botToken string) error {
	if !s.config.Notifications.Enabled {
		return nil
	}

	// Получаем пользователей с истекающими подписками
	users, err := s.repo.GetExpiringSubscriptions(s.config.Notifications.ExpiringDaysBefore)
	if err != nil {
		return fmt.Errorf("ошибка получения пользователей с истекающими подписками: %v", err)
	}

	for _, user := range users {
		// Получаем подписки пользователя
		subscriptions, err := s.subscriptionRepo.GetByUserID(user.ID)
		if err != nil {
			continue
		}

		for _, subscription := range subscriptions {
			if subscription.IsActive() {
				daysLeft := subscription.GetDaysLeft()
				if daysLeft <= s.config.Notifications.ExpiringDaysBefore && daysLeft > 0 {
					title := "⚠️ Подписка истекает"
					message := fmt.Sprintf("Ваша подписка истекает через %d дн. Продлите её, чтобы не потерять доступ.", daysLeft)

					// Создаем уведомление
					notification := &models.Notification{
						UserID:  &user.ID,
						Type:    "subscription_expiring",
						Title:   title,
						Message: message,
						IsRead:  false,
						IsSent:  false,
					}

					if err := s.repo.Create(notification); err != nil {
						continue
					}

					// Отправляем уведомление
					messageText := fmt.Sprintf("🔔 *%s*\n\n%s", title, message)
					if err := sendMessage(user.TelegramID, messageText, botToken); err != nil {
						fmt.Printf("Ошибка отправки уведомления пользователю %d: %v\n", user.TelegramID, err)
					}
				}
			}
		}
	}

	return nil
}

// GetNotificationsByUserID получает уведомления пользователя
func (s *NotificationService) GetNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// MarkAsRead помечает уведомление как прочитанное
func (s *NotificationService) MarkAsRead(notificationID uuid.UUID) error {
	return s.repo.MarkAsRead(notificationID)
}

// GetUnreadCount получает количество непрочитанных уведомлений пользователя
func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.repo.CountUnreadByUserID(userID)
}

// sendMessage отправляет сообщение через tgbotapi
func sendMessage(chatID int64, text string, botToken string) error {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)
	return err
}
