package value

import (
	"errors"
	"fmt"
)

// Money represents a monetary value with currency
type Money struct {
	amount   int64  // stored in cents to avoid floating point issues
	currency string
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}
	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount returns the amount in cents
func (m *Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency code
func (m *Money) Currency() string {
	return m.currency
}

// Add adds two Money values (must be same currency)
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, errors.New("cannot add money with different currencies")
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// Multiply multiplies the money by a quantity
func (m *Money) Multiply(quantity int) (*Money, error) {
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}
	return NewMoney(m.amount*int64(quantity), m.currency)
}

// String returns a string representation
func (m *Money) String() string {
	dollars := m.amount / 100
	cents := m.amount % 100
	return fmt.Sprintf("%s %d.%02d", m.currency, dollars, cents)
}

// Equals checks if two Money values are equal
func (m *Money) Equals(other *Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
