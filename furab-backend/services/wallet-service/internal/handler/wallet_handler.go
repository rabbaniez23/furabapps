// Package handler provides HTTP handlers for wallet-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// WalletHandler handles HTTP requests for wallet-service.
type WalletHandler struct {
	// TODO: add service dependency
}

// NewWalletHandler creates a new WalletHandler.
func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

// RegisterRoutes registers all wallet-service routes.
func (h *WalletHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/wallets", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "wallet-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
