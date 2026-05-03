package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"furab-backend/services/otp-service/internal/model"
	mock_repository "furab-backend/services/otp-service/internal/repository/mock"
	"furab-backend/services/otp-service/internal/service"
)

// 1. Generate OTP

func TestGenerateOTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := service.NewOTPService(mockRepo)

	res, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
		Contact: "08123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Message != "OTP generated" {
		t.Errorf("expected message 'OTP generated', got %s", res.Message)
	}
}

func TestGenerateOTP_InputKosong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	svc := service.NewOTPService(mockRepo)

	_, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
		Contact: "",
	})

	if err == nil || err.Error() != "phone/email required" {
		t.Fatalf("expected error 'phone/email required', got %v", err)
	}
}

func TestGenerateOTP_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	svc := service.NewOTPService(mockRepo)

	_, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
		Contact: "08123",
	})

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected 'db error', got %v", err)
	}
}

// 2. Verify OTP

func TestVerifyOTP_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		FindByPhone(gomock.Any(), "08123").
		Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}, nil)

	svc := service.NewOTPService(mockRepo)

	res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
		Contact: "08123",
		Code:    "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Message != "OTP valid" {
		t.Errorf("expected message 'OTP valid', got %s", res.Message)
	}
	if !res.Valid {
		t.Errorf("expected Valid to be true")
	}
}

func TestVerifyOTP_InvalidOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		FindByPhone(gomock.Any(), "08123").
		Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}, nil)

	svc := service.NewOTPService(mockRepo)

	res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
		Contact: "08123",
		Code:    "000000",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Message != "OTP invalid" {
		t.Errorf("expected message 'OTP invalid', got %s", res.Message)
	}
	if res.Valid {
		t.Errorf("expected Valid to be false")
	}
}

func TestVerifyOTP_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		FindByPhone(gomock.Any(), "08123").
		Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(-5 * time.Minute).Unix(),
		}, nil)

	svc := service.NewOTPService(mockRepo)

	res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
		Contact: "08123",
		Code:    "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Message != "OTP expired" {
		t.Errorf("expected message 'OTP expired', got %s", res.Message)
	}
	if res.Valid {
		t.Errorf("expected Valid to be false")
	}
}

func TestVerifyOTP_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockOTPRepository(ctrl)

	mockRepo.EXPECT().
		FindByPhone(gomock.Any(), "99999").
		Return(nil, nil)

	svc := service.NewOTPService(mockRepo)

	_, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
		Contact: "99999",
		Code:    "123456",
	})

	if err == nil || err.Error() != "otp not found" {
		t.Fatalf("expected 'otp not found' error, got %v", err)
	}
}
