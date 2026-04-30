package repository

import (
	"context"

	"furab-backend/services/audit-log-service/internal/model"
)

// AuditLogRepository defines the interface for audit log data operations.
type AuditLogRepository interface {
	Save(ctx context.Context, log model.AuditLog) error
	GetByID(ctx context.Context, logID string) (model.AuditLog, error)
	Search(ctx context.Context, filter map[string]interface{}, page, limit int) ([]model.AuditLog, int, error)
}

// postgresAuditLogRepository is a dummy implementation of AuditLogRepository.
// It is useful to satisfy interfaces during development.
type postgresAuditLogRepository struct {
	// TODO: add connection dependencies (e.g., *sql.DB)
}

// NewPostgresAuditLogRepository creates a new postgresAuditLogRepository.
func NewPostgresAuditLogRepository() AuditLogRepository {
	return &postgresAuditLogRepository{}
}

func (r *postgresAuditLogRepository) Save(ctx context.Context, log model.AuditLog) error {
	return nil
}

func (r *postgresAuditLogRepository) GetByID(ctx context.Context, logID string) (model.AuditLog, error) {
	return model.AuditLog{}, nil
}

func (r *postgresAuditLogRepository) Search(ctx context.Context, filter map[string]interface{}, page, limit int) ([]model.AuditLog, int, error) {
	return nil, 0, nil
}
