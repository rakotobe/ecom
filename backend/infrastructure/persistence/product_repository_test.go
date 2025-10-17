package persistence

import (
	"context"
	"database/sql"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/value"
	"testing"

	_ "github.com/lib/pq"
)

// setupTestDB creates a test database connection
// NOTE: This requires PostgreSQL to be running
// Run with: go test -v ./infrastructure/persistence/
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=ecom_test sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping integration test (DB not available): %v", err)
	}

	// Create tables for testing
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price_amount BIGINT NOT NULL,
			price_currency VARCHAR(3) NOT NULL,
			stock INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	db.Exec("TRUNCATE TABLE products CASCADE")
	db.Close()
}

func TestProductRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := NewProductRepository(db)
	ctx := context.Background()

	t.Run("Save and FindByID", func(t *testing.T) {
		// Arrange
		price, _ := value.NewMoney(1999, "USD")
		stock, _ := value.NewQuantity(10)
		product, _ := entity.NewProduct("Test Product", "Test Description", price, stock)

		// Act - Save
		err := repo.Save(ctx, product)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Act - Find
		found, err := repo.FindByID(ctx, product.ID())
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}

		// Assert
		if found.ID() != product.ID() {
			t.Errorf("Expected ID %s, got %s", product.ID(), found.ID())
		}
		if found.Name() != product.Name() {
			t.Errorf("Expected name %s, got %s", product.Name(), found.Name())
		}
		if found.Price().Amount() != product.Price().Amount() {
			t.Errorf("Expected price %d, got %d", product.Price().Amount(), found.Price().Amount())
		}
	})

	t.Run("Update", func(t *testing.T) {
		// Arrange
		price, _ := value.NewMoney(1999, "USD")
		stock, _ := value.NewQuantity(10)
		product, _ := entity.NewProduct("Original Name", "Original Description", price, stock)
		repo.Save(ctx, product)

		// Act
		newPrice, _ := value.NewMoney(2499, "USD")
		product.UpdateDetails("Updated Name", "Updated Description", newPrice)
		err := repo.Update(ctx, product)

		// Assert
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		found, _ := repo.FindByID(ctx, product.ID())
		if found.Name() != "Updated Name" {
			t.Errorf("Expected updated name, got %s", found.Name())
		}
		if found.Price().Amount() != 2499 {
			t.Errorf("Expected updated price 2499, got %d", found.Price().Amount())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		// Arrange
		price, _ := value.NewMoney(1999, "USD")
		stock, _ := value.NewQuantity(10)
		product, _ := entity.NewProduct("To Delete", "Description", price, stock)
		repo.Save(ctx, product)

		// Act
		err := repo.Delete(ctx, product.ID())

		// Assert
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repo.FindByID(ctx, product.ID())
		if err == nil {
			t.Error("Expected error when finding deleted product")
		}
	})

	t.Run("FindAll", func(t *testing.T) {
		// Clean up first
		db.Exec("TRUNCATE TABLE products CASCADE")

		// Arrange
		price, _ := value.NewMoney(1999, "USD")
		stock, _ := value.NewQuantity(10)
		product1, _ := entity.NewProduct("Product 1", "Description 1", price, stock)
		product2, _ := entity.NewProduct("Product 2", "Description 2", price, stock)

		repo.Save(ctx, product1)
		repo.Save(ctx, product2)

		// Act
		products, err := repo.FindAll(ctx)

		// Assert
		if err != nil {
			t.Fatalf("FindAll failed: %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})

	t.Run("ExistsByID", func(t *testing.T) {
		// Arrange
		price, _ := value.NewMoney(1999, "USD")
		stock, _ := value.NewQuantity(10)
		product, _ := entity.NewProduct("Exists Test", "Description", price, stock)
		repo.Save(ctx, product)

		// Act
		exists, err := repo.ExistsByID(ctx, product.ID())

		// Assert
		if err != nil {
			t.Fatalf("ExistsByID failed: %v", err)
		}
		if !exists {
			t.Error("Expected product to exist")
		}

		// Test non-existent
		exists, err = repo.ExistsByID(ctx, "non-existent-id")
		if err != nil {
			t.Fatalf("ExistsByID failed: %v", err)
		}
		if exists {
			t.Error("Expected product not to exist")
		}
	})
}
