// Package service implements the business logic for otp-service.
package service

import "context"

// OTPService defines the interface for otp-service business logic.
type OTPService interface {

	// SendOTP implements the business logic for SendOTP.
	SendOTP(ctx context.Context) error

	// VerifyOTP implements the business logic for VerifyOTP.
	VerifyOTP(ctx context.Context) error

	// ResendOTP implements the business logic for ResendOTP.
	ResendOTP(ctx context.Context) error
}

// otpServiceImpl is the concrete implementation of OTPService.
type otpServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewOTPService creates a new OTPService.
func NewOTPService() OTPService {
	return &otpServiceImpl{}
}
