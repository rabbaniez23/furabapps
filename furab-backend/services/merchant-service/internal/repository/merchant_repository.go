// Package repository provides data access layer for merchant-service.
package repository

import "context"

// MerchantRepository defines the interface for merchant-service data access.
type MerchantRepository interface {

	// Register performs the Register operation.
	Register(ctx context.Context) error

	// GetMerchant performs the GetMerchant operation.
	GetMerchant(ctx context.Context) error

	// UpdateProfile performs the UpdateProfile operation.
	UpdateProfile(ctx context.Context) error

	// SetOperatingHours performs the SetOperatingHours operation.
	SetOperatingHours(ctx context.Context) error
}

// postgresMerchantRepository implements MerchantRepository using PostgreSQL.
type postgresMerchantRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresMerchantRepository creates a new PostgreSQL-based repository.
func NewPostgresMerchantRepository() MerchantRepository {
	return &postgresMerchantRepository{}
}
