package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) GetAllUsers() []entity.User {
	return u.repo.GetAllUsers()
}

func (u *UserUseCase) GetUserById(id int) entity.User {
	return u.repo.GetUserById(id)
}

func (u *UserUseCase) SignUp(user entity.User) error {
	return u.repo.CreateUser(user)
}

func (u *UserUseCase) Login(name, password string) (*entity.User, error) {
	user, err := u.repo.FindByName(name)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
