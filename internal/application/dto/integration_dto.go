package dto

import (
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"
	"strings"
)

// IntegrationDTO represents the Data Transfer Object for an Integration.
type IntegrationDTO struct {
	ID         string        `json:"id"`          // Unique identifier for the integration
	Name       string        `json:"name"`        // Name of the integration
	Type       string        `json:"type"`        // Type of the integration (e.g., REST, gRPC)
	BaseURL    string        `json:"base_url"`    // Base URL for the integration
	AuthHeader string        `json:"auth_header"` // Authorization header type (e.g., "Authorization")
	AuthToken  string        `json:"auth_token"`  // Token used for authentication
	Currency   string        `json:"currency"`    // Currency used in the integration
	Endpoints  []EndpointDTO `json:"endpoints"`   // List of endpoints for the integration
}

// EndpointDTO represents the Data Transfer Object for an endpoint in an integration.
type EndpointDTO struct {
	Action string            `json:"action"` // Action name (e.g., "authorize", "capture")
	Method string            `json:"method"` // HTTP method (e.g., "POST", "GET")
	Path   string            `json:"path"`   // Path for the endpoint
	Params map[string]string `json:"params"` // Parameters for the request, using placeholders
}

// IntegrationRequestDTO represents the request body for creating a new integration.
type IntegrationRequestDTO struct {
	Name      string                `json:"name" binding:"required"`         // Name of the integration
	Type      string                `json:"type" binding:"required"`         // Type of integration (e.g., REST, gRPC)
	BaseURL   string                `json:"base_url" binding:"required,url"` // Base URL for API requests
	AuthType  string                `json:"auth_type"`                       // Type of authentication (e.g., Bearer, Basic)
	AuthToken string                `json:"auth_token,omitempty"`            // The authentication token (optional)
	Currency  string                `json:"currency" binding:"required"`     // Currency for transactions
	Endpoints []*EndpointRequestDTO `json:"endpoints" binding:"required"`    // List of endpoints associated with this integration
}

// EndpointRequestDTO represents the request body for defining an endpoint in an integration.
type EndpointRequestDTO struct {
	Name    string `json:"name" binding:"required"`   // Name of the endpoint
	Method  string `json:"method" binding:"required"` // HTTP method (e.g., GET, POST)
	Path    string `json:"path" binding:"required"`   // Path of the endpoint
	Headers string `json:"headers,omitempty"`         // Additional headers (optional)
}

// ToDomain maps IntegrationRequestDTO to Integration domain model.
func (dto IntegrationRequestDTO) ToDomain() integration.Integration {
	endpoints := make([]*endpoint.Endpoint, len(dto.Endpoints))
	for i, endpointDTO := range dto.Endpoints {
		endpoints[i] = &endpoint.Endpoint{
			Action:  endpointDTO.Name, // Assuming Action corresponds to Name
			Method:  endpointDTO.Method,
			Path:    endpointDTO.Path,
			Headers: map[string]string{}, // You might need to parse headers from string to map
		}
	}

	return integration.Integration{
		Name:      dto.Name,
		Type:      dto.Type,
		BaseURL:   dto.BaseURL,
		AuthType:  dto.AuthType,
		AuthToken: dto.AuthToken,
		Endpoints: endpoints,
	}
}

// ToResponseDTO maps Integration domain model to IntegrationResponseDTO.
func ToResponseDTO(integration integration.Integration) IntegrationResponseDTO {
	endpoints := make([]*EndpointResponseDTO, len(integration.Endpoints))
	for i, endpoint := range integration.Endpoints {
		endpoints[i] = &EndpointResponseDTO{
			Name:   endpoint.Action, // Assuming Action corresponds to Name
			Method: endpoint.Method,
			Path:   endpoint.Path,
		}
	}

	return IntegrationResponseDTO{
		ID:        integration.ID,
		Name:      integration.Name,
		Type:      integration.Type,
		BaseURL:   integration.BaseURL,
		AuthType:  integration.AuthType,
		Endpoints: endpoints,
	}
}

// ToDomain maps EndpointRequestDTO to Endpoint domain model.
func (dto EndpointRequestDTO) ToDomain() endpoint.Endpoint {
	return endpoint.Endpoint{
		Action:  dto.Name, // Assuming Action corresponds to Name
		Method:  dto.Method,
		Path:    dto.Path,
		Headers: parseHeaders(dto.Headers), // Parse string headers if needed
	}
}

