package mocks

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

// UserQueriesInterface defines the interface for user-related database operations
type UserQueriesInterface interface {
	GetUserByID(ctx context.Context, id pgtype.UUID) (database.User, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
	CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error)
	UpdateUserName(ctx context.Context, params database.UpdateUserNameParams) (database.User, error)
	UpdateUserEmail(ctx context.Context, params database.UpdateUserEmailParams) (database.User, error)
	UpdateUserCoins(ctx context.Context, params database.UpdateUserCoinsParams) (database.User, error)
	UpdateUserPassword(ctx context.Context, params database.UpdateUserPasswordParams) error
	DeleteUser(ctx context.Context, id pgtype.UUID) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckEmailExistsForOtherUser(ctx context.Context, params database.CheckEmailExistsForOtherUserParams) (bool, error)
}

// MockUserQueries is a mock implementation for user-related database operations
type MockUserQueries struct {
	mock.Mock
}

func (m *MockUserQueries) GetUserByID(ctx context.Context, id pgtype.UUID) (database.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) CreateUser(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) UpdateUserName(ctx context.Context, params database.UpdateUserNameParams) (database.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) UpdateUserEmail(ctx context.Context, params database.UpdateUserEmailParams) (database.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) UpdateUserCoins(ctx context.Context, params database.UpdateUserCoinsParams) (database.User, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.User), args.Error(1)
}

func (m *MockUserQueries) UpdateUserPassword(ctx context.Context, params database.UpdateUserPasswordParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockUserQueries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserQueries) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserQueries) CheckEmailExistsForOtherUser(ctx context.Context, params database.CheckEmailExistsForOtherUserParams) (bool, error) {
	args := m.Called(ctx, params)
	return args.Bool(0), args.Error(1)
}

// ProductQueriesInterface defines the interface for product-related database operations
type ProductQueriesInterface interface {
	ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error)
	GetProductByID(ctx context.Context, id int32) (database.Product, error)
	ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error)
	UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error)
}

// MockProductQueries is a mock implementation for product-related database operations
type MockProductQueries struct {
	mock.Mock
}

func (m *MockProductQueries) ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Product), args.Error(1)
}

func (m *MockProductQueries) GetProductByID(ctx context.Context, id int32) (database.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Product), args.Error(1)
}

func (m *MockProductQueries) ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Product), args.Error(1)
}

func (m *MockProductQueries) UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Product), args.Error(1)
}

type CartQueriesInterface interface {
	AddToCart(ctx context.Context, userID pgtype.UUID, productID int32, quantity int32) error
	GetCartItems(ctx context.Context, userID pgtype.UUID) ([]database.CartItem, error)
	RemoveFromCart(ctx context.Context, userID pgtype.UUID, productID int32) error
	ClearCart(ctx context.Context, userID pgtype.UUID) error
	GetCartItem(ctx context.Context, userID pgtype.UUID, productID int32) (database.CartItem, error)
	UpdateCartItemQuantity(ctx context.Context, userID pgtype.UUID, productID int32, quantity int32) error
}

// MockCartQueries is a mock implementation for cart-related database operations
type MockCartQueries struct {
	mock.Mock
}

func (m *MockCartQueries) AddToCart(ctx context.Context, userID pgtype.UUID, productID int32, quantity int32) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}

func (m *MockCartQueries) GetCartItems(ctx context.Context, userID pgtype.UUID) ([]database.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]database.CartItem), args.Error(1)
}

func (m *MockCartQueries) RemoveFromCart(ctx context.Context, userID pgtype.UUID, productID int32) error {
	args := m.Called(ctx, userID, productID)
	return args.Error(0)
}

func (m *MockCartQueries) ClearCart(ctx context.Context, userID pgtype.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCartQueries) GetCartItem(ctx context.Context, userID pgtype.UUID, productID int32) (database.CartItem, error) {
	args := m.Called(ctx, userID, productID)
	return args.Get(0).(database.CartItem), args.Error(1)
}

func (m *MockCartQueries) UpdateCartItemQuantity(ctx context.Context, userID pgtype.UUID, productID int32, quantity int32) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}

// CoinTransactionQueriesInterface defines the interface for coin transaction database operations
type CoinTransactionQueriesInterface interface {
	CreateCoinTransaction(ctx context.Context, arg database.CreateCoinTransactionParams) (database.CoinTransaction, error)
	GetCoinTransactionsByUserID(ctx context.Context, arg database.GetCoinTransactionsByUserIDParams) ([]database.CoinTransaction, error)
	GetCoinTransactionByID(ctx context.Context, id int32) (database.CoinTransaction, error)
}

