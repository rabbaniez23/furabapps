// Package service implements the business logic for wallet-service.
package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/wallet-service/internal/model"
	"furab-backend/services/wallet-service/internal/repository"

	"github.com/google/uuid"
)

// WalletService defines the interface for wallet-service business logic.
type WalletService interface {
	HoldBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error)
	ReleaseBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error)
	DebitBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error)
	CreditBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error)
	Refund(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error)
}

// walletServiceImpl is the concrete implementation of WalletService.
type walletServiceImpl struct {
	repo repository.WalletRepository
}

// NewWalletService creates a new WalletService.
func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletServiceImpl{repo: repo}
}

func (s *walletServiceImpl) HoldBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error) {
	return s.changeBalance(ctx, userID, -amount, amount, referenceID, model.TypeHold)
}

func (s *walletServiceImpl) ReleaseBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error) {
	return s.changeBalance(ctx, userID, amount, amount, referenceID, model.TypeRelease)
}

func (s *walletServiceImpl) DebitBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error) {
	return s.changeBalance(ctx, userID, -amount, amount, referenceID, model.TypeDebit)
}

func (s *walletServiceImpl) CreditBalance(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error) {
	return s.changeBalance(ctx, userID, amount, amount, referenceID, model.TypeCredit)
}

func (s *walletServiceImpl) Refund(ctx context.Context, userID string, amount float64, referenceID string) (*model.WalletResult, error) {
	return s.changeBalance(ctx, userID, amount, amount, referenceID, model.TypeRefund)
}

func (s *walletServiceImpl) changeBalance(
	ctx context.Context,
	userID string,
	delta float64,
	amount float64,
	referenceID string,
	typ model.TransactionType,
) (*model.WalletResult, error) {
	if userID == "" || amount <= 0 {
		return nil, errors.New("invalid request")
	}
	if s.repo == nil {
		return nil, errors.New("missing repository")
	}
	if referenceID != "" {
		existing, err := s.repo.GetTransactionByReference(ctx, referenceID, typ)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return &model.WalletResult{
				Status:         existing.Status,
				CurrentBalance: existing.CurrentBalance,
				TransactionID:  existing.ID,
			}, nil
		}
	}

	w, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	newBal := w.Balance + delta
	if newBal < 0 {
		return &model.WalletResult{
			Status:         model.StatusFailed,
			CurrentBalance: w.Balance,
			TransactionID:  "",
		}, errors.New("insufficient balance")
	}

	if err := s.repo.UpdateBalance(ctx, w.ID, newBal); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		ID:             uuid.New().String(),
		WalletID:       w.ID,
		ReferenceID:    referenceID,
		Type:           typ,
		Amount:         amount,
		Status:         model.StatusSuccess,
		CurrentBalance: newBal,
		CreatedAt:      time.Now().UTC(),
	}
	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return &model.WalletResult{
		Status:         model.StatusSuccess,
		CurrentBalance: newBal,
		TransactionID:  tx.ID,
	}, nil
}
