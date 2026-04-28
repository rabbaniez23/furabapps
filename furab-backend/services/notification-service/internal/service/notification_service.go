// Package service implements the business logic for notification-service.
package service

import "context"

// NotificationService defines the interface for notification-service business logic.
type NotificationService interface {

	// Send implements the business logic for Send.
	Send(ctx context.Context) error

	// GetAll implements the business logic for GetAll.
	GetAll(ctx context.Context) error

	// MarkAsRead implements the business logic for MarkAsRead.
	MarkAsRead(ctx context.Context) error

	// GetUnreadCount implements the business logic for GetUnreadCount.
	GetUnreadCount(ctx context.Context) error
}

// notificationServiceImpl is the concrete implementation of NotificationService.
type notificationServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService() NotificationService {
	return &notificationServiceImpl{}
}
