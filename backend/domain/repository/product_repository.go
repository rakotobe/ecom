package repository

import (
	"context"
	"ecom-backend/domain/entity"
)

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Save persists a product
	Save(ctx context.Context, product *entity.Product) error

	// FindByID retrieves a product by ID
	FindByID(ctx context.Context, id string) (*entity.Product, error)

	// FindAll retrieves all products
	FindAll(ctx context.Context) ([]*entity.Product, error)

	// Update updates an existing product
	Update(ctx context.Context, product *entity.Product) error

	// Delete removes a product
	Delete(ctx context.Context, id string) error

	// ExistsByID checks if a product exists
	ExistsByID(ctx context.Context, id string) (bool, error)
}
