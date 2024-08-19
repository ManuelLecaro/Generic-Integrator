package command

import (
	"time"
)

// RefundCommand encapsula los datos necesarios para procesar un reembolso.
type RefundCommand struct {
	ID            string    `json:"id"`
	TransactionID string    // ID de la transacción que se va a reembolsar
	Amount        float64   // Monto del reembolso
	RefundReason  string    // Motivo del reembolso
	RequestedAt   time.Time // Fecha y hora en que se solicitó el reembolso
}

// NewRefundCommand es un constructor que crea una nueva instancia de RefundCommand.
func NewRefundCommand(transactionID string, amount float64, reason string) *RefundCommand {
	return &RefundCommand{
		TransactionID: transactionID,
		Amount:        amount,
		RefundReason:  reason,
		RequestedAt:   time.Now(),
	}
}
