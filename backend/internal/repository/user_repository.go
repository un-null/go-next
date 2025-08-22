package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error)
	UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error)
	UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error
	UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	queries *database.Queries
}

func NewUserRepository(queries *database.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

func (r *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, database.UUIDToPgtype(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:           database.PgtypeToUUID(dbUser.ID),
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Coins:        int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:           database.PgtypeToUUID(dbUser.ID),
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Coins:        int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	// Check if email already exists
	exists, err := r.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user in database
	dbUser, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Coins:        database.Int32ToPgtype(1000),
	})
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        database.PgtypeToUUID(dbUser.ID),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUserName(ctx, database.UpdateUserNameParams{
		ID:   database.UUIDToPgtype(id),
		Name: name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        database.PgtypeToUUID(dbUser.ID),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error) {
	// Check if new email already exists for another user
	exists, err := r.queries.CheckEmailExistsForOtherUser(ctx, database.CheckEmailExistsForOtherUserParams{
		Email: email,
		ID:    database.UUIDToPgtype(id),
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	dbUser, err := r.queries.UpdateUserEmail(ctx, database.UpdateUserEmailParams{
		ID:    database.UUIDToPgtype(id),
		Email: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        database.PgtypeToUUID(dbUser.ID),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUserCoins(ctx, database.UpdateUserCoinsParams{
		ID:    database.UUIDToPgtype(id),
		Coins: database.Int32ToPgtype(int32(coinsDelta)),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        database.PgtypeToUUID(dbUser.ID),
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(database.PgtypeToInt32(dbUser.Coins)),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:           database.UUIDToPgtype(id),
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, database.UUIDToPgtype(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}
