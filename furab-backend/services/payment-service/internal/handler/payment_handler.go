// Package handler provides HTTP handlers for payment-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// PaymentHandler handles HTTP requests for payment-service.
type PaymentHandler struct {
	// TODO: add service dependency
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// RegisterRoutes registers all payment-service routes.
func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/payments", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "payment-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
