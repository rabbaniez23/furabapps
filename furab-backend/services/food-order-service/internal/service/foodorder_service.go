// Package service implements the business logic for food-order-service.
package service

import "context"

// FoodOrderService defines the interface for food-order-service business logic.
type FoodOrderService interface {

	// CreateOrder implements the business logic for CreateOrder.
	CreateOrder(ctx context.Context) error

	// ConfirmOrder implements the business logic for ConfirmOrder.
	ConfirmOrder(ctx context.Context) error

	// PrepareOrder implements the business logic for PrepareOrder.
	PrepareOrder(ctx context.Context) error

	// CompleteOrder implements the business logic for CompleteOrder.
	CompleteOrder(ctx context.Context) error

	// CancelOrder implements the business logic for CancelOrder.
	CancelOrder(ctx context.Context) error
}

// foodorderServiceImpl is the concrete implementation of FoodOrderService.
type foodorderServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewFoodOrderService creates a new FoodOrderService.
func NewFoodOrderService() FoodOrderService {
	return &foodorderServiceImpl{}
}
