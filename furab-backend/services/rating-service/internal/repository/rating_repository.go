// Package repository provides data access layer for rating-service.
package repository

import "context"

// RatingRepository defines the interface for rating-service data access.
type RatingRepository interface {

	// SubmitRating performs the SubmitRating operation.
	SubmitRating(ctx context.Context) error

	// GetAverage performs the GetAverage operation.
	GetAverage(ctx context.Context) error

	// GetRatings performs the GetRatings operation.
	GetRatings(ctx context.Context) error

	// GetDriverRating performs the GetDriverRating operation.
	GetDriverRating(ctx context.Context) error
}

// postgresRatingRepository implements RatingRepository using PostgreSQL.
type postgresRatingRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresRatingRepository creates a new PostgreSQL-based repository.
func NewPostgresRatingRepository() RatingRepository {
	return &postgresRatingRepository{}
}
