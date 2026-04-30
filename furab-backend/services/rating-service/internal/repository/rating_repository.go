package repository

import (
	"context"

	"furab-backend/services/rating-service/internal/model"
)

// RatingRepository defines the interface for rating-service data access.
type RatingRepository interface {
	SaveRating(ctx context.Context, rating model.Rating) error
	CheckDuplicate(ctx context.Context, reviewerID, targetID, targetType, orderID string) (bool, string, error)
	GetStatistics(ctx context.Context, targetID, targetType string) (model.RatingSummary, error)
	UpdateStatistics(ctx context.Context, targetID, targetType string, score int) error
	GetHistory(ctx context.Context, reviewerID string, page, limit int) ([]model.Rating, int, error)
}

// postgresRatingRepository implements RatingRepository using PostgreSQL.
type postgresRatingRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresRatingRepository creates a new PostgreSQL-based repository.
func NewPostgresRatingRepository() RatingRepository {
	return &postgresRatingRepository{}
}

// Dummy Implementations to satisfy interface
func (r *postgresRatingRepository) SaveRating(ctx context.Context, rating model.Rating) error {
	return nil
}

func (r *postgresRatingRepository) CheckDuplicate(ctx context.Context, reviewerID, targetID, targetType, orderID string) (bool, string, error) {
	return false, "", nil
}

func (r *postgresRatingRepository) GetStatistics(ctx context.Context, targetID, targetType string) (model.RatingSummary, error) {
	return model.RatingSummary{}, nil
}

func (r *postgresRatingRepository) UpdateStatistics(ctx context.Context, targetID, targetType string, score int) error {
	return nil
}

func (r *postgresRatingRepository) GetHistory(ctx context.Context, reviewerID string, page, limit int) ([]model.Rating, int, error) {
	return nil, 0, nil
}
