package validators

import (
	"errors"
	"generic-integration-platform/internal/application/dto"
	"net/url"
	"strings"
)

// ValidateIntegration validates the IntegrationDTO input.
func ValidateIntegration(integration dto.IntegrationDTO) error {
	if strings.TrimSpace(integration.Name) == "" {
		return errors.New("integration name is required")
	}

	if strings.TrimSpace(integration.Type) == "" {
		return errors.New("integration type is required")
	}

	if strings.TrimSpace(integration.BaseURL) == "" {
		return errors.New("base URL is required")
	}

	if !isValidURL(integration.BaseURL) {
		return errors.New("invalid base URL format")
	}

	if len(integration.Endpoints) == 0 {
		return errors.New("at least one endpoint is required")
	}

	for _, endpoint := range integration.Endpoints {
		if err := ValidateEndpoint(endpoint); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEndpoint validates an endpoint in the integration.
func ValidateEndpoint(endpoint dto.EndpointDTO) error {
	if strings.TrimSpace(endpoint.Action) == "" {
		return errors.New("action is required for endpoint")
	}

	if strings.TrimSpace(endpoint.Method) == "" {
		return errors.New("HTTP method is required for endpoint")
	}

	if strings.TrimSpace(endpoint.Path) == "" {
		return errors.New("path is required for endpoint")
	}

	return nil
}

// isValidURL checks if a URL is valid.
func isValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
