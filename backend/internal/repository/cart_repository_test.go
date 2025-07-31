package repository

import (
	"testing"
)

func TestAddToCart_NewItem(t *testing.T) {
	repo := NewCartRepository()

	err := repo.AddToCart(1, 100, 2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	items, _ := repo.GetCartItems(1)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ProductID != 100 || items[0].Quantity != 2 {
		t.Errorf("unexpected item: %+v", items[0])
	}
}

func TestAddToCart_ExistingItem(t *testing.T) {
	repo := NewCartRepository()
	repo.AddToCart(1, 100, 2)

	err := repo.AddToCart(1, 100, 3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	items, _ := repo.GetCartItems(1)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ProductID != 100 || items[0].Quantity != 5 {
		t.Errorf("unexpected item: %+v", items[0])
	}
}

func TestGetCartItems(t *testing.T) {
	repo := NewCartRepository()
	repo.AddToCart(1, 100, 2)
	repo.AddToCart(1, 101, 3)

	items, _ := repo.GetCartItems(1)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].ProductID != 100 || items[0].Quantity != 2 {
		t.Errorf("unexpected item: %+v", items[0])
	}
	if items[1].ProductID != 101 || items[1].Quantity != 3 {
		t.Errorf("unexpected item: %+v", items[1])
	}
}

func TestRemoveFromCart(t *testing.T) {
	repo := NewCartRepository()
	repo.AddToCart(1, 100, 2)
	repo.AddToCart(1, 101, 3)

	err := repo.RemoveFromCart(1, 101)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	items, _ := repo.GetCartItems(1)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ProductID != 100 || items[0].Quantity != 2 {
		t.Errorf("unexpected item: %+v", items[0])
	}
}

func TestRemoveFromCart_NotFound(t *testing.T) {
	repo := NewCartRepository()
	err := repo.RemoveFromCart(1, 100)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	items, _ := repo.GetCartItems(1)
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestClearCart(t *testing.T) {
	repo := NewCartRepository()
	repo.AddToCart(1, 100, 2)
	repo.AddToCart(1, 101, 3)

	err := repo.ClearCart(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	items, _ := repo.GetCartItems(1)
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
