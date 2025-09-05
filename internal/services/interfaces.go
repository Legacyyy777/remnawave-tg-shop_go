package services

import (
	"remnawave-tg-shop/internal/models"

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
	GetReferrals(userID uuid.UUID) ([]models.User, error)
	SearchUsers(query string, limit int) ([]models.User, error)
	IsAdmin(telegramID int64) bool
}

// SubscriptionService интерфейс для работы с подписками
type SubscriptionService interface {
	CreateSubscription(userID uuid.UUID, serverID, planID int) (*models.Subscription, error)
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
