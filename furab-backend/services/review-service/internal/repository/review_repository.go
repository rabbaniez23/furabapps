package repository

import (
	"context"

	"furab-backend/services/review-service/internal/model"
)

// ReviewRepository defines the interface for review data operations.
type ReviewRepository interface {
	Create(ctx context.Context, review model.Review) error
	GetByTarget(ctx context.Context, targetID, targetType string, page, limit int) ([]model.Review, int, error)
	GetByOrderID(ctx context.Context, orderID, targetType string) (model.Review, error)
	CreateReport(ctx context.Context, report model.ReviewReport) error
	UpdateStatus(ctx context.Context, reviewID string, status string) error
	GetHistory(ctx context.Context, userID string, targetType string, page, limit int) ([]model.Review, int, error)
}

// postgresReviewRepository is a dummy implementation of ReviewRepository.
// It is useful to satisfy interfaces during development.
type postgresReviewRepository struct {
	// TODO: add connection dependencies (e.g., *sql.DB)
}

// NewPostgresReviewRepository creates a new postgresReviewRepository.
func NewPostgresReviewRepository() ReviewRepository {
	return &postgresReviewRepository{}
}

func (r *postgresReviewRepository) Create(ctx context.Context, review model.Review) error {
	return nil
}

func (r *postgresReviewRepository) GetByTarget(ctx context.Context, targetID, targetType string, page, limit int) ([]model.Review, int, error) {
	return nil, 0, nil
}

func (r *postgresReviewRepository) GetByOrderID(ctx context.Context, orderID, targetType string) (model.Review, error) {
	return model.Review{}, nil
}

func (r *postgresReviewRepository) CreateReport(ctx context.Context, report model.ReviewReport) error {
	return nil
}

func (r *postgresReviewRepository) UpdateStatus(ctx context.Context, reviewID string, status string) error {
	return nil
}

func (r *postgresReviewRepository) GetHistory(ctx context.Context, userID string, targetType string, page, limit int) ([]model.Review, int, error) {
	return nil, 0, nil
}
