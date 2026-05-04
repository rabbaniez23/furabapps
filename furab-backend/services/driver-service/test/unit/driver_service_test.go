// Package unit contains unit tests for the driver service.
// All dependencies are mocked using gomock. No database access.
package unit

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"furab-backend/services/driver-service/internal/model"
	"furab-backend/services/driver-service/internal/repository"
	mock_repository "furab-backend/services/driver-service/internal/repository/mock"
	"furab-backend/services/driver-service/internal/service"

	"go.uber.org/mock/gomock"
)

// --- Helper Functions ---

func newTestService(t *testing.T) (service.DriverService, *mock_repository.MockDriverRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)
	return svc, mockRepo, ctrl
}

func sampleDriver() *model.Driver {
	return &model.Driver{
		DriverID:    "driver-123",
		Name:        "John Driver",
		Phone:       "081234567890",
		VehicleType: "motorcycle",
		Status:      model.DriverStatusOffline,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// ========================================
// CreateDriver
// ========================================

func TestCreateDriver_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)

	resp, err := svc.CreateDriver(ctx, &model.CreateDriverRequest{
		DriverID:    "driver-123",
		Name:        "John Driver",
		Phone:       "081234567890",
		VehicleType: "motorcycle",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
	if resp.DriverID != "driver-123" {
		t.Errorf("expected driver-123, got: %s", resp.DriverID)
	}
}

func TestCreateDriver_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CreateDriver(context.Background(), nil)
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestCreateDriver_EmptyDriverID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CreateDriver(context.Background(), &model.CreateDriverRequest{
		Name: "John", Phone: "08123", VehicleType: "car",
	})
	if err == nil {
		t.Fatal("expected error for empty driver_id")
	}
}

func TestCreateDriver_EmptyName(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CreateDriver(context.Background(), &model.CreateDriverRequest{
		DriverID: "d1", Phone: "08123", VehicleType: "car",
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreateDriver_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	repoErr := errors.New("db error")
	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(repoErr)

	_, err := svc.CreateDriver(context.Background(), &model.CreateDriverRequest{
		DriverID: "d1", Name: "John", Phone: "08123", VehicleType: "car",
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

func TestCreateDriver_DefaultStatusOffline(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, driver *model.Driver) error {
			if driver.Status != model.DriverStatusOffline {
				t.Errorf("expected OFFLINE status, got: %s", driver.Status)
			}
			return nil
		})

	_, err := svc.CreateDriver(context.Background(), &model.CreateDriverRequest{
		DriverID: "d1", Name: "John", Phone: "08123", VehicleType: "car",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ========================================
// GetDriver
// ========================================

func TestGetDriver_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	expected := sampleDriver()
	mockRepo.EXPECT().FindByID(gomock.Any(), expected.DriverID).Return(expected, nil)

	driver, err := svc.GetDriver(context.Background(), expected.DriverID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if driver.DriverID != expected.DriverID {
		t.Errorf("expected %s, got: %s", expected.DriverID, driver.DriverID)
	}
}

func TestGetDriver_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByID(gomock.Any(), "unknown").Return(nil, repository.ErrDriverNotFound)

	_, err := svc.GetDriver(context.Background(), "unknown")
	if err != service.ErrDriverNotFound {
		t.Fatalf("expected ErrDriverNotFound, got: %v", err)
	}
}

func TestGetDriver_EmptyID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GetDriver(context.Background(), "")
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// ========================================
// UpdateDriver
// ========================================

func TestUpdateDriver_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	resp, err := svc.UpdateDriver(context.Background(), existing.DriverID, &model.UpdateDriverRequest{
		Name: "Jane", Phone: "089876", VehicleType: "car",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestUpdateDriver_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByID(gomock.Any(), "unknown").Return(nil, repository.ErrDriverNotFound)

	_, err := svc.UpdateDriver(context.Background(), "unknown", &model.UpdateDriverRequest{
		Name: "Jane", Phone: "089876", VehicleType: "car",
	})
	if err != service.ErrDriverNotFound {
		t.Fatalf("expected ErrDriverNotFound, got: %v", err)
	}
}

func TestUpdateDriver_EmptyDriverID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateDriver(context.Background(), "", &model.UpdateDriverRequest{
		Name: "Jane", Phone: "089876", VehicleType: "car",
	})
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestUpdateDriver_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateDriver(context.Background(), "d1", nil)
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestUpdateDriver_EmptyName(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateDriver(context.Background(), "d1", &model.UpdateDriverRequest{
		Name: "", Phone: "08123", VehicleType: "car",
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUpdateDriver_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	repoErr := errors.New("update failed")
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(repoErr)

	_, err := svc.UpdateDriver(context.Background(), existing.DriverID, &model.UpdateDriverRequest{
		Name: "Jane", Phone: "089876", VehicleType: "car",
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// ========================================
// UpdateStatus
// ========================================

func TestUpdateStatus_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), existing.DriverID, "ONLINE").Return(nil)

	resp, err := svc.UpdateStatus(context.Background(), existing.DriverID, "ONLINE")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestUpdateStatus_CaseInsensitive(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().UpdateStatus(gomock.Any(), existing.DriverID, "OFFLINE").Return(nil)

	resp, err := svc.UpdateStatus(context.Background(), existing.DriverID, "offline")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateStatus(context.Background(), "d1", "INVALID")
	if err != service.ErrStatusInvalid {
		t.Fatalf("expected ErrStatusInvalid, got: %v", err)
	}
}

func TestUpdateStatus_EmptyStatus(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateStatus(context.Background(), "d1", "")
	if err != service.ErrStatusInvalid {
		t.Fatalf("expected ErrStatusInvalid, got: %v", err)
	}
}

func TestUpdateStatus_DriverNotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByID(gomock.Any(), "unknown").Return(nil, repository.ErrDriverNotFound)

	_, err := svc.UpdateStatus(context.Background(), "unknown", "ONLINE")
	if err != service.ErrDriverNotFound {
		t.Fatalf("expected ErrDriverNotFound, got: %v", err)
	}
}

func TestUpdateStatus_EmptyDriverID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateStatus(context.Background(), "", "ONLINE")
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// ========================================
// UpdateLocation
// ========================================

func TestUpdateLocation_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().UpdateLocation(gomock.Any(), existing.DriverID, -6.2, 106.8).Return(nil)

	resp, err := svc.UpdateLocation(context.Background(), existing.DriverID, -6.2, 106.8)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Status != "success" {
		t.Errorf("expected success, got: %s", resp.Status)
	}
}

func TestUpdateLocation_InvalidLatitude(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateLocation(context.Background(), "d1", 91.0, 106.8)
	if err != service.ErrLocationInvalid {
		t.Fatalf("expected ErrLocationInvalid, got: %v", err)
	}
}

func TestUpdateLocation_InvalidLongitude(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateLocation(context.Background(), "d1", -6.2, 181.0)
	if err != service.ErrLocationInvalid {
		t.Fatalf("expected ErrLocationInvalid, got: %v", err)
	}
}

func TestUpdateLocation_NaN(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateLocation(context.Background(), "d1", math.NaN(), 106.8)
	if err != service.ErrLocationInvalid {
		t.Fatalf("expected ErrLocationInvalid, got: %v", err)
	}
}

func TestUpdateLocation_DriverNotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().FindByID(gomock.Any(), "unknown").Return(nil, repository.ErrDriverNotFound)

	_, err := svc.UpdateLocation(context.Background(), "unknown", -6.2, 106.8)
	if err != service.ErrDriverNotFound {
		t.Fatalf("expected ErrDriverNotFound, got: %v", err)
	}
}

func TestUpdateLocation_EmptyDriverID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateLocation(context.Background(), "", -6.2, 106.8)
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

func TestUpdateLocation_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	existing := sampleDriver()
	repoErr := errors.New("location update failed")
	mockRepo.EXPECT().FindByID(gomock.Any(), existing.DriverID).Return(existing, nil)
	mockRepo.EXPECT().UpdateLocation(gomock.Any(), existing.DriverID, -6.2, 106.8).Return(repoErr)

	_, err := svc.UpdateLocation(context.Background(), existing.DriverID, -6.2, 106.8)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// ========================================
// Model Validation
// ========================================

func TestDriverStatus_IsValid(t *testing.T) {
	tests := []struct {
		status model.DriverStatus
		valid  bool
	}{
		{model.DriverStatusOnline, true},
		{model.DriverStatusOffline, true},
		{model.DriverStatusBusy, true},
		{model.DriverStatus("INVALID"), false},
		{model.DriverStatus(""), false},
	}
	for _, tc := range tests {
		if tc.status.IsValid() != tc.valid {
			t.Errorf("status %q: expected IsValid=%v, got %v", tc.status, tc.valid, !tc.valid)
		}
	}
}
