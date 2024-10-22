package routes

import (
	"generic-integration-platform/internal/infra/http/handler"

	"go.uber.org/fx"
)

type Routes []Route

type Route interface {
	Load()
}

type NewRoutesParams struct {
	fx.In
	HealthRouter *GeneralRouter
}

func NewRoutes(rp NewRoutesParams) Routes {
	return Routes{
		rp.HealthRouter,
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
	),
)
