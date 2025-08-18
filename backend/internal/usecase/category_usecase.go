package usecase

import (
	"context"

	"backend/internal/entity"
	"backend/internal/repository"
)

type CategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewCategoryUseCase(r repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: r}
}

func (u *CategoryUseCase) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	return u.repo.GetAllCategories(ctx)
}

func (u *CategoryUseCase) GetCategoryByID(ctx context.Context, id int) (*entity.Category, error) {
	return u.repo.GetCategoryByID(ctx, id)
}
