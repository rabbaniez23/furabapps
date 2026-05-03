// Package service implements the business logic for chat-service.
package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/chat-service/internal/model"
	"furab-backend/services/chat-service/internal/repository"
)

// NotificationClient defines the interface for calling Notification Service.
type NotificationClient interface {
	SendNotification(ctx context.Context, receiverID string, messagePreview string) error
}

// ChatService defines the interface for chat-service business logic.
type ChatService interface {
	SendMessage(ctx context.Context, req model.SendMessageRequest) (*model.SendMessageResponse, error)
	UpdateReadStatus(ctx context.Context, req model.ReadReceiptRequest) error
	GetChatHistory(ctx context.Context, orderID string) ([]model.Message, error)
}

// chatServiceImpl is the concrete implementation of ChatService.
type chatServiceImpl struct {
	repo        repository.ChatRepository
	notifClient NotificationClient
}

// NewChatService creates a new ChatService.
func NewChatService(repo repository.ChatRepository, notifClient NotificationClient) ChatService {
	return &chatServiceImpl{
		repo:        repo,
		notifClient: notifClient,
	}
}

// =======================
// SEND MESSAGE
// =======================
func (s *chatServiceImpl) SendMessage(ctx context.Context, req model.SendMessageRequest) (*model.SendMessageResponse, error) {

	// ✅ Validasi input
	if req.MessageText == "" {
		return nil, errors.New("message cannot be empty")
	}

	if req.SenderType != "user" && req.SenderType != "driver" {
		return nil, errors.New("invalid sender type")
	}

	if req.ReceiverID == "" {
		return nil, errors.New("receiver required")
	}

	// ✅ Validasi chat session
	session, err := s.repo.GetChatSession(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, errors.New("chat session not found")
	}

	if session.ClosedAt != nil {
		return nil, errors.New("chat session closed")
	}

	// ✅ Create message
	msg := model.Message{
		OrderID:   req.OrderID,
		SenderID:  req.SenderID,
		Content:   req.MessageText,
		Timestamp: time.Now(),
	}

	// ✅ Save ke repository
	if err := s.repo.SaveMessage(ctx, msg); err != nil {
		return nil, err
	}

	// ✅ Trigger notification
	if err := s.notifClient.SendNotification(ctx, req.ReceiverID, req.MessageText); err != nil {
		return nil, err
	}

	// ✅ Response sukses
	return &model.SendMessageResponse{
		Status: "success",
	}, nil
}

// =======================
// UPDATE READ STATUS
// =======================
func (s *chatServiceImpl) UpdateReadStatus(ctx context.Context, req model.ReadReceiptRequest) error {

	// ✅ Validasi status
	if req.Status != "delivered" && req.Status != "read" {
		return errors.New("invalid status")
	}

	// ✅ Update ke repository
	err := s.repo.UpdateMessageStatus(ctx, req.MessageID, req.Status)
	if err != nil {
		return err
	}

	return nil
}

// =======================
// GET CHAT HISTORY
// =======================
func (s *chatServiceImpl) GetChatHistory(ctx context.Context, orderID string) ([]model.Message, error) {

	if orderID == "" {
		return nil, errors.New("order id required")
	}

	messages, err := s.repo.GetMessagesByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return messages, nil
}