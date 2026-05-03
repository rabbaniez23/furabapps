// Package model defines the domain models for driver-service.
package model

import "time"

// Driver represents a driver entity.
type Driver struct {
	DriverID    string    `json:"driver_id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	VehicleType string    `json:"vehicle_type"`
	Status      string    `json:"status"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateDriverRequest holds the input for creating a driver.
type CreateDriverRequest struct {
	DriverID    string `json:"driver_id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	VehicleType string `json:"vehicle_type"`
}

// UpdateDriverRequest holds the input for updating a driver.
type UpdateDriverRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	VehicleType string `json:"vehicle_type"`
}

// DriverResponse is a generic response for driver operations.
type DriverResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DriverID string `json:"driver_id,omitempty"`
}
