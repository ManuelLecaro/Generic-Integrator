package repository

import (
	model "agnostic-payment-platform/internal/domain/merchant"
	"context"
)

// MerchantRepository defines the operations for managing merchants in the storage layer.
type MerchantRepository interface {
	// Save stores a new merchant in the repository.
	Save(ctx context.Context, merchant *model.Merchant) error

	// FindByID retrieves a merchant by its unique ID.
	FindByID(ctx context.Context, id string) (*model.Merchant, error)

	// FindByEmail retrieves a merchant by its email address.
	FindByEmail(ctx context.Context, email string) (*model.Merchant, error)

	// Update updates the information of an existing merchant.
	Update(ctx context.Context, merchant *model.Merchant) error

	// Delete removes a merchant by its unique ID.
	Delete(ctx context.Context, id string) error
}
