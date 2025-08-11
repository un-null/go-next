package repository

import (
	"backend/internal/entity"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGetAllUsers(t *testing.T) {
	repo := NewUserRepository()

	users := repo.GetAllUsers()
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", users[0].Name)
	}
}

func TestGetUserById_Found(t *testing.T) {
	repo := NewUserRepository()

	user := repo.GetUserById(1)
	if user.ID == 0 {
		t.Fatalf("expected to find user with ID 1, got zero value")
	}

	if user.Name != "Alice" {
		t.Errorf("expected user name 'Alice', got '%s'", user.Name)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	repo := NewUserRepository()

	user := repo.GetUserById(99)
	if user.ID != 0 {
		t.Fatalf("expected zero value user for not found, got ID %d", user.ID)
	}

	if user.Name != "" {
		t.Errorf("expected empty name for not found, got '%s'", user.Name)
	}
}

func TestFindByName_Found(t *testing.T) {
	repo := NewUserRepository()

	user, err := repo.GetUserByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", user.Name)
	}
}

func TestFindByName_NotFound(t *testing.T) {
	repo := NewUserRepository()

	_, err := repo.GetUserByEmail("bob@example.com")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreateUser_Success(t *testing.T) {
	repo := NewUserRepository()

	err := repo.CreateUser(entity.User{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user, err := repo.GetUserByEmail("bob@example.com")
	if err != nil {
		t.Fatalf("expected user found, got error %v", err)
	}

	if user.Name != "Bob" {
		t.Errorf("expected name 'Bob', got '%s'", user.Name)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret")) != nil {
		t.Errorf("password hash does not match")
	}

	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Errorf("expected timestamps to be set")
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	repo := NewUserRepository()

	err := repo.CreateUser(entity.User{Name: "Alice", Email: "alice@example.com", Password: "secret"})
	if err == nil {
		t.Errorf("expected error for duplicate user, got nil")
	}
}
