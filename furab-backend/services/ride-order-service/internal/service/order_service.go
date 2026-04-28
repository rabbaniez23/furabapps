// Package service implements the business logic for ride order management.
package service

import (
	"context"
	"errors"
	"math"
	"time"

	"furab-backend/services/ride-order-service/internal/model"
	"furab-backend/services/ride-order-service/internal/repository"
	"furab-backend/shared/event"

	"github.com/google/uuid"
)

// Common service errors.
var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidTransition  = errors.New("invalid status transition")
	ErrDriverAlreadyAssigned = errors.New("driver already assigned")
)

// OrderService defines the interface for ride order business logic.
// This interface is used for dependency injection in handlers and can be mocked in tests.
type OrderService interface {
	// CreateOrder creates a new ride order, estimating fare and publishing ride.created event.
	CreateOrder(ctx context.Context, req *model.CreateRideOrderRequest) (*model.RideOrder, error)

	// GetOrder retrieves a ride order by its ID.
	GetOrder(ctx context.Context, id string) (*model.RideOrder, error)

	// AssignDriver assigns a driver to a pending ride order and publishes ride.assigned event.
	AssignDriver(ctx context.Context, orderID, driverID string) (*model.RideOrder, error)

	// StartRide transitions an assigned ride to started status and publishes ride.started event.
	StartRide(ctx context.Context, orderID string) (*model.RideOrder, error)

	// CompleteRide transitions a started ride to completed status and publishes ride.completed event.
	CompleteRide(ctx context.Context, orderID string) (*model.RideOrder, error)

	// CancelRide cancels a pending or assigned ride order.
	CancelRide(ctx context.Context, orderID string) (*model.RideOrder, error)

	// GetUserOrders retrieves all ride orders for a specific user with pagination.
	GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*model.RideOrder, int, error)
}

// orderServiceImpl is the concrete implementation of OrderService.
type orderServiceImpl struct {
	repo      repository.OrderRepository
	publisher event.Publisher
}

// NewOrderService creates a new OrderService with the given dependencies.
func NewOrderService(repo repository.OrderRepository, publisher event.Publisher) OrderService {
	return &orderServiceImpl{
		repo:      repo,
		publisher: publisher,
	}
}