// MockCoinTransactionQueries is a mock implementation for coin transaction database operations
type MockCoinTransactionQueries struct {
	mock.Mock
}

func (m *MockCoinTransactionQueries) CreateCoinTransaction(ctx context.Context, arg database.CreateCoinTransactionParams) (database.CoinTransaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.CoinTransaction), args.Error(1)
}

func (m *MockCoinTransactionQueries) GetCoinTransactionsByUserID(ctx context.Context, arg database.GetCoinTransactionsByUserIDParams) ([]database.CoinTransaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.CoinTransaction), args.Error(1)
}

func (m *MockCoinTransactionQueries) GetCoinTransactionByID(ctx context.Context, id int32) (database.CoinTransaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.CoinTransaction), args.Error(1)
}

// MockOrderQuerier is a mock implementation for order-related database operations
// that matches repository.OrderQuerier
type MockOrderQuerier struct {
	mock.Mock
}

func (m *MockOrderQuerier) GetOrderByID(ctx context.Context, id int32) (database.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Order), args.Error(1)
}

func (m *MockOrderQuerier) GetOrdersByUserID(ctx context.Context, userID pgtype.UUID) ([]database.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]database.Order), args.Error(1)
}

func (m *MockOrderQuerier) CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Order), args.Error(1)
}

func (m *MockOrderQuerier) WithTx(tx pgx.Tx) *database.Queries { return database.New(tx) }

// MockTx is a mock implementation for database transactions
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockDB is a mock implementation for database pool operations
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockDB) Close() {
	m.Called()
}

// MockOrderItemQuerier implements repository.OrderItemQuerier for repository tests
type MockOrderItemQuerier struct {
	mock.Mock
}

func (m *MockOrderItemQuerier) GetOrderItemByID(ctx context.Context, id int32) (database.OrderItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.OrderItem), args.Error(1)
}

func (m *MockOrderItemQuerier) GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]database.OrderItem, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]database.OrderItem), args.Error(1)
}

func (m *MockOrderItemQuerier) CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.OrderItem), args.Error(1)
}

func (m *MockOrderItemQuerier) WithTx(tx pgx.Tx) *database.Queries { return database.New(tx) }

// MockOrderRepository implements repository.OrderRepository for usecase tests
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, req entity.CreateOrderRequest) (*entity.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) CreateOrderWithTx(ctx context.Context, tx pgx.Tx, req entity.CreateOrderRequest) (*entity.Order, error) {
	args := m.Called(ctx, tx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

// MockOrderItemRepository implements repository.OrderItemRepository for usecase tests
type MockOrderItemRepository struct {
	mock.Mock
}

func (m *MockOrderItemRepository) GetOrderItemById(ctx context.Context, id int) (*entity.OrderItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]entity.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepository) CreateOrderItem(ctx context.Context, req entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepository) CreateOrderItemWithTx(ctx context.Context, tx pgx.Tx, req entity.OrderItem) (*entity.OrderItem, error) {
	args := m.Called(ctx, tx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *MockOrderItemRepository) CreateOrderItemsWithTx(ctx context.Context, tx pgx.Tx, items []entity.OrderItem) ([]entity.OrderItem, error) {
	args := m.Called(ctx, tx, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.OrderItem), args.Error(1)
}

// ---- Cart Repository and Querier Mocks ----

// MockCartRepository implements repository.CartRepository for usecase tests
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

func (m *MockCartRepository) ClearCartWithTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	args := m.Called(ctx, tx, userID)
	return args.Error(0)
}

// MockCartQuerier implements repository.CartQuerier for repository tests
type MockCartQuerier struct {
	mock.Mock
}

func (m *MockCartQuerier) CheckCartItemExists(ctx context.Context, arg database.CheckCartItemExistsParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Bool(0), args.Error(1)
}

func (m *MockCartQuerier) CreateCartItem(ctx context.Context, arg database.CreateCartItemParams) (database.CartItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.CartItem), args.Error(1)
}

func (m *MockCartQuerier) UpdateCartItemQuantity(ctx context.Context, arg database.UpdateCartItemQuantityParams) (database.CartItem, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.CartItem), args.Error(1)
}

