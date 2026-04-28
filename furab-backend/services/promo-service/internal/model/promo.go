// Package model defines the domain models for promo-service.
package model

import "time"

// Promo represents the Promo model in promo-service.
type Promo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Promo-specific fields
}

// PromoUsage represents the PromoUsage model in promo-service.
type PromoUsage struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add PromoUsage-specific fields
}

