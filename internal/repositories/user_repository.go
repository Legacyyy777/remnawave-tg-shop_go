package repositories

import (
	"fmt"

	"remnawave-tg-shop/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository реализация UserRepository
type userRepository struct {
	db *gorm.DB
}

// Убеждаемся, что userRepository реализует UserRepository
var _ UserRepository = (*userRepository)(nil)

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByTelegramID получает пользователя по Telegram ID
func (r *userRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "telegram_id = ?", telegramID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by Telegram ID: %w", err)
	}
	return &user, nil
}

// GetByReferralCode получает пользователя по реферальному коду
func (r *userRepository) GetByReferralCode(code string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "referral_code = ?", code).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by referral code: %w", err)
	}
	return &user, nil
}

// Update обновляет пользователя
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete удаляет пользователя
func (r *userRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List получает список пользователей с пагинацией
func (r *userRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// GetReferrals получает рефералов пользователя
func (r *userRepository) GetReferrals(userID uuid.UUID) ([]models.User, error) {
	var users []models.User
	if err := r.db.Where("referred_by = ?", userID).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get referrals: %w", err)
	}
	return users, nil
}

// GetByUsername получает пользователя по username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

// Search ищет пользователей по запросу
func (r *userRepository) Search(query string, limit int) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + query + "%"

	if err := r.db.Where(
		"username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
		searchPattern, searchPattern, searchPattern,
	).Limit(limit).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// GetAll получает всех пользователей
func (r *userRepository) GetAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	query := r.db.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}
