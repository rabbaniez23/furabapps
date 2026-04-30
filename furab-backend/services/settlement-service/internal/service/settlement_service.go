// Package service implements the business logic for settlement-service.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"furab-backend/services/settlement-service/internal/model"
	"furab-backend/services/settlement-service/internal/repository"

	"github.com/google/uuid"
)

type WalletClient interface {
	CreditBalance(ctx context.Context, walletID string, amount float64, referenceID string) error
}

type DriverClient interface {
	GetDriverWalletIDByOrderID(ctx context.Context, orderID string) (string, error)
}

type MerchantClient interface {
	GetMerchantWalletIDByOrderID(ctx context.Context, orderID string) (string, error)
}

// SettlementService defines settlement orchestration.
type SettlementService interface {
	ProcessSettlement(ctx context.Context, req *model.ProcessSettlementRequest) (*model.ProcessSettlementResponse, error)
}

type settlementServiceImpl struct {
	repo        repository.SettlementRepository
	walletCli   WalletClient
	driverCli   DriverClient
	merchantCli MerchantClient
}

func NewSettlementService(
	repo repository.SettlementRepository,
	walletCli WalletClient,
	driverCli DriverClient,
	merchantCli MerchantClient,
) SettlementService {
	return &settlementServiceImpl{
		repo:        repo,
		walletCli:   walletCli,
		driverCli:   driverCli,
		merchantCli: merchantCli,
	}
}

func (s *settlementServiceImpl) ProcessSettlement(ctx context.Context, req *model.ProcessSettlementRequest) (*model.ProcessSettlementResponse, error) {
	if req == nil || req.PaymentID == "" || req.OrderID == "" || req.TotalAmount <= 0 {
		return nil, errors.New("invalid request")
	}

	existing, err := s.repo.GetSettlementByPaymentID(ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &model.ProcessSettlementResponse{
			Status:         mapStatus(existing.Status),
			DriverAmount:   existing.DriverAmount,
			MerchantAmount: existing.MerchantAmount,
			PlatformFee:    existing.PlatformFee,
		}, nil
	}

	driverAmount := req.TotalAmount * 0.80
	merchantAmount := req.TotalAmount * 0.15
	platformFee := req.TotalAmount - driverAmount - merchantAmount

	now := time.Now().UTC()
	settlement := &model.Settlement{
		ID:             uuid.New().String(),
		PaymentID:      req.PaymentID,
		OrderID:        req.OrderID,
		TotalAmount:    req.TotalAmount,
		DriverAmount:   driverAmount,
		MerchantAmount: merchantAmount,
		PlatformFee:    platformFee,
		Status:         model.StatusPending,
		IdempotencyKey: req.PaymentID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.repo.CreateSettlement(ctx, settlement); err != nil {
		return nil, err
	}

	driverWalletID, err := s.driverCli.GetDriverWalletIDByOrderID(ctx, req.OrderID)
	if err != nil {
		_ = s.repo.UpdateSettlementStatus(ctx, settlement.ID, model.StatusFailed)
		return &model.ProcessSettlementResponse{Status: "Failed", DriverAmount: driverAmount, MerchantAmount: merchantAmount, PlatformFee: platformFee}, err
	}
	merchantWalletID, err := s.merchantCli.GetMerchantWalletIDByOrderID(ctx, req.OrderID)
	if err != nil {
		_ = s.repo.UpdateSettlementStatus(ctx, settlement.ID, model.StatusFailed)
		return &model.ProcessSettlementResponse{Status: "Failed", DriverAmount: driverAmount, MerchantAmount: merchantAmount, PlatformFee: platformFee}, err
	}

	if err := s.walletCli.CreditBalance(ctx, driverWalletID, driverAmount, fmt.Sprintf("SETTLE-DRV-%s", req.PaymentID)); err != nil {
		_ = s.repo.UpdateSettlementStatus(ctx, settlement.ID, model.StatusFailed)
		return &model.ProcessSettlementResponse{Status: "Failed", DriverAmount: driverAmount, MerchantAmount: merchantAmount, PlatformFee: platformFee}, err
	}
	if err := s.walletCli.CreditBalance(ctx, merchantWalletID, merchantAmount, fmt.Sprintf("SETTLE-MER-%s", req.PaymentID)); err != nil {
		_ = s.repo.UpdateSettlementStatus(ctx, settlement.ID, model.StatusFailed)
		return &model.ProcessSettlementResponse{Status: "Failed", DriverAmount: driverAmount, MerchantAmount: merchantAmount, PlatformFee: platformFee}, err
	}

	if err := s.repo.UpdateSettlementStatus(ctx, settlement.ID, model.StatusSuccess); err != nil {
		return nil, err
	}

	return &model.ProcessSettlementResponse{
		Status:         "Success",
		DriverAmount:   driverAmount,
		MerchantAmount: merchantAmount,
		PlatformFee:    platformFee,
	}, nil
}

func mapStatus(s model.SettlementStatus) string {
	switch s {
	case model.StatusSuccess:
		return "Success"
	case model.StatusFailed:
		return "Failed"
	default:
		return "Failed"
	}
}
