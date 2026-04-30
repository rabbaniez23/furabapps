package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"furab-backend/services/emergency-service/internal/model"
	"furab-backend/services/emergency-service/internal/service"
	"furab-backend/services/emergency-service/test/unit/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTriggerEmergency(t *testing.T) {
	t.Run("Success - Emergency berhasil dibuat", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

		now := time.Now()
		req := model.TriggerEmergencyRequest{
			ActorID:       "user123",
			ActorType:     "user",
			OrderID:       "order001",
			Latitude:      -6.2000,
			Longitude:     106.8166,
			EmergencyType: "accident",
			Timestamp:     now,
		}

		mockActor.EXPECT().ValidateActor(gomock.Any(), "user123", "user").Return(true, nil)
		mockActor.EXPECT().ValidateOrder(gomock.Any(), "order001").Return(true, nil)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), "user123", "user").Return(&model.EmergencyLocation{
			Latitude:  -6.2000,
			Longitude: 106.8166,
			Timestamp: now,
			Accuracy:  7,
		}, nil)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil)
		mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "user123", "user").Return(&model.EmergencyContact{
			ReceiverID: "ops-team",
		}, nil)
		mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyContact{}), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)

		res, err := svc.TriggerEmergency(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.Equal(t, "emergency created", res.Message)
		assert.NotEmpty(t, res.EmergencyID)
	})

	t.Run("Error - Actor tidak valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

		req := model.TriggerEmergencyRequest{
			ActorID:       "",
			ActorType:     "user",
			Latitude:      -6.2000,
			Longitude:     106.8166,
			EmergencyType: "accident",
			Timestamp:     time.Now(),
		}

		mockActor.EXPECT().ValidateActor(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.Any()).Times(0)
		mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.Any()).Times(0)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.TriggerEmergency(context.Background(), req)
		require.Error(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "invalid actor", err.Error())
		assert.Equal(t, "failed", res.Status)
		assert.Equal(t, "invalid actor", res.Message)
	})

	t.Run("Error - Location tidak tersedia", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

		req := model.TriggerEmergencyRequest{
			ActorID:       "user123",
			ActorType:     "user",
			Latitude:      -6.1,
			Longitude:     106.9,
			EmergencyType: "accident",
			Timestamp:     time.Now(),
		}

		mockActor.EXPECT().ValidateActor(gomock.Any(), "user123", "user").Return(true, nil).Times(1)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), "user123", "user").Return(nil, errors.New("location service unavailable")).Times(1)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil).Times(1)
		mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "user123", "user").Return(nil, errors.New("contact unavailable")).Times(1)
		mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.TriggerEmergency(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.NotEmpty(t, res.EmergencyID)
	})

	t.Run("Success - Emergency diproses walau order invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)
		now := time.Now()

		req := model.TriggerEmergencyRequest{
			ActorID:       "driver123",
			ActorType:     "driver",
			OrderID:       "invalid-order",
			Latitude:      -6.2000,
			Longitude:     106.8166,
			EmergencyType: "unsafe",
			Timestamp:     now,
		}

		mockActor.EXPECT().ValidateActor(gomock.Any(), "driver123", "driver").Return(true, nil)
		mockActor.EXPECT().ValidateOrder(gomock.Any(), "invalid-order").Return(false, nil)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), "driver123", "driver").Return(&model.EmergencyLocation{
			Latitude:  -6.2000,
			Longitude: 106.8166,
			Timestamp: now,
			Accuracy:  5,
		}, nil)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil)
		mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "driver123", "driver").Return(nil, errors.New("contact unavailable"))
		mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.TriggerEmergency(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
	})
}

