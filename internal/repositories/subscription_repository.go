package repositories

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// subscriptionRepository реализация SubscriptionRepository
type subscriptionRepository struct {
	db *gorm.DB
}

// Убеждаемся, что subscriptionRepository реализует SubscriptionRepository
var _ SubscriptionRepository = (*subscriptionRepository)(nil)

// NewSubscriptionRepository создает новый репозиторий подписок
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create создает новую подписку
func (r *subscriptionRepository) Create(subscription *models.Subscription) error {
	if err := r.db.Create(subscription).Error; err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

// GetByID получает подписку по ID
func (r *subscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := r.db.Preload("User").First(&subscription, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by ID: %w", err)
	}
	return &subscription, nil
}

// GetByUserID получает подписки пользователя
func (r *subscriptionRepository) GetByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by user ID: %w", err)
	}
	return subscriptions, nil
}

// GetActiveByUserID получает активные подписки пользователя
func (r *subscriptionRepository) GetActiveByUserID(userID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("user_id = ? AND status = ? AND expires_at > ?",
		userID, "active", time.Now()).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions by user ID: %w", err)
	}
	return subscriptions, nil
}

// Update обновляет подписку
func (r *subscriptionRepository) Update(subscription *models.Subscription) error {
	if err := r.db.Save(subscription).Error; err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

// Delete удаляет подписку
func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

// List получает список подписок с пагинацией
func (r *subscriptionRepository) List(limit, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetExpired получает истекшие подписки
func (r *subscriptionRepository) GetExpired() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("expires_at < ? AND status = ?", time.Now(), "active").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get expired subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetExpiringSoon получает подписки, истекающие в ближайшие дни
func (r *subscriptionRepository) GetExpiringSoon(days int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	expiryDate := time.Now().AddDate(0, 0, days)

	if err := r.db.Where("expires_at BETWEEN ? AND ? AND status = ?",
		time.Now(), expiryDate, "active").Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetUsersWithActiveSubscriptions получает пользователей с активными подписками
func (r *subscriptionRepository) GetUsersWithActiveSubscriptions() ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN subscriptions ON users.id = subscriptions.user_id").
		Where("subscriptions.status = ? AND subscriptions.expires_at > ?", "active", time.Now()).
		Group("users.id").
		Find(&users).Error
	return users, err
}

// GetUsersWithExpiredSubscriptions получает пользователей с истекшими подписками
func (r *subscriptionRepository) GetUsersWithExpiredSubscriptions() ([]models.User, error) {
	var users []models.User
	err := r.db.Joins("JOIN subscriptions ON users.id = subscriptions.user_id").
		Where("subscriptions.status = ? AND subscriptions.expires_at <= ?", "active", time.Now()).
		Group("users.id").
		Find(&users).Error
	return users, err
}
