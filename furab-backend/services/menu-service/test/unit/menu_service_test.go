package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"furab-backend/services/menu-service/internal/model"
	"furab-backend/services/menu-service/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ==========================================
// Mock Menu Repository (Manual gomock merujuk ke repository asli di package internal)
// ==========================================

type MockMenuRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMenuRepositoryMockRecorder
}

type MockMenuRepositoryMockRecorder struct {
	mock *MockMenuRepository
}

func NewMockMenuRepository(ctrl *gomock.Controller) *MockMenuRepository {
	mock := &MockMenuRepository{ctrl: ctrl}
	mock.recorder = &MockMenuRepositoryMockRecorder{mock}
	return mock
}

func (m *MockMenuRepository) EXPECT() *MockMenuRepositoryMockRecorder {
	return m.recorder
}

// Create mock
func (m *MockMenuRepository) Create(ctx context.Context, menu model.Menu) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, menu)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMenuRepositoryMockRecorder) Create(ctx, menu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMenuRepository)(nil).Create), ctx, menu)
}

// Update mock
func (m *MockMenuRepository) Update(ctx context.Context, menu model.Menu) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, menu)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMenuRepositoryMockRecorder) Update(ctx, menu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMenuRepository)(nil).Update), ctx, menu)
}

// Delete mock
func (m *MockMenuRepository) Delete(ctx context.Context, menuID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, menuID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMenuRepositoryMockRecorder) Delete(ctx, menuID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMenuRepository)(nil).Delete), ctx, menuID)
}

// UpdateStock mock
func (m *MockMenuRepository) UpdateStock(ctx context.Context, menuID string, jumlah int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStock", ctx, menuID, jumlah)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMenuRepositoryMockRecorder) UpdateStock(ctx, menuID, jumlah interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStock", reflect.TypeOf((*MockMenuRepository)(nil).UpdateStock), ctx, menuID, jumlah)
}

// GetByID mock
func (m *MockMenuRepository) GetByID(ctx context.Context, menuID string) (model.Menu, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, menuID)
	ret0, _ := ret[0].(model.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMenuRepositoryMockRecorder) GetByID(ctx, menuID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMenuRepository)(nil).GetByID), ctx, menuID)
}

// ListByMerchant mock
func (m *MockMenuRepository) ListByMerchant(ctx context.Context, merchantID string) ([]model.Menu, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByMerchant", ctx, merchantID)
	ret0, _ := ret[0].([]model.Menu)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMenuRepositoryMockRecorder) ListByMerchant(ctx, merchantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByMerchant", reflect.TypeOf((*MockMenuRepository)(nil).ListByMerchant), ctx, merchantID)
}

// SetAvailability mock
func (m *MockMenuRepository) SetAvailability(ctx context.Context, menuID string, status bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvailability", ctx, menuID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMenuRepositoryMockRecorder) SetAvailability(ctx, menuID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvailability", reflect.TypeOf((*MockMenuRepository)(nil).SetAvailability), ctx, menuID, status)
}

// ==========================================
// Unit Tests (Table-Driven)
// ==========================================

