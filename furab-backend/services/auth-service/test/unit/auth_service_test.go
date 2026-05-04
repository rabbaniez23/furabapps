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

// Sentinel errors returned by mocks so assertions use errors.Is, not string equality.
var (
	errDependencyUserCreate    = errors.New("dependency: user create failed")
	errDependencyUserGet       = errors.New("dependency: user get failed")
	errDependencyOTPGenerate   = errors.New("dependency: otp generate failed")
	errDependencyOTPVerify     = errors.New("dependency: otp verify failed")
	errDependencyTokenGenerate = errors.New("dependency: token generate failed")
	errDependencyTokenValidate = errors.New("dependency: token validate failed")
)

const (
	validPhone = "08123456789"
	validOTP   = "123456"
)

func setupAuthService(t *testing.T) (
	mockUser *mock_service.MockUserService,
	mockOTP *mock_service.MockOTPService,
	mockToken *mock_service.MockTokenGenerator,
	svc service.AuthService,
) {
	t.Helper()
	// gomock.NewController registers Finish on t.Cleanup for *testing.T — no defer needed.
	ctrl := gomock.NewController(t)
	mockUser = mock_service.NewMockUserService(ctrl)
	mockOTP = mock_service.NewMockOTPService(ctrl)
	mockToken = mock_service.NewMockTokenGenerator(ctrl)
	svc = service.NewAuthService(mockUser, mockOTP, mockToken)
	return mockUser, mockOTP, mockToken, svc
}

