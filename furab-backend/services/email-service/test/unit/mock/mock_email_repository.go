// Package mock provides mock implementations for email-service testing.
package mock

import (
	"context"

	"furab-backend/services/email-service/internal/model"
)

// MockEmailRepository is a lightweight test double for EmailRepository.
type MockEmailRepository struct {
	SaveEmailLogFn   func(ctx context.Context, log model.EmailLog) error
	SaveEmailLogCall int
	LastSavedLog     model.EmailLog
}

// SaveEmailLog records calls and delegates behavior to SaveEmailLogFn.
func (m *MockEmailRepository) SaveEmailLog(ctx context.Context, log model.EmailLog) error {
	m.SaveEmailLogCall++
	m.LastSavedLog = log
	if m.SaveEmailLogFn != nil {
		return m.SaveEmailLogFn(ctx, log)
	}
	return nil
}
