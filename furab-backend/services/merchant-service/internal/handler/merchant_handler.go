// Package handler provides HTTP handlers for merchant-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// MerchantHandler handles HTTP requests for merchant-service.
type MerchantHandler struct {
	// TODO: add service dependency
}

// NewMerchantHandler creates a new MerchantHandler.
func NewMerchantHandler() *MerchantHandler {
	return &MerchantHandler{}
}

// RegisterRoutes registers all merchant-service routes.
func (h *MerchantHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/merchants", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "merchant-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
