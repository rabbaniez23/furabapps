// Package model defines the domain models for matching-service.
package model

import "time"

// MatchRequest represents the MatchRequest model in matching-service.
type MatchRequest struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add MatchRequest-specific fields
}

// MatchResult represents the MatchResult model in matching-service.
type MatchResult struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add MatchResult-specific fields
}

