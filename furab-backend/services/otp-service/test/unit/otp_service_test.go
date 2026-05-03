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

var (
	errRepoSave = errors.New("repo save error")
	errRepoFind = errors.New("repo find error")
)

type otpArgMatcher struct {
	match func(*model.OTP) bool
}

func (m otpArgMatcher) Matches(x any) bool {
	o, ok := x.(*model.OTP)
	if !ok {
		return false
	}
	return m.match(o)
}

func (m otpArgMatcher) String() string {
	return "matches *model.OTP predicate"
}

func matchOTP(match func(*model.OTP) bool) gomock.Matcher {
	return otpArgMatcher{match: match}
}

func setupOTPService(t *testing.T) (*gomock.Controller, *mock_repository.MockOTPRepository, service.OTPService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockOTPRepository(ctrl)
	return ctrl, mockRepo, service.NewOTPService(mockRepo)
}

func TestOTPService_GenerateOTP(t *testing.T) {
	ctrl, mockRepo, svc := setupOTPService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		now := time.Now().Unix()
		mockRepo.EXPECT().
			Save(gomock.Any(), matchOTP(func(o *model.OTP) bool {
				// Phone trimmed + dummy code + expiry in the future.
				return o != nil && o.Phone == "08123" && o.Code == "123456" && o.ExpiresAt > now
			})).
			Return(nil)

		res, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
			Contact: " 08123 ",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil response")
		}
		if res.OTPCode == "" {
			t.Fatal("expected non-empty OTPCode")
		}
	})

	t.Run("validation_error_contact_required", func(t *testing.T) {
		res, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
			Contact: " ",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactRequired) {
			t.Fatalf("expected ErrContactRequired, got %v", err)
		}
	})

	t.Run("repository_error_save", func(t *testing.T) {
		mockRepo.EXPECT().
			Save(gomock.Any(), matchOTP(func(o *model.OTP) bool {
				return o != nil && o.Phone == "08123" && o.Code == "123456"
			})).
			Return(errRepoSave)

		res, err := svc.GenerateOTP(context.Background(), service.GenerateOTPRequest{
			Contact: "08123",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errRepoSave) {
			t.Fatalf("expected repo save error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.GenerateOTP(cancelledCtx, service.GenerateOTPRequest{
			Contact: "08123",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestOTPService_VerifyOTP(t *testing.T) {
	ctrl, mockRepo, svc := setupOTPService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByPhone(gomock.Any(), "08123").Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}, nil)

		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: " 08123 ",
			Code:    " 123456 ",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil || !res.Valid {
			t.Fatalf("expected valid response, got %#v", res)
		}
	})

	t.Run("validation_error_contact_required", func(t *testing.T) {
		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: " ",
			Code:    "123456",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactRequired) {
			t.Fatalf("expected ErrContactRequired, got %v", err)
		}
	})

	t.Run("validation_error_otp_required", func(t *testing.T) {
		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: "08123",
			Code:    " ",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrOTPRequired) {
			t.Fatalf("expected ErrOTPRequired, got %v", err)
		}
	})

	t.Run("business_error_not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByPhone(gomock.Any(), "99999").Return(nil, nil)

		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: "99999",
			Code:    "123456",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrOTPNotFound) {
			t.Fatalf("expected ErrOTPNotFound, got %v", err)
		}
	})

	t.Run("business_error_invalid_code", func(t *testing.T) {
		mockRepo.EXPECT().FindByPhone(gomock.Any(), "08123").Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		}, nil)

		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: "08123",
			Code:    "000000",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrOTPInvalid) {
			t.Fatalf("expected ErrOTPInvalid, got %v", err)
		}
	})

	t.Run("business_error_expired", func(t *testing.T) {
		mockRepo.EXPECT().FindByPhone(gomock.Any(), "08123").Return(&model.OTP{
			Phone:     "08123",
			Code:      "123456",
			ExpiresAt: time.Now().Add(-5 * time.Minute).Unix(),
		}, nil)

		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: "08123",
			Code:    "123456",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrOTPExpired) {
			t.Fatalf("expected ErrOTPExpired, got %v", err)
		}
	})

	t.Run("repository_error_find_by_phone", func(t *testing.T) {
		mockRepo.EXPECT().FindByPhone(gomock.Any(), "08123").Return(nil, errRepoFind)

		res, err := svc.VerifyOTP(context.Background(), service.VerifyOTPRequest{
			Contact: "08123",
			Code:    "123456",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errRepoFind) {
			t.Fatalf("expected repo find error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.VerifyOTP(cancelledCtx, service.VerifyOTPRequest{
			Contact: "08123",
			Code:    "123456",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}