// Package model defines the domain models for food-order-service.
package model

import "time"

// FoodOrder represents the FoodOrder model in food-order-service.
type FoodOrder struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add FoodOrder-specific fields
}

// OrderItem represents the OrderItem model in food-order-service.
type OrderItem struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add OrderItem-specific fields
}

