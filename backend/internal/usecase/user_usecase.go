package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: r,
	}
}

func (u *UserUseCase) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetUserById(ctx, id)
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}

func (u *UserUseCase) SignUp(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	return u.repo.CreateUser(ctx, req)
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (u *UserUseCase) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	return u.repo.UpdateUserName(ctx, id, name)
}

func (u *UserUseCase) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	return u.repo.UpdateUserEmail(ctx, id, email)
}

func (u *UserUseCase) UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error) {
	user, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	newBalance := user.Coins + coinsDelta
	if newBalance < 0 {
		return nil, errors.New("insufficient coins")
	}

	return u.repo.UpdateUserCoins(ctx, id, coinsDelta)
}

func (u *UserUseCase) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	user, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect")
	}

	return u.repo.UpdateUserPassword(ctx, id, newPassword)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(ctx, id)
}

func (u *UserUseCase) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	return u.repo.CheckEmailExists(ctx, email)
}

func (u *UserUseCase) PurchaseItem(ctx context.Context, userID uuid.UUID, itemPrice int) (*entity.User, error) {
	if itemPrice <= 0 {
		return nil, errors.New("invalid item price")
	}

	return u.UpdateUserCoins(ctx, userID, -itemPrice)
}
