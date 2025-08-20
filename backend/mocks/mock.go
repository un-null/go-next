package mocks

import (
	"backend/internal/database"
	"context"

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
