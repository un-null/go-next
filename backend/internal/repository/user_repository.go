package repository

import (
	"backend/internal/entity"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAllUsers() []entity.User
	GetUserById(id int) entity.User
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(user entity.User) error
}

type userRepository struct {
	users []entity.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: []entity.User{
			{
				ID:        1,
				Name:      "Alice",
				Email:     "alice@example.com",
				Password:  "$2a$12$examplehash",
				Coins:     100,
				CreatedAt: time.Now().AddDate(-1, 0, 0), // 1 year ago
				UpdatedAt: time.Now(),
			},
		},
	}
}

func (r *userRepository) GetAllUsers() []entity.User {
	return r.users
}

func (r *userRepository) GetUserById(id int) entity.User {
	for _, u := range r.users {
		if u.ID == id {
			return u
		}
	}

	return entity.User{}
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	for i := range r.users {
		if r.users[i].Email == email {
			return &r.users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *userRepository) CreateUser(user entity.User) error {
	if _, err := r.GetUserByEmail(user.Email); err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()

	newUser := entity.User{
		ID:        len(r.users) + 1,
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hashedPassword),
		Coins:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.users = append(r.users, newUser)
	return nil
}
