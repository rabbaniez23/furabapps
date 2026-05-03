// Package repository provides data access layer for promo-service.
package repository

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/promo-service/internal/model"
)

var (
	ErrPromoNotFound      = errors.New("promo not found")
	ErrPromoUsageExceeded = errors.New("promo usage limit exceeded")
)

// PromoRepository defines the interface for promo-service promo data access.
type PromoRepository interface {
	GetPromoByCode(ctx context.Context, promoCode string) (*model.Promo, error)
	IncrementUsage(ctx context.Context, promoID string) error
}

type inMemoryPromoRepository struct {
	promos []*model.Promo
}

// NewInMemoryPromoRepository creates a fixed promo repository for testing and development.
func NewInMemoryPromoRepository() PromoRepository {
	return &inMemoryPromoRepository{
		promos: []*model.Promo{
			{
				PromoID:       "promo-001",
				Code:          "DISKONHEMAT",
				DiscountType:  "percentage",
				DiscountValue: 0.1,
				MinPurchase:   50000,
				MaxDiscount:   20000,
				ExpiryDate:    time.Now().AddDate(0, 1, 0),
				UsageLimit:    100,
				UsageCount:    0,
			},
			{
				PromoID:       "promo-002",
				Code:          "FIXED50",
				DiscountType:  "fixed",
				DiscountValue: 50000,
				MinPurchase:   100000,
				MaxDiscount:   50000,
				ExpiryDate:    time.Now().AddDate(0, 2, 0),
				UsageLimit:    50,
				UsageCount:    0,
			},
		},
	}
}

func (r *inMemoryPromoRepository) GetPromoByCode(ctx context.Context, promoCode string) (*model.Promo, error) {
	for _, promo := range r.promos {
		if promo.Code == promoCode {
			return promo, nil
		}
	}

	return nil, ErrPromoNotFound
}

func (r *inMemoryPromoRepository) IncrementUsage(ctx context.Context, promoID string) error {
	for _, promo := range r.promos {
		if promo.PromoID == promoID {
			if promo.UsageLimit > 0 && promo.UsageCount >= promo.UsageLimit {
				return ErrPromoUsageExceeded
			}
			promo.UsageCount++
			return nil
		}
	}

	return ErrPromoNotFound
}

type postgresPromoRepository struct {
}

// NewPostgresPromoRepository creates a new PostgreSQL-based repository.
func NewPostgresPromoRepository() PromoRepository {
	return &postgresPromoRepository{}
}

func (*postgresPromoRepository) GetPromoByCode(ctx context.Context, promoCode string) (*model.Promo, error) {
	return nil, errors.New("postgres promo repository not implemented")
}

func (*postgresPromoRepository) IncrementUsage(ctx context.Context, promoID string) error {
	return errors.New("postgres promo repository not implemented")
}
