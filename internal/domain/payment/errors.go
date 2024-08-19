package payment

import "errors"

// Errores posibles en el dominio de pagos.
var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidAmount   = errors.New("invalid amount")
)
