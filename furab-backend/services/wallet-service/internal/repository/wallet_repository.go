// Package repository provides data access layer for wallet-service.
package repository

import "context"

// WalletRepository defines the interface for wallet-service data access.
type WalletRepository interface {

	// GetBalance performs the GetBalance operation.
	GetBalance(ctx context.Context) error

	// TopUp performs the TopUp operation.
	TopUp(ctx context.Context) error

	// Debit performs the Debit operation.
	Debit(ctx context.Context) error

	// Transfer performs the Transfer operation.
	Transfer(ctx context.Context) error

	// GetHistory performs the GetHistory operation.
	GetHistory(ctx context.Context) error
}

// postgresWalletRepository implements WalletRepository using PostgreSQL.
type postgresWalletRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresWalletRepository creates a new PostgreSQL-based repository.
func NewPostgresWalletRepository() WalletRepository {
	return &postgresWalletRepository{}
}
