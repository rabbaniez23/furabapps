// Package handler provides HTTP handlers for menu-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// MenuHandler handles HTTP requests for menu-service.
type MenuHandler struct {
	// TODO: add service dependency
}

// NewMenuHandler creates a new MenuHandler.
func NewMenuHandler() *MenuHandler {
	return &MenuHandler{}
}

// RegisterRoutes registers all menu-service routes.
func (h *MenuHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/menus", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "menu-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
