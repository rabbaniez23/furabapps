// Package unit contains unit tests for email-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"testing"

	"furab-backend/services/email-service/internal/model"
	"furab-backend/services/email-service/internal/service"
	"furab-backend/services/email-service/test/unit/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewEmailService_Creation tests that the service can be created.
func TestNewEmailService_Creation(t *testing.T) {
	repo := &mock.MockEmailRepository{}
	sender := &mock.MockEmailSender{}
	svc := service.NewEmailService(repo, sender)
	require.NotNil(t, svc)
}

func TestSendEmail(t *testing.T) {
	t.Run("Success - Email berhasil dikirim", func(t *testing.T) {
		repo := &mock.MockEmailRepository{}
		sender := &mock.MockEmailSender{}
		svc := service.NewEmailService(repo, sender)

		res, err := svc.SendEmail(context.Background(), model.SendEmailRequest{
			ReceiverEmail: "user@mail.com",
			Subject:       "Invoice",
			Body:          "Detail transaksi",
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.Equal(t, "email berhasil dikirim", res.Message)
		assert.Equal(t, 1, sender.SendCall)
		assert.Equal(t, 1, repo.SaveEmailLogCall)
		assert.Equal(t, "user@mail.com", repo.LastSavedLog.ReceiverEmail)
		assert.Equal(t, "Invoice", repo.LastSavedLog.Subject)
		assert.Equal(t, "sent", repo.LastSavedLog.Status)
	})

	t.Run("Error - Email kosong", func(t *testing.T) {
		repo := &mock.MockEmailRepository{}
		sender := &mock.MockEmailSender{}
		svc := service.NewEmailService(repo, sender)

		res, err := svc.SendEmail(context.Background(), model.SendEmailRequest{
			ReceiverEmail: "",
			Subject:       "Invoice",
			Body:          "Detail transaksi",
		})

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "email required", err.Error())
		assert.Equal(t, 0, sender.SendCall)
		assert.Equal(t, 0, repo.SaveEmailLogCall)
	})

	t.Run("Error - Subject kosong", func(t *testing.T) {
		repo := &mock.MockEmailRepository{}
		sender := &mock.MockEmailSender{}
		svc := service.NewEmailService(repo, sender)

		res, err := svc.SendEmail(context.Background(), model.SendEmailRequest{
			ReceiverEmail: "user@mail.com",
			Subject:       "",
			Body:          "Detail transaksi",
		})

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "subject required", err.Error())
		assert.Equal(t, 0, sender.SendCall)
		assert.Equal(t, 0, repo.SaveEmailLogCall)
	})
}

func TestTriggerEventEmail(t *testing.T) {
	t.Run("Success - Event diproses", func(t *testing.T) {
		repo := &mock.MockEmailRepository{}
		sender := &mock.MockEmailSender{}
		svc := service.NewEmailService(repo, sender)

		res, err := svc.TriggerEventEmail(context.Background(), model.EventEmailRequest{
			EventType:     "payment.success",
			ReceiverEmail: "user@mail.com",
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "success", res.Status)
		assert.Equal(t, 1, sender.SendCall)
		assert.Equal(t, "user@mail.com", sender.LastEmail)
		assert.Equal(t, "Invoice", sender.LastSubj)
		assert.Equal(t, "Detail transaksi", sender.LastBody)
		assert.Equal(t, 1, repo.SaveEmailLogCall)
	})

	t.Run("Error - Event tidak valid", func(t *testing.T) {
		repo := &mock.MockEmailRepository{}
		sender := &mock.MockEmailSender{}
		svc := service.NewEmailService(repo, sender)

		res, err := svc.TriggerEventEmail(context.Background(), model.EventEmailRequest{
			EventType:     "invalid.event",
			ReceiverEmail: "user@mail.com",
		})

		require.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "invalid event", err.Error())
		assert.Equal(t, 0, sender.SendCall)
		assert.Equal(t, 0, repo.SaveEmailLogCall)
	})
}
