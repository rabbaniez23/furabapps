// Package handler provides HTTP handlers for driver-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// DriverHandler handles HTTP requests for driver-service.
type DriverHandler struct {
	// TODO: add service dependency
}

// NewDriverHandler creates a new DriverHandler.
func NewDriverHandler() *DriverHandler {
	return &DriverHandler{}
}

// RegisterRoutes registers all driver-service routes.
func (h *DriverHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/drivers", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "driver-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
