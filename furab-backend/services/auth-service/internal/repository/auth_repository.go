// Package repository provides data access layer for auth-service.
package repository

import "context"

// AuthRepository defines the interface for auth-service data access.
type AuthRepository interface {

	// Login performs the Login operation.
	Login(ctx context.Context) error

	// Register performs the Register operation.
	Register(ctx context.Context) error

	// RefreshToken performs the RefreshToken operation.
	RefreshToken(ctx context.Context) error

	// Logout performs the Logout operation.
	Logout(ctx context.Context) error
}

// postgresAuthRepository implements AuthRepository using PostgreSQL.
type postgresAuthRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresAuthRepository creates a new PostgreSQL-based repository.
func NewPostgresAuthRepository() AuthRepository {
	return &postgresAuthRepository{}
}
