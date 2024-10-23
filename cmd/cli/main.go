package main

import (
	"fmt"
	"generic-integration-platform/internal/infra/config"
	"generic-integration-platform/internal/infra/db"
	"generic-integration-platform/internal/infra/eventstore"
	"generic-integration-platform/internal/infra/http"
	"generic-integration-platform/internal/infra/http/routes"
	"generic-integration-platform/internal/infra/monitoring"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

// RunCLI initializes and runs the CLI application
func RunCLI() error {
	cliApp := &cli.App{
		Name:  "SaaS CLI",
		Usage: "A command-line interface for managing integrations",
		Commands: []*cli.Command{
			{
				Name:  "execute",
				Usage: "Execute a flow",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "id",
						Usage:    "ID of the flow to execute",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	return cliApp.Run(os.Args)
}

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	fx.New(
		fx.Provide(func() *config.Config {
			return cfg
		}),
		monitoring.Module,
		http.Module,
		routes.Module,
		eventstore.Module,
		db.Module,
		fx.Invoke(func() {
			if err := RunCLI(); err != nil {
				fmt.Printf("Error running CLI: %v\n", err)
			}
		}),
	).Run()
}
