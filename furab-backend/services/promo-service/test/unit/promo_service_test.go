// Package unit contains unit tests for promo-service.
package unit

import (
	"context"
	"testing"

	"furab-backend/services/promo-service/internal/client"
	"furab-backend/services/promo-service/internal/repository"
	"furab-backend/services/promo-service/internal/service"
)

func TestNewPromoService_Creation(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestValidatePromo_Success(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "DISKONHEMAT", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsValid {
		t.Fatalf("expected status IsValid=true, got false. error: %s", result.ErrorMessage)
	}

	if result.DiscountAmount <= 0 {
		t.Fatalf("expected discount amount > 0, got %v", result.DiscountAmount)
	}

	if result.FinalAmount != 90000 {
		t.Fatalf("expected final amount 90000, got %v", result.FinalAmount)
	}
}

func TestValidatePromo_InvalidCode(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "UNKNOWN", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected status IsValid=false, got true")
	}

	if result.ErrorMessage == "" {
		t.Fatalf("expected error message for invalid promo")
	}

	if result.DiscountAmount != 0 {
		t.Fatalf("expected discount amount 0, got %v", result.DiscountAmount)
	}
}

// mockUserClient allows controlling the output of ValidateUserPromo for testing
type mockUserClient struct {
	valid bool
}

func (m *mockUserClient) ValidateUserPromo(ctx context.Context, userID, promoCode string) (bool, error) {
	return m.valid, nil
}

// mockOrderClient allows controlling the output of ValidateOrderPromo for testing
type mockOrderClient struct {
	valid bool
}

func (m *mockOrderClient) ValidateOrderPromo(ctx context.Context, orderID, promoCode string) (bool, error) {
	return m.valid, nil
}

// 1. Logika Pembatasan (Limits & Constraints)

func TestValidatePromo_Expired(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "EXPIRED", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected promo to be invalid, but got valid")
	}

	if result.ErrorMessage != "promo has expired" {
		t.Fatalf("expected error 'promo has expired', got '%s'", result.ErrorMessage)
	}
}

func TestValidatePromo_BelowMinimumPurchase(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "DISKONHEMAT", "user-1", "order-1", 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected promo to be invalid, but got valid")
	}

	if result.ErrorMessage != "minimum purchase not met. minimum: 50000.00" {
		t.Fatalf("expected error regarding minimum purchase, got '%s'", result.ErrorMessage)
	}
}

func TestValidatePromo_UsageLimitReached(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "FULL", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected promo to be invalid, but got valid")
	}

	if result.ErrorMessage != "promo usage limit exceeded" {
		t.Fatalf("expected error 'promo usage limit exceeded', got '%s'", result.ErrorMessage)
	}
}

// 2. Akurasi Perhitungan Diskon (Calculation Logic)

func TestCalculateDiscount_PercentageWithCap(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	// BIGPERCENT = 50% discount, max cap = 10000
	// 50% of 100000 = 50000, but cap is 10000.
	result, err := svc.ValidatePromo(context.Background(), "BIGPERCENT", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsValid {
		t.Fatalf("expected promo to be valid")
	}

	if result.DiscountAmount != 10000 {
		t.Fatalf("expected discount amount to be capped at 10000, got %v", result.DiscountAmount)
	}

	if result.FinalAmount != 90000 {
		t.Fatalf("expected final amount 90000, got %v", result.FinalAmount)
	}
}

func TestCalculateDiscount_FixedAmount(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyUserClient(),
	)

	// FIXED50 = fixed 50000 discount
	result, err := svc.ValidatePromo(context.Background(), "FIXED50", "user-1", "order-1", 150000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.IsValid {
		t.Fatalf("expected promo to be valid")
	}

	if result.DiscountAmount != 50000 {
		t.Fatalf("expected discount amount 50000, got %v", result.DiscountAmount)
	}

	if result.FinalAmount != 100000 {
		t.Fatalf("expected final amount 100000, got %v", result.FinalAmount)
	}
}

// 3. Validasi User & Order (Eligibility)

func TestValidatePromo_UserNotEligible(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		client.NewDummyOrderClient(),
		&mockUserClient{valid: false},
	)

	result, err := svc.ValidatePromo(context.Background(), "DISKONHEMAT", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected promo to be invalid due to user eligibility")
	}

	if result.ErrorMessage != "user is not eligible for this promo" {
		t.Fatalf("expected error 'user is not eligible for this promo', got '%s'", result.ErrorMessage)
	}
}

func TestValidatePromo_InvalidOrderType(t *testing.T) {
	svc := service.NewPromoService(
		repository.NewInMemoryPromoRepository(),
		&mockOrderClient{valid: false},
		client.NewDummyUserClient(),
	)

	result, err := svc.ValidatePromo(context.Background(), "DISKONHEMAT", "user-1", "order-1", 100000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.IsValid {
		t.Fatalf("expected promo to be invalid due to order eligibility")
	}

	if result.ErrorMessage != "order does not meet promo conditions" {
		t.Fatalf("expected error 'order does not meet promo conditions', got '%s'", result.ErrorMessage)
	}
}
