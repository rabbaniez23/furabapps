// Package repository provides data access layer for audit-log-service.
package repository

import "context"

// AuditLogRepository defines the interface for audit-log-service data access.
type AuditLogRepository interface {

	// LogAction performs the LogAction operation.
	LogAction(ctx context.Context) error

	// GetLogs performs the GetLogs operation.
	GetLogs(ctx context.Context) error

	// SearchLogs performs the SearchLogs operation.
	SearchLogs(ctx context.Context) error

	// GetLogsByUser performs the GetLogsByUser operation.
	GetLogsByUser(ctx context.Context) error
}

// postgresAuditLogRepository implements AuditLogRepository using PostgreSQL.
type postgresAuditLogRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresAuditLogRepository creates a new PostgreSQL-based repository.
func NewPostgresAuditLogRepository() AuditLogRepository {
	return &postgresAuditLogRepository{}
}
