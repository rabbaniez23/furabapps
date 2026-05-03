package repository

import (
	"context"

	"furab-backend/services/merchant-service/internal/model"
)

// MerchantRepository defines the interface for merchant data operations.
type MerchantRepository interface {
	Create(ctx context.Context, merchant model.Merchant) error
	Update(ctx context.Context, merchant model.Merchant) error
	GetByID(ctx context.Context, merchantID string) (model.Merchant, error)
	UpdateStatus(ctx context.Context, merchantID string, status string) error
	Deactivate(ctx context.Context, merchantID string) error
	Search(ctx context.Context, filter map[string]interface{}) ([]model.Merchant, error)
	SetStatusCache(ctx context.Context, merchantID string, status string) error
	GetStatusCache(ctx context.Context, merchantID string) (string, error)
}

// postgresMerchantRepository is a dummy implementation of MerchantRepository.
// It is useful to satisfy interfaces during development.
type postgresMerchantRepository struct {
	// TODO: add connection pool or cache dependencies
}

// NewPostgresMerchantRepository creates a new postgresMerchantRepository.
func NewPostgresMerchantRepository() MerchantRepository {
	return &postgresMerchantRepository{}
}

func (r *postgresMerchantRepository) Create(ctx context.Context, merchant model.Merchant) error {
	return nil
}

func (r *postgresMerchantRepository) Update(ctx context.Context, merchant model.Merchant) error {
	return nil
}

func (r *postgresMerchantRepository) GetByID(ctx context.Context, merchantID string) (model.Merchant, error) {
	return model.Merchant{}, nil
}

func (r *postgresMerchantRepository) UpdateStatus(ctx context.Context, merchantID string, status string) error {
	return nil
}

func (r *postgresMerchantRepository) Deactivate(ctx context.Context, merchantID string) error {
	return nil
}

func (r *postgresMerchantRepository) Search(ctx context.Context, filter map[string]interface{}) ([]model.Merchant, error) {
	return nil, nil
}

func (r *postgresMerchantRepository) SetStatusCache(ctx context.Context, merchantID string, status string) error {
	return nil
}

func (r *postgresMerchantRepository) GetStatusCache(ctx context.Context, merchantID string) (string, error) {
	return "", nil
}
