package service

import (
	"context"
	"errors"
	"time"

	"furab-backend/services/audit-log-service/internal/model"
	"furab-backend/services/audit-log-service/internal/repository"
)

// SearchLogFilter represents the available filters for searching audit logs.
type SearchLogFilter struct {
	StartDate   *time.Time
	EndDate     *time.Time
	ServiceName string
	ActorID     string
	Action      string
	TargetID    string
	Status      string
}

// AuditLogService defines the interface for audit log business logic.
type AuditLogService interface {
	RecordLog(ctx context.Context, log model.AuditLog) error
	SearchLog(ctx context.Context, filter SearchLogFilter, page, limit int) ([]model.AuditLog, int, error)
	GetLogDetail(ctx context.Context, logID string) (model.AuditLog, error)
}

// auditLogServiceImpl is the concrete implementation of AuditLogService.
type auditLogServiceImpl struct {
	repo repository.AuditLogRepository
}

// NewAuditLogService creates a new AuditLogService.
func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogServiceImpl{
		repo: repo,
	}
}

// RecordLog validates payload and appends a new log.
func (s *auditLogServiceImpl) RecordLog(ctx context.Context, log model.AuditLog) error {
	// 1. Validasi field wajib
	if log.ServiceName == "" || log.ActorID == "" || log.Action == "" || log.Status == "" {
		return errors.New("INVALID_PAYLOAD")
	}

	// 2. Set Timestamp otomatis jika belum diisi (Zero Value)
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	// 3. Simpan ke repository (append-only)
	return s.repo.Save(ctx, log)
}

// SearchLog fetches a list of logs applying the requested filters.
func (s *auditLogServiceImpl) SearchLog(ctx context.Context, filter SearchLogFilter, page, limit int) ([]model.AuditLog, int, error) {
	// Pagination default fallback
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Konversi struct filter menjadi map interface untuk repository query builder
	repoFilter := make(map[string]interface{})

	if filter.StartDate != nil {
		repoFilter["start_date"] = *filter.StartDate
	}
	if filter.EndDate != nil {
		repoFilter["end_date"] = *filter.EndDate
	}
	if filter.ServiceName != "" {
		repoFilter["service_name"] = filter.ServiceName
	}
	if filter.ActorID != "" {
		repoFilter["actor_id"] = filter.ActorID
	}
	if filter.Action != "" {
		repoFilter["action"] = filter.Action
	}
	if filter.TargetID != "" {
		repoFilter["target_id"] = filter.TargetID
	}
	if filter.Status != "" {
		repoFilter["status"] = filter.Status
	}

	return s.repo.Search(ctx, repoFilter, page, limit)
}

// GetLogDetail retrieves the details of a single audit log by its ID.
func (s *auditLogServiceImpl) GetLogDetail(ctx context.Context, logID string) (model.AuditLog, error) {
	return s.repo.GetByID(ctx, logID)
}
