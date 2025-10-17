package entity

import (
	"ecom-backend/domain/value"
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusDelivered OrderStatus = "DELIVERED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// OrderItem represents an item in an order
type OrderItem struct {
	productID string
	quantity  *value.Quantity
	price     *value.Money
}

// NewOrderItem creates a new order item
func NewOrderItem(productID string, quantity *value.Quantity, price *value.Money) (*OrderItem, error) {
	if productID == "" {
		return nil, errors.New("product ID cannot be empty")
	}
	if quantity == nil || quantity.IsZero() {
		return nil, errors.New("quantity must be greater than zero")
	}
	if price == nil {
		return nil, errors.New("price cannot be nil")
	}

	return &OrderItem{
		productID: productID,
		quantity:  quantity,
		price:     price,
	}, nil
}

// ProductID returns the product ID
func (oi *OrderItem) ProductID() string {
	return oi.productID
}

// Quantity returns the quantity
func (oi *OrderItem) Quantity() *value.Quantity {
	return oi.quantity
}

// Price returns the price
func (oi *OrderItem) Price() *value.Money {
	return oi.price
}

// Subtotal calculates the subtotal for this item
func (oi *OrderItem) Subtotal() (*value.Money, error) {
	return oi.price.Multiply(oi.quantity.Value())
}

// Order represents a customer order
type Order struct {
	id        string
	items     []*OrderItem
	total     *value.Money
	status    OrderStatus
	createdAt time.Time
	updatedAt time.Time
}

// NewOrder creates a new order from basket items
func NewOrder(basketItems []*BasketItem) (*Order, error) {
	if len(basketItems) == 0 {
		return nil, errors.New("cannot create order with empty basket")
	}

	// Convert basket items to order items
	orderItems := make([]*OrderItem, 0, len(basketItems))
	for _, bi := range basketItems {
		orderItem, err := NewOrderItem(bi.ProductID(), bi.Quantity(), bi.Price())
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}

	// Calculate total
	total, err := value.NewMoney(0, basketItems[0].Price().Currency())
	if err != nil {
		return nil, err
	}

	for _, item := range orderItems {
		subtotal, err := item.Subtotal()
		if err != nil {
			return nil, err
		}
		total, err = total.Add(subtotal)
		if err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &Order{
		id:        uuid.New().String(),
		items:     orderItems,
		total:     total,
		status:    OrderStatusPending,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructOrder reconstructs an Order from persistence
func ReconstructOrder(id string, items []*OrderItem, total *value.Money, status OrderStatus, createdAt, updatedAt time.Time) *Order {
	return &Order{
		id:        id,
		items:     items,
		total:     total,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the order ID
func (o *Order) ID() string {
	return o.id
}

// Items returns the order items
func (o *Order) Items() []*OrderItem {
	return o.items
}

// Total returns the order total
func (o *Order) Total() *value.Money {
	return o.total
}

// Status returns the order status
func (o *Order) Status() OrderStatus {
	return o.status
}

// CreatedAt returns the creation time
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// UpdatedAt returns the last update time
func (o *Order) UpdatedAt() time.Time {
	return o.updatedAt
}

// Confirm confirms the order
func (o *Order) Confirm() error {
	if o.status != OrderStatusPending {
		return errors.New("only pending orders can be confirmed")
	}
	o.status = OrderStatusConfirmed
	o.updatedAt = time.Now()
	return nil
}

// Ship marks the order as shipped
func (o *Order) Ship() error {
	if o.status != OrderStatusConfirmed {
		return errors.New("only confirmed orders can be shipped")
	}
	o.status = OrderStatusShipped
	o.updatedAt = time.Now()
	return nil
}

// Deliver marks the order as delivered
func (o *Order) Deliver() error {
	if o.status != OrderStatusShipped {
		return errors.New("only shipped orders can be delivered")
	}
	o.status = OrderStatusDelivered
	o.updatedAt = time.Now()
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if o.status == OrderStatusDelivered {
		return errors.New("delivered orders cannot be cancelled")
	}
	if o.status == OrderStatusCancelled {
		return errors.New("order is already cancelled")
	}
	o.status = OrderStatusCancelled
	o.updatedAt = time.Now()
	return nil
}

// IsCancellable checks if the order can be cancelled
func (o *Order) IsCancellable() bool {
	return o.status != OrderStatusDelivered && o.status != OrderStatusCancelled
}
