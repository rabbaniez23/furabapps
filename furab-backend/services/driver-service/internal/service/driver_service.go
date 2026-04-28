// Package service implements the business logic for driver-service.
package service

import "context"

// DriverService defines the interface for driver-service business logic.
type DriverService interface {

	// GetDriver implements the business logic for GetDriver.
	GetDriver(ctx context.Context) error

	// UpdateLocation implements the business logic for UpdateLocation.
	UpdateLocation(ctx context.Context) error

	// SetAvailability implements the business logic for SetAvailability.
	SetAvailability(ctx context.Context) error

	// GetNearbyDrivers implements the business logic for GetNearbyDrivers.
	GetNearbyDrivers(ctx context.Context) error
}

// driverServiceImpl is the concrete implementation of DriverService.
type driverServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewDriverService creates a new DriverService.
func NewDriverService() DriverService {
	return &driverServiceImpl{}
}
