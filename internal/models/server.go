package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Server представляет сервер Remnawave
type Server struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Связи
	Plans []Plan `gorm:"foreignKey:ServerID" json:"plans,omitempty"`
}

// Plan представляет тарифный план
type Plan struct {
	ID          int     `gorm:"primaryKey" json:"id"`
	ServerID    int     `gorm:"not null;index" json:"server_id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Duration    int     `gorm:"not null" json:"duration"` // в днях
	IsActive    bool    `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Связи
	Server Server `gorm:"foreignKey:ServerID" json:"server,omitempty"`
}

// GetPricePerDay возвращает цену за день
func (p *Plan) GetPricePerDay() float64 {
	if p.Duration <= 0 {
		return 0
	}
	return p.Price / float64(p.Duration)
}

// GetFormattedPrice возвращает отформатированную цену
func (p *Plan) GetFormattedPrice() string {
	return formatPrice(p.Price)
}

// formatPrice форматирует цену
func formatPrice(price float64) string {
	return fmt.Sprintf("%.2f ₽", price)
}
