// Package repository provides data access layer for driver-service.
package repository

import "context"

// DriverRepository defines the interface for driver-service data access.
type DriverRepository interface {

	// GetDriver performs the GetDriver operation.
	GetDriver(ctx context.Context) error

	// UpdateLocation performs the UpdateLocation operation.
	UpdateLocation(ctx context.Context) error

	// SetAvailability performs the SetAvailability operation.
	SetAvailability(ctx context.Context) error

	// GetNearbyDrivers performs the GetNearbyDrivers operation.
	GetNearbyDrivers(ctx context.Context) error
}

// postgresDriverRepository implements DriverRepository using PostgreSQL.
type postgresDriverRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresDriverRepository creates a new PostgreSQL-based repository.
func NewPostgresDriverRepository() DriverRepository {
	return &postgresDriverRepository{}
}
