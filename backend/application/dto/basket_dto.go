package dto

import "time"

// AddItemRequest represents the request to add an item to basket
type AddItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// UpdateItemQuantityRequest represents the request to update item quantity
type UpdateItemQuantityRequest struct {
	Quantity int `json:"quantity"`
}

// BasketItemResponse represents a basket item in responses
type BasketItemResponse struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int64  `json:"price"`     // price in cents
	Currency  string `json:"currency"`
	Subtotal  int64  `json:"subtotal"`  // subtotal in cents
}

// BasketResponse represents a basket in responses
type BasketResponse struct {
	ID        string               `json:"id"`
	Items     []BasketItemResponse `json:"items"`
	Total     int64                `json:"total"`     // total in cents
	Currency  string               `json:"currency"`
	ItemCount int                  `json:"item_count"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
