package integration

import (
	"errors"
	"generic-integration-platform/internal/domain/endpoint"
)

// Integration represents a payment integration with a service provider.
type Integration struct {
	ID        string               `json:"id,omitempty"`
	Name      string               // The name of the integration
	Type      string               // The type of integration (e.g., REST, gRPC)
	BaseURL   string               // The base URL for API requests
	AuthType  string               // The type of authentication (e.g., Bearer, Basic)
	AuthToken string               // The authentication token
	Endpoints []*endpoint.Endpoint // List of endpoints associated with this integration
}

// NewIntegration creates a new Integration instance.
func NewIntegration(name, integrationType, baseURL, authType, authToken, currency string, endpoints []*endpoint.Endpoint) *Integration {
	return &Integration{
		Name:      name,
		Type:      integrationType,
		BaseURL:   baseURL,
		AuthType:  authType,
		AuthToken: authToken,
		Endpoints: endpoints,
	}
}

// AddEndpoint adds a new endpoint to the integration.
func (i *Integration) AddEndpoint(endpoint *endpoint.Endpoint) {
	i.Endpoints = append(i.Endpoints, endpoint)
}

// Validate checks if the integration has the necessary fields set.
func (i *Integration) Validate() error {
	if i.Name == "" {
		return errors.New("integration name cannot be empty")
	}
	if i.Type == "" {
		return errors.New("integration type cannot be empty")
	}
	if i.BaseURL == "" {
		return errors.New("base URL cannot be empty")
	}
	if i.AuthType == "" {
		return errors.New("auth type cannot be empty")
	}

	return nil
}
