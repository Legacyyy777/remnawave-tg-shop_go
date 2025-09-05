package repositories

import (
	"remnawave-tg-shop/internal/models"
	"time"

	"github.com/google/uuid"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByTelegramID(telegramID int64) (*models.User, error)
	GetByReferralCode(code string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]models.User, error)
	Search(query string, limit int) ([]models.User, error)
	GetAll(limit, offset int) ([]models.User, error)
	GetReferrals(userID uuid.UUID) ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
}

// SubscriptionRepository интерфейс для работы с подписками
type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	GetByID(id uuid.UUID) (*models.Subscription, error)
	GetByUserID(userID uuid.UUID) ([]models.Subscription, error)
	GetActiveByUserID(userID uuid.UUID) ([]models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]models.Subscription, error)
	GetExpired() ([]models.Subscription, error)
	GetExpiringSoon(days int) ([]models.Subscription, error)
	GetUsersWithActiveSubscriptions() ([]models.User, error)
	GetUsersWithExpiredSubscriptions() ([]models.User, error)
}

// PaymentRepository интерфейс для работы с платежами
type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id uuid.UUID) (*models.Payment, error)
	GetByUserID(userID uuid.UUID) ([]models.Payment, error)
	GetByExternalID(externalID string) (*models.Payment, error)
	Update(payment *models.Payment) error
	GetByStatus(status string) ([]models.Payment, error)
	GetByMethod(method string) ([]models.Payment, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.Payment, error)
}

// PromoCodeRepository интерфейс для работы с промокодами
type PromoCodeRepository interface {
	Create(promoCode *models.PromoCode) error
	GetByCode(code string) (*models.PromoCode, error)
	GetByID(id uuid.UUID) (*models.PromoCode, error)
	GetAll(limit, offset int) ([]models.PromoCode, error)
	Update(promoCode *models.PromoCode) error
	Delete(id uuid.UUID) error
	IncrementUsage(id uuid.UUID) error
	GetValidPromoCodes() ([]models.PromoCode, error)
	CreateUsage(usage *models.PromoCodeUsage) error
	GetUsageByUserAndPromoCode(userID, promoCodeID uuid.UUID) (*models.PromoCodeUsage, error)
	GetUsageCountByPromoCode(promoCodeID uuid.UUID) (int64, error)
}

// NotificationRepository интерфейс для работы с уведомлениями
type NotificationRepository interface {
	Create(notification *models.Notification) error
	GetByID(id uuid.UUID) (*models.Notification, error)
	GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	GetUnreadByUserID(userID uuid.UUID) ([]models.Notification, error)
	GetUnsentNotifications(limit int) ([]models.Notification, error)
	GetByType(notificationType string, limit, offset int) ([]models.Notification, error)
	MarkAsRead(id uuid.UUID) error
	MarkAsSent(id uuid.UUID) error
	Update(notification *models.Notification) error
	Delete(id uuid.UUID) error
	GetExpiringSubscriptions(daysBefore int) ([]models.User, error)
	CreateBulk(notifications []models.Notification) error
	CountUnreadByUserID(userID uuid.UUID) (int64, error)
	DeleteOldNotifications(beforeDate time.Time) error
}

// ActivityLogRepository интерфейс для работы с логами активности
type ActivityLogRepository interface {
	Create(log *models.ActivityLog) error
	GetByUserID(userID uuid.UUID, limit, offset int) ([]models.ActivityLog, error)
	GetByAction(action string, limit, offset int) ([]models.ActivityLog, error)
	GetAll(limit, offset int) ([]models.ActivityLog, error)
	GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error)
	GetByUserAndDateRange(userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error)
	CountByUserID(userID uuid.UUID) (int64, error)
	CountByAction(action string) (int64, error)
	DeleteOldLogs(beforeDate time.Time) error
}
