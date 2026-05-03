// Package repository provides data access layer for payment-service.
package repository

import (
	"context"
	"database/sql"

	"furab-backend/services/payment-service/internal/model"
)

// PaymentRepository defines the interface for payment-service data access.
type PaymentRepository interface {
	CreatePayment(ctx context.Context, p *model.Payment) error
	GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error)
	GetPaymentByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status model.PaymentStatus) error
	CreatePaymentLog(ctx context.Context, paymentID string, status model.PaymentStatus) error
}

// postgresPaymentRepository implements PaymentRepository using PostgreSQL.
type postgresPaymentRepository struct {
	db *sql.DB
}

// NewPostgresPaymentRepository creates a new PostgreSQL-based repository.
func NewPostgresPaymentRepository(db *sql.DB) PaymentRepository {
	return &postgresPaymentRepository{db: db}
}

func (r *postgresPaymentRepository) CreatePayment(ctx context.Context, p *model.Payment) error {
	// TODO: implement DB insert
	return nil
}

func (r *postgresPaymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	// TODO: implement DB read
	return nil, sql.ErrNoRows
}

func (r *postgresPaymentRepository) GetPaymentByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error) {
	// TODO: implement DB read
	return nil, nil
}

func (r *postgresPaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	// TODO: implement DB update
	return nil
}

func (r *postgresPaymentRepository) CreatePaymentLog(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	// TODO: implement DB insert payment_logs
	return nil
}
