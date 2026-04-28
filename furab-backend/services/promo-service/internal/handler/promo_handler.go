// Package handler provides HTTP handlers for promo-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// PromoHandler handles HTTP requests for promo-service.
type PromoHandler struct {
	// TODO: add service dependency
}

// NewPromoHandler creates a new PromoHandler.
func NewPromoHandler() *PromoHandler {
	return &PromoHandler{}
}

// RegisterRoutes registers all promo-service routes.
func (h *PromoHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/promos", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "promo-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
