// Package repository provides data access layer for match requests.
package repository

import (
	"context"
	"errors"

	"furab-backend/services/matching-service/internal/model"
)

var (
	ErrMatchNotFound  = errors.New("match request not found")
)

// MatchRepository defines the interface for match request data access.
type MatchRepository interface {
	Create(ctx context.Context, match *model.MatchRequest) error
	GetByID(ctx context.Context, id string) (*model.MatchRequest, error)
	Update(ctx context.Context, match *model.MatchRequest) error
	GetByOrderID(ctx context.Context, orderID string) (*model.MatchRequest, error)
}
