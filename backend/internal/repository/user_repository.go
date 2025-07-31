package repository

import (
	"backend/internal/entity"
)

type UserRepository interface {
	FindAllUser() []entity.User
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindAllUser() []entity.User {
	// Simulate DB with in-memory data
	return []entity.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
}
