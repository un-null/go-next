package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) ListUsers() []entity.User {
	return u.repo.FindAllUser()
}
