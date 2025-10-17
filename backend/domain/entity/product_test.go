package entity

import (
	"ecom-backend/domain/value"
	"testing"
)

func TestNewProduct(t *testing.T) {
	price, _ := value.NewMoney(1000, "USD")
	stock, _ := value.NewQuantity(10)

	tests := []struct {
		name        string
		productName string
		description string
		price       *value.Money
		stock       *value.Quantity
		wantError   bool
	}{
		{"valid product", "Test Product", "Description", price, stock, false},
		{"empty name", "", "Description", price, stock, true},
		{"nil price", "Test Product", "Description", nil, stock, true},
		{"nil stock", "Test Product", "Description", price, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.description, tt.price, tt.stock)
			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if product.Name() != tt.productName {
					t.Errorf("expected name %s, got %s", tt.productName, product.Name())
				}
			}
		})
	}
}

func TestProduct_UpdateDetails(t *testing.T) {
	price, _ := value.NewMoney(1000, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := NewProduct("Original", "Original desc", price, stock)

	newPrice, _ := value.NewMoney(1500, "USD")
	err := product.UpdateDetails("Updated", "Updated desc", newPrice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if product.Name() != "Updated" {
		t.Errorf("expected name 'Updated', got %s", product.Name())
	}
	if product.Price().Amount() != 1500 {
		t.Errorf("expected price 1500, got %d", product.Price().Amount())
	}
}

func TestProduct_ReduceStock(t *testing.T) {
	price, _ := value.NewMoney(1000, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := NewProduct("Test", "Test", price, stock)

	reduceBy, _ := value.NewQuantity(3)
	err := product.ReduceStock(reduceBy)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if product.Stock().Value() != 7 {
		t.Errorf("expected stock 7, got %d", product.Stock().Value())
	}

	// Try to reduce more than available
	reduceBy, _ = value.NewQuantity(10)
	err = product.ReduceStock(reduceBy)
	if err == nil {
		t.Error("expected error when reducing more than available stock")
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	price, _ := value.NewMoney(1000, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := NewProduct("Test", "Test", price, stock)

	if !product.IsAvailable() {
		t.Error("expected product to be available")
	}

	zeroStock, _ := value.NewQuantity(0)
	product.UpdateStock(zeroStock)

	if product.IsAvailable() {
		t.Error("expected product not to be available")
	}
}
