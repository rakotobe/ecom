package entity

import (
	"ecom-backend/domain/value"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the catalog
type Product struct {
	id          string
	name        string
	description string
	price       *value.Money
	stock       *value.Quantity
	createdAt   time.Time
	updatedAt   time.Time
}

// NewProduct creates a new Product entity
func NewProduct(name, description string, price *value.Money, stock *value.Quantity) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price == nil {
		return nil, errors.New("product price cannot be nil")
	}
	if stock == nil {
		return nil, errors.New("product stock cannot be nil")
	}

	now := time.Now()
	return &Product{
		id:          uuid.New().String(),
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ReconstructProduct reconstructs a Product from persistence
func ReconstructProduct(id, name, description string, price *value.Money, stock *value.Quantity, createdAt, updatedAt time.Time) *Product {
	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the product ID
func (p *Product) ID() string {
	return p.id
}

// Name returns the product name
func (p *Product) Name() string {
	return p.name
}

// Description returns the product description
func (p *Product) Description() string {
	return p.description
}

// Price returns the product price
func (p *Product) Price() *value.Money {
	return p.price
}

// Stock returns the product stock
func (p *Product) Stock() *value.Quantity {
	return p.stock
}

// CreatedAt returns the creation time
func (p *Product) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update time
func (p *Product) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateDetails updates product details
func (p *Product) UpdateDetails(name, description string, price *value.Money) error {
	if name == "" {
		return errors.New("product name cannot be empty")
	}
	if price == nil {
		return errors.New("product price cannot be nil")
	}

	p.name = name
	p.description = description
	p.price = price
	p.updatedAt = time.Now()
	return nil
}

// UpdateStock updates the product stock
func (p *Product) UpdateStock(stock *value.Quantity) error {
	if stock == nil {
		return errors.New("product stock cannot be nil")
	}
	p.stock = stock
	p.updatedAt = time.Now()
	return nil
}

// ReduceStock reduces stock by the given quantity
func (p *Product) ReduceStock(quantity *value.Quantity) error {
	newStock, err := p.stock.Subtract(quantity)
	if err != nil {
		return errors.New("insufficient stock")
	}
	p.stock = newStock
	p.updatedAt = time.Now()
	return nil
}

// IsAvailable checks if the product has stock
func (p *Product) IsAvailable() bool {
	return !p.stock.IsZero()
}
