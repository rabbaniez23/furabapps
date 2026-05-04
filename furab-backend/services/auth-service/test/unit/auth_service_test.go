// Package unit contains unit tests for the auth service.
// All dependencies are mocked using gomock. No database access.
package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/auth-service/internal/model"
	"furab-backend/services/auth-service/internal/service"
	mock_service "furab-backend/services/auth-service/internal/service/mock"

	"go.uber.org/mock/gomock"
)

// --- Helper Functions ---

func newTestService(t *testing.T) (
	service.AuthService,
	*mock_service.MockUserService,
	*mock_service.MockOTPService,
	*mock_service.MockTokenGenerator,
	*gomock.Controller,
) {
	ctrl := gomock.NewController(t)
	mockUser := mock_service.NewMockUserService(ctrl)
	mockOTP := mock_service.NewMockOTPService(ctrl)
	mockToken := mock_service.NewMockTokenGenerator(ctrl)
	svc := service.NewAuthService(mockUser, mockOTP, mockToken)
	return svc, mockUser, mockOTP, mockToken, ctrl
}

// ========================================
// Register
// ========================================

func TestRegister_Success(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockUser.EXPECT().CreateUser(gomock.Any(), "081234567890").Return(nil)
	mockOTP.EXPECT().GenerateOTP(gomock.Any(), "081234567890").Return(nil)

	resp, err := svc.Register(context.Background(), "081234567890")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestRegister_EmptyContact(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.Register(context.Background(), "")
	if err != service.ErrContactRequired {
		t.Fatalf("expected ErrContactRequired, got: %v", err)
	}
}

func TestRegister_InvalidFormat(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.Register(context.Background(), "abc")
	if err != service.ErrContactInvalidFormat {
		t.Fatalf("expected ErrContactInvalidFormat, got: %v", err)
	}
}

func TestRegister_CreateUserError(t *testing.T) {
	svc, mockUser, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	createErr := errors.New("user creation failed")
	mockUser.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(createErr)

	_, err := svc.Register(context.Background(), "081234567890")
	if !errors.Is(err, createErr) {
		t.Fatalf("expected create error, got: %v", err)
	}
}

func TestRegister_GenerateOTPError(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	otpErr := errors.New("OTP generation failed")
	mockUser.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
	mockOTP.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Return(otpErr)

	_, err := svc.Register(context.Background(), "081234567890")
	if !errors.Is(err, otpErr) {
		t.Fatalf("expected OTP error, got: %v", err)
	}
}

func TestRegister_EmailContact(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockUser.EXPECT().CreateUser(gomock.Any(), "user@example.com").Return(nil)
	mockOTP.EXPECT().GenerateOTP(gomock.Any(), "user@example.com").Return(nil)

	resp, err := svc.Register(context.Background(), "  user@example.com  ")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

// ========================================
// RequestOTP
// ========================================

func TestRequestOTP_Success(t *testing.T) {
	svc, _, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockOTP.EXPECT().GenerateOTP(gomock.Any(), "081234567890").Return(nil)

	resp, err := svc.RequestOTP(context.Background(), "081234567890")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestRequestOTP_EmptyContact(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.RequestOTP(context.Background(), "")
	if err != service.ErrContactRequired {
		t.Fatalf("expected ErrContactRequired, got: %v", err)
	}
}

func TestRequestOTP_OTPError(t *testing.T) {
	svc, _, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	otpErr := errors.New("OTP service down")
	mockOTP.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Return(otpErr)

	_, err := svc.RequestOTP(context.Background(), "081234567890")
	if !errors.Is(err, otpErr) {
		t.Fatalf("expected OTP error, got: %v", err)
	}
}

// ========================================
// VerifyOTP (Login)
// ========================================

func TestVerifyOTP_Success(t *testing.T) {
	svc, mockUser, mockOTP, mockToken, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockOTP.EXPECT().VerifyOTP(gomock.Any(), "081234567890", "123456").Return(true, nil)
	mockUser.EXPECT().GetUser(gomock.Any(), "081234567890").Return(&model.User{
		ID: "user-abc-123", Contact: "081234567890",
	}, nil)
	mockToken.EXPECT().GenerateToken("user-abc-123").Return("jwt-token-xyz", nil)

	resp, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
	if resp.AccessToken != "jwt-token-xyz" {
		t.Errorf("expected jwt-token-xyz, got: %s", resp.AccessToken)
	}
}

func TestVerifyOTP_InvalidOTP(t *testing.T) {
	svc, _, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockOTP.EXPECT().VerifyOTP(gomock.Any(), "081234567890", "000000").Return(false, nil)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "000000")
	if err != service.ErrOTPInvalid {
		t.Fatalf("expected ErrOTPInvalid, got: %v", err)
	}
}

func TestVerifyOTP_OTPServiceError(t *testing.T) {
	svc, _, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	otpErr := errors.New("OTP verification failed")
	mockOTP.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, otpErr)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if !errors.Is(err, otpErr) {
		t.Fatalf("expected OTP error, got: %v", err)
	}
}

func TestVerifyOTP_UserNotFound(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockOTP.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	mockUser.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, nil)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if err != service.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestVerifyOTP_UserIDMissing(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockOTP.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	mockUser.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&model.User{
		ID: "", Contact: "081234567890",
	}, nil)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if err != service.ErrUserIDMissing {
		t.Fatalf("expected ErrUserIDMissing, got: %v", err)
	}
}

