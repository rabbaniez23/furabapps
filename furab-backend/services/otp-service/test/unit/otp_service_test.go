// Package unit contains unit tests for the OTP service.
// Unit tests do NOT access any database or external service.
// All dependencies are mocked using gomock.
//
// Tests cover both OTP operations:
// GenerateOTP and VerifyOTP
// Each operation is tested for success, error, and edge cases.
package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"furab-backend/services/otp-service/internal/model"
	"furab-backend/services/otp-service/internal/repository"
	mock_repository "furab-backend/services/otp-service/internal/repository/mock"
	"furab-backend/services/otp-service/internal/service"

	"go.uber.org/mock/gomock"
)

// --- Helper Functions ---

// newTestService creates a new OTPService with mocked dependencies.
func newTestService(t *testing.T) (service.OTPService, *mock_repository.MockOTPRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockOTPRepository(ctrl)
	svc := service.NewOTPService(mockRepo)
	return svc, mockRepo, ctrl
}

// sampleOTP returns a sample valid (non-expired) OTP for testing.
func sampleOTP() *model.OTP {
	return &model.OTP{
		OTPID:     "otp-abc-123",
		Target:    "081234567890",
		OTPCode:   "123456",
		ExpiredAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now().UTC(),
	}
}

// expiredOTP returns a sample expired OTP for testing.
func expiredOTP() *model.OTP {
	return &model.OTP{
		OTPID:     "otp-expired-123",
		Target:    "081234567890",
		OTPCode:   "654321",
		ExpiredAt: time.Now().Add(-5 * time.Minute), // expired 5 minutes ago
		CreatedAt: time.Now().Add(-10 * time.Minute),
	}
}

// ========================================
// Test Cases: GenerateOTP
// ========================================

