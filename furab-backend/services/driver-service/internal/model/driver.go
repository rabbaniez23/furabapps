// Package model defines the domain models for driver-service.
package model

import "time"

// DriverProfile represents the DriverProfile model in driver-service.
type DriverProfile struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add DriverProfile-specific fields
}

// DriverLocation represents the DriverLocation model in driver-service.
type DriverLocation struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add DriverLocation-specific fields
}

