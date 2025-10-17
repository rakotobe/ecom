package service

import (
	"context"
	"ecom-backend/application/dto"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"ecom-backend/domain/value"
	"errors"
)

// ProductService handles product-related business logic
type ProductService struct {
	productRepo repository.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if req.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// Create value objects
	price, err := value.NewMoney(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	stock, err := value.NewQuantity(req.Stock)
	if err != nil {
		return nil, err
	}

	// Create entity
	product, err := entity.NewProduct(req.Name, req.Description, price, stock)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, s.toProductResponse(product))
	}

	return responses, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("product name is required")
	}
	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// Retrieve existing product
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create new price
	price, err := value.NewMoney(req.Price, req.Currency)
	if err != nil {
		return nil, err
	}

	// Update product
	if err := product.UpdateDetails(req.Name, req.Description, price); err != nil {
		return nil, err
	}

	// Persist
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UpdateStock updates product stock
func (s *ProductService) UpdateStock(ctx context.Context, id string, req *dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	if req.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	stock, err := value.NewQuantity(req.Stock)
	if err != nil {
		return nil, err
	}

	if err := product.UpdateStock(stock); err != nil {
		return nil, err
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	exists, err := s.productRepo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(ctx, id)
}

// toProductResponse converts a Product entity to ProductResponse DTO
func (s *ProductService) toProductResponse(product *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID(),
		Name:        product.Name(),
		Description: product.Description(),
		Price:       product.Price().Amount(),
		Currency:    product.Price().Currency(),
		Stock:       product.Stock().Value(),
		CreatedAt:   product.CreatedAt(),
		UpdatedAt:   product.UpdatedAt(),
	}
}
