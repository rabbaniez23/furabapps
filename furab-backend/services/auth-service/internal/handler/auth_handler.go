// Package handler provides HTTP handlers for auth-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// AuthHandler handles HTTP requests for auth-service.
type AuthHandler struct {
	// TODO: add service dependency
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// RegisterRoutes registers all auth-service routes.
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/auths", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "auth-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
