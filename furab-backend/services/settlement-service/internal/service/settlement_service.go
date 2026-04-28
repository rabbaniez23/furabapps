// Package service implements the business logic for settlement-service.
package service

import "context"

// SettlementService defines the interface for settlement-service business logic.
type SettlementService interface {

	// CreateSettlement implements the business logic for CreateSettlement.
	CreateSettlement(ctx context.Context) error

	// ProcessSettlement implements the business logic for ProcessSettlement.
	ProcessSettlement(ctx context.Context) error

	// GetSettlement implements the business logic for GetSettlement.
	GetSettlement(ctx context.Context) error
}

// settlementServiceImpl is the concrete implementation of SettlementService.
type settlementServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewSettlementService creates a new SettlementService.
func NewSettlementService() SettlementService {
	return &settlementServiceImpl{}
}
