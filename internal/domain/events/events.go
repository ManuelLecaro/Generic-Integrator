package events

// DomainEvent represents a generic domain event in the application.
type DomainEvent interface {
	EventType() string // Returns the type of the event
}

// IntegrationCreatedEvent represents an event that is triggered when a new integration is created.
type IntegrationCreatedEvent struct {
	IntegrationID string // The unique identifier of the created integration
}

// EventType returns the type of the IntegrationCreatedEvent.
func (e IntegrationCreatedEvent) EventType() string {
	return "IntegrationCreated"
}

// IntegrationUpdatedEvent represents an event that is triggered when an existing integration is updated.
type IntegrationUpdatedEvent struct {
	IntegrationID string // The unique identifier of the updated integration
}

// EventType returns the type of the IntegrationUpdatedEvent.
func (e IntegrationUpdatedEvent) EventType() string {
	return "IntegrationUpdated"
}

// FlowExecutedEvent represents an event that is triggered when a flow is executed.
type FlowExecutedEvent struct {
	FlowID string // The unique identifier of the executed flow
}

// EventType returns the type of the FlowExecutedEvent.
func (e FlowExecutedEvent) EventType() string {
	return "FlowExecuted"
}

// FlowExecutionFailedEvent represents an event that is triggered when a flow execution fails.
type FlowExecutionFailedEvent struct {
	FlowID string // The unique identifier of the flow that failed
	Error  error  // The error that caused the failure
}

// EventType returns the type of the FlowExecutionFailedEvent.
func (e FlowExecutionFailedEvent) EventType() string {
	return "FlowExecutionFailed"
}
