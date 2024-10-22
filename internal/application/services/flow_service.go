package services

import (
	"context"
	"fmt"
	"generic-integration-platform/internal/application/dto"
	"generic-integration-platform/internal/domain/flow"
	"generic-integration-platform/internal/domain/integration"
	"generic-integration-platform/internal/infra/db"
	"generic-integration-platform/internal/infra/eventstore"
	"time"
)

// FlowService provides methods for managing flows.
type FlowService struct {
	Repository            db.FlowRepository
	IntegrationRepository db.IntegrationRepository
	EventStore            eventstore.FlowEventStore
}

// NewFlowService creates a new instance of FlowService.
func NewFlowService(repository db.FlowRepository, integrationRepo db.IntegrationRepository, es eventstore.FlowEventStore) *FlowService {
	return &FlowService{
		Repository:            repository,
		IntegrationRepository: integrationRepo,
		EventStore:            es,
	}
}

// ListFlows retrieves all flows.
func (s *FlowService) ListFlows(ctx context.Context) ([]dto.FlowDTO, error) {
	flows, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var flowResponses []dto.FlowDTO
	for _, flow := range flows {
		flowResponses = append(flowResponses, dto.FromFlowDomain(flow))
	}

	return flowResponses, nil
}

// CreateFlow adds a new flow.
func (s *FlowService) CreateFlow(ctx context.Context, input dto.FlowDTO) (dto.FlowDTO, error) {
	// Convert the input DTO to a domain model
	newFlow := input.ToDomain()

	// Save the new flow to the repository
	if err := s.Repository.Create(ctx, newFlow); err != nil {
		return dto.FlowDTO{}, err
	}

	// Append flow created event to event store
	if err := s.EventStore.AppendFlowCreatedEvent(ctx, eventstore.FromFlow(newFlow)); err != nil {
		return dto.FlowDTO{}, err
	}

	// Return the created flow as a response DTO
	return dto.FromFlowDomain(newFlow), nil
}

// GetFlowByID retrieves a specific flow by its ID.
func (s *FlowService) GetFlowByID(ctx context.Context, id string) (dto.FlowDTO, error) {
	flow, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return dto.FlowDTO{}, err
	}

	return dto.FromFlowDomain(flow), nil
}

// UpdateFlow updates an existing flow by its ID.
func (s *FlowService) UpdateFlow(ctx context.Context, id string, input dto.FlowDTO) (dto.FlowDTO, error) {
	// Retrieve the flow by ID
	existingFlow, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return dto.FlowDTO{}, err
	}

	// Update the flow with new data
	updatedFlow := input.ToDomain()
	updatedFlow.ID = existingFlow.ID

	// Save the updated flow to the repository
	if err := s.Repository.Update(ctx, updatedFlow); err != nil {
		return dto.FlowDTO{}, err
	}

	// Append flow updated event to event store
	if err := s.EventStore.AppendFlowUpdatedEvent(ctx, eventstore.FromUpdatedFlow(updatedFlow)); err != nil {
		return dto.FlowDTO{}, err
	}

	// Return the updated flow as a response DTO
	return dto.FromFlowDomain(updatedFlow), nil
}

// DeleteFlow removes a flow by its ID.
func (s *FlowService) DeleteFlow(ctx context.Context, id string) error {
	flow, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete the flow from the repository
	if err := s.Repository.Delete(ctx, id); err != nil {
		return err
	}

	// Append flow deleted event to event store
	if err := s.EventStore.AppendFlowDeletedEvent(ctx, eventstore.FromDeletedFlow(flow)); err != nil {
		return err
	}

	return nil
}

// ExecuteFlow executes a specific flow by its ID.
func (s *FlowService) ExecuteFlow(ctx context.Context, id string) (dto.FlowDTO, error) {
	// Retrieve the flow by ID from the repository
	flow, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return dto.FlowDTO{}, fmt.Errorf("failed to retrieve flow by ID: %w", err)
	}

	// Ensure the flow is valid before execution
	if err := flow.Validate(); err != nil {
		return dto.FlowDTO{}, fmt.Errorf("flow validation failed: %w", err)
	}

	// Logic to execute the flow
	// Iterating through each step and executing the action
	for _, step := range flow.Steps {
		// Execute each step (this could involve calling an external service)
		_, err := s.executeStep(ctx, step)
		if err != nil {
			// If a step fails, append the FlowStepFailedEvent
			stepFailedEvent := eventstore.FlowStepFailedEvent{
				FlowID:    flow.ID,
				StepID:    step.ID,
				StepName:  step.Name,
				Error:     err.Error(),
				Params:    step.Params,
				Timestamp: time.Now(),
			}
			_ = s.EventStore.AppendFlowStepFailedEvent(ctx, stepFailedEvent)
			return dto.FlowDTO{}, fmt.Errorf("failed to execute step '%s' in flow '%s': %w", step.Name, flow.Name, err)
		}

		// Append FlowStepCompletedEvent after successful execution of the step
		stepCompletedEvent := eventstore.FlowStepCompletedEvent{
			FlowID:     flow.ID,
			StepID:     step.ID,
			StepName:   step.Name,
			Params:     step.Params,
			NextStepID: step.NextStepID,
			Timestamp:  time.Now(),
		}
		_ = s.EventStore.AppendFlowStepCompletedEvent(ctx, stepCompletedEvent)
	}

	// After executing all steps, append FlowExecutedEvent to the EventStore
	flowExecutedEvent := eventstore.FlowExecutedEvent{
		FlowID:    flow.ID,
		Name:      flow.Name,
		Timestamp: time.Now(),
	}
	if err := s.EventStore.AppendFlowExecutedEvent(ctx, flowExecutedEvent); err != nil {
		return dto.FlowDTO{}, fmt.Errorf("failed to append FlowExecutedEvent: %w", err)
	}

	// Convert the flow back to a FlowDTO to return to the caller
	return dto.FromFlowDomain(flow), nil
}

// executeStep executes a specific step in a flow.
func (s *FlowService) executeStep(ctx context.Context, step *flow.Step) (map[string]interface{}, error) {
	// Retrieve the integration associated with the step
	integration, err := s.IntegrationRepository.GetByID(ctx, step.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve integration for step %s: %w", step.ID, err)
	}

	// Execute the action associated with the step using the integration and params
	result, err := s.performAction(ctx, integration, step.Action, step.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute action for step %s: %w", step.ID, err)
	}

	// If the step is successfully executed, append the FlowStepCompleted event
	completedEvent := eventstore.FlowStepCompletedEvent{
		FlowID:     step.ID,
		StepID:     step.ID,
		StepName:   step.Name,
		Params:     step.Params,
		NextStepID: step.NextStepID,
		Timestamp:  time.Now(),
	}
	if err := s.EventStore.AppendFlowStepCompletedEvent(ctx, completedEvent); err != nil {
		return nil, fmt.Errorf("failed to append FlowStepCompletedEvent for step %s: %w", step.ID, err)
	}

	return result, nil
}

// performAction performs the action associated with a step using the provided integration.
func (s *FlowService) performAction(ctx context.Context, integration *integration.Integration, action string, params map[string]interface{}) (map[string]interface{}, error) {
	// The actual logic of how the integration performs the action depends on your business needs.
	// This could involve making an API call, interacting with an external system, etc.
	// For simplicity, we'll assume this is an HTTP call to the integration's API.

	// Example: Make an HTTP request to the integration's base URL with the action and parameters
	/*response, err := integration.ExecuteAction(ctx, action, params)
	if err != nil {
		return nil, err
	}

	// Process the response and return the result
	return response, nil*/
	return nil, nil
}
