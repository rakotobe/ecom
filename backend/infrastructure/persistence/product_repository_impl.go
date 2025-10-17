package persistence

import (
	"context"
	"database/sql"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"ecom-backend/domain/value"
	"errors"
)

// ProductRepositoryImpl implements ProductRepository using PostgreSQL
type ProductRepositoryImpl struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepositoryImpl
func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

// Save persists a new product
func (r *ProductRepositoryImpl) Save(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (id, name, description, price_amount, price_currency, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		product.ID(),
		product.Name(),
		product.Description(),
		product.Price().Amount(),
		product.Price().Currency(),
		product.Stock().Value(),
		product.CreatedAt(),
		product.UpdatedAt(),
	)

	return err
}

// FindByID retrieves a product by ID
func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `
		SELECT id, name, description, price_amount, price_currency, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var (
		productID, name, description, currency string
		priceAmount                            int64
		stock                                  int
		createdAt, updatedAt                   sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&productID, &name, &description, &priceAmount, &currency, &stock, &createdAt, &updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	price, err := value.NewMoney(priceAmount, currency)
	if err != nil {
		return nil, err
	}

	stockQty, err := value.NewQuantity(stock)
	if err != nil {
		return nil, err
	}

	return entity.ReconstructProduct(
		productID, name, description, price, stockQty,
		createdAt.Time, updatedAt.Time,
	), nil
}

// FindAll retrieves all products
func (r *ProductRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Product, error) {
	query := `
		SELECT id, name, description, price_amount, price_currency, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*entity.Product, 0)

	for rows.Next() {
		var (
			productID, name, description, currency string
			priceAmount                            int64
			stock                                  int
			createdAt, updatedAt                   sql.NullTime
		)

		if err := rows.Scan(
			&productID, &name, &description, &priceAmount, &currency, &stock, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}

		price, err := value.NewMoney(priceAmount, currency)
		if err != nil {
			return nil, err
		}

		stockQty, err := value.NewQuantity(stock)
		if err != nil {
			return nil, err
		}

		product := entity.ReconstructProduct(
			productID, name, description, price, stockQty,
			createdAt.Time, updatedAt.Time,
		)

		products = append(products, product)
	}

	return products, rows.Err()
}

// Update updates an existing product
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price_amount = $4, price_currency = $5, stock = $6, updated_at = $7
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		product.ID(),
		product.Name(),
		product.Description(),
		product.Price().Amount(),
		product.Price().Currency(),
		product.Stock().Value(),
		product.UpdatedAt(),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// Delete removes a product
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

// ExistsByID checks if a product exists
func (r *ProductRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)

	return exists, err
}
