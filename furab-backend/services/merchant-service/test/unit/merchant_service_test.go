package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"furab-backend/services/merchant-service/internal/model"
	"furab-backend/services/merchant-service/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ==========================================
// Mock Merchant Repository
// ==========================================

type MockMerchantRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMerchantRepositoryMockRecorder
}

type MockMerchantRepositoryMockRecorder struct {
	mock *MockMerchantRepository
}

func NewMockMerchantRepository(ctrl *gomock.Controller) *MockMerchantRepository {
	mock := &MockMerchantRepository{ctrl: ctrl}
	mock.recorder = &MockMerchantRepositoryMockRecorder{mock}
	return mock
}

func (m *MockMerchantRepository) EXPECT() *MockMerchantRepositoryMockRecorder {
	return m.recorder
}

// Create
func (m *MockMerchantRepository) Create(ctx context.Context, merchant model.Merchant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, merchant)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMerchantRepositoryMockRecorder) Create(ctx, merchant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMerchantRepository)(nil).Create), ctx, merchant)
}

// Update
func (m *MockMerchantRepository) Update(ctx context.Context, merchant model.Merchant) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, merchant)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMerchantRepositoryMockRecorder) Update(ctx, merchant interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMerchantRepository)(nil).Update), ctx, merchant)
}

// GetByID
func (m *MockMerchantRepository) GetByID(ctx context.Context, merchantID string) (model.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, merchantID)
	ret0, _ := ret[0].(model.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMerchantRepositoryMockRecorder) GetByID(ctx, merchantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMerchantRepository)(nil).GetByID), ctx, merchantID)
}

// UpdateStatus
func (m *MockMerchantRepository) UpdateStatus(ctx context.Context, merchantID string, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, merchantID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMerchantRepositoryMockRecorder) UpdateStatus(ctx, merchantID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockMerchantRepository)(nil).UpdateStatus), ctx, merchantID, status)
}

// Deactivate
func (m *MockMerchantRepository) Deactivate(ctx context.Context, merchantID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deactivate", ctx, merchantID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMerchantRepositoryMockRecorder) Deactivate(ctx, merchantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deactivate", reflect.TypeOf((*MockMerchantRepository)(nil).Deactivate), ctx, merchantID)
}

// Search
func (m *MockMerchantRepository) Search(ctx context.Context, filter map[string]interface{}) ([]model.Merchant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, filter)
	ret0, _ := ret[0].([]model.Merchant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMerchantRepositoryMockRecorder) Search(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockMerchantRepository)(nil).Search), ctx, filter)
}

// SetStatusCache
func (m *MockMerchantRepository) SetStatusCache(ctx context.Context, merchantID string, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatusCache", ctx, merchantID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockMerchantRepositoryMockRecorder) SetStatusCache(ctx, merchantID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatusCache", reflect.TypeOf((*MockMerchantRepository)(nil).SetStatusCache), ctx, merchantID, status)
}

// GetStatusCache
func (m *MockMerchantRepository) GetStatusCache(ctx context.Context, merchantID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatusCache", ctx, merchantID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockMerchantRepositoryMockRecorder) GetStatusCache(ctx, merchantID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatusCache", reflect.TypeOf((*MockMerchantRepository)(nil).GetStatusCache), ctx, merchantID)
}

// ==========================================
// Unit Tests (Table-Driven)
// ==========================================

