package command

import (
	model "agnostic-payment-platform/internal/domain/payment"
	"time"
)

// CreatePaymentCommand represents the command to create a new payment.
type CreatePaymentCommand struct {
	ID          string        `json:"id"`
	Amount      float64       `json:"amount"`
	Currency    string        `json:"currency"`
	MerchantID  string        `json:"merchant_id"`
	Integration string        `json:"integration"`
	Status      string        `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	Payment     model.Payment `json:"payment"`
}

// NewCreatePaymentCommand creates a new instance of CreatePaymentCommand.
func NewCreatePaymentCommand(id, merchantID string, amount float64, currency string, integration string) *CreatePaymentCommand {
	return &CreatePaymentCommand{
		ID:          id,
		Amount:      amount,
		Currency:    currency,
		MerchantID:  merchantID,
		Integration: integration,
		Status:      "Pending",
		CreatedAt:   time.Now(),
	}
}