func (m *MockCartQuerier) GetCartItemsByUser(ctx context.Context, userID pgtype.UUID) ([]database.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]database.CartItem), args.Error(1)
}

func (m *MockCartQuerier) DeleteCartItem(ctx context.Context, arg database.DeleteCartItemParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockCartQuerier) DeleteAllCartItemsByUser(ctx context.Context, userID pgtype.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCartQuerier) WithTx(tx pgx.Tx) *database.Queries { return database.New(tx) }

// Empty implementations to satisfy database.Querier
func (m *MockCartQuerier) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *MockCartQuerier) CheckEmailExistsForOtherUser(ctx context.Context, arg database.CheckEmailExistsForOtherUserParams) (bool, error) {
	return false, nil
}
func (m *MockCartQuerier) CreateCoinTransaction(ctx context.Context, arg database.CreateCoinTransactionParams) (database.CoinTransaction, error) {
	return database.CoinTransaction{}, nil
}
func (m *MockCartQuerier) CreateOrder(ctx context.Context, arg database.CreateOrderParams) (database.Order, error) {
	return database.Order{}, nil
}
func (m *MockCartQuerier) CreateOrderItem(ctx context.Context, arg database.CreateOrderItemParams) (database.OrderItem, error) {
	return database.OrderItem{}, nil
}
func (m *MockCartQuerier) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	return database.CreateUserRow{}, nil
}
func (m *MockCartQuerier) DeleteUser(ctx context.Context, id pgtype.UUID) error { return nil }
func (m *MockCartQuerier) GetAllCategories(ctx context.Context) ([]database.Category, error) {
	return nil, nil
}
func (m *MockCartQuerier) GetCategoryByID(ctx context.Context, id int32) (database.Category, error) {
	return database.Category{}, nil
}
func (m *MockCartQuerier) GetCoinTransactionByID(ctx context.Context, id int32) (database.CoinTransaction, error) {
	return database.CoinTransaction{}, nil
}
func (m *MockCartQuerier) GetCoinTransactionsByUserID(ctx context.Context, arg database.GetCoinTransactionsByUserIDParams) ([]database.CoinTransaction, error) {
	return nil, nil
}
func (m *MockCartQuerier) GetOrderByID(ctx context.Context, id int32) (database.Order, error) {
	return database.Order{}, nil
}
func (m *MockCartQuerier) GetOrderItemByID(ctx context.Context, id int32) (database.OrderItem, error) {
	return database.OrderItem{}, nil
}
func (m *MockCartQuerier) GetOrderItemsByOrderID(ctx context.Context, orderID int32) ([]database.OrderItem, error) {
	return nil, nil
}
func (m *MockCartQuerier) GetOrdersByUserID(ctx context.Context, userID pgtype.UUID) ([]database.Order, error) {
	return nil, nil
}
func (m *MockCartQuerier) GetProductByID(ctx context.Context, id int32) (database.Product, error) {
	return database.Product{}, nil
}
func (m *MockCartQuerier) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	return database.User{}, nil
}
func (m *MockCartQuerier) GetUserByID(ctx context.Context, id pgtype.UUID) (database.User, error) {
	return database.User{}, nil
}
func (m *MockCartQuerier) ListProducts(ctx context.Context, arg database.ListProductsParams) ([]database.Product, error) {
	return nil, nil
}
func (m *MockCartQuerier) ListProductsByCategory(ctx context.Context, arg database.ListProductsByCategoryParams) ([]database.Product, error) {
	return nil, nil
}
func (m *MockCartQuerier) UpdateProductStock(ctx context.Context, arg database.UpdateProductStockParams) (database.Product, error) {
	return database.Product{}, nil
}
func (m *MockCartQuerier) UpdateUserCoins(ctx context.Context, arg database.UpdateUserCoinsParams) (database.UpdateUserCoinsRow, error) {
	return database.UpdateUserCoinsRow{}, nil
}
func (m *MockCartQuerier) UpdateUserEmail(ctx context.Context, arg database.UpdateUserEmailParams) (database.UpdateUserEmailRow, error) {
	return database.UpdateUserEmailRow{}, nil
}
func (m *MockCartQuerier) UpdateUserName(ctx context.Context, arg database.UpdateUserNameParams) (database.UpdateUserNameRow, error) {
	return database.UpdateUserNameRow{}, nil
}
func (m *MockCartQuerier) UpdateUserPassword(ctx context.Context, arg database.UpdateUserPasswordParams) error {
	return nil
}
