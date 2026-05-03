// Package unit contains unit tests for pricing-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"testing"

	"furab-backend/services/pricing-service/internal/client"
	"furab-backend/services/pricing-service/internal/repository"
	"furab-backend/services/pricing-service/internal/service"
)

func TestNewPriceService_Creation(t *testing.T) {
	svc := service.NewPriceService(
		repository.NewInMemoryPriceRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyLocationClient(),
	)

	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestCalculatePrice_Success(t *testing.T) {
	svc := service.NewPriceService(
		repository.NewInMemoryPriceRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyLocationClient(),
	)

	result, err := svc.CalculatePrice(context.Background(), "order-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OrderID != "order-123" {
		t.Fatalf("expected order id order-123, got %s", result.OrderID)
	}

	if result.ItemPrice != 53000 {
		t.Fatalf("expected item price 53000, got %v", result.ItemPrice)
	}

	if result.DeliveryFee != 26000 {
		t.Fatalf("expected delivery fee 26000, got %v", result.DeliveryFee)
	}

	if result.ServiceFee != 2650 {
		t.Fatalf("expected service fee 2650, got %v", result.ServiceFee)
	}

	if result.TotalAmount != 81650 {
		t.Fatalf("expected total amount 81650, got %v", result.TotalAmount)
	}
}

func TestCalculatePrice_MissingOrderID(t *testing.T) {
	svc := service.NewPriceService(
		repository.NewInMemoryPriceRepository(),
		client.NewDummyOrderClient(),
		client.NewDummyLocationClient(),
	)

	_, err := svc.CalculatePrice(context.Background(), "")
	if err == nil {
		t.Fatal("expected an error when order ID is missing")
	}

	if err != service.ErrOrderIDRequired {
		t.Fatalf("expected ErrOrderIDRequired, got %v", err)
	}
}
