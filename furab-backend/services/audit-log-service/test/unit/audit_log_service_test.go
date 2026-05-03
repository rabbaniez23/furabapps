package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"furab-backend/services/audit-log-service/internal/model"
	"furab-backend/services/audit-log-service/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ==========================================
// Mock Audit Log Repository
// ==========================================

type MockAuditLogRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuditLogRepositoryMockRecorder
}

type MockAuditLogRepositoryMockRecorder struct {
	mock *MockAuditLogRepository
}

func NewMockAuditLogRepository(ctrl *gomock.Controller) *MockAuditLogRepository {
	mock := &MockAuditLogRepository{ctrl: ctrl}
	mock.recorder = &MockAuditLogRepositoryMockRecorder{mock}
	return mock
}

func (m *MockAuditLogRepository) EXPECT() *MockAuditLogRepositoryMockRecorder {
	return m.recorder
}

// Save
func (m *MockAuditLogRepository) Save(ctx context.Context, log model.AuditLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, log)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockAuditLogRepositoryMockRecorder) Save(ctx, log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAuditLogRepository)(nil).Save), ctx, log)
}

// GetByID
func (m *MockAuditLogRepository) GetByID(ctx context.Context, logID string) (model.AuditLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, logID)
	ret0, _ := ret[0].(model.AuditLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuditLogRepositoryMockRecorder) GetByID(ctx, logID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAuditLogRepository)(nil).GetByID), ctx, logID)
}

// Search
func (m *MockAuditLogRepository) Search(ctx context.Context, filter map[string]interface{}, page, limit int) ([]model.AuditLog, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, filter, page, limit)
	ret0, _ := ret[0].([]model.AuditLog)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockAuditLogRepositoryMockRecorder) Search(ctx, filter, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockAuditLogRepository)(nil).Search), ctx, filter, page, limit)
}

// ==========================================
// Unit Tests (Table-Driven)
// ==========================================

