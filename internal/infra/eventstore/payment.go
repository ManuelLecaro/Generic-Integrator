package eventstore

import (
	"agnostic-payment-platform/internal/application/payment/ports/repository"
	model "agnostic-payment-platform/internal/domain/payment"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type PaymentEventStore struct {
	es *EventStore
}

func NewPaymentEventStore(es *EventStore) repository.PaymentEventStore {
	return &PaymentEventStore{es: es}
}

// Save guarda un evento de pago en el EventStore.
func (pes *PaymentEventStore) Save(ctx context.Context, event model.PaymentEvent) error {
	streamID := getStreamID(event)

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	esEvent := esdb.EventData{
		EventType:   event.EventType(),
		Data:        eventData,
		ContentType: esdb.JsonContentType,
	}

	_, err = pes.es.DB.AppendToStream(ctx, streamID, esdb.AppendToStreamOptions{}, esEvent)
	return err
}

// GetByID obtiene todos los eventos de pago asociados con un ID de pago.
func (pes *PaymentEventStore) GetByID(ctx context.Context, paymentID string) ([]model.PaymentEvent, error) {
	return pes.getEvents(ctx, paymentID)
}

// GetByTransactionID obtiene todos los eventos de pago asociados con un ID de transacción.
func (pes *PaymentEventStore) GetByTransactionID(ctx context.Context, transactionID string) ([]model.PaymentEvent, error) {
	return pes.getEvents(ctx, transactionID)
}

// getEvents es una función auxiliar para obtener eventos de un stream en EventStoreDB.
func (pes *PaymentEventStore) getEvents(ctx context.Context, streamID string) ([]model.PaymentEvent, error) {
	stream, err := pes.es.DB.ReadStream(ctx, streamID, esdb.ReadStreamOptions{}, 1000)
	if err != nil {
		if errors.Is(err, esdb.ErrStreamNotFound) {
			return nil, nil
		}
		return nil, err
	}
	defer stream.Close()

	var events []model.PaymentEvent
	for {
		resolvedEvent, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) || errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		var event model.PaymentEvent
		switch resolvedEvent.Event.EventType {
		case model.PaymentCreated:
			event = &model.PaymentCreatedEvent{}
		case model.PaymentFailed:
			event = &model.PaymentCreatedEventFailed{}
		case model.PaymentStatusRefunded:
			event = &model.PaymentRefundedEvent{}
		case model.PaymentUpdated:
			event = &model.PaymentStatusUpdatedEvent{}
		default:
			continue
		}

		if err := json.Unmarshal(resolvedEvent.Event.Data, event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

// getStreamID devuelve el ID del stream basado en el tipo de evento.
func getStreamID(event interface{}) string {
	switch e := event.(type) {
	case *model.PaymentCreatedEvent:
		return "payment-" + e.ID
	case *model.PaymentRefundedEvent:
		return "payment-" + e.ID
	case *model.PaymentStatusUpdatedEvent:
		return "payment-" + e.ID
	default:
		return "unknown"
	}
}
