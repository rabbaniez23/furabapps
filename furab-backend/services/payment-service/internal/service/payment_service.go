// Package service implements the business logic for payment-service.
package service

import "context"

// PaymentService defines the interface for payment-service business logic.
type PaymentService interface {

	// Authorize implements the business logic for Authorize.
	Authorize(ctx context.Context) error

	// Capture implements the business logic for Capture.
	Capture(ctx context.Context) error

	// Refund implements the business logic for Refund.
	Refund(ctx context.Context) error

	// GetPayment implements the business logic for GetPayment.
	GetPayment(ctx context.Context) error
}

// paymentServiceImpl is the concrete implementation of PaymentService.
type paymentServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService() PaymentService {
	return &paymentServiceImpl{}
}
