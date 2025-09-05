package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User представляет пользователя бота
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TelegramID   int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     string    `gorm:"size:255" json:"username"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	LastName     string    `gorm:"size:255" json:"last_name"`
	LanguageCode string    `gorm:"size:10;default:'ru'" json:"language_code"`
	IsBlocked    bool      `gorm:"default:false" json:"is_blocked"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	Balance      float64   `gorm:"default:0" json:"balance"`
	ReferralCode string    `gorm:"size:20;uniqueIndex" json:"referral_code"`
	ReferredBy   *uuid.UUID `gorm:"type:uuid" json:"referred_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Связи
	Subscriptions []Subscription `gorm:"foreignKey:UserID" json:"subscriptions,omitempty"`
	Payments      []Payment      `gorm:"foreignKey:UserID" json:"payments,omitempty"`
	Referrals     []User         `gorm:"foreignKey:ReferredBy" json:"referrals,omitempty"`
}

// BeforeCreate выполняется перед созданием пользователя
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.ReferralCode == "" {
		u.ReferralCode = generateReferralCode()
	}
	return nil
}

// generateReferralCode генерирует уникальный реферальный код
func generateReferralCode() string {
	return uuid.New().String()[:8]
}

// GetFullName возвращает полное имя пользователя
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.Username != "" {
		return "@" + u.Username
	}
	return "Пользователь"
}

// GetDisplayName возвращает отображаемое имя пользователя
func (u *User) GetDisplayName() string {
	if u.Username != "" {
		return "@" + u.Username
	}
	return u.GetFullName()
}
