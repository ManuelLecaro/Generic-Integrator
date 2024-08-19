package dto

import "time"

// PaymentRequestDTO represents the data sent to process a new payment.
type PaymentRequestDTO struct {
	MerchantID string  `json:"merchant_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
	Currency   string  `json:"currency" validate:"required,len=3"`
	CardNumber string  `json:"card_number" validate:"required"`
	ExpiryDate string  `json:"expiry_date" validate:"required"`
	CVV        string  `json:"cvv" validate:"required,len=3"`
	Type       string  `json:"type" validate:"required"`
}

// PaymentResponseDTO represents the response after processing a payment.
type PaymentResponseDTO struct {
	TransactionID string  `json:"transaction_id" validate:"required"`
	PaymentID     string  `json:"payment_id"`
	PaymentStatus string  `json:"payment_status"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Message       string  `json:"message"`
}

// PaymentsByIDDTO represents the data used to query payment details by ID.
type PaymentsByIDDTO struct {
	ID string `path:"id"`
}

// PaymentDetailsDTO represents the details of a retrieved payment.
type PaymentDetailsDTO struct {
	ID             string    `json:"id" example:"uuid123421"`
	MerchantID     string    `json:"merchant_id" example:"uuidasfhtu12"`
	Amount         float64   `json:"amount" example:"12.45"`
	Currency       string    `json:"currency" example:"USD"`
	PaymentStatus  string    `json:"payment_status"`
	CreatedAt      time.Time `json:"created_at"`
	RefundedAmount float64   `json:"refunded_amount,omitempty"`
}

// RefundDTO represents the data sent to process a refund.
type RefundDTO struct {
	PaymentID     string  `json:"payment_id" validate:"required"`
	TransactionID string  `json:"transaction_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Reason        string  `json:"reason"`
}

// RefundResponseDTO represents the response after processing a refund.
type RefundResponseDTO struct {
	TransactionID  string  `json:"transaction_id"`
	RefundedAmount float64 `json:"refunded_amount"`
	Status         string  `json:"status"`
	Message        string  `json:"message"`
}
