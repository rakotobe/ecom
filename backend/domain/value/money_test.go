package value

import "testing"

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name      string
		amount    int64
		currency  string
		wantError bool
	}{
		{"valid money", 1000, "USD", false},
		{"zero amount", 0, "USD", false},
		{"negative amount", -100, "USD", true},
		{"empty currency", 1000, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if money.Amount() != tt.amount {
					t.Errorf("expected amount %d, got %d", tt.amount, money.Amount())
				}
				if money.Currency() != tt.currency {
					t.Errorf("expected currency %s, got %s", tt.currency, money.Currency())
				}
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	m1, _ := NewMoney(1000, "USD")
	m2, _ := NewMoney(500, "USD")
	m3, _ := NewMoney(1000, "EUR")

	result, err := m1.Add(m2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.Amount() != 1500 {
		t.Errorf("expected 1500, got %d", result.Amount())
	}

	_, err = m1.Add(m3)
	if err == nil {
		t.Error("expected error when adding different currencies")
	}
}

func TestMoney_Multiply(t *testing.T) {
	m, _ := NewMoney(1000, "USD")

	result, err := m.Multiply(3)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result.Amount() != 3000 {
		t.Errorf("expected 3000, got %d", result.Amount())
	}

	_, err = m.Multiply(-1)
	if err == nil {
		t.Error("expected error when multiplying by negative")
	}
}

func TestMoney_Equals(t *testing.T) {
	m1, _ := NewMoney(1000, "USD")
	m2, _ := NewMoney(1000, "USD")
	m3, _ := NewMoney(500, "USD")

	if !m1.Equals(m2) {
		t.Error("expected m1 to equal m2")
	}
	if m1.Equals(m3) {
		t.Error("expected m1 not to equal m3")
	}
}
