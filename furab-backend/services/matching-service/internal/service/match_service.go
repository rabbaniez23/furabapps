// Package service implements the business logic for matching-service.
package service

import "context"

// MatchService defines the interface for matching-service business logic.
type MatchService interface {

	// FindDriver implements the business logic for FindDriver.
	FindDriver(ctx context.Context) error

	// AcceptMatch implements the business logic for AcceptMatch.
	AcceptMatch(ctx context.Context) error

	// RejectMatch implements the business logic for RejectMatch.
	RejectMatch(ctx context.Context) error

	// GetMatchStatus implements the business logic for GetMatchStatus.
	GetMatchStatus(ctx context.Context) error
}

// matchServiceImpl is the concrete implementation of MatchService.
type matchServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewMatchService creates a new MatchService.
func NewMatchService() MatchService {
	return &matchServiceImpl{}
}
