package unit

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

// ============================================================================
// Catatan: Struct dan interface di bawah ini adalah representasi dari desain.
// Pada proyek nyata, definisi ini akan berada di internal/service dan mock
// digenerate otomatis menggunakan gomock (mockgen).
// ============================================================================

// 1. DTO (Data Transfer Objects)
type AuthResponse struct {
	Status  string
	Message string
}

type LoginResponse struct {
	Status      string
	Message     string
	AccessToken string
}

type TokenResponse struct {
	Status  string
	Message string
}

// 2. Mocked Dependencies Interface

// MockUserService mensimulasikan interaksi dengan User Service
type MockUserService struct {
	SaveUserFunc func(ctx context.Context, contact string) error
}

func (m *MockUserService) SaveUser(ctx context.Context, contact string) error {
	if m.SaveUserFunc != nil {
		return m.SaveUserFunc(ctx, contact)
	}
	return nil
}

// MockOTPService mensimulasikan interaksi dengan OTP Service
type MockOTPService struct {
	GenerateOTPFunc func(ctx context.Context, contact string) error
	VerifyOTPFunc   func(ctx context.Context, contact, otpCode string) (bool, error)
}

func (m *MockOTPService) GenerateOTP(ctx context.Context, contact string) error {
	if m.GenerateOTPFunc != nil {
		return m.GenerateOTPFunc(ctx, contact)
	}
	return nil
}

func (m *MockOTPService) VerifyOTP(ctx context.Context, contact, otpCode string) (bool, error) {
	if m.VerifyOTPFunc != nil {
		return m.VerifyOTPFunc(ctx, contact, otpCode)
	}
	return false, nil
}

// MockTokenGenerator mensimulasikan pembuatan dan validasi JWT/Token
type MockTokenGenerator struct {
	GenerateTokenFunc func(userID string) (string, error)
	ValidateTokenFunc func(token string) (bool, error)
}

func (m *MockTokenGenerator) GenerateToken(userID string) (string, error) {
	if m.GenerateTokenFunc != nil {
		return m.GenerateTokenFunc(userID)
	}
	return "", nil
}

func (m *MockTokenGenerator) ValidateToken(token string) (bool, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return false, nil
}

// 3. AuthService Interface & Implementation

type AuthService interface {
	Register(ctx context.Context, contact string) (*AuthResponse, error)
	RequestOTP(ctx context.Context, contact string) (*AuthResponse, error)
	VerifyOTP(ctx context.Context, contact, otpCode string) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*TokenResponse, error)
}

type authServiceImpl struct {
	userService    *MockUserService
	otpService     *MockOTPService
	tokenGenerator *MockTokenGenerator
}

func (s *authServiceImpl) Register(ctx context.Context, contact string) (*AuthResponse, error) {
	if contact == "" {
		return &AuthResponse{Status: "failed", Message: "input required"}, errors.New("input required")
	}

	// Panggil dependency User Service
	err := s.userService.SaveUser(ctx, contact)
	if err != nil {
		return &AuthResponse{Status: "failed", Message: "gagal menyimpan user"}, err
	}

	return &AuthResponse{Status: "success", Message: "registrasi berhasil"}, nil
}

func (s *authServiceImpl) RequestOTP(ctx context.Context, contact string) (*AuthResponse, error) {
	if contact == "" {
		return &AuthResponse{Status: "failed", Message: "input required"}, errors.New("input required")
	}

	// Panggil dependency OTP Service
	err := s.otpService.GenerateOTP(ctx, contact)
	if err != nil {
		return &AuthResponse{Status: "failed", Message: "gagal request OTP"}, err
	}

	return &AuthResponse{Status: "success", Message: "OTP dikirim"}, nil
}

func (s *authServiceImpl) VerifyOTP(ctx context.Context, contact, otpCode string) (*LoginResponse, error) {
	valid, err := s.otpService.VerifyOTP(ctx, contact, otpCode)
	if err != nil || !valid {
		return &LoginResponse{Status: "failed", Message: "OTP tidak valid"}, errors.New("OTP tidak valid")
	}

	// Dummy get user ID for token generation
	userID := "user-123"

	token, err := s.tokenGenerator.GenerateToken(userID)
	if err != nil {
		return &LoginResponse{Status: "failed", Message: "gagal membuat token"}, err
	}

	return &LoginResponse{
		Status:      "success",
		Message:     "login berhasil",
		AccessToken: token,
	}, nil
}

func (s *authServiceImpl) ValidateToken(ctx context.Context, token string) (*TokenResponse, error) {
	valid, err := s.tokenGenerator.ValidateToken(token)
	if err != nil || !valid {
		return &TokenResponse{Status: "invalid", Message: "token tidak valid"}, errors.New("token tidak valid")
	}

	return &TokenResponse{Status: "valid", Message: "token valid"}, nil
}

