package repository

import (
	model "agnostic-payment-platform/internal/domain/payment"
	"context"
)

// PaymentRepository define las operaciones que se pueden realizar en el almacenamiento de pagos.
type PaymentRepository interface {
	// Save guarda un pago en el repositorio.
	Save(ctx context.Context, payment *model.Payment) error

	// FindByID busca un pago por su ID único.
	FindByID(ctx context.Context, id string) (*model.Payment, error)

	// GetPaymentByTransactionID busca un pago por su ID de transacción.
	GetPaymentByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error)

	// ListPayments lista los pagos según los filtros opcionales.
	ListPayments(ctx context.Context, merchantID, status string, limit, offset int) ([]*model.Payment, error)

	// UpdateStatus actualiza el estado de un pago.
	UpdateStatus(ctx context.Context, id string, status string) error

	UpdateStatusByTransactionID(ctx context.Context, transactionID string, status string) error

	// SaveEvent almacena un evento de pago en el repositorio.
	SaveEvent(ctx context.Context, event *model.PaymentEvent) error

	// ListEvents lista todos los eventos asociados con un pago específico.
	ListEvents(ctx context.Context, paymentID string) ([]*model.PaymentEvent, error)

	// SaveSnapshot guarda un snapshot del estado actual de un pago.
	SaveSnapshot(ctx context.Context, snapshot *model.PaymentSnapshot) error

	// GetLatestSnapshot obtiene el snapshot más reciente de un pago por su ID.
	GetLatestSnapshot(ctx context.Context, paymentID string) (*model.PaymentSnapshot, error)
}

// PaymentEventStore define los métodos para almacenar y recuperar eventos de pago.
type PaymentEventStore interface {
	// Save persiste un evento de pago en el almacén de eventos.
	Save(ctx context.Context, event model.PaymentEvent) error

	// GetByID recupera eventos de pago por el identificador del pago.
	GetByID(ctx context.Context, paymentID string) ([]model.PaymentEvent, error)

	// GetByTransactionID recupera eventos de pago por el identificador de transacción.
	GetByTransactionID(ctx context.Context, transactionID string) ([]model.PaymentEvent, error)
}
