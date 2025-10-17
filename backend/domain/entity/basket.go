package entity

import (
	"ecom-backend/domain/value"
	"errors"
	"time"

	"github.com/google/uuid"
)

// BasketItem represents an item in the basket
type BasketItem struct {
	productID string
	quantity  *value.Quantity
	price     *value.Money // price at the time of adding to basket
}

// NewBasketItem creates a new basket item
func NewBasketItem(productID string, quantity *value.Quantity, price *value.Money) (*BasketItem, error) {
	if productID == "" {
		return nil, errors.New("product ID cannot be empty")
	}
	if quantity == nil || quantity.IsZero() {
		return nil, errors.New("quantity must be greater than zero")
	}
	if price == nil {
		return nil, errors.New("price cannot be nil")
	}

	return &BasketItem{
		productID: productID,
		quantity:  quantity,
		price:     price,
	}, nil
}

// ProductID returns the product ID
func (bi *BasketItem) ProductID() string {
	return bi.productID
}

// Quantity returns the quantity
func (bi *BasketItem) Quantity() *value.Quantity {
	return bi.quantity
}

// Price returns the price
func (bi *BasketItem) Price() *value.Money {
	return bi.price
}

// Subtotal calculates the subtotal for this item
func (bi *BasketItem) Subtotal() (*value.Money, error) {
	return bi.price.Multiply(bi.quantity.Value())
}

// Basket represents a shopping basket
type Basket struct {
	id        string
	items     []*BasketItem
	createdAt time.Time
	updatedAt time.Time
}

// NewBasket creates a new empty basket
func NewBasket() *Basket {
	now := time.Now()
	return &Basket{
		id:        uuid.New().String(),
		items:     make([]*BasketItem, 0),
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructBasket reconstructs a Basket from persistence
func ReconstructBasket(id string, items []*BasketItem, createdAt, updatedAt time.Time) *Basket {
	return &Basket{
		id:        id,
		items:     items,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the basket ID
func (b *Basket) ID() string {
	return b.id
}

// Items returns the basket items
func (b *Basket) Items() []*BasketItem {
	return b.items
}

// CreatedAt returns the creation time
func (b *Basket) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns the last update time
func (b *Basket) UpdatedAt() time.Time {
	return b.updatedAt
}

// AddItem adds an item to the basket or updates quantity if item already exists
func (b *Basket) AddItem(productID string, quantity *value.Quantity, price *value.Money) error {
	// Check if item already exists
	for i, item := range b.items {
		if item.productID == productID {
			// Update quantity
			newQuantity, err := item.quantity.Add(quantity)
			if err != nil {
				return err
			}
			newItem, err := NewBasketItem(productID, newQuantity, price)
			if err != nil {
				return err
			}
			b.items[i] = newItem
			b.updatedAt = time.Now()
			return nil
		}
	}

	// Add new item
	item, err := NewBasketItem(productID, quantity, price)
	if err != nil {
		return err
	}
	b.items = append(b.items, item)
	b.updatedAt = time.Now()
	return nil
}

// RemoveItem removes an item from the basket
func (b *Basket) RemoveItem(productID string) error {
	for i, item := range b.items {
		if item.productID == productID {
			b.items = append(b.items[:i], b.items[i+1:]...)
			b.updatedAt = time.Now()
			return nil
		}
	}
	return errors.New("item not found in basket")
}

// UpdateItemQuantity updates the quantity of an item
func (b *Basket) UpdateItemQuantity(productID string, quantity *value.Quantity) error {
	if quantity.IsZero() {
		return b.RemoveItem(productID)
	}

	for i, item := range b.items {
		if item.productID == productID {
			newItem, err := NewBasketItem(productID, quantity, item.price)
			if err != nil {
				return err
			}
			b.items[i] = newItem
			b.updatedAt = time.Now()
			return nil
		}
	}
	return errors.New("item not found in basket")
}

// Clear removes all items from the basket
func (b *Basket) Clear() {
	b.items = make([]*BasketItem, 0)
	b.updatedAt = time.Now()
}

// IsEmpty checks if the basket is empty
func (b *Basket) IsEmpty() bool {
	return len(b.items) == 0
}

// Total calculates the total price of all items in the basket
func (b *Basket) Total() (*value.Money, error) {
	if b.IsEmpty() {
		return value.NewMoney(0, "USD")
	}

	total, err := value.NewMoney(0, b.items[0].price.Currency())
	if err != nil {
		return nil, err
	}

	for _, item := range b.items {
		subtotal, err := item.Subtotal()
		if err != nil {
			return nil, err
		}
		total, err = total.Add(subtotal)
		if err != nil {
			return nil, err
		}
	}

	return total, nil
}

// ItemCount returns the total number of items in the basket
func (b *Basket) ItemCount() int {
	count := 0
	for _, item := range b.items {
		count += item.quantity.Value()
	}
	return count
}
