// Package unit contains unit tests for the ride order service.
// Unit tests do NOT access any database or external service.
// All dependencies are mocked using gomock.
package unit

import (
	"context"
	"testing"
	"time"

	"furab-backend/services/ride-order-service/internal/model"
	"furab-backend/services/ride-order-service/internal/repository"
	"furab-backend/services/ride-order-service/internal/service"
	"furab-backend/services/ride-order-service/test/unit/mock"
	"furab-backend/shared/event"

	"go.uber.org/mock/gomock"
)

// --- Helper Functions ---

// newTestService creates a new OrderService with mocked dependencies.
func newTestService(t *testing.T) (service.OrderService, *mock.MockOrderRepository, *mock.MockEventPublisher, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockOrderRepository(ctrl)
	mockPublisher := mock.NewMockEventPublisher(ctrl)
	svc := service.NewOrderService(mockRepo, mockPublisher)
	return svc, mockRepo, mockPublisher, ctrl
}

// validCreateRequest returns a valid CreateRideOrderRequest for testing.
func validCreateRequest() *model.CreateRideOrderRequest {
	return &model.CreateRideOrderRequest{
		UserID: "user-123",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Monas, Jakarta Pusat",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.1751,
			Longitude: 106.8650,
			Address:   "Ancol, Jakarta Utara",
		},
	}
}

