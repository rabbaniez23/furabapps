// Package service implements the business logic for menu-service.
package service

import "context"

// MenuService defines the interface for menu-service business logic.
type MenuService interface {

	// GetMenu implements the business logic for GetMenu.
	GetMenu(ctx context.Context) error

	// AddItem implements the business logic for AddItem.
	AddItem(ctx context.Context) error

	// UpdateItem implements the business logic for UpdateItem.
	UpdateItem(ctx context.Context) error

	// DeleteItem implements the business logic for DeleteItem.
	DeleteItem(ctx context.Context) error

	// GetCategories implements the business logic for GetCategories.
	GetCategories(ctx context.Context) error
}

// menuServiceImpl is the concrete implementation of MenuService.
type menuServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewMenuService creates a new MenuService.
func NewMenuService() MenuService {
	return &menuServiceImpl{}
}
