package usecase

import (
	"backend/internal/repository"
	"testing"
)

func TestGetAllProducts(t *testing.T) {
	repo := repository.NewProductRepository()
	uc := NewProductUseCase(repo)

	products := uc.GetAllProducts()
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	if products[0].Name != "Apple" {
		t.Errorf("expected first product to be 'Apple', got '%s'", products[0].Name)
	}
	if products[1].Name != "Banana" {
		t.Errorf("expected second product to be 'Banana', got '%s'", products[1].Name)
	}
}

func TestGetProductById_Found(t *testing.T) {
	repo := repository.NewProductRepository()
	uc := NewProductUseCase(repo)

	product := uc.GetProductById(1)
	if product.ID == 0 {
		t.Fatalf("expected to find product with ID 1, got zero value")
	}

	if product.Name != "Apple" {
		t.Errorf("expected product name 'Apple', got '%s'", product.Name)
	}
}

func TestGetProductById_NotFound(t *testing.T) {
	repo := repository.NewProductRepository()
	uc := NewProductUseCase(repo)

	product := uc.GetProductById(999)
	if product.ID != 0 {
		t.Errorf("expected zero value product for non-existing ID, got %+v", product)
	}
}
