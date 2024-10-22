package validators

import (
	"errors"
	"generic-integration-platform/internal/application/dto"
	"strings"
)

// ValidateFlow validates the FlowDTO input.
func ValidateFlow(flow dto.FlowDTO) error {
	if strings.TrimSpace(flow.Name) == "" {
		return errors.New("flow name is required")
	}

	if len(flow.Steps) == 0 {
		return errors.New("at least one step is required in the flow")
	}

	for _, step := range flow.Steps {
		if err := ValidateStep(step); err != nil {
			return err
		}
	}

	return nil
}

// ValidateStep validates a step in the flow.
func ValidateStep(step dto.StepDTO) error {
	if strings.TrimSpace(step.Action) == "" {
		return errors.New("action is required for step")
	}

	if strings.TrimSpace(step.IntegrationID) == "" {
		return errors.New("integration ID is required for step")
	}

	return nil
}
