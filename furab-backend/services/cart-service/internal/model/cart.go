// Package model defines the domain models for cart-service.
package model

import "time"

// Cart represents the Cart model in cart-service.
type Cart struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Cart-specific fields
}

// CartItem represents the CartItem model in cart-service.
type CartItem struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add CartItem-specific fields
}

