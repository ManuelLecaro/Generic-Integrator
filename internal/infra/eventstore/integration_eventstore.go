package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"generic-integration-platform/internal/domain/integration"
	"time"

	esdb "github.com/EventStore/EventStore-Client-Go/esdb"
)

// IntegrationEventStore handles storing integration events in EventStoreDB.
type IntegrationEventStore struct {
	client *esdb.Client
}

// NewIntegrationEventStore creates a new instance of IntegrationEventStore.
func NewIntegrationEventStore(client *esdb.Client) *IntegrationEventStore {
	return &IntegrationEventStore{
		client: client,
	}
}

// AppendIntegrationCreatedEvent stores the IntegrationCreated event in EventStore.
func (store *IntegrationEventStore) AppendIntegrationCreatedEvent(ctx context.Context, event IntegrationEvent) error {
	return store.appendEvent(ctx, "integration-"+event.IntegrationID, event, "IntegrationCreatedEvent")
}

// AppendIntegrationUpdatedEvent stores the IntegrationUpdated event in EventStore.
func (store *IntegrationEventStore) AppendIntegrationUpdatedEvent(ctx context.Context, event IntegrationEvent) error {
	return store.appendEvent(ctx, "integration-"+event.IntegrationID, event, "IntegrationUpdatedEvent")
}

// appendEvent serializes and stores a generic event in EventStoreDB.
func (store *IntegrationEventStore) appendEvent(ctx context.Context, streamID string, event interface{}, eventType string) error {
	// Serialize the event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Create a new event
	esEvent := esdb.EventData{
		ContentType: esdb.JsonContentType,
		EventType:   eventType,
		Data:        eventData,
	}

	// Write the event to the stream
	_, err = store.client.AppendToStream(ctx, streamID, esdb.AppendToStreamOptions{}, esEvent)
	if err != nil {
		return fmt.Errorf("failed to append event to stream: %w", err)
	}

	return nil
}

// IntegrationEvent defines the structure of the event when an integration is created.
type IntegrationEvent struct {
	IntegrationID string           `json:"integration_id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	BaseURL       string           `json:"base_url"`
	AuthType      string           `json:"auth_type"`
	AuthToken     string           `json:"auth_token"`
	Endpoints     []EndpointConfig `json:"endpoints"`
	Timestamp     time.Time        `json:"timestamp"`
}

// EndpointConfig contains the configuration of a specific endpoint within the integration.
type EndpointConfig struct {
	Action string            `json:"action"`
	Method string            `json:"method"`
	Path   string            `json:"path"`
	Params map[string]string `json:"params"`
}

// FromIntegration converts an Integration entity to an IntegrationEvent.
func FromIntegration(integration *integration.Integration) IntegrationEvent {
	endpoints := make([]EndpointConfig, len(integration.Endpoints))
	for i, ep := range integration.Endpoints {
		endpoints[i] = EndpointConfig{
			Action: ep.Action,
			Method: ep.Method,
			Path:   ep.Path,
			Params: ep.Params,
		}
	}

	return IntegrationEvent{
		IntegrationID: integration.Name, // If you have a specific ID field, use it here
		Name:          integration.Name,
		Type:          integration.Type,
		BaseURL:       integration.BaseURL,
		AuthType:      integration.AuthType,
		AuthToken:     integration.AuthToken,
		Endpoints:     endpoints,
		Timestamp:     time.Now(),
	}
}
