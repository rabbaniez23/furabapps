// Package handler provides HTTP handlers for audit-log-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// AuditLogHandler handles HTTP requests for audit-log-service.
type AuditLogHandler struct {
	// TODO: add service dependency
}

// NewAuditLogHandler creates a new AuditLogHandler.
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{}
}

// RegisterRoutes registers all audit-log-service routes.
func (h *AuditLogHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auditlogs", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "audit-log-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
