package service

import (
	"context"
	"ecom-backend/application/dto"
	"ecom-backend/domain/entity"
	"ecom-backend/domain/repository"
	"errors"
)

// OrderService handles order-related business logic
type OrderService struct {
	orderRepo   repository.OrderRepository
	basketRepo  repository.BasketRepository
	productRepo repository.ProductRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo repository.OrderRepository, basketRepo repository.BasketRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		basketRepo:  basketRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates an order from a basket (checkout)
func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if req.BasketID == "" {
		return nil, errors.New("basket ID is required")
	}

	// Retrieve basket
	basket, err := s.basketRepo.FindByID(ctx, req.BasketID)
	if err != nil {
		return nil, err
	}

	if basket.IsEmpty() {
		return nil, errors.New("cannot create order from empty basket")
	}

	// Verify stock availability for all items
	for _, item := range basket.Items() {
		product, err := s.productRepo.FindByID(ctx, item.ProductID())
		if err != nil {
			return nil, err
		}

		if product.Stock().Value() < item.Quantity().Value() {
			return nil, errors.New("insufficient stock for product: " + product.Name())
		}
	}

	// Reduce stock for all items
	for _, item := range basket.Items() {
		product, err := s.productRepo.FindByID(ctx, item.ProductID())
		if err != nil {
			return nil, err
		}

		if err := product.ReduceStock(item.Quantity()); err != nil {
			return nil, err
		}

		if err := s.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}
	}

	// Create order
	order, err := entity.NewOrder(basket.Items())
	if err != nil {
		return nil, err
	}

	// Persist order
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	// Clear basket after successful order
	basket.Clear()
	if err := s.basketRepo.Update(ctx, basket); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders(ctx context.Context) ([]*dto.OrderResponse, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, s.toOrderResponse(order))
	}

	return responses, nil
}

// ConfirmOrder confirms a pending order
func (s *OrderService) ConfirmOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Confirm(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// ShipOrder marks an order as shipped
func (s *OrderService) ShipOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Ship(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// DeliverOrder marks an order as delivered
func (s *OrderService) DeliverOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Deliver(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := order.Cancel(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return s.toOrderResponse(order), nil
}

// toOrderResponse converts an Order entity to OrderResponse DTO
func (s *OrderService) toOrderResponse(order *entity.Order) *dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items()))

	for _, item := range order.Items() {
		subtotal, _ := item.Subtotal()

		items = append(items, dto.OrderItemResponse{
			ProductID: item.ProductID(),
			Quantity:  item.Quantity().Value(),
			Price:     item.Price().Amount(),
			Currency:  item.Price().Currency(),
			Subtotal:  subtotal.Amount(),
		})
	}

	return &dto.OrderResponse{
		ID:        order.ID(),
		Items:     items,
		Total:     order.Total().Amount(),
		Currency:  order.Total().Currency(),
		Status:    string(order.Status()),
		CreatedAt: order.CreatedAt(),
		UpdatedAt: order.UpdatedAt(),
	}
}
