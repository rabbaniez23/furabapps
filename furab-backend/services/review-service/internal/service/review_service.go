// Package service implements the business logic for review-service.
package service

import "context"

// ReviewService defines the interface for review-service business logic.
type ReviewService interface {

	// SubmitReview implements the business logic for SubmitReview.
	SubmitReview(ctx context.Context) error

	// GetReviews implements the business logic for GetReviews.
	GetReviews(ctx context.Context) error

	// FlagReview implements the business logic for FlagReview.
	FlagReview(ctx context.Context) error

	// GetReviewStats implements the business logic for GetReviewStats.
	GetReviewStats(ctx context.Context) error
}

// reviewServiceImpl is the concrete implementation of ReviewService.
type reviewServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewReviewService creates a new ReviewService.
func NewReviewService() ReviewService {
	return &reviewServiceImpl{}
}
