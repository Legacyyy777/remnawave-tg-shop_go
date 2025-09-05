package repositories

import (
	"remnawave-tg-shop/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// Create создает новое уведомление
func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

// GetByID получает уведомление по ID
func (r *notificationRepository) GetByID(id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetByUserID получает уведомления пользователя
func (r *notificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// GetUnreadByUserID получает непрочитанные уведомления пользователя
func (r *notificationRepository) GetUnreadByUserID(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// GetUnsentNotifications получает неотправленные уведомления
func (r *notificationRepository) GetUnsentNotifications(limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("is_sent = ?", false).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

// GetByType получает уведомления по типу
func (r *notificationRepository) GetByType(notificationType string, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("type = ?", notificationType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

// MarkAsRead помечает уведомление как прочитанное
func (r *notificationRepository) MarkAsRead(id uuid.UUID) error {
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAsSent помечает уведомление как отправленное
func (r *notificationRepository) MarkAsSent(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_sent": true,
		"sent_at": &now,
	}).Error
}

// Update обновляет уведомление
func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

// Delete удаляет уведомление
func (r *notificationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, id).Error
}

// GetExpiringSubscriptions получает пользователей с истекающими подписками
func (r *notificationRepository) GetExpiringSubscriptions(daysBefore int) ([]models.User, error) {
	var users []models.User
	expiringDate := time.Now().AddDate(0, 0, daysBefore)

	err := r.db.Joins("JOIN subscriptions ON users.id = subscriptions.user_id").
		Where("subscriptions.status = ? AND subscriptions.expires_at BETWEEN ? AND ?",
			"active", time.Now(), expiringDate).
		Group("users.id").
		Find(&users).Error

	return users, err
}

// CreateBulk создает несколько уведомлений
func (r *notificationRepository) CreateBulk(notifications []models.Notification) error {
	return r.db.CreateInBatches(notifications, 100).Error
}

// CountUnreadByUserID подсчитывает количество непрочитанных уведомлений пользователя
func (r *notificationRepository) CountUnreadByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// DeleteOldNotifications удаляет старые уведомления
func (r *notificationRepository) DeleteOldNotifications(beforeDate time.Time) error {
	return r.db.Where("created_at < ?", beforeDate).Delete(&models.Notification{}).Error
}
