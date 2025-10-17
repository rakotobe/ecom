package value

import "errors"

// Quantity represents a quantity of items
type Quantity struct {
	value int
}

// NewQuantity creates a new Quantity value object
func NewQuantity(value int) (*Quantity, error) {
	if value < 0 {
		return nil, errors.New("quantity cannot be negative")
	}
	return &Quantity{value: value}, nil
}

// Value returns the quantity value
func (q *Quantity) Value() int {
	return q.value
}

// Add adds two quantities
func (q *Quantity) Add(other *Quantity) (*Quantity, error) {
	return NewQuantity(q.value + other.value)
}

// Subtract subtracts a quantity
func (q *Quantity) Subtract(other *Quantity) (*Quantity, error) {
	result := q.value - other.value
	if result < 0 {
		return nil, errors.New("resulting quantity cannot be negative")
	}
	return NewQuantity(result)
}

// IsZero checks if quantity is zero
func (q *Quantity) IsZero() bool {
	return q.value == 0
}

// Equals checks if two quantities are equal
func (q *Quantity) Equals(other *Quantity) bool {
	return q.value == other.value
}
