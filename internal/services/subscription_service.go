package services

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"
	"remnawave-tg-shop/internal/services/remnawave"

	"github.com/google/uuid"
)

// subscriptionService реализация SubscriptionService
type subscriptionService struct {
	subscriptionRepo repositories.SubscriptionRepository
	remnawaveClient  *remnawave.Client
	logger           logger.Logger
}

// NewSubscriptionService создает новый сервис подписок
func NewSubscriptionService(subscriptionRepo repositories.SubscriptionRepository, remnawaveClient *remnawave.Client, log logger.Logger) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		remnawaveClient:  remnawaveClient,
		logger:           log,
	}
}

// CreateSubscription создает новую подписку
func (s *subscriptionService) CreateSubscription(userID uuid.UUID, serverID, planID int) (*models.Subscription, error) {
	// Получаем информацию о сервере и плане из Remnawave
	servers, err := s.remnawaveClient.GetServers()
	if err != nil {
		return nil, fmt.Errorf("failed to get servers: %w", err)
	}

	var serverName string
	for _, server := range servers {
		if server.ID == serverID {
			serverName = server.Name
			break
		}
	}

	plans, err := s.remnawaveClient.GetPlans(serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}

	var planName string
	var planDuration int
	for _, plan := range plans {
		if plan.ID == planID {
			planName = plan.Name
			planDuration = plan.Duration
			break
		}
	}

	if planName == "" {
		return nil, fmt.Errorf("plan not found")
	}

	// Создаем подписку в Remnawave
	_, err = s.remnawaveClient.CreateSubscription(0, serverID, planID) // userID будет 0, так как мы работаем с Telegram пользователями
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription in Remnawave: %w", err)
	}

	// Создаем подписку в нашей БД
	subscription := &models.Subscription{
		UserID:     userID,
		ServerID:   serverID,
		ServerName: serverName,
		PlanID:     planID,
		PlanName:   planName,
		Status:     "active",
		ExpiresAt:  time.Now().AddDate(0, 0, planDuration),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info("Subscription created", "user_id", userID, "server_id", serverID, "plan_id", planID)
	return subscription, nil
}

// GetUserSubscriptions получает подписки пользователя
func (s *subscriptionService) GetUserSubscriptions(userID uuid.UUID) ([]models.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetActiveSubscriptions получает активные подписки пользователя
func (s *subscriptionService) GetActiveSubscriptions(userID uuid.UUID) ([]models.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetSubscription получает подписку по ID
func (s *subscriptionService) GetSubscription(id uuid.UUID) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return subscription, nil
}

// UpdateSubscription обновляет подписку
func (s *subscriptionService) UpdateSubscription(subscription *models.Subscription) error {
	subscription.UpdatedAt = time.Now()
	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

// CancelSubscription отменяет подписку
func (s *subscriptionService) CancelSubscription(id uuid.UUID) error {
	subscription, err := s.subscriptionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription == nil {
		return fmt.Errorf("subscription not found")
	}

	subscription.Status = "cancelled"
	subscription.UpdatedAt = time.Now()

	if err := s.subscriptionRepo.Update(subscription); err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	// Отменяем подписку в Remnawave
	if err := s.remnawaveClient.DeleteSubscription(0); err != nil { // Здесь нужно передать правильный ID подписки из Remnawave
		s.logger.Warn("Failed to cancel subscription in Remnawave", "error", err)
	}

	s.logger.Info("Subscription cancelled", "subscription_id", id)
	return nil
}

// GetExpiredSubscriptions получает истекшие подписки
func (s *subscriptionService) GetExpiredSubscriptions() ([]models.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.GetExpired()
	if err != nil {
		return nil, fmt.Errorf("failed to get expired subscriptions: %w", err)
	}
	return subscriptions, nil
}

// GetExpiringSoon получает подписки, истекающие в ближайшие дни
func (s *subscriptionService) GetExpiringSoon(days int) ([]models.Subscription, error) {
	subscriptions, err := s.subscriptionRepo.GetExpiringSoon(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring subscriptions: %w", err)
	}
	return subscriptions, nil
}

// CreateSubscriptionByPlan создает подписку по плану
func (s *subscriptionService) CreateSubscriptionByPlan(userID uuid.UUID, planName string, durationMonths, price int) error {
	// Создаем подписку в нашей БД
	subscription := &models.Subscription{
		UserID:     userID,
		ServerID:   1, // По умолчанию сервер 1
		ServerName: "Default Server",
		PlanID:     1, // По умолчанию план 1
		PlanName:   planName,
		Status:     "active",
		ExpiresAt:  time.Now().AddDate(0, durationMonths, 0),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	s.logger.Info("Subscription created by plan", "user_id", userID, "plan_name", planName, "duration_months", durationMonths, "price", price)
	return nil
}

// CreateTrialSubscription создает пробную подписку
func (s *subscriptionService) CreateTrialSubscription(userID uuid.UUID, durationDays, trafficLimitGB int, trafficStrategy string) error {
	// Создаем пробную подписку в нашей БД
	subscription := &models.Subscription{
		UserID:     userID,
		ServerID:   1, // По умолчанию сервер 1
		ServerName: "Trial Server",
		PlanID:     0, // 0 для пробной подписки
		PlanName:   "Trial Plan",
		Status:     "active",
		ExpiresAt:  time.Now().AddDate(0, 0, durationDays),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return fmt.Errorf("failed to create trial subscription: %w", err)
	}

	s.logger.Info("Trial subscription created", "user_id", userID, "duration_days", durationDays, "traffic_limit_gb", trafficLimitGB, "traffic_strategy", trafficStrategy)
	return nil
}

// HasUsedTrial проверяет, использовал ли пользователь пробный период
func (s *subscriptionService) HasUsedTrial(userID uuid.UUID) (bool, error) {
	subscriptions, err := s.subscriptionRepo.GetByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	// Проверяем, есть ли у пользователя пробная подписка (PlanID = 0)
	for _, sub := range subscriptions {
		if sub.PlanID == 0 {
			return true, nil
		}
	}

	return false, nil
}