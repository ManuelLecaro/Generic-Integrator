package integration

import (
	"agnostic-payment-platform/internal/application/payment/ports/services"
	"agnostic-payment-platform/internal/infra/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/fx"
)

// ProvideProcessors creates instances of GenericProcessor for each payment provider.
func ProvideProcessors(cfg *config.IntegrationConfig) (*ProcessorManager, error) {
	processors := make(map[string]services.Processor)
	for _, provider := range cfg.PaymentProviders {
		processors[provider.Name] = NewGenericIntegrationProcessor(&provider)
	}

	return NewProcessorManager(processors), nil
}

type GenericIntegrationProcessor struct {
	providerConfig *config.PaymentProvider
}

func NewGenericIntegrationProcessor(cfg *config.PaymentProvider) *GenericIntegrationProcessor {
	return &GenericIntegrationProcessor{providerConfig: cfg}
}

func (g *GenericIntegrationProcessor) Process(ctx context.Context, action string, params map[string]string) (string, error) {
	endpoint := g.findEndpoint(action)
	if endpoint == nil {
		return "", fmt.Errorf("action %s not supported by provider %s", action, g.providerConfig.Name)
	}

	// Construct the URL with replaced parameters
	url := g.providerConfig.BaseURL + endpoint.Path
	for key, value := range params {
		url = strings.ReplaceAll(url, "{{"+key+"}}", value)
	}

	// Prepare the body
	body := make(map[string]interface{})
	for key, value := range endpoint.Params {
		// Directly replace parameters in the body values
		for paramKey, paramValue := range params {
			placeholder := "{{" + paramKey + "}}"
			value = strings.ReplaceAll(value, placeholder, paramValue)
		}
		body[key] = value
	}

	// Convert the body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body to JSON: %v", err)
	}

	// Create and send the HTTP request
	req, err := http.NewRequest(endpoint.Method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set(g.providerConfig.AuthHeader, g.providerConfig.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract transaction ID
	transactionID := extractTransactionID(result)
	return transactionID, nil
}

func (g *GenericIntegrationProcessor) findEndpoint(action string) *config.EndpointConfig {
	for _, endpoint := range g.providerConfig.Endpoints {
		if endpoint.Action == action {
			return &endpoint
		}
	}
	return nil
}

func extractTransactionID(response map[string]interface{}) string {
	if id, ok := response["transaction_id"].(string); ok {
		return id
	}

	return ""
}

var Module = fx.Options(
	fx.Provide(
		ProvideProcessors,
	),
)
