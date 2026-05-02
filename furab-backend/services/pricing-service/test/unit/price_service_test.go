// Package unit contains unit tests for pricing-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"errors"
	"strings"
	"testing"

	"furab-backend/services/pricing-service/internal/model"
	"furab-backend/services/pricing-service/internal/repository"
	"furab-backend/services/pricing-service/internal/service"
	"furab-backend/services/pricing-service/test/unit/mock"

	"go.uber.org/mock/gomock"
)

type fakeOrderClient struct {
	items []model.OrderItem
	err   error
}

func (f *fakeOrderClient) GetOrderItems(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

type fakeLocationClient struct {
	distance float64
	err      error
}

func (f *fakeLocationClient) GetDeliveryDistance(ctx context.Context, orderID string) (float64, error) {
	if f.err != nil {
		return 0, f.err
	}
	return f.distance, nil
}

func newTestService(t *testing.T) (service.PriceService, *mock.MockPriceRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockPriceRepository(ctrl)
	svc := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{
			items: []model.OrderItem{
				{ProductID: "item-1", Quantity: 2, UnitPrice: 10000},
				{ProductID: "item-2", Quantity: 1, UnitPrice: 15000},
			},
		},
		&fakeLocationClient{distance: 5},
	)
	return svc, mockRepo, ctrl
}

func TestCalculatePrice_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().GetPricingRuleByType(ctx, "delivery").Return(&model.PriceRule{
		RuleID: "delivery-per-km",
		Type:   "delivery",
		Value:  5000,
	}, nil)
	mockRepo.EXPECT().GetPricingRuleByType(ctx, "service").Return(&model.PriceRule{
		RuleID: "service-percent",
		Type:   "service",
		Value:  0.05,
	}, nil)

	result, err := svc.CalculatePrice(ctx, "order-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.OrderID != "order-123" {
		t.Fatalf("expected order id order-123, got %s", result.OrderID)
	}

	if result.ItemPrice != 35000 {
		t.Fatalf("expected item price 35000, got %v", result.ItemPrice)
	}

	if result.DeliveryFee != 25000 {
		t.Fatalf("expected delivery fee 25000, got %v", result.DeliveryFee)
	}

	if result.ServiceFee != 1750 {
		t.Fatalf("expected service fee 1750, got %v", result.ServiceFee)
	}

	if result.TotalAmount != 61750 {
		t.Fatalf("expected total amount 61750, got %v", result.TotalAmount)
	}
}

func TestCalculatePrice_InvalidRequest_EmptyOrderID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CalculatePrice(context.Background(), "")
	if err == nil {
		t.Fatal("expected an error when order ID is missing")
	}

	if !errors.Is(err, service.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestCalculatePrice_NoOrderItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPriceRepository(ctrl)
	svc := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{items: []model.OrderItem{}},
		&fakeLocationClient{distance: 2},
	)

	_, err := svc.CalculatePrice(context.Background(), "order-empty")
	if !errors.Is(err, service.ErrNoOrderItems) {
		t.Fatalf("expected ErrNoOrderItems, got %v", err)
	}
}

func TestCalculatePrice_MissingPricingRule(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().GetPricingRuleByType(ctx, "delivery").Return(nil, repository.ErrPriceRuleNotFound)

	_, err := svc.CalculatePrice(ctx, "order-123")
	if err == nil {
		t.Fatal("expected error when delivery pricing rule is missing")
	}
}

func TestCalculatePrice_DistanceLogic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPriceRepository(ctrl)
	items := []model.OrderItem{
		{ProductID: "item-1", Quantity: 1, UnitPrice: 10000},
	}
	svcNear := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{items: items},
		&fakeLocationClient{distance: 1},
	)
	svcFar := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{items: items},
		&fakeLocationClient{distance: 10},
	)

	mockRepo.EXPECT().GetPricingRuleByType(gomock.Any(), "delivery").Return(&model.PriceRule{Type: "delivery", Value: 5000}, nil).Times(2)
	mockRepo.EXPECT().GetPricingRuleByType(gomock.Any(), "service").Return(&model.PriceRule{Type: "service", Value: 0.05}, nil).Times(2)

	nearRes, err := svcNear.CalculatePrice(context.Background(), "order-near")
	if err != nil {
		t.Fatalf("unexpected error for near distance: %v", err)
	}
	farRes, err := svcFar.CalculatePrice(context.Background(), "order-far")
	if err != nil {
		t.Fatalf("unexpected error for far distance: %v", err)
	}

	if nearRes.DeliveryFee != 5000 {
		t.Fatalf("expected near delivery fee 5000, got %v", nearRes.DeliveryFee)
	}
	if farRes.DeliveryFee != 50000 {
		t.Fatalf("expected far delivery fee 50000, got %v", farRes.DeliveryFee)
	}
	if farRes.TotalAmount <= nearRes.TotalAmount {
		t.Fatalf("expected far total (%v) > near total (%v)", farRes.TotalAmount, nearRes.TotalAmount)
	}
}

func TestCalculatePrice_TotalSumValidation(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().GetPricingRuleByType(ctx, "delivery").Return(&model.PriceRule{Type: "delivery", Value: 5000}, nil)
	mockRepo.EXPECT().GetPricingRuleByType(ctx, "service").Return(&model.PriceRule{Type: "service", Value: 0.05}, nil)

	res, err := svc.CalculatePrice(ctx, "order-sum-check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTotal := res.ItemPrice + res.DeliveryFee + res.ServiceFee
	if res.TotalAmount != expectedTotal {
		t.Fatalf("total mismatch: total=%v item=%v delivery=%v service=%v",
			res.TotalAmount, res.ItemPrice, res.DeliveryFee, res.ServiceFee)
	}
}

func TestCalculatePrice_LocationServiceTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPriceRepository(ctrl)
	svc := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{
			items: []model.OrderItem{{ProductID: "item-1", Quantity: 1, UnitPrice: 10000}},
		},
		&fakeLocationClient{err: errors.New("location timeout")},
	)

	_, err := svc.CalculatePrice(context.Background(), "order-timeout")
	if err == nil {
		t.Fatal("expected error when location service times out")
	}
	if !strings.Contains(err.Error(), "failed to fetch delivery distance") {
		t.Fatalf("expected delivery distance error wrapping, got: %v", err)
	}
}

func TestCalculatePrice_OrderServiceDown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPriceRepository(ctrl)
	svc := service.NewPriceService(
		mockRepo,
		&fakeOrderClient{err: errors.New("order service unavailable")},
		&fakeLocationClient{distance: 2},
	)

	_, err := svc.CalculatePrice(context.Background(), "order-down")
	if err == nil {
		t.Fatal("expected error when order service is down")
	}
	if !strings.Contains(err.Error(), "failed to fetch order items") {
		t.Fatalf("expected order items error wrapping, got: %v", err)
	}
}