func TestAuditLogService_RecordLog(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	validLog := model.AuditLog{
		ServiceName: "user-service",
		ActorID:     "admin-1",
		Action:      "UPDATE_USER",
		Status:      "SUCCESS",
		Metadata:    map[string]interface{}{"ip": "192.168.1.1"},
		Timestamp:   now,
	}

	tests := []struct {
		name          string
		input         model.AuditLog
		mockSetup     func(mockRepo *MockAuditLogRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Sukses record log dengan metadata",
			input: validLog,
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				mockRepo.EXPECT().Save(ctx, validLog).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:  "Positive Case: Sukses record log, Timestamp otomatis diset jika kosong",
			input: model.AuditLog{
				ServiceName: "order-service",
				ActorID:     "system",
				Action:      "CRON_JOB",
				Status:      "FAILED",
			},
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				// gomock.Any() allows matching the struct despite Timestamp being dynamically generated inside the service
				mockRepo.EXPECT().Save(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, log model.AuditLog) error {
					assert.False(t, log.Timestamp.IsZero(), "Timestamp should be auto-generated")
					return nil
				}).Times(1)
			},
			expectedError: "",
		},
		{
			name: "Negative Case: ServiceName kosong",
			input: model.AuditLog{ActorID: "sys", Action: "RUN", Status: "OK"},
			mockSetup: func(mockRepo *MockAuditLogRepository) {},
			expectedError: "INVALID_PAYLOAD",
		},
		{
			name: "Negative Case: ActorID kosong",
			input: model.AuditLog{ServiceName: "svc", Action: "RUN", Status: "OK"},
			mockSetup: func(mockRepo *MockAuditLogRepository) {},
			expectedError: "INVALID_PAYLOAD",
		},
		{
			name: "Negative Case: Action kosong",
			input: model.AuditLog{ServiceName: "svc", ActorID: "sys", Status: "OK"},
			mockSetup: func(mockRepo *MockAuditLogRepository) {},
			expectedError: "INVALID_PAYLOAD",
		},
		{
			name: "Negative Case: Status kosong",
			input: model.AuditLog{ServiceName: "svc", ActorID: "sys", Action: "RUN"},
			mockSetup: func(mockRepo *MockAuditLogRepository) {},
			expectedError: "INVALID_PAYLOAD",
		},
		{
			name:  "Negative Case: Error DB saat save",
			input: validLog,
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				mockRepo.EXPECT().Save(ctx, validLog).Return(errors.New("db error")).Times(1)
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockAuditLogRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewAuditLogService(mockRepo)
			err := svc.RecordLog(ctx, tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestAuditLogService_SearchLog(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	expectedLogs := []model.AuditLog{
		{ServiceName: "auth-service", Action: "LOGIN"},
	}

	tests := []struct {
		name          string
		filter        service.SearchLogFilter
		page          int
		limit         int
		mockSetup     func(mockRepo *MockAuditLogRepository)
		expectedError string
		expectedCount int
	}{
		{
			name: "Positive Case: Filter by ServiceName dan date range",
			filter: service.SearchLogFilter{
				StartDate:   &startDate,
				EndDate:     &endDate,
				ServiceName: "auth-service",
			},
			page:  1,
			limit: 10,
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				expectedRepoFilter := map[string]interface{}{
					"start_date":   startDate,
					"end_date":     endDate,
					"service_name": "auth-service",
				}
				mockRepo.EXPECT().Search(ctx, expectedRepoFilter, 1, 10).Return(expectedLogs, 1, nil).Times(1)
			},
			expectedError: "",
			expectedCount: 1,
		},
		{
			name: "Positive Case: Filter all properties with fallback pagination",
			filter: service.SearchLogFilter{
				ActorID:  "u1",
				Action:   "CREATE",
				TargetID: "t1",
				Status:   "SUCCESS",
			},
			page:  0, // will fallback to 1
			limit: 0, // will fallback to 10
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				expectedRepoFilter := map[string]interface{}{
					"actor_id":  "u1",
					"action":    "CREATE",
					"target_id": "t1",
					"status":    "SUCCESS",
				}
				mockRepo.EXPECT().Search(ctx, expectedRepoFilter, 1, 10).Return(expectedLogs, 1, nil).Times(1)
			},
			expectedError: "",
			expectedCount: 1,
		},
		{
			name: "Negative Case: DB Error saat Search",
			filter: service.SearchLogFilter{},
			page:  1,
			limit: 10,
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				expectedRepoFilter := map[string]interface{}{}
				mockRepo.EXPECT().Search(ctx, expectedRepoFilter, 1, 10).Return(nil, 0, errors.New("db timeout")).Times(1)
			},
			expectedError: "db timeout",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockAuditLogRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewAuditLogService(mockRepo)
			res, count, err := svc.SearchLog(ctx, tt.filter, tt.page, tt.limit)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, expectedLogs, res)
				assert.Equal(t, tt.expectedCount, count)
			} else {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, res)
			}
		})
	}
}

func TestAuditLogService_GetLogDetail(t *testing.T) {
	ctx := context.Background()

	expectedLog := model.AuditLog{
		LogID:       "log-1",
		ServiceName: "payment-service",
		Metadata:    map[string]interface{}{"tx_id": "tx-123"},
	}

	tests := []struct {
		name          string
		logID         string
		mockSetup     func(mockRepo *MockAuditLogRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Detail log ditemukan (lengkap metadata)",
			logID: "log-1",
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				mockRepo.EXPECT().GetByID(ctx, "log-1").Return(expectedLog, nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:  "Negative Case: Log tidak ditemukan",
			logID: "log-404",
			mockSetup: func(mockRepo *MockAuditLogRepository) {
				mockRepo.EXPECT().GetByID(ctx, "log-404").Return(model.AuditLog{}, errors.New("not found")).Times(1)
			},
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockAuditLogRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewAuditLogService(mockRepo)
			res, err := svc.GetLogDetail(ctx, tt.logID)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, expectedLog, res)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
