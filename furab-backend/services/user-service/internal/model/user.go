// Package model defines the domain models for user-service.
package model

import "time"

// UserProfile represents the UserProfile model in user-service.
type UserProfile struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add UserProfile-specific fields
}

// UserAddress represents the UserAddress model in user-service.
type UserAddress struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents the core user entity.
type User struct {
	UserID string
	Name   string
	Email  string
	Phone  string
	Status string
}
