package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription представляет подписку пользователя
type Subscription struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ServerID     int       `gorm:"not null" json:"server_id"`
	ServerName   string    `gorm:"size:255" json:"server_name"`
	PlanID       int       `gorm:"not null" json:"plan_id"`
	PlanName     string    `gorm:"size:255" json:"plan_name"`
	Status       string    `gorm:"size:50;default:'active'" json:"status"` // active, expired, cancelled, suspended
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// IsActive проверяет, активна ли подписка
func (s *Subscription) IsActive() bool {
	return s.Status == "active" && s.ExpiresAt.After(time.Now())
}

// IsExpired проверяет, истекла ли подписка
func (s *Subscription) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

// GetDaysLeft возвращает количество дней до истечения подписки
func (s *Subscription) GetDaysLeft() int {
	if s.IsExpired() {
		return 0
	}
	days := int(time.Until(s.ExpiresAt).Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}

// GetStatusText возвращает текстовое описание статуса
func (s *Subscription) GetStatusText() string {
	switch s.Status {
	case "active":
		if s.IsExpired() {
			return "Истекла"
		}
		return "Активна"
	case "expired":
		return "Истекла"
	case "cancelled":
		return "Отменена"
	case "suspended":
		return "Приостановлена"
	default:
		return "Неизвестно"
	}
}
