package models

import (
	"time"

	"github.com/google/uuid"
)

// ActivityLog представляет лог активности пользователя
type ActivityLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Action    string    `gorm:"size:100;not null" json:"action"` // command, message, callback, payment, etc.
	Data      string    `gorm:"type:text" json:"data"`           // JSON с данными действия
	IPAddress string    `gorm:"size:45" json:"ip_address"`       // IPv4 или IPv6
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GetActionText возвращает текстовое описание действия
func (al *ActivityLog) GetActionText() string {
	switch al.Action {
	case "command":
		return "Команда"
	case "message":
		return "Сообщение"
	case "callback":
		return "Callback"
	case "payment":
		return "Платеж"
	case "subscription":
		return "Подписка"
	case "promo_code":
		return "Промокод"
	case "referral":
		return "Реферал"
	default:
		return al.Action
	}
}
