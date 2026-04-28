// Package service implements the business logic for promo-service.
package service

import "context"

// PromoService defines the interface for promo-service business logic.
type PromoService interface {

	// ValidatePromo implements the business logic for ValidatePromo.
	ValidatePromo(ctx context.Context) error

	// ApplyPromo implements the business logic for ApplyPromo.
	ApplyPromo(ctx context.Context) error

	// CreatePromo implements the business logic for CreatePromo.
	CreatePromo(ctx context.Context) error

	// GetPromos implements the business logic for GetPromos.
	GetPromos(ctx context.Context) error
}

// promoServiceImpl is the concrete implementation of PromoService.
type promoServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewPromoService creates a new PromoService.
func NewPromoService() PromoService {
	return &promoServiceImpl{}
}
