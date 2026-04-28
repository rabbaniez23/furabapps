// Package handler provides HTTP handlers for settlement-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// SettlementHandler handles HTTP requests for settlement-service.
type SettlementHandler struct {
	// TODO: add service dependency
}

// NewSettlementHandler creates a new SettlementHandler.
func NewSettlementHandler() *SettlementHandler {
	return &SettlementHandler{}
}

// RegisterRoutes registers all settlement-service routes.
func (h *SettlementHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/settlements", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "settlement-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
