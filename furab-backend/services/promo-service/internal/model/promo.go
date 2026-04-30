// Package model defines the domain models for promo-service.
package model

import "time"

// PromoValidationRequest is the input payload for promo validation.
type PromoValidationRequest struct {
	PromoCode   string  `json:"promo_code"`
	UserID      string  `json:"user_id"`
	OrderID     string  `json:"order_id"`
	TotalAmount float64 `json:"total_amount"`
}

// PromoValidationResponse is the response payload for promo validation.
type PromoValidationResponse struct {
	Status         string  `json:"status"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalAmount    float64 `json:"final_amount"`
}

// Promo represents promotion metadata and rules.
type Promo struct {
	PromoID       string    `json:"promo_id"`
	Code          string    `json:"code"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	MinPurchase   float64   `json:"min_purchase"`
	MaxDiscount   float64   `json:"max_discount"`
	ExpiryDate    time.Time `json:"expiry_date"`
	UsageLimit    int       `json:"usage_limit"`
	UsageCount    int       `json:"usage_count"`
}

const (
	PromoStatusValid   = "Valid"
	PromoStatusInvalid = "Invalid"
)

