package repository

import (
	"context"
	"ecom-backend/domain/entity"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	// Save persists an order
	Save(ctx context.Context, order *entity.Order) error

	// FindByID retrieves an order by ID
	FindByID(ctx context.Context, id string) (*entity.Order, error)

	// FindAll retrieves all orders
	FindAll(ctx context.Context) ([]*entity.Order, error)

	// Update updates an existing order
	Update(ctx context.Context, order *entity.Order) error

	// ExistsByID checks if an order exists
	ExistsByID(ctx context.Context, id string) (bool, error)
}
