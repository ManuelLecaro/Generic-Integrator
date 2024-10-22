package routes

import (
	"generic-integration-platform/internal/infra/config"
	"generic-integration-platform/internal/infra/http/handler"
	"generic-integration-platform/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
)

type IntegrationRouter struct {
	handler handler.IntegrationHandler
	engine  *gin.Engine
	config  *config.Config
}

func NewIntegrationRouter(handler *handler.IntegrationHandler, engine *gin.Engine, config *config.Config) *IntegrationRouter {
	return &IntegrationRouter{
		handler: *handler,
		engine:  engine,
		config:  config,
	}
}

func (ir *IntegrationRouter) Load() {
	group := ir.engine.Group("/integrations")
	group.Use(middleware.APIKeyMiddleware(*ir.config))

	group.GET("/", ir.handler.GetIntegrations)          // List all integrations
	group.POST("/", ir.handler.CreateIntegration)       // Create a new integration
	group.GET("/:id", ir.handler.GetIntegrationDetails) // Get a specific integration by ID
	group.PUT("/:id", ir.handler.UpdateIntegration)     // Update an existing integration by ID
	group.DELETE("/:id", ir.handler.DeleteIntegration)  // Delete an existing integration by ID
}
