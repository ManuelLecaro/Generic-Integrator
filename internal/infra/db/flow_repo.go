package db

import (
	"context"
	"errors"
	"generic-integration-platform/internal/domain/flow"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FlowRepository defines the interface for flow repository methods.
type FlowRepository interface {
	Create(ctx context.Context, f *flow.Flow) error
	GetByID(ctx context.Context, id string) (*flow.Flow, error)
	GetAll(ctx context.Context) ([]*flow.Flow, error)
	Update(ctx context.Context, f *flow.Flow) error
	Delete(ctx context.Context, id string) error
}

// flowRepo implements FlowRepository interface.
type flowRepo struct {
	collection *mongo.Collection
}

// NewFlowRepository creates a new flow repository.
func NewFlowRepository(mdb *MongoDB) FlowRepository {
	return &flowRepo{
		collection: mdb.Database.Collection("flows"),
	}
}

// Create inserts a new flow into the database.
func (r *flowRepo) Create(ctx context.Context, f *flow.Flow) error {
	if f == nil {
		return errors.New("flow cannot be nil")
	}

	_, err := r.collection.InsertOne(ctx, f)
	return err
}

// GetByID retrieves a flow by its ID.
func (r *flowRepo) GetByID(ctx context.Context, id string) (*flow.Flow, error) {
	var f flow.Flow
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// GetAll retrieves all flows from the database.
func (r *flowRepo) GetAll(ctx context.Context) ([]*flow.Flow, error) {
	var flows []*flow.Flow

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var f flow.Flow
		if err := cursor.Decode(&f); err != nil {
			return nil, err
		}
		flows = append(flows, &f)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return flows, nil
}

// Update modifies an existing flow in the database.
func (r *flowRepo) Update(ctx context.Context, f *flow.Flow) error {
	if f == nil {
		return errors.New("flow cannot be nil")
	}

	filter := bson.M{"_id": f.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, f)
	return err
}

// Delete removes a flow from the database by its ID.
func (r *flowRepo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
