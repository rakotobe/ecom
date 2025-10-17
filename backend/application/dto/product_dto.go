package dto

import "time"

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`        // price in cents
	Currency    string `json:"currency"`     // e.g., "USD"
	Stock       int    `json:"stock"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`        // price in cents
	Currency    string `json:"currency"`
}

// UpdateStockRequest represents the request to update stock
type UpdateStockRequest struct {
	Stock int `json:"stock"`
}

// ProductResponse represents a product in responses
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`        // price in cents
	Currency    string    `json:"currency"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
