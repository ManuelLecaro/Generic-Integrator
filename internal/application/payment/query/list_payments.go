package query

import (
	"agnostic-payment-platform/internal/application/payment/dto"
	"context"
)

// ListPaymentsQuery encapsula los parámetros opcionales para la consulta de la lista de pagos.
type ListPaymentsQuery struct {
	MerchantID string
	Status     string
	Limit      int
	Offset     int
}

// HandleListPayments maneja la consulta para listar los pagos.
func (h *QueryHandler) HandleListPayments(ctx context.Context, query *ListPaymentsQuery) ([]*dto.PaymentDetailsDTO, error) {
	payments, err := h.paymentRepo.ListPayments(ctx, query.MerchantID, query.Status, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	// Mapear los datos del modelo de dominio a los DTOs
	paymentDetailsList := make([]*dto.PaymentDetailsDTO, len(payments))
	for i, payment := range payments {
		paymentDetailsList[i] = &dto.PaymentDetailsDTO{
			ID:             payment.ID,
			MerchantID:     payment.MerchantID,
			Amount:         payment.Amount,
			Currency:       payment.Currency,
			PaymentStatus:  string(payment.Status),
			CreatedAt:      payment.CreatedAt,
			RefundedAmount: payment.Amount,
		}
	}

	return paymentDetailsList, nil
}
