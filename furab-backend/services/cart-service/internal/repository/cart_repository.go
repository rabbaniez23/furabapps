// Package repository provides data access layer for cart-service.
package repository

import "context"

// CartRepository defines the interface for cart-service data access.
type CartRepository interface {

	// AddItem performs the AddItem operation.
	AddItem(ctx context.Context) error

	// RemoveItem performs the RemoveItem operation.
	RemoveItem(ctx context.Context) error

	// UpdateQuantity performs the UpdateQuantity operation.
	UpdateQuantity(ctx context.Context) error

	// GetCart performs the GetCart operation.
	GetCart(ctx context.Context) error

	// ClearCart performs the ClearCart operation.
	ClearCart(ctx context.Context) error
}

// postgresCartRepository implements CartRepository using PostgreSQL.
type postgresCartRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresCartRepository creates a new PostgreSQL-based repository.
func NewPostgresCartRepository() CartRepository {
	return &postgresCartRepository{}
}
