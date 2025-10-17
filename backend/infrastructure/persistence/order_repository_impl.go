package persistence

import (
	"context"
	"database/sql"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"ecom-backend/domain/value"
	"errors"
)

// OrderRepositoryImpl implements OrderRepository using PostgreSQL
type OrderRepositoryImpl struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepositoryImpl
func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

// Save persists a new order
func (r *OrderRepositoryImpl) Save(ctx context.Context, order *entity.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	query := `
		INSERT INTO orders (id, total_amount, total_currency, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(ctx, query,
		order.ID(),
		order.Total().Amount(),
		order.Total().Currency(),
		string(order.Status()),
		order.CreatedAt(),
		order.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	// Insert order items
	if err := r.saveOrderItems(ctx, tx, order); err != nil {
		return err
	}

	return tx.Commit()
}

// FindByID retrieves an order by ID
func (r *OrderRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	// Get order
	query := `
		SELECT id, total_amount, total_currency, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var orderID, currency, status string
	var totalAmount int64
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&orderID, &totalAmount, &currency, &status, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Get order items
	items, err := r.findOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	total, err := value.NewMoney(totalAmount, currency)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructOrder(
		orderID, items, total, entity.OrderStatus(status),
		createdAt.Time, updatedAt.Time,
	), nil
}

// FindAll retrieves all orders
func (r *OrderRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Order, error) {
	query := `
		SELECT id, total_amount, total_currency, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)

	for rows.Next() {
		var orderID, currency, status string
		var totalAmount int64
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(&orderID, &totalAmount, &currency, &status, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		items, err := r.findOrderItems(ctx, orderID)
		if err != nil {
			return nil, err
		}

		total, err := value.NewMoney(totalAmount, currency)
		if err != nil {
			return nil, err
		}

		order := entity.ReconstructOrder(
			orderID, items, total, entity.OrderStatus(status),
			createdAt.Time, updatedAt.Time,
		)

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// Update updates an existing order
func (r *OrderRepositoryImpl) Update(ctx context.Context, order *entity.Order) error {
	query := `
		UPDATE orders
		SET total_amount = $2, total_currency = $3, status = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		order.ID(),
		order.Total().Amount(),
		order.Total().Currency(),
		string(order.Status()),
		order.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil
}

// ExistsByID checks if an order exists
func (r *OrderRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)

	return exists, err
}

// saveOrderItems saves order items within a transaction
func (r *OrderRepositoryImpl) saveOrderItems(ctx context.Context, tx *sql.Tx, order *entity.Order) error {
	if len(order.Items()) == 0 {
		return nil
	}

	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price_amount, price_currency)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items() {
		_, err := tx.ExecContext(ctx, query,
			order.ID(),
			item.ProductID(),
			item.Quantity().Value(),
			item.Price().Amount(),
			item.Price().Currency(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// findOrderItems retrieves order items
func (r *OrderRepositoryImpl) findOrderItems(ctx context.Context, orderID string) ([]*entity.OrderItem, error) {
	query := `
		SELECT product_id, quantity, price_amount, price_currency
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.OrderItem, 0)

	for rows.Next() {
		var productID, currency string
		var quantity int
		var priceAmount int64

		if err := rows.Scan(&productID, &quantity, &priceAmount, &currency); err != nil {
			return nil, err
		}

		price, err := value.NewMoney(priceAmount, currency)
		if err != nil {
			return nil, err
		}

		qty, err := value.NewQuantity(quantity)
		if err != nil {
			return nil, err
		}

		item, err := entity.NewOrderItem(productID, qty, price)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
