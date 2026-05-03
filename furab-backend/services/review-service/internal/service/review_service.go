package service

import (
	"context"
	"errors"

	"furab-backend/services/review-service/internal/model"
	"furab-backend/services/review-service/internal/repository"
)

// ReviewService defines the interface for review business logic.
type ReviewService interface {
	Create(ctx context.Context, review model.Review) error
	GetByTarget(ctx context.Context, targetID, targetType string, page, limit int) ([]model.Review, int, error)
	CreateReport(ctx context.Context, report model.ReviewReport) error
	UpdateStatus(ctx context.Context, reviewID string, status string) error
	GetHistory(ctx context.Context, userID string, targetType string, page, limit int) ([]model.Review, int, error)
}

// reviewServiceImpl is the concrete implementation of ReviewService.
type reviewServiceImpl struct {
	repo repository.ReviewRepository
}

// NewReviewService creates a new ReviewService.
func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewServiceImpl{
		repo: repo,
	}
}

// simulateOrderCheck is a dummy method to simulate checking order completion via another service
func (s *reviewServiceImpl) simulateOrderCheck(orderID string) bool {
	// Simulasi: jika orderID kosong atau mengandung kata 'invalid', anggap order belum selesai
	if orderID == "" || orderID == "invalid_order" {
		return false
	}
	// Default: asumsikan order sudah selesai
	return true
}

func (s *reviewServiceImpl) Create(ctx context.Context, review model.Review) error {
	// 1. Validasi Target Type
	if review.TargetType != "driver" && review.TargetType != "merchant" {
		return errors.New("INVALID_TARGET_TYPE")
	}

	// 2. Simulasi cek status order
	if !s.simulateOrderCheck(review.OrderID) {
		return errors.New("ORDER_NOT_COMPLETED")
	}

	// 3. Cek apakah order ini sudah diulas untuk tipe target yang sama
	existingReview, err := s.repo.GetByOrderID(ctx, review.OrderID, review.TargetType)
	// Jika query berhasil dan id terisi, artinya sudah pernah di-review
	if err == nil && existingReview.ReviewID != "" {
		return errors.New("ALREADY_REVIEWED")
	}

	// 4. Set default status jika belum diisi
	if review.Status == "" {
		review.Status = "active"
	}

	return s.repo.Create(ctx, review)
}

func (s *reviewServiceImpl) GetByTarget(ctx context.Context, targetID, targetType string, page, limit int) ([]model.Review, int, error) {
	// Pagination default fallback
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.GetByTarget(ctx, targetID, targetType, page, limit)
}

func (s *reviewServiceImpl) CreateReport(ctx context.Context, report model.ReviewReport) error {
	// 1. Simpan laporan ke database
	err := s.repo.CreateReport(ctx, report)
	if err != nil {
		return err
	}

	// 2. Ubah status ulasan menjadi 'flagged' sebagai respons pelaporan
	return s.repo.UpdateStatus(ctx, report.ReviewID, "flagged")
}

func (s *reviewServiceImpl) UpdateStatus(ctx context.Context, reviewID string, status string) error {
	return s.repo.UpdateStatus(ctx, reviewID, status)
}

func (s *reviewServiceImpl) GetHistory(ctx context.Context, userID string, targetType string, page, limit int) ([]model.Review, int, error) {
	// Pagination default fallback
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.GetHistory(ctx, userID, targetType, page, limit)
}
