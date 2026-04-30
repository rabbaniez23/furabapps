package service

import (
	"context"
	"errors"

	"furab-backend/services/merchant-service/internal/model"
	"furab-backend/services/merchant-service/internal/repository"
)

// MerchantService defines the interface for merchant business logic.
type MerchantService interface {
	Create(ctx context.Context, merchant model.Merchant) error
	Update(ctx context.Context, merchant model.Merchant) error
	GetByID(ctx context.Context, merchantID string) (model.Merchant, error)
	UpdateStatus(ctx context.Context, merchantID string, status string) error
	Deactivate(ctx context.Context, merchantID string) error
	Search(ctx context.Context, keyword string, kategori string) ([]model.Merchant, error)
	CheckMerchantStatus(ctx context.Context, merchantID string) (bool, error)
}

// merchantServiceImpl is the concrete implementation of MerchantService.
type merchantServiceImpl struct {
	repo repository.MerchantRepository
}

// NewMerchantService creates a new MerchantService.
func NewMerchantService(repo repository.MerchantRepository) MerchantService {
	return &merchantServiceImpl{
		repo: repo,
	}
}

func (s *merchantServiceImpl) Create(ctx context.Context, merchant model.Merchant) error {
	if merchant.UserID == "" {
		return errors.New("user_id tidak boleh kosong")
	}
	if merchant.NamaToko == "" {
		return errors.New("nama toko tidak boleh kosong")
	}
	// Initial default values
	merchant.IsActive = true
	if merchant.StatusOperasional == "" {
		merchant.StatusOperasional = "closed"
	}
	return s.repo.Create(ctx, merchant)
}

func (s *merchantServiceImpl) Update(ctx context.Context, merchant model.Merchant) error {
	if merchant.NamaToko == "" {
		return errors.New("nama toko tidak boleh kosong")
	}
	return s.repo.Update(ctx, merchant)
}

func (s *merchantServiceImpl) GetByID(ctx context.Context, merchantID string) (model.Merchant, error) {
	return s.repo.GetByID(ctx, merchantID)
}

func (s *merchantServiceImpl) UpdateStatus(ctx context.Context, merchantID string, status string) error {
	if status != "open" && status != "closed" {
		return errors.New("status tidak valid")
	}

	err := s.repo.UpdateStatus(ctx, merchantID, status)
	if err != nil {
		return err
	}

	// Sync dengan cache setelah berhasil di database
	_ = s.repo.SetStatusCache(ctx, merchantID, status)

	return nil
}

func (s *merchantServiceImpl) Deactivate(ctx context.Context, merchantID string) error {
	// 1. Set IsActive = false di database
	err := s.repo.Deactivate(ctx, merchantID)
	if err != nil {
		return err
	}

	// 2. Otomatis set StatusOperasional = closed
	err = s.repo.UpdateStatus(ctx, merchantID, "closed")
	if err != nil {
		return err
	}

	// 3. Update status closed ke Redis/Cache
	_ = s.repo.SetStatusCache(ctx, merchantID, "closed")

	return nil
}

func (s *merchantServiceImpl) Search(ctx context.Context, keyword string, kategori string) ([]model.Merchant, error) {
	filter := make(map[string]interface{})
	if keyword != "" {
		filter["keyword"] = keyword
	}
	if kategori != "" {
		filter["kategori"] = kategori
	}

	return s.repo.Search(ctx, filter)
}

func (s *merchantServiceImpl) CheckMerchantStatus(ctx context.Context, merchantID string) (bool, error) {
	// Optimization: Cek cache terlebih dahulu. 
	// Jika cache mengatakan closed, return false tanpa ke DB.
	cachedStatus, err := s.repo.GetStatusCache(ctx, merchantID)
	if err == nil && cachedStatus == "closed" {
		return false, nil
	}

	// Jika tidak di cache atau statusnya open, butuh verifikasi IsActive ke DB
	merchant, err := s.repo.GetByID(ctx, merchantID)
	if err != nil {
		return false, err
	}

	// Syarat aktif: IsActive == true DAN StatusOperasional == "open"
	isOpen := merchant.IsActive && merchant.StatusOperasional == "open"
	return isOpen, nil
}
