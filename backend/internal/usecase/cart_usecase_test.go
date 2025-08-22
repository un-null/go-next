package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"backend/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mock Repository ----

// MockCartRepository is a mock implementation of repository.CartRepository
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) AddToCart(ctx context.Context, userID uuid.UUID, productID int, quantity int) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}

func (m *MockCartRepository) GetCartItems(ctx context.Context, userID uuid.UUID) ([]entity.CartItem, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.CartItem), args.Error(1)
}

func (m *MockCartRepository) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// ---- Helper Functions ----

func sampleEntityCartItem(userID uuid.UUID, productID int, quantity int) entity.CartItem {
	now := time.Now()
	return entity.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ---- Tests ----

func TestCartUseCase_AddToCart_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 100
	quantity := 2

	mockRepo.On("AddToCart", ctx, userID, productID, quantity).Return(nil)

	err := useCase.AddToCart(ctx, userID, productID, quantity)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_AddToCart_Error(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 100
	quantity := 2

	expectedErr := errors.New("repository error")
	mockRepo.On("AddToCart", ctx, userID, productID, quantity).Return(expectedErr)

	err := useCase.AddToCart(ctx, userID, productID, quantity)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_GetCartItems_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	expectedItems := []entity.CartItem{
		sampleEntityCartItem(userID, 100, 2),
		sampleEntityCartItem(userID, 101, 3),
	}

	mockRepo.On("GetCartItems", ctx, userID).Return(expectedItems, nil)

	items, err := useCase.GetCartItems(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, expectedItems, items)
	assert.Equal(t, userID, items[0].UserID)
	assert.Equal(t, 100, items[0].ProductID)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, userID, items[1].UserID)
	assert.Equal(t, 101, items[1].ProductID)
	assert.Equal(t, 3, items[1].Quantity)

	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_GetCartItems_Empty(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetCartItems", ctx, userID).Return([]entity.CartItem{}, nil)

	items, err := useCase.GetCartItems(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, items, 0)
	assert.NotNil(t, items) // Should return empty slice, not nil

	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_GetCartItems_Error(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	expectedErr := errors.New("repository error")
	mockRepo.On("GetCartItems", ctx, userID).Return(nil, expectedErr)

	items, err := useCase.GetCartItems(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_RemoveFromCart_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 100

	mockRepo.On("RemoveFromCart", ctx, userID, productID).Return(nil)

	err := useCase.RemoveFromCart(ctx, userID, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_RemoveFromCart_Error(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 999

	expectedErr := errors.New("item not found")
	mockRepo.On("RemoveFromCart", ctx, userID, productID).Return(expectedErr)

	err := useCase.RemoveFromCart(ctx, userID, productID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_ClearCart_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("ClearCart", ctx, userID).Return(nil)

	err := useCase.ClearCart(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_ClearCart_Error(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	expectedErr := errors.New("database error")
	mockRepo.On("ClearCart", ctx, userID).Return(expectedErr)

	err := useCase.ClearCart(ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}

// ---- Edge Cases and Additional Tests ----

func TestCartUseCase_AddToCart_ZeroQuantity(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 100
	quantity := 0

	mockRepo.On("AddToCart", ctx, userID, productID, quantity).Return(nil)

	err := useCase.AddToCart(ctx, userID, productID, quantity)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_AddToCart_NegativeQuantity(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()
	productID := 100
	quantity := -1

	mockRepo.On("AddToCart", ctx, userID, productID, quantity).Return(nil)

	err := useCase.AddToCart(ctx, userID, productID, quantity)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_Operations_WithEmptyUUID(t *testing.T) {
	mockRepo := new(MockCartRepository)
	useCase := NewCartUseCase(mockRepo)

	ctx := context.Background()
	var emptyUUID uuid.UUID // Zero value UUID
	productID := 100
	quantity := 2

	mockRepo.On("AddToCart", ctx, emptyUUID, productID, quantity).Return(nil)
	mockRepo.On("GetCartItems", ctx, emptyUUID).Return([]entity.CartItem{}, nil)
	mockRepo.On("RemoveFromCart", ctx, emptyUUID, productID).Return(nil)
	mockRepo.On("ClearCart", ctx, emptyUUID).Return(nil)

	// Test all operations with empty UUID
	err := useCase.AddToCart(ctx, emptyUUID, productID, quantity)
	assert.NoError(t, err)

	items, err := useCase.GetCartItems(ctx, emptyUUID)
	assert.NoError(t, err)
	assert.Len(t, items, 0)

	err = useCase.RemoveFromCart(ctx, emptyUUID, productID)
	assert.NoError(t, err)

	err = useCase.ClearCart(ctx, emptyUUID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCartUseCase_Constructor(t *testing.T) {
	mockRepo := new(MockCartRepository)

	useCase := NewCartUseCase(mockRepo)

	assert.NotNil(t, useCase)
	assert.Equal(t, mockRepo, useCase.repo)
}

func TestCartUseCase_Constructor_NilRepository(t *testing.T) {
	useCase := NewCartUseCase(nil)

	assert.NotNil(t, useCase)
	assert.Nil(t, useCase.repo)
}
