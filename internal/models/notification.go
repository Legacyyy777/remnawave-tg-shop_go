package models

import (
	"time"

	"github.com/google/uuid"
)

// Notification представляет уведомление
type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"` // nil для глобальных уведомлений
	Type      string     `gorm:"size:50;not null" json:"type"`             // subscription_expiring, payment_success, referral_bonus, etc.
	Title     string     `gorm:"size:255;not null" json:"title"`
	Message   string     `gorm:"type:text;not null" json:"message"`
	IsRead    bool       `gorm:"default:false" json:"is_read"`
	IsSent    bool       `gorm:"default:false" json:"is_sent"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GetTypeText возвращает текстовое описание типа уведомления
func (n *Notification) GetTypeText() string {
	switch n.Type {
	case "subscription_expiring":
		return "Истекает подписка"
	case "payment_success":
		return "Успешный платеж"
	case "referral_bonus":
		return "Реферальный бонус"
	case "promo_code_applied":
		return "Применен промокод"
	case "admin_message":
		return "Сообщение от администратора"
	default:
		return n.Type
	}
}