func TestMerchantService_Create(t *testing.T) {
	ctx := context.Background()
	validMerchant := model.Merchant{UserID: "u1", NamaToko: "Toko A", StatusOperasional: "closed", IsActive: true}

	tests := []struct {
		name          string
		input         model.Merchant
		mockSetup     func(mockRepo *MockMerchantRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Register berhasil",
			input: model.Merchant{UserID: "u1", NamaToko: "Toko A"},
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().Create(ctx, validMerchant).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name: "Negative Case: UserID kosong",
			input: model.Merchant{NamaToko: "Toko B"}, // UserID empty
			mockSetup: func(mockRepo *MockMerchantRepository) {},
			expectedError: "user_id tidak boleh kosong",
		},
		{
			name: "Negative Case: Nama Toko kosong",
			input: model.Merchant{UserID: "u2"}, // NamaToko empty
			mockSetup: func(mockRepo *MockMerchantRepository) {},
			expectedError: "nama toko tidak boleh kosong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockMerchantRepository(ctrl)
			tt.mockSetup(mockRepo)
			svc := service.NewMerchantService(mockRepo)
			err := svc.Create(ctx, tt.input)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMerchantService_UpdateStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		merchantID    string
		status        string
		mockSetup     func(mockRepo *MockMerchantRepository)
		expectedError string
	}{
		{
			name:       "Positive Case: Update status sukses ke open",
			merchantID: "m1",
			status:     "open",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				// Pastikan kedua repository method dipanggil berurutan/bersamaan
				mockRepo.EXPECT().UpdateStatus(ctx, "m1", "open").Return(nil).Times(1)
				mockRepo.EXPECT().SetStatusCache(ctx, "m1", "open").Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:       "Negative Case: Status invalid",
			merchantID: "m2",
			status:     "invalid_status",
			mockSetup: func(mockRepo *MockMerchantRepository) {},
			expectedError: "status tidak valid",
		},
		{
			name:       "Negative Case: Error DB pada UpdateStatus",
			merchantID: "m3",
			status:     "closed",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().UpdateStatus(ctx, "m3", "closed").Return(errors.New("db err")).Times(1)
			},
			expectedError: "db err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockMerchantRepository(ctrl)
			tt.mockSetup(mockRepo)
			svc := service.NewMerchantService(mockRepo)
			err := svc.UpdateStatus(ctx, tt.merchantID, tt.status)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMerchantService_Deactivate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		merchantID    string
		mockSetup     func(mockRepo *MockMerchantRepository)
		expectedError string
	}{
		{
			name:       "Positive Case: Deactivate berhasil",
			merchantID: "m1",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().Deactivate(ctx, "m1").Return(nil).Times(1)
				mockRepo.EXPECT().UpdateStatus(ctx, "m1", "closed").Return(nil).Times(1)
				mockRepo.EXPECT().SetStatusCache(ctx, "m1", "closed").Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:       "Negative Case: Error saat Set IsActive false di DB",
			merchantID: "m2",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().Deactivate(ctx, "m2").Return(errors.New("db err")).Times(1)
			},
			expectedError: "db err",
		},
		{
			name:       "Negative Case: Error saat ubah status menjadi closed",
			merchantID: "m3",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().Deactivate(ctx, "m3").Return(nil).Times(1)
				mockRepo.EXPECT().UpdateStatus(ctx, "m3", "closed").Return(errors.New("db update err")).Times(1)
			},
			expectedError: "db update err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockMerchantRepository(ctrl)
			tt.mockSetup(mockRepo)
			svc := service.NewMerchantService(mockRepo)
			err := svc.Deactivate(ctx, tt.merchantID)
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestMerchantService_Search(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		keyword     string
		kategori    string
		mockSetup   func(mockRepo *MockMerchantRepository)
		expectedRes []model.Merchant
	}{
		{
			name:     "Positive Case: Search dengan keyword dan kategori",
			keyword:  "nasi",
			kategori: "makanan",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				expectedFilter := map[string]interface{}{"keyword": "nasi", "kategori": "makanan"}
				mockRepo.EXPECT().Search(ctx, expectedFilter).Return([]model.Merchant{{NamaToko: "Nasi Uduk"}}, nil).Times(1)
			},
			expectedRes: []model.Merchant{{NamaToko: "Nasi Uduk"}},
		},
		{
			name:     "Positive Case: Search tanpa filter",
			keyword:  "",
			kategori: "",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				expectedFilter := map[string]interface{}{} // empty map
				mockRepo.EXPECT().Search(ctx, expectedFilter).Return([]model.Merchant{{NamaToko: "A"}, {NamaToko: "B"}}, nil).Times(1)
			},
			expectedRes: []model.Merchant{{NamaToko: "A"}, {NamaToko: "B"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockMerchantRepository(ctrl)
			tt.mockSetup(mockRepo)
			svc := service.NewMerchantService(mockRepo)
			res, err := svc.Search(ctx, tt.keyword, tt.kategori)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestMerchantService_CheckMerchantStatus(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		merchantID    string
		mockSetup     func(mockRepo *MockMerchantRepository)
		expectedRes   bool
		expectedError string
	}{
		{
			name:       "Positive Case: Dari cache langsung terbaca closed (bypass DB)",
			merchantID: "m1",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				// Cache hit returning 'closed', service bypasses DB!
				mockRepo.EXPECT().GetStatusCache(ctx, "m1").Return("closed", nil).Times(1)
			},
			expectedRes:   false,
			expectedError: "",
		},
		{
			name:       "Positive Case: Cache terbaca open, harus cek DB untuk memastikan IsActive",
			merchantID: "m2",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().GetStatusCache(ctx, "m2").Return("open", nil).Times(1)
				// DB dicek dan terkonfirmasi aktif
				mockRepo.EXPECT().GetByID(ctx, "m2").Return(model.Merchant{IsActive: true, StatusOperasional: "open"}, nil).Times(1)
			},
			expectedRes:   true,
			expectedError: "",
		},
		{
			name:       "Positive Case: Cache miss (error), fallback ke DB dan open",
			merchantID: "m3",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().GetStatusCache(ctx, "m3").Return("", errors.New("cache miss")).Times(1)
				mockRepo.EXPECT().GetByID(ctx, "m3").Return(model.Merchant{IsActive: true, StatusOperasional: "open"}, nil).Times(1)
			},
			expectedRes:   true,
			expectedError: "",
		},
		{
			name:       "Negative Case: Cache miss, DB Error",
			merchantID: "m4",
			mockSetup: func(mockRepo *MockMerchantRepository) {
				mockRepo.EXPECT().GetStatusCache(ctx, "m4").Return("", errors.New("cache miss")).Times(1)
				mockRepo.EXPECT().GetByID(ctx, "m4").Return(model.Merchant{}, errors.New("db error")).Times(1)
			},
			expectedRes:   false,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := NewMockMerchantRepository(ctrl)
			tt.mockSetup(mockRepo)
			svc := service.NewMerchantService(mockRepo)
			res, err := svc.CheckMerchantStatus(ctx, tt.merchantID)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
