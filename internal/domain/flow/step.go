package flow

import (
	"errors"
)

// Step represents an individual step in a flow of an integration process.
type Step struct {
	ID            string                 // Unique identifier for the step
	Name          string                 // Name of the step
	IntegrationID string                 // ID of the integration to use
	Action        string                 // Action to be performed (e.g., "authorize", "capture", etc.)
	Params        map[string]interface{} // Parameters to be sent to the endpoint
	NextStepID    string                 // ID of the next step (for transitioning)
}

// New creates a new Step instance.
func New(id, name, integrationID, action string, params map[string]interface{}, nextStepID string) *Step {
	return &Step{
		ID:            id,
		Name:          name,
		IntegrationID: integrationID,
		Action:        action,
		Params:        params,
		NextStepID:    nextStepID,
	}
}

// Validate checks if the step has the necessary fields set.
func (s *Step) Validate() error {
	if s.Action == "" {
		return errors.New("step action cannot be empty")
	}
	if s.IntegrationID == "" {
		return errors.New("step integration cannot be empty")
	}
	return nil
}
