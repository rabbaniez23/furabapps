// Package repository provides data access layer for settlement-service.
package repository

import (
	"context"
	"database/sql"

	"furab-backend/services/settlement-service/internal/model"
)

// SettlementRepository defines the interface for settlement-service data access.
type SettlementRepository interface {
	CreateSettlement(ctx context.Context, s *model.Settlement) error
	GetSettlementByPaymentID(ctx context.Context, paymentID string) (*model.Settlement, error)
	UpdateSettlementStatus(ctx context.Context, settlementID string, status model.SettlementStatus) error
}

// postgresSettlementRepository implements SettlementRepository using PostgreSQL.
type postgresSettlementRepository struct {
	db *sql.DB
}

// NewPostgresSettlementRepository creates a new PostgreSQL-based repository.
func NewPostgresSettlementRepository(db *sql.DB) SettlementRepository {
	return &postgresSettlementRepository{db: db}
}

func (r *postgresSettlementRepository) CreateSettlement(ctx context.Context, s *model.Settlement) error {
	// TODO: implement insert settlements
	return nil
}

func (r *postgresSettlementRepository) GetSettlementByPaymentID(ctx context.Context, paymentID string) (*model.Settlement, error) {
	// TODO: implement select settlement by payment_id
	return nil, nil
}

func (r *postgresSettlementRepository) UpdateSettlementStatus(ctx context.Context, settlementID string, status model.SettlementStatus) error {
	// TODO: implement update status
	return nil
}
