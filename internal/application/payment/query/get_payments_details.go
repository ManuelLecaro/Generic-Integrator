package query

import (
	"agnostic-payment-platform/internal/application/payment/dto"
	"agnostic-payment-platform/internal/domain/payment"
	"context"

	"errors"
)

// GetPaymentDetailsQuery es la estructura que encapsula los parámetros de la consulta.
type GetPaymentDetailsQuery struct {
	PaymentID string
}

// HandleGetPaymentDetails maneja la consulta para obtener los detalles de un pago.
func (h *QueryHandler) HandleGetPaymentDetails(ctx context.Context, query *GetPaymentDetailsQuery) (*dto.PaymentDetailsDTO, error) {
	pm, err := h.paymentRepo.FindByID(ctx, query.PaymentID)
	if err != nil {
		if errors.Is(err, payment.ErrPaymentNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	// Mapear los datos del modelo de dominio al DTO
	paymentDetails := &dto.PaymentDetailsDTO{
		ID:             pm.ID,
		Amount:         pm.Amount,
		Currency:       pm.Currency,
		MerchantID:     pm.MerchantID,
		PaymentStatus:  pm.Status,
		CreatedAt:      pm.CreatedAt,
		RefundedAmount: pm.Amount,
	}

	return paymentDetails, nil
}
