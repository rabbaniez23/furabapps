// Package handler provides HTTP handlers for user-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// UserHandler handles HTTP requests for user-service.
type UserHandler struct {
	// TODO: add service dependency
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// RegisterRoutes registers all user-service routes.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "user-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
