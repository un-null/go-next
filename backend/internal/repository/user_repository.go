package repository

import (
	"backend/internal/entity"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	FindAllUser() []entity.User
	FindByName(name string) (*entity.User, error)
	CreateUser(user entity.User) error
}

type userRepository struct {
	users []entity.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: []entity.User{
			{ID: 1, Name: "Alice", Password: "$2a$12$examplehash"},
		},
	}
}

func (r *userRepository) FindAllUser() []entity.User {
	return r.users
}

func (r *userRepository) FindByName(name string) (*entity.User, error) {
	for _, user := range r.users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) CreateUser(user entity.User) error {
	if _, err := r.FindByName(user.Name); err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := entity.User{
		ID:       len(r.users) + 1,
		Name:     user.Name,
		Password: string(hashedPassword),
	}

	r.users = append(r.users, newUser)
	return nil
}