func TestMenuService_Create(t *testing.T) {
	ctx := context.Background()

	validMenu := model.Menu{
		MenuID:     "menu-123",
		MerchantID: "merch-1",
		NamaMenu:   "Nasi Goreng",
		Harga:      25000,
		Stok:       10,
	}

	tests := []struct {
		name          string
		input         model.Menu
		mockSetup     func(mockRepo *MockMenuRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Berhasil tambah menu",
			input: validMenu,
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().Create(ctx, validMenu).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name: "Negative Case: Harga negatif",
			input: model.Menu{
				MenuID:   "menu-123",
				NamaMenu: "Nasi Goreng",
				Harga:    -5000, // invalid
			},
			mockSetup: func(mockRepo *MockMenuRepository) {
				// repo not called due to validation
			},
			expectedError: "harga tidak boleh negatif",
		},
		{
			name: "Negative Case: Nama menu kosong",
			input: model.Menu{
				MenuID:   "menu-123",
				NamaMenu: "", // invalid
				Harga:    25000,
			},
			mockSetup: func(mockRepo *MockMenuRepository) {
				// repo not called due to validation
			},
			expectedError: "nama menu tidak boleh kosong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockMenuRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			err := svc.Create(ctx, tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMenuService_UpdateStock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		menuID        string
		jumlah        int
		mockSetup     func(mockRepo *MockMenuRepository)
		expectedError string
	}{
		{
			name:   "Positive Case: Update stok biasa (stok bertambah)",
			menuID: "menu-1",
			jumlah: 5,
			mockSetup: func(mockRepo *MockMenuRepository) {
				// Current stock = 10
				mockRepo.EXPECT().GetByID(ctx, "menu-1").Return(model.Menu{MenuID: "menu-1", Stok: 10}, nil).Times(1)
				// Success update (newStock = 15)
				mockRepo.EXPECT().UpdateStock(ctx, "menu-1", 5).Return(nil).Times(1)
				// No SetAvailability called since stock is > 0 and was > 0
			},
			expectedError: "",
		},
		{
			name:   "Positive Case: Stok menjadi 0, auto SetAvailability(false)",
			menuID: "menu-2",
			jumlah: -10,
			mockSetup: func(mockRepo *MockMenuRepository) {
				// Current stock = 10
				mockRepo.EXPECT().GetByID(ctx, "menu-2").Return(model.Menu{MenuID: "menu-2", Stok: 10}, nil).Times(1)
				// Success update (newStock = 0)
				mockRepo.EXPECT().UpdateStock(ctx, "menu-2", -10).Return(nil).Times(1)
				// Auto SetAvailability
				mockRepo.EXPECT().SetAvailability(ctx, "menu-2", false).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:   "Positive Case: Stok dari 0 ditambah, auto SetAvailability(true)",
			menuID: "menu-3",
			jumlah: 5,
			mockSetup: func(mockRepo *MockMenuRepository) {
				// Current stock = 0
				mockRepo.EXPECT().GetByID(ctx, "menu-3").Return(model.Menu{MenuID: "menu-3", Stok: 0}, nil).Times(1)
				// Success update (newStock = 5)
				mockRepo.EXPECT().UpdateStock(ctx, "menu-3", 5).Return(nil).Times(1)
				// Auto SetAvailability
				mockRepo.EXPECT().SetAvailability(ctx, "menu-3", true).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:   "Negative Case: Stok tidak cukup",
			menuID: "menu-4",
			jumlah: -15, // Want to deduct 15, but only has 10
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().GetByID(ctx, "menu-4").Return(model.Menu{MenuID: "menu-4", Stok: 10}, nil).Times(1)
			},
			expectedError: "insufficient stock",
		},
		{
			name:   "Negative Case: Error dari repository GetByID",
			menuID: "menu-5",
			jumlah: -5,
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().GetByID(ctx, "menu-5").Return(model.Menu{}, errors.New("db connection failed")).Times(1)
			},
			expectedError: "db connection failed",
		},
		{
			name:   "Negative Case: Error dari repository UpdateStock",
			menuID: "menu-6",
			jumlah: -5,
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().GetByID(ctx, "menu-6").Return(model.Menu{MenuID: "menu-6", Stok: 10}, nil).Times(1)
				mockRepo.EXPECT().UpdateStock(ctx, "menu-6", -5).Return(errors.New("db update failed")).Times(1)
			},
			expectedError: "db update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockMenuRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			err := svc.UpdateStock(ctx, tt.menuID, tt.jumlah)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMenuService_GetByID(t *testing.T) {
	ctx := context.Background()

	expectedMenu := model.Menu{
		MenuID:      "menu-123",
		NamaMenu:    "Nasi Goreng",
		IsAvailable: true,
	}

	tests := []struct {
		name          string
		menuID        string
		mockSetup     func(mockRepo *MockMenuRepository)
		expectedError string
		expectedRes   model.Menu
	}{
		{
			name:   "Positive Case: Berhasil ambil menu",
			menuID: "menu-123",
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().GetByID(ctx, "menu-123").Return(expectedMenu, nil).Times(1)
			},
			expectedError: "",
			expectedRes:   expectedMenu,
		},
		{
			name:   "Negative Case: Menu tidak ditemukan",
			menuID: "invalid-id",
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().GetByID(ctx, "invalid-id").Return(model.Menu{}, errors.New("menu not found")).Times(1)
			},
			expectedError: "menu not found",
			expectedRes:   model.Menu{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockMenuRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			res, err := svc.GetByID(ctx, tt.menuID)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMenuService_ListByMerchant(t *testing.T) {
	ctx := context.Background()

	menuList := []model.Menu{
		{MenuID: "m1", NamaMenu: "A"},
		{MenuID: "m2", NamaMenu: "B"},
	}

	tests := []struct {
		name          string
		merchantID    string
		mockSetup     func(mockRepo *MockMenuRepository)
		expectedError string
		expectedRes   []model.Menu
	}{
		{
			name:       "Positive Case: List menu merchant",
			merchantID: "merch-1",
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().ListByMerchant(ctx, "merch-1").Return(menuList, nil).Times(1)
			},
			expectedError: "",
			expectedRes:   menuList,
		},
		{
			name:       "Positive Case: List kosong",
			merchantID: "merch-2",
			mockSetup: func(mockRepo *MockMenuRepository) {
				mockRepo.EXPECT().ListByMerchant(ctx, "merch-2").Return([]model.Menu{}, nil).Times(1)
			},
			expectedError: "",
			expectedRes:   []model.Menu{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockMenuRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewMenuService(mockRepo)
			res, err := svc.ListByMerchant(ctx, tt.merchantID)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
