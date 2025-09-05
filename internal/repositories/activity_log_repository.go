package repositories

import (
	"remnawave-tg-shop/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

// Create создает новую запись в логе активности
func (r *activityLogRepository) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

// GetByUserID получает логи активности пользователя
func (r *activityLogRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByAction получает логи по типу действия
func (r *activityLogRepository) GetByAction(action string, limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetAll получает все логи с пагинацией
func (r *activityLogRepository) GetAll(limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByDateRange получает логи за период
func (r *activityLogRepository) GetByDateRange(startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetByUserAndDateRange получает логи пользователя за период
func (r *activityLogRepository) GetByUserAndDateRange(userID uuid.UUID, startDate, endDate time.Time, limit, offset int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// CountByUserID подсчитывает количество записей пользователя
func (r *activityLogRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByAction подсчитывает количество записей по действию
func (r *activityLogRepository) CountByAction(action string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActivityLog{}).Where("action = ?", action).Count(&count).Error
	return count, err
}

// DeleteOldLogs удаляет старые записи (старше указанной даты)
func (r *activityLogRepository) DeleteOldLogs(beforeDate time.Time) error {
	return r.db.Where("created_at < ?", beforeDate).Delete(&models.ActivityLog{}).Error
}
