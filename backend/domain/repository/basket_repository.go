package repository

import (
	"context"
	"ecom-backend/domain/entity"
)

// BasketRepository defines the interface for basket persistence
type BasketRepository interface {
	// Save persists a basket
	Save(ctx context.Context, basket *entity.Basket) error

	// FindByID retrieves a basket by ID
	FindByID(ctx context.Context, id string) (*entity.Basket, error)

	// Update updates an existing basket
	Update(ctx context.Context, basket *entity.Basket) error

	// Delete removes a basket
	Delete(ctx context.Context, id string) error

	// ExistsByID checks if a basket exists
	ExistsByID(ctx context.Context, id string) (bool, error)
}
