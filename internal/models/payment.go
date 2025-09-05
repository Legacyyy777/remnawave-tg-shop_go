package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Payment представляет платеж пользователя
type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	Currency      string    `gorm:"size:10;default:'RUB'" json:"currency"`
	PaymentMethod string    `gorm:"size:50;not null" json:"payment_method"` // stars, tribute, yookassa
	Status        string    `gorm:"size:50;default:'pending'" json:"status"` // pending, completed, failed, cancelled
	ExternalID    string    `gorm:"size:255;uniqueIndex" json:"external_id"` // ID платежа в внешней системе
	Description   string    `gorm:"size:500" json:"description"`
	Metadata      string    `gorm:"type:text" json:"metadata"` // JSON с дополнительными данными
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// IsCompleted проверяет, завершен ли платеж
func (p *Payment) IsCompleted() bool {
	return p.Status == "completed"
}

// IsPending проверяет, ожидает ли платеж
func (p *Payment) IsPending() bool {
	return p.Status == "pending"
}

// IsFailed проверяет, неудачен ли платеж
func (p *Payment) IsFailed() bool {
	return p.Status == "failed"
}

// GetStatusText возвращает текстовое описание статуса
func (p *Payment) GetStatusText() string {
	switch p.Status {
	case "pending":
		return "Ожидает оплаты"
	case "completed":
		return "Завершен"
	case "failed":
		return "Неудачен"
	case "cancelled":
		return "Отменен"
	default:
		return "Неизвестно"
	}
}

// GetPaymentMethodText возвращает текстовое описание способа оплаты
func (p *Payment) GetPaymentMethodText() string {
	switch p.PaymentMethod {
	case "stars":
		return "Telegram Stars"
	case "tribute":
		return "Tribute"
	case "yookassa":
		return "ЮKassa"
	default:
		return p.PaymentMethod
	}
}
