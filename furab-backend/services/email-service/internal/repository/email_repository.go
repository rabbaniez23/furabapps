// Package repository provides data access layer for email-service.
package repository

import "context"

// EmailRepository defines the interface for email-service data access.
type EmailRepository interface {

	// SendEmail performs the SendEmail operation.
	SendEmail(ctx context.Context) error

	// SendBulk performs the SendBulk operation.
	SendBulk(ctx context.Context) error

	// GetStatus performs the GetStatus operation.
	GetStatus(ctx context.Context) error
}

// postgresEmailRepository implements EmailRepository using PostgreSQL.
type postgresEmailRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresEmailRepository creates a new PostgreSQL-based repository.
func NewPostgresEmailRepository() EmailRepository {
	return &postgresEmailRepository{}
}
