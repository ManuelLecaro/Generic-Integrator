package endpoint

import (
	"errors"
)

// Endpoint represents a specific endpoint in an integration.
type Endpoint struct {
	Action           string            // Action to execute (e.g., authorize, capture, refund)
	Method           string            // HTTP method (e.g., POST, GET)
	Path             string            // The path of the endpoint
	Params           map[string]string // Parameters for the request (e.g., amount, currency)
	Headers          map[string]string // Additional headers for the request
	ResponseMappings map[string]string // Mappings for the response fields
}

// NewEndpoint creates a new Endpoint instance.
func NewEndpoint(action, method, path string, params, headers, responseMappings map[string]string) *Endpoint {
	return &Endpoint{
		Action:           action,
		Method:           method,
		Path:             path,
		Params:           params,
		Headers:          headers,
		ResponseMappings: responseMappings,
	}
}

// Validate checks if the endpoint has the necessary fields set.
func (e *Endpoint) Validate() error {
	if e.Action == "" {
		return errors.New("endpoint action cannot be empty")
	}
	if e.Method == "" {
		return errors.New("HTTP method cannot be empty")
	}
	if e.Path == "" {
		return errors.New("path cannot be empty")
	}
	return nil
}
