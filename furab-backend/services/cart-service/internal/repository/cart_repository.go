// Package repository provides data access layer for cart.
package repository

import (
	"context"
	"errors"

	"furab-backend/services/cart-service/internal/model"
)

// Common repository errors.
var (
	ErrCartNotFound = errors.New("cart not found")
	ErrItemNotFound = errors.New("item not found in cart")
)

// CartRepository defines the interface for cart data access.
type CartRepository interface {
	// GetByUserID retrieves a cart by user ID.
	GetByUserID(ctx context.Context, userID string) (*model.Cart, error)

	// Save creates or updates a cart.
	Save(ctx context.Context, cart *model.Cart) error

	// Delete removes a cart by user ID.
	Delete(ctx context.Context, userID string) error
}
