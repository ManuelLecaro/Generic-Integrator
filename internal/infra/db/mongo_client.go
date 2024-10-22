package db

import (
	"context"
	"fmt"
	"generic-integration-platform/internal/infra/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB struct holds the MongoDB client and database name.
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB client and connects to the database.
func NewMongoDB(config *config.Config) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(config.DB.ConnectionString)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	db := client.Database(config.DB.Name)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// Close closes the MongoDB client connection.
func (mdb *MongoDB) Close(ctx context.Context) error {
	return mdb.Client.Disconnect(ctx)
}
