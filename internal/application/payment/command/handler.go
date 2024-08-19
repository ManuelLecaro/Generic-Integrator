package command

import (
	"agnostic-payment-platform/internal/application/integration"
	"agnostic-payment-platform/internal/application/payment/ports/repository"
	model "agnostic-payment-platform/internal/domain/payment"
	"context"
	"fmt"
	"time"

	"errors"
)

// CommandHandler maneja la ejecución de comandos relacionados con pagos.
type CommandHandler struct {
	paymentRepo          repository.PaymentRepository
	paymentEventStore    repository.PaymentEventStore
	integrationProcessor integration.ProcessorManager
}

// NewCommandHandler crea una nueva instancia de CommandHandler.
func NewCommandHandler(paymentRepo repository.PaymentRepository, paymentEventStore repository.PaymentEventStore, pm *integration.ProcessorManager) *CommandHandler {
	return &CommandHandler{
		paymentRepo:          paymentRepo,
		paymentEventStore:    paymentEventStore,
		integrationProcessor: *pm,
	}
}

// Handle procesa un comando dado.
func (h *CommandHandler) Handle(ctx context.Context, cmd interface{}) (string, error) {
	switch c := cmd.(type) {
	case CreatePaymentCommand:
		return h.handleCreatePaymentCommand(ctx, c)
	case RefundCommand:
		return h.handleRefundCommand(ctx, c)
	default:
		return "", errors.New("unknown command type")
	}
}

func (h *CommandHandler) handleCreatePaymentCommand(ctx context.Context, cmd CreatePaymentCommand) (string, error) {
	payment := &model.Payment{
		ID:              cmd.ID,
		MerchantID:      cmd.MerchantID,
		Amount:          cmd.Amount,
		Currency:        cmd.Currency,
		Status:          model.PaymentStatusPending,
		Integration:     cmd.Integration,
		CreatedAt:       time.Now(),
		PaymentMetadata: cmd.Payment.PaymentMetadata,
	}

	transactionID, err := h.integrate(ctx, payment, "authorize")
	if err != nil {
		fmt.Println("integration error:", err)
		go h.saveError(ctx, payment, err)

		return "", err
	}

	payment.TransactionID = transactionID
	payment.Status = model.PaymentStatusCompleted

	if err := h.paymentRepo.Save(ctx, payment); err != nil {
		go h.saveError(ctx, payment, err)
		return "", err
	}

	event := model.PaymentCreatedEvent{
		ID:            cmd.ID,
		Amount:        cmd.Amount,
		Currency:      cmd.Currency,
		TransactionId: transactionID,
		Payment:       payment,
		CreatedAt:     time.Now(),
	}

	if err := h.paymentEventStore.Save(ctx, &event); err != nil {
		go h.saveError(ctx, payment, err)
		return "", err
	}

	return transactionID, nil
}

func (h *CommandHandler) handleRefundCommand(ctx context.Context, cmd RefundCommand) (string, error) {
	payment, err := h.paymentRepo.GetPaymentByTransactionID(ctx, cmd.TransactionID)
	if err != nil {
		return "", err
	}

	payment.Status = model.PaymentStatusRefunded
	payment.Amount = cmd.Amount

	_, err = h.integrate(ctx, payment, "refund")
	if err != nil {
		go h.saveError(ctx, payment, err)

		return "", err
	}

	if err := h.paymentRepo.UpdateStatusByTransactionID(ctx, cmd.TransactionID, model.PaymentStatusRefunded); err != nil {
		go h.saveError(ctx, payment, err)
		return "", err
	}

	event := model.PaymentRefundedEvent{
		ID:            cmd.ID,
		RefundedAt:    time.Now(),
		Amount:        cmd.Amount,
		TransactionId: payment.TransactionID,
		MerchantID:    payment.MerchantID,
		Payment:       payment,
	}

	if err := h.paymentEventStore.Save(ctx, &event); err != nil {
		go h.saveError(ctx, payment, err)
		return "", err
	}

	return cmd.TransactionID, nil
}

func (h *CommandHandler) integrate(ctx context.Context, payment *model.Payment, action string) (string, error) {
	params := map[string]string{
		"merchant_id": payment.MerchantID,
		"amount":      fmt.Sprintf("%.2f", payment.Amount),
		"currency":    payment.Currency,
	}

	// Add specific parameters based on the action type
	switch action {
	case "authorize", "capture":
		params["card_number"] = payment.PaymentMetadata.CardNumber
		params["expiry_date"] = payment.PaymentMetadata.ExpiryDate
		params["cvv"] = payment.PaymentMetadata.CVV

	case "refund":
		params["transaction_id"] = payment.TransactionID
		params["refund_amount"] = fmt.Sprintf("%.2f", payment.Amount)

	case "status":
		params["payment_id"] = payment.ID
	}

	return h.integrationProcessor.Process(ctx, payment.Integration, action, params)
}

func (h *CommandHandler) saveError(ctx context.Context, payment *model.Payment, err error) error {
	return h.paymentEventStore.Save(ctx, &model.PaymentCreatedEventFailed{
		ID:        payment.ID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Payment:   payment,
		Error:     err.Error(),
		CreatedAt: time.Now(),
	})
}
