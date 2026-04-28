// Package model defines the domain models for merchant-service.
package model

import "time"

// Merchant represents the Merchant model in merchant-service.
type Merchant struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Merchant-specific fields
}

// MerchantProfile represents the MerchantProfile model in merchant-service.
type MerchantProfile struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add MerchantProfile-specific fields
}

