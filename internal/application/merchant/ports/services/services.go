package service

import (
	model "agnostic-payment-platform/internal/domain/merchant"
	"context"
)

// MerchantService defines the operations for managing merchants.
type MerchantService interface {
	// Register registers a new merchant.
	Register(ctx context.Context, email, password, name string) error

	// Login authenticates a merchant by email and password.
	Login(ctx context.Context, email, password string) (*model.Merchant, error)

	// GetMerchantByID retrieves a merchant by its unique ID.
	GetMerchantByID(ctx context.Context, id string) (*model.Merchant, error)

	// UpdateMerchant updates the details of an existing merchant.
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error

	// DeleteMerchant removes a merchant by its unique ID.
	DeleteMerchant(ctx context.Context, id string) error
}
