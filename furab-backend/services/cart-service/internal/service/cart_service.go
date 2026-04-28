// Package service implements the business logic for cart-service.
package service

import "context"

// CartService defines the interface for cart-service business logic.
type CartService interface {

	// AddItem implements the business logic for AddItem.
	AddItem(ctx context.Context) error

	// RemoveItem implements the business logic for RemoveItem.
	RemoveItem(ctx context.Context) error

	// UpdateQuantity implements the business logic for UpdateQuantity.
	UpdateQuantity(ctx context.Context) error

	// GetCart implements the business logic for GetCart.
	GetCart(ctx context.Context) error

	// ClearCart implements the business logic for ClearCart.
	ClearCart(ctx context.Context) error
}

// cartServiceImpl is the concrete implementation of CartService.
type cartServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewCartService creates a new CartService.
func NewCartService() CartService {
	return &cartServiceImpl{}
}
