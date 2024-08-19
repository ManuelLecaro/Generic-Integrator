package service

import (
	"agnostic-payment-platform/internal/application/payment/command"
	"agnostic-payment-platform/internal/application/payment/dto"
	query "agnostic-payment-platform/internal/application/payment/query"
	"agnostic-payment-platform/internal/domain/payment"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentService es la implementación de la interfaz de servicio para manejar operaciones de pago.
type PaymentService struct {
	commandHandler command.CommandHandler
	queryHandler   query.QueryHandler
}

// NewPaymentService crea una nueva instancia de PaymentService.
func NewPaymentService(commandHandler *command.CommandHandler, queryHandler *query.QueryHandler) *PaymentService {
	return &PaymentService{
		commandHandler: *commandHandler,
		queryHandler:   *queryHandler,
	}
}

// ProcessPayment procesa un nuevo pago basado en el DTO de solicitud.
func (s *PaymentService) ProcessPayment(ctx context.Context, input *dto.PaymentRequestDTO) (*dto.PaymentResponseDTO, error) {
	createPaymentCmd := command.CreatePaymentCommand{
		ID:          uuid.NewString(),
		Amount:      input.Amount,
		Currency:    input.Currency,
		MerchantID:  input.MerchantID,
		Integration: input.Type,
		Status:      payment.PaymentStatusPending,
		CreatedAt:   time.Now(),
		Payment: payment.Payment{
			PaymentMetadata: payment.PaymentMetadata{
				CardNumber: input.CardNumber,
				ExpiryDate: input.ExpiryDate,
				CVV:        input.CVV,
			},
		},
	}

	tID, err := s.commandHandler.Handle(ctx, createPaymentCmd)
	if err != nil {
		return nil, err
	}

	response := &dto.PaymentResponseDTO{
		TransactionID: tID,
		PaymentID:     createPaymentCmd.ID,
		PaymentStatus: string(payment.PaymentStatusPending),
		Amount:        input.Amount,
		Currency:      input.Currency,
	}

	return response, nil
}

// GetPaymentDetails devuelve los detalles de un pago basado en el ID.
func (s *PaymentService) GetPaymentDetails(ctx context.Context, id string) ([]dto.PaymentDetailsDTO, error) {
	getPaymentDetailsQuery := query.GetPaymentDetailsQuery{
		PaymentID: id,
	}

	paymentDetails, err := s.queryHandler.Handle(ctx, getPaymentDetailsQuery)
	if err != nil {
		return nil, err
	}

	details, ok := paymentDetails.([]dto.PaymentDetailsDTO)
	if !ok {
		return nil, fmt.Errorf("not a valid payment")

	}

	return details, nil
}

// ProcessRefund maneja el reembolso de un pago.
func (s *PaymentService) ProcessRefund(ctx context.Context, input *dto.RefundDTO) (*dto.PaymentResponseDTO, error) {
	refundCmd := command.RefundCommand{
		ID:            input.PaymentID,
		TransactionID: input.TransactionID,
		Amount:        input.Amount,
		RefundReason:  input.Reason,
		RequestedAt:   time.Now(),
	}

	_, err := s.commandHandler.Handle(ctx, refundCmd)
	if err != nil {
		return nil, err
	}

	response := &dto.PaymentResponseDTO{
		PaymentID:     refundCmd.TransactionID,
		PaymentStatus: payment.PaymentStatusRefunded,
	}

	return response, nil
}
