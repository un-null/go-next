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

	users := uc.GetAllUsers()
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Name != "Alice" {
		t.Errorf("expected Alice, got %s", users[0].Name)
	}
}

func TestGetUserById_Found(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	product := uc.GetUserById(1)
	if product.ID == 0 {
		t.Fatalf("expected to find product with ID 1, got zero value")
	}

	if product.Name != "Alice" {
		t.Errorf("expected product name 'Apple', got '%s'", product.Name)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	product := uc.GetUserById(99)
	if product.ID != 0 {
		t.Errorf("expected zero value product for non-existing ID, got %+v", product)
	}
}

func TestSignUp_Success(t *testing.T) {
	repo := repository.NewUserRepository()
	uc := NewUserUseCase(repo)

	err := uc.SignUp(entity.User{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify Bob exists by email
	user, err := repo.GetUserByEmail("bob@example.com")
	if err != nil {
		t.Fatalf("expected user found, got error %v", err)
	}

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
