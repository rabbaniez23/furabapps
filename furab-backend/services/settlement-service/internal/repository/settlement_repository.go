// Package repository provides data access layer for settlement-service.
package repository

import "context"

// SettlementRepository defines the interface for settlement-service data access.
type SettlementRepository interface {

	// CreateSettlement performs the CreateSettlement operation.
	CreateSettlement(ctx context.Context) error

	// ProcessSettlement performs the ProcessSettlement operation.
	ProcessSettlement(ctx context.Context) error

	// GetSettlement performs the GetSettlement operation.
	GetSettlement(ctx context.Context) error
}

// postgresSettlementRepository implements SettlementRepository using PostgreSQL.
type postgresSettlementRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresSettlementRepository creates a new PostgreSQL-based repository.
func NewPostgresSettlementRepository() SettlementRepository {
	return &postgresSettlementRepository{}
}
