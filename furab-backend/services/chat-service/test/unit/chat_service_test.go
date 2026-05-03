package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"furab-backend/services/chat-service/internal/model"
	"furab-backend/services/chat-service/internal/repository/mock_repository"
	"furab-backend/services/chat-service/internal/service"
	"furab-backend/services/chat-service/internal/service/mock_service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockChatRepository(ctrl)
	mockNotif := mock_service.NewMockNotificationClient(ctrl)
	svc := service.NewChatService(mockRepo, mockNotif)

	t.Run("Success - Pesan berhasil dikirim", func(t *testing.T) {
		req := model.SendMessageRequest{
			OrderID:     "order1",
			SenderID:    "user1",
			SenderType:  "user",
			ReceiverID:  "driver1",
			MessageText: "Halo, saya sudah di titik penjemputan",
		}

		activeSession := &model.ChatSession{
			OrderID:  "order1",
			UserID:   "user1",
			DriverID: "driver1",
		}

		mockRepo.EXPECT().GetChatSession(gomock.Any(), "order1").Return(activeSession, nil)
		mockRepo.EXPECT().SaveMessage(gomock.Any(), gomock.Any()).Return(nil)
		mockNotif.EXPECT().SendNotification(gomock.Any(), "driver1", gomock.Any()).Return(nil)

		res, err := svc.SendMessage(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
	})

	t.Run("Error - Pesan kosong", func(t *testing.T) {
		req := model.SendMessageRequest{
			OrderID:     "order1",
			SenderID:    "user1",
			SenderType:  "user",
			ReceiverID:  "driver1",
			MessageText: "",
		}

		mockRepo.EXPECT().GetChatSession(gomock.Any(), gomock.Any()).Times(0)
		mockRepo.EXPECT().SaveMessage(gomock.Any(), gomock.Any()).Times(0)
		mockNotif.EXPECT().SendNotification(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.SendMessage(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "message cannot be empty", err.Error())
	})

	t.Run("Error - Sender tidak valid", func(t *testing.T) {
		req := model.SendMessageRequest{
			OrderID:     "order1",
			SenderID:    "admin1",
			SenderType:  "admin", // Invalid
			ReceiverID:  "driver1",
			MessageText: "Halo",
		}

		res, err := svc.SendMessage(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid sender type", err.Error())
	})

	t.Run("Error - Chat session tidak aktif", func(t *testing.T) {
		req := model.SendMessageRequest{
			OrderID:     "order1",
			SenderID:    "user1",
			SenderType:  "user",
			ReceiverID:  "driver1",
			MessageText: "Halo",
		}

		// Simulate closed session
		closedTime := time.Now()
		closedSession := &model.ChatSession{
			OrderID:  "order1",
			UserID:   "user1",
			DriverID: "driver1",
			ClosedAt: &closedTime,
		}

		mockRepo.EXPECT().GetChatSession(gomock.Any(), "order1").Return(closedSession, nil)
		mockRepo.EXPECT().SaveMessage(gomock.Any(), gomock.Any()).Times(0)

		res, err := svc.SendMessage(context.Background(), req)
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "chat session closed", err.Error())
	})
}

func TestUpdateReadStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockChatRepository(ctrl)
	mockNotif := mock_service.NewMockNotificationClient(ctrl)
	svc := service.NewChatService(mockRepo, mockNotif)

	t.Run("Success - Status pesan diperbarui", func(t *testing.T) {
		req := model.ReadReceiptRequest{
			MessageID: "msg1",
			OrderID:   "order1",
			Status:    "read",
		}

		mockRepo.EXPECT().UpdateMessageStatus(gomock.Any(), "msg1", "read").Return(nil)

		err := svc.UpdateReadStatus(context.Background(), req)
		require.NoError(t, err)
	})

	t.Run("Error - Status tidak valid", func(t *testing.T) {
		req := model.ReadReceiptRequest{
			MessageID: "msg1",
			OrderID:   "order1",
			Status:    "seen", // Invalid
		}

		mockRepo.EXPECT().UpdateMessageStatus(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		err := svc.UpdateReadStatus(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "invalid status", err.Error())
	})

	t.Run("Error - Message tidak ditemukan", func(t *testing.T) {
		req := model.ReadReceiptRequest{
			MessageID: "msg_unknown",
			OrderID:   "order1",
			Status:    "read",
		}

		mockRepo.EXPECT().UpdateMessageStatus(gomock.Any(), "msg_unknown", "read").Return(errors.New("message not found"))

		err := svc.UpdateReadStatus(context.Background(), req)
		require.Error(t, err)
		assert.Equal(t, "message not found", err.Error())
	})
}

func TestGetChatHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockChatRepository(ctrl)
	mockNotif := mock_service.NewMockNotificationClient(ctrl)
	svc := service.NewChatService(mockRepo, mockNotif)

	t.Run("Success - Riwayat ditemukan", func(t *testing.T) {
		mockMessages := []model.Message{
			{MessageID: "msg1", Content: "Halo"},
			{MessageID: "msg2", Content: "Ya, tunggu sebentar"},
		}

		mockRepo.EXPECT().GetMessagesByOrderID(gomock.Any(), "order1").Return(mockMessages, nil)

		res, err := svc.GetChatHistory(context.Background(), "order1")
		require.NoError(t, err)
		require.Len(t, res, 2)
		assert.Equal(t, "msg1", res[0].MessageID)
	})

	t.Run("Success - Tidak ada riwayat", func(t *testing.T) {
		mockRepo.EXPECT().GetMessagesByOrderID(gomock.Any(), "order1").Return([]model.Message{}, nil)

		res, err := svc.GetChatHistory(context.Background(), "order1")
		require.NoError(t, err)
		assert.Len(t, res, 0)
	})

	t.Run("Error - Order tidak valid", func(t *testing.T) {
		mockRepo.EXPECT().GetMessagesByOrderID(gomock.Any(), "invalid_order").Return(nil, errors.New("order not found"))

		res, err := svc.GetChatHistory(context.Background(), "invalid_order")
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "order not found", err.Error())
	})
}
