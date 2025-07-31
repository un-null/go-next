package usecase

import (
	"backend/internal/entity"
	"testing"
)

type mockCartRepo struct {
	AddCalled    bool
	RemoveCalled bool
	ClearCalled  bool
	GetCalled    bool
	MockItems    []entity.CartItem
	MockError    error
}

func (m *mockCartRepo) AddToCart(userID, productID, quantity int) error {
	m.AddCalled = true
	return m.MockError
}

func (m *mockCartRepo) GetCartItems(userID int) ([]entity.CartItem, error) {
	m.GetCalled = true
	return m.MockItems, m.MockError
}

func (m *mockCartRepo) RemoveFromCart(userID, productID int) error {
	m.RemoveCalled = true
	return m.MockError
}

func (m *mockCartRepo) ClearCart(userID int) error {
	m.ClearCalled = true
	return m.MockError
}

func TestAddToCart(t *testing.T) {
	mockRepo := &mockCartRepo{}
	uc := NewCartUseCase(mockRepo)

	err := uc.AddToCart(1, 100, 2)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !mockRepo.AddCalled {
		t.Error("expected AddToCart to be called on repository")
	}
}

func TestGetCartItems(t *testing.T) {
	mockRepo := &mockCartRepo{
		MockItems: []entity.CartItem{
			{UserID: 1, ProductID: 100, Quantity: 2},
		},
	}
	uc := NewCartUseCase(mockRepo)

	items, err := uc.GetCartItems(1)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !mockRepo.GetCalled {
		t.Error("expected GetCartItems to be called on repository")
	}
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}
	if items[0].UserID != 1 || items[0].ProductID != 100 || items[0].Quantity != 2 {
		t.Errorf("expected specific item, got %+v", items[0])
	}
}

func TestRemoveFromCart(t *testing.T) {
	mockRepo := &mockCartRepo{}
	uc := NewCartUseCase(mockRepo)

	err := uc.RemoveFromCart(1, 100)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !mockRepo.RemoveCalled {
		t.Error("expected RemoveFromCart to be called on repository")
	}
}

func TestClearCart(t *testing.T) {
	mockRepo := &mockCartRepo{}
	uc := NewCartUseCase(mockRepo)

	err := uc.ClearCart(1)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !mockRepo.ClearCalled {
		t.Error("expected ClearCart to be called on repository")
	}
}