// CreateOrder creates a new ride order.
func (s *orderServiceImpl) CreateOrder(ctx context.Context, req *model.CreateRideOrderRequest) (*model.RideOrder, error) {
	// Validate request
	if req == nil {
		return nil, ErrInvalidRequest
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Calculate estimated fare and distance
	distance := calculateDistance(
		req.PickupLocation.Latitude, req.PickupLocation.Longitude,
		req.DropoffLocation.Latitude, req.DropoffLocation.Longitude,
	)
	fare := estimateFare(distance)
	duration := estimateDuration(distance)

	// Create order
	now := time.Now().UTC()
	order := &model.RideOrder{
		ID:                uuid.New().String(),
		UserID:            req.UserID,
		PickupLocation:    req.PickupLocation,
		DropoffLocation:   req.DropoffLocation,
		Status:            model.RideStatusPending,
		Fare:              fare,
		Distance:          distance,
		EstimatedDuration: duration,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Save to database
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish ride.created event
	evt, err := event.NewEvent(event.TopicRideCreated, "ride-order-service", model.RideCreatedEvent{
		OrderID:         order.ID,
		UserID:          order.UserID,
		PickupLocation:  order.PickupLocation,
		DropoffLocation: order.DropoffLocation,
		EstimatedFare:   order.Fare,
	})
	if err == nil && s.publisher != nil {
		_ = s.publisher.Publish(ctx, event.TopicRideCreated, evt)
	}

	return order, nil
}

// GetOrder retrieves a ride order by its ID.
func (s *orderServiceImpl) GetOrder(ctx context.Context, id string) (*model.RideOrder, error) {
	if id == "" {
		return nil, ErrInvalidRequest
	}

	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}

// AssignDriver assigns a driver to a pending ride order.
func (s *orderServiceImpl) AssignDriver(ctx context.Context, orderID, driverID string) (*model.RideOrder, error) {
	if orderID == "" || driverID == "" {
		return nil, ErrInvalidRequest
	}

	// Get current order
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate status transition: only PENDING -> ASSIGNED
	if !order.Status.CanTransitionTo(model.RideStatusAssigned) {
		return nil, ErrInvalidTransition
	}

	// Check if already assigned
	if order.DriverID != "" {
		return nil, ErrDriverAlreadyAssigned
	}

	// Update order
	order.DriverID = driverID
	order.Status = model.RideStatusAssigned
	order.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish ride.assigned event
	evt, err := event.NewEvent(event.TopicRideAssigned, "ride-order-service", model.RideAssignedEvent{
		OrderID:  order.ID,
		DriverID: driverID,
		UserID:   order.UserID,
	})
	if err == nil && s.publisher != nil {
		_ = s.publisher.Publish(ctx, event.TopicRideAssigned, evt)
	}

	return order, nil
}

// StartRide transitions an assigned ride to started status.
func (s *orderServiceImpl) StartRide(ctx context.Context, orderID string) (*model.RideOrder, error) {
	if orderID == "" {
		return nil, ErrInvalidRequest
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate transition: only ASSIGNED -> STARTED
	if !order.Status.CanTransitionTo(model.RideStatusStarted) {
		return nil, ErrInvalidTransition
	}

	order.Status = model.RideStatusStarted
	order.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish ride.started event
	evt, err := event.NewEvent(event.TopicRideStarted, "ride-order-service", model.RideStartedEvent{
		OrderID:  order.ID,
		DriverID: order.DriverID,
		UserID:   order.UserID,
	})
	if err == nil && s.publisher != nil {
		_ = s.publisher.Publish(ctx, event.TopicRideStarted, evt)
	}

	return order, nil
}

// CompleteRide transitions a started ride to completed status.
func (s *orderServiceImpl) CompleteRide(ctx context.Context, orderID string) (*model.RideOrder, error) {
	if orderID == "" {
		return nil, ErrInvalidRequest
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate transition: only STARTED -> COMPLETED
	if !order.Status.CanTransitionTo(model.RideStatusCompleted) {
		return nil, ErrInvalidTransition
	}

	order.Status = model.RideStatusCompleted
	order.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	// Publish ride.completed event
	evt, err := event.NewEvent(event.TopicRideCompleted, "ride-order-service", model.RideCompletedEvent{
		OrderID:  order.ID,
		DriverID: order.DriverID,
		UserID:   order.UserID,
		Fare:     order.Fare,
		Distance: order.Distance,
	})
	if err == nil && s.publisher != nil {
		_ = s.publisher.Publish(ctx, event.TopicRideCompleted, evt)
	}

	return order, nil
}

// CancelRide cancels a pending or assigned ride order.
func (s *orderServiceImpl) CancelRide(ctx context.Context, orderID string) (*model.RideOrder, error) {
	if orderID == "" {
		return nil, ErrInvalidRequest
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Validate transition: only PENDING/ASSIGNED -> CANCELLED
	if !order.Status.CanTransitionTo(model.RideStatusCancelled) {
		return nil, ErrInvalidTransition
	}

	order.Status = model.RideStatusCancelled
	order.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetUserOrders retrieves all ride orders for a user with pagination.
func (s *orderServiceImpl) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*model.RideOrder, int, error) {
	if userID == "" {
		return nil, 0, ErrInvalidRequest
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	orders, err := s.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// --- Helper Functions ---

// calculateDistance calculates the distance between two coordinates using the Haversine formula.
// Returns distance in kilometers.
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// estimateFare calculates the estimated fare based on distance.
// Base fare: Rp 8,000 + Rp 2,500/km
func estimateFare(distanceKm float64) float64 {
	baseFare := 8000.0
	perKm := 2500.0
	return baseFare + (distanceKm * perKm)
}

// estimateDuration estimates ride duration in minutes based on distance.
// Assumes average speed of 30 km/h in city traffic.
func estimateDuration(distanceKm float64) int {
	avgSpeedKmH := 30.0
	minutes := (distanceKm / avgSpeedKmH) * 60
	return int(math.Ceil(minutes))
}
