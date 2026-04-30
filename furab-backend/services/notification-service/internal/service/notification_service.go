// Package service implements the business logic for notification-service.
package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/notification-service/internal/model"
	"furab-backend/services/notification-service/internal/repository"
)

// EmailClient defines the interface for calling the Email Service.
type EmailClient interface {
	SendEmail(ctx context.Context, receiverID string, title string, message string) error
}

// NotificationService defines the interface for notification-service business logic.
type NotificationService interface {
	ProcessEventNotification(ctx context.Context, req model.EventNotificationRequest) (*model.NotificationResponse, error)
	GenerateNotificationTemplate(ctx context.Context, eventType string) (*model.NotifTemplate, error)
}

// notificationServiceImpl is the concrete implementation of NotificationService.
type notificationServiceImpl struct {
	repo       repository.NotificationRepository
	emailClient EmailClient
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo repository.NotificationRepository, emailClient EmailClient) NotificationService {
	return &notificationServiceImpl{
		repo:       repo,
		emailClient: emailClient,
	}
}

func (s *notificationServiceImpl) ProcessEventNotification(ctx context.Context, req model.EventNotificationRequest) (*model.NotificationResponse, error) {
	if req.EventType == "" {
		return nil, errors.New("event type is required")
	}

	if req.Channel != "push" && req.Channel != "email" {
		return nil, errors.New("invalid channel")
	}

	template, err := s.GenerateNotificationTemplate(ctx, req.EventType)
	if err != nil {
		return nil, err
	}

	log := model.NotificationLog{
		ReceiverID:  req.ReceiverID,
		Title:       template.TitleTemplate,
		Message:     template.MessageTemplate,
		Channel:     req.Channel,
		ReferenceID: req.ReferenceID,
		Timestamp:   time.Now(),
		Status:      "sent",
	}

	if err := s.repo.SaveNotificationLog(ctx, log); err != nil {
		return nil, err
	}

	if req.Channel == "email" {
		if err := s.emailClient.SendEmail(ctx, req.ReceiverID, template.TitleTemplate, template.MessageTemplate); err != nil {
			return nil, err
		}
	}

	return &model.NotificationResponse{
		Status:  "success",
		Message: "notifikasi berhasil",
	}, nil
}

func (s *notificationServiceImpl) GenerateNotificationTemplate(ctx context.Context, eventType string) (*model.NotifTemplate, error) {
	template, err := s.repo.GetTemplateByEventType(ctx, eventType)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.New("template not found")
	}
	return template, nil
}
