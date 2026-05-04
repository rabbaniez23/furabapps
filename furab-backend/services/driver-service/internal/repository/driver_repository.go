// Package repository provides data access layer for driver-service.
package repository

import (
	"context"

	"furab-backend/services/driver-service/internal/model"
)

// DriverRepository defines the interface for driver-service data access.
type DriverRepository interface {
	// Save persists a new driver.
	Save(ctx context.Context, driver *model.Driver) error

	// FindByID retrieves a driver by its ID.
	FindByID(ctx context.Context, driverID string) (*model.Driver, error)

	// Update modifies an existing driver record.
	Update(ctx context.Context, driver *model.Driver) error

	// UpdateStatus changes the availability status of a driver.
	UpdateStatus(ctx context.Context, driverID, status string) error

	// UpdateLocation updates the GPS coordinates of a driver.
	UpdateLocation(ctx context.Context, driverID string, lat, long float64) error
}

// postgresDriverRepository implements DriverRepository using PostgreSQL.
type postgresDriverRepository struct {
	// TODO: add *sql.DB / pgxpool field
}

// NewPostgresDriverRepository creates a new PostgreSQL-based repository.
func NewPostgresDriverRepository() DriverRepository {
	return &postgresDriverRepository{}
}

func (r *postgresDriverRepository) Save(ctx context.Context, driver *model.Driver) error {
	// TODO: implement
	return nil
}

func (r *postgresDriverRepository) FindByID(ctx context.Context, driverID string) (*model.Driver, error) {
	// TODO: implement
	return nil, nil
}

func (r *postgresDriverRepository) Update(ctx context.Context, driver *model.Driver) error {
	// TODO: implement
	return nil
}

func (r *postgresDriverRepository) UpdateStatus(ctx context.Context, driverID, status string) error {
	// TODO: implement
	return nil
}

func (r *postgresDriverRepository) UpdateLocation(ctx context.Context, driverID string, lat, long float64) error {
	// TODO: implement
	return nil
}
