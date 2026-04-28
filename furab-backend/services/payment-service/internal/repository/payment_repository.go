// Package repository provides data access layer for payment-service.
package repository

import "context"

// PaymentRepository defines the interface for payment-service data access.
type PaymentRepository interface {

	// Authorize performs the Authorize operation.
	Authorize(ctx context.Context) error

	// Capture performs the Capture operation.
	Capture(ctx context.Context) error

	// Refund performs the Refund operation.
	Refund(ctx context.Context) error

	// GetPayment performs the GetPayment operation.
	GetPayment(ctx context.Context) error
}

// postgresPaymentRepository implements PaymentRepository using PostgreSQL.
type postgresPaymentRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresPaymentRepository creates a new PostgreSQL-based repository.
func NewPostgresPaymentRepository() PaymentRepository {
	return &postgresPaymentRepository{}
}
