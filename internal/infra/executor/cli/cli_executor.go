package cli

import (
	"bytes"
	"fmt"
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"
	"os/exec"
)

// CLIExecutor is an implementation of the executor for CLI commands.
type CLIExecutor struct{}

func NewCLIExecutor() *CLIExecutor {
	return &CLIExecutor{}
}

// Execute runs a command in the command line interface and returns the output.
func (e *CLIExecutor) Execute(integration integration.Integration, endpoint endpoint.Endpoint) (map[string]interface{}, error) {
	cmd := exec.Command(endpoint.Action) // Assuming Action contains the command to run

	// Prepare the command parameters
	var params []string
	for key, value := range endpoint.Params {
		params = append(params, fmt.Sprintf("--%s=%s", key, value)) // Assuming parameters are passed as flags
	}
	cmd.Args = append(cmd.Args, params...)

	// Execute the command and capture the output
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output // Capture standard error as well

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute command '%s': %w, output: %s", endpoint.Action, err, output.String())
	}

	// Process the output into a map
	response := map[string]interface{}{
		"output": output.String(),
	}

	return response, nil
}
