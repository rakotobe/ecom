package dto

import "time"

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	BasketID string `json:"basket_id"`
}

// OrderItemResponse represents an order item in responses
type OrderItemResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int64  `json:"price"`     // price in cents
	Currency  string `json:"currency"`
	Subtotal  int64  `json:"subtotal"`  // subtotal in cents
}

// OrderResponse represents an order in responses
type OrderResponse struct {
	ID        string              `json:"id"`
	Items     []OrderItemResponse `json:"items"`
	Total     int64               `json:"total"`     // total in cents
	Currency  string              `json:"currency"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}
