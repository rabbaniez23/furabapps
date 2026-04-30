package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"furab-backend/services/review-service/internal/model"
	"furab-backend/services/review-service/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ==========================================
// Mock Review Repository
// ==========================================

type MockReviewRepository struct {
	ctrl     *gomock.Controller
	recorder *MockReviewRepositoryMockRecorder
}

type MockReviewRepositoryMockRecorder struct {
	mock *MockReviewRepository
}

func NewMockReviewRepository(ctrl *gomock.Controller) *MockReviewRepository {
	mock := &MockReviewRepository{ctrl: ctrl}
	mock.recorder = &MockReviewRepositoryMockRecorder{mock}
	return mock
}

func (m *MockReviewRepository) EXPECT() *MockReviewRepositoryMockRecorder {
	return m.recorder
}

// Create
func (m *MockReviewRepository) Create(ctx context.Context, review model.Review) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, review)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockReviewRepositoryMockRecorder) Create(ctx, review interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReviewRepository)(nil).Create), ctx, review)
}

// GetByTarget
func (m *MockReviewRepository) GetByTarget(ctx context.Context, targetID, targetType string, page, limit int) ([]model.Review, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTarget", ctx, targetID, targetType, page, limit)
	ret0, _ := ret[0].([]model.Review)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockReviewRepositoryMockRecorder) GetByTarget(ctx, targetID, targetType, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTarget", reflect.TypeOf((*MockReviewRepository)(nil).GetByTarget), ctx, targetID, targetType, page, limit)
}

// GetByOrderID
func (m *MockReviewRepository) GetByOrderID(ctx context.Context, orderID, targetType string) (model.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOrderID", ctx, orderID, targetType)
	ret0, _ := ret[0].(model.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockReviewRepositoryMockRecorder) GetByOrderID(ctx, orderID, targetType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOrderID", reflect.TypeOf((*MockReviewRepository)(nil).GetByOrderID), ctx, orderID, targetType)
}

// CreateReport
func (m *MockReviewRepository) CreateReport(ctx context.Context, report model.ReviewReport) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReport", ctx, report)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockReviewRepositoryMockRecorder) CreateReport(ctx, report interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReport", reflect.TypeOf((*MockReviewRepository)(nil).CreateReport), ctx, report)
}

// UpdateStatus
func (m *MockReviewRepository) UpdateStatus(ctx context.Context, reviewID string, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, reviewID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockReviewRepositoryMockRecorder) UpdateStatus(ctx, reviewID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockReviewRepository)(nil).UpdateStatus), ctx, reviewID, status)
}

// GetHistory
func (m *MockReviewRepository) GetHistory(ctx context.Context, userID string, targetType string, page, limit int) ([]model.Review, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", ctx, userID, targetType, page, limit)
	ret0, _ := ret[0].([]model.Review)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockReviewRepositoryMockRecorder) GetHistory(ctx, userID, targetType, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockReviewRepository)(nil).GetHistory), ctx, userID, targetType, page, limit)
}

// ==========================================
// Unit Tests (Table-Driven)
// ==========================================

