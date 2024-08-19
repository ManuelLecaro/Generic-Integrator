package routes

import (
	"agnostic-payment-platform/internal/infra/config"
	"agnostic-payment-platform/internal/infra/http/handler"
	"agnostic-payment-platform/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
)

type Merchantouter struct {
	handler handler.MerchantHandler
	engine  *gin.Engine
	config  *config.Config
}

func NewMerchantouter(handler *handler.MerchantHandler, engine *gin.Engine, config *config.Config) *Merchantouter {
	return &Merchantouter{
		handler: *handler,
		engine:  engine,
		config:  config,
	}
}

func (hr *Merchantouter) Load() {
	group := hr.engine.Group("/merchant")

	group.Use(middleware.APIKeyMiddleware(*hr.config))

	group.POST("/signup", hr.handler.Signup)
	group.POST("/login", hr.handler.Login)
}
