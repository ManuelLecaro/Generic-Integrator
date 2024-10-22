package db

import (
	"context"
	"errors"
	"generic-integration-platform/internal/domain/integration"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// IntegrationRepository defines the interface for integration repository methods.
type IntegrationRepository interface {
	Create(ctx context.Context, i *integration.Integration) error
	GetByID(ctx context.Context, id string) (*integration.Integration, error)
	GetAll(ctx context.Context) ([]*integration.Integration, error)
	Update(ctx context.Context, i *integration.Integration) error
	Delete(ctx context.Context, id string) error
}

// integrationRepo implements IntegrationRepository interface.
type integrationRepo struct {
	collection *mongo.Collection
}

// NewIntegrationRepository creates a new integration repository.
func NewIntegrationRepository(mdb *MongoDB) IntegrationRepository {
	return &integrationRepo{
		collection: mdb.Database.Collection("integrations"),
	}
}

// Create inserts a new integration into the database.
func (r *integrationRepo) Create(ctx context.Context, i *integration.Integration) error {
	if i == nil {
		return errors.New("integration cannot be nil")
	}

	_, err := r.collection.InsertOne(ctx, i)
	return err
}

// GetByID retrieves an integration by its ID.
func (r *integrationRepo) GetByID(ctx context.Context, id string) (*integration.Integration, error) {
	var i integration.Integration
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// GetAll retrieves all integrations from the database.
func (r *integrationRepo) GetAll(ctx context.Context) ([]*integration.Integration, error) {
	var integrations []*integration.Integration

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var i integration.Integration
		if err := cursor.Decode(&i); err != nil {
			return nil, err
		}
		integrations = append(integrations, &i)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return integrations, nil
}

// Update modifies an existing integration in the database.
func (r *integrationRepo) Update(ctx context.Context, i *integration.Integration) error {
	if i == nil {
		return errors.New("integration cannot be nil")
	}

	filter := bson.M{"_id": i.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, i)
	return err
}

// Delete removes an integration from the database by its ID.
func (r *integrationRepo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
