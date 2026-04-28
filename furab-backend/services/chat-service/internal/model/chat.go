// Package model defines the domain models for chat-service.
package model

import "time"

// Conversation represents the Conversation model in chat-service.
type Conversation struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Conversation-specific fields
}

// Message represents the Message model in chat-service.
type Message struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Message-specific fields
}

