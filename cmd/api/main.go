package main

import (
	"agnostic-payment-platform/internal/application/integration"
	"agnostic-payment-platform/internal/application/merchant"
	"agnostic-payment-platform/internal/application/payment"
	"agnostic-payment-platform/internal/infra/config"
	"agnostic-payment-platform/internal/infra/db"
	"agnostic-payment-platform/internal/infra/eventstore"
	"agnostic-payment-platform/internal/infra/http"
	"agnostic-payment-platform/internal/infra/http/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	_ "agnostic-payment-platform/docs"
)

// @title Agnostic Payment Platform API
// @version 1.0
// @description This is an agnostic payment platform API that processes and retrieves payment details.
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

	fmt.Println("DATA : : ", cfg)

	iCfg, err := config.LoadIntegrationConfig("payments.toml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	fmt.Println("DATA : : ", iCfg)

	fx.New(
		fx.Provide(func() (*config.Config, *config.IntegrationConfig) {
			return cfg, iCfg
		}),
		http.Module,
		routes.Module,
		payment.Module,
		merchant.Module,
		integration.Module,
		eventstore.Module,
		db.Module,
		fx.Invoke(
			routes.Routes.Load,
			func(r *gin.Engine) {},
		),
	).Run()
}
