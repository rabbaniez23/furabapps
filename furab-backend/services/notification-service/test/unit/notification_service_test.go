package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/notification-service/internal/model"
	"furab-backend/services/notification-service/internal/repository/mock_repository"
	"furab-backend/services/notification-service/internal/service"
	"furab-backend/services/notification-service/internal/service/mock_service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProcessEventNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockNotificationRepository(ctrl)
	mockEmail := mock_service.NewMockEmailClient(ctrl)
	svc := service.NewNotificationService(mockRepo, mockEmail)

	t.Run("Success - Push Notification", func(t *testing.T) {
		req := model.EventNotificationRequest{
			EventType:   "ride.created",
			ReceiverID:  "user1",
			ReferenceID: "order1",
			Channel:     "push",
		}

		mockTemplate := &model.NotifTemplate{
			EventType:       "ride.created",
			TitleTemplate:   "Ride Created",
			MessageTemplate: "Your ride is created",
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), "ride.created").Return(mockTemplate, nil)
		mockRepo.EXPECT().SaveNotificationLog(gomock.Any(), gomock.Any()).Return(nil)
		mockEmail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.ProcessEventNotification(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.Equal(t, "notifikasi berhasil", res.Message)
	})

	t.Run("Success - Email Notification", func(t *testing.T) {
		req := model.EventNotificationRequest{
			EventType:   "payment.success",
			ReceiverID:  "user1",
			ReferenceID: "order1",
			Channel:     "email",
		}

		mockTemplate := &model.NotifTemplate{
			EventType:       "payment.success",
			TitleTemplate:   "Payment Success",
			MessageTemplate: "Your payment is successful",
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), "payment.success").Return(mockTemplate, nil)
		mockRepo.EXPECT().SaveNotificationLog(gomock.Any(), gomock.Any()).Return(nil)
		mockEmail.EXPECT().SendEmail(gomock.Any(), "user1", "Payment Success", "Your payment is successful").Return(nil)

		res, err := svc.ProcessEventNotification(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.Equal(t, "notifikasi berhasil", res.Message)
	})

	t.Run("Error - Channel tidak valid", func(t *testing.T) {
		req := model.EventNotificationRequest{
			EventType:   "ride.created",
			ReceiverID:  "user1",
			ReferenceID: "order1",
			Channel:     "sms", // Invalid
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().SaveNotificationLog(gomock.Any(), gomock.Any()).Times(0)
		mockEmail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.ProcessEventNotification(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid channel", err.Error())
	})

	t.Run("Error - Event kosong", func(t *testing.T) {
		req := model.EventNotificationRequest{
			EventType:   "",
			ReceiverID:  "user1",
			ReferenceID: "order1",
			Channel:     "push",
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().SaveNotificationLog(gomock.Any(), gomock.Any()).Times(0)
		mockEmail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.ProcessEventNotification(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "event type is required", err.Error())
	})

	t.Run("Error - Template tidak ditemukan", func(t *testing.T) {
		req := model.EventNotificationRequest{
			EventType:   "unknown.event",
			ReceiverID:  "user1",
			ReferenceID: "order1",
			Channel:     "push",
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), "unknown.event").Return(nil, errors.New("template not found"))
		mockRepo.EXPECT().SaveNotificationLog(gomock.Any(), gomock.Any()).Times(0)
		mockEmail.EXPECT().SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.ProcessEventNotification(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "template not found", err.Error())
	})
}

func TestGenerateNotificationTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockNotificationRepository(ctrl)
	mockEmail := mock_service.NewMockEmailClient(ctrl)
	svc := service.NewNotificationService(mockRepo, mockEmail)

	t.Run("Success - Template ditemukan", func(t *testing.T) {
		mockTemplate := &model.NotifTemplate{
			EventType:       "ride.created",
			TitleTemplate:   "Title",
			MessageTemplate: "Message",
		}

		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), "ride.created").Return(mockTemplate, nil)

		res, err := svc.GenerateNotificationTemplate(context.Background(), "ride.created")
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "Title", res.TitleTemplate)
	})

	t.Run("Error - Template tidak ditemukan", func(t *testing.T) {
		mockRepo.EXPECT().GetTemplateByEventType(gomock.Any(), "invalid.event").Return(nil, errors.New("template not found"))

		res, err := svc.GenerateNotificationTemplate(context.Background(), "invalid.event")
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "template not found", err.Error())
	})
}
