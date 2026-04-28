// Package service implements the business logic for rating-service.
package service

import "context"

// RatingService defines the interface for rating-service business logic.
type RatingService interface {

	// SubmitRating implements the business logic for SubmitRating.
	SubmitRating(ctx context.Context) error

	// GetAverage implements the business logic for GetAverage.
	GetAverage(ctx context.Context) error

	// GetRatings implements the business logic for GetRatings.
	GetRatings(ctx context.Context) error

	// GetDriverRating implements the business logic for GetDriverRating.
	GetDriverRating(ctx context.Context) error
}

// ratingServiceImpl is the concrete implementation of RatingService.
type ratingServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewRatingService creates a new RatingService.
func NewRatingService() RatingService {
	return &ratingServiceImpl{}
}
