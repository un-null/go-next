package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCoinTransactionRepository matches your repository interface
type MockCoinTransactionRepository struct {
	mock.Mock
}

func (m *MockCoinTransactionRepository) CreateCoinTransaction(ctx context.Context, params repository.CreateCoinTransactionParams) (*entity.CoinTransaction, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CoinTransaction), args.Error(1)
}

func (m *MockCoinTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*entity.CoinTransaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.CoinTransaction), args.Error(1)
}

func (m *MockCoinTransactionRepository) GetTransactionByID(ctx context.Context, id int32) (*entity.CoinTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CoinTransaction), args.Error(1)
}

func (m *MockCoinTransactionRepository) ChargeUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	args := m.Called(ctx, userID, amount, description, orderID)
	if args.Get(0) == nil || args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.User), args.Get(1).(*entity.CoinTransaction), args.Error(2)
}

func (m *MockCoinTransactionRepository) SpendUserCoins(ctx context.Context, userID uuid.UUID, amount int, description string, orderID *int32) (*entity.User, *entity.CoinTransaction, error) {
	args := m.Called(ctx, userID, amount, description, orderID)
	if args.Get(0) == nil || args.Get(1) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.User), args.Get(1).(*entity.CoinTransaction), args.Error(2)
}

// Helper functions
func createCoinTransactionUser(id uuid.UUID, name, email string, coins int) *entity.User {
	return &entity.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Coins:     coins,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createCoinTransaction(id int32, userID uuid.UUID, txType string, amount, balanceAfter int) *entity.CoinTransaction {
	return &entity.CoinTransaction{
		ID:              id,
		UserID:          userID,
		TransactionType: txType,
		Amount:          amount,
		BalanceAfter:    balanceAfter,
		Description:     "Test transaction",
		CreatedAt:       time.Now(),
	}
}

func setupCoinTransactionUseCase() (CoinTransactionUseCase, *MockCoinTransactionRepository) {
	mockRepo := new(MockCoinTransactionRepository)
	useCase := NewCoinTransactionUseCase(mockRepo)
	return useCase, mockRepo
}

// Tests for ChargeUserCoins
func TestChargeUserCoins_Success(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 100
	description := "Test charge"
	var orderID *int32

	expectedUser := createCoinTransactionUser(userID, "Alice", "alice@example.com", 200)
	expectedTransaction := createCoinTransaction(1, userID, "charge", amount, 200)

	mockRepo.On("ChargeUserCoins", ctx, userID, amount, description, orderID).
		Return(expectedUser, expectedTransaction, nil)

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, orderID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 200, user.Coins)
	assert.Equal(t, "charge", transaction.TransactionType)
	assert.Equal(t, amount, transaction.Amount)

	mockRepo.AssertExpectations(t)
}

func TestChargeUserCoins_WithOrderID(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 50
	description := "Order payment"
	orderID := int32(123)

	expectedUser := createCoinTransactionUser(userID, "Bob", "bob@example.com", 150)
	expectedTransaction := createCoinTransaction(2, userID, "charge", amount, 150)
	expectedTransaction.OrderID = &orderID

	mockRepo.On("ChargeUserCoins", ctx, userID, amount, description, &orderID).
		Return(expectedUser, expectedTransaction, nil)

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, &orderID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, transaction)
	assert.NotNil(t, transaction.OrderID)
	assert.Equal(t, orderID, *transaction.OrderID)

	mockRepo.AssertExpectations(t)
}

func TestChargeUserCoins_ZeroAmount(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 0
	description := "Test charge"

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "amount must be positive", err.Error())
}

func TestChargeUserCoins_NegativeAmount(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := -50
	description := "Test charge"

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "amount must be positive", err.Error())
}

func TestChargeUserCoins_EmptyDescription(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 100
	description := ""

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "description is required", err.Error())
}

func TestChargeUserCoins_RepositoryError(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 100
	description := "Test charge"

	mockRepo.On("ChargeUserCoins", ctx, userID, amount, description, (*int32)(nil)).
		Return(nil, nil, errors.New("user not found"))

	user, transaction, err := uc.ChargeUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "user not found", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for SpendUserCoins
func TestSpendUserCoins_Success(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 50
	description := "Test purchase"
	var orderID *int32

	expectedUser := createCoinTransactionUser(userID, "Alice", "alice@example.com", 50)
	expectedTransaction := createCoinTransaction(3, userID, "purchase", -amount, 50)

	mockRepo.On("SpendUserCoins", ctx, userID, amount, description, orderID).
		Return(expectedUser, expectedTransaction, nil)

	user, transaction, err := uc.SpendUserCoins(ctx, userID, amount, description, orderID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, 50, user.Coins)
	assert.Equal(t, "purchase", transaction.TransactionType)
	assert.Equal(t, -amount, transaction.Amount)

	mockRepo.AssertExpectations(t)
}

func TestSpendUserCoins_ZeroAmount(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 0
	description := "Test spend"

	user, transaction, err := uc.SpendUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "amount must be positive", err.Error())
}

