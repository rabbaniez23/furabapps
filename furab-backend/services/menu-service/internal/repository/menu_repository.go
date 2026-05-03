package repository

import (
	"context"

	"furab-backend/services/menu-service/internal/model"
)

// MenuRepository defines the interface for menu data operations.
type MenuRepository interface {
	Create(ctx context.Context, menu model.Menu) error
	Update(ctx context.Context, menu model.Menu) error
	Delete(ctx context.Context, menuID string) error
	UpdateStock(ctx context.Context, menuID string, jumlah int) error
	GetByID(ctx context.Context, menuID string) (model.Menu, error)
	ListByMerchant(ctx context.Context, merchantID string) ([]model.Menu, error)
	SetAvailability(ctx context.Context, menuID string, status bool) error
}

// postgresMenuRepository is a dummy implementation of MenuRepository
type postgresMenuRepository struct {
	// TODO: add *sql.DB field or other connection dependency
}

// NewPostgresMenuRepository creates a new postgresMenuRepository
func NewPostgresMenuRepository() MenuRepository {
	return &postgresMenuRepository{}
}

// Dummy methods to satisfy the MenuRepository interface
func (r *postgresMenuRepository) Create(ctx context.Context, menu model.Menu) error {
	return nil
}

func (r *postgresMenuRepository) Update(ctx context.Context, menu model.Menu) error {
	return nil
}

func (r *postgresMenuRepository) Delete(ctx context.Context, menuID string) error {
	return nil
}

func (r *postgresMenuRepository) UpdateStock(ctx context.Context, menuID string, jumlah int) error {
	return nil
}

func (r *postgresMenuRepository) GetByID(ctx context.Context, menuID string) (model.Menu, error) {
	return model.Menu{}, nil
}

func (r *postgresMenuRepository) ListByMerchant(ctx context.Context, merchantID string) ([]model.Menu, error) {
	return nil, nil
}

func (r *postgresMenuRepository) SetAvailability(ctx context.Context, menuID string, status bool) error {
	return nil
}
