// Package handler provides HTTP handlers for rating-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// RatingHandler handles HTTP requests for rating-service.
type RatingHandler struct {
	// TODO: add service dependency
}

// NewRatingHandler creates a new RatingHandler.
func NewRatingHandler() *RatingHandler {
	return &RatingHandler{}
}

// RegisterRoutes registers all rating-service routes.
func (h *RatingHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/ratings", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "rating-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