// Helper function to parse headers from string to map[string]string
func parseHeaders(headers string) map[string]string {
	result := make(map[string]string)
	if headers == "" {
		return result // Return empty map if no headers
	}

	// Split headers by comma
	pairs := strings.Split(headers, ",")
	for _, pair := range pairs {
		// Split each pair by the first colon to separate key and value
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])   // Trim spaces from key
			value := strings.TrimSpace(kv[1]) // Trim spaces from value
			result[key] = value
		}
	}

	return result
}

// IntegrationResponseDTO represents the response body for an integration.
type IntegrationResponseDTO struct {
	ID        string                 `json:"id"`        // Unique identifier for the integration
	Name      string                 `json:"name"`      // Name of the integration
	Type      string                 `json:"type"`      // Type of integration (e.g., REST, gRPC)
	BaseURL   string                 `json:"base_url"`  // Base URL for API requests
	AuthType  string                 `json:"auth_type"` // Type of authentication (e.g., Bearer, Basic)
	Currency  string                 `json:"currency"`  // Currency for transactions
	Endpoints []*EndpointResponseDTO `json:"endpoints"` // List of endpoints associated with this integration
}

// EndpointResponseDTO represents the response body for an endpoint in an integration.
type EndpointResponseDTO struct {
	Name   string `json:"name"`   // Name of the endpoint
	Method string `json:"method"` // HTTP method (e.g., GET, POST)
	Path   string `json:"path"`   // Path of the endpoint
}

// ToDomain converts IntegrationResponseDTO to Integration domain model.
func (dto IntegrationResponseDTO) ToDomain() integration.Integration {
	endpoints := make([]*endpoint.Endpoint, len(dto.Endpoints))
	for i, endpointDTO := range dto.Endpoints {
		endpoints[i] = &endpoint.Endpoint{
			Action:           endpointDTO.Name, // Assuming the Name is the action
			Method:           endpointDTO.Method,
			Path:             endpointDTO.Path,
			Params:           make(map[string]string), // Initialize empty Params map
			Headers:          make(map[string]string), // Initialize empty Headers map
			ResponseMappings: make(map[string]string), // Initialize empty ResponseMappings map
		}
	}

	return integration.Integration{
		ID:        dto.ID,
		Name:      dto.Name,
		Type:      dto.Type,
		BaseURL:   dto.BaseURL,
		AuthType:  dto.AuthType,
		Endpoints: endpoints,
	}
}

// FromDomain converts Integration domain model to IntegrationResponseDTO.
func FromDomain(integration integration.Integration) IntegrationResponseDTO {
	endpoints := make([]*EndpointResponseDTO, len(integration.Endpoints))
	for i, endpoint := range integration.Endpoints {
		endpoints[i] = &EndpointResponseDTO{
			Name:   endpoint.Action, // Assuming Action is the name
			Method: endpoint.Method,
			Path:   endpoint.Path,
		}
	}

	return IntegrationResponseDTO{
		ID:        integration.ID,
		Name:      integration.Name,
		Type:      integration.Type,
		BaseURL:   integration.BaseURL,
		AuthType:  integration.AuthType,
		Endpoints: endpoints,
	}
}

// ToDomainEndpoint converts EndpointResponseDTO to Endpoint domain model.
func (dto EndpointResponseDTO) ToDomain() endpoint.Endpoint {
	return endpoint.Endpoint{
		Action:           dto.Name, // Assuming Name is the action
		Method:           dto.Method,
		Path:             dto.Path,
		Params:           make(map[string]string), // Initialize empty Params map
		Headers:          make(map[string]string), // Initialize empty Headers map
		ResponseMappings: make(map[string]string), // Initialize empty ResponseMappings map
	}
}

// FromDomainEndpoint converts Endpoint domain model to EndpointResponseDTO.
func FromDomainEndpoint(endpoint endpoint.Endpoint) EndpointResponseDTO {
	return EndpointResponseDTO{
		Name:   endpoint.Action, // Assuming Action is the name
		Method: endpoint.Method,
		Path:   endpoint.Path,
	}
}