func TestSpendUserCoins_EmptyDescription(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 50
	description := ""

	user, transaction, err := uc.SpendUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Equal(t, "description is required", err.Error())
}

func TestSpendUserCoins_InsufficientCoins(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	amount := 200
	description := "Expensive item"

	mockRepo.On("SpendUserCoins", ctx, userID, amount, description, (*int32)(nil)).
		Return(nil, nil, errors.New("insufficient coins: have 100, need 200"))

	user, transaction, err := uc.SpendUserCoins(ctx, userID, amount, description, nil)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "insufficient coins")

	mockRepo.AssertExpectations(t)
}

// Tests for GetUserTransactions
func TestGetUserTransactions_Success(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	page := int32(1)
	limit := int32(10)
	offset := int32(0)

	expectedTransactions := []*entity.CoinTransaction{
		createCoinTransaction(1, userID, "charge", 100, 100),
		createCoinTransaction(2, userID, "purchase", -50, 50),
	}

	mockRepo.On("GetTransactionsByUserID", ctx, userID, limit, offset).
		Return(expectedTransactions, nil)

	transactions, err := uc.GetUserTransactions(ctx, userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Len(t, transactions, 2)
	assert.Equal(t, "charge", transactions[0].TransactionType)
	assert.Equal(t, "purchase", transactions[1].TransactionType)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_WithPagination(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	page := int32(2)
	limit := int32(5)
	offset := int32(5)

	expectedTransactions := []*entity.CoinTransaction{
		createCoinTransaction(6, userID, "refund", 25, 75),
	}

	mockRepo.On("GetTransactionsByUserID", ctx, userID, limit, offset).
		Return(expectedTransactions, nil)

	transactions, err := uc.GetUserTransactions(ctx, userID, page, limit)

	assert.NoError(t, err)
	assert.Len(t, transactions, 1)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_DefaultValues(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	page := int32(0)
	limit := int32(0)

	expectedTransactions := []*entity.CoinTransaction{}

	mockRepo.On("GetTransactionsByUserID", ctx, userID, int32(20), int32(0)).
		Return(expectedTransactions, nil)

	transactions, err := uc.GetUserTransactions(ctx, userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_LimitTooHigh(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	page := int32(1)
	limit := int32(150)

	expectedTransactions := []*entity.CoinTransaction{}

	mockRepo.On("GetTransactionsByUserID", ctx, userID, int32(20), int32(0)).
		Return(expectedTransactions, nil)

	transactions, err := uc.GetUserTransactions(ctx, userID, page, limit)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)

	mockRepo.AssertExpectations(t)
}

func TestGetUserTransactions_RepositoryError(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	userID := uuid.New()
	page := int32(1)
	limit := int32(10)

	mockRepo.On("GetTransactionsByUserID", ctx, userID, limit, int32(0)).
		Return(nil, errors.New("database error"))

	transactions, err := uc.GetUserTransactions(ctx, userID, page, limit)

	assert.Error(t, err)
	assert.Nil(t, transactions)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}

// Tests for GetTransactionByID
func TestGetTransactionByID_Success(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	transactionID := int32(1)
	userID := uuid.New()
	expectedTransaction := createCoinTransaction(transactionID, userID, "charge", 100, 200)

	mockRepo.On("GetTransactionByID", ctx, transactionID).
		Return(expectedTransaction, nil)

	transaction, err := uc.GetTransactionByID(ctx, transactionID)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, transactionID, transaction.ID)
	assert.Equal(t, "charge", transaction.TransactionType)

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionByID_InvalidID(t *testing.T) {
	uc, _ := setupCoinTransactionUseCase()
	ctx := context.Background()

	// Test zero ID
	transaction, err := uc.GetTransactionByID(ctx, 0)
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "invalid transaction ID", err.Error())

	// Test negative ID
	transaction, err = uc.GetTransactionByID(ctx, -1)
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "invalid transaction ID", err.Error())
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	uc, mockRepo := setupCoinTransactionUseCase()
	ctx := context.Background()

	transactionID := int32(999)

	mockRepo.On("GetTransactionByID", ctx, transactionID).
		Return(nil, errors.New("transaction not found"))

	transaction, err := uc.GetTransactionByID(ctx, transactionID)

	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Equal(t, "transaction not found", err.Error())

	mockRepo.AssertExpectations(t)
}
