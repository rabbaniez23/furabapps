// Package model defines the domain models for review-service.
package model

import "time"

// Review represents the Review model in review-service.
type Review struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Review-specific fields
}

