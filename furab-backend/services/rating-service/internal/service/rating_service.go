// Package service implements the business logic for rating-service.
package service

import (
	"context"
	"errors"

	"furab-backend/services/rating-service/internal/model"
	"furab-backend/services/rating-service/internal/repository"
)

// RatingService defines the interface for rating-service business logic.
type RatingService interface {
	SubmitRating(ctx context.Context, rating model.Rating) error
	GetStatistics(ctx context.Context, targetID, targetType string) (model.RatingSummary, error)
	GetHistory(ctx context.Context, reviewerID string, page, limit int) ([]model.Rating, int, error)
}

// ratingServiceImpl is the concrete implementation of RatingService.
type ratingServiceImpl struct {
	repo repository.RatingRepository
}

// NewRatingService creates a new RatingService.
func NewRatingService(repo repository.RatingRepository) RatingService {
	return &ratingServiceImpl{
		repo: repo,
	}
}

// SubmitRating validates, checks for duplicates, saves rating and updates statistics.
func (s *ratingServiceImpl) SubmitRating(ctx context.Context, rating model.Rating) error {
	// 1. Validasi Score harus 1-5
	if rating.Score < 1 || rating.Score > 5 {
		return errors.New("INVALID_SCORE")
	}

	// 2. Panggil CheckDuplicate
	isDuplicate, _, err := s.repo.CheckDuplicate(ctx, rating.ReviewerID, rating.TargetID, rating.TargetType, rating.OrderID)
	if err != nil {
		return err
	}

	if isDuplicate {
		return errors.New("ALREADY_RATED")
	}

	// 3. Simpan ke repository
	err = s.repo.SaveRating(ctx, rating)
	if err != nil {
		return err
	}

	// 4. UpdateStatistics untuk memperbarui rating_summary
	return s.repo.UpdateStatistics(ctx, rating.TargetID, rating.TargetType, rating.Score)
}

// GetStatistics mengambil data dari rating_summary
func (s *ratingServiceImpl) GetStatistics(ctx context.Context, targetID, targetType string) (model.RatingSummary, error) {
	return s.repo.GetStatistics(ctx, targetID, targetType)
}

// GetHistory mengambil daftar rating dengan sistem pagination
func (s *ratingServiceImpl) GetHistory(ctx context.Context, reviewerID string, page, limit int) ([]model.Rating, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.GetHistory(ctx, reviewerID, page, limit)
}
