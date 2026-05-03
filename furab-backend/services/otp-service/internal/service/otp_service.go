// Package service implements the business logic for otp-service.
package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/otp-service/internal/model"
	"furab-backend/services/otp-service/internal/repository"
)

type OTP = model.OTP

type GenerateOTPRequest struct {
	Contact string
}

type GenerateOTPResponse struct {
	OTPCode string
	Message string
}

type VerifyOTPRequest struct {
	Contact string
	Code    string
}

type VerifyOTPResponse struct {
	Valid   bool
	Message string
}

// OTPService defines the interface for otp-service business logic.
type OTPService interface {
	GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error)
	VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error)
}

// otpServiceImpl is the concrete implementation of OTPService.
type otpServiceImpl struct {
	repo repository.OTPRepository
}

// NewOTPService creates a new OTPService.
func NewOTPService(repo repository.OTPRepository) OTPService {
	return &otpServiceImpl{repo: repo}
}

func (s *otpServiceImpl) GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	if req.Contact == "" {
		return nil, errors.New("phone/email required")
	}

	otp := &OTP{
		Phone:     req.Contact,
		Code:      "123456",
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	if err := s.repo.Save(ctx, otp); err != nil {
		return nil, err
	}

	return &GenerateOTPResponse{
		OTPCode: otp.Code,
		Message: "OTP generated",
	}, nil
}

func (s *otpServiceImpl) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	otp, err := s.repo.FindByPhone(ctx, req.Contact)
	if err != nil {
		return nil, err
	}
	if otp == nil {
		return nil, errors.New("otp not found")
	}

	if otp.Code != req.Code {
		return &VerifyOTPResponse{
			Valid:   false,
			Message: "OTP invalid",
		}, nil
	}

	if time.Now().Unix() > otp.ExpiresAt {
		return &VerifyOTPResponse{
			Valid:   false,
			Message: "OTP expired",
		}, nil
	}

	return &VerifyOTPResponse{
		Valid:   true,
		Message: "OTP valid",
	}, nil
}