// TestGenerateOTP_Success tests generating an OTP with valid input.
// Expected: OTP saved to repository, response with status "success".
func TestGenerateOTP_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Expect repository Save to be called with a valid OTP
	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, otp *model.OTP) error {
			if otp.OTPID == "" {
				t.Error("expected non-empty otp_id")
			}
			if otp.Target != "081234567890" {
				t.Errorf("expected target 081234567890, got: %s", otp.Target)
			}
			if len(otp.OTPCode) != 6 {
				t.Errorf("expected 6-digit OTP code, got: %s", otp.OTPCode)
			}
			if otp.ExpiredAt.Before(time.Now()) {
				t.Error("expected future expiration time")
			}
			if otp.CreatedAt.IsZero() {
				t.Error("expected non-zero created_at")
			}
			return nil
		})

	resp, err := svc.GenerateOTP(ctx, &model.GenerateOTPRequest{
		Target: "081234567890",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Status != "success" {
		t.Errorf("expected status 'success', got: %s", resp.Status)
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

// TestGenerateOTP_NilRequest tests generating an OTP with nil request.
func TestGenerateOTP_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GenerateOTP(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestGenerateOTP_EmptyTarget tests generating an OTP with empty target.
func TestGenerateOTP_EmptyTarget(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GenerateOTP(context.Background(), &model.GenerateOTPRequest{
		Target: "",
	})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

// TestGenerateOTP_WhitespaceOnlyTarget tests generating an OTP with whitespace-only target.
func TestGenerateOTP_WhitespaceOnlyTarget(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GenerateOTP(context.Background(), &model.GenerateOTPRequest{
		Target: "   ",
	})
	if err == nil {
		t.Fatal("expected error for whitespace-only target")
	}
}

// TestGenerateOTP_WhitespaceNormalization tests that target is trimmed.
func TestGenerateOTP_WhitespaceNormalization(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, otp *model.OTP) error {
			if otp.Target != "081234567890" {
				t.Errorf("expected trimmed target, got: %q", otp.Target)
			}
			return nil
		})

	resp, err := svc.GenerateOTP(ctx, &model.GenerateOTPRequest{
		Target: "  081234567890  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

// TestGenerateOTP_RepositoryError tests generating an OTP when repository Save fails.
func TestGenerateOTP_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database connection failed")

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(repoErr)

	_, err := svc.GenerateOTP(ctx, &model.GenerateOTPRequest{
		Target: "081234567890",
	})

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// TestGenerateOTP_EmailTarget tests generating an OTP with email target.
func TestGenerateOTP_EmailTarget(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, otp *model.OTP) error {
			if otp.Target != "john@example.com" {
				t.Errorf("expected email target, got: %s", otp.Target)
			}
			return nil
		})

	resp, err := svc.GenerateOTP(ctx, &model.GenerateOTPRequest{
		Target: "john@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

// ========================================
// Test Cases: VerifyOTP
// ========================================

// TestVerifyOTP_Success tests verifying a valid OTP (correct code, not expired).
// Expected: status "valid", OTP deleted after verification.
func TestVerifyOTP_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := sampleOTP()

	mockRepo.EXPECT().
		GetByTarget(ctx, otp.Target).
		Return(otp, nil)

	// Expect Delete to be called (one-time use)
	mockRepo.EXPECT().
		Delete(ctx, otp.OTPID).
		Return(nil)

	resp, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  otp.Target,
		OTPCode: otp.OTPCode,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Status != "valid" {
		t.Errorf("expected status 'valid', got: %s", resp.Status)
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

// TestVerifyOTP_InvalidCode tests verifying with wrong OTP code.
func TestVerifyOTP_InvalidCode(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := sampleOTP()

	mockRepo.EXPECT().
		GetByTarget(ctx, otp.Target).
		Return(otp, nil)

	_, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  otp.Target,
		OTPCode: "000000", // wrong code
	})

	if err != service.ErrOTPInvalid {
		t.Fatalf("expected ErrOTPInvalid, got: %v", err)
	}
}

// TestVerifyOTP_Expired tests verifying an expired OTP.
func TestVerifyOTP_Expired(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := expiredOTP()

	mockRepo.EXPECT().
		GetByTarget(ctx, otp.Target).
		Return(otp, nil)

	_, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  otp.Target,
		OTPCode: otp.OTPCode, // correct code but expired
	})

	if err != service.ErrOTPExpired {
		t.Fatalf("expected ErrOTPExpired, got: %v", err)
	}
}

// TestVerifyOTP_NotFound tests verifying when no OTP exists for the target.
func TestVerifyOTP_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByTarget(ctx, "unknown-target").
		Return(nil, repository.ErrOTPNotFound)

	_, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  "unknown-target",
		OTPCode: "123456",
	})

	if err != service.ErrOTPNotFound {
		t.Fatalf("expected ErrOTPNotFound, got: %v", err)
	}
}

// TestVerifyOTP_NilRequest tests verifying with nil request.
func TestVerifyOTP_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.VerifyOTP(context.Background(), nil)
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestVerifyOTP_EmptyTarget tests verifying with empty target.
func TestVerifyOTP_EmptyTarget(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.VerifyOTP(context.Background(), &model.VerifyOTPRequest{
		Target:  "",
		OTPCode: "123456",
	})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

// TestVerifyOTP_EmptyCode tests verifying with empty OTP code.
func TestVerifyOTP_EmptyCode(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.VerifyOTP(context.Background(), &model.VerifyOTPRequest{
		Target:  "081234567890",
		OTPCode: "",
	})
	if err == nil {
		t.Fatal("expected error for empty OTP code")
	}
}

// TestVerifyOTP_WhitespaceNormalization tests that inputs are trimmed before verification.
func TestVerifyOTP_WhitespaceNormalization(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	otp := sampleOTP()

	// Mock expects the trimmed target
	mockRepo.EXPECT().
		GetByTarget(ctx, "081234567890").
		Return(otp, nil)

	mockRepo.EXPECT().
		Delete(ctx, otp.OTPID).
		Return(nil)

	// Pass target and code with whitespace
	resp, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  "  081234567890  ",
		OTPCode: "  123456  ",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "valid" {
		t.Errorf("expected valid, got: %s", resp.Status)
	}
}

// TestVerifyOTP_RepositoryError tests verifying when GetByTarget fails.
func TestVerifyOTP_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database timeout")

	mockRepo.EXPECT().
		GetByTarget(ctx, "081234567890").
		Return(nil, repoErr)

	_, err := svc.VerifyOTP(ctx, &model.VerifyOTPRequest{
		Target:  "081234567890",
		OTPCode: "123456",
	})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// ========================================
// Test Cases: Model Validation
// ========================================

// TestGenerateOTPRequest_Validate tests the GenerateOTPRequest validation.
func TestGenerateOTPRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.GenerateOTPRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     model.GenerateOTPRequest{Target: "081234567890"},
			wantErr: false,
		},
		{
			name:    "empty target",
			req:     model.GenerateOTPRequest{Target: ""},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestVerifyOTPRequest_Validate tests the VerifyOTPRequest validation.
func TestVerifyOTPRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.VerifyOTPRequest
		wantErr bool
	}{
		{
			name:    "valid request",
			req:     model.VerifyOTPRequest{Target: "08123", OTPCode: "123456"},
			wantErr: false,
		},
		{
			name:    "empty target",
			req:     model.VerifyOTPRequest{Target: "", OTPCode: "123456"},
			wantErr: true,
		},
		{
			name:    "empty otp_code",
			req:     model.VerifyOTPRequest{Target: "08123", OTPCode: ""},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestOTP_IsExpired tests the OTP expiration check.
func TestOTP_IsExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		otp := sampleOTP()
		if otp.IsExpired() {
			t.Error("expected OTP to NOT be expired")
		}
	})

	t.Run("expired", func(t *testing.T) {
		otp := expiredOTP()
		if !otp.IsExpired() {
			t.Error("expected OTP to be expired")
		}
	})
}