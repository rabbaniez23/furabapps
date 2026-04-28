// Package repository provides data access layer for location-service.
package repository

import "context"

// LocationRepository defines the interface for location-service data access.
type LocationRepository interface {

	// UpdateLocation performs the UpdateLocation operation.
	UpdateLocation(ctx context.Context) error

	// GetNearby performs the GetNearby operation.
	GetNearby(ctx context.Context) error

	// TrackDriver performs the TrackDriver operation.
	TrackDriver(ctx context.Context) error

	// GetGeoFence performs the GetGeoFence operation.
	GetGeoFence(ctx context.Context) error
}

// postgresLocationRepository implements LocationRepository using PostgreSQL.
type postgresLocationRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresLocationRepository creates a new PostgreSQL-based repository.
func NewPostgresLocationRepository() LocationRepository {
	return &postgresLocationRepository{}
}
