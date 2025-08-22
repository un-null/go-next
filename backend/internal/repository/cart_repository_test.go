package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

// Test repository wrapper that uses the interface
type testCartRepository struct {
	queries mocks.CartQueriesInterface
}

func NewTestCartRepository(queries mocks.CartQueriesInterface) CartRepository {
	return &testCartRepository{
		queries: queries,
	}
}

// Implement the cart repository methods matching your exact interface
func (r *testCartRepository) AddToCart(ctx context.Context, userID uuid.UUID, productID int, quantity int) error {
	// Check if item already exists
	existingItem, err := r.queries.GetCartItem(ctx, pgtype.UUID{Bytes: userID, Valid: true}, int32(productID))

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if err == nil {
		// Item exists, update quantity
		newQuantity := existingItem.Quantity + int32(quantity)
		return r.queries.UpdateCartItemQuantity(ctx, pgtype.UUID{Bytes: userID, Valid: true}, int32(productID), newQuantity)
	}

	// Item doesn't exist, create new
	return r.queries.AddToCart(ctx, pgtype.UUID{Bytes: userID, Valid: true}, int32(productID), int32(quantity))
}

func (r *testCartRepository) GetCartItems(ctx context.Context, userID uuid.UUID) ([]entity.CartItem, error) {
	dbItems, err := r.queries.GetCartItems(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	var items []entity.CartItem
	for _, dbItem := range dbItems {
		item := entity.CartItem{
			ID:        int(dbItem.ID),
			UserID:    dbItem.UserID.Bytes,
			ProductID: int(dbItem.ProductID),
			Quantity:  int(dbItem.Quantity),
			CreatedAt: dbItem.CreatedAt.Time,
			UpdatedAt: dbItem.UpdatedAt.Time,
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *testCartRepository) RemoveFromCart(ctx context.Context, userID uuid.UUID, productID int) error {
	return r.queries.RemoveFromCart(ctx, pgtype.UUID{Bytes: userID, Valid: true}, int32(productID))
}

func (r *testCartRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return r.queries.ClearCart(ctx, pgtype.UUID{Bytes: userID, Valid: true})
}

// Helper functions for tests
func createMockDBCartItem(id, userID uuid.UUID, productID int32, quantity int32) database.CartItem {
	return database.CartItem{
		ID:        int32(time.Now().UnixNano() % 1000000), // Use a simple int32 ID
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func setupCartTestRepository() (CartRepository, *mocks.MockCartQueries) {
	mockQueries := new(mocks.MockCartQueries)
	repo := NewTestCartRepository(mockQueries)
	return repo, mockQueries
}

func TestAddToCart_NewItem(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := 2

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock GetCartItem to return sql.ErrNoRows (item doesn't exist)
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(database.CartItem{}, sql.ErrNoRows)

	// Mock AddToCart
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID), int32(quantity)).
		Return(nil)

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestAddToCart_ExistingItem(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	existingQuantity := int32(3)
	additionalQuantity := 2
	expectedTotalQuantity := existingQuantity + int32(additionalQuantity)

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock GetCartItem to return existing item
	existingItem := createMockDBCartItem(uuid.New(), userID, int32(productID), existingQuantity)
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(existingItem, nil)

	// Mock UpdateCartItemQuantity
	mockQueries.On("UpdateCartItemQuantity", ctx, pgUUID, int32(productID), expectedTotalQuantity).
		Return(nil)

	// Execute
	err := repo.AddToCart(ctx, userID, productID, additionalQuantity)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestGetCartItems_Success(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Create mock cart items
	dbItems := []database.CartItem{
		createMockDBCartItem(uuid.New(), userID, 100, 2),
		createMockDBCartItem(uuid.New(), userID, 101, 3),
	}

	mockQueries.On("GetCartItems", ctx, pgUUID).
		Return(dbItems, nil)

	// Execute
	items, err := repo.GetCartItems(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, 100, items[0].ProductID)
	assert.Equal(t, 2, items[0].Quantity)
	assert.Equal(t, 101, items[1].ProductID)
	assert.Equal(t, 3, items[1].Quantity)

	mockQueries.AssertExpectations(t)
}

func TestGetCartItems_EmptyCart(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock empty result
	mockQueries.On("GetCartItems", ctx, pgUUID).
		Return([]database.CartItem{}, nil)

	// Execute
	items, err := repo.GetCartItems(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, items, 0)

	mockQueries.AssertExpectations(t)
}

func TestRemoveFromCart_Success(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("RemoveFromCart", ctx, pgUUID, int32(productID)).
		Return(nil)

	// Execute
	err := repo.RemoveFromCart(ctx, userID, productID)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestRemoveFromCart_ItemNotFound(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("RemoveFromCart", ctx, pgUUID, int32(productID)).
		Return(sql.ErrNoRows)

	// Execute
	err := repo.RemoveFromCart(ctx, userID, productID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	mockQueries.AssertExpectations(t)
}

func TestClearCart_Success(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("ClearCart", ctx, pgUUID).
		Return(nil)

	// Execute
	err := repo.ClearCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestAddToCart_DatabaseError(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := 2

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock GetCartItem to return database error
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(database.CartItem{}, errors.New("database connection failed"))

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database connection failed", err.Error())
	mockQueries.AssertExpectations(t)
}

func TestGetCartItems_DatabaseError(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("GetCartItems", ctx, pgUUID).
		Return([]database.CartItem{}, errors.New("database timeout"))

	// Execute
	items, err := repo.GetCartItems(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Equal(t, "database timeout", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestClearCart_DatabaseError(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("ClearCart", ctx, pgUUID).
		Return(errors.New("foreign key constraint violation"))

	// Execute
	err := repo.ClearCart(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "foreign key constraint violation", err.Error())
	mockQueries.AssertExpectations(t)
}

// Integration-style test with multiple operations
func TestCartWorkflow_AddAndRemove(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID1 := 100
	productID2 := 101

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Step 1: Add first item
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID1)).
		Return(database.CartItem{}, sql.ErrNoRows)
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID1), int32(2)).
		Return(nil)

	err := repo.AddToCart(ctx, userID, productID1, 2)
	assert.NoError(t, err)

	// Step 2: Add second item
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID2)).
		Return(database.CartItem{}, sql.ErrNoRows)
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID2), int32(3)).
		Return(nil)

	err = repo.AddToCart(ctx, userID, productID2, 3)
	assert.NoError(t, err)

	// Step 3: Get cart items
	cartItem1 := createMockDBCartItem(uuid.New(), userID, int32(productID1), 2)
	cartItem2 := createMockDBCartItem(uuid.New(), userID, int32(productID2), 3)
	dbItems := []database.CartItem{cartItem1, cartItem2}
	mockQueries.On("GetCartItems", ctx, pgUUID).
		Return(dbItems, nil)

	items, err := repo.GetCartItems(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, items, 2)

	// Step 4: Remove first item
	mockQueries.On("RemoveFromCart", ctx, pgUUID, int32(productID1)).
		Return(nil)

	err = repo.RemoveFromCart(ctx, userID, productID1)
	assert.NoError(t, err)

	// Step 5: Clear cart
	mockQueries.On("ClearCart", ctx, pgUUID).
		Return(nil)

	err = repo.ClearCart(ctx, userID)
	assert.NoError(t, err)

	mockQueries.AssertExpectations(t)
}

func TestAddToCart_ZeroQuantity(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := 0

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock GetCartItem to return sql.ErrNoRows (item doesn't exist)
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(database.CartItem{}, sql.ErrNoRows)

	// Mock AddToCart with zero quantity
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID), int32(quantity)).
		Return(nil)

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestAddToCart_ExistingItemTwice(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// First add - item doesn't exist
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(database.CartItem{}, sql.ErrNoRows).Once()
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID), int32(2)).
		Return(nil)

	err := repo.AddToCart(ctx, userID, productID, 2)
	assert.NoError(t, err)

	// Second add - item exists
	cartItem := createMockDBCartItem(uuid.New(), userID, int32(productID), 2)
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(cartItem, nil).Once()
	mockQueries.On("UpdateCartItemQuantity", ctx, pgUUID, int32(productID), int32(5)).
		Return(nil)

	err = repo.AddToCart(ctx, userID, productID, 3)
	assert.NoError(t, err)

	mockQueries.AssertExpectations(t)
}

// Test edge cases
func TestAddToCart_NegativeQuantity(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100
	quantity := -1

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	// Mock GetCartItem to return sql.ErrNoRows (item doesn't exist)
	mockQueries.On("GetCartItem", ctx, pgUUID, int32(productID)).
		Return(database.CartItem{}, sql.ErrNoRows)

	// Mock AddToCart with negative quantity
	mockQueries.On("AddToCart", ctx, pgUUID, int32(productID), int32(quantity)).
		Return(nil)

	// Execute
	err := repo.AddToCart(ctx, userID, productID, quantity)

	// Assert
	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestRemoveFromCart_EmptyCart(t *testing.T) {
	repo, mockQueries := setupCartTestRepository()
	ctx := context.Background()

	userID := uuid.New()
	productID := 100

	pgUUID := pgtype.UUID{Bytes: userID, Valid: true}

	mockQueries.On("RemoveFromCart", ctx, pgUUID, int32(productID)).
		Return(errors.New("cart is empty"))

	// Execute
	err := repo.RemoveFromCart(ctx, userID, productID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cart is empty", err.Error())
	mockQueries.AssertExpectations(t)
}
