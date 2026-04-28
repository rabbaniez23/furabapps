// Package service implements the business logic for pricing-service.
package service

import "context"

// PriceService defines the interface for pricing-service business logic.
type PriceService interface {

	// EstimatePrice implements the business logic for EstimatePrice.
	EstimatePrice(ctx context.Context) error

	// GetSurgeMultiplier implements the business logic for GetSurgeMultiplier.
	GetSurgeMultiplier(ctx context.Context) error

	// UpdatePriceRule implements the business logic for UpdatePriceRule.
	UpdatePriceRule(ctx context.Context) error
}

// priceServiceImpl is the concrete implementation of PriceService.
type priceServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewPriceService creates a new PriceService.
func NewPriceService() PriceService {
	return &priceServiceImpl{}
}
