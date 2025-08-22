package usecase

import (
	"backend/internal/entity"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error) {
	args := m.Called(ctx, id, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	args := m.Called(ctx, id, newPassword)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error) {
	args := m.Called(ctx, id, coinsDelta)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper functions for tests
func createTestUser(id uuid.UUID, name, email, passwordHash string, coins int) *entity.User {
	return &entity.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Coins:        coins,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func setupUserUseCase() (*UserUseCase, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	useCase := NewUserUseCase(mockRepo)
	return useCase, mockRepo
}

// Tests for GetUserById
func TestGetUserById_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	expectedUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 100)

	mockRepo.On("GetUserById", ctx, userID).Return(expectedUser, nil)

	// Execute
	user, err := uc.GetUserById(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.On("GetUserById", ctx, userID).Return(nil, errors.New("user not found"))

	// Execute
	user, err := uc.GetUserById(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for GetUserByEmail
func TestGetUserByEmail_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "alice@example.com"
	expectedUser := createTestUser(uuid.New(), "Alice", email, "hashedpassword", 100)

	mockRepo.On("GetUserByEmail", ctx, email).Return(expectedUser, nil)

	// Execute
	user, err := uc.GetUserByEmail(ctx, email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "nonexistent@example.com"

	mockRepo.On("GetUserByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// Execute
	user, err := uc.GetUserByEmail(ctx, email)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

// Tests for SignUp
func TestSignUp_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret123",
		Coins:    0,
	}

	expectedUser := createTestUser(uuid.New(), req.Name, req.Email, "hashedpassword", req.Coins)

	mockRepo.On("CreateUser", ctx, req).Return(expectedUser, nil)

	// Execute
	user, err := uc.SignUp(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestSignUp_EmptyName(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "", // Empty name
		Email:    "bob@example.com",
		Password: "secret123",
	}

	// Execute
	user, err := uc.SignUp(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "name is required", err.Error())
}

func TestSignUp_EmptyEmail(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "", // Empty email
		Password: "secret123",
	}

	// Execute
	user, err := uc.SignUp(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email is required", err.Error())
}

func TestSignUp_ShortPassword(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "short", // Less than 8 characters
	}

	// Execute
	user, err := uc.SignUp(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "password must be at least 8 characters", err.Error())
}

func TestSignUp_RepositoryError(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret123",
	}

	mockRepo.On("CreateUser", ctx, req).Return(nil, errors.New("email already exists"))

	// Execute
	user, err := uc.SignUp(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already exists", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for Login
func TestLogin_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "alice@example.com"
	password := "secret123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	expectedUser := createTestUser(uuid.New(), "Alice", email, string(hashedPassword), 100)

	mockRepo.On("GetUserByEmail", ctx, email).Return(expectedUser, nil)

	// Execute
	user, err := uc.Login(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestLogin_EmptyEmail(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	// Execute
	user, err := uc.Login(ctx, "", "password")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email is required", err.Error())
}

func TestLogin_EmptyPassword(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	// Execute
	user, err := uc.Login(ctx, "alice@example.com", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "password is required", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "nonexistent@example.com"

	mockRepo.On("GetUserByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// Execute
	user, err := uc.Login(ctx, email, "password")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "alice@example.com"
	correctPassword := "secret123"
	wrongPassword := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	expectedUser := createTestUser(uuid.New(), "Alice", email, string(hashedPassword), 100)

	mockRepo.On("GetUserByEmail", ctx, email).Return(expectedUser, nil)

	// Execute
	user, err := uc.Login(ctx, email, wrongPassword)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for UpdateUserName
func TestUpdateUserName_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	newName := "Updated Name"
	updatedUser := createTestUser(userID, newName, "alice@example.com", "hashedpassword", 100)

	mockRepo.On("UpdateUserName", ctx, userID, newName).Return(updatedUser, nil)

	// Execute
	user, err := uc.UpdateUserName(ctx, userID, newName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newName, user.Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	// Execute
	user, err := uc.UpdateUserName(ctx, userID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "name cannot be empty", err.Error())
}

// Tests for UpdateUserEmail
func TestUpdateUserEmail_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	newEmail := "newemail@example.com"
	updatedUser := createTestUser(userID, "Alice", newEmail, "hashedpassword", 100)

	mockRepo.On("UpdateUserEmail", ctx, userID, newEmail).Return(updatedUser, nil)

	// Execute
	user, err := uc.UpdateUserEmail(ctx, userID, newEmail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newEmail, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserEmail_EmptyEmail(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	// Execute
	user, err := uc.UpdateUserEmail(ctx, userID, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email cannot be empty", err.Error())
}

// Tests for UpdateUserCoins
func TestUpdateUserCoins_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	currentUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 100)
	coinsDelta := 50
	updatedUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 150)

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)
	mockRepo.On("UpdateUserCoins", ctx, userID, coinsDelta).Return(updatedUser, nil)

	// Execute
	user, err := uc.UpdateUserCoins(ctx, userID, coinsDelta)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 150, user.Coins)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserCoins_InsufficientCoins(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	currentUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 50)
	coinsDelta := -100 // Trying to deduct more than available

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)

	// Execute
	user, err := uc.UpdateUserCoins(ctx, userID, coinsDelta)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "insufficient coins", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for ChangePassword
func TestChangePassword_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	currentPassword := "oldpassword"
	newPassword := "newpassword123"
	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword([]byte(currentPassword), bcrypt.DefaultCost)

	currentUser := createTestUser(userID, "Alice", "alice@example.com", string(hashedCurrentPassword), 100)

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)
	mockRepo.On("UpdateUserPassword", ctx, userID, newPassword).Return(nil)

	// Execute
	err := uc.ChangePassword(ctx, userID, currentPassword, newPassword)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestChangePassword_ShortNewPassword(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	// Execute
	err := uc.ChangePassword(ctx, userID, "oldpass", "short")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "new password must be at least 8 characters", err.Error())
}

func TestChangePassword_IncorrectCurrentPassword(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	correctCurrentPassword := "correctpassword"
	wrongCurrentPassword := "wrongpassword"
	newPassword := "newpassword123"
	hashedCurrentPassword, _ := bcrypt.GenerateFromPassword([]byte(correctCurrentPassword), bcrypt.DefaultCost)

	currentUser := createTestUser(userID, "Alice", "alice@example.com", string(hashedCurrentPassword), 100)

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)

	// Execute
	err := uc.ChangePassword(ctx, userID, wrongCurrentPassword, newPassword)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "current password is incorrect", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for DeleteUser
func TestDeleteUser_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	mockRepo.On("DeleteUser", ctx, userID).Return(nil)

	// Execute
	err := uc.DeleteUser(ctx, userID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Tests for CheckEmailExists
func TestCheckEmailExists_True(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "existing@example.com"

	mockRepo.On("CheckEmailExists", ctx, email).Return(true, nil)

	// Execute
	exists, err := uc.CheckEmailExists(ctx, email)

	// Assert
	assert.NoError(t, err)
	assert.True(t, exists)

	mockRepo.AssertExpectations(t)
}

func TestCheckEmailExists_False(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	email := "nonexistent@example.com"

	mockRepo.On("CheckEmailExists", ctx, email).Return(false, nil)

	// Execute
	exists, err := uc.CheckEmailExists(ctx, email)

	// Assert
	assert.NoError(t, err)
	assert.False(t, exists)

	mockRepo.AssertExpectations(t)
}

// Tests for PurchaseItem
func TestPurchaseItem_Success(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	itemPrice := 50
	currentUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 100)
	updatedUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 50)

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)
	mockRepo.On("UpdateUserCoins", ctx, userID, -itemPrice).Return(updatedUser, nil)

	// Execute
	user, err := uc.PurchaseItem(ctx, userID, itemPrice)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 50, user.Coins)

	mockRepo.AssertExpectations(t)
}

func TestPurchaseItem_InvalidPrice(t *testing.T) {
	uc, _ := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()

	// Execute with zero price
	user, err := uc.PurchaseItem(ctx, userID, 0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid item price", err.Error())

	// Execute with negative price
	user, err = uc.PurchaseItem(ctx, userID, -10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid item price", err.Error())
}

func TestPurchaseItem_InsufficientCoins(t *testing.T) {
	uc, mockRepo := setupUserUseCase()
	ctx := context.Background()

	userID := uuid.New()
	itemPrice := 150
	currentUser := createTestUser(userID, "Alice", "alice@example.com", "hashedpassword", 100)

	mockRepo.On("GetUserById", ctx, userID).Return(currentUser, nil)

	// Execute
	user, err := uc.PurchaseItem(ctx, userID, itemPrice)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "insufficient coins", err.Error())

	mockRepo.AssertExpectations(t)
}
