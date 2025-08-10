package usecase

import (
	"backend/internal/entity"
	"backend/internal/repository"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestListUsers(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	users := uc.ListUsers()
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
}

func TestSignUp_Success(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	err := uc.SignUp(entity.User{Name: "Bob", Password: "secret"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify Bob exists
	user, err := repo.FindByName("Bob")
	if err != nil {
		t.Fatalf("expected user found, got error %v", err)
	}
	// Verify password is hashed correctly
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("secret")) != nil {
		t.Errorf("password hash mismatch")
	}
}

func TestSignUp_Duplicate(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	err := uc.SignUp(entity.User{Name: "Alice", Password: "secret"})
	if err == nil {
		t.Errorf("expected error for duplicate username, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	repo.CreateUser(entity.User{Name: "Charlie", Password: "mypassword"})

	// Login with correct password
	user, err := uc.Login("Charlie", "mypassword")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Charlie" {
		t.Errorf("expected Charlie, got %s", user.Name)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	repo.CreateUser(entity.User{Name: "David", Password: "correctpass"})

	_, err := uc.Login("David", "wrongpass")
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	_, err := uc.Login("NonExistent", "pass")
	if err == nil {
		t.Errorf("expected error for missing user, got nil")
	}
}
