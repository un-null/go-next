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

// Test repository wrapper that uses the interface
type testUserRepository struct {
	queries mocks.UserQueriesInterface
}

func NewTestUserRepository(queries mocks.UserQueriesInterface) UserRepository {
	return &testUserRepository{
		queries: queries,
	}
}

// Implement the same methods as the original repository
func (r *testUserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:           dbUser.ID.Bytes,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Coins:        int(dbUser.Coins.Int32),
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
		ID:           dbUser.ID.Bytes,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Coins:        int(dbUser.Coins.Int32),
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
		Coins:        pgtype.Int4{Int32: int32(req.Coins), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(dbUser.Coins.Int32),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *testUserRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUserName(ctx, database.UpdateUserNameParams{
		ID:   pgtype.UUID{Bytes: id, Valid: true},
		Name: name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(dbUser.Coins.Int32),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *testUserRepository) UpdateUserEmail(ctx context.Context, id uuid.UUID, email string) (*entity.User, error) {
	// Check if new email already exists for another user
	exists, err := r.queries.CheckEmailExistsForOtherUser(ctx, database.CheckEmailExistsForOtherUserParams{
		Email: email,
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	dbUser, err := r.queries.UpdateUserEmail(ctx, database.UpdateUserEmailParams{
		ID:    pgtype.UUID{Bytes: id, Valid: true},
		Email: email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(dbUser.Coins.Int32),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *testUserRepository) UpdateUserCoins(ctx context.Context, id uuid.UUID, coinsDelta int) (*entity.User, error) {
	dbUser, err := r.queries.UpdateUserCoins(ctx, database.UpdateUserCoinsParams{
		ID:    pgtype.UUID{Bytes: id, Valid: true},
		Coins: pgtype.Int4{Int32: int32(coinsDelta), Valid: true},
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user := &entity.User{
		ID:        dbUser.ID.Bytes,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Coins:     int(dbUser.Coins.Int32),
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	return user, nil
}

func (r *testUserRepository) UpdateUserPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = r.queries.UpdateUserPassword(ctx, database.UpdateUserPasswordParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *testUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
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
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Coins:        pgtype.Int4{Int32: int32(coins), Valid: true},
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
}

func setupTestRepository() (UserRepository, *mocks.MockUserQueries) {
	mockQueries := new(mocks.MockUserQueries)
	repo := NewTestUserRepository(mockQueries)
	return repo, mockQueries
}

func TestGetUserById_Found(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	// Test data
	testUserID := uuid.New()
	expectedDBUser := createMockDBUser(testUserID, "Alice", "alice@example.com", "hashedpassword", 100)

	// Set up mock expectation
	mockQueries.On("GetUserByID", ctx, pgtype.UUID{Bytes: testUserID, Valid: true}).
		Return(expectedDBUser, nil)

	// Execute
	user, err := repo.GetUserById(ctx, testUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUserID, user.ID)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.PasswordHash)
	assert.Equal(t, 100, user.Coins)

	// Verify mock expectations
	mockQueries.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()

	// Set up mock to return sql.ErrNoRows
	mockQueries.On("GetUserByID", ctx, pgtype.UUID{Bytes: testUserID, Valid: true}).
		Return(database.User{}, sql.ErrNoRows)

	// Execute
	user, err := repo.GetUserById(ctx, testUserID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestGetUserByEmail_Found(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testEmail := "alice@example.com"
	testUserID := uuid.New()
	expectedDBUser := createMockDBUser(testUserID, "Alice", testEmail, "hashedpassword", 50)

	mockQueries.On("GetUserByEmail", ctx, testEmail).
		Return(expectedDBUser, nil)

	// Execute
	user, err := repo.GetUserByEmail(ctx, testEmail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testEmail, user.Email)
	assert.Equal(t, "Alice", user.Name)

	mockQueries.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	mockQueries.On("GetUserByEmail", ctx, testEmail).
		Return(database.User{}, sql.ErrNoRows)

	// Execute
	user, err := repo.GetUserByEmail(ctx, testEmail)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret123",
		Coins:    0,
	}

	// Mock CheckEmailExists to return false (email doesn't exist)
	mockQueries.On("CheckEmailExists", ctx, req.Email).
		Return(false, nil)

	// Mock CreateUser
	createdUserID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// We need to use mock.MatchedBy to handle the hashed password comparison
	mockQueries.On("CreateUser", ctx, mock.MatchedBy(func(params database.CreateUserParams) bool {
		return params.Name == req.Name &&
			params.Email == req.Email &&
			params.Coins.Int32 == int32(req.Coins) &&
			len(params.PasswordHash) > 0 // Just check that password is hashed
	})).Return(createMockDBUser(createdUserID, req.Name, req.Email, string(hashedPassword), req.Coins), nil)

	// Execute
	user, err := repo.CreateUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Coins, user.Coins)
	assert.NotEmpty(t, user.ID)

	mockQueries.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	req := entity.CreateUserRequest{
		Name:     "Bob",
		Email:    "existing@example.com",
		Password: "secret123",
		Coins:    0,
	}

	// Mock CheckEmailExists to return true (email already exists)
	mockQueries.On("CheckEmailExists", ctx, req.Email).
		Return(true, nil)

	// Execute
	user, err := repo.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user already exists", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestCheckEmailExists_True(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testEmail := "existing@example.com"

	mockQueries.On("CheckEmailExists", ctx, testEmail).
		Return(true, nil)

	// Execute
	exists, err := repo.CheckEmailExists(ctx, testEmail)

	// Assert
	assert.NoError(t, err)
	assert.True(t, exists)

	mockQueries.AssertExpectations(t)
}

func TestCheckEmailExists_False(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testEmail := "nonexistent@example.com"

	mockQueries.On("CheckEmailExists", ctx, testEmail).
		Return(false, nil)

	// Execute
	exists, err := repo.CheckEmailExists(ctx, testEmail)

	// Assert
	assert.NoError(t, err)
	assert.False(t, exists)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserName_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newName := "Updated Name"
	expectedParams := database.UpdateUserNameParams{
		ID:   pgtype.UUID{Bytes: testUserID, Valid: true},
		Name: newName,
	}

	updatedUser := createMockDBUser(testUserID, newName, "test@example.com", "hashedpassword", 100)

	mockQueries.On("UpdateUserName", ctx, expectedParams).
		Return(updatedUser, nil)

	// Execute
	user, err := repo.UpdateUserName(ctx, testUserID, newName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newName, user.Name)
	assert.Equal(t, testUserID, user.ID)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserName_NotFound(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newName := "Updated Name"
	expectedParams := database.UpdateUserNameParams{
		ID:   pgtype.UUID{Bytes: testUserID, Valid: true},
		Name: newName,
	}

	mockQueries.On("UpdateUserName", ctx, expectedParams).
		Return(database.User{}, sql.ErrNoRows)

	// Execute
	user, err := repo.UpdateUserName(ctx, testUserID, newName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserEmail_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newEmail := "newemail@example.com"

	// Mock CheckEmailExistsForOtherUser
	checkParams := database.CheckEmailExistsForOtherUserParams{
		Email: newEmail,
		ID:    pgtype.UUID{Bytes: testUserID, Valid: true},
	}
	mockQueries.On("CheckEmailExistsForOtherUser", ctx, checkParams).
		Return(false, nil)

	// Mock UpdateUserEmail
	updateParams := database.UpdateUserEmailParams{
		ID:    pgtype.UUID{Bytes: testUserID, Valid: true},
		Email: newEmail,
	}
	updatedUser := createMockDBUser(testUserID, "Test User", newEmail, "hashedpassword", 100)

	mockQueries.On("UpdateUserEmail", ctx, updateParams).
		Return(updatedUser, nil)

	// Execute
	user, err := repo.UpdateUserEmail(ctx, testUserID, newEmail)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newEmail, user.Email)
	assert.Equal(t, testUserID, user.ID)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserCoins_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	coinsDelta := 50
	expectedParams := database.UpdateUserCoinsParams{
		ID:    pgtype.UUID{Bytes: testUserID, Valid: true},
		Coins: pgtype.Int4{Int32: int32(coinsDelta), Valid: true},
	}

	updatedUser := createMockDBUser(testUserID, "Test User", "test@example.com", "hashedpassword", 150)

	mockQueries.On("UpdateUserCoins", ctx, expectedParams).
		Return(updatedUser, nil)

	// Execute
	user, err := repo.UpdateUserCoins(ctx, testUserID, coinsDelta)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 150, user.Coins)
	assert.Equal(t, testUserID, user.ID)

	mockQueries.AssertExpectations(t)
}

func TestUpdateUserPassword_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	newPassword := "newpassword123"

	// Mock UpdateUserPassword - we use mock.MatchedBy to handle hashed password
	mockQueries.On("UpdateUserPassword", ctx, mock.MatchedBy(func(params database.UpdateUserPasswordParams) bool {
		return params.ID.Bytes == testUserID && len(params.PasswordHash) > 0
	})).Return(nil)

	// Execute
	err := repo.UpdateUserPassword(ctx, testUserID, newPassword)

	// Assert
	assert.NoError(t, err)

	mockQueries.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	repo, mockQueries := setupTestRepository()
	ctx := context.Background()

	testUserID := uuid.New()
	expectedPgUUID := pgtype.UUID{Bytes: testUserID, Valid: true}

	mockQueries.On("DeleteUser", ctx, expectedPgUUID).
		Return(nil)

	// Execute
	err := repo.DeleteUser(ctx, testUserID)

	// Assert
	assert.NoError(t, err)

	mockQueries.AssertExpectations(t)
}
