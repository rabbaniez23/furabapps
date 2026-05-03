// Package repository provides data access layer for otp-service.
package repository

import (
	"context"

	"furab-backend/services/otp-service/internal/model"
)

// OTPRepository defines the interface for otp-service data access.
type OTPRepository interface {
	SendOTP(ctx context.Context) error
	VerifyOTP(ctx context.Context) error
	ResendOTP(ctx context.Context) error

	Save(ctx context.Context, otp *model.OTP) error
	FindByPhone(ctx context.Context, phone string) (*model.OTP, error)
}

// postgresOTPRepository implements OTPRepository using PostgreSQL.
type postgresOTPRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresOTPRepository creates a new PostgreSQL-based repository.
func NewPostgresOTPRepository() OTPRepository {
	return &postgresOTPRepository{}
}

func (r *postgresOTPRepository) SendOTP(ctx context.Context) error   { return nil }
func (r *postgresOTPRepository) VerifyOTP(ctx context.Context) error { return nil }
func (r *postgresOTPRepository) ResendOTP(ctx context.Context) error { return nil }
func (r *postgresOTPRepository) Save(ctx context.Context, otp *model.OTP) error { return nil }
func (r *postgresOTPRepository) FindByPhone(ctx context.Context, phone string) (*model.OTP, error) { return nil, nil }
