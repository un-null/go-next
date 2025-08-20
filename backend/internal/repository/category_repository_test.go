package repository

import (
	"context"
	"errors"
	"testing"

	"backend/internal/database"
	"backend/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- Category-specific interface for testing ----

type CategoryQueriesInterface interface {
	ListActiveCategories(ctx context.Context) ([]database.Category, error)
	GetCategoryByID(ctx context.Context, id int32) (database.Category, error)
}

type MockCategoryQueries struct {
	mock.Mock
}

func (m *MockCategoryQueries) ListActiveCategories(ctx context.Context) ([]database.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]database.Category), args.Error(1)
}

func (m *MockCategoryQueries) GetCategoryByID(ctx context.Context, id int32) (database.Category, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Category), args.Error(1)
}

// ---- Test Repository Implementation ----

type testCategoryRepository struct {
	queries CategoryQueriesInterface
}

func NewTestCategoryRepository(queries CategoryQueriesInterface) CategoryRepository {
	return &testCategoryRepository{
		queries: queries,
	}
}

func (r *testCategoryRepository) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	dbCategories, err := r.queries.ListActiveCategories(ctx)
	if err != nil {
		return nil, err
	}

	categories := make([]entity.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		categories[i] = entity.Category{
			ID:   int(dbCategory.ID),
			Name: dbCategory.Name,
		}
	}

	return categories, nil
}

func (r *testCategoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	dbCategory, err := r.queries.GetCategoryByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	category := &entity.Category{
		ID:   int(dbCategory.ID),
		Name: dbCategory.Name,
	}

	return category, nil
}

// ---- Helper Functions ----

func sampleCategory(id int32, name string) database.Category {
	return database.Category{
		ID:   id,
		Name: name,
	}
}

// ---- Tests ----

func TestGetAllCategories(t *testing.T) {
	mockQ := new(MockCategoryQueries)
	repo := NewTestCategoryRepository(mockQ) // ✅ Use test repository

	dbCategories := []database.Category{
		sampleCategory(1, "Electronics"),
		sampleCategory(2, "Books"),
	}

	mockQ.On("ListActiveCategories", mock.Anything).Return(dbCategories, nil)

	categories, err := repo.GetAllCategories(context.Background())
	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Electronics", categories[0].Name)
	assert.Equal(t, "Books", categories[1].Name)

	mockQ.AssertExpectations(t)
}

func TestGetAllCategories_Error(t *testing.T) {
	mockQ := new(MockCategoryQueries)
	repo := NewTestCategoryRepository(mockQ)

	mockQ.On("ListActiveCategories", mock.Anything).Return([]database.Category(nil), errors.New("db error"))

	categories, err := repo.GetAllCategories(context.Background())
	assert.Error(t, err)
	assert.Nil(t, categories)
	assert.Contains(t, err.Error(), "db error")

	mockQ.AssertExpectations(t)
}

func TestGetCategoryByID_Found(t *testing.T) {
	mockQ := new(MockCategoryQueries)
	repo := NewTestCategoryRepository(mockQ)

	dbCategory := sampleCategory(1, "Food")
	mockQ.On("GetCategoryByID", mock.Anything, int32(1)).Return(dbCategory, nil)

	category, err := repo.GetCategoryByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, 1, category.ID)
	assert.Equal(t, "Food", category.Name)

	mockQ.AssertExpectations(t)
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	mockQ := new(MockCategoryQueries)
	repo := NewTestCategoryRepository(mockQ)

	mockQ.On("GetCategoryByID", mock.Anything, int32(999)).Return(database.Category{}, errors.New("not found"))

	category, err := repo.GetCategoryByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "not found")

	mockQ.AssertExpectations(t)
}
