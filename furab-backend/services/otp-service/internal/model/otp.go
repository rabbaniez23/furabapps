// Package model defines the domain models for otp-service.
package model

import "time"

// OTPRequest represents the OTPRequest model in otp-service.
type OTPRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add OTPRequest-specific fields
}

// OTPVerification represents the OTPVerification model in otp-service.
type OTPVerification struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add OTPVerification-specific fields
}

