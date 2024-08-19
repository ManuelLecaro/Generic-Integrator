package routes

import (
	"agnostic-payment-platform/internal/infra/config"
	"agnostic-payment-platform/internal/infra/http/handler"
	"agnostic-payment-platform/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
)

type PaymentRouter struct {
	handler handler.PaymentHandler
	engine  *gin.Engine
	config  *config.Config
}

func NewPaymentRouter(handler *handler.PaymentHandler, engine *gin.Engine, config *config.Config) *PaymentRouter {
	return &PaymentRouter{
		handler: *handler,
		engine:  engine,
		config:  config,
	}
}

func (hr *PaymentRouter) Load() {
	group := hr.engine.Group("/payments")
	group.Use(middleware.APIKeyMiddleware(*hr.config))

	group.POST("/", hr.handler.ProcessPayment)
	group.GET("/:id", hr.handler.GetPaymentDetails)
	group.POST("/refund", hr.handler.Refund)
}
