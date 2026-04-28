// Package model defines the domain models and DTOs for the ride order service.
package model

import (
	"errors"
	"time"
)

// RideStatus represents the current state of a ride order.
type RideStatus string

const (
	// RideStatusPending indicates the ride order has been created but not yet assigned.
	RideStatusPending RideStatus = "PENDING"
	// RideStatusAssigned indicates a driver has been assigned to the ride.
	RideStatusAssigned RideStatus = "ASSIGNED"
	// RideStatusStarted indicates the ride is in progress.
	RideStatusStarted RideStatus = "STARTED"
	// RideStatusCompleted indicates the ride has been completed successfully.
	RideStatusCompleted RideStatus = "COMPLETED"
	// RideStatusCancelled indicates the ride has been cancelled.
	RideStatusCancelled RideStatus = "CANCELLED"
)

// ValidTransitions defines the allowed state transitions for ride orders.
var ValidTransitions = map[RideStatus][]RideStatus{
	RideStatusPending:   {RideStatusAssigned, RideStatusCancelled},
	RideStatusAssigned:  {RideStatusStarted, RideStatusCancelled},
	RideStatusStarted:   {RideStatusCompleted},
	RideStatusCompleted: {},
	RideStatusCancelled: {},
}

// CanTransitionTo checks if the current status can transition to the target status.
func (s RideStatus) CanTransitionTo(target RideStatus) bool {
	allowed, exists := ValidTransitions[s]
	if !exists {
		return false
	}
	for _, status := range allowed {
		if status == target {
			return true
		}
	}
	return false
}

// IsValid checks if the ride status is a known valid status.
func (s RideStatus) IsValid() bool {
	switch s {
	case RideStatusPending, RideStatusAssigned, RideStatusStarted,
		RideStatusCompleted, RideStatusCancelled:
		return true
	}
	return false
}

// Location represents a geographical coordinate with an address.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// Validate checks if the location has valid coordinates.
func (l *Location) Validate() error {
	if l.Latitude < -90 || l.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if l.Longitude < -180 || l.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	if l.Address == "" {
		return errors.New("address is required")
	}
	return nil
}

// RideOrder represents a ride order in the system.
type RideOrder struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	DriverID        string     `json:"driver_id,omitempty"`
	PickupLocation  Location   `json:"pickup_location"`
	DropoffLocation Location   `json:"dropoff_location"`
	Status          RideStatus `json:"status"`
	Fare            float64    `json:"fare"`
	Distance        float64    `json:"distance"`         // in kilometers
	EstimatedDuration int      `json:"estimated_duration"` // in minutes
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// --- Request DTOs ---

// CreateRideOrderRequest is the request body for creating a new ride order.
type CreateRideOrderRequest struct {
	UserID          string   `json:"user_id" validate:"required"`
	PickupLocation  Location `json:"pickup_location" validate:"required"`
	DropoffLocation Location `json:"dropoff_location" validate:"required"`
}

// Validate validates the create ride order request.
func (r *CreateRideOrderRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	if err := r.PickupLocation.Validate(); err != nil {
		return errors.New("invalid pickup location: " + err.Error())
	}
	if err := r.DropoffLocation.Validate(); err != nil {
		return errors.New("invalid dropoff location: " + err.Error())
	}
	return nil
}

// AssignDriverRequest is the request body for assigning a driver.
type AssignDriverRequest struct {
	DriverID string `json:"driver_id" validate:"required"`
}

// Validate validates the assign driver request.
func (r *AssignDriverRequest) Validate() error {
	if r.DriverID == "" {
		return errors.New("driver_id is required")
	}
	return nil
}

// --- Response DTOs ---

// RideOrderResponse wraps a ride order for API response.
type RideOrderResponse struct {
	Order         *RideOrder `json:"order"`
	EstimatedFare float64    `json:"estimated_fare,omitempty"`
}

// --- Event Payloads ---

// RideCreatedEvent is the payload for the ride.created event.
type RideCreatedEvent struct {
	OrderID         string   `json:"order_id"`
	UserID          string   `json:"user_id"`
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`
	EstimatedFare   float64  `json:"estimated_fare"`
}

// RideAssignedEvent is the payload for the ride.assigned event.
type RideAssignedEvent struct {
	OrderID  string `json:"order_id"`
	DriverID string `json:"driver_id"`
	UserID   string `json:"user_id"`
}

// RideStartedEvent is the payload for the ride.started event.
type RideStartedEvent struct {
	OrderID  string `json:"order_id"`
	DriverID string `json:"driver_id"`
	UserID   string `json:"user_id"`
}

// RideCompletedEvent is the payload for the ride.completed event.
type RideCompletedEvent struct {
	OrderID  string  `json:"order_id"`
	DriverID string  `json:"driver_id"`
	UserID   string  `json:"user_id"`
	Fare     float64 `json:"fare"`
	Distance float64 `json:"distance"`
}
