// Package service implements the business logic for driver-service.
package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"furab-backend/services/driver-service/internal/model"
	"furab-backend/services/driver-service/internal/repository"
)

var (
	ErrValidation       = errors.New("validation error")
	ErrDriverIDRequired = errors.New("driver id required")
	ErrNameRequired     = errors.New("name required")
	ErrPhoneRequired    = errors.New("phone required")
	ErrVehicleRequired  = errors.New("vehicle type required")
	ErrStatusRequired   = errors.New("status required")
	ErrStatusInvalid    = errors.New("status invalid")
	ErrLocationInvalid  = errors.New("location invalid")
	ErrDriverNotFound   = errors.New("driver not found")
)

const (
	driverStatusActive   = "ACTIVE"
	driverStatusInactive = "INACTIVE"
	driverStatusOnline   = "ONLINE"
	driverStatusOffline  = "OFFLINE"
	driverStatusBusy     = "BUSY"
)

var validDriverStatuses = map[string]struct{}{
	driverStatusActive:   {},
	driverStatusInactive: {},
	driverStatusOnline:   {},
	driverStatusOffline:  {},
	driverStatusBusy:     {},
}

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

func normalizeInput(v string) string {
	return strings.TrimSpace(v)
}

func validateDriverID(driverID string) error {
	if normalizeInput(driverID) == "" {
		return ErrDriverIDRequired
	}
	return nil
}

func validateCreateRequest(req model.CreateDriverRequest) error {
	if err := validateDriverID(req.DriverID); err != nil {
		return err
	}
	if normalizeInput(req.Name) == "" {
		return ErrNameRequired
	}
	if normalizeInput(req.Phone) == "" {
		return ErrPhoneRequired
	}
	if normalizeInput(req.VehicleType) == "" {
		return ErrVehicleRequired
	}
	return nil
}

func validateUpdateRequest(req model.UpdateDriverRequest) error {
	if normalizeInput(req.Name) == "" || normalizeInput(req.Phone) == "" || normalizeInput(req.VehicleType) == "" {
		return ErrValidation
	}
	return nil
}

func normalizeStatus(status string) string {
	return strings.ToUpper(normalizeInput(status))
}

func validateStatus(status string) (string, error) {
	normalized := normalizeStatus(status)
	if normalized == "" {
		return "", ErrStatusRequired
	}
	if _, ok := validDriverStatuses[normalized]; !ok {
		return "", ErrStatusInvalid
	}
	return normalized, nil
}

func validateLocation(lat, long float64) error {
	if math.IsNaN(lat) || math.IsInf(lat, 0) || math.IsNaN(long) || math.IsInf(long, 0) {
		return ErrLocationInvalid
	}
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return ErrLocationInvalid
	}
	return nil
}

func (s *driverServiceImpl) findDriver(ctx context.Context, driverID string) (*model.Driver, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, ErrDriverNotFound
	}
	return driver, nil
}

func (s *driverServiceImpl) CreateDriver(ctx context.Context, req model.CreateDriverRequest) (*model.DriverResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	req.DriverID = normalizeInput(req.DriverID)
	req.Name = normalizeInput(req.Name)
	req.Phone = normalizeInput(req.Phone)
	req.VehicleType = normalizeInput(req.VehicleType)
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	driver := &model.Driver{
		DriverID:    req.DriverID,
		Name:        req.Name,
		Phone:       req.Phone,
		VehicleType: req.VehicleType,
		Status:      driverStatusOffline,
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	driverID = normalizeInput(driverID)
	if err := validateDriverID(driverID); err != nil {
		return nil, err
	}
	return s.findDriver(ctx, driverID)
}

func (s *driverServiceImpl) UpdateDriver(ctx context.Context, driverID string, req model.UpdateDriverRequest) (*model.DriverResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	driverID = normalizeInput(driverID)
	if err := validateDriverID(driverID); err != nil {
		return nil, err
	}
	req.Name = normalizeInput(req.Name)
	req.Phone = normalizeInput(req.Phone)
	req.VehicleType = normalizeInput(req.VehicleType)
	if err := validateUpdateRequest(req); err != nil {
		return nil, err
	}

	driver, err := s.findDriver(ctx, driverID)
	if err != nil {
		return nil, err
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	driverID = normalizeInput(driverID)
	if err := validateDriverID(driverID); err != nil {
		return nil, err
	}
	normalizedStatus, err := validateStatus(status)
	if err != nil {
		return nil, err
	}
	if _, err := s.findDriver(ctx, driverID); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStatus(ctx, driverID, normalizedStatus); err != nil {
		return nil, err
	}

	return &model.DriverResponse{Status: "success", Message: "status berhasil diperbarui"}, nil
}

func (s *driverServiceImpl) UpdateLocation(ctx context.Context, driverID string, lat, long float64) (*model.DriverResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	driverID = normalizeInput(driverID)
	if err := validateDriverID(driverID); err != nil {
		return nil, err
	}
	if err := validateLocation(lat, long); err != nil {
		return nil, err
	}
	if _, err := s.findDriver(ctx, driverID); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateLocation(ctx, driverID, lat, long); err != nil {
		return nil, err
	}

	return &model.DriverResponse{Status: "success", Message: "lokasi berhasil diperbarui"}, nil
}
