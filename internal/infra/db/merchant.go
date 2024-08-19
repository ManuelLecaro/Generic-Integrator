package db

import (
	"agnostic-payment-platform/internal/application/merchant/ports/repository"
	model "agnostic-payment-platform/internal/domain/merchant"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoMerchantRepository struct {
	collection *mongo.Collection
}

// NewMongoPaymentRepository crea una nueva instancia de MongoMerchantRepository.
func NewMongoMerchantRepository(db *MongoDB) repository.MerchantRepository {
	return &MongoMerchantRepository{
		collection: db.Database.Collection("merchant"),
	}
}

// Save stores a new merchant in the repository.
func (r *MongoMerchantRepository) Save(ctx context.Context, merchant *model.Merchant) error {
	merchant.ID = primitive.NewObjectID().Hex()
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, merchant)
	return err
}

// FindByID retrieves a merchant by its unique ID.
func (r *MongoMerchantRepository) FindByID(ctx context.Context, id string) (*model.Merchant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var merchant model.Merchant
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&merchant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("merchant not found")
		}
		return nil, err
	}

	return &merchant, nil
}

// FindByEmail retrieves a merchant by its email address.
func (r *MongoMerchantRepository) FindByEmail(ctx context.Context, email string) (*model.Merchant, error) {
	var merchant model.Merchant
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&merchant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("merchant not found")
		}
		return nil, err
	}

	return &merchant, nil
}

// Update updates the information of an existing merchant.
func (r *MongoMerchantRepository) Update(ctx context.Context, merchant *model.Merchant) error {
	objectID, err := primitive.ObjectIDFromHex(merchant.ID)
	if err != nil {
		return err
	}

	merchant.UpdatedAt = time.Now()
	update := bson.M{
		"$set": merchant,
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// Delete removes a merchant by its unique ID.
func (r *MongoMerchantRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
