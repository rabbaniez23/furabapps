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
// Mock akan digenerate otomatis menggunakan gomock (mockgen).
// ============================================================================

// 1. Models & DTOs

type Driver struct {
	DriverID    string
	Name        string
	Phone       string
	VehicleType string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateDriverRequest struct {
	DriverID    string
	Name        string
	Phone       string
	VehicleType string
}

type UpdateDriverRequest struct {
	Name        string
	Phone       string
	VehicleType string
}

type DriverResponse struct {
	Status  string
	Message string
}

// 2. Mocked Repository Interface

type MockDriverRepository struct {
	SaveFunc           func(ctx context.Context, driver *Driver) error
	FindByIDFunc       func(ctx context.Context, driverID string) (*Driver, error)
	UpdateFunc         func(ctx context.Context, driver *Driver) error
	UpdateStatusFunc   func(ctx context.Context, driverID, status string) error
	UpdateLocationFunc func(ctx context.Context, driverID string, lat, long float64) error
}

func (m *MockDriverRepository) Save(ctx context.Context, driver *Driver) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, driver)
	}
	return nil
}

func (m *MockDriverRepository) FindByID(ctx context.Context, driverID string) (*Driver, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, driverID)
	}
	return nil, nil
}

func (m *MockDriverRepository) Update(ctx context.Context, driver *Driver) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, driver)
	}
	return nil
}

func (m *MockDriverRepository) UpdateStatus(ctx context.Context, driverID, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, driverID, status)
	}
	return nil
}

func (m *MockDriverRepository) UpdateLocation(ctx context.Context, driverID string, lat, long float64) error {
	if m.UpdateLocationFunc != nil {
		return m.UpdateLocationFunc(ctx, driverID, lat, long)
	}
	return nil
}

// 3. DriverService Interface & Implementation

type DriverService interface {
	CreateDriver(ctx context.Context, req CreateDriverRequest) (*DriverResponse, error)
	GetDriver(ctx context.Context, driverID string) (*Driver, error)
	UpdateDriver(ctx context.Context, driverID string, req UpdateDriverRequest) (*DriverResponse, error)
	UpdateStatus(ctx context.Context, driverID, status string) (*DriverResponse, error)
	UpdateLocation(ctx context.Context, driverID string, lat, long float64) (*DriverResponse, error)
}

type driverServiceImpl struct {
	repo *MockDriverRepository
}

func (s *driverServiceImpl) CreateDriver(ctx context.Context, req CreateDriverRequest) (*DriverResponse, error) {
	if req.Name == "" || req.Phone == "" || req.VehicleType == "" {
		return &DriverResponse{Status: "failed", Message: "validation error"}, errors.New("validation error")
	}

	driver := &Driver{
		DriverID:    req.DriverID,
		Name:        req.Name,
		Phone:       req.Phone,
		VehicleType: req.VehicleType,
		Status:      "offline",
	}

	if err := s.repo.Save(ctx, driver); err != nil {
		return &DriverResponse{Status: "failed", Message: "db error"}, err
	}

	return &DriverResponse{Status: "success", Message: "driver berhasil dibuat"}, nil
}

func (s *driverServiceImpl) GetDriver(ctx context.Context, driverID string) (*Driver, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}
	return driver, nil
}

func (s *driverServiceImpl) UpdateDriver(ctx context.Context, driverID string, req UpdateDriverRequest) (*DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	driver.Name = req.Name
	driver.Phone = req.Phone
	driver.VehicleType = req.VehicleType

	if err := s.repo.Update(ctx, driver); err != nil {
		return nil, err
	}

	return &DriverResponse{Status: "success", Message: "update berhasil"}, nil
}

func (s *driverServiceImpl) UpdateStatus(ctx context.Context, driverID, status string) (*DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	if err := s.repo.UpdateStatus(ctx, driverID, status); err != nil {
		return nil, err
	}

	return &DriverResponse{Status: "success", Message: "status berhasil diperbarui"}, nil
}

func (s *driverServiceImpl) UpdateLocation(ctx context.Context, driverID string, lat, long float64) (*DriverResponse, error) {
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver == nil {
		return nil, errors.New("driver not found")
	}

	if err := s.repo.UpdateLocation(ctx, driverID, lat, long); err != nil {
		return nil, err
	}

	return &DriverResponse{Status: "success", Message: "lokasi berhasil diperbarui"}, nil
}

// ============================================================================
// UNIT TESTS MULAI DARI SINI
// ============================================================================

