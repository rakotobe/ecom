package service

import (
	"context"
	"ecom-backend/application/dto"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/value"
	"errors"
	"testing"
)

// Mock repository for service testing
type mockProductRepo struct {
	products  map[string]*entity.Product
	saveErr   error
	findErr   error
	updateErr error
	deleteErr error
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[string]*entity.Product),
	}
}

func (m *mockProductRepo) Save(ctx context.Context, product *entity.Product) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.products[product.ID()] = product
	return nil
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	product, ok := m.products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (m *mockProductRepo) FindAll(ctx context.Context) ([]*entity.Product, error) {
	products := make([]*entity.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *mockProductRepo) Update(ctx context.Context, product *entity.Product) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.products[product.ID()]; !ok {
		return errors.New("product not found")
	}
	m.products[product.ID()] = product
	return nil
}

func (m *mockProductRepo) Delete(ctx context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.products[id]; !ok {
		return errors.New("product not found")
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	_, ok := m.products[id]
	return ok, nil
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	t.Run("Valid product creation", func(t *testing.T) {
		req := &dto.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       1999,
			Currency:    "USD",
			Stock:       10,
		}

		response, err := service.CreateProduct(ctx, req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if response == nil {
			t.Fatal("Expected response, got nil")
		}
		if response.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, response.Name)
		}
		if response.Price != req.Price {
			t.Errorf("Expected price %d, got %d", req.Price, response.Price)
		}
	})

	t.Run("Invalid product - empty name", func(t *testing.T) {
		req := &dto.CreateProductRequest{
			Name:     "",
			Price:    1999,
			Currency: "USD",
			Stock:    10,
		}

		_, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error for empty name, got nil")
		}
	})

	t.Run("Invalid product - negative price", func(t *testing.T) {
		req := &dto.CreateProductRequest{
			Name:     "Test",
			Price:    -100,
			Currency: "USD",
			Stock:    10,
		}

		_, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error for negative price, got nil")
		}
	})

	t.Run("Repository save error", func(t *testing.T) {
		repo.saveErr = errors.New("database error")
		defer func() { repo.saveErr = nil }()

		req := &dto.CreateProductRequest{
			Name:     "Test",
			Price:    1999,
			Currency: "USD",
			Stock:    10,
		}

		_, err := service.CreateProduct(ctx, req)

		if err == nil {
			t.Error("Expected error from repository, got nil")
		}
	})
}

func TestProductService_GetProduct(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	// Create a test product
	price, _ := value.NewMoney(1999, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := entity.NewProduct("Test Product", "Description", price, stock)
	repo.Save(ctx, product)

	t.Run("Get existing product", func(t *testing.T) {
		response, err := service.GetProduct(ctx, product.ID())

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if response.ID != product.ID() {
			t.Errorf("Expected ID %s, got %s", product.ID(), response.ID)
		}
		if response.Name != product.Name() {
			t.Errorf("Expected name %s, got %s", product.Name(), response.Name)
		}
	})

	t.Run("Get non-existent product", func(t *testing.T) {
		_, err := service.GetProduct(ctx, "non-existent-id")

		if err == nil {
			t.Error("Expected error for non-existent product, got nil")
		}
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	// Create a test product
	price, _ := value.NewMoney(1999, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := entity.NewProduct("Original Name", "Original Description", price, stock)
	repo.Save(ctx, product)

	t.Run("Valid update", func(t *testing.T) {
		req := &dto.UpdateProductRequest{
			Name:        "Updated Name",
			Description: "Updated Description",
			Price:       2499,
			Currency:    "USD",
		}

		response, err := service.UpdateProduct(ctx, product.ID(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if response.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, response.Name)
		}
		if response.Price != req.Price {
			t.Errorf("Expected price %d, got %d", req.Price, response.Price)
		}
	})

	t.Run("Update non-existent product", func(t *testing.T) {
		req := &dto.UpdateProductRequest{
			Name:     "Test",
			Price:    1999,
			Currency: "USD",
		}

		_, err := service.UpdateProduct(ctx, "non-existent-id", req)

		if err == nil {
			t.Error("Expected error for non-existent product, got nil")
		}
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	// Create a test product
	price, _ := value.NewMoney(1999, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := entity.NewProduct("To Delete", "Description", price, stock)
	repo.Save(ctx, product)

	t.Run("Delete existing product", func(t *testing.T) {
		err := service.DeleteProduct(ctx, product.ID())

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify it's deleted
		exists, _ := repo.ExistsByID(ctx, product.ID())
		if exists {
			t.Error("Expected product to be deleted")
		}
	})

	t.Run("Delete non-existent product", func(t *testing.T) {
		err := service.DeleteProduct(ctx, "non-existent-id")

		if err == nil {
			t.Error("Expected error for non-existent product, got nil")
		}
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	// Create test products
	price, _ := value.NewMoney(1999, "USD")
	stock, _ := value.NewQuantity(10)
	product1, _ := entity.NewProduct("Product 1", "Description 1", price, stock)
	product2, _ := entity.NewProduct("Product 2", "Description 2", price, stock)
	repo.Save(ctx, product1)
	repo.Save(ctx, product2)

	t.Run("Get all products", func(t *testing.T) {
		products, err := service.GetAllProducts(ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(products) != 2 {
			t.Errorf("Expected 2 products, got %d", len(products))
		}
	})
}

func TestProductService_UpdateStock(t *testing.T) {
	repo := newMockProductRepo()
	service := NewProductService(repo)
	ctx := context.Background()

	// Create a test product
	price, _ := value.NewMoney(1999, "USD")
	stock, _ := value.NewQuantity(10)
	product, _ := entity.NewProduct("Test Product", "Description", price, stock)
	repo.Save(ctx, product)

	t.Run("Valid stock update", func(t *testing.T) {
		req := &dto.UpdateStockRequest{Stock: 20}

		response, err := service.UpdateStock(ctx, product.ID(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if response.Stock != 20 {
			t.Errorf("Expected stock 20, got %d", response.Stock)
		}
	})

	t.Run("Invalid stock - negative", func(t *testing.T) {
		req := &dto.UpdateStockRequest{Stock: -5}

		_, err := service.UpdateStock(ctx, product.ID(), req)

		if err == nil {
			t.Error("Expected error for negative stock, got nil")
		}
	})
}
