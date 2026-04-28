// Package model defines the domain models for notification-service.
package model

import "time"

// Notification represents the Notification model in notification-service.
type Notification struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Notification-specific fields
}

// NotifTemplate represents the NotifTemplate model in notification-service.
type NotifTemplate struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add NotifTemplate-specific fields
}

