package command

import (
	"agnostic-payment-platform/internal/application/integration"
	"agnostic-payment-platform/internal/application/payment/ports/repository"
	"context"
	"testing"
)

func TestCommandHandler_Handle(t *testing.T) {
	type fields struct {
		paymentRepo          repository.PaymentRepository
		paymentEventStore    repository.PaymentEventStore
		integrationProcessor integration.ProcessorManager
	}
	type args struct {
		ctx context.Context
		cmd interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CommandHandler{
				paymentRepo:          tt.fields.paymentRepo,
				paymentEventStore:    tt.fields.paymentEventStore,
				integrationProcessor: tt.fields.integrationProcessor,
			}
			got, err := h.Handle(tt.args.ctx, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CommandHandler.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