func TestVerifyOTP_TokenGenerationError(t *testing.T) {
	svc, mockUser, mockOTP, mockToken, ctrl := newTestService(t)
	defer ctrl.Finish()

	tokenErr := errors.New("token generation failed")
	mockOTP.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	mockUser.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(&model.User{
		ID: "user-123", Contact: "081234567890",
	}, nil)
	mockToken.EXPECT().GenerateToken("user-123").Return("", tokenErr)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if !errors.Is(err, tokenErr) {
		t.Fatalf("expected token error, got: %v", err)
	}
}

func TestVerifyOTP_EmptyContact(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.VerifyOTP(context.Background(), "", "123456")
	if err != service.ErrInputRequired {
		t.Fatalf("expected ErrInputRequired, got: %v", err)
	}
}

func TestVerifyOTP_EmptyOTPCode(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "")
	if err != service.ErrInputRequired {
		t.Fatalf("expected ErrInputRequired, got: %v", err)
	}
}

func TestVerifyOTP_GetUserError(t *testing.T) {
	svc, mockUser, mockOTP, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	userErr := errors.New("user service error")
	mockOTP.EXPECT().VerifyOTP(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
	mockUser.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, userErr)

	_, err := svc.VerifyOTP(context.Background(), "081234567890", "123456")
	if !errors.Is(err, userErr) {
		t.Fatalf("expected user error, got: %v", err)
	}
}

// ========================================
// ValidateToken
// ========================================

func TestValidateToken_Valid(t *testing.T) {
	svc, _, _, mockToken, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockToken.EXPECT().ValidateToken("valid-token").Return(true, nil)

	resp, err := svc.ValidateToken(context.Background(), "valid-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "valid" {
		t.Errorf("expected valid, got: %s", resp.Status)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc, _, _, mockToken, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockToken.EXPECT().ValidateToken("bad-token").Return(false, nil)

	resp, err := svc.ValidateToken(context.Background(), "bad-token")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "invalid" {
		t.Errorf("expected invalid, got: %s", resp.Status)
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	svc, _, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	resp, err := svc.ValidateToken(context.Background(), "")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "invalid" {
		t.Errorf("expected invalid, got: %s", resp.Status)
	}
}

func TestValidateToken_TokenGeneratorError(t *testing.T) {
	svc, _, _, mockToken, ctrl := newTestService(t)
	defer ctrl.Finish()

	tokenErr := errors.New("token validation error")
	mockToken.EXPECT().ValidateToken("some-token").Return(false, tokenErr)

	_, err := svc.ValidateToken(context.Background(), "some-token")
	if !errors.Is(err, tokenErr) {
		t.Fatalf("expected token error, got: %v", err)
	}
}
