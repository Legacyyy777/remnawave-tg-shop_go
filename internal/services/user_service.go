package services

import (
	"fmt"
	"time"

	"remnawave-tg-shop/internal/config"
	"remnawave-tg-shop/internal/logger"
	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/repositories"
	"remnawave-tg-shop/internal/services/remnawave"

	"github.com/google/uuid"
)


// userService реализация UserService
type userService struct {
	userRepo        repositories.UserRepository
	remnawaveClient *remnawave.Client
	logger          logger.Logger
	config          *config.Config
}

// NewUserService создает новый сервис пользователей
func NewUserService(userRepo repositories.UserRepository, remnawaveClient *remnawave.Client, log logger.Logger, cfg *config.Config) UserService {
	return &userService{
		userRepo:        userRepo,
		remnawaveClient: remnawaveClient,
		logger:          log,
		config:          cfg,
	}
}

// CreateOrGetUser создает или получает пользователя
func (s *userService) CreateOrGetUser(telegramID int64, username, firstName, lastName, languageCode string) (*models.User, error) {
	// Сначала пытаемся найти существующего пользователя
	user, err := s.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Если пользователь найден, обновляем его данные
	if user != nil {
		user.Username = username
		user.FirstName = firstName
		user.LastName = lastName
		user.LanguageCode = languageCode
		user.UpdatedAt = time.Now()

		if err := s.userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		return user, nil
	}

	// Создаем нового пользователя
	user = &models.User{
		TelegramID:   telegramID,
		Username:     username,
		FirstName:    firstName,
		LastName:     lastName,
		LanguageCode: languageCode,
		IsBlocked:    false,
		IsAdmin:      false,
		Balance:      0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("New user created", "telegram_id", telegramID, "username", username)
	return user, nil
}

// GetUser получает пользователя по Telegram ID
func (s *userService) GetUser(telegramID int64) (*models.User, error) {
	user, err := s.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByReferralCode получает пользователя по реферальному коду
func (s *userService) GetUserByReferralCode(code string) (*models.User, error) {
	user, err := s.userRepo.GetByReferralCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by referral code: %w", err)
	}
	return user, nil
}

// UpdateUser обновляет пользователя
func (s *userService) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// BlockUser блокирует пользователя
func (s *userService) BlockUser(telegramID int64) error {
	user, err := s.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsBlocked = true
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	s.logger.Info("User blocked", "telegram_id", telegramID)
	return nil
}

// UnblockUser разблокирует пользователя
func (s *userService) UnblockUser(telegramID int64) error {
	user, err := s.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.IsBlocked = false
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	s.logger.Info("User unblocked", "telegram_id", telegramID)
	return nil
}

// AddBalance добавляет средства на баланс пользователя
func (s *userService) AddBalance(userID uuid.UUID, amount float64) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	user.Balance += amount
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to add balance: %w", err)
	}

	s.logger.Info("Balance added", "user_id", userID, "amount", amount)
	return nil
}

// SubtractBalance списывает средства с баланса пользователя
func (s *userService) SubtractBalance(userID uuid.UUID, amount float64) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	user.Balance -= amount
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to subtract balance: %w", err)
	}

	s.logger.Info("Balance subtracted", "user_id", userID, "amount", amount)
	return nil
}

// GetReferrals получает рефералов пользователя
func (s *userService) GetReferrals(userID uuid.UUID) ([]models.User, error) {
	referrals, err := s.userRepo.GetReferrals(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referrals: %w", err)
	}
	return referrals, nil
}

// SearchUsers ищет пользователей
func (s *userService) SearchUsers(query string, limit int) ([]models.User, error) {
	users, err := s.userRepo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return users, nil
}

// IsAdmin проверяет, является ли пользователь администратором
func (s *userService) IsAdmin(telegramID int64) bool {
	// Сначала проверяем ADMIN_TELEGRAM_ID из конфигурации
	if s.config.Admin.TelegramID != 0 && s.config.Admin.TelegramID == telegramID {
		return true
	}

	// Затем проверяем поле IsAdmin в базе данных
	user, err := s.userRepo.GetByTelegramID(telegramID)
	if err != nil || user == nil {
		return false
	}
	return user.IsAdmin
}
