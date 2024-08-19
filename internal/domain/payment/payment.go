package payment

import (
	"time"
)

const (
	PaymentStatusPending   = "Pending"
	PaymentStatusRefunded  = "Refunded"
	PaymentStatusCompleted = "Completed"
	PaymentCreated         = "Created"
	PaymentUpdated         = "Updated"
	PaymentFailed          = "Failed"
)

// Payment representa un pago en el sistema.
type Payment struct {
	ID              string          `json:"id"`
	Amount          float64         `json:"amount"`
	MerchantID      string          `json:"merchant_id"`
	Currency        string          `json:"currency"`
	Status          string          `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Integration     string          `json:"type"`
	TransactionID   string          `json:"transaction_id"`
	PaymentMetadata PaymentMetadata `json:"metadata"`
}

type PaymentMetadata struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

// NewPayment crea una nueva instancia de Payment con los valores iniciales.
func NewPayment(id string, amount float64, currency string) *Payment {
	return &Payment{
		ID:        id,
		Amount:    amount,
		Currency:  currency,
		Status:    PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateStatus actualiza el estado del pago.
func (p *Payment) UpdateStatus(status string) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// ApplyEvent aplica un evento de dominio al pago.
func (p *Payment) ApplyEvent(event PaymentEvent) {
	switch e := event.(type) {
	case *PaymentCreatedEvent:
		p.ID = e.ID
		p.Amount = e.Amount
		p.Currency = e.Currency
		p.Status = "Pending"
		p.CreatedAt = e.CreatedAt
		p.UpdatedAt = e.CreatedAt
	case *PaymentRefundedEvent:
		p.Status = "Refunded"
		p.UpdatedAt = time.Now()
	}
}