// sampleOrder returns a sample RideOrder for testing.
func sampleOrder() *model.RideOrder {
	return &model.RideOrder{
		ID:     "order-abc-123",
		UserID: "user-123",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Monas, Jakarta Pusat",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.1751,
			Longitude: 106.8650,
			Address:   "Ancol, Jakarta Utara",
		},
		Status:    model.RideStatusPending,
		Fare:      18500,
		Distance:  4.2,
		EstimatedDuration: 9,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// --- Test Cases: CreateOrder ---

// TestCreateOrder_Success tests creating a ride order with valid data.
// Expected: order created with PENDING status, ride.created event published.
func TestCreateOrder_Success(t *testing.T) {
	svc, mockRepo, mockPublisher, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := validCreateRequest()

	// Expect repository Create to be called
	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	// Expect event to be published
	mockPublisher.EXPECT().
		Publish(ctx, event.TopicRideCreated, gomock.Any()).
		Return(nil)

	order, err := svc.CreateOrder(ctx, req)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if order == nil {
		t.Fatal("expected order, got nil")
	}
	if order.Status != model.RideStatusPending {
		t.Errorf("expected status PENDING, got: %s", order.Status)
	}
	if order.UserID != req.UserID {
		t.Errorf("expected user ID %s, got: %s", req.UserID, order.UserID)
	}
	if order.Fare <= 0 {
		t.Error("expected fare > 0")
	}
	if order.ID == "" {
		t.Error("expected non-empty order ID")
	}
}

// TestCreateOrder_NilRequest tests creating a ride order with nil request.
// Expected: error returned.
func TestCreateOrder_NilRequest(t *testing.T) {
	svc, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CreateOrder(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
}

// TestCreateOrder_InvalidPickup tests creating a ride order with empty pickup address.
// Expected: validation error returned.
func TestCreateOrder_InvalidPickup(t *testing.T) {
	svc, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := validCreateRequest()
	req.PickupLocation.Address = "" // invalid

	_, err := svc.CreateOrder(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid pickup")
	}
}

// TestCreateOrder_InvalidDropoff tests creating a ride order with empty dropoff address.
// Expected: validation error returned.
func TestCreateOrder_InvalidDropoff(t *testing.T) {
	svc, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := validCreateRequest()
	req.DropoffLocation.Address = "" // invalid

	_, err := svc.CreateOrder(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid dropoff")
	}
}

// TestCreateOrder_EmptyUserID tests creating a ride order with empty user ID.
// Expected: validation error returned.
func TestCreateOrder_EmptyUserID(t *testing.T) {
	svc, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := validCreateRequest()
	req.UserID = ""

	_, err := svc.CreateOrder(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
}

// --- Test Cases: GetOrder ---

// TestGetOrder_Success tests retrieving an existing ride order.
// Expected: order returned successfully.
func TestGetOrder_Success(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expected := sampleOrder()

	mockRepo.EXPECT().
		GetByID(ctx, expected.ID).
		Return(expected, nil)

	order, err := svc.GetOrder(ctx, expected.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if order.ID != expected.ID {
		t.Errorf("expected order ID %s, got: %s", expected.ID, order.ID)
	}
}

// TestGetOrder_NotFound tests retrieving a non-existent order.
// Expected: ErrOrderNotFound returned.
func TestGetOrder_NotFound(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByID(ctx, "non-existent").
		Return(nil, repository.ErrOrderNotFound)

	_, err := svc.GetOrder(ctx, "non-existent")
	if err != service.ErrOrderNotFound {
		t.Fatalf("expected ErrOrderNotFound, got: %v", err)
	}
}

// TestGetOrder_EmptyID tests retrieving an order with empty ID.
// Expected: ErrInvalidRequest returned.
func TestGetOrder_EmptyID(t *testing.T) {
	svc, _, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GetOrder(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

// --- Test Cases: AssignDriver ---

// TestAssignDriver_Success tests assigning a driver to a PENDING order.
// Expected: status transitions to ASSIGNED, ride.assigned event published.
func TestAssignDriver_Success(t *testing.T) {
	svc, mockRepo, mockPublisher, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusPending

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil)

	mockPublisher.EXPECT().
		Publish(ctx, event.TopicRideAssigned, gomock.Any()).
		Return(nil)

	result, err := svc.AssignDriver(ctx, order.ID, "driver-456")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Status != model.RideStatusAssigned {
		t.Errorf("expected status ASSIGNED, got: %s", result.Status)
	}
	if result.DriverID != "driver-456" {
		t.Errorf("expected driver ID driver-456, got: %s", result.DriverID)
	}
}

// TestAssignDriver_InvalidStatus tests assigning a driver to a COMPLETED order.
// Expected: ErrInvalidTransition returned.
func TestAssignDriver_InvalidStatus(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusCompleted // cannot assign

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	_, err := svc.AssignDriver(ctx, order.ID, "driver-456")
	if err != service.ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

// TestAssignDriver_AlreadyAssigned tests assigning a driver when one is already assigned.
// Expected: ErrDriverAlreadyAssigned returned.
func TestAssignDriver_AlreadyAssigned(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusPending
	order.DriverID = "existing-driver" // already assigned

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	_, err := svc.AssignDriver(ctx, order.ID, "new-driver")
	if err != service.ErrDriverAlreadyAssigned {
		t.Fatalf("expected ErrDriverAlreadyAssigned, got: %v", err)
	}
}

// --- Test Cases: StartRide ---

// TestStartRide_Success tests starting an ASSIGNED ride.
// Expected: status transitions to STARTED, ride.started event published.
func TestStartRide_Success(t *testing.T) {
	svc, mockRepo, mockPublisher, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusAssigned
	order.DriverID = "driver-456"

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil)

	mockPublisher.EXPECT().
		Publish(ctx, event.TopicRideStarted, gomock.Any()).
		Return(nil)

	result, err := svc.StartRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Status != model.RideStatusStarted {
		t.Errorf("expected status STARTED, got: %s", result.Status)
	}
}

// TestStartRide_InvalidStatus tests starting a PENDING ride (should fail).
// Expected: ErrInvalidTransition returned.
func TestStartRide_InvalidStatus(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusPending // cannot start from PENDING

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	_, err := svc.StartRide(ctx, order.ID)
	if err != service.ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

// --- Test Cases: CompleteRide ---

// TestCompleteRide_Success tests completing a STARTED ride.
// Expected: status transitions to COMPLETED, ride.completed event published.
func TestCompleteRide_Success(t *testing.T) {
	svc, mockRepo, mockPublisher, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusStarted
	order.DriverID = "driver-456"

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil)

	mockPublisher.EXPECT().
		Publish(ctx, event.TopicRideCompleted, gomock.Any()).
		Return(nil)

	result, err := svc.CompleteRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Status != model.RideStatusCompleted {
		t.Errorf("expected status COMPLETED, got: %s", result.Status)
	}
}

// TestCompleteRide_InvalidStatus tests completing a PENDING ride (should fail).
// Expected: ErrInvalidTransition returned.
func TestCompleteRide_InvalidStatus(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusPending // cannot complete from PENDING

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	_, err := svc.CompleteRide(ctx, order.ID)
	if err != service.ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

// --- Test Cases: CancelRide ---

// TestCancelRide_Success tests cancelling a PENDING ride.
// Expected: status transitions to CANCELLED.
func TestCancelRide_Success(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusPending

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil)

	result, err := svc.CancelRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if result.Status != model.RideStatusCancelled {
		t.Errorf("expected status CANCELLED, got: %s", result.Status)
	}
}

// TestCancelRide_AlreadyCompleted tests cancelling a COMPLETED ride (should fail).
// Expected: ErrInvalidTransition returned.
func TestCancelRide_AlreadyCompleted(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	order := sampleOrder()
	order.Status = model.RideStatusCompleted // cannot cancel

	mockRepo.EXPECT().
		GetByID(ctx, order.ID).
		Return(order, nil)

	_, err := svc.CancelRide(ctx, order.ID)
	if err != service.ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition, got: %v", err)
	}
}

// --- Test Cases: GetUserOrders ---

// TestGetUserOrders_Success tests retrieving orders for a valid user.
// Expected: list of orders returned with count.
func TestGetUserOrders_Success(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedOrders := []*model.RideOrder{sampleOrder(), sampleOrder()}

	mockRepo.EXPECT().
		GetByUserID(ctx, "user-123", 10, 0).
		Return(expectedOrders, nil)

	mockRepo.EXPECT().
		CountByUserID(ctx, "user-123").
		Return(2, nil)

	orders, total, err := svc.GetUserOrders(ctx, "user-123", 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got: %d", len(orders))
	}
	if total != 2 {
		t.Errorf("expected total 2, got: %d", total)
	}
}

// TestGetUserOrders_Empty tests retrieving orders for a user with no orders.
// Expected: empty list, total = 0.
func TestGetUserOrders_Empty(t *testing.T) {
	svc, mockRepo, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByUserID(ctx, "user-empty", 10, 0).
		Return([]*model.RideOrder{}, nil)

	mockRepo.EXPECT().
		CountByUserID(ctx, "user-empty").
		Return(0, nil)

	orders, total, err := svc.GetUserOrders(ctx, "user-empty", 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got: %d", len(orders))
	}
	if total != 0 {
		t.Errorf("expected total 0, got: %d", total)
	}
}
