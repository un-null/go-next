package mocks

import (
	"backend/internal/database"
	"context"

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
