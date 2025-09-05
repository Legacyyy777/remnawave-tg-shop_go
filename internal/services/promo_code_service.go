package services

import (
	"fmt"
	"math/rand"
	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PromoCodeService struct {
	repo   repositories.PromoCodeRepository
	config *config.Config
}

func NewPromoCodeService(repo repositories.PromoCodeRepository, config *config.Config) *PromoCodeService {
	return &PromoCodeService{
		repo:   repo,
		config: config,
	}
}

// CreatePromoCode создает новый промокод
func (s *PromoCodeService) CreatePromoCode(code, promoType string, value float64, maxUses int, validFrom, validUntil *time.Time, description string, createdBy uuid.UUID) (*models.PromoCode, error) {
	// Валидация
	if err := s.validatePromoCode(code, promoType, value); err != nil {
		return nil, err
	}

	// Проверяем, не существует ли уже такой код
	existing, _ := s.repo.GetByCode(code)
	if existing != nil {
		return nil, fmt.Errorf("промокод с таким кодом уже существует")
	}

	promoCode := &models.PromoCode{
		Code:        strings.ToUpper(code),
		Type:        promoType,
		Value:       value,
		MaxUses:     maxUses,
		IsActive:    true,
		Description: description,
		CreatedBy:   createdBy,
	}

	if validFrom != nil {
		promoCode.ValidFrom = *validFrom
	} else {
		promoCode.ValidFrom = time.Now()
	}

	if validUntil != nil {
		promoCode.ValidUntil = validUntil
	}

	return promoCode, s.repo.Create(promoCode)
}

// GeneratePromoCode генерирует случайный промокод
func (s *PromoCodeService) GeneratePromoCode(promoType string, value float64, maxUses int, validFrom, validUntil *time.Time, description string, createdBy uuid.UUID) (*models.PromoCode, error) {
	code := s.generateRandomCode()

	// Проверяем уникальность
	for {
		existing, _ := s.repo.GetByCode(code)
		if existing == nil {
			break
		}
		code = s.generateRandomCode()
	}

	return s.CreatePromoCode(code, promoType, value, maxUses, validFrom, validUntil, description, createdBy)
}

// ApplyPromoCode применяет промокод к пользователю
func (s *PromoCodeService) ApplyPromoCode(userID uuid.UUID, code string) (*models.PromoCode, error) {
	// Получаем промокод
	promoCode, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, fmt.Errorf("промокод не найден")
	}

	// Проверяем валидность
	if !promoCode.CanBeUsed() {
		return nil, fmt.Errorf("промокод недействителен или истек")
	}

	// Проверяем, не использовал ли уже пользователь этот промокод
	usage, err := s.repo.GetUsageByUserAndPromoCode(userID, promoCode.ID)
	if err == nil && usage != nil {
		return nil, fmt.Errorf("вы уже использовали этот промокод")
	}

	// Создаем запись об использовании
	usage = &models.PromoCodeUsage{
		PromoCodeID: promoCode.ID,
		UserID:      userID,
		UsedAt:      time.Now(),
	}

	if err := s.repo.CreateUsage(usage); err != nil {
		return nil, fmt.Errorf("ошибка при применении промокода: %v", err)
	}

	// Увеличиваем счетчик использований
	if err := s.repo.IncrementUsage(promoCode.ID); err != nil {
		return nil, fmt.Errorf("ошибка при обновлении счетчика использований: %v", err)
	}

	return promoCode, nil
}

// GetPromoCode получает промокод по коду
func (s *PromoCodeService) GetPromoCode(code string) (*models.PromoCode, error) {
	return s.repo.GetByCode(code)
}

// GetPromoCodeByID получает промокод по ID
func (s *PromoCodeService) GetPromoCodeByID(id uuid.UUID) (*models.PromoCode, error) {
	return s.repo.GetByID(id)
}

// GetAllPromoCodes получает все промокоды
func (s *PromoCodeService) GetAllPromoCodes(limit, offset int) ([]models.PromoCode, error) {
	return s.repo.GetAll(limit, offset)
}

// UpdatePromoCode обновляет промокод
func (s *PromoCodeService) UpdatePromoCode(promoCode *models.PromoCode) error {
	return s.repo.Update(promoCode)
}

// DeletePromoCode удаляет промокод
func (s *PromoCodeService) DeletePromoCode(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// GetValidPromoCodes получает все действительные промокоды
func (s *PromoCodeService) GetValidPromoCodes() ([]models.PromoCode, error) {
	return s.repo.GetValidPromoCodes()
}

// validatePromoCode валидирует данные промокода
func (s *PromoCodeService) validatePromoCode(code, promoType string, value float64) error {
	// Проверяем длину кода
	if len(code) < s.config.PromoCodes.MinCodeLength || len(code) > s.config.PromoCodes.MaxCodeLength {
		return fmt.Errorf("длина кода должна быть от %d до %d символов",
			s.config.PromoCodes.MinCodeLength, s.config.PromoCodes.MaxCodeLength)
	}

	// Проверяем тип промокода
	validTypes := []string{"bonus_days", "discount_percent", "discount_amount"}
	if !contains(validTypes, promoType) {
		return fmt.Errorf("недопустимый тип промокода")
	}

	// Проверяем значение
	if value <= 0 {
		return fmt.Errorf("значение промокода должно быть больше 0")
	}

	if promoType == "discount_percent" && value > 100 {
		return fmt.Errorf("скидка в процентах не может быть больше 100")
	}

	return nil
}

// generateRandomCode генерирует случайный код
func (s *PromoCodeService) generateRandomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := s.config.PromoCodes.MinCodeLength + rand.Intn(s.config.PromoCodes.MaxCodeLength-s.config.PromoCodes.MinCodeLength+1)

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

// contains проверяет, содержится ли строка в слайсе
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
