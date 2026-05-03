// Package service implements the business logic for location-service.
package service

import (
	"context"
	"errors"

	"furab-backend/services/location-service/internal/model"
	"furab-backend/services/location-service/internal/repository"
)

// LocationService defines the interface for location-service business logic.
type LocationService interface {
	UpdateLocation(ctx context.Context, req model.UpdateLocationRequest) error
	UpdateStatus(ctx context.Context, req model.UpdateStatusRequest) error
	SearchNearbyDrivers(ctx context.Context, req model.SearchDriverRequest) ([]model.DriverLocationResponse, error)
	TrackDriver(ctx context.Context, driverID string) (*model.TrackLocationResponse, error)
}

// locationServiceImpl is the concrete implementation of LocationService.
type locationServiceImpl struct {
	repo repository.LocationRepository
}

// NewLocationService creates a new LocationService.
func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationServiceImpl{
		repo: repo,
	}
}

func (s *locationServiceImpl) UpdateLocation(ctx context.Context, req model.UpdateLocationRequest) error {
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		return errors.New("invalid coordinate")
	}
	return s.repo.UpdateLocation(ctx, req)
}

func (s *locationServiceImpl) UpdateStatus(ctx context.Context, req model.UpdateStatusRequest) error {
	if req.DriverStatus != "available" && req.DriverStatus != "busy" {
		return errors.New("invalid status")
	}
	return s.repo.UpdateStatus(ctx, req)
}

func (s *locationServiceImpl) SearchNearbyDrivers(ctx context.Context, req model.SearchDriverRequest) ([]model.DriverLocationResponse, error) {
	if req.LatitudeOrigin < -90 || req.LatitudeOrigin > 90 || req.LongitudeOrigin < -180 || req.LongitudeOrigin > 180 || req.Radius <= 0 {
		return nil, errors.New("invalid input")
	}

	geos, err := s.repo.SearchNearbyDrivers(ctx, req)
	if err != nil {
		return nil, err
	}

	var results []model.DriverLocationResponse
	for _, geo := range geos {
		driverID := geo.Name

		// Check if active (TTL exists)
		isActive, err := s.repo.IsDriverActive(ctx, driverID)
		if err != nil || !isActive {
			continue // skip inactive drivers
		}

		// Check status (available/busy)
		status, err := s.repo.GetStatus(ctx, driverID)
		if err != nil || status == "busy" {
			continue // skip busy drivers
		}

		results = append(results, model.DriverLocationResponse{
			DriverID:     driverID,
			Longitude:    geo.Longitude,
			Latitude:     geo.Latitude,
			Distance:     geo.Dist,
			DriverStatus: status,
		})
	}

	return results, nil
}

func (s *locationServiceImpl) TrackDriver(ctx context.Context, driverID string) (*model.TrackLocationResponse, error) {
	return s.repo.TrackDriver(ctx, driverID)
}
