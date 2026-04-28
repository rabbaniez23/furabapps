// Package handler provides HTTP handlers for ride order API endpoints.
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"furab-backend/services/ride-order-service/internal/model"
	"furab-backend/services/ride-order-service/internal/service"
	"furab-backend/shared/utils"

	"github.com/go-chi/chi/v5"
)

// OrderHandler handles HTTP requests for ride order operations.
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler creates a new OrderHandler with the given service.
func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{service: svc}
}

// RegisterRoutes registers all ride order routes on the given chi router.
func (h *OrderHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/rides", func(r chi.Router) {
		r.Post("/", h.CreateOrder)           // POST /api/v1/rides
		r.Get("/{orderID}", h.GetOrder)      // GET  /api/v1/rides/{orderID}
		r.Put("/{orderID}/assign", h.AssignDriver)  // PUT  /api/v1/rides/{orderID}/assign
		r.Put("/{orderID}/start", h.StartRide)      // PUT  /api/v1/rides/{orderID}/start
		r.Put("/{orderID}/complete", h.CompleteRide) // PUT  /api/v1/rides/{orderID}/complete
		r.Put("/{orderID}/cancel", h.CancelRide)     // PUT  /api/v1/rides/{orderID}/cancel
		r.Get("/user/{userID}", h.GetUserOrders)     // GET  /api/v1/rides/user/{userID}
	})
}

// CreateOrder handles POST /api/v1/rides
// Creates a new ride order with pickup and dropoff locations.
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRideOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	order, err := h.service.CreateOrder(r.Context(), &req)
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to create order")
		}
		return
	}

	utils.SuccessMessageResponse(w, http.StatusCreated, "ride order created successfully", model.RideOrderResponse{
		Order:         order,
		EstimatedFare: order.Fare,
	})
}

// GetOrder handles GET /api/v1/rides/{orderID}
// Retrieves a ride order by its ID.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			utils.ErrorResponse(w, http.StatusNotFound, "order not found")
		case service.ErrInvalidRequest:
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to get order")
		}
		return
	}

	utils.SuccessResponse(w, http.StatusOK, order)
}

// AssignDriver handles PUT /api/v1/rides/{orderID}/assign
// Assigns a driver to a pending ride order.
func (h *OrderHandler) AssignDriver(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "order ID is required")
		return
	}

	var req model.AssignDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.AssignDriver(r.Context(), orderID, req.DriverID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			utils.ErrorResponse(w, http.StatusNotFound, "order not found")
		case service.ErrInvalidTransition:
			utils.ErrorResponse(w, http.StatusConflict, "cannot assign driver to order with current status")
		case service.ErrDriverAlreadyAssigned:
			utils.ErrorResponse(w, http.StatusConflict, "driver already assigned to this order")
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to assign driver")
		}
		return
	}

	utils.SuccessMessageResponse(w, http.StatusOK, "driver assigned successfully", order)
}

// StartRide handles PUT /api/v1/rides/{orderID}/start
// Starts an assigned ride.
func (h *OrderHandler) StartRide(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.service.StartRide(r.Context(), orderID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			utils.ErrorResponse(w, http.StatusNotFound, "order not found")
		case service.ErrInvalidTransition:
			utils.ErrorResponse(w, http.StatusConflict, "cannot start ride with current status")
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to start ride")
		}
		return
	}

	utils.SuccessMessageResponse(w, http.StatusOK, "ride started successfully", order)
}

// CompleteRide handles PUT /api/v1/rides/{orderID}/complete
// Completes a started ride.
func (h *OrderHandler) CompleteRide(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.service.CompleteRide(r.Context(), orderID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			utils.ErrorResponse(w, http.StatusNotFound, "order not found")
		case service.ErrInvalidTransition:
			utils.ErrorResponse(w, http.StatusConflict, "cannot complete ride with current status")
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to complete ride")
		}
		return
	}

	utils.SuccessMessageResponse(w, http.StatusOK, "ride completed successfully", order)
}

// CancelRide handles PUT /api/v1/rides/{orderID}/cancel
// Cancels a pending or assigned ride.
func (h *OrderHandler) CancelRide(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderID")
	if orderID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "order ID is required")
		return
	}

	order, err := h.service.CancelRide(r.Context(), orderID)
	if err != nil {
		switch err {
		case service.ErrOrderNotFound:
			utils.ErrorResponse(w, http.StatusNotFound, "order not found")
		case service.ErrInvalidTransition:
			utils.ErrorResponse(w, http.StatusConflict, "cannot cancel ride with current status")
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to cancel ride")
		}
		return
	}

	utils.SuccessMessageResponse(w, http.StatusOK, "ride cancelled successfully", order)
}

// GetUserOrders handles GET /api/v1/rides/user/{userID}
// Retrieves all ride orders for a specific user with pagination.
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	orders, total, err := h.service.GetUserOrders(r.Context(), userID, limit, offset)
	if err != nil {
		switch err {
		case service.ErrInvalidRequest:
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			utils.ErrorResponse(w, http.StatusInternalServerError, "failed to get orders")
		}
		return
	}

	utils.PaginatedSuccessResponse(w, orders, page, limit, total)
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "ride-order-service",
	})
}
