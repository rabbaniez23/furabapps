// Package handler provides HTTP handlers for location-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// LocationHandler handles HTTP requests for location-service.
type LocationHandler struct {
	// TODO: add service dependency
}

// NewLocationHandler creates a new LocationHandler.
func NewLocationHandler() *LocationHandler {
	return &LocationHandler{}
}

// RegisterRoutes registers all location-service routes.
func (h *LocationHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/locations", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "location-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
