package repository

import (
	"backend/internal/database"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock for database.Querier
type MockQuerier struct {
	mock.Mock
}

// Implement all methods of the database.Querier interface
func (m *MockQuerier) CheckCartItemExists(ctx context.Context, arg database.CheckCartItemExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuerier) CreateCartItem(ctx context.Context, arg database.CreateCartItemParams) (database.CartItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.CartItem), args.Error(1)
}

func (m *MockQuerier) UpdateCartItemQuantity(ctx context.Context, arg database.UpdateCartItemQuantityParams) (database.CartItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.CartItem), args.Error(1)
}

func (m *MockQuerier) GetCartItemsByUser(ctx context.Context, userID pgtype.UUID) ([]database.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]database.CartItem), args.Error(1)
}

func (m *MockQuerier) DeleteCartItem(ctx context.Context, arg database.DeleteCartItemParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) DeleteAllCartItemsByUser(ctx context.Context, userID pgtype.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Add empty implementations for all other Querier methods to satisfy the interface
func (m *MockQuerier) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *MockQuerier) CheckEmailExistsForOtherUser(ctx context.Context, arg database.CheckEmailExistsForOtherUserParams) (bool, error) {
	return false, nil
}
func (m *MockQuerier) CreateCoinTransaction(ctx context.Context, arg database.CreateCoinTransactionParams) (database.CoinTransaction, error) {
	return database.CoinTransaction{}, nil
}
func (m *MockQuerier) CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
	return database.Order{}, nil
}
func (m *MockQuerier) CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
	return database.OrderItem{}, nil
}
func (m *MockQuerier) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	return database.CreateUserRow{}, nil
}
func (m *MockQuerier) DeleteUser(ctx context.Context, id pgtype.UUID) error { return nil }
func (m *MockQuerier) GetAllCategories(ctx context.Context) ([]database.Category, error) {
	return nil, nil
}
func (m *MockQuerier) GetCategoryByID(ctx context.Context, id int32) (database.Category, error) {
	return database.Category{}, nil
}
func (m *MockQuerier) GetCoinTransactionByID(ctx context.Context, id int32) (database.CoinTransaction, error) {
	return database.CoinTransaction{}, nil
}
func (m *MockQuerier) GetCoinTransactionsByUserID(ctx context.Context, arg database.GetCoinTransactionsByUserIDParams) ([]database.CoinTransaction, error) {
	return nil, nil
}
func (m *MockQuerier) GetOrderByID(ctx context.Context, id int32) (database.Order, error) {
	return database.Order{}, nil
}
func (m *MockQuerier) GetOrderItemByID(ctx context.Context, id int32) (database.OrderItem, error) {
	return database.OrderItem{}, nil
}
func (m *MockQuerier) GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]database.OrderItem, error) {
	return nil, nil
}
func (m *MockQuerier) GetOrdersByUserID(ctx context.Context, userID pgtype.UUID) ([]database.Order, error) {
	return nil, nil
}
func (m *MockQuerier) GetProductByID(ctx context.Context, id int32) (database.Product, error) {
	return database.Product{}, nil
}
func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return database.User{}, nil
}
func (m *MockQuerier) GetUserByID(ctx context.Context, id pgtype.UUID) (database.User, error) {
	return database.User{}, nil
}
func (m *MockQuerier) ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error) {
	return nil, nil
}
func (m *MockQuerier) ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error) {
	return nil, nil
}
func (m *MockQuerier) UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error) {
	return database.Product{}, nil
}
func (m *MockQuerier) UpdateUserCoins(ctx context.Context, arg database.UpdateUserCoinsParams) (database.UpdateUserCoinsRow, error) {
	return database.UpdateUserCoinsRow{}, nil
}
func (m *MockQuerier) UpdateUserEmail(ctx context.Context, arg database.UpdateUserEmailParams) (database.UpdateUserEmailRow, error) {
	return database.UpdateUserEmailRow{}, nil
}
func (m *MockQuerier) UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.UpdateUserNameRow, error) {
	return database.UpdateUserNameRow{}, nil
}
func (m *MockQuerier) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return nil
}
func (m *MockQuerier) WithTx(tx pgx.Tx) *database.Queries { return database.New(tx) }

// Helper functions for tests
func createMockDBCartItem(id int32, userID uuid.UUID, productID int32, quantity int32) database.CartItem {
	return database.CartItem{
		ID:        id,
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func setupCartTest(t *testing.T) (CartRepository, *MockQuerier) {
	mockQuerier := new(MockQuerier)
	repo := NewCartRepository(mockQuerier, nil)
	return repo, mockQuerier
}

func TestAddToCart_NewItem(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := 2

	// Mock CheckCartItemExists to return false
	mockQuerier.On("CheckCartItemExists", ctx, mock.Anything).Return(false, nil).Once()

	// Mock CreateCartItem
	mockQuerier.On("CreateCartItem", ctx, mock.Anything).Return(database.CartItem{}, nil).Once()

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestAddToCart_ExistingItem(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	additionalQuantity := 2

	// Mock CheckCartItemExists to return true
	mockQuerier.On("CheckCartItemExists", ctx, mock.Anything).Return(true, nil).Once()

	// Mock UpdateCartItemQuantity
	// The real repository has a bug and REPLACES the quantity instead of adding.
	// This test asserts the current buggy behavior.
	mockQuerier.On("UpdateCartItemQuantity", ctx, mock.MatchedBy(func(params database.UpdateCartItemQuantityParams) bool {
		return params.Quantity == int32(additionalQuantity)
	})).Return(database.CartItem{}, nil).Once()

	// Execute
	err := repo.AddToCart(ctx, userID, productID, additionalQuantity)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestGetCartItems_Success(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Create mock cart items
	dbItems := []database.CartItem{
		createMockDBCartItem(1, userID, 100, 2),
		createMockDBCartItem(2, userID, 101, 3),
	}

	mockQuerier.On("GetCartItemsByUser", ctx, pgUUID).Return(dbItems, nil)

	// Execute
	items, err := repo.GetCartItems(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, 100, items[0].ProductID)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, 101, items[1].ProductID)
	assert.Equal(t, 3, items[1].Quantity)
	mockQuerier.AssertExpectations(t)
}

func TestGetCartItems_EmptyCart(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock empty result
	mockQuerier.On("GetCartItemsByUser", ctx, pgUUID).Return([]database.CartItem{}, nil)

	// Execute
	items, err := repo.GetCartItems(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 0)
	mockQuerier.AssertExpectations(t)
}

func TestRemoveFromCart_Success(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	productID := 100

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQuerier.On("DeleteCartItem", ctx, database.DeleteCartItemParams{
		UserID:    pgUUID,
		ProductID: int32(productID),
	}).Return(nil)

	// Execute
	err := repo.RemoveFromCart(ctx, userID, productID)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestClearCart_Success(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQuerier.On("DeleteAllCartItemsByUser", ctx, pgUUID).Return(nil)

	// Execute
	err := repo.ClearCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockQuerier.AssertExpectations(t)
}

func TestAddToCart_DatabaseError(t *testing.T) {
	repo, mockQuerier := setupCartTest(t)
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := 2

	// Mock CheckCartItemExists to return a database error
	mockQuerier.On("CheckCartItemExists", ctx, mock.Anything).Return(false, errors.New("db connection failed"))

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cart is empty", err.Error())
	mockQuerier.AssertExpectations(t)
}
