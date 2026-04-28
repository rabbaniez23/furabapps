// Package service implements the business logic for merchant-service.
package service

import "context"

// MerchantService defines the interface for merchant-service business logic.
type MerchantService interface {

	// Register implements the business logic for Register.
	Register(ctx context.Context) error

	// GetMerchant implements the business logic for GetMerchant.
	GetMerchant(ctx context.Context) error

	// UpdateProfile implements the business logic for UpdateProfile.
	UpdateProfile(ctx context.Context) error

	// SetOperatingHours implements the business logic for SetOperatingHours.
	SetOperatingHours(ctx context.Context) error
}

// merchantServiceImpl is the concrete implementation of MerchantService.
type merchantServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewMerchantService creates a new MerchantService.
func NewMerchantService() MerchantService {
	return &merchantServiceImpl{}
}
