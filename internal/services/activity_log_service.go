package services

import (
	"encoding/json"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type ActivityLogService struct {
	repo   repositories.ActivityLogRepository
	config *config.Config
}

func NewActivityLogService(repo repositories.ActivityLogRepository, config *config.Config) *ActivityLogService {
	return &ActivityLogService{
		repo:   repo,
		config: config,
	}
}

// LogActivity логирует активность пользователя
func (s *ActivityLogService) LogActivity(userID uuid.UUID, action string, data interface{}, ipAddress, userAgent string) error {
	var dataJSON string
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			dataJSON = ""
		} else {
			dataJSON = string(jsonData)
		}
	}

	log := &models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Data:      dataJSON,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(log)
}

// LogCommand логирует выполнение команды
func (s *ActivityLogService) LogCommand(userID uuid.UUID, command, args string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"command": command,
		"args":    args,
	}

	return s.LogActivity(userID, "command", data, ipAddress, userAgent)
}

// LogMessage логирует отправку сообщения
func (s *ActivityLogService) LogMessage(userID uuid.UUID, message string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"message": message,
	}

	return s.LogActivity(userID, "message", data, ipAddress, userAgent)
}

// LogCallback логирует callback запрос
func (s *ActivityLogService) LogCallback(userID uuid.UUID, callbackData string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"callback_data": callbackData,
	}

	return s.LogActivity(userID, "callback", data, ipAddress, userAgent)
}

// LogPayment логирует платеж
func (s *ActivityLogService) LogPayment(userID uuid.UUID, paymentID uuid.UUID, amount float64, method string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"payment_id": paymentID,
		"amount":     amount,
		"method":     method,
	}

	return s.LogActivity(userID, "payment", data, ipAddress, userAgent)
}

// LogSubscription логирует действия с подпиской
func (s *ActivityLogService) LogSubscription(userID uuid.UUID, subscriptionID uuid.UUID, action string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"subscription_id": subscriptionID,
		"action":          action,
	}

	return s.LogActivity(userID, "subscription", data, ipAddress, userAgent)
}

// LogPromoCode логирует использование промокода
func (s *ActivityLogService) LogPromoCode(userID uuid.UUID, promoCodeID uuid.UUID, code string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"promo_code_id": promoCodeID,
		"code":          code,
	}

	return s.LogActivity(userID, "promo_code", data, ipAddress, userAgent)
}

// LogReferral логирует реферальную активность
func (s *ActivityLogService) LogReferral(userID uuid.UUID, referredUserID uuid.UUID, action string, ipAddress, userAgent string) error {
	data := map[string]interface{}{
		"referred_user_id": referredUserID,
		"action":           action,
	}

	return s.LogActivity(userID, "referral", data, ipAddress, userAgent)
}

// GetUserActivity получает активность пользователя
func (s *ActivityLogService) GetUserActivity(userID uuid.UUID, limit, offset int) ([]models.ActivityLog, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// GetActivityByAction получает активность по типу действия
func (s *ActivityLogService) GetActivityByAction(action string, limit, offset int) ([]models.ActivityLog, error) {
	return s.repo.GetByAction(action, limit, offset)
}

// GetAllActivity получает всю активность
func (s *ActivityLogService) GetAllActivity(limit, offset int) ([]models.ActivityLog, error) {
	return s.repo.GetAll(limit, offset)
}

// GetActivityByDateRange получает активность за период
func (s *ActivityLogService) GetActivityByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error) {
	return s.repo.GetByDateRange(startDate, endDate, limit, offset)
}

// GetUserActivityByDateRange получает активность пользователя за период
func (s *ActivityLogService) GetUserActivityByDateRange(userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error) {
	return s.repo.GetByUserAndDateRange(userID, startDate, endDate, limit, offset)
}

// GetUserActivityCount получает количество записей активности пользователя
func (s *ActivityLogService) GetUserActivityCount(userID uuid.UUID) (int64, error) {
	return s.repo.CountByUserID(userID)
}

// GetActionCount получает количество записей по действию
func (s *ActivityLogService) GetActionCount(action string) (int64, error) {
	return s.repo.CountByAction(action)
}

// CleanupOldLogs удаляет старые записи активности
func (s *ActivityLogService) CleanupOldLogs(daysToKeep int) error {
	cutoffDate := time.Now().AddDate(0, 0, -daysToKeep)
	return s.repo.DeleteOldLogs(cutoffDate)
}
