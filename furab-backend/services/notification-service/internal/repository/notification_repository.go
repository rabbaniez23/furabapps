// Package repository provides data access layer for notification-service.
package repository

import "context"

// NotificationRepository defines the interface for notification-service data access.
type NotificationRepository interface {

	// Send performs the Send operation.
	Send(ctx context.Context) error

	// GetAll performs the GetAll operation.
	GetAll(ctx context.Context) error

	// MarkAsRead performs the MarkAsRead operation.
	MarkAsRead(ctx context.Context) error

	// GetUnreadCount performs the GetUnreadCount operation.
	GetUnreadCount(ctx context.Context) error
}

// postgresNotificationRepository implements NotificationRepository using PostgreSQL.
type postgresNotificationRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresNotificationRepository creates a new PostgreSQL-based repository.
func NewPostgresNotificationRepository() NotificationRepository {
	return &postgresNotificationRepository{}
}
