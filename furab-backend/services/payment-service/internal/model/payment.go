// Package model defines the domain models for payment-service.
package model

import "time"

// Payment represents the Payment model in payment-service.
type Payment struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Payment-specific fields
}

// PaymentMethod represents the PaymentMethod model in payment-service.
type PaymentMethod struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add PaymentMethod-specific fields
}

