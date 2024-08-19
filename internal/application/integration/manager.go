package integration

import (
	"agnostic-payment-platform/internal/application/payment/ports/services"
	"context"
	"fmt"
)

// ProcessorManager manages available payment processors.
type ProcessorManager struct {
	processors map[string]services.Processor
}

// NewProcessorManager creates a new ProcessorManager.
func NewProcessorManager(processors map[string]services.Processor) *ProcessorManager {
	return &ProcessorManager{processors: processors}
}

// Process delegates the processing request to the appropriate processor.
func (pm *ProcessorManager) Process(ctx context.Context, name string, action string, params map[string]string) (string, error) {
	processor, ok := pm.processors[name]
	if !ok {
		return "", fmt.Errorf("processor %s not found", name)
	}
	return processor.Process(ctx, action, params)
}
