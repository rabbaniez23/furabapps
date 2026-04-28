// Package repository provides data access layer for food orders.
package repository

import (
	"context"
	"errors"

	"furab-backend/services/food-order-service/internal/model"
)

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrDuplicateOrder = errors.New("duplicate order")
)

// FoodOrderRepository defines the interface for food order data access.
type FoodOrderRepository interface {
	Create(ctx context.Context, order *model.FoodOrder) error
	GetByID(ctx context.Context, id string) (*model.FoodOrder, error)
	Update(ctx context.Context, order *model.FoodOrder) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*model.FoodOrder, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}
