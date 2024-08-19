package services

import (
	"agnostic-payment-platform/internal/domain/payment"
	"context"
)

// Service es la interfaz que define los métodos para manejar pagos.
type Service interface {
	// ProcessPayment procesa un pago y devuelve una respuesta con el estado del pago.
	ProcessPayment(ctx context.Context, request *payment.Payment) (*payment.Payment, error)

	// GetPaymentDetails obtiene los detalles de un pago usando su ID.
	GetPaymentDetails(ctx context.Context, id string) (*payment.Payment, error)
}

type Processor interface {
	Process(ctx context.Context, action string, params map[string]string) (string, error)
}
