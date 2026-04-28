// Package repository provides data access layer for pricing-service.
package repository

import "context"

// PriceRepository defines the interface for pricing-service data access.
type PriceRepository interface {

	// EstimatePrice performs the EstimatePrice operation.
	EstimatePrice(ctx context.Context) error

	// GetSurgeMultiplier performs the GetSurgeMultiplier operation.
	GetSurgeMultiplier(ctx context.Context) error

	// UpdatePriceRule performs the UpdatePriceRule operation.
	UpdatePriceRule(ctx context.Context) error
}

// postgresPriceRepository implements PriceRepository using PostgreSQL.
type postgresPriceRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresPriceRepository creates a new PostgreSQL-based repository.
func NewPostgresPriceRepository() PriceRepository {
	return &postgresPriceRepository{}
}
