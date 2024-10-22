package main

import (
	"generic-integration-platform/internal/infra/config"
	"generic-integration-platform/internal/infra/db"
	"generic-integration-platform/internal/infra/eventstore"
	"generic-integration-platform/internal/infra/http"
	"generic-integration-platform/internal/infra/http/routes"
	"generic-integration-platform/internal/infra/monitoring"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	_ "generic-integration-platform/docs"
)

// @title Generic Integration Platform API
// @version 1.0
// @description This is an Generic Integration Platform API that processes multiple integration and integration flows.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
// @description To access protected routes, add the x-api-key header with your API key. Example: "x-api-key: your-api-key"

// @security ApiKeyAuth
func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	iCfg, err := config.LoadIntegrationConfig("payments.toml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	fx.New(
		fx.Provide(func() (*config.Config, *config.IntegrationConfig) {
			return cfg, iCfg
		}),
		monitoring.Module,
		http.Module,
		routes.Module,
		eventstore.Module,
		db.Module,
		fx.Invoke(
			routes.Routes.Load,
			func(r *gin.Engine) {},
		),
	).Run()
}
