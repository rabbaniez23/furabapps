package service

import (
	"context"
	"errors"

	"furab-backend/services/auth-service/internal/model"
)

// AuthResponse represents standard auth response
type AuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

// TokenResponse represents token validation response
type TokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UserService defines the interface for interacting with user data
type UserService interface {
	CreateUser(ctx context.Context, contact string) error
	GetUser(ctx context.Context, contact string) (*model.User, error)
}

// OTPService defines the interface for OTP operations
type OTPService interface {
	GenerateOTP(ctx context.Context, contact string) error
	VerifyOTP(ctx context.Context, contact, otpCode string) (bool, error)
}

// TokenGenerator defines the interface for token operations
type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (bool, error)
}

// AuthService defines the interface for auth-service business logic.
type AuthService interface {
	Register(ctx context.Context, contact string) (*AuthResponse, error)
	RequestOTP(ctx context.Context, contact string) (*AuthResponse, error)
	VerifyOTP(ctx context.Context, contact, otpCode string) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*TokenResponse, error)
}

// authServiceImpl is the concrete implementation of AuthService.
type authServiceImpl struct {
	userService    UserService
	otpService     OTPService
	tokenGenerator TokenGenerator
}

// NewAuthService creates a new AuthService.
func NewAuthService(userService UserService, otpService OTPService, tokenGenerator TokenGenerator) AuthService {
	return &authServiceImpl{
		userService:    userService,
		otpService:     otpService,
		tokenGenerator: tokenGenerator,
	}
}

func (s *authServiceImpl) Register(ctx context.Context, contact string) (*AuthResponse, error) {
	if contact == "" {
		return nil, errors.New("phone/email required")
	}

	err := s.userService.CreateUser(ctx, contact)
	if err != nil {
		return nil, err
	}

	err = s.otpService.GenerateOTP(ctx, contact)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Status: "success", Message: "register success"}, nil
}

func (s *authServiceImpl) RequestOTP(ctx context.Context, contact string) (*AuthResponse, error) {
	if contact == "" {
		return nil, errors.New("phone/email required")
	}

	err := s.otpService.GenerateOTP(ctx, contact)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Status: "success", Message: "OTP dikirim"}, nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, contact, otpCode string) (*LoginResponse, error) {
	if contact == "" || otpCode == "" {
		return nil, errors.New("input required")
	}

	valid, err := s.otpService.VerifyOTP(ctx, contact, otpCode)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("OTP tidak valid")
	}

	user, err := s.userService.GetUser(ctx, contact)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	token, err := s.tokenGenerator.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Status:      "success",
		Message:     "login berhasil",
		AccessToken: token,
	}, nil
}

func (s *authServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenResponse, error) {
	if token == "" {
		return &TokenResponse{Status: "invalid", Message: "token invalid"}, nil
	}
	
	valid, err := s.tokenGenerator.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	if !valid {
		return &TokenResponse{Status: "invalid", Message: "token invalid"}, nil
	}

	return &TokenResponse{Status: "valid", Message: "token valid"}, nil
}
