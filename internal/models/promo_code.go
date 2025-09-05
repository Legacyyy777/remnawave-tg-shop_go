package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PromoCode представляет промокод
type PromoCode struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Type        string         `gorm:"size:20;default:'bonus_days'" json:"type"` // bonus_days, discount_percent, discount_amount
	Value       float64        `gorm:"not null" json:"value"`                    // количество дней или размер скидки
	MaxUses     int            `gorm:"default:0" json:"max_uses"`                // 0 = без ограничений
	UsedCount   int            `gorm:"default:0" json:"used_count"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	ValidFrom   time.Time      `json:"valid_from"`
	ValidUntil  *time.Time     `json:"valid_until,omitempty"`
	Description string         `gorm:"size:500" json:"description"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"` // ID администратора, создавшего промокод
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Связи
	Usages []PromoCodeUsage `gorm:"foreignKey:PromoCodeID" json:"usages,omitempty"`
}

// IsValid проверяет, действителен ли промокод
func (pc *PromoCode) IsValid() bool {
	now := time.Now()

	// Проверяем активность
	if !pc.IsActive {
		return false
	}

	// Проверяем дату начала действия
	if now.Before(pc.ValidFrom) {
		return false
	}

	// Проверяем дату окончания действия
	if pc.ValidUntil != nil && now.After(*pc.ValidUntil) {
		return false
	}

	// Проверяем лимит использований
	if pc.MaxUses > 0 && pc.UsedCount >= pc.MaxUses {
		return false
	}

	return true
}

// CanBeUsed проверяет, можно ли использовать промокод
func (pc *PromoCode) CanBeUsed() bool {
	return pc.IsValid()
}

// GetTypeText возвращает текстовое описание типа промокода
func (pc *PromoCode) GetTypeText() string {
	switch pc.Type {
	case "bonus_days":
		return "Бонусные дни"
	case "discount_percent":
		return "Скидка в процентах"
	case "discount_amount":
		return "Скидка в рублях"
	default:
		return "Неизвестно"
	}
}

// PromoCodeUsage представляет использование промокода
type PromoCodeUsage struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PromoCodeID uuid.UUID `gorm:"type:uuid;not null;index" json:"promo_code_id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	UsedAt      time.Time `json:"used_at"`

	// Связи
	PromoCode PromoCode `gorm:"foreignKey:PromoCodeID" json:"promo_code,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
