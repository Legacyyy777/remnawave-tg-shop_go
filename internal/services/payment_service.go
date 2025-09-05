package services

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"

	"github.com/google/uuid"
)

// paymentService реализация PaymentService
type paymentService struct {
	paymentRepo repositories.PaymentRepository
	userService UserService
	logger      logger.Logger
}

// NewPaymentService создает новый сервис платежей
func NewPaymentService(paymentRepo repositories.PaymentRepository, userService UserService, log logger.Logger) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		userService: userService,
		logger:      log,
	}
}

// CreatePayment создает новый платеж
func (s *paymentService) CreatePayment(userID uuid.UUID, amount float64, method, description string) (*models.Payment, error) {
	payment := &models.Payment{
		UserID:        userID,
		Amount:        amount,
		Currency:      "RUB",
		PaymentMethod: method,
		Status:        "pending",
		Description:   description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	s.logger.Info("Payment created", "user_id", userID, "amount", amount, "method", method)
	return payment, nil
}

// GetPayment получает платеж по ID
func (s *paymentService) GetPayment(id uuid.UUID) (*models.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

// UpdatePaymentStatus обновляет статус платежа
func (s *paymentService) UpdatePaymentStatus(id uuid.UUID, status string) error {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return fmt.Errorf("payment not found")
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()

	if status == "completed" {
		now := time.Now()
		payment.CompletedAt = &now
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// Если платеж завершен, добавляем средства на баланс
	if status == "completed" {
		if err := s.userService.AddBalance(payment.UserID, payment.Amount); err != nil {
			s.logger.Error("Failed to add balance after payment completion", "error", err, "user_id", payment.UserID, "amount", payment.Amount)
		}
	}

	s.logger.Info("Payment status updated", "payment_id", id, "status", status)
	return nil
}

// GetUserPayments получает платежи пользователя
func (s *paymentService) GetUserPayments(userID uuid.UUID) ([]models.Payment, error) {
	payments, err := s.paymentRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user payments: %w", err)
	}
	return payments, nil
}

// ProcessStarsPayment обрабатывает платеж через Telegram Stars
func (s *paymentService) ProcessStarsPayment(userID uuid.UUID, amount float64) error {
	payment, err := s.CreatePayment(userID, amount, "stars", "Пополнение баланса через Telegram Stars")
	if err != nil {
		return fmt.Errorf("failed to create stars payment: %w", err)
	}

	// В реальном приложении здесь должна быть интеграция с Telegram Stars API
	// Пока что просто помечаем как завершенный
	if err := s.UpdatePaymentStatus(payment.ID, "completed"); err != nil {
		return fmt.Errorf("failed to complete stars payment: %w", err)
	}

	s.logger.Info("Stars payment processed", "user_id", userID, "amount", amount)
	return nil
}

// ProcessTributePayment обрабатывает платеж через Tribute
func (s *paymentService) ProcessTributePayment(userID uuid.UUID, amount float64) error {
	payment, err := s.CreatePayment(userID, amount, "tribute", "Пополнение баланса через Tribute")
	if err != nil {
		return fmt.Errorf("failed to create tribute payment: %w", err)
	}

	// В реальном приложении здесь должна быть интеграция с Tribute API
	// Пока что просто помечаем как завершенный
	if err := s.UpdatePaymentStatus(payment.ID, "completed"); err != nil {
		return fmt.Errorf("failed to complete tribute payment: %w", err)
	}

	s.logger.Info("Tribute payment processed", "user_id", userID, "amount", amount)
	return nil
}

// ProcessYooKassaPayment обрабатывает платеж через ЮKassa
func (s *paymentService) ProcessYooKassaPayment(userID uuid.UUID, amount float64) error {
	payment, err := s.CreatePayment(userID, amount, "yookassa", "Пополнение баланса через ЮKassa")
	if err != nil {
		return fmt.Errorf("failed to create yookassa payment: %w", err)
	}

	// В реальном приложении здесь должна быть интеграция с ЮKassa API
	// Пока что просто помечаем как завершенный
	if err := s.UpdatePaymentStatus(payment.ID, "completed"); err != nil {
		return fmt.Errorf("failed to complete yookassa payment: %w", err)
	}

	s.logger.Info("YooKassa payment processed", "user_id", userID, "amount", amount)
	return nil
}
