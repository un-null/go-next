package usecase

import (
	"context"
	"errors"
	"testing"

	"backend/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mock Repository ----

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAllProducts(ctx context.Context, page, limit int) ([]entity.Product, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductsByCategory(ctx context.Context, categoryID, page, limit int) ([]entity.Product, error) {
	args := m.Called(ctx, categoryID, page, limit)
	return args.Get(0).([]entity.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, productID, newStock int) error {
	args := m.Called(ctx, productID, newStock)
	return args.Error(0)
}

// ---- Helper Functions ----

func sampleProduct(id int, name string) entity.Product {
	return entity.Product{
		ID:            id,
		CategoryID:    1,
		Name:          name,
		Description:   "Sample description",
		Price:         1.2,
		StockQuantity: 10,
		ImageURL:      "http://example.com/image.jpg",
		AverageRating: 4.5,
		TotalComments: 5,
	}
}

// ---- Tests ----

func TestGetAllProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	expectedProducts := []entity.Product{
		sampleProduct(1, "Apple"),
		sampleProduct(2, "Banana"),
	}

	mockRepo.On("GetAllProducts", mock.Anything, 1, 10).Return(expectedProducts, nil)

	products, err := uc.GetAllProducts(context.Background(), 1, 10)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Apple", products[0].Name)
	assert.Equal(t, "Banana", products[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_Found(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	expectedProduct := sampleProduct(1, "Apple")
	mockRepo.On("GetProductByID", mock.Anything, 1).Return(&expectedProduct, nil)

	product, err := uc.GetProductByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, 1, product.ID)
	assert.Equal(t, "Apple", product.Name)

	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	mockRepo.On("GetProductByID", mock.Anything, 999).Return(nil, errors.New("product not found"))

	product, err := uc.GetProductByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestGetProductsByCategory(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	expectedProducts := []entity.Product{
		sampleProduct(1, "iPhone"),
		sampleProduct(2, "MacBook"),
	}

	mockRepo.On("GetProductsByCategory", mock.Anything, 1, 1, 10).Return(expectedProducts, nil)

	products, err := uc.GetProductsByCategory(context.Background(), 1, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "iPhone", products[0].Name)
	assert.Equal(t, "MacBook", products[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	mockRepo.On("UpdateStock", mock.Anything, 1, 5).Return(nil)

	err := uc.UpdateStock(context.Background(), 1, 5)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateStock_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewProductUseCase(mockRepo)

	mockRepo.On("UpdateStock", mock.Anything, 999, 5).Return(errors.New("product not found"))

	err := uc.UpdateStock(context.Background(), 999, 5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}
