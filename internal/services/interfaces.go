package services

import (
	"remnawave-tg-shop/internal/models"
	"time"

	"github.com/google/uuid"
)

// UserService интерфейс для работы с пользователями
type UserService interface {
	CreateOrGetUser(telegramID int64, username, firstName, lastName, languageCode string) (*models.User, error)
	GetUser(telegramID int64) (*models.User, error)
	GetUserByReferralCode(code string) (*models.User, error)
	UpdateUser(user *models.User) error
	BlockUser(telegramID int64) error
	UnblockUser(telegramID int64) error
	AddBalance(userID uuid.UUID, amount float64) error
	SubtractBalance(userID uuid.UUID, amount float64) error
	DeductBalance(userID uuid.UUID, amount float64) error
	GetReferrals(userID uuid.UUID) ([]models.User, error)
	SearchUsers(query string, limit int) ([]models.User, error)
	IsAdmin(telegramID int64) bool
}

// SubscriptionService интерфейс для работы с подписками
type SubscriptionService interface {
	CreateSubscription(userID uuid.UUID, serverID, planID int) (*models.Subscription, error)
	CreateSubscriptionByPlan(userID uuid.UUID, planName string, durationMonths, price int) error
	CreateTrialSubscription(userID uuid.UUID, durationDays, trafficLimitGB int, trafficStrategy string) error
	HasUsedTrial(userID uuid.UUID) (bool, error)
	GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error)
	GetActiveSubscriptions(userID uuid.UUID) ([]models.Subscription, error)
	GetSubscription(id uuid.UUID) (*models.Subscription, error)
	UpdateSubscription(subscription *models.Subscription) error
	CancelSubscription(id uuid.UUID) error
	GetExpiredSubscriptions() ([]models.Subscription, error)
	GetExpiringSoon(days int) ([]models.Subscription, error)
}

// PaymentService интерфейс для работы с платежами
type PaymentService interface {
	CreatePayment(userID uuid.UUID, amount float64, method, description string) (*models.Payment, error)
	GetPayment(id uuid.UUID) (*models.Payment, error)
	UpdatePaymentStatus(id uuid.UUID, status string) error
	GetUserPayments(userID uuid.UUID) ([]models.Payment, error)
	ProcessStarsPayment(userID uuid.UUID, amount float64) error
	ProcessTributePayment(userID uuid.UUID, amount float64) error
	ProcessYooKassaPayment(userID uuid.UUID, amount float64) error
}

// ServerService интерфейс для работы с серверами
type ServerService interface {
	GetServers() ([]models.Server, error)
	GetServer(id int) (*models.Server, error)
	GetPlans(serverID int) ([]models.Plan, error)
	GetPlan(id int) (*models.Plan, error)
	SyncServers() error
	SyncPlans(serverID int) error
}

// IPromoCodeService интерфейс для работы с промокодами
type IPromoCodeService interface {
	CreatePromoCode(code, promoType string, value float64, maxUses int, validFrom, validUntil *time.Time, description string, createdBy uuid.UUID) (*models.PromoCode, error)
	GeneratePromoCode(promoType string, value float64, maxUses int, validFrom, validUntil *time.Time, description string, createdBy uuid.UUID) (*models.PromoCode, error)
	ApplyPromoCode(userID uuid.UUID, code string) (*models.PromoCode, error)
	GetPromoCode(code string) (*models.PromoCode, error)
	GetPromoCodeByID(id uuid.UUID) (*models.PromoCode, error)
	GetAllPromoCodes(limit, offset int) ([]models.PromoCode, error)
	UpdatePromoCode(promoCode *models.PromoCode) error
	DeletePromoCode(id uuid.UUID) error
	GetValidPromoCodes() ([]models.PromoCode, error)
}

// INotificationService интерфейс для работы с уведомлениями
type INotificationService interface {
	CreateNotification(userID *uuid.UUID, notificationType, title, message string) (*models.Notification, error)
	SendNotification(notificationID uuid.UUID, botToken string) error
	SendBulkNotification(notificationType, title, message string, botToken string) error
	SendToUsersWithActiveSubscriptions(notificationType, title, message string, botToken string) error
	SendToUsersWithExpiredSubscriptions(notificationType, title, message string, botToken string) error
	CheckExpiringSubscriptions(botToken string) error
	GetNotificationsByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
	GetUnreadCount(userID uuid.UUID) (int64, error)
}

// IActivityLogService интерфейс для работы с логами активности
type IActivityLogService interface {
	LogActivity(userID uuid.UUID, action string, data interface{}, ipAddress, userAgent string) error
	LogCommand(userID uuid.UUID, command, args string, ipAddress, userAgent string) error
	LogMessage(userID uuid.UUID, message string, ipAddress, userAgent string) error
	LogCallback(userID uuid.UUID, callbackData string, ipAddress, userAgent string) error
	LogPayment(userID uuid.UUID, paymentID uuid.UUID, amount float64, method string, ipAddress, userAgent string) error
	LogSubscription(userID uuid.UUID, subscriptionID uuid.UUID, action string, ipAddress, userAgent string) error
	LogPromoCode(userID uuid.UUID, promoCodeID uuid.UUID, code string, ipAddress, userAgent string) error
	LogReferral(userID uuid.UUID, referredUserID uuid.UUID, action string, ipAddress, userAgent string) error
	GetUserActivity(userID uuid.UUID, limit, offset int) ([]models.ActivityLog, error)
	GetActivityByAction(action string, limit, offset int) ([]models.ActivityLog, error)
	GetAllActivity(limit, offset int) ([]models.ActivityLog, error)
	GetActivityByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error)
	GetUserActivityByDateRange(userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error)
	GetUserActivityCount(userID uuid.UUID) (int64, error)
	GetActionCount(action string) (int64, error)
	CleanupOldLogs(daysToKeep int) error
}
