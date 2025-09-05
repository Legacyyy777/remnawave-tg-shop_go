package repositories

import (
	"remnawave-tg-shop/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type promoCodeRepository struct {
	db *gorm.DB
}

// Убеждаемся, что promoCodeRepository реализует PromoCodeRepository
var _ PromoCodeRepository = (*promoCodeRepository)(nil)

func NewPromoCodeRepository(db *gorm.DB) PromoCodeRepository {
	return &promoCodeRepository{db: db}
}

// Create создает новый промокод
func (r *promoCodeRepository) Create(promoCode *models.PromoCode) error {
	return r.db.Create(promoCode).Error
}

// GetByCode получает промокод по коду
func (r *promoCodeRepository) GetByCode(code string) (*models.PromoCode, error) {
	var promoCode models.PromoCode
	err := r.db.Where("code = ?", code).First(&promoCode).Error
	if err != nil {
		return nil, err
	}
	return &promoCode, nil
}

// GetByID получает промокод по ID
func (r *promoCodeRepository) GetByID(id uuid.UUID) (*models.PromoCode, error) {
	var promoCode models.PromoCode
	err := r.db.Where("id = ?", id).First(&promoCode).Error
	if err != nil {
		return nil, err
	}
	return &promoCode, nil
}

// GetAll получает все промокоды с пагинацией
func (r *promoCodeRepository) GetAll(limit, offset int) ([]models.PromoCode, error) {
	var promoCodes []models.PromoCode
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&promoCodes).Error
	return promoCodes, err
}

// Update обновляет промокод
func (r *promoCodeRepository) Update(promoCode *models.PromoCode) error {
	return r.db.Save(promoCode).Error
}

// Delete удаляет промокод
func (r *promoCodeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PromoCode{}, id).Error
}

// IncrementUsage увеличивает счетчик использований
func (r *promoCodeRepository) IncrementUsage(id uuid.UUID) error {
	return r.db.Model(&models.PromoCode{}).Where("id = ?", id).Update("used_count", gorm.Expr("used_count + 1")).Error
}

// GetValidPromoCodes получает все действительные промокоды
func (r *promoCodeRepository) GetValidPromoCodes() ([]models.PromoCode, error) {
	var promoCodes []models.PromoCode
	now := time.Now()
	err := r.db.Where("is_active = ? AND valid_from <= ? AND (valid_until IS NULL OR valid_until > ?) AND (max_uses = 0 OR used_count < max_uses)",
		true, now, now).Find(&promoCodes).Error
	return promoCodes, err
}

// CreateUsage создает запись об использовании промокода
func (r *promoCodeRepository) CreateUsage(usage *models.PromoCodeUsage) error {
	return r.db.Create(usage).Error
}

// GetUsageByUserAndPromoCode получает использование промокода пользователем
func (r *promoCodeRepository) GetUsageByUserAndPromoCode(userID, promoCodeID uuid.UUID) (*models.PromoCodeUsage, error) {
	var usage models.PromoCodeUsage
	err := r.db.Where("user_id = ? AND promo_code_id = ?", userID, promoCodeID).First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// GetUsageCountByPromoCode получает количество использований промокода
func (r *promoCodeRepository) GetUsageCountByPromoCode(promoCodeID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.PromoCodeUsage{}).Where("promo_code_id = ?", promoCodeID).Count(&count).Error
	return count, err
}
