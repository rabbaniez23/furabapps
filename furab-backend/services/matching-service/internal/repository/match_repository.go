// Package repository provides data access layer for matching-service.
package repository

import "context"

// MatchRepository defines the interface for matching-service data access.
type MatchRepository interface {

	// FindDriver performs the FindDriver operation.
	FindDriver(ctx context.Context) error

	// AcceptMatch performs the AcceptMatch operation.
	AcceptMatch(ctx context.Context) error

	// RejectMatch performs the RejectMatch operation.
	RejectMatch(ctx context.Context) error

	// GetMatchStatus performs the GetMatchStatus operation.
	GetMatchStatus(ctx context.Context) error
}

// postgresMatchRepository implements MatchRepository using PostgreSQL.
type postgresMatchRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresMatchRepository creates a new PostgreSQL-based repository.
func NewPostgresMatchRepository() MatchRepository {
	return &postgresMatchRepository{}
}
