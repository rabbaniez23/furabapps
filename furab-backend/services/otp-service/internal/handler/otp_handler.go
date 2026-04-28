// Package handler provides HTTP handlers for otp-service.
package handler

import (
	"net/http"

	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// OTPHandler handles HTTP requests for otp-service.
type OTPHandler struct {
	// TODO: add service dependency
}

// NewOTPHandler creates a new OTPHandler.
func NewOTPHandler() *OTPHandler {
	return &OTPHandler{}
}

// RegisterRoutes registers all otp-service routes.
func (h *OTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/otps", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, map[string]string{
				"status":  "healthy",
				"service": "otp-service",
			})
		})
		// TODO: Register endpoint routes
	})
}
