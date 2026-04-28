// Package repository provides data access layer for promo-service.
package repository

import "context"

// PromoRepository defines the interface for promo-service data access.
type PromoRepository interface {

	// ValidatePromo performs the ValidatePromo operation.
	ValidatePromo(ctx context.Context) error

	// ApplyPromo performs the ApplyPromo operation.
	ApplyPromo(ctx context.Context) error

	// CreatePromo performs the CreatePromo operation.
	CreatePromo(ctx context.Context) error

	// GetPromos performs the GetPromos operation.
	GetPromos(ctx context.Context) error
}

// postgresPromoRepository implements PromoRepository using PostgreSQL.
type postgresPromoRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresPromoRepository creates a new PostgreSQL-based repository.
func NewPostgresPromoRepository() PromoRepository {
	return &postgresPromoRepository{}
}
