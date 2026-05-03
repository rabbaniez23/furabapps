package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"furab-backend/services/otp-service/internal/model"
	"furab-backend/services/otp-service/internal/repository"
)

type OTP = model.OTP

var (
	ErrContactRequired = errors.New("contact required")
	ErrOTPRequired     = errors.New("otp required")
	ErrOTPInvalid      = errors.New("otp invalid")
	ErrOTPExpired      = errors.New("otp expired")
	ErrOTPNotFound     = errors.New("otp not found")
)

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

type OTPService interface {
	GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error)
	VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error)
}

type otpServiceImpl struct {
	repo repository.OTPRepository
}

func NewOTPService(repo repository.OTPRepository) OTPService {
	return &otpServiceImpl{repo: repo}
}

func normalizeInput(input string) string {
	return strings.TrimSpace(input)
}

func validateContact(contact string) error {
	if contact == "" {
		return ErrContactRequired
	}
	return nil
}

func validateOTP(code string) error {
	if code == "" {
		return ErrOTPRequired
	}
	return nil
}

func (s *otpServiceImpl) GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req.Contact = normalizeInput(req.Contact)
	if err := validateContact(req.Contact); err != nil {
		return nil, err
	}

	otp := &OTP{
		Phone:     req.Contact,
		Code:      "123456", // TODO: generate random OTP
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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req.Contact = normalizeInput(req.Contact)
	req.Code = normalizeInput(req.Code)
	if err := validateContact(req.Contact); err != nil {
		return nil, err
	}
	if err := validateOTP(req.Code); err != nil {
		return nil, err
	}

	otp, err := s.repo.FindByPhone(ctx, req.Contact)
	if err != nil {
		return nil, err
	}
	if otp == nil {
		return nil, ErrOTPNotFound
	}

	if otp.Code != req.Code {
		return nil, ErrOTPInvalid
	}

	if time.Now().Unix() > otp.ExpiresAt {
		return nil, ErrOTPExpired
	}

	return &VerifyOTPResponse{
		Valid:   true,
		Message: "OTP valid",
	}, nil
}