package unit

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"furab-backend/services/auth-service/internal/model"
	"furab-backend/services/auth-service/internal/service"
	mock_service "furab-backend/services/auth-service/test/unit/mock"
)

// ============================================================================
// TestAuthService_Register
// ============================================================================

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mock_service.NewMockUserService(ctrl)
	mockOTP := mock_service.NewMockOTPService(ctrl)
	mockToken := mock_service.NewMockTokenGenerator(ctrl)

	svc := service.NewAuthService(mockUser, mockOTP, mockToken)

	ctx := context.Background()

	t.Run("Success - Register berhasil", func(t *testing.T) {
		contact := "08123"

		mockUser.EXPECT().CreateUser(gomock.Any(), contact).Return(nil)
		mockOTP.EXPECT().GenerateOTP(gomock.Any(), contact).Return(nil)

		res, err := svc.Register(ctx, contact)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.Message != "register success" {
			t.Errorf("Expected message 'register success', got %v", res.Message)
		}
	})

	t.Run("Error - Input kosong", func(t *testing.T) {
		res, err := svc.Register(ctx, "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "phone/email required" {
			t.Errorf("Expected error 'phone/email required', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - User service gagal", func(t *testing.T) {
		contact := "08123"
		expectedErr := errors.New("user service error")

		mockUser.EXPECT().CreateUser(gomock.Any(), contact).Return(expectedErr)

		res, err := svc.Register(ctx, contact)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestAuthService_RequestOTP
// ============================================================================

func TestAuthService_RequestOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mock_service.NewMockUserService(ctrl)
	mockOTP := mock_service.NewMockOTPService(ctrl)
	mockToken := mock_service.NewMockTokenGenerator(ctrl)

	svc := service.NewAuthService(mockUser, mockOTP, mockToken)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		contact := "08123"

		mockOTP.EXPECT().GenerateOTP(gomock.Any(), contact).Return(nil)

		res, err := svc.RequestOTP(ctx, contact)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
	})

	t.Run("Error - Input kosong", func(t *testing.T) {
		res, err := svc.RequestOTP(ctx, "")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "phone/email required" {
			t.Errorf("Expected error 'phone/email required', got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - OTP gagal", func(t *testing.T) {
		contact := "08123"
		expectedErr := errors.New("otp error")

		mockOTP.EXPECT().GenerateOTP(gomock.Any(), contact).Return(expectedErr)

		res, err := svc.RequestOTP(ctx, contact)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestAuthService_VerifyOTP
// ============================================================================

func TestAuthService_VerifyOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mock_service.NewMockUserService(ctrl)
	mockOTP := mock_service.NewMockOTPService(ctrl)
	mockToken := mock_service.NewMockTokenGenerator(ctrl)

	svc := service.NewAuthService(mockUser, mockOTP, mockToken)

	ctx := context.Background()

	t.Run("Success - OTP valid", func(t *testing.T) {
		contact := "08123"
		otpCode := "123456"
		expectedUser := &model.User{ID: "user-123"}
		expectedToken := "access_token_123"

		mockOTP.EXPECT().VerifyOTP(gomock.Any(), contact, otpCode).Return(true, nil)
		mockUser.EXPECT().GetUser(gomock.Any(), contact).Return(expectedUser, nil)
		mockToken.EXPECT().GenerateToken(expectedUser.ID).Return(expectedToken, nil)

		res, err := svc.VerifyOTP(ctx, contact, otpCode)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.AccessToken != expectedToken {
			t.Errorf("Expected token %v, got %v", expectedToken, res.AccessToken)
		}
	})

	t.Run("Error - OTP invalid", func(t *testing.T) {
		contact := "08123"
		otpCode := "wrong_otp"

		mockOTP.EXPECT().VerifyOTP(gomock.Any(), contact, otpCode).Return(false, nil)

		res, err := svc.VerifyOTP(ctx, contact, otpCode)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "OTP tidak valid" {
			t.Errorf("Expected error 'OTP tidak valid', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - OTP service error", func(t *testing.T) {
		contact := "08123"
		otpCode := "123456"
		expectedErr := errors.New("otp service error")

		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), contact, otpCode).
			Return(false, expectedErr)

		res, err := svc.VerifyOTP(ctx, contact, otpCode)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - User tidak ditemukan", func(t *testing.T) {
		contact := "08123"
		otpCode := "123456"

		mockOTP.EXPECT().VerifyOTP(gomock.Any(), contact, otpCode).Return(true, nil)
		mockUser.EXPECT().GetUser(gomock.Any(), contact).Return(nil, nil)

		res, err := svc.VerifyOTP(ctx, contact, otpCode)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "user not found" {
			t.Errorf("Expected error 'user not found', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Token gagal dibuat", func(t *testing.T) {
		contact := "08123"
		otpCode := "123456"
		expectedUser := &model.User{ID: "user-123"}
		expectedErr := errors.New("token generation error")

		mockOTP.EXPECT().VerifyOTP(gomock.Any(), contact, otpCode).Return(true, nil)
		mockUser.EXPECT().GetUser(gomock.Any(), contact).Return(expectedUser, nil)
		mockToken.EXPECT().GenerateToken(expectedUser.ID).Return("", expectedErr)

		res, err := svc.VerifyOTP(ctx, contact, otpCode)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestAuthService_ValidateToken
// ============================================================================

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := mock_service.NewMockUserService(ctrl)
	mockOTP := mock_service.NewMockOTPService(ctrl)
	mockToken := mock_service.NewMockTokenGenerator(ctrl)

	svc := service.NewAuthService(mockUser, mockOTP, mockToken)

	ctx := context.Background()

	t.Run("Success - Token valid", func(t *testing.T) {
		token := "valid_token"

		mockToken.EXPECT().ValidateToken(token).Return(true, nil)

		res, err := svc.ValidateToken(ctx, token)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "valid" {
			t.Errorf("Expected status valid, got %v", res.Status)
		}
	})

	t.Run("Error - Token invalid", func(t *testing.T) {
		token := "invalid_token"

		mockToken.EXPECT().ValidateToken(token).Return(false, nil)

		res, err := svc.ValidateToken(ctx, token)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "invalid" {
			t.Errorf("Expected status invalid, got %v", res.Status)
		}
	})

	t.Run("Error - Token service error", func(t *testing.T) {
		token := "token"
		expectedErr := errors.New("token error")

		mockToken.EXPECT().
			ValidateToken(token).
			Return(false, expectedErr)

		res, err := svc.ValidateToken(ctx, token)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}
