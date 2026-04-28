// Package service implements the business logic for location-service.
package service

import "context"

// LocationService defines the interface for location-service business logic.
type LocationService interface {

	// UpdateLocation implements the business logic for UpdateLocation.
	UpdateLocation(ctx context.Context) error

	// GetNearby implements the business logic for GetNearby.
	GetNearby(ctx context.Context) error

	// TrackDriver implements the business logic for TrackDriver.
	TrackDriver(ctx context.Context) error

	// GetGeoFence implements the business logic for GetGeoFence.
	GetGeoFence(ctx context.Context) error
}

// locationServiceImpl is the concrete implementation of LocationService.
type locationServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewLocationService creates a new LocationService.
func NewLocationService() LocationService {
	return &locationServiceImpl{}
}
