// Package handler provides HTTP handlers for pricing-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// PriceHandler handles HTTP requests for pricing-service.
type PriceHandler struct {
	// TODO: add service dependency
}

// NewPriceHandler creates a new PriceHandler.
func NewPriceHandler() *PriceHandler {
	return &PriceHandler{}
}

// RegisterRoutes registers all pricing-service routes.
func (h *PriceHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/prices", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "pricing-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
