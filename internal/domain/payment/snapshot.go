package payment

import "time"

// PaymentSnapshot represents a snapshot of the payment state at a specific point in time.
type PaymentSnapshot struct {
	ID        string    `json:"id"`         // Payment identifier
	Amount    float64   `json:"amount"`     // Payment amount
	Currency  string    `json:"currency"`   // Payment currency
	Status    string    `json:"status"`     // Payment status
	CreatedAt time.Time `json:"created_at"` // Date and time the payment was created
	UpdatedAt time.Time `json:"updated_at"` // Date and time the payment was last updated
}

// NewPaymentSnapshot creates a new instance of PaymentSnapshot with the current Payment values.
func NewPaymentSnapshot(payment *Payment) *PaymentSnapshot {
	return &PaymentSnapshot{
		ID:        payment.ID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}
}
