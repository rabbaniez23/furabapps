// Package model defines the domain models for rating-service.
package model

import "time"

// Rating represents the Rating model in rating-service.
type Rating struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Rating-specific fields
}

// RatingStats represents the RatingStats model in rating-service.
type RatingStats struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add RatingStats-specific fields
}

