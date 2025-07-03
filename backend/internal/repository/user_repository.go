package repository

import (
	"backend/internal/entity"
)

type UserRepository interface {
	FindAll() []entity.User
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindAll() []entity.User {
	// Simulate DB with in-memory data
	return []entity.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
}
