// Package repository provides data access layer for food-order-service.
package repository

import "context"

// FoodOrderRepository defines the interface for food-order-service data access.
type FoodOrderRepository interface {

	// CreateOrder performs the CreateOrder operation.
	CreateOrder(ctx context.Context) error

	// ConfirmOrder performs the ConfirmOrder operation.
	ConfirmOrder(ctx context.Context) error

	// PrepareOrder performs the PrepareOrder operation.
	PrepareOrder(ctx context.Context) error

	// CompleteOrder performs the CompleteOrder operation.
	CompleteOrder(ctx context.Context) error

	// CancelOrder performs the CancelOrder operation.
	CancelOrder(ctx context.Context) error
}

// postgresFoodOrderRepository implements FoodOrderRepository using PostgreSQL.
type postgresFoodOrderRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresFoodOrderRepository creates a new PostgreSQL-based repository.
func NewPostgresFoodOrderRepository() FoodOrderRepository {
	return &postgresFoodOrderRepository{}
}
