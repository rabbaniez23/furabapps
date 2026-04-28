// Package repository provides data access layer for otp-service.
package repository

import "context"

// OTPRepository defines the interface for otp-service data access.
type OTPRepository interface {

	// SendOTP performs the SendOTP operation.
	SendOTP(ctx context.Context) error

	// VerifyOTP performs the VerifyOTP operation.
	VerifyOTP(ctx context.Context) error

	// ResendOTP performs the ResendOTP operation.
	ResendOTP(ctx context.Context) error
}

// postgresOTPRepository implements OTPRepository using PostgreSQL.
type postgresOTPRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresOTPRepository creates a new PostgreSQL-based repository.
func NewPostgresOTPRepository() OTPRepository {
	return &postgresOTPRepository{}
}