func TestEmergencyNotificationFlow(t *testing.T) {
	req := model.TriggerEmergencyRequest{
		ActorID:       "user123",
		ActorType:     "user",
		OrderID:       "order001",
		Latitude:      -6.2000,
		Longitude:     106.8166,
		EmergencyType: "accident",
		Timestamp:     time.Now(),
	}

	t.Run("Success - Notification terkirim", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

		mockActor.EXPECT().ValidateActor(gomock.Any(), "user123", "user").Return(true, nil)
		mockActor.EXPECT().ValidateOrder(gomock.Any(), "order001").Return(true, nil)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), "user123", "user").Return(&model.EmergencyLocation{
			Latitude:  -6.2,
			Longitude: 106.8166,
			Timestamp: req.Timestamp,
		}, nil)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil)
		mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "user123", "user").Return(&model.EmergencyContact{ReceiverID: "dispatcher-1"}, nil)
		mockNotification.EXPECT().
			SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).
			DoAndReturn(func(ctx context.Context, notif model.EmergencyNotification) error {
				require.NotEmpty(t, notif.ReceiverID)
				assert.Equal(t, "Emergency Alert", notif.Title)
				assert.Equal(t, "high", notif.Priority)
				return nil
			}).Times(1)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyContact{}), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)

		_, err := svc.TriggerEmergency(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("Error - Notification gagal tidak menghentikan flow", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock.NewMockEmergencyRepository(ctrl)
		mockLocation := mock.NewMockLocationClient(ctrl)
		mockActor := mock.NewMockActorClient(ctrl)
		mockNotification := mock.NewMockNotificationClient(ctrl)
		svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

		mockActor.EXPECT().ValidateActor(gomock.Any(), "user123", "user").Return(true, nil)
		mockActor.EXPECT().ValidateOrder(gomock.Any(), "order001").Return(true, nil)
		mockLocation.EXPECT().GetLastLocation(gomock.Any(), "user123", "user").Return(&model.EmergencyLocation{
			Latitude:  -6.2,
			Longitude: 106.8166,
			Timestamp: req.Timestamp,
		}, nil)
		mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil)
		mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "user123", "user").Return(&model.EmergencyContact{ReceiverID: "dispatcher-1"}, nil)
		mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(errors.New("publish failed")).Times(1)
		mockNotification.EXPECT().SendEmergencyContact(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyContact{}), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)

		res, err := svc.TriggerEmergency(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
	})
}

func TestEmergencyContactEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockEmergencyRepository(ctrl)
	mockLocation := mock.NewMockLocationClient(ctrl)
	mockActor := mock.NewMockActorClient(ctrl)
	mockNotification := mock.NewMockNotificationClient(ctrl)
	svc := service.NewEmergencyServiceWithDependencies(mockRepo, mockLocation, mockActor, mockNotification)

	req := model.TriggerEmergencyRequest{
		ActorID:       "driver42",
		ActorType:     "driver",
		OrderID:       "order999",
		Latitude:      -6.3,
		Longitude:     106.9,
		EmergencyType: "other",
		Timestamp:     time.Now(),
	}

	mockActor.EXPECT().ValidateActor(gomock.Any(), "driver42", "driver").Return(true, nil)
	mockActor.EXPECT().ValidateOrder(gomock.Any(), "order999").Return(true, nil)
	mockLocation.EXPECT().GetLastLocation(gomock.Any(), "driver42", "driver").Return(&model.EmergencyLocation{
		Latitude:  -6.3,
		Longitude: 106.9,
		Timestamp: req.Timestamp,
	}, nil)
	mockRepo.EXPECT().SaveEmergencyEvent(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyEvent{})).Return(nil)
	mockActor.EXPECT().GetEmergencyContact(gomock.Any(), "driver42", "driver").Return(&model.EmergencyContact{
		ReceiverID: "family-001",
		Phone:      "08123",
		Email:      "family@example.com",
	}, nil)
	mockNotification.EXPECT().SendNotification(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyNotification{})).Return(nil).Times(1)
	mockNotification.EXPECT().
		SendEmergencyContact(gomock.Any(), gomock.AssignableToTypeOf(model.EmergencyContact{}), gomock.AssignableToTypeOf(model.EmergencyNotification{})).
		DoAndReturn(func(ctx context.Context, contact model.EmergencyContact, notif model.EmergencyNotification) error {
			assert.Equal(t, "family-001", contact.ReceiverID)
			assert.NotEmpty(t, contact.Phone)
			assert.NotEmpty(t, notif.LocationURL)
			return nil
		}).Times(1)

	res, err := svc.TriggerEmergency(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "success", res.Status)
}
