package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	swaggerFiles "github.com/swaggo/files"
)

type GeneralRouter struct {
	engine *gin.Engine
}

func NewHealthRouter(engine *gin.Engine) *GeneralRouter {
	return &GeneralRouter{
		engine: engine,
	}
}

type HealthInput struct{}

type HealthOutput struct {
	Body string `json:"message"`
}

func (hr *GeneralRouter) Load() {
	hr.engine.Use(otelgin.Middleware("generic-integration"))
	hr.engine.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	hr.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
