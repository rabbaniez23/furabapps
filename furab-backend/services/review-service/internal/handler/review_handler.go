// Package handler provides HTTP handlers for review-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// ReviewHandler handles HTTP requests for review-service.
type ReviewHandler struct {
	// TODO: add service dependency
}

// NewReviewHandler creates a new ReviewHandler.
func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

// RegisterRoutes registers all review-service routes.
func (h *ReviewHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/reviews", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "review-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
