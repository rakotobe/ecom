package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Products table
		`CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price_amount BIGINT NOT NULL,
			price_currency VARCHAR(3) NOT NULL,
			stock INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,

		// Baskets table
		`CREATE TABLE IF NOT EXISTS baskets (
			id VARCHAR(36) PRIMARY KEY,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,

		// Basket items table
		`CREATE TABLE IF NOT EXISTS basket_items (
			id SERIAL PRIMARY KEY,
			basket_id VARCHAR(36) NOT NULL REFERENCES baskets(id) ON DELETE CASCADE,
			product_id VARCHAR(36) NOT NULL,
			quantity INTEGER NOT NULL,
			price_amount BIGINT NOT NULL,
			price_currency VARCHAR(3) NOT NULL,
			UNIQUE(basket_id, product_id)
		)`,

		// Orders table
		`CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(36) PRIMARY KEY,
			total_amount BIGINT NOT NULL,
			total_currency VARCHAR(3) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)`,

		// Order items table
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id VARCHAR(36) NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			product_id VARCHAR(36) NOT NULL,
			quantity INTEGER NOT NULL,
			price_amount BIGINT NOT NULL,
			price_currency VARCHAR(3) NOT NULL
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_basket_items_basket_id ON basket_items(basket_id)`,
		`CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}

	return nil
}
