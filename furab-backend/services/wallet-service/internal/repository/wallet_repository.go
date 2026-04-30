// Package repository provides data access layer for wallet-service.
package repository

import (
	"context"
	"database/sql"

	"furab-backend/services/wallet-service/internal/model"
)

// WalletRepository defines the interface for wallet-service data access.
type WalletRepository interface {
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	UpdateBalance(ctx context.Context, walletID string, newBalance float64) error
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
	GetTransactionByReference(ctx context.Context, referenceID string, typ model.TransactionType) (*model.Transaction, error)
}

// postgresWalletRepository implements WalletRepository using PostgreSQL.
type postgresWalletRepository struct {
	db *sql.DB
}

// NewPostgresWalletRepository creates a new PostgreSQL-based repository.
func NewPostgresWalletRepository(db *sql.DB) WalletRepository {
	return &postgresWalletRepository{db: db}
}

func (r *postgresWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	// TODO: implement DB query
	return nil, sql.ErrNoRows
}

func (r *postgresWalletRepository) UpdateBalance(ctx context.Context, walletID string, newBalance float64) error {
	// TODO: implement DB update
	return nil
}

func (r *postgresWalletRepository) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	// TODO: implement DB insert
	return nil
}

func (r *postgresWalletRepository) GetTransactionByReference(ctx context.Context, referenceID string, typ model.TransactionType) (*model.Transaction, error) {
	// TODO: implement DB query
	return nil, nil
}
