// Package service implements the business logic for wallet-service.
package service

import "context"

// WalletService defines the interface for wallet-service business logic.
type WalletService interface {

	// GetBalance implements the business logic for GetBalance.
	GetBalance(ctx context.Context) error

	// TopUp implements the business logic for TopUp.
	TopUp(ctx context.Context) error

	// Debit implements the business logic for Debit.
	Debit(ctx context.Context) error

	// Transfer implements the business logic for Transfer.
	Transfer(ctx context.Context) error

	// GetHistory implements the business logic for GetHistory.
	GetHistory(ctx context.Context) error
}

// walletServiceImpl is the concrete implementation of WalletService.
type walletServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewWalletService creates a new WalletService.
func NewWalletService() WalletService {
	return &walletServiceImpl{}
}
