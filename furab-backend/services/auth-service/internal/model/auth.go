// Package model defines the domain models for auth-service.
package model

import "time"

// User represents the User model in auth-service.
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add User-specific fields
}

// Token represents the Token model in auth-service.
type Token struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Token-specific fields
}

// Session represents the Session model in auth-service.
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Session-specific fields
}

