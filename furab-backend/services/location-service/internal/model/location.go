// Package model defines the domain models for location-service.
package model

import "time"

// Location represents the Location model in location-service.
type Location struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Location-specific fields
}

// GeoFence represents the GeoFence model in location-service.
type GeoFence struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add GeoFence-specific fields
}

