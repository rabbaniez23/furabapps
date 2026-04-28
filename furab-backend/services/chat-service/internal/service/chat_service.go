// Package service implements the business logic for chat-service.
package service

import "context"

// ChatService defines the interface for chat-service business logic.
type ChatService interface {

	// SendMessage implements the business logic for SendMessage.
	SendMessage(ctx context.Context) error

	// GetMessages implements the business logic for GetMessages.
	GetMessages(ctx context.Context) error

	// GetConversation implements the business logic for GetConversation.
	GetConversation(ctx context.Context) error

	// MarkAsRead implements the business logic for MarkAsRead.
	MarkAsRead(ctx context.Context) error
}

// chatServiceImpl is the concrete implementation of ChatService.
type chatServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewChatService creates a new ChatService.
func NewChatService() ChatService {
	return &chatServiceImpl{}
}
