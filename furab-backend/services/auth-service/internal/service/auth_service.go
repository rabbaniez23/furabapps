// Package service implements the business logic for auth-service.
package service

import "context"

// AuthService defines the interface for auth-service business logic.
type AuthService interface {

	// Login implements the business logic for Login.
	Login(ctx context.Context) error

	// Register implements the business logic for Register.
	Register(ctx context.Context) error

	// RefreshToken implements the business logic for RefreshToken.
	RefreshToken(ctx context.Context) error

	// Logout implements the business logic for Logout.
	Logout(ctx context.Context) error
}

// authServiceImpl is the concrete implementation of AuthService.
type authServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewAuthService creates a new AuthService.
func NewAuthService() AuthService {
	return &authServiceImpl{}
}
