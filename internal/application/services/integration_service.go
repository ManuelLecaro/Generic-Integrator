package services

import (
	"context"
	"errors"
	"fmt"
	"generic-integration-platform/internal/application/dto"
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"
	"generic-integration-platform/internal/infra/db"
	"generic-integration-platform/internal/infra/eventstore"
	"generic-integration-platform/internal/infra/executor"
)

var ErrIntegrationNotFound = errors.New("integration not found")

// IntegrationService provides methods for managing integrations.
type IntegrationService struct {
	Repository db.IntegrationRepository
	EventStore eventstore.IntegrationEventStore
	Executor   map[string]executor.Executor
}

// NewIntegrationService creates a new instance of IntegrationService.
func NewIntegrationService(repository db.IntegrationRepository, store eventstore.IntegrationEventStore, executors map[string]executor.Executor) *IntegrationService {
	return &IntegrationService{
		Repository: repository,
		EventStore: store,
		Executor:   executors,
	}
}

// ListIntegrations retrieves all integrations from the repository.
func (s *IntegrationService) ListIntegrations(ctx context.Context) ([]dto.IntegrationResponseDTO, error) {
	integrations, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var integrationResponses []dto.IntegrationResponseDTO
	for _, integration := range integrations {
		integrationResponses = append(integrationResponses, dto.FromDomain(*integration))
	}

	return integrationResponses, nil
}

// CreateIntegration adds a new integration to the repository.
func (s *IntegrationService) CreateIntegration(ctx context.Context, input dto.IntegrationRequestDTO) (dto.IntegrationResponseDTO, error) {
	newIntegration := input.ToDomain()

	if err := s.Repository.Create(ctx, &newIntegration); err != nil {
		return dto.IntegrationResponseDTO{}, err
	}

	if err := s.EventStore.AppendIntegrationCreatedEvent(ctx, eventstore.FromIntegration(&newIntegration)); err != nil {
		return dto.IntegrationResponseDTO{}, err
	}

	return dto.IntegrationResponseDTO{
		ID:       newIntegration.ID,
		Name:     newIntegration.Name,
		Type:     newIntegration.Type,
		BaseURL:  newIntegration.BaseURL,
		AuthType: newIntegration.AuthType,
	}, nil
}

// GetIntegrationByID retrieves a specific integration by its ID.
func (s *IntegrationService) GetIntegrationByID(ctx context.Context, id string) (dto.IntegrationResponseDTO, error) {
	integration, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return dto.IntegrationResponseDTO{}, err
	}

	return dto.IntegrationResponseDTO{
		ID:       integration.ID,
		Name:     integration.Name,
		Type:     integration.Type,
		BaseURL:  integration.BaseURL,
		AuthType: integration.AuthType,
	}, nil
}

// UpdateIntegration updates an existing integration by its ID.
func (s *IntegrationService) UpdateIntegration(ctx context.Context, id string, input dto.IntegrationRequestDTO) (dto.IntegrationResponseDTO, error) {
	updatedIntegration := input.ToDomain()

	if err := s.Repository.Update(ctx, &updatedIntegration); err != nil {
		return dto.IntegrationResponseDTO{}, err
	}

	if err := s.EventStore.AppendIntegrationUpdatedEvent(ctx, eventstore.FromIntegration(&updatedIntegration)); err != nil {
		return dto.IntegrationResponseDTO{}, err
	}

	return dto.IntegrationResponseDTO{
		ID:       updatedIntegration.ID,
		Name:     updatedIntegration.Name,
		Type:     updatedIntegration.Type,
		BaseURL:  updatedIntegration.BaseURL,
		AuthType: updatedIntegration.AuthType,
	}, nil
}

// DeleteIntegration removes an integration by its ID.
func (s *IntegrationService) DeleteIntegration(ctx context.Context, id string) error {
	return s.Repository.Delete(ctx, id)
}

func (s *IntegrationService) RunIntegration(integration integration.Integration, endpoint endpoint.Endpoint) (map[string]interface{}, error) {
	exec, ok := s.Executor[integration.Type]
	if !ok {
		return nil, fmt.Errorf("executor not found for type: %s", integration.Type)
	}

	// Execute the integration with the specific endpoint
	response, err := exec.Execute(integration, endpoint)
	if err != nil {
		return nil, err
	}

	return response, nil
}
