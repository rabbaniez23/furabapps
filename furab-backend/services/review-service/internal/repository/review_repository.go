// Package repository provides data access layer for review-service.
package repository

import "context"

// ReviewRepository defines the interface for review-service data access.
type ReviewRepository interface {

	// SubmitReview performs the SubmitReview operation.
	SubmitReview(ctx context.Context) error

	// GetReviews performs the GetReviews operation.
	GetReviews(ctx context.Context) error

	// FlagReview performs the FlagReview operation.
	FlagReview(ctx context.Context) error

	// GetReviewStats performs the GetReviewStats operation.
	GetReviewStats(ctx context.Context) error
}

// postgresReviewRepository implements ReviewRepository using PostgreSQL.
type postgresReviewRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresReviewRepository creates a new PostgreSQL-based repository.
func NewPostgresReviewRepository() ReviewRepository {
	return &postgresReviewRepository{}
}
