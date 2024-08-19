package merchant

import (
	"agnostic-payment-platform/internal/application/merchant/service"

	"go.uber.org/fx"
)

var Module = fx.Option(
	fx.Provide(
		service.NewMerchantService,
	),
)
