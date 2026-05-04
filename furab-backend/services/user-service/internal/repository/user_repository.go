// Package repository provides data access layer for user-service.
package repository

import (
	"context"

	"furab-backend/services/user-service/internal/model"
)

// UserRepository defines the interface for user-service data access.
type UserRepository interface {
	GetProfile(ctx context.Context) error
	UpdateProfile(ctx context.Context) error
	AddAddress(ctx context.Context) error
	DeleteAddress(ctx context.Context) error

	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Deactivate(ctx context.Context, userID string) error
}

// postgresUserRepository implements UserRepository using PostgreSQL.
type postgresUserRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresUserRepository creates a new PostgreSQL-based repository.
func NewPostgresUserRepository() UserRepository {
	return &postgresUserRepository{}
}

func (r *postgresUserRepository) GetProfile(ctx context.Context) error { return nil }
func (r *postgresUserRepository) UpdateProfile(ctx context.Context) error { return nil }
func (r *postgresUserRepository) AddAddress(ctx context.Context) error { return nil }
func (r *postgresUserRepository) DeleteAddress(ctx context.Context) error { return nil }
func (r *postgresUserRepository) Save(ctx context.Context, user *model.User) error { return nil }
func (r *postgresUserRepository) FindByID(ctx context.Context, userID string) (*model.User, error) { return nil, nil }
func (r *postgresUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (r *postgresUserRepository) Deactivate(ctx context.Context, userID string) error { return nil }