func TestReviewService_Create(t *testing.T) {
	ctx := context.Background()

	validReview := model.Review{
		UserID:     "u1",
		TargetID:   "d1",
		TargetType: "driver",
		OrderID:    "ord-123",
		Comment:    "Mantap!",
	}

	tests := []struct {
		name          string
		input         model.Review
		mockSetup     func(mockRepo *MockReviewRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Berhasil kirim ulasan",
			input: validReview,
			mockSetup: func(mockRepo *MockReviewRepository) {
				// Cek apakah order sdh diulas
				mockRepo.EXPECT().GetByOrderID(ctx, "ord-123", "driver").Return(model.Review{}, errors.New("not found")).Times(1)
				
				// Create the review, with status injected by service
				expectedCreate := validReview
				expectedCreate.Status = "active"
				mockRepo.EXPECT().Create(ctx, expectedCreate).Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name: "Negative Case: TargetType invalid (user)",
			input: model.Review{
				TargetType: "user",
				OrderID:    "ord-123",
			},
			mockSetup: func(mockRepo *MockReviewRepository) {},
			expectedError: "INVALID_TARGET_TYPE",
		},
		{
			name: "Negative Case: Order belum selesai (OrderID kosong/invalid)",
			input: model.Review{
				TargetType: "driver",
				OrderID:    "invalid_order", // matches simulateOrderCheck condition
			},
			mockSetup: func(mockRepo *MockReviewRepository) {},
			expectedError: "ORDER_NOT_COMPLETED",
		},
		{
			name:  "Negative Case: Sudah pernah ulas (ALREADY_REVIEWED)",
			input: validReview,
			mockSetup: func(mockRepo *MockReviewRepository) {
				// Return existing review with ID
				mockRepo.EXPECT().GetByOrderID(ctx, "ord-123", "driver").Return(model.Review{ReviewID: "rev-1"}, nil).Times(1)
			},
			expectedError: "ALREADY_REVIEWED",
		},
		{
			name:  "Negative Case: Error saat Create",
			input: validReview,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetByOrderID(ctx, "ord-123", "driver").Return(model.Review{}, errors.New("not found")).Times(1)
				expectedCreate := validReview
				expectedCreate.Status = "active"
				mockRepo.EXPECT().Create(ctx, expectedCreate).Return(errors.New("db error")).Times(1)
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockReviewRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewReviewService(mockRepo)
			err := svc.Create(ctx, tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestReviewService_GetByTarget(t *testing.T) {
	ctx := context.Background()

	expectedReviews := []model.Review{
		{ReviewID: "r1", Comment: "Bagus"},
	}

	tests := []struct {
		name          string
		targetID      string
		targetType    string
		page          int
		limit         int
		mockSetup     func(mockRepo *MockReviewRepository)
		expectedCount int
		expectedError string
	}{
		{
			name:       "Positive Case: Dengan pagination valid",
			targetID:   "d1",
			targetType: "driver",
			page:       2,
			limit:      5,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetByTarget(ctx, "d1", "driver", 2, 5).Return(expectedReviews, 1, nil).Times(1)
			},
			expectedCount: 1,
			expectedError: "",
		},
		{
			name:       "Positive Case: Default pagination if negative/zero",
			targetID:   "m1",
			targetType: "merchant",
			page:       -1, // should default to 1
			limit:      0,  // should default to 10
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetByTarget(ctx, "m1", "merchant", 1, 10).Return(expectedReviews, 1, nil).Times(1)
			},
			expectedCount: 1,
			expectedError: "",
		},
		{
			name:       "Negative Case: Error dari repository",
			targetID:   "d1",
			targetType: "driver",
			page:       1,
			limit:      10,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetByTarget(ctx, "d1", "driver", 1, 10).Return(nil, 0, errors.New("db error")).Times(1)
			},
			expectedCount: 0,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockReviewRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewReviewService(mockRepo)
			res, count, err := svc.GetByTarget(ctx, tt.targetID, tt.targetType, tt.page, tt.limit)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, expectedReviews, res)
				assert.Equal(t, tt.expectedCount, count)
			} else {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, res)
			}
		})
	}
}

func TestReviewService_CreateReport(t *testing.T) {
	ctx := context.Background()

	validReport := model.ReviewReport{
		ReviewID:   "rev-1",
		ReportedBy: "u2",
		Reason:     "Spam",
	}

	tests := []struct {
		name          string
		input         model.ReviewReport
		mockSetup     func(mockRepo *MockReviewRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Report sukses dan status diubah jadi flagged",
			input: validReport,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().CreateReport(ctx, validReport).Return(nil).Times(1)
				// Pastikan setelahnya memanggil UpdateStatus menjadi "flagged"
				mockRepo.EXPECT().UpdateStatus(ctx, "rev-1", "flagged").Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:  "Negative Case: Gagal create report",
			input: validReport,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().CreateReport(ctx, validReport).Return(errors.New("db error")).Times(1)
				// UpdateStatus seharusnya tidak dipanggil
			},
			expectedError: "db error",
		},
		{
			name:  "Negative Case: Gagal update status review",
			input: validReport,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().CreateReport(ctx, validReport).Return(nil).Times(1)
				mockRepo.EXPECT().UpdateStatus(ctx, "rev-1", "flagged").Return(errors.New("update err")).Times(1)
			},
			expectedError: "update err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockReviewRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewReviewService(mockRepo)
			err := svc.CreateReport(ctx, tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestReviewService_GetHistory(t *testing.T) {
	ctx := context.Background()

	expectedReviews := []model.Review{
		{ReviewID: "r1", Comment: "Ok"},
	}

	tests := []struct {
		name          string
		userID        string
		targetType    string
		page          int
		limit         int
		mockSetup     func(mockRepo *MockReviewRepository)
		expectedError string
	}{
		{
			name:       "Positive Case: GetHistory dengan target_type spesifik",
			userID:     "u1",
			targetType: "merchant",
			page:       1,
			limit:      10,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetHistory(ctx, "u1", "merchant", 1, 10).Return(expectedReviews, 1, nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:       "Positive Case: GetHistory tanpa target_type (empty string)",
			userID:     "u1",
			targetType: "",
			page:       0, // fallback to 1
			limit:      0, // fallback to 10
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetHistory(ctx, "u1", "", 1, 10).Return(expectedReviews, 1, nil).Times(1)
			},
			expectedError: "",
		},
		{
			name:       "Negative Case: Error dari repo",
			userID:     "u1",
			targetType: "driver",
			page:       1,
			limit:      10,
			mockSetup: func(mockRepo *MockReviewRepository) {
				mockRepo.EXPECT().GetHistory(ctx, "u1", "driver", 1, 10).Return(nil, 0, errors.New("db err")).Times(1)
			},
			expectedError: "db err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockReviewRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewReviewService(mockRepo)
			res, count, err := svc.GetHistory(ctx, tt.userID, tt.targetType, tt.page, tt.limit)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, expectedReviews, res)
				assert.Equal(t, 1, count)
			} else {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, res)
			}
		})
	}
}
