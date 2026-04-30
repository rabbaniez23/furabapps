package service

import (
	"context"
	"errors"

	"furab-backend/services/menu-service/internal/model"
	"furab-backend/services/menu-service/internal/repository"
)

// MenuService defines the interface for menu-service business logic.
type MenuService interface {
	Create(ctx context.Context, menu model.Menu) error
	Update(ctx context.Context, menu model.Menu) error
	Delete(ctx context.Context, menuID string) error
	UpdateStock(ctx context.Context, menuID string, jumlah int) error
	GetByID(ctx context.Context, menuID string) (model.Menu, error)
	ListByMerchant(ctx context.Context, merchantID string) ([]model.Menu, error)
	SetAvailability(ctx context.Context, menuID string, status bool) error
}

// menuServiceImpl is the concrete implementation of MenuService.
type menuServiceImpl struct {
	repo repository.MenuRepository
}

// NewMenuService creates a new MenuService.
func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuServiceImpl{
		repo: repo,
	}
}

// validateMenu performs basic validation on Menu data.
func (s *menuServiceImpl) validateMenu(menu model.Menu) error {
	if menu.NamaMenu == "" {
		return errors.New("nama menu tidak boleh kosong")
	}
	if menu.Harga < 0 {
		return errors.New("harga tidak boleh negatif")
	}
	return nil
}

func (s *menuServiceImpl) Create(ctx context.Context, menu model.Menu) error {
	if err := s.validateMenu(menu); err != nil {
		return err
	}
	return s.repo.Create(ctx, menu)
}

func (s *menuServiceImpl) Update(ctx context.Context, menu model.Menu) error {
	if err := s.validateMenu(menu); err != nil {
		return err
	}
	return s.repo.Update(ctx, menu)
}

func (s *menuServiceImpl) Delete(ctx context.Context, menuID string) error {
	return s.repo.Delete(ctx, menuID)
}

func (s *menuServiceImpl) UpdateStock(ctx context.Context, menuID string, jumlah int) error {
	// 1. Dapatkan menu saat ini untuk mengecek stok
	menu, err := s.repo.GetByID(ctx, menuID)
	if err != nil {
		return err
	}

	// 2. Kalkulasi stok baru (jumlah diasumsikan sebagai delta: positif untuk tambah, negatif untuk kurangi)
	newStock := menu.Stok + jumlah

	// 3. Validasi stok tidak boleh negatif
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	// 4. Update stok ke database melalui repository
	err = s.repo.UpdateStock(ctx, menuID, jumlah)
	if err != nil {
		return err
	}

	// 5. UX Enhancement: Auto-set availability
	if newStock == 0 {
		// Jika stok habis, otomatis set ketersediaan menjadi false
		_ = s.repo.SetAvailability(ctx, menuID, false)
	} else if menu.Stok == 0 && newStock > 0 {
		// Jika stok sebelumnya habis lalu diisi kembali, otomatis set menjadi true
		_ = s.repo.SetAvailability(ctx, menuID, true)
	}

	return nil
}

func (s *menuServiceImpl) GetByID(ctx context.Context, menuID string) (model.Menu, error) {
	return s.repo.GetByID(ctx, menuID)
}

func (s *menuServiceImpl) ListByMerchant(ctx context.Context, merchantID string) ([]model.Menu, error) {
	return s.repo.ListByMerchant(ctx, merchantID)
}

func (s *menuServiceImpl) SetAvailability(ctx context.Context, menuID string, status bool) error {
	return s.repo.SetAvailability(ctx, menuID, status)
}
