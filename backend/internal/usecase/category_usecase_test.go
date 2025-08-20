package usecase

import (
	"context"
	"errors"
	"testing"

	"backend/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Mock Category Repository ----

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

// ---- Helper Functions ----

func sampleCategory(id int, name string) entity.Category {
	return entity.Category{
		ID:   id,
		Name: name,
	}
}

// ---- Tests ----

func TestCategoryUseCase_GetAllCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	expectedCategories := []entity.Category{
		sampleCategory(1, "Electronics"),
		sampleCategory(2, "Books"),
		sampleCategory(3, "Clothing"),
	}

	mockRepo.On("GetAllCategories", mock.Anything).Return(expectedCategories, nil)

	categories, err := uc.GetAllCategories(context.Background())

	assert.NoError(t, err)
	assert.Len(t, categories, 3)
	assert.Equal(t, "Electronics", categories[0].Name)
	assert.Equal(t, "Books", categories[1].Name)
	assert.Equal(t, "Clothing", categories[2].Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetAllCategories_Empty(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	emptyCategories := []entity.Category{}

	mockRepo.On("GetAllCategories", mock.Anything).Return(emptyCategories, nil)

	categories, err := uc.GetAllCategories(context.Background())

	assert.NoError(t, err)
	assert.Len(t, categories, 0)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetAllCategories_Error(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	mockRepo.On("GetAllCategories", mock.Anything).Return([]entity.Category(nil), errors.New("database connection failed"))

	categories, err := uc.GetAllCategories(context.Background())

	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Contains(t, err.Error(), "database connection failed")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetCategoryByID_Found(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	expectedCategory := sampleCategory(1, "Electronics")

	mockRepo.On("GetCategoryByID", mock.Anything, 1).Return(&expectedCategory, nil)

	category, err := uc.GetCategoryByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Electronics", category.Name)

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetCategoryByID_NotFound(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	mockRepo.On("GetCategoryByID", mock.Anything, 999).Return(nil, errors.New("category not found"))

	category, err := uc.GetCategoryByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "not found")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetCategoryByID_InvalidID(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	mockRepo.On("GetCategoryByID", mock.Anything, 0).Return(nil, errors.New("invalid category ID"))

	category, err := uc.GetCategoryByID(context.Background(), 0)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "invalid")

	mockRepo.AssertExpectations(t)
}

func TestCategoryUseCase_GetCategoryByID_DatabaseError(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	mockRepo.On("GetCategoryByID", mock.Anything, 1).Return(nil, errors.New("database timeout"))

	category, err := uc.GetCategoryByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "timeout")

	mockRepo.AssertExpectations(t)
}

// ---- Benchmark Tests (Optional) ----

func BenchmarkCategoryUseCase_GetAllCategories(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	categories := []entity.Category{
		sampleCategory(1, "Electronics"),
		sampleCategory(2, "Books"),
	}

	mockRepo.On("GetAllCategories", mock.Anything).Return(categories, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.GetAllCategories(context.Background())
	}
}

func BenchmarkCategoryUseCase_GetCategoryByID(b *testing.B) {
	mockRepo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(mockRepo)

	category := sampleCategory(1, "Electronics")
	mockRepo.On("GetCategoryByID", mock.Anything, 1).Return(&category, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.GetCategoryByID(context.Background(), 1)
	}
}
