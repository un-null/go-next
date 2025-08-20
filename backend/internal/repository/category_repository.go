package repository

import (
	"context"
	"fmt"

	"backend/internal/database"
	"backend/internal/entity"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*entity.Category, error)
}

type categoryRepository struct {
	queries *database.Queries
}

func NewCategoryRepository(queries *database.Queries) CategoryRepository {
	return &categoryRepository{
		queries: queries,
	}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	dbCategories, err := r.queries.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
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

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	dbCategory, err := r.queries.GetCategoryByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	category := &entity.Category{
		ID:   int(dbCategory.ID),
		Name: dbCategory.Name,
	}

	return category, nil
}
