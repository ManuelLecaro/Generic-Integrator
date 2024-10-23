package executor

import (
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"
	"generic-integration-platform/internal/infra/executor/cli"
	"generic-integration-platform/internal/infra/executor/rest"

	"go.uber.org/fx"
)

type Executor interface {
	Execute(integration integration.Integration, endpoint endpoint.Endpoint) (map[string]interface{}, error)
}

// NewExecutorProvider initializes the executors and provides them as a map.
func NewExecutorProvider() fx.Option {
	return fx.Provide(
		rest.NewRestExecutor,
		cli.NewCLIExecutor,
		fx.Annotate(
			aggregateExecutors,
			fx.ResultTags(`group:"executor"`),
		),
	)
}

func aggregateExecutors(
	restExec *rest.RestExecutor,
	cliExec *cli.CLIExecutor,
) (map[string]Executor, error) {
	executors := make(map[string]Executor)
	executors["rest"] = restExec
	executors["cli"] = cliExec
	return executors, nil
}
