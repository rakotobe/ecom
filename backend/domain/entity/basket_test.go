package entity

import (
	"ecom-backend/domain/value"
	"testing"
)

func TestNewBasket(t *testing.T) {
	basket := NewBasket()

	if basket.ID() == "" {
		t.Error("expected basket to have an ID")
	}
	if !basket.IsEmpty() {
		t.Error("expected new basket to be empty")
	}
}

func TestBasket_AddItem(t *testing.T) {
	basket := NewBasket()
	price, _ := value.NewMoney(1000, "USD")
	qty, _ := value.NewQuantity(2)

	err := basket.AddItem("product-1", qty, price)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if basket.IsEmpty() {
		t.Error("expected basket not to be empty")
	}
	if len(basket.Items()) != 1 {
		t.Errorf("expected 1 item, got %d", len(basket.Items()))
	}
	if basket.ItemCount() != 2 {
		t.Errorf("expected item count 2, got %d", basket.ItemCount())
	}
}

func TestBasket_AddItem_IncrementsQuantity(t *testing.T) {
	basket := NewBasket()
	price, _ := value.NewMoney(1000, "USD")
	qty1, _ := value.NewQuantity(2)
	qty2, _ := value.NewQuantity(3)

	basket.AddItem("product-1", qty1, price)
	basket.AddItem("product-1", qty2, price)

	if len(basket.Items()) != 1 {
		t.Errorf("expected 1 unique item, got %d", len(basket.Items()))
	}
	if basket.ItemCount() != 5 {
		t.Errorf("expected total quantity 5, got %d", basket.ItemCount())
	}
}

func TestBasket_RemoveItem(t *testing.T) {
	basket := NewBasket()
	price, _ := value.NewMoney(1000, "USD")
	qty, _ := value.NewQuantity(2)

	basket.AddItem("product-1", qty, price)
	err := basket.RemoveItem("product-1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !basket.IsEmpty() {
		t.Error("expected basket to be empty")
	}

	err = basket.RemoveItem("non-existent")
	if err == nil {
		t.Error("expected error when removing non-existent item")
	}
}

func TestBasket_Total(t *testing.T) {
	basket := NewBasket()
	price1, _ := value.NewMoney(1000, "USD")
	price2, _ := value.NewMoney(1500, "USD")
	qty1, _ := value.NewQuantity(2)
	qty2, _ := value.NewQuantity(1)

	basket.AddItem("product-1", qty1, price1)
	basket.AddItem("product-2", qty2, price2)

	total, err := basket.Total()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedTotal := int64(2000 + 1500)
	if total.Amount() != expectedTotal {
		t.Errorf("expected total %d, got %d", expectedTotal, total.Amount())
	}
}

func TestBasket_Clear(t *testing.T) {
	basket := NewBasket()
	price, _ := value.NewMoney(1000, "USD")
	qty, _ := value.NewQuantity(2)

	basket.AddItem("product-1", qty, price)
	basket.Clear()

	if !basket.IsEmpty() {
		t.Error("expected basket to be empty after clear")
	}
}
