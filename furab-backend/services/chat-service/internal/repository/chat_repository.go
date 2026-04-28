// Package repository provides data access layer for chat-service.
package repository

import "context"

// ChatRepository defines the interface for chat-service data access.
type ChatRepository interface {

	// SendMessage performs the SendMessage operation.
	SendMessage(ctx context.Context) error

	// GetMessages performs the GetMessages operation.
	GetMessages(ctx context.Context) error

	// GetConversation performs the GetConversation operation.
	GetConversation(ctx context.Context) error

	// MarkAsRead performs the MarkAsRead operation.
	MarkAsRead(ctx context.Context) error
}

// postgresChatRepository implements ChatRepository using PostgreSQL.
type postgresChatRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresChatRepository creates a new PostgreSQL-based repository.
func NewPostgresChatRepository() ChatRepository {
	return &postgresChatRepository{}
}
