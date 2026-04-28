// Package model defines the domain models for emergency-service.
package model

import "time"

// EmergencyRequest represents the EmergencyRequest model in emergency-service.
type EmergencyRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add EmergencyRequest-specific fields
}

// SOSAlert represents the SOSAlert model in emergency-service.
type SOSAlert struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add SOSAlert-specific fields
}

