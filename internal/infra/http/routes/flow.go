package routes

import (
	"generic-integration-platform/internal/infra/config"
	"generic-integration-platform/internal/infra/http/handler"
	"generic-integration-platform/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
)

type FlowRouter struct {
	handler handler.FlowHandler
	engine  *gin.Engine
	config  *config.Config
}

func NewFlowRouter(handler *handler.FlowHandler, engine *gin.Engine, config *config.Config) *FlowRouter {
	return &FlowRouter{
		handler: *handler,
		engine:  engine,
		config:  config,
	}
}

func (fr *FlowRouter) Load() {
	group := fr.engine.Group("/flows")
	group.Use(middleware.APIKeyMiddleware(*fr.config))

	group.GET("/", fr.handler.GetFlows)                // List all flows
	group.POST("/", fr.handler.CreateFlow)             // Create a new flow
	group.GET("/:id", fr.handler.GetFlowDetails)       // Get a specific flow by ID
	group.PUT("/:id", fr.handler.UpdateFlow)           // Update an existing flow by ID
	group.DELETE("/:id", fr.handler.DeleteFlow)        // Delete an existing flow by ID
	group.POST("/:id/execute", fr.handler.ExecuteFlow) // Execute a specific flow
}
