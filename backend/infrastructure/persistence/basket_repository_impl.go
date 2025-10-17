package persistence

import (
	"context"
	"database/sql"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"ecom-backend/domain/value"
	"errors"
)

// BasketRepositoryImpl implements BasketRepository using PostgreSQL
type BasketRepositoryImpl struct {
	db *sql.DB
}

// NewBasketRepository creates a new BasketRepositoryImpl
func NewBasketRepository(db *sql.DB) repository.BasketRepository {
	return &BasketRepositoryImpl{db: db}
}

// Save persists a new basket
func (r *BasketRepositoryImpl) Save(ctx context.Context, basket *entity.Basket) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert basket
	query := `INSERT INTO baskets (id, created_at, updated_at) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, query, basket.ID(), basket.CreatedAt(), basket.UpdatedAt())
	if err != nil {
		return err
	}

	// Insert basket items
	if err := r.saveBasketItems(ctx, tx, basket); err != nil {
		return err
	}

	return tx.Commit()
}

// FindByID retrieves a basket by ID
func (r *BasketRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Basket, error) {
	// Get basket
	query := `SELECT id, created_at, updated_at FROM baskets WHERE id = $1`

	var basketID string
	var createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(&basketID, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("basket not found")
		}
		return nil, err
	}

	// Get basket items
	items, err := r.findBasketItems(ctx, basketID)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructBasket(basketID, items, createdAt.Time, updatedAt.Time), nil
}

// Update updates an existing basket
func (r *BasketRepositoryImpl) Update(ctx context.Context, basket *entity.Basket) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update basket
	query := `UPDATE baskets SET updated_at = $2 WHERE id = $1`
	result, err := tx.ExecContext(ctx, query, basket.ID(), basket.UpdatedAt())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("basket not found")
	}

	// Delete existing items
	deleteQuery := `DELETE FROM basket_items WHERE basket_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, basket.ID())
	if err != nil {
		return err
	}

	// Insert updated items
	if err := r.saveBasketItems(ctx, tx, basket); err != nil {
		return err
	}

	return tx.Commit()
}

// Delete removes a basket
func (r *BasketRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM baskets WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("basket not found")
	}

	return nil
}

// ExistsByID checks if a basket exists
func (r *BasketRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM baskets WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)

	return exists, err
}

// saveBasketItems saves basket items within a transaction
func (r *BasketRepositoryImpl) saveBasketItems(ctx context.Context, tx *sql.Tx, basket *entity.Basket) error {
	if len(basket.Items()) == 0 {
		return nil
	}

	query := `
		INSERT INTO basket_items (basket_id, product_id, quantity, price_amount, price_currency)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range basket.Items() {
		_, err := tx.ExecContext(ctx, query,
			basket.ID(),
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

// findBasketItems retrieves basket items
func (r *BasketRepositoryImpl) findBasketItems(ctx context.Context, basketID string) ([]*entity.BasketItem, error) {
	query := `
		SELECT product_id, quantity, price_amount, price_currency
		FROM basket_items
		WHERE basket_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, basketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.BasketItem, 0)

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

		item, err := entity.NewBasketItem(productID, qty, price)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, rows.Err()
}
