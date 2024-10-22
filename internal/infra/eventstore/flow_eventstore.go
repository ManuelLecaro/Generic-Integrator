package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"generic-integration-platform/internal/domain/flow"
	"time"

	esdb "github.com/EventStore/EventStore-Client-Go/esdb"
)

// FlowEventStore handles storing flow events in EventStoreDB.
type FlowEventStore struct {
	client *esdb.Client
}

// NewFlowEventStore creates a new instance of FlowEventStore.
func NewFlowEventStore(client *esdb.Client) *FlowEventStore {
	return &FlowEventStore{
		client: client,
	}
}

// AppendFlowCreatedEvent stores the FlowCreated event in EventStore.
func (store *FlowEventStore) AppendFlowCreatedEvent(ctx context.Context, event FlowCreatedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowCreatedEvent")
}

// AppendFlowStepCompletedEvent stores the FlowStepCompleted event in EventStore.
func (store *FlowEventStore) AppendFlowStepCompletedEvent(ctx context.Context, event FlowStepCompletedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowStepCompletedEvent")
}

// AppendFlowUpdatedEvent stores the FlowUpdated event in EventStore.
func (store *FlowEventStore) AppendFlowUpdatedEvent(ctx context.Context, event FlowUpdatedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowUpdatedEvent")
}

// AppendFlowDeletedEvent stores the FlowDeleted event in EventStore.
func (store *FlowEventStore) AppendFlowDeletedEvent(ctx context.Context, event FlowDeletedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowDeletedEvent")
}

// AppendFlowStepFailedEvent stores the FlowStepFailed event in EventStore.
func (store *FlowEventStore) AppendFlowStepFailedEvent(ctx context.Context, event FlowStepFailedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowStepFailedEvent")
}

// AppendFlowExecutedEvent stores the FlowExecuted event in EventStore.
func (store *FlowEventStore) AppendFlowExecutedEvent(ctx context.Context, event FlowExecutedEvent) error {
	return store.appendEvent(ctx, "flow-"+event.FlowID, event, "FlowExecutedEvent")
}

// appendEvent serializes and stores a generic event in EventStoreDB.
func (store *FlowEventStore) appendEvent(ctx context.Context, streamID string, event interface{}, eventType string) error {
	// Serialize the event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Create a new event
	esEvent := esdb.EventData{
		ContentType: esdb.JsonContentType,
		EventType:   eventType,
		Data:        eventData,
	}

	// Write the event to the stream
	_, err = store.client.AppendToStream(ctx, streamID, esdb.AppendToStreamOptions{}, esEvent)
	if err != nil {
		return fmt.Errorf("failed to append event to stream: %w", err)
	}

	return nil
}

// FlowCreatedEvent defines the structure of the event when a flow is created.
type FlowCreatedEvent struct {
	FlowID      string       `json:"flow_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Steps       []StepConfig `json:"steps"`
	Timestamp   time.Time    `json:"timestamp"`
}

// FlowStepCompletedEvent defines the structure of the event when a step in the flow is completed.
type FlowStepCompletedEvent struct {
	FlowID     string                 `json:"flow_id"`
	StepID     string                 `json:"step_id"`
	StepName   string                 `json:"step_name"`
	Params     map[string]interface{} `json:"params"`
	NextStepID string                 `json:"next_step_id"`
	Timestamp  time.Time              `json:"timestamp"`
}

// StepConfig contains the configuration of a specific step within the flow.
type StepConfig struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	IntegrationID string                 `json:"integration_id"`
	Action        string                 `json:"action"`
	Params        map[string]interface{} `json:"params"`
	NextStepID    string                 `json:"next_step_id"`
}

// FromFlow converts a Flow entity to a FlowCreatedEvent.
func FromFlow(f *flow.Flow) FlowCreatedEvent {
	steps := make([]StepConfig, len(f.Steps))
	for i, step := range f.Steps {
		steps[i] = StepConfig{
			ID:            step.ID,
			Name:          step.Name,
			IntegrationID: step.IntegrationID,
			Action:        step.Action,
			Params:        step.Params,
			NextStepID:    step.NextStepID,
		}
	}

	return FlowCreatedEvent{
		FlowID:      f.Name, // If you have a specific ID field, use it here
		Name:        f.Name,
		Description: f.Description,
		Steps:       steps,
		Timestamp:   time.Now(),
	}
}

// FlowUpdatedEvent defines the structure of the event when a flow is updated.
type FlowUpdatedEvent struct {
	FlowID      string       `json:"flow_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Steps       []StepConfig `json:"steps"`
	Timestamp   time.Time    `json:"timestamp"`
}

// FromUpdatedFlow converts a Flow entity to a FlowUpdatedEvent.
func FromUpdatedFlow(f *flow.Flow) FlowUpdatedEvent {
	steps := make([]StepConfig, len(f.Steps))
	for i, step := range f.Steps {
		steps[i] = StepConfig{
			ID:            step.ID,
			Name:          step.Name,
			IntegrationID: step.IntegrationID,
			Action:        step.Action,
			Params:        step.Params,
			NextStepID:    step.NextStepID,
		}
	}

	return FlowUpdatedEvent{
		FlowID:      f.ID, // Assuming Flow ID is used here
		Name:        f.Name,
		Description: f.Description,
		Steps:       steps,
		Timestamp:   time.Now(),
	}
}

// FlowDeletedEvent defines the structure of the event when a flow is deleted.
type FlowDeletedEvent struct {
	FlowID    string    `json:"flow_id"`
	Timestamp time.Time `json:"timestamp"`
}

// FromDeletedFlow converts a Flow entity to a FlowDeletedEvent.
func FromDeletedFlow(f *flow.Flow) FlowDeletedEvent {
	return FlowDeletedEvent{
		FlowID:    f.ID, // Assuming Flow ID is used here
		Timestamp: time.Now(),
	}
}

// FlowStepFailedEvent defines the structure of the event when a step in the flow fails.
type FlowStepFailedEvent struct {
	FlowID    string                 `json:"flow_id"`
	StepID    string                 `json:"step_id"`
	StepName  string                 `json:"step_name"`
	Error     string                 `json:"error"`
	Params    map[string]interface{} `json:"params"`
	Timestamp time.Time              `json:"timestamp"`
}

// FromFailedStep converts a Flow.Step entity to a FlowStepFailedEvent.
func FromFailedStep(flowID string, step *flow.Step, err error) FlowStepFailedEvent {
	return FlowStepFailedEvent{
		FlowID:    flowID,
		StepID:    step.ID,
		StepName:  step.Name,
		Error:     err.Error(),
		Params:    step.Params,
		Timestamp: time.Now(),
	}
}

// FlowExecutedEvent defines the structure of the event when a flow has been executed successfully.
type FlowExecutedEvent struct {
	FlowID    string    `json:"flow_id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}
