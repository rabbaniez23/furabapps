// Package repository provides data access layer for email-service.
package repository

import (
	"context"

	"furab-backend/services/email-service/internal/model"
)

// EmailRepository defines the interface for email-service data access.
type EmailRepository interface {
	// SaveEmailLog stores email delivery result for monitoring/audit.
	SaveEmailLog(ctx context.Context, log model.EmailLog) error
}

// postgresEmailRepository implements EmailRepository using PostgreSQL.
type postgresEmailRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresEmailRepository creates a new PostgreSQL-based repository.
func NewPostgresEmailRepository() EmailRepository {
	return &postgresEmailRepository{}
}

// SaveEmailLog stores email logs in database.
func (r *postgresEmailRepository) SaveEmailLog(ctx context.Context, log model.EmailLog) error {
	_ = ctx
	_ = log
	// TODO: implement database persistence.
	return nil
}