func TestDriverService_CreateDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Driver berhasil dibuat", func(t *testing.T) {
		repo := &MockDriverRepository{
			SaveFunc: func(ctx context.Context, driver *Driver) error {
				if driver.DriverID != "1" || driver.Name != "Driver A" || driver.Phone != "08123" || driver.VehicleType != "motor" {
					t.Errorf("Unexpected driver data saved: %+v", driver)
				}
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		req := CreateDriverRequest{
			DriverID:    "1",
			Name:        "Driver A",
			Phone:       "08123",
			VehicleType: "motor",
		}

		res, err := service.CreateDriver(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "driver berhasil dibuat" {
			t.Errorf("Expected success response, got %+v", res)
		}
	})

	t.Run("Error - Data tidak lengkap", func(t *testing.T) {
		repo := &MockDriverRepository{
			SaveFunc: func(ctx context.Context, driver *Driver) error {
				t.Fatal("Repository Save() should not be called")
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		req := CreateDriverRequest{
			DriverID:    "1",
			Name:        "", // Kosong memicu error
			Phone:       "08123",
			VehicleType: "motor",
		}

		res, err := service.CreateDriver(context.Background(), req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if res == nil || res.Status != "failed" || res.Message != "validation error" {
			t.Errorf("Expected failed response, got %+v", res)
		}
	})
}

func TestDriverService_GetDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Driver ditemukan", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				if driverID != "1" {
					t.Errorf("Expected DriverID '1', got %s", driverID)
				}
				return &Driver{DriverID: "1", Name: "Driver A"}, nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		driver, err := service.GetDriver(context.Background(), "1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if driver == nil || driver.DriverID != "1" {
			t.Errorf("Expected Driver 1, got %+v", driver)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return nil, errors.New("driver not found")
			},
		}
		service := &driverServiceImpl{repo: repo}

		driver, err := service.GetDriver(context.Background(), "99")

		if err == nil || err.Error() != "driver not found" {
			t.Fatalf("Expected 'driver not found' error, got %v", err)
		}
		if driver != nil {
			t.Errorf("Expected nil driver, got %+v", driver)
		}
	})
}

func TestDriverService_UpdateDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Data berhasil diupdate", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return &Driver{DriverID: "1", Name: "Old Name"}, nil
			},
			UpdateFunc: func(ctx context.Context, driver *Driver) error {
				if driver.Name != "Driver Update" || driver.Phone != "08124" || driver.VehicleType != "mobil" {
					t.Errorf("Unexpected updated driver data: %+v", driver)
				}
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		req := UpdateDriverRequest{
			Name:        "Driver Update",
			Phone:       "08124",
			VehicleType: "mobil",
		}

		res, err := service.UpdateDriver(context.Background(), "1", req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "update berhasil" {
			t.Errorf("Expected success response, got %+v", res)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return nil, errors.New("driver not found")
			},
			UpdateFunc: func(ctx context.Context, driver *Driver) error {
				t.Fatal("Repository Update() should not be called")
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		req := UpdateDriverRequest{
			Name:        "Driver Update",
			Phone:       "08124",
			VehicleType: "mobil",
		}

		res, err := service.UpdateDriver(context.Background(), "99", req)

		if err == nil || err.Error() != "driver not found" {
			t.Fatalf("Expected 'driver not found' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %+v", res)
		}
	})
}

func TestDriverService_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Status diupdate", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return &Driver{DriverID: "1", Status: "offline"}, nil
			},
			UpdateStatusFunc: func(ctx context.Context, driverID, status string) error {
				if driverID != "1" || status != "online" {
					t.Errorf("Unexpected driverID %s or status %s", driverID, status)
				}
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		res, err := service.UpdateStatus(context.Background(), "1", "online")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "status berhasil diperbarui" {
			t.Errorf("Expected success response, got %+v", res)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return nil, errors.New("driver not found")
			},
		}
		service := &driverServiceImpl{repo: repo}

		res, err := service.UpdateStatus(context.Background(), "99", "online")

		if err == nil || err.Error() != "driver not found" {
			t.Fatalf("Expected 'driver not found' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %+v", res)
		}
	})
}

func TestDriverService_UpdateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Lokasi diupdate", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return &Driver{DriverID: "1"}, nil
			},
			UpdateLocationFunc: func(ctx context.Context, driverID string, lat, long float64) error {
				if driverID != "1" || lat != -6.2 || long != 106.8 {
					t.Errorf("Unexpected location update: id=%s, lat=%f, long=%f", driverID, lat, long)
				}
				return nil
			},
		}
		service := &driverServiceImpl{repo: repo}

		res, err := service.UpdateLocation(context.Background(), "1", -6.2, 106.8)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.Status != "success" || res.Message != "lokasi berhasil diperbarui" {
			t.Errorf("Expected success response, got %+v", res)
		}
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		repo := &MockDriverRepository{
			FindByIDFunc: func(ctx context.Context, driverID string) (*Driver, error) {
				return nil, errors.New("driver not found")
			},
		}
		service := &driverServiceImpl{repo: repo}

		res, err := service.UpdateLocation(context.Background(), "99", -6.2, 106.8)

		if err == nil || err.Error() != "driver not found" {
			t.Fatalf("Expected 'driver not found' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %+v", res)
		}
	})
}
