package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"furab-backend/services/rating-service/internal/model"
	"furab-backend/services/rating-service/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// ==========================================
// Mock Repository (Manual gomock merujuk ke repository asli di package internal)
// ==========================================

type MockRatingRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRatingRepositoryMockRecorder
}

type MockRatingRepositoryMockRecorder struct {
	mock *MockRatingRepository
}

func NewMockRatingRepository(ctrl *gomock.Controller) *MockRatingRepository {
	mock := &MockRatingRepository{ctrl: ctrl}
	mock.recorder = &MockRatingRepositoryMockRecorder{mock}
	return mock
}

func (m *MockRatingRepository) EXPECT() *MockRatingRepositoryMockRecorder {
	return m.recorder
}

// CheckDuplicate mock
func (m *MockRatingRepository) CheckDuplicate(ctx context.Context, reviewerID, targetID, targetType, orderID string) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDuplicate", ctx, reviewerID, targetID, targetType, orderID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockRatingRepositoryMockRecorder) CheckDuplicate(ctx, reviewerID, targetID, targetType, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDuplicate", reflect.TypeOf((*MockRatingRepository)(nil).CheckDuplicate), ctx, reviewerID, targetID, targetType, orderID)
}

// SaveRating mock
func (m *MockRatingRepository) SaveRating(ctx context.Context, rating model.Rating) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRating", ctx, rating)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockRatingRepositoryMockRecorder) SaveRating(ctx, rating interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRating", reflect.TypeOf((*MockRatingRepository)(nil).SaveRating), ctx, rating)
}

// GetStatistics mock
func (m *MockRatingRepository) GetStatistics(ctx context.Context, targetID, targetType string) (model.RatingSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatistics", ctx, targetID, targetType)
	ret0, _ := ret[0].(model.RatingSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockRatingRepositoryMockRecorder) GetStatistics(ctx, targetID, targetType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatistics", reflect.TypeOf((*MockRatingRepository)(nil).GetStatistics), ctx, targetID, targetType)
}

// UpdateStatistics mock
func (m *MockRatingRepository) UpdateStatistics(ctx context.Context, targetID, targetType string, score int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatistics", ctx, targetID, targetType, score)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockRatingRepositoryMockRecorder) UpdateStatistics(ctx, targetID, targetType, score interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatistics", reflect.TypeOf((*MockRatingRepository)(nil).UpdateStatistics), ctx, targetID, targetType, score)
}

// GetHistory mock
func (m *MockRatingRepository) GetHistory(ctx context.Context, reviewerID string, page, limit int) ([]model.Rating, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", ctx, reviewerID, page, limit)
	ret0, _ := ret[0].([]model.Rating)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockRatingRepositoryMockRecorder) GetHistory(ctx, reviewerID, page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockRatingRepository)(nil).GetHistory), ctx, reviewerID, page, limit)
}

// ==========================================
// Unit Tests (Table-Driven)
// ==========================================

