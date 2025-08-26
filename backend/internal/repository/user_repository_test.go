package repository

import (
	"backend/internal/database"
	"backend/internal/entity"
	"backend/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Test repository wrapper that implements UserRepository interface
type testUserRepository struct {
	queries mocks.UserQueriesInterface
}

func NewTestUserRepository(queries mocks.UserQueriesInterface) UserRepository {
	return &testUserRepository{
		queries: queries,
	}
}

func (r *testUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
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

func (r *testUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
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

func (r *testUserRepository) CreateUser(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	// Check if email already exists
	exists, err := r.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
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
		Coins:        database.Int32ToPgtype(int32(req.Coins)),
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

func (r *testUserRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
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

func (r *testUserRepository) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error) {
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

func (r *testUserRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	err = r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:           database.UUIDToPgtype(id),
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *testUserRepository) UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error) {
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

func (r *testUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, database.UUIDToPgtype(id))
	if err != nil {
		return err
	}
	return nil
}

func (r *testUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Helper functions for tests
func createMockDBUser(id uuid.UUID, name, email, passwordHash string, coins int) database.User {
	return database.User{
		ID:           database.UUIDToPgtype(id),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Coins:        database.Int32ToPgtype(int32(coins)),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func setupUserTestRepository() (UserRepository, *mocks.MockUserQueries) {
	mockQueries := new(mocks.MockUserQueries)
	repo := NewTestUserRepository(mockQueries)
	return repo, mockQueries
}

// Tests
func TestGetUserById_Found(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	expectedDBUser := createMockDBUser(testUserID, "Alice", "alice@example.com", "hashedpassword", 100)

	mockQueries.On("GetUserByID", ctx, database.UUIDToPgtype(testUserID)).
		Return(expectedDBUser, nil)

	user, err := repo.GetUserById(ctx, testUserID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUserID, user.ID)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, 100, user.Coins)

	mockQueries.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()

	mockQueries.On("GetUserByID", ctx, database.UUIDToPgtype(testUserID)).
		Return(database.User{}, sql.ErrNoRows)

	user, err := repo.GetUserById(ctx, testUserID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestGetUserByEmail_Found(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testEmail := "alice@example.com"
	testUserID := uuid.New()
	expectedDBUser := createMockDBUser(testUserID, "Alice", testEmail, "hashedpassword", 50)

	mockQueries.On("GetUserByEmail", ctx, testEmail).
		Return(expectedDBUser, nil)

	user, err := repo.GetUserByEmail(ctx, testEmail)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, "Alice", user.Name)

	mockQueries.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	mockQueries.On("GetUserByEmail", ctx, testEmail).
		Return(database.User{}, sql.ErrNoRows)

	user, err := repo.GetUserByEmail(ctx, testEmail)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret123",
		Coins:    0,
	}

	// Mock CheckEmailExists to return false
	mockQueries.On("CheckEmailExists", ctx, req.Email).
		Return(false, nil)

	// Mock CreateUser
	createdUserID := uuid.New()
	mockQueries.On("CreateUser", ctx, mock.MatchedBy(func(params database.CreateUserParams) bool {
		return params.Name == req.Name &&
			params.Email == req.Email &&
			database.PgtypeToInt32(params.Coins) == int32(req.Coins) &&
			len(params.PasswordHash) > 0
	})).Return(createMockDBUser(createdUserID, req.Name, req.Email, "hashedpassword", req.Coins), nil)

	user, err := repo.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Coins, user.Coins)

	mockQueries.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "existing@example.com",
		Password: "secret123",
		Coins:    0,
	}

	mockQueries.On("CheckEmailExists", ctx, req.Email).
		Return(true, nil)

	user, err := repo.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserName_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newName := "Updated Name"

	updatedUser := createMockDBUser(testUserID, newName, "test@example.com", "hashedpassword", 100)

	mockQueries.On("UpdateUserName", ctx, database.UpdateUserNameParams{
		ID:   database.UUIDToPgtype(testUserID),
		Name: newName,
	}).Return(updatedUser, nil)

	user, err := repo.UpdateUserName(ctx, testUserID, newName)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newName, user.Name)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserEmail_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newEmail := "newemail@example.com"

	// Mock email check
	mockQueries.On("CheckEmailExistsForOtherUser", ctx, database.CheckEmailExistsForOtherUserParams{
		Email: newEmail,
		ID:    database.UUIDToPgtype(testUserID),
	}).Return(false, nil)

	// Mock update
	updatedUser := createMockDBUser(testUserID, "Test User", newEmail, "hashedpassword", 100)
	mockQueries.On("UpdateUserEmail", ctx, database.UpdateUserEmailParams{
		ID:    database.UUIDToPgtype(testUserID),
		Email: newEmail,
	}).Return(updatedUser, nil)

	user, err := repo.UpdateUserEmail(ctx, testUserID, newEmail)

	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserEmail_EmailTaken(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newEmail := "taken@example.com"

	mockQueries.On("CheckEmailExistsForOtherUser", ctx, database.CheckEmailExistsForOtherUserParams{
		Email: newEmail,
		ID:    database.UUIDToPgtype(testUserID),
	}).Return(true, nil)

	user, err := repo.UpdateUserEmail(ctx, testUserID, newEmail)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newPassword := "newpassword123"

	// Mock password update
	mockQueries.On("UpdateUserPassword", ctx, mock.MatchedBy(func(params database.UpdateUserPasswordParams) bool {
		return database.PgtypeToUUID(params.ID) == testUserID && len(params.PasswordHash) > 0
	})).Return(nil)

	err := repo.UpdateUserPassword(ctx, testUserID, newPassword)

	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestUpdateUserCoins_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	coinsDelta := 50

	updatedUser := createMockDBUser(testUserID, "Test User", "test@example.com", "hashedpassword", 150)

	mockQueries.On("UpdateUserCoins", ctx, database.UpdateUserCoinsParams{
		ID:    database.UUIDToPgtype(testUserID),
		Coins: database.Int32ToPgtype(int32(coinsDelta)),
	}).Return(updatedUser, nil)

	user, err := repo.UpdateUserCoins(ctx, testUserID, coinsDelta)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 150, user.Coins)
	assert.Equal(t, testUserID, user.ID)

	mockQueries.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()

	mockQueries.On("DeleteUser", ctx, database.UUIDToPgtype(testUserID)).
		Return(nil)

	err := repo.DeleteUser(ctx, testUserID)

	assert.NoError(t, err)
	mockQueries.AssertExpectations(t)
}

func TestCheckEmailExists_True(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testEmail := "existing@example.com"

	mockQueries.On("CheckEmailExists", ctx, testEmail).
		Return(true, nil)

	exists, err := repo.CheckEmailExists(ctx, testEmail)

	assert.NoError(t, err)
	assert.True(t, exists)

	mockQueries.AssertExpectations(t)
}

func TestCheckEmailExists_False(t *testing.T) {
	repo, mockQueries := setupUserTestRepository()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	mockQueries.On("CheckEmailExists", ctx, testEmail).
		Return(false, nil)

	exists, err := repo.CheckEmailExists(ctx, testEmail)

	assert.NoError(t, err)
	assert.False(t, exists)

	mockQueries.AssertExpectations(t)
}
