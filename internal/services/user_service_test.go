package services

import (
	"testing"
	"time"

	"remnawave-tg-shop/internal/models"
	"remnawave-tg-shop/internal/services/remnawave"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository мок для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	args := m.Called(telegramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByReferralCode(code string) (*models.User, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(limit, offset int) ([]models.User, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetReferrals(userID uuid.UUID) ([]models.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Search(query string, limit int) ([]models.User, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]models.User), args.Error(1)
}

// MockLogger мок для Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args...)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args...)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args...)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args...)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args...)
}

func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.Called(format, args...)
}

func (m *MockLogger) WithField(key string, value interface{}) Logger {
	args := m.Called(key, value)
	return args.Get(0).(Logger)
}

func (m *MockLogger) WithFields(fields map[string]interface{}) Logger {
	args := m.Called(fields)
	return args.Get(0).(Logger)
}

func TestUserService_CreateOrGetUser_NewUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	mockRemnawaveClient := &remnawave.Client{} // В реальном тесте можно создать мок

	service := NewUserService(mockRepo, mockRemnawaveClient, mockLogger)

	telegramID := int64(123456789)
	username := "testuser"
	firstName := "Test"
	lastName := "User"
	languageCode := "ru"

	// Мокаем, что пользователь не найден
	mockRepo.On("GetByTelegramID", telegramID).Return(nil, nil)

	// Мокаем создание пользователя
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Мокаем логирование
	mockLogger.On("Info", "New user created", "telegram_id", telegramID, "username", username).Return()

	// Act
	user, err := service.CreateOrGetUser(telegramID, username, firstName, lastName, languageCode)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, telegramID, user.TelegramID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, firstName, user.FirstName)
	assert.Equal(t, lastName, user.LastName)
	assert.Equal(t, languageCode, user.LanguageCode)
	assert.False(t, user.IsBlocked)
	assert.False(t, user.IsAdmin)
	assert.Equal(t, 0.0, user.Balance)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUserService_CreateOrGetUser_ExistingUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	mockRemnawaveClient := &remnawave.Client{}

	service := NewUserService(mockRepo, mockRemnawaveClient, mockLogger)

	telegramID := int64(123456789)
	username := "testuser"
	firstName := "Test"
	lastName := "User"
	languageCode := "ru"

	existingUser := &models.User{
		ID:           uuid.New(),
		TelegramID:   telegramID,
		Username:     "oldusername",
		FirstName:    "OldFirst",
		LastName:     "OldLast",
		LanguageCode: "en",
		IsBlocked:    false,
		IsAdmin:      false,
		Balance:      100.0,
		CreatedAt:    time.Now().Add(-24 * time.Hour),
		UpdatedAt:    time.Now().Add(-24 * time.Hour),
	}

	// Мокаем, что пользователь найден
	mockRepo.On("GetByTelegramID", telegramID).Return(existingUser, nil)

	// Мокаем обновление пользователя
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	user, err := service.CreateOrGetUser(telegramID, username, firstName, lastName, languageCode)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, telegramID, user.TelegramID)
	assert.Equal(t, username, user.Username) // Обновлено
	assert.Equal(t, firstName, user.FirstName) // Обновлено
	assert.Equal(t, lastName, user.LastName) // Обновлено
	assert.Equal(t, languageCode, user.LanguageCode) // Обновлено
	assert.Equal(t, 100.0, user.Balance) // Сохранено

	mockRepo.AssertExpectations(t)
}

func TestUserService_AddBalance(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	mockRemnawaveClient := &remnawave.Client{}

	service := NewUserService(mockRepo, mockRemnawaveClient, mockLogger)

	userID := uuid.New()
	amount := 50.0

	user := &models.User{
		ID:        userID,
		Balance:   100.0,
		UpdatedAt: time.Now(),
	}

	// Мокаем получение пользователя
	mockRepo.On("GetByID", userID).Return(user, nil)

	// Мокаем обновление пользователя
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)

	// Мокаем логирование
	mockLogger.On("Info", "Balance added", "user_id", userID, "amount", amount).Return()

	// Act
	err := service.AddBalance(userID, amount)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 150.0, user.Balance) // 100 + 50

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUserService_SubtractBalance_InsufficientBalance(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	mockRemnawaveClient := &remnawave.Client{}

	service := NewUserService(mockRepo, mockRemnawaveClient, mockLogger)

	userID := uuid.New()
	amount := 150.0

	user := &models.User{
		ID:        userID,
		Balance:   100.0, // Недостаточно средств
		UpdatedAt: time.Now(),
	}

	// Мокаем получение пользователя
	mockRepo.On("GetByID", userID).Return(user, nil)

	// Act
	err := service.SubtractBalance(userID, amount)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient balance", err.Error())
	assert.Equal(t, 100.0, user.Balance) // Баланс не изменился

	mockRepo.AssertExpectations(t)
}

func TestUserService_IsAdmin(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	mockRemnawaveClient := &remnawave.Client{}

	service := NewUserService(mockRepo, mockRemnawaveClient, mockLogger)

	telegramID := int64(123456789)

	// Test case 1: User is admin
	adminUser := &models.User{
		TelegramID: telegramID,
		IsAdmin:    true,
	}

	mockRepo.On("GetByTelegramID", telegramID).Return(adminUser, nil)

	// Act
	isAdmin := service.IsAdmin(telegramID)

	// Assert
	assert.True(t, isAdmin)

	// Test case 2: User is not admin
	regularUser := &models.User{
		TelegramID: telegramID,
		IsAdmin:    false,
	}

	mockRepo.On("GetByTelegramID", telegramID).Return(regularUser, nil)

	// Act
	isAdmin = service.IsAdmin(telegramID)

	// Assert
	assert.False(t, isAdmin)

	mockRepo.AssertExpectations(t)
}
