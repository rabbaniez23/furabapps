// Package model defines the domain models for audit-log-service.
package model

import "time"

// AuditLog represents the AuditLog model in audit-log-service.
type AuditLog struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add AuditLog-specific fields
}

