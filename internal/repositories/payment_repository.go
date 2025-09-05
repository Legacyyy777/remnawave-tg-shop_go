package repositories

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// paymentRepository реализация PaymentRepository
type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository создает новый репозиторий платежей
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create создает новый платеж
func (r *paymentRepository) Create(payment *models.Payment) error {
	if err := r.db.Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

// GetByID получает платеж по ID
func (r *paymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}
	return &payment, nil
}

// GetByUserID получает платежи пользователя
func (r *paymentRepository) GetByUserID(userID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by user ID: %w", err)
	}
	return payments, nil
}

// GetByExternalID получает платеж по внешнему ID
func (r *paymentRepository) GetByExternalID(externalID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "external_id = ?", externalID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by external ID: %w", err)
	}
	return &payment, nil
}

// Update обновляет платеж
func (r *paymentRepository) Update(payment *models.Payment) error {
	if err := r.db.Save(payment).Error; err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return nil
}

// Delete удаляет платеж
func (r *paymentRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.Payment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete payment: %w", err)
	}
	return nil
}

// List получает список платежей с пагинацией
func (r *paymentRepository) List(limit, offset int) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Preload("User").Limit(limit).Offset(offset).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}
	return payments, nil
}

// GetByStatus получает платежи по статусу
func (r *paymentRepository) GetByStatus(status string) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("status = ?", status).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by status: %w", err)
	}
	return payments, nil
}

// GetByMethod получает платежи по способу оплаты
func (r *paymentRepository) GetByMethod(method string) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("payment_method = ?", method).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by method: %w", err)
	}
	return payments, nil
}

// GetCompletedPayments получает завершенные платежи
func (r *paymentRepository) GetCompletedPayments() ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("status = ?", "completed").Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get completed payments: %w", err)
	}
	return payments, nil
}

// GetByDateRange получает платежи за определенный период
func (r *paymentRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("failed to get payments by date range: %w", err)
	}
	return payments, nil
}
