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

	if result.Status != "Valid" {
		t.Fatalf("expected status Valid, got %s", result.Status)
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

	if result.Status != "Invalid" {
		t.Fatalf("expected status Invalid, got %s", result.Status)
	}

	if result.DiscountAmount != 0 {
		t.Fatalf("expected discount amount 0, got %v", result.DiscountAmount)
	}
}
