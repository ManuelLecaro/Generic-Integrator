package extender

import (
	"context"
	"generic-integration-platform/internal/domain/integration"
)

// IntegrationExtender defines the interface for adding new integration types.
type IntegrationExtender interface {
	// Initialize sets up the integration with the given configuration.
	Initialize(ctx context.Context, config *integration.Integration) error

	// Execute performs the defined action on the integration with given parameters.
	Execute(ctx context.Context, action string, params map[string]interface{}) (interface{}, error)

	// Validate validates the configuration of the integration.
	Validate(ctx context.Context) error

	// Close cleans up any resources used by the integration.
	Close(ctx context.Context) error
}
