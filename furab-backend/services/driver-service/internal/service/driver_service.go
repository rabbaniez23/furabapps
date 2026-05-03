// Package service implements the business logic for driver-service.
package service

import (
	"context"
	"errors"

	"furab-backend/services/driver-service/internal/model"
	"furab-backend/services/driver-service/internal/repository"
)

// DriverService defines the interface for driver-service business logic.
type DriverService interface {
	// CreateDriver registers a new driver.
	CreateDriver(ctx context.Context, req model.CreateDriverRequest) (*model.DriverResponse, error)

	// GetDriver retrieves a driver by ID.
	GetDriver(ctx context.Context, driverID string) (*model.Driver, error)

	// UpdateDriver updates driver profile data.
	UpdateDriver(ctx context.Context, driverID string, req model.UpdateDriverRequest) (*model.DriverResponse, error)

	// UpdateStatus changes a driver's availability status.
	UpdateStatus(ctx context.Context, driverID, status string) (*model.DriverResponse, error)

	// UpdateLocation updates a driver's GPS coordinates.
	UpdateLocation(ctx context.Context, driverID string, lat, long float64) (*model.DriverResponse, error)
}

// driverServiceImpl is the concrete implementation of DriverService.
type driverServiceImpl struct {
	repo repository.DriverRepository
}

// NewDriverService creates a new DriverService.
func NewDriverService(repo repository.DriverRepository) DriverService {
	return &driverServiceImpl{repo: repo}
}

func (s *driverServiceImpl) CreateDriver(ctx context.Context, req model.CreateDriverRequest) (*model.DriverResponse, error) {
	if req.Name == "" || req.Phone == "" || req.VehicleType == "" {
		return nil, errors.New("validation error")
	}

	driver := &model.Driver{
		DriverID:    req.DriverID,
		Name:        req.Name,
		Phone:       req.Phone,
		VehicleType: req.VehicleType,
		Status:      "offline",
	}

	if err := s.repo.Save(ctx, driver); err != nil {
		return nil, err
	}

	return &model.DriverResponse{
		Status:   "success",
		Message:  "driver berhasil dibuat",
		DriverID: driver.DriverID,
	}, nil
}

func (s *driverServiceImpl) GetDriver(ctx context.Context, driverID string) (*model.Driver, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}
	return driver, nil
}

func (s *driverServiceImpl) UpdateDriver(ctx context.Context, driverID string, req model.UpdateDriverRequest) (*model.DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	driver.Name = req.Name
	driver.Phone = req.Phone
	driver.VehicleType = req.VehicleType

	if err := s.repo.Update(ctx, driver); err != nil {
		return nil, err
	}

	return &model.DriverResponse{Status: "success", Message: "update berhasil"}, nil
}

func (s *driverServiceImpl) UpdateStatus(ctx context.Context, driverID, status string) (*model.DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	if err := s.repo.UpdateStatus(ctx, driverID, status); err != nil {
		return nil, err
	}

	return &model.DriverResponse{Status: "success", Message: "status berhasil diperbarui"}, nil
}

func (s *driverServiceImpl) UpdateLocation(ctx context.Context, driverID string, lat, long float64) (*model.DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	if err := s.repo.UpdateLocation(ctx, driverID, lat, long); err != nil {
		return nil, err
	}

	return &model.DriverResponse{Status: "success", Message: "lokasi berhasil diperbarui"}, nil
}
