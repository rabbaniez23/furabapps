// Package model defines the domain models for email-service.
package model

import "time"

// EmailRequest represents the EmailRequest model in email-service.
type EmailRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add EmailRequest-specific fields
}

// EmailTemplate represents the EmailTemplate model in email-service.
type EmailTemplate struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add EmailTemplate-specific fields
}

