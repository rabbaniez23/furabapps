// Package repository provides data access layer for user-service.
package repository

import "context"

// UserRepository defines the interface for user-service data access.
type UserRepository interface {

	// GetProfile performs the GetProfile operation.
	GetProfile(ctx context.Context) error

	// UpdateProfile performs the UpdateProfile operation.
	UpdateProfile(ctx context.Context) error

	// AddAddress performs the AddAddress operation.
	AddAddress(ctx context.Context) error

	// DeleteAddress performs the DeleteAddress operation.
	DeleteAddress(ctx context.Context) error
}

// postgresUserRepository implements UserRepository using PostgreSQL.
type postgresUserRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresUserRepository creates a new PostgreSQL-based repository.
func NewPostgresUserRepository() UserRepository {
	return &postgresUserRepository{}
}
