// Package service implements the business logic for user-service.
package service

import "context"

// UserService defines the interface for user-service business logic.
type UserService interface {

	// GetProfile implements the business logic for GetProfile.
	GetProfile(ctx context.Context) error

	// UpdateProfile implements the business logic for UpdateProfile.
	UpdateProfile(ctx context.Context) error

	// AddAddress implements the business logic for AddAddress.
	AddAddress(ctx context.Context) error

	// DeleteAddress implements the business logic for DeleteAddress.
	DeleteAddress(ctx context.Context) error
}

// userServiceImpl is the concrete implementation of UserService.
type userServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewUserService creates a new UserService.
func NewUserService() UserService {
	return &userServiceImpl{}
}
