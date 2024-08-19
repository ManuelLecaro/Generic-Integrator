package service

import (
	"agnostic-payment-platform/internal/application/merchant/ports/repository"
	service "agnostic-payment-platform/internal/application/merchant/ports/services"
	model "agnostic-payment-platform/internal/domain/merchant"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type MerchantService struct {
	repo repository.MerchantRepository
}

func NewMerchantService(repo repository.MerchantRepository) service.MerchantService {
	return &MerchantService{
		repo: repo,
	}
}

// RegisterMerchant registers a new merchant with a hashed password.
func (s *MerchantService) RegisterMerchant(ctx context.Context, email, password, name string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	merchant := &model.Merchant{
		Email:        email,
		PasswordHash: string(passwordHash),
		Name:         name,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return s.repo.Save(ctx, merchant)
}

// LoginMerchant authenticates a merchant by email and password.
func (s *MerchantService) LoginMerchant(ctx context.Context, email, password string) (*model.Merchant, error) {
	merchant, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(merchant.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return merchant, nil
}

// Register registers a new merchant.
func (s *MerchantService) Register(ctx context.Context, email, password, name string) error {
	existingMerchant, _ := s.repo.FindByEmail(ctx, email)
	if existingMerchant != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	merchant := &model.Merchant{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Name:         name,
	}

	return s.repo.Save(ctx, merchant)
}

// Login authenticates a merchant by email and password.
func (s *MerchantService) Login(ctx context.Context, email, password string) (*model.Merchant, error) {
	merchant, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(merchant.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return merchant, nil
}

// GetMerchantByID retrieves a merchant by its unique ID.
func (s *MerchantService) GetMerchantByID(ctx context.Context, id string) (*model.Merchant, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateMerchant updates the details of an existing merchant.
func (s *MerchantService) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	return s.repo.Update(ctx, merchant)
}

// DeleteMerchant removes a merchant by its unique ID.
func (s *MerchantService) DeleteMerchant(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