func TestRatingService_SubmitRating(t *testing.T) {
	ctx := context.Background()

	validRating := model.Rating{
		ReviewerID: "user-123",
		TargetID:   "driver-456",
		TargetType: "driver",
		OrderID:    "order-789",
		Score:      5,
	}

	tests := []struct {
		name          string
		input         model.Rating
		mockSetup     func(mockRepo *MockRatingRepository)
		expectedError string
	}{
		{
			name:  "Positive Case: Rating berhasil disimpan dan statistik diupdate",
			input: validRating,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					CheckDuplicate(ctx, validRating.ReviewerID, validRating.TargetID, validRating.TargetType, validRating.OrderID).
					Return(false, "", nil).Times(1)
				
				mockRepo.EXPECT().
					SaveRating(ctx, validRating).
					Return(nil).Times(1)

				mockRepo.EXPECT().
					UpdateStatistics(ctx, validRating.TargetID, validRating.TargetType, validRating.Score).
					Return(nil).Times(1)
			},
			expectedError: "",
		},
		{
			name: "Negative Case: Skor 0 (di luar batas bawah)",
			input: model.Rating{
				ReviewerID: "user-123",
				TargetID:   "driver-456",
				TargetType: "driver",
				OrderID:    "order-789",
				Score:      0, // invalid score
			},
			mockSetup: func(mockRepo *MockRatingRepository) {
				// Tidak ada interaksi repository karena gagal validasi
			},
			expectedError: "INVALID_SCORE",
		},
		{
			name: "Negative Case: Skor 6 (di luar batas atas)",
			input: model.Rating{
				ReviewerID: "user-123",
				TargetID:   "driver-456",
				TargetType: "driver",
				OrderID:    "order-789",
				Score:      6, // invalid score
			},
			mockSetup: func(mockRepo *MockRatingRepository) {
				// Tidak ada interaksi repository karena gagal validasi
			},
			expectedError: "INVALID_SCORE",
		},
		{
			name:  "Negative Case: Duplicate Error",
			input: validRating,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					CheckDuplicate(ctx, validRating.ReviewerID, validRating.TargetID, validRating.TargetType, validRating.OrderID).
					Return(true, "existing-rating-id", nil).Times(1)
			},
			expectedError: "ALREADY_RATED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRatingRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewRatingService(mockRepo)
			err := svc.SubmitRating(ctx, tt.input)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestRatingService_GetStatistics(t *testing.T) {
	ctx := context.Background()

	expectedSummary := model.RatingSummary{
		TargetType:   "driver",
		TargetID:     "driver-123",
		AverageScore: 4.8,
		TotalCount:   15,
	}

	tests := []struct {
		name          string
		targetID      string
		targetType    string
		mockSetup     func(mockRepo *MockRatingRepository)
		expectedError string
		expectedRes   model.RatingSummary
	}{
		{
			name:       "Positive Case: GetStatistics berhasil",
			targetID:   "driver-123",
			targetType: "driver",
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					GetStatistics(ctx, "driver-123", "driver").
					Return(expectedSummary, nil).Times(1)
			},
			expectedError: "",
			expectedRes:   expectedSummary,
		},
		{
			name:       "Negative Case: Repository error",
			targetID:   "driver-123",
			targetType: "driver",
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					GetStatistics(ctx, "driver-123", "driver").
					Return(model.RatingSummary{}, errors.New("db error")).Times(1)
			},
			expectedError: "db error",
			expectedRes:   model.RatingSummary{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRatingRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewRatingService(mockRepo)
			res, err := svc.GetStatistics(ctx, tt.targetID, tt.targetType)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedRes, res)
			}
		})
	}
}

func TestRatingService_GetHistory(t *testing.T) {
	ctx := context.Background()

	expectedHistory := []model.Rating{
		{ReviewerID: "user-1", Score: 5},
		{ReviewerID: "user-1", Score: 4},
	}

	tests := []struct {
		name          string
		reviewerID    string
		page          int
		limit         int
		mockSetup     func(mockRepo *MockRatingRepository)
		expectedError string
		expectedCount int
		expectedRes   []model.Rating
	}{
		{
			name:       "Positive Case: GetHistory dengan pagination valid",
			reviewerID: "user-1",
			page:       2,
			limit:      10,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					GetHistory(ctx, "user-1", 2, 10).
					Return(expectedHistory, 20, nil).Times(1)
			},
			expectedError: "",
			expectedCount: 20,
			expectedRes:   expectedHistory,
		},
		{
			name:       "Positive Case: GetHistory dengan page/limit negatif harus dinormalkan",
			reviewerID: "user-1",
			page:       0, // invalid, should default to 1
			limit:      0, // invalid, should default to 10
			mockSetup: func(mockRepo *MockRatingRepository) {
				// Service harus menormalisasi ke page=1 dan limit=10
				mockRepo.EXPECT().
					GetHistory(ctx, "user-1", 1, 10).
					Return(expectedHistory, 20, nil).Times(1)
			},
			expectedError: "",
			expectedCount: 20,
			expectedRes:   expectedHistory,
		},
		{
			name:       "Negative Case: Repository error",
			reviewerID: "user-1",
			page:       1,
			limit:      10,
			mockSetup: func(mockRepo *MockRatingRepository) {
				mockRepo.EXPECT().
					GetHistory(ctx, "user-1", 1, 10).
					Return(nil, 0, errors.New("db error")).Times(1)
			},
			expectedError: "db error",
			expectedCount: 0,
			expectedRes:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRatingRepository(ctrl)
			tt.mockSetup(mockRepo)

			svc := service.NewRatingService(mockRepo)
			res, count, err := svc.GetHistory(ctx, tt.reviewerID, tt.page, tt.limit)

			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
				assert.Equal(t, tt.expectedCount, count)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
