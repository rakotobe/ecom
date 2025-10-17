package service

import (
	"context"
	"ecom-backend/application/dto"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"ecom-backend/domain/value"
	"errors"
)

// BasketService handles basket-related business logic
type BasketService struct {
	basketRepo  repository.BasketRepository
	productRepo repository.ProductRepository
}

// NewBasketService creates a new BasketService
func NewBasketService(basketRepo repository.BasketRepository, productRepo repository.ProductRepository) *BasketService {
	return &BasketService{
		basketRepo:  basketRepo,
		productRepo: productRepo,
	}
}

// CreateBasket creates a new empty basket
func (s *BasketService) CreateBasket(ctx context.Context) (*dto.BasketResponse, error) {
	basket := entity.NewBasket()

	if err := s.basketRepo.Save(ctx, basket); err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// GetBasket retrieves a basket by ID
func (s *BasketService) GetBasket(ctx context.Context, id string) (*dto.BasketResponse, error) {
	basket, err := s.basketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// AddItem adds an item to the basket
func (s *BasketService) AddItem(ctx context.Context, basketID string, req *dto.AddItemRequest) (*dto.BasketResponse, error) {
	// Validate request
	if req.ProductID == "" {
		return nil, errors.New("product ID is required")
	}
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	// Retrieve basket
	basket, err := s.basketRepo.FindByID(ctx, basketID)
	if err != nil {
		return nil, err
	}

	// Retrieve product to get current price and verify availability
	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	// Check if product has sufficient stock
	requestedQty, err := value.NewQuantity(req.Quantity)
	if err != nil {
		return nil, err
	}

	if product.Stock().Value() < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Add item to basket
	if err := basket.AddItem(product.ID(), requestedQty, product.Price()); err != nil {
		return nil, err
	}

	// Persist
	if err := s.basketRepo.Update(ctx, basket); err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// RemoveItem removes an item from the basket
func (s *BasketService) RemoveItem(ctx context.Context, basketID, productID string) (*dto.BasketResponse, error) {
	basket, err := s.basketRepo.FindByID(ctx, basketID)
	if err != nil {
		return nil, err
	}

	if err := basket.RemoveItem(productID); err != nil {
		return nil, err
	}

	if err := s.basketRepo.Update(ctx, basket); err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// UpdateItemQuantity updates the quantity of an item in the basket
func (s *BasketService) UpdateItemQuantity(ctx context.Context, basketID, productID string, req *dto.UpdateItemQuantityRequest) (*dto.BasketResponse, error) {
	if req.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	basket, err := s.basketRepo.FindByID(ctx, basketID)
	if err != nil {
		return nil, err
	}

	// If quantity is 0, remove the item
	if req.Quantity == 0 {
		if err := basket.RemoveItem(productID); err != nil {
			return nil, err
		}
	} else {
		// Verify product availability
		product, err := s.productRepo.FindByID(ctx, productID)
		if err != nil {
			return nil, err
		}

		if product.Stock().Value() < req.Quantity {
			return nil, errors.New("insufficient stock")
		}

		quantity, err := value.NewQuantity(req.Quantity)
		if err != nil {
			return nil, err
		}

		if err := basket.UpdateItemQuantity(productID, quantity); err != nil {
			return nil, err
		}
	}

	if err := s.basketRepo.Update(ctx, basket); err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// ClearBasket removes all items from the basket
func (s *BasketService) ClearBasket(ctx context.Context, basketID string) (*dto.BasketResponse, error) {
	basket, err := s.basketRepo.FindByID(ctx, basketID)
	if err != nil {
		return nil, err
	}

	basket.Clear()

	if err := s.basketRepo.Update(ctx, basket); err != nil {
		return nil, err
	}

	return s.toBasketResponse(basket)
}

// toBasketResponse converts a Basket entity to BasketResponse DTO
func (s *BasketService) toBasketResponse(basket *entity.Basket) (*dto.BasketResponse, error) {
	items := make([]dto.BasketItemResponse, 0, len(basket.Items()))

	for _, item := range basket.Items() {
		subtotal, err := item.Subtotal()
		if err != nil {
			return nil, err
		}

		items = append(items, dto.BasketItemResponse{
			ProductID: item.ProductID(),
			Quantity:  item.Quantity().Value(),
			Price:     item.Price().Amount(),
			Currency:  item.Price().Currency(),
			Subtotal:  subtotal.Amount(),
		})
	}

	total, err := basket.Total()
	if err != nil {
		return nil, err
	}

	currency := "USD"
	if len(basket.Items()) > 0 {
		currency = basket.Items()[0].Price().Currency()
	}

	return &dto.BasketResponse{
		ID:        basket.ID(),
		Items:     items,
		Total:     total.Amount(),
		Currency:  currency,
		ItemCount: basket.ItemCount(),
		CreatedAt: basket.CreatedAt(),
		UpdatedAt: basket.UpdatedAt(),
	}, nil
}