// ============================================================================
// UNIT TESTS MULAI DARI SINI
// ============================================================================

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Registrasi berhasil", func(t *testing.T) {
		userService := &MockUserService{
			SaveUserFunc: func(ctx context.Context, contact string) error {
				if contact != "08123456789" {
					t.Errorf("Expected contact '08123456789', got %s", contact)
				}
				return nil
			},
		}
		service := &authServiceImpl{userService: userService}

		res, err := service.Register(context.Background(), "08123456789")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "registrasi berhasil" {
			t.Errorf("Expected success response, got %v", res)
		}
	})

	t.Run("Error - Input kosong", func(t *testing.T) {
		userService := &MockUserService{
			SaveUserFunc: func(ctx context.Context, contact string) error {
				t.Fatal("UserService should not be called when input is empty")
				return nil
			},
		}
		service := &authServiceImpl{userService: userService}

		res, err := service.Register(context.Background(), "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if res == nil || res.Status != "failed" || res.Message != "input required" {
			t.Errorf("Expected failed response due to empty input, got %v", res)
		}
	})
}

func TestAuthService_RequestOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - OTP berhasil dikirim", func(t *testing.T) {
		otpService := &MockOTPService{
			GenerateOTPFunc: func(ctx context.Context, contact string) error {
				if contact != "08123456789" {
					t.Errorf("Expected contact '08123456789', got %s", contact)
				}
				return nil
			},
		}
		service := &authServiceImpl{otpService: otpService}

		res, err := service.RequestOTP(context.Background(), "08123456789")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "OTP dikirim" {
			t.Errorf("Expected success response, got %v", res)
		}
	})

	t.Run("Error - Input kosong", func(t *testing.T) {
		otpService := &MockOTPService{
			GenerateOTPFunc: func(ctx context.Context, contact string) error {
				t.Fatal("OTPService should not be called when input is empty")
				return nil
			},
		}
		service := &authServiceImpl{otpService: otpService}

		res, err := service.RequestOTP(context.Background(), "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if res == nil || res.Status != "failed" || res.Message != "input required" {
			t.Errorf("Expected failed response due to empty input, got %v", res)
		}
	})
}

func TestAuthService_VerifyOTPAndLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Login berhasil", func(t *testing.T) {
		otpService := &MockOTPService{
			VerifyOTPFunc: func(ctx context.Context, contact, otpCode string) (bool, error) {
				return true, nil // OTP Valid
			},
		}
		tokenGenerator := &MockTokenGenerator{
			GenerateTokenFunc: func(userID string) (string, error) {
				return "access_token_123", nil
			},
		}

		service := &authServiceImpl{otpService: otpService, tokenGenerator: tokenGenerator}

		res, err := service.VerifyOTP(context.Background(), "08123456789", "123456")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "login berhasil" {
			t.Errorf("Expected success response, got %v", res)
		}
		if res.AccessToken != "access_token_123" {
			t.Errorf("Expected AccessToken, got %s", res.AccessToken)
		}
	})

	t.Run("Error - OTP tidak valid", func(t *testing.T) {
		otpService := &MockOTPService{
			VerifyOTPFunc: func(ctx context.Context, contact, otpCode string) (bool, error) {
				return false, errors.New("invalid code") // OTP Invalid
			},
		}
		tokenGenerator := &MockTokenGenerator{
			GenerateTokenFunc: func(userID string) (string, error) {
				t.Fatal("Token generator should not be called if OTP is invalid")
				return "", nil
			},
		}

		service := &authServiceImpl{otpService: otpService, tokenGenerator: tokenGenerator}

		res, err := service.VerifyOTP(context.Background(), "08123456789", "999999")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if res == nil || res.Status != "failed" || res.Message != "OTP tidak valid" {
			t.Errorf("Expected failed response due to invalid OTP, got %v", res)
		}
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Token valid", func(t *testing.T) {
		tokenGen := &MockTokenGenerator{
			ValidateTokenFunc: func(token string) (bool, error) {
				if token != "valid_token" {
					t.Errorf("Expected 'valid_token', got %s", token)
				}
				return true, nil
			},
		}

		service := &authServiceImpl{tokenGenerator: tokenGen}
		res, err := service.ValidateToken(context.Background(), "valid_token")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "valid" || res.Message != "token valid" {
			t.Errorf("Expected valid token response, got %v", res)
		}
	})

	t.Run("Error - Token tidak valid", func(t *testing.T) {
		tokenGen := &MockTokenGenerator{
			ValidateTokenFunc: func(token string) (bool, error) {
				return false, errors.New("token expired")
			},
		}

		service := &authServiceImpl{tokenGenerator: tokenGen}
		res, err := service.ValidateToken(context.Background(), "invalid_token")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if res == nil || res.Status != "invalid" || res.Message != "token tidak valid" {
			t.Errorf("Expected invalid token response, got %v", res)
		}
	})
}
