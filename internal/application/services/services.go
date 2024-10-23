package services

import (
	"context"
	"generic-integration-platform/internal/application/dto"
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"

	"go.uber.org/fx"
)

type IFlowService interface {
	// ListFlows retrieves all flows.
	ListFlows(ctx context.Context) ([]dto.FlowDTO, error)

	// CreateFlow adds a new flow.
	CreateFlow(ctx context.Context, input dto.FlowDTO) (dto.FlowDTO, error)

	// GetFlowByID retrieves a specific flow by its ID.
	GetFlowByID(ctx context.Context, id string) (dto.FlowDTO, error)

	// UpdateFlow updates an existing flow by its ID.
	UpdateFlow(ctx context.Context, id string, input dto.FlowDTO) (dto.FlowDTO, error)

	// DeleteFlow removes a flow by its ID.
	DeleteFlow(ctx context.Context, id string) error

	// ExecuteFlow executes a specific flow by its ID.
	ExecuteFlow(ctx context.Context, id string) (dto.FlowDTO, error)
}

type IIntegrationService interface {
	// ListIntegrations retrieves all integrations.
	ListIntegrations(ctx context.Context) ([]dto.IntegrationResponseDTO, error)

	// CreateIntegration adds a new integration.
	CreateIntegration(ctx context.Context, input dto.IntegrationRequestDTO) (dto.IntegrationResponseDTO, error)

	// GetIntegrationByID retrieves a specific integration by its ID.
	GetIntegrationByID(ctx context.Context, id string) (dto.IntegrationResponseDTO, error)

	// UpdateIntegration updates an existing integration by its ID.
	UpdateIntegration(ctx context.Context, id string, input dto.IntegrationRequestDTO) (dto.IntegrationResponseDTO, error)

	// DeleteIntegration removes an integration by its ID.
	DeleteIntegration(ctx context.Context, id string) error

	// RunIntegration runs the integration.
	RunIntegration(integration integration.Integration, endpoint endpoint.Endpoint) (map[string]interface{}, error)
}

var Module = fx.Provide(
	NewFlowService,
	NewIntegrationService,
)
