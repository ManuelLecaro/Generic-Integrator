package db

import (
	"agnostic-payment-platform/internal/application/payment/ports/repository"
	model "agnostic-payment-platform/internal/domain/payment"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPaymentRepository struct {
	collection          *mongo.Collection
	eventsCollection    *mongo.Collection
	snapshotsCollection *mongo.Collection
}

// NewMongoPaymentRepository crea una nueva instancia de MongoPaymentRepository.
func NewMongoPaymentRepository(db *MongoDB) repository.PaymentRepository {
	return &MongoPaymentRepository{
		collection:          db.Database.Collection("payments"),
		eventsCollection:    db.Database.Collection("payment_events"),
		snapshotsCollection: db.Database.Collection("payment_snapshots"),
	}
}

func (r *MongoPaymentRepository) Save(ctx context.Context, payment *model.Payment) error {
	_, err := r.collection.InsertOne(ctx, payment)
	return err
}

func (r *MongoPaymentRepository) FindByID(ctx context.Context, id string) (*model.Payment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var payment model.Payment
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&payment)
	return &payment, err
}

func (r *MongoPaymentRepository) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.collection.FindOne(ctx, bson.M{"transactionid": transactionID}).Decode(&payment)
	return &payment, err
}

func (r *MongoPaymentRepository) ListPayments(ctx context.Context, merchantID, status string, limit, offset int) ([]*model.Payment, error) {
	filter := bson.M{}
	if merchantID != "" {
		filter["merchant_id"] = merchantID
	}
	if status != "" {
		filter["status"] = status
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*model.Payment
	for cursor.Next(ctx) {
		var payment model.Payment
		if err := cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}

func (r *MongoPaymentRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *MongoPaymentRepository) UpdateStatusByTransactionID(ctx context.Context, transactionID string, status string) error {
	filter := bson.M{"transactionid": transactionID}

	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r *MongoPaymentRepository) SaveEvent(ctx context.Context, event *model.PaymentEvent) error {
	_, err := r.eventsCollection.InsertOne(ctx, event)
	return err
}

func (r *MongoPaymentRepository) ListEvents(ctx context.Context, paymentID string) ([]*model.PaymentEvent, error) {
	cursor, err := r.eventsCollection.Find(ctx, bson.M{"payment_id": paymentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []*model.PaymentEvent
	for cursor.Next(ctx) {
		var event model.PaymentEvent
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

func (r *MongoPaymentRepository) SaveSnapshot(ctx context.Context, snapshot *model.PaymentSnapshot) error {
	_, err := r.snapshotsCollection.InsertOne(ctx, snapshot)
	return err
}

func (r *MongoPaymentRepository) GetLatestSnapshot(ctx context.Context, paymentID string) (*model.PaymentSnapshot, error) {
	var snapshot model.PaymentSnapshot
	err := r.snapshotsCollection.
		FindOne(
			ctx,
			bson.M{"payment_id": paymentID},
			options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}}),
		).
		Decode(&snapshot)
	return &snapshot, err
}
