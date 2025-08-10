package repository

import (
	"backend/internal/entity"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestFindAllUser(t *testing.T) {
	repo := NewUserRepository()

	users := repo.FindAllUser()
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", users[0].Name)
	}
}

func TestFindByName_Found(t *testing.T) {
	repo := NewUserRepository()

	user, err := repo.FindByName("Alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected 'Alice', got '%s'", user.Name)
	}
}

func TestFindByName_NotFound(t *testing.T) {
	repo := NewUserRepository()

	_, err := repo.FindByName("Bob")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreateUser_Success(t *testing.T) {
	repo := NewUserRepository()

	err := repo.CreateUser(entity.User{Name: "Bob", Password: "secret"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify Bob exists now
	user, err := repo.FindByName("Bob")
	if err != nil {
		t.Fatalf("expected user found, got error %v", err)
	}

	// Password should be hashed
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret")) != nil {
		t.Errorf("password hash does not match")
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	repo := NewUserRepository()

	err := repo.CreateUser(entity.User{Name: "Alice", Password: "secret"})
	if err == nil {
		t.Errorf("expected error for duplicate user, got nil")
	}
}
