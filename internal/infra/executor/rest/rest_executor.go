package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"generic-integration-platform/internal/domain/endpoint"
	"generic-integration-platform/internal/domain/integration"
	"net/http"
)

// RestExecutor is an implementation of the executor for REST services.
type RestExecutor struct{}

func NewRestExecutor() *RestExecutor {
	return &RestExecutor{}
}

// Execute performs a REST call using the integration and endpoint details.
func (e *RestExecutor) Execute(integration integration.Integration, endpoint endpoint.Endpoint) (map[string]interface{}, error) {
	// Build the endpoint URL
	url := fmt.Sprintf("%s/%s", integration.BaseURL, endpoint.Path)

	// Create the request body
	body, err := json.Marshal(endpoint.Params) // Use the Params map from the Endpoint
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(endpoint.Method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add necessary headers
	req.Header.Set("Authorization", integration.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Process the response
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return responseBody, nil
}
