package query

import (
	"agnostic-payment-platform/internal/application/payment/dto"
	"agnostic-payment-platform/internal/application/payment/ports/repository"
	model "agnostic-payment-platform/internal/domain/payment"
	"context"
	"fmt"
)

// QueryHandler handles the execution of queries.
type QueryHandler struct {
	paymentRepo repository.PaymentRepository
	storeDB     repository.PaymentEventStore
}

// NewQueryHandler is the constructor for QueryHandler.
func NewQueryHandler(paymentRepo repository.PaymentRepository, storeDB repository.PaymentEventStore) *QueryHandler {
	return &QueryHandler{
		paymentRepo: paymentRepo,
		storeDB:     storeDB,
	}
}

// Handle handles specific queries and delegates to the appropriate service.
func (h *QueryHandler) Handle(ctx context.Context, query interface{}) (interface{}, error) {
	switch q := query.(type) {
	case GetPaymentDetailsQuery:
		return h.handleGetPaymentDetailsQuery(ctx, &q)
	case ListPaymentsQuery:
		return h.handleListPaymentsQuery(ctx, &q)
	default:
		return nil, model.ErrPaymentNotFound
	}
}

// handleGetPaymentDetailsQuery handles the query to get the details of a payment.
func (h *QueryHandler) handleGetPaymentDetailsQuery(ctx context.Context, query *GetPaymentDetailsQuery) ([]dto.PaymentDetailsDTO, error) {
	events, err := h.storeDB.GetByTransactionID(ctx, "payment-"+query.PaymentID)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no payment found with transaction ID: %s", query.PaymentID)
	}

	var details []dto.PaymentDetailsDTO

	for _, event := range events {
		switch e := event.(type) {
		case *model.PaymentCreatedEvent:
			details = append(details, dto.PaymentDetailsDTO{
				ID:            e.ID,
				Amount:        e.Amount,
				MerchantID:    e.MerchantID,
				Currency:      e.Currency,
				CreatedAt:     e.CreatedAt,
				PaymentStatus: model.PaymentCreated,
			})
		case *model.PaymentCreatedEventFailed:
			details = append(details, dto.PaymentDetailsDTO{
				ID:            e.ID,
				Amount:        e.Amount,
				MerchantID:    e.MerchantID,
				Currency:      e.Currency,
				CreatedAt:     e.CreatedAt,
				PaymentStatus: model.PaymentFailed,
			})
		case *model.PaymentRefundedEvent:
			details = append(details, dto.PaymentDetailsDTO{
				ID:             e.ID,
				Amount:         e.Amount,
				Currency:       e.Currency,
				MerchantID:     e.MerchantID,
				CreatedAt:      e.RefundedAt,
				PaymentStatus:  model.PaymentStatusRefunded,
				RefundedAmount: e.Amount,
			})
		case *model.PaymentStatusUpdatedEvent:
			details = append(details, dto.PaymentDetailsDTO{
				ID:            e.ID,
				Amount:        e.Amount,
				Currency:      e.Currency,
				MerchantID:    e.MerchantID,
				CreatedAt:     e.UpdatedAt,
				PaymentStatus: e.Status,
			})
		default:
			return nil, fmt.Errorf("unknown event type: %T", e)
		}
	}

	return details, nil
}

// handleListPaymentsQuery handles the query to list payments (you can adapt this according to your needs).
func (h *QueryHandler) handleListPaymentsQuery(ctx context.Context, query *ListPaymentsQuery) ([]dto.PaymentDetailsDTO, error) {
	payments, err := h.paymentRepo.ListPayments(ctx, query.MerchantID, query.Status, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	var dtos []dto.PaymentDetailsDTO
	for _, payment := range payments {
		dtos = append(dtos, dto.PaymentDetailsDTO{
			ID:            payment.ID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			CreatedAt:     payment.CreatedAt,
			PaymentStatus: payment.Status,
		})
	}
	return dtos, nil
}
