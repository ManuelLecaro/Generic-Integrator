package flow

import (
	"errors"
)

// Flow represents a sequence of steps that define the integration process.
type Flow struct {
	ID          string  `json:"id,omitempty"`
	Name        string  // The name of the flow
	Steps       []*Step // List of steps in the flow
	Description string  // A brief description of the flow
}

// NewFlow creates a new Flow instance.
func NewFlow(name, description string, steps []*Step) *Flow {
	return &Flow{
		Name:        name,
		Steps:       steps,
		Description: description,
	}
}

// AddStep adds a new step to the flow.
func (f *Flow) AddStep(step *Step) {
	f.Steps = append(f.Steps, step)
}

// Validate checks if the flow has the necessary fields set and if steps are valid.
func (f *Flow) Validate() error {
	if f.Name == "" {
		return errors.New("flow name cannot be empty")
	}
	if len(f.Steps) == 0 {
		return errors.New("flow must contain at least one step")
	}
	for _, step := range f.Steps {
		if err := step.Validate(); err != nil {
			return err
		}
	}
	return nil
}
