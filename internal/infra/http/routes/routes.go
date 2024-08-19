package routes

import (
	"agnostic-payment-platform/internal/infra/http/handler"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	Load()
}

type NewRoutesParams struct {
	fx.In
	HealthRouter  *GeneralRouter
	PaymentRouter *PaymentRouter
	Merchantouter *Merchantouter
}

func NewRoutes(rp NewRoutesParams) Routes {
	return Routes{
		rp.HealthRouter,
		rp.PaymentRouter,
		rp.Merchantouter,
	}
}

func (r Routes) Load() {
	for _, route := range r {
		route.Load()
	}
}

var Module = fx.Options(
	handler.Module,
	fx.Provide(
		NewRoutes,
		NewHealthRouter,
		NewPaymentRouter,
		NewMerchantouter,
	),
)
