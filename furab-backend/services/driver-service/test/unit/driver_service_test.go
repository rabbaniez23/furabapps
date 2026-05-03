package unit

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"furab-backend/services/driver-service/internal/model"
	"furab-backend/services/driver-service/internal/service"
	mock_repository "furab-backend/services/driver-service/internal/repository/mock"
)

// ============================================================================
// TestDriverService_CreateDriver
// ============================================================================

func TestDriverService_CreateDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)

	ctx := context.Background()

	t.Run("Success - Driver berhasil dibuat", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "1",
			Name:        "Driver A",
			Phone:       "08123",
			VehicleType: "motor",
		}

		mockRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		res, err := svc.CreateDriver(ctx, req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.Message != "driver berhasil dibuat" {
			t.Errorf("Expected message 'driver berhasil dibuat', got %v", res.Message)
		}
		if res.DriverID != "1" {
			t.Errorf("Expected DriverID '1', got %v", res.DriverID)
		}
	})

	t.Run("Error - Input kosong (Name)", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "1",
			Name:        "", // kosong
			Phone:       "08123",
			VehicleType: "motor",
		}

		res, err := svc.CreateDriver(ctx, req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "validation error" {
			t.Errorf("Expected 'validation error', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Input kosong (Phone)", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "1",
			Name:        "Driver A",
			Phone:       "", // kosong
			VehicleType: "motor",
		}

		res, err := svc.CreateDriver(ctx, req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "validation error" {
			t.Errorf("Expected 'validation error', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Repository gagal", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "1",
			Name:        "Driver A",
			Phone:       "08123",
			VehicleType: "motor",
		}
		expectedErr := errors.New("db error")

		mockRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(expectedErr)

		res, err := svc.CreateDriver(ctx, req)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestDriverService_GetDriver
// ============================================================================

func TestDriverService_GetDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)

	ctx := context.Background()

	t.Run("Success - Driver ditemukan", func(t *testing.T) {
		expectedDriver := &model.Driver{DriverID: "1", Name: "Driver A"}

		mockRepo.EXPECT().
			FindByID(gomock.Any(), "1").
			Return(expectedDriver, nil)

		driver, err := svc.GetDriver(ctx, "1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if driver == nil || driver.DriverID != "1" {
			t.Errorf("Expected driver with ID '1', got %+v", driver)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(gomock.Any(), "99").
			Return(nil, nil)

		driver, err := svc.GetDriver(ctx, "99")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "driver not found" {
			t.Errorf("Expected 'driver not found', got %v", err.Error())
		}
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
	})

	t.Run("Error - Repository error", func(t *testing.T) {
		expectedErr := errors.New("db error")

		mockRepo.EXPECT().
			FindByID(gomock.Any(), "1").
			Return(nil, expectedErr)

		driver, err := svc.GetDriver(ctx, "1")

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
	})
}

// ============================================================================
// TestDriverService_UpdateDriver
// ============================================================================

func TestDriverService_UpdateDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)

	ctx := context.Background()

	t.Run("Success - Data berhasil diupdate", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1", Name: "Old Name"}
		req := model.UpdateDriverRequest{
			Name:        "Driver Update",
			Phone:       "08124",
			VehicleType: "mobil",
		}

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		res, err := svc.UpdateDriver(ctx, "1", req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.Message != "update berhasil" {
			t.Errorf("Expected message 'update berhasil', got %v", res.Message)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "08124", VehicleType: "mobil"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		res, err := svc.UpdateDriver(ctx, "99", req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "driver not found" {
			t.Errorf("Expected 'driver not found', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Update gagal", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1", Name: "Old Name"}
		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "08124", VehicleType: "mobil"}
		expectedErr := errors.New("db error")

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedErr)

		res, err := svc.UpdateDriver(ctx, "1", req)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestDriverService_UpdateStatus
// ============================================================================

func TestDriverService_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)

	ctx := context.Background()

	t.Run("Success - Status berhasil diperbarui", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1", Status: "offline"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "1", "online").Return(nil)

		res, err := svc.UpdateStatus(ctx, "1", "online")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.Message != "status berhasil diperbarui" {
			t.Errorf("Expected message 'status berhasil diperbarui', got %v", res.Message)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		res, err := svc.UpdateStatus(ctx, "99", "online")

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "driver not found" {
			t.Errorf("Expected 'driver not found', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Repository gagal", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1", Status: "offline"}
		expectedErr := errors.New("db error")

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "1", "online").Return(expectedErr)

		res, err := svc.UpdateStatus(ctx, "1", "online")

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

// ============================================================================
// TestDriverService_UpdateLocation
// ============================================================================

func TestDriverService_UpdateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	svc := service.NewDriverService(mockRepo)

	ctx := context.Background()

	t.Run("Success - Lokasi berhasil diperbarui", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateLocation(gomock.Any(), "1", -6.2, 106.8).Return(nil)

		res, err := svc.UpdateLocation(ctx, "1", -6.2, 106.8)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.Message != "lokasi berhasil diperbarui" {
			t.Errorf("Expected message 'lokasi berhasil diperbarui', got %v", res.Message)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		res, err := svc.UpdateLocation(ctx, "99", -6.2, 106.8)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != "driver not found" {
			t.Errorf("Expected 'driver not found', got %v", err.Error())
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Repository gagal", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "1"}
		expectedErr := errors.New("db error")

		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateLocation(gomock.Any(), "1", -6.2, 106.8).Return(expectedErr)

		res, err := svc.UpdateLocation(ctx, "1", -6.2, 106.8)

		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}
