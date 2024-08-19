package eventstore

import (
	"agnostic-payment-platform/internal/infra/config"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"go.uber.org/fx"
)

type EventStore struct {
	DB *esdb.Client
}

func NewEventStoreClient(config *config.Config) *EventStore {
	settings, err := esdb.ParseConnectionString(config.EventStore.ConnectionString)

	if err != nil {
		panic(err)
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		panic(err)
	}

	return &EventStore{
		DB: db,
	}
}

var Module = fx.Option(
	fx.Provide(
		NewEventStoreClient,
		NewPaymentEventStore,
	),
)
