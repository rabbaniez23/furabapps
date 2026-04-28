// Package repository provides data access layer for menu-service.
package repository

import "context"

// MenuRepository defines the interface for menu-service data access.
type MenuRepository interface {

	// GetMenu performs the GetMenu operation.
	GetMenu(ctx context.Context) error

	// AddItem performs the AddItem operation.
	AddItem(ctx context.Context) error

	// UpdateItem performs the UpdateItem operation.
	UpdateItem(ctx context.Context) error

	// DeleteItem performs the DeleteItem operation.
	DeleteItem(ctx context.Context) error

	// GetCategories performs the GetCategories operation.
	GetCategories(ctx context.Context) error
}

// postgresMenuRepository implements MenuRepository using PostgreSQL.
type postgresMenuRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresMenuRepository creates a new PostgreSQL-based repository.
func NewPostgresMenuRepository() MenuRepository {
	return &postgresMenuRepository{}
}
