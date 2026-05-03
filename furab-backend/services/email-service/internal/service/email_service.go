// Package service implements the business logic for email-service.
package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/email-service/internal/model"
	"furab-backend/services/email-service/internal/repository"
)

var eventTemplates = map[string]struct {
	subject string
	body    string
}{
	"payment.success": {subject: "Invoice", body: "Detail transaksi"},
}

// EmailSender defines outbound email provider client behavior.
type EmailSender interface {
	Send(ctx context.Context, receiverEmail, subject, body string) error
}

// EmailService defines the interface for email-service business logic.
type EmailService interface {
	SendEmail(ctx context.Context, req model.SendEmailRequest) (*model.EmailResponse, error)
	TriggerEventEmail(ctx context.Context, req model.EventEmailRequest) (*model.EmailResponse, error)
}

// emailServiceImpl is the concrete implementation of EmailService.
type emailServiceImpl struct {
	repo   repository.EmailRepository
	sender EmailSender
}

// NewEmailService creates a new EmailService.
func NewEmailService(repo repository.EmailRepository, sender EmailSender) EmailService {
	return &emailServiceImpl{
		repo:   repo,
		sender: sender,
	}
}

func (s *emailServiceImpl) SendEmail(ctx context.Context, req model.SendEmailRequest) (*model.EmailResponse, error) {
	if req.ReceiverEmail == "" {
		return nil, errors.New("email required")
	}
	if req.Subject == "" {
		return nil, errors.New("subject required")
	}

	if err := s.sender.Send(ctx, req.ReceiverEmail, req.Subject, req.Body); err != nil {
		return nil, err
	}

	if err := s.repo.SaveEmailLog(ctx, model.EmailLog{
		ReceiverEmail: req.ReceiverEmail,
		Subject:       req.Subject,
		Status:        "sent",
		Timestamp:     time.Now(),
		ReceiverID:    req.ReceiverID,
		ReferenceID:   req.ReferenceID,
	}); err != nil {
		return nil, err
	}

	return &model.EmailResponse{
		Status:  "success",
		Message: "email berhasil dikirim",
	}, nil
}

func (s *emailServiceImpl) TriggerEventEmail(ctx context.Context, req model.EventEmailRequest) (*model.EmailResponse, error) {
	template, ok := eventTemplates[req.EventType]
	if !ok {
		return nil, errors.New("invalid event")
	}

	return s.SendEmail(ctx, model.SendEmailRequest{
		ReceiverEmail: req.ReceiverEmail,
		Subject:       template.subject,
		Body:          template.body,
		ReceiverID:    req.ReceiverID,
		ReferenceID:   req.ReferenceID,
	})
}
