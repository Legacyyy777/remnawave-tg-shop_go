package repositories

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	GetByStatus(status string) ([]models.Subscription, error)
	GetByServerID(serverID int) ([]models.Subscription, error)
	GetByPlanID(planID int) ([]models.Subscription, error)
	GetExpiringSoon(days int) ([]models.Subscription, error)
}

// subscriptionRepository реализация SubscriptionRepository
type subscriptionRepository struct {
	db *gorm.DB
}

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

// GetByStatus получает подписки по статусу
func (r *subscriptionRepository) GetByStatus(status string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("status = ?", status).Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by status: %w", err)
	}
	return subscriptions, nil
}

// GetByServerID получает подписки по ID сервера
func (r *subscriptionRepository) GetByServerID(serverID int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("server_id = ?", serverID).Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by server ID: %w", err)
	}
	return subscriptions, nil
}

// GetByPlanID получает подписки по ID плана
func (r *subscriptionRepository) GetByPlanID(planID int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := r.db.Where("plan_id = ?", planID).Find(&subscriptions).Error; err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by plan ID: %w", err)
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
