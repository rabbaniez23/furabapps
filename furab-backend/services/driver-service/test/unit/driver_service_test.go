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

var (
	errRepoSave           = errors.New("repo save error")
	errRepoFindByID       = errors.New("repo find error")
	errRepoUpdate         = errors.New("repo update error")
	errRepoUpdateStatus   = errors.New("repo update status error")
	errRepoUpdateLocation = errors.New("repo update location error")
)

type driverArgMatcher struct {
	match func(*model.Driver) bool
}

func (m driverArgMatcher) Matches(x any) bool {
	d, ok := x.(*model.Driver)
	if !ok {
		return false
	}
	return m.match(d)
}

func (m driverArgMatcher) String() string {
	return "matches *model.Driver predicate"
}

func matchDriver(match func(*model.Driver) bool) gomock.Matcher {
	return driverArgMatcher{match: match}
}

func setupDriverService(t *testing.T) (*gomock.Controller, *mock_repository.MockDriverRepository, service.DriverService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockDriverRepository(ctrl)
	return ctrl, mockRepo, service.NewDriverService(mockRepo)
}

func TestDriverService_CreateDriver(t *testing.T) {
	ctrl, mockRepo, svc := setupDriverService(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "drv-1",
			Name:        "Driver A",
			Phone:       "0812345678",
			VehicleType: "motor",
		}

		mockRepo.EXPECT().
			Save(gomock.Any(), matchDriver(func(d *model.Driver) bool {
				return d != nil &&
					d.DriverID == "drv-1" &&
					d.Name == "Driver A" &&
					d.Phone == "0812345678" &&
					d.VehicleType == "motor" &&
					d.Status == "OFFLINE"
			})).
			Return(nil)

		res, err := svc.CreateDriver(ctx, req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
		if res.DriverID != "drv-1" {
			t.Errorf("Expected DriverID 'drv-1', got %v", res.DriverID)
		}
	})

	t.Run("validation_error_missing_driver_id", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "",
			Name:        "Driver A",
			Phone:       "0812345678",
			VehicleType: "motor",
		}

		res, err := svc.CreateDriver(ctx, req)

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverIDRequired) {
			t.Errorf("Expected ErrDriverIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_missing_name", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "drv-1",
			Name:        "",
			Phone:       "0812345678",
			VehicleType: "motor",
		}

		res, err := svc.CreateDriver(ctx, req)

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrNameRequired) {
			t.Errorf("Expected ErrNameRequired, got %v", err)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		req := model.CreateDriverRequest{
			DriverID:    "drv-1",
			Name:        "Driver A",
			Phone:       "0812345678",
			VehicleType: "motor",
		}

		mockRepo.EXPECT().
			Save(gomock.Any(), matchDriver(func(d *model.Driver) bool {
				return d != nil &&
					d.DriverID == "drv-1" &&
					d.Name == "Driver A" &&
					d.Phone == "0812345678" &&
					d.VehicleType == "motor" &&
					d.Status == "OFFLINE"
			})).
			Return(errRepoSave)

		res, err := svc.CreateDriver(ctx, req)

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, errRepoSave) {
			t.Errorf("Expected repo save error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		req := model.CreateDriverRequest{
			DriverID:    "drv-1",
			Name:        "Driver A",
			Phone:       "0812345678",
			VehicleType: "motor",
		}
		res, err := svc.CreateDriver(cancelledCtx, req)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestDriverService_GetDriver(t *testing.T) {
	ctrl, mockRepo, svc := setupDriverService(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedDriver := &model.Driver{DriverID: "drv-1", Name: "Driver A"}

		mockRepo.EXPECT().
			FindByID(gomock.Any(), "drv-1").
			Return(expectedDriver, nil)

		driver, err := svc.GetDriver(ctx, "drv-1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if driver == nil || driver.DriverID != "drv-1" {
			t.Errorf("Expected driver with ID 'drv-1', got %+v", driver)
		}
	})

	t.Run("validation_error_missing_driver_id", func(t *testing.T) {
		driver, err := svc.GetDriver(ctx, " ")
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
		if !errors.Is(err, service.ErrDriverIDRequired) {
			t.Errorf("Expected ErrDriverIDRequired, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(gomock.Any(), "drv-99").
			Return(nil, nil)

		driver, err := svc.GetDriver(ctx, "drv-99")
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
		if !errors.Is(err, service.ErrDriverNotFound) {
			t.Errorf("Expected ErrDriverNotFound, got %v", err)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(gomock.Any(), "drv-1").
			Return(nil, errRepoFindByID)

		driver, err := svc.GetDriver(ctx, "drv-1")

		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
		if !errors.Is(err, errRepoFindByID) {
			t.Errorf("Expected repo find error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		driver, err := svc.GetDriver(cancelledCtx, "drv-1")
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestDriverService_UpdateDriver(t *testing.T) {
	ctrl, mockRepo, svc := setupDriverService(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1", Name: "Old Name"}
		req := model.UpdateDriverRequest{
			Name:        "Driver Update",
			Phone:       "0812345678",
			VehicleType: "mobil",
		}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().
			Update(gomock.Any(), matchDriver(func(d *model.Driver) bool {
				return d != nil &&
					d.DriverID == "drv-1" &&
					d.Name == "Driver Update" &&
					d.Phone == "0812345678" &&
					d.VehicleType == "mobil"
			})).
			Return(nil)

		res, err := svc.UpdateDriver(ctx, "drv-1", req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
	})

	t.Run("validation_error_missing_driver_id", func(t *testing.T) {
		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "0812345678", VehicleType: "mobil"}
		res, err := svc.UpdateDriver(ctx, "", req)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverIDRequired) {
			t.Errorf("Expected ErrDriverIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_empty_update_payload", func(t *testing.T) {
		req := model.UpdateDriverRequest{Name: "", Phone: "0812345678", VehicleType: "mobil"}
		res, err := svc.UpdateDriver(ctx, "drv-1", req)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrValidation) {
			t.Errorf("Expected ErrValidation, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "0812345678", VehicleType: "mobil"}
		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-99").Return(nil, nil)
		res, err := svc.UpdateDriver(ctx, "drv-99", req)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverNotFound) {
			t.Errorf("Expected ErrDriverNotFound, got %v", err)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1", Name: "Old Name"}
		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "0812345678", VehicleType: "mobil"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().
			Update(gomock.Any(), matchDriver(func(d *model.Driver) bool {
				return d != nil &&
					d.DriverID == "drv-1" &&
					d.Name == "Driver Update" &&
					d.Phone == "0812345678" &&
					d.VehicleType == "mobil"
			})).
			Return(errRepoUpdate)

		res, err := svc.UpdateDriver(ctx, "drv-1", req)

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, errRepoUpdate) {
			t.Errorf("Expected repo update error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		req := model.UpdateDriverRequest{Name: "Driver Update", Phone: "0812345678", VehicleType: "mobil"}
		res, err := svc.UpdateDriver(cancelledCtx, "drv-1", req)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestDriverService_UpdateStatus(t *testing.T) {
	ctrl, mockRepo, svc := setupDriverService(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1", Status: "OFFLINE"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "drv-1", "ACTIVE").Return(nil)

		res, err := svc.UpdateStatus(ctx, "drv-1", "active")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
	})

	t.Run("validation_error_missing_driver_id", func(t *testing.T) {
		res, err := svc.UpdateStatus(ctx, "", "ACTIVE")
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverIDRequired) {
			t.Errorf("Expected ErrDriverIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_invalid_status", func(t *testing.T) {
		res, err := svc.UpdateStatus(ctx, "drv-1", "READY")
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrStatusInvalid) {
			t.Errorf("Expected ErrStatusInvalid, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-99").Return(nil, nil)
		res, err := svc.UpdateStatus(ctx, "drv-99", "ACTIVE")
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverNotFound) {
			t.Errorf("Expected ErrDriverNotFound, got %v", err)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1", Status: "OFFLINE"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), "drv-1", "ACTIVE").Return(errRepoUpdateStatus)

		res, err := svc.UpdateStatus(ctx, "drv-1", "ACTIVE")

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, errRepoUpdateStatus) {
			t.Errorf("Expected repo update status error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.UpdateStatus(cancelledCtx, "drv-1", "ACTIVE")
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}

func TestDriverService_UpdateLocation(t *testing.T) {
	ctrl, mockRepo, svc := setupDriverService(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateLocation(gomock.Any(), "drv-1", -6.2, 106.8).Return(nil)

		res, err := svc.UpdateLocation(ctx, "drv-1", -6.2, 106.8)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" {
			t.Errorf("Expected status success, got %v", res.Status)
		}
	})

	t.Run("validation_error_missing_driver_id", func(t *testing.T) {
		res, err := svc.UpdateLocation(ctx, "", -6.2, 106.8)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverIDRequired) {
			t.Errorf("Expected ErrDriverIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_invalid_location", func(t *testing.T) {
		res, err := svc.UpdateLocation(ctx, "drv-1", -91, 106.8)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrLocationInvalid) {
			t.Errorf("Expected ErrLocationInvalid, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-99").Return(nil, nil)
		res, err := svc.UpdateLocation(ctx, "drv-99", -6.2, 106.8)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, service.ErrDriverNotFound) {
			t.Errorf("Expected ErrDriverNotFound, got %v", err)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		existingDriver := &model.Driver{DriverID: "drv-1"}

		mockRepo.EXPECT().FindByID(gomock.Any(), "drv-1").Return(existingDriver, nil)
		mockRepo.EXPECT().UpdateLocation(gomock.Any(), "drv-1", -6.2, 106.8).Return(errRepoUpdateLocation)

		res, err := svc.UpdateLocation(ctx, "drv-1", -6.2, 106.8)

		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, errRepoUpdateLocation) {
			t.Errorf("Expected repo update location error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.UpdateLocation(cancelledCtx, "drv-1", -6.2, 106.8)
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})
}
