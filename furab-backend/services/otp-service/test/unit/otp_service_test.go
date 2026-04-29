package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

// ============================================================================
// Catatan: Struct dan interface di bawah ini adalah representasi dari desain.
// Pada proyek nyata, definisi ini akan berada di internal/model dan internal/service.
// ============================================================================

type OTP struct {
	UserID    string
	Code      string
	ExpiresAt int64
}

type GenerateOTPRequest struct {
	UserID string
	// Bisa berupa Email atau Phone
	Contact string 
}

type VerifyOTPRequest struct {
	UserID string
	Code   string
}

type VerifyOTPResponse struct {
	Valid   bool
	Message string
}

type GenerateOTPResponse struct {
	OTPCode string
	Message string
}

// Representasi manual dari MockOTPRepository (seharusnya digenerate otomatis oleh gomock)
type MockOTPRepository struct {
	SaveFunc func(ctx context.Context, otp *OTP) error
	FindFunc func(ctx context.Context, userID string) (*OTP, error)
}

func (m *MockOTPRepository) Save(ctx context.Context, otp *OTP) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, otp)
	}
	return nil
}

func (m *MockOTPRepository) Find(ctx context.Context, userID string) (*OTP, error) {
	if m.FindFunc != nil {
		return m.FindFunc(ctx, userID)
	}
	return nil, nil
}

// Representasi dari OTPService yang akan di-test
type OTPService interface {
	GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error)
	VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error)
}

// Implementasi Dummy Service untuk Testing
type otpServiceImpl struct {
	repo *MockOTPRepository
}

func (s *otpServiceImpl) GenerateOTP(ctx context.Context, req GenerateOTPRequest) (*GenerateOTPResponse, error) {
	// Validasi input kosong
	if req.UserID == "" || req.Contact == "" {
		return nil, errors.New("input kosong")
	}

	// Logic generate OTP (dummy untuk test)
	generatedCode := "123456" 

	otp := &OTP{
		UserID:    req.UserID,
		Code:      generatedCode,
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	// Simpan ke repository
	if err := s.repo.Save(ctx, otp); err != nil {
		return nil, err
	}

	return &GenerateOTPResponse{
		OTPCode: generatedCode,
		Message: "success",
	}, nil
}

func (s *otpServiceImpl) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, error) {
	if req.UserID == "" || req.Code == "" {
		return nil, errors.New("input kosong")
	}

	otp, err := s.repo.Find(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if otp == nil {
		return nil, errors.New("OTP tidak valid")
	}

	if otp.Code != req.Code {
		return nil, errors.New("OTP tidak valid")
	}

	if time.Now().Unix() > otp.ExpiresAt {
		return nil, errors.New("OTP expired")
	}

	return &VerifyOTPResponse{
		Valid:   true,
		Message: "OTP valid",
	}, nil
}

// ============================================================================
// UNIT TESTS MULAI DARI SINI
// ============================================================================

func TestOTPService_GenerateOTP(t *testing.T) {
	// gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - input valid", func(t *testing.T) {
		repo := &MockOTPRepository{
			SaveFunc: func(ctx context.Context, otp *OTP) error {
				if otp.UserID != "user-123" {
					t.Errorf("Expected UserID 'user-123', got %s", otp.UserID)
				}
				if len(otp.Code) != 6 {
					t.Errorf("Expected OTP Code to be 6 chars, got %s", otp.Code)
				}
				return nil // Mengembalikan nil (tidak ada error dari DB)
			},
		}

		service := &otpServiceImpl{repo: repo}

		req := GenerateOTPRequest{
			UserID:  "user-123",
			Contact: "user@mail.com",
		}

		res, err := service.GenerateOTP(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res == nil {
			t.Fatal("Expected response, got nil")
		}
		if res.Message != "success" {
			t.Errorf("Expected success message, got %s", res.Message)
		}
	})

	t.Run("Error - input kosong", func(t *testing.T) {
		repo := &MockOTPRepository{
			SaveFunc: func(ctx context.Context, otp *OTP) error {
				t.Fatal("Repository Save() should not be called when validation fails")
				return nil
			},
		}

		service := &otpServiceImpl{repo: repo}

		// Skenario: UserID kosong, atau Contact kosong
		invalidRequests := []GenerateOTPRequest{
			{UserID: "", Contact: "user@mail.com"},
			{UserID: "user-123", Contact: ""},
			{UserID: "", Contact: ""},
		}

		for _, req := range invalidRequests {
			res, err := service.GenerateOTP(context.Background(), req)

			if err == nil {
				t.Fatalf("Expected validation error for req: %+v, got none", req)
			}
			if err.Error() != "input kosong" {
				t.Errorf("Expected 'input kosong' error, got: %v", err)
			}
			if res != nil {
				t.Errorf("Expected nil response, got %v", res)
			}
		}
	})
}

func TestOTPService_VerifyOTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - OTP valid", func(t *testing.T) {
		repo := &MockOTPRepository{
			FindFunc: func(ctx context.Context, userID string) (*OTP, error) {
				return &OTP{
					UserID:    "user-123",
					Code:      "123456",
					ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				}, nil
			},
		}

		service := &otpServiceImpl{repo: repo}
		req := VerifyOTPRequest{UserID: "user-123", Code: "123456"}

		res, err := service.VerifyOTP(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res == nil || !res.Valid {
			t.Fatal("Expected OTP to be valid")
		}
	})

	t.Run("Error - OTP tidak valid", func(t *testing.T) {
		repo := &MockOTPRepository{
			FindFunc: func(ctx context.Context, userID string) (*OTP, error) {
				// Return OTP with different code
				return &OTP{
					UserID:    "user-123",
					Code:      "654321",
					ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				}, nil
			},
		}

		service := &otpServiceImpl{repo: repo}
		req := VerifyOTPRequest{UserID: "user-123", Code: "123456"}

		res, err := service.VerifyOTP(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "OTP tidak valid" {
			t.Errorf("Expected 'OTP tidak valid' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - OTP expired", func(t *testing.T) {
		repo := &MockOTPRepository{
			FindFunc: func(ctx context.Context, userID string) (*OTP, error) {
				// Return expired OTP
				return &OTP{
					UserID:    "user-123",
					Code:      "123456",
					ExpiresAt: time.Now().Add(-5 * time.Minute).Unix(),
				}, nil
			},
		}

		service := &otpServiceImpl{repo: repo}
		req := VerifyOTPRequest{UserID: "user-123", Code: "123456"}

		res, err := service.VerifyOTP(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "OTP expired" {
			t.Errorf("Expected 'OTP expired' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}
