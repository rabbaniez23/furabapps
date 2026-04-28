// Package service implements the business logic for email-service.
package service

import "context"

// EmailService defines the interface for email-service business logic.
type EmailService interface {

	// SendEmail implements the business logic for SendEmail.
	SendEmail(ctx context.Context) error

	// SendBulk implements the business logic for SendBulk.
	SendBulk(ctx context.Context) error

	// GetStatus implements the business logic for GetStatus.
	GetStatus(ctx context.Context) error
}

// emailServiceImpl is the concrete implementation of EmailService.
type emailServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewEmailService creates a new EmailService.
func NewEmailService() EmailService {
	return &emailServiceImpl{}
}
