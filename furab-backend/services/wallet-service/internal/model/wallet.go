// Package model defines the domain models for wallet-service.
package model

import "time"

// Wallet represents the Wallet model in wallet-service.
type Wallet struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Wallet-specific fields
}

// Transaction represents the Transaction model in wallet-service.
type Transaction struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Transaction-specific fields
}

