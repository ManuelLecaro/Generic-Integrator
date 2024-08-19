package payment

import (
	"time"
)

// PaymentEvent represents a domain event related to a payment.
type PaymentEvent interface {
	EventType() string
}

// PaymentCreatedEvent represents the event of creating a payment.
type PaymentCreatedEvent struct {
	ID            string      `json:"id"`
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	TransactionId string      `json:"transaction_id"`
	Payment       interface{} `json:"payment"`
	MerchantID    string      `json:"merchant_id"`
	CreatedAt     time.Time   `json:"created_at"`
}

func (e *PaymentCreatedEvent) EventType() string {
	return PaymentCreated
}

// PaymentCreatedEventFailed represents the event of a failed payment creation.
type PaymentCreatedEventFailed struct {
	ID            string      `json:"id"`
	Amount        float64     `json:"amount"`
	Currency      string      `json:"currency"`
	TransactionId string      `json:"transaction_id"`
	Payment       interface{} `json:"payment"`
	Error         string      `json:"error"`
	MerchantID    string      `json:"merchant_id"`
	CreatedAt     time.Time   `json:"created_at"`
}

func (e *PaymentCreatedEventFailed) EventType() string {
	return PaymentFailed
}

// PaymentRefundedEvent represents the event of refunding a payment.
type PaymentRefundedEvent struct {
	ID            string      `json:"id"`
	RefundedAt    time.Time   `json:"refunded_at"`
	Amount        float64     `json:"amount"`
	TransactionId string      `json:"transaction_id"`
	Currency      string      `json:"currency"`
	MerchantID    string      `json:"merchant_id"`
	Payment       interface{} `json:"payment"`
}

func (e *PaymentRefundedEvent) EventType() string {
	return PaymentStatusRefunded
}

// PaymentStatusUpdatedEvent represents the event of updating the status of a payment.
type PaymentStatusUpdatedEvent struct {
	ID         string      `json:"id"`
	Status     string      `json:"status"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Amount     float64     `json:"amount"`
	MerchantID string      `json:"merchant_id"`
	Currency   string      `json:"currency"`
	Payment    interface{} `json:"payment"`
}

func (e *PaymentStatusUpdatedEvent) EventType() string {
	return "PaymentStatusUpdated"
}
