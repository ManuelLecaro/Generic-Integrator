package payment

import (
	"agnostic-payment-platform/internal/application/payment/command"
	"agnostic-payment-platform/internal/application/payment/query"
	"agnostic-payment-platform/internal/application/payment/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		service.NewPaymentService,
		query.NewQueryHandler,
		command.NewCommandHandler,
	),
)