func TestAuthService_Register(t *testing.T) {
	ctx := context.Background()

	t.Run("success_creates_user_and_sends_otp", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)
		mockUser.EXPECT().
			CreateUser(gomock.Any(), validPhone).
			Return(nil)
		mockOTP.EXPECT().
			GenerateOTP(gomock.Any(), validPhone).
			Return(nil)

		res, err := svc.Register(ctx, validPhone)
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		if res == nil || res.Status != "success" {
			t.Fatalf("Register: got res %#v", res)
		}
	})

	t.Run("success_phone_input_with_separators_calls_deps_with_canonical_contact", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)

		input := "08123 456789"
		mockUser.EXPECT().
			CreateUser(gomock.Any(), validPhone).
			Return(nil)
		mockOTP.EXPECT().
			GenerateOTP(gomock.Any(), validPhone).
			Return(nil)

		res, err := svc.Register(ctx, input)
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		if res == nil || res.Status != "success" {
			t.Fatalf("Register: got res %#v", res)
		}
	})

	t.Run("validation_error_missing_contact", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.Register(ctx, "")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactRequired) {
			t.Fatalf("expected ErrContactRequired, got %v", err)
		}
	})

	t.Run("validation_error_invalid_contact_format", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.Register(ctx, "not-a-phone-or-email")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactInvalidFormat) {
			t.Fatalf("expected ErrContactInvalidFormat, got %v", err)
		}
	})

	t.Run("validation_error_phone_with_invalid_character", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)

		res, err := svc.Register(ctx, "08123x456789")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactInvalidFormat) {
			t.Fatalf("expected ErrContactInvalidFormat, got %v", err)
		}
	})

	t.Run("dependency_error_user_service", func(t *testing.T) {
		mockUser, _, _, svc := setupAuthService(t)
		mockUser.EXPECT().
			CreateUser(gomock.Any(), validPhone).
			Return(errDependencyUserCreate)

		res, err := svc.Register(ctx, validPhone)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyUserCreate) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("dependency_error_otp_service", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)
		mockUser.EXPECT().
			CreateUser(gomock.Any(), validPhone).
			Return(nil)
		mockOTP.EXPECT().
			GenerateOTP(gomock.Any(), validPhone).
			Return(errDependencyOTPGenerate)

		res, err := svc.Register(ctx, validPhone)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyOTPGenerate) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.Register(cancelledCtx, validPhone)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestAuthService_RequestOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("success_sends_otp", func(t *testing.T) {
		_, mockOTP, _, svc := setupAuthService(t)
		mockOTP.EXPECT().
			GenerateOTP(gomock.Any(), validPhone).
			Return(nil)

		res, err := svc.RequestOTP(ctx, validPhone)
		if err != nil {
			t.Fatalf("RequestOTP: %v", err)
		}
		if res == nil || res.Status != "success" {
			t.Fatalf("RequestOTP: got res %#v", res)
		}
	})

	t.Run("validation_error_missing_contact", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.RequestOTP(ctx, "")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactRequired) {
			t.Fatalf("expected ErrContactRequired, got %v", err)
		}
	})

	t.Run("validation_error_invalid_contact_format", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.RequestOTP(ctx, "@@@")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactInvalidFormat) {
			t.Fatalf("expected ErrContactInvalidFormat, got %v", err)
		}
	})

	t.Run("dependency_error_otp_service", func(t *testing.T) {
		_, mockOTP, _, svc := setupAuthService(t)
		mockOTP.EXPECT().
			GenerateOTP(gomock.Any(), validPhone).
			Return(errDependencyOTPGenerate)

		res, err := svc.RequestOTP(ctx, validPhone)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyOTPGenerate) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.RequestOTP(cancelledCtx, validPhone)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestAuthService_VerifyOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("success_returns_access_token", func(t *testing.T) {
		mockUser, mockOTP, mockToken, svc := setupAuthService(t)
		user := &model.User{ID: "user-123"}
		wantToken := "access-token"

		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(true, nil)
		mockUser.EXPECT().
			GetUser(gomock.Any(), validPhone).
			Return(user, nil)
		mockToken.EXPECT().
			GenerateToken(user.ID).
			Return(wantToken, nil)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if err != nil {
			t.Fatalf("VerifyOTP: %v", err)
		}
		if res == nil || res.Status != "success" || res.AccessToken != wantToken {
			t.Fatalf("VerifyOTP: got %#v", res)
		}
	})

	t.Run("validation_error_missing_inputs", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.VerifyOTP(ctx, validPhone, "")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrInputRequired) {
			t.Fatalf("expected ErrInputRequired, got %v", err)
		}
	})

	t.Run("validation_error_invalid_contact_format", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		res, err := svc.VerifyOTP(ctx, "bad-contact", validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrContactInvalidFormat) {
			t.Fatalf("expected ErrContactInvalidFormat, got %v", err)
		}
	})

	t.Run("dependency_error_otp_service", func(t *testing.T) {
		_, mockOTP, _, svc := setupAuthService(t)
		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(false, errDependencyOTPVerify)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyOTPVerify) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("otp_not_valid_returns_ErrOTPInvalid", func(t *testing.T) {
		_, mockOTP, _, svc := setupAuthService(t)
		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, "wrong").
			Return(false, nil)

		res, err := svc.VerifyOTP(ctx, validPhone, "wrong")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrOTPInvalid) {
			t.Fatalf("expected ErrOTPInvalid, got %v", err)
		}
	})

	t.Run("not_found_user_nil_after_valid_otp", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)
		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(true, nil)
		mockUser.EXPECT().
			GetUser(gomock.Any(), validPhone).
			Return(nil, nil)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("dependency_error_user_service_on_get", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)

		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(true, nil)
		mockUser.EXPECT().
			GetUser(gomock.Any(), validPhone).
			Return(nil, errDependencyUserGet)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyUserGet) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("error_user_record_without_id", func(t *testing.T) {
		mockUser, mockOTP, _, svc := setupAuthService(t)

		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(true, nil)
		mockUser.EXPECT().
			GetUser(gomock.Any(), validPhone).
			Return(&model.User{ID: ""}, nil)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrUserIDMissing) {
			t.Fatalf("expected ErrUserIDMissing, got %v", err)
		}
	})

	t.Run("dependency_error_token_generator", func(t *testing.T) {
		mockUser, mockOTP, mockToken, svc := setupAuthService(t)
		user := &model.User{ID: "user-123"}

		mockOTP.EXPECT().
			VerifyOTP(gomock.Any(), validPhone, validOTP).
			Return(true, nil)
		mockUser.EXPECT().
			GetUser(gomock.Any(), validPhone).
			Return(user, nil)
		mockToken.EXPECT().
			GenerateToken(user.ID).
			Return("", errDependencyTokenGenerate)

		res, err := svc.VerifyOTP(ctx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyTokenGenerate) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.VerifyOTP(cancelledCtx, validPhone, validOTP)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success_token_valid", func(t *testing.T) {
		_, _, mockToken, svc := setupAuthService(t)
		token := "signed-token"
		mockToken.EXPECT().
			ValidateToken(token).
			Return(true, nil)

		res, err := svc.ValidateToken(ctx, token)
		if err != nil {
			t.Fatalf("ValidateToken: %v", err)
		}
		if res == nil || res.Status != "valid" {
			t.Fatalf("ValidateToken: got %#v", res)
		}
	})

	t.Run("validation_empty_token_marked_invalid_without_dependency_call", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		// No TokenGenerator EXPECT — service short-circuits on empty token.

		res, err := svc.ValidateToken(ctx, "")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil || res.Status != "invalid" {
			t.Fatalf("expected invalid status, got %#v", res)
		}
	})

	t.Run("dependency_error_token_validator", func(t *testing.T) {
		_, _, mockToken, svc := setupAuthService(t)
		token := "any-token"
		mockToken.EXPECT().
			ValidateToken(token).
			Return(false, errDependencyTokenValidate)

		res, err := svc.ValidateToken(ctx, token)
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errDependencyTokenValidate) {
			t.Fatalf("expected dependency error, got %v", err)
		}
	})

	t.Run("invalid_token_response_without_error", func(t *testing.T) {
		_, _, mockToken, svc := setupAuthService(t)
		token := "expired-token"
		mockToken.EXPECT().
			ValidateToken(token).
			Return(false, nil)

		res, err := svc.ValidateToken(ctx, token)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil || res.Status != "invalid" {
			t.Fatalf("expected invalid status, got %#v", res)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		_, _, _, svc := setupAuthService(t)
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.ValidateToken(cancelledCtx, "any-token")
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}
