// Package model defines the domain models for settlement-service.
package model

import "time"

// Settlement represents the Settlement model in settlement-service.
type Settlement struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Settlement-specific fields
}

// SettlementItem represents the SettlementItem model in settlement-service.
type SettlementItem struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add SettlementItem-specific fields
}

