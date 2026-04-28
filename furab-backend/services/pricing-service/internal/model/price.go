// Package model defines the domain models for pricing-service.
package model

import "time"

// PriceEstimate represents the PriceEstimate model in pricing-service.
type PriceEstimate struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add PriceEstimate-specific fields
}

// PriceRule represents the PriceRule model in pricing-service.
type PriceRule struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add PriceRule-specific fields
}

// SurgeZone represents the SurgeZone model in pricing-service.
type SurgeZone struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add SurgeZone-specific fields
}

