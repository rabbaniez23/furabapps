// Package service implements the business logic for promo-service.
package service

import (
	"context"
	"fmt"
	"time"

	"furab-backend/services/promo-service/internal/client"
	"furab-backend/services/promo-service/internal/model"
	"furab-backend/services/promo-service/internal/repository"
)

var (
	ErrPromoCodeRequired = fmt.Errorf("promo code is required")
)

// PromoService defines the interface for promo validation and application.
type PromoService interface {
	ValidatePromo(ctx context.Context, promoCode, userID, orderID string, totalAmount float64) (*model.PromoValidationResponse, error)
}

// promoServiceImpl is the concrete implementation of PromoService.
type promoServiceImpl struct {
	repo        repository.PromoRepository
	orderClient client.OrderClient
	userClient  client.UserClient
}

// NewPromoService creates a new PromoService.
func NewPromoService(repo repository.PromoRepository, orderClient client.OrderClient, userClient client.UserClient) PromoService {
	return &promoServiceImpl{
		repo:        repo,
		orderClient: orderClient,
		userClient:  userClient,
	}
}

func (s *promoServiceImpl) ValidatePromo(ctx context.Context, promoCode, userID, orderID string, totalAmount float64) (*model.PromoValidationResponse, error) {
	response := &model.PromoValidationResponse{
		Status:         model.PromoStatusInvalid,
		DiscountAmount: 0,
		FinalAmount:    totalAmount,
	}

	if promoCode == "" {
		return response, nil
	}

	promo, err := s.repo.GetPromoByCode(ctx, promoCode)
	if err != nil {
		return response, nil
	}

	if time.Now().After(promo.ExpiryDate) {
		return response, nil
	}

	if promo.UsageLimit > 0 && promo.UsageCount >= promo.UsageLimit {
		return response, nil
	}

	if totalAmount < promo.MinPurchase {
		return response, nil
	}

	orderValid, err := s.orderClient.ValidateOrderPromo(ctx, orderID, promoCode)
	if err != nil || !orderValid {
		return response, nil
	}

	userValid, err := s.userClient.ValidateUserPromo(ctx, userID, promoCode)
	if err != nil || !userValid {
		return response, nil
	}

	discountAmount := calculateDiscount(totalAmount, promo)
	finalAmount := totalAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	if err := s.repo.IncrementUsage(ctx, promo.PromoID); err != nil {
		return nil, err
	}

	response.Status = model.PromoStatusValid
	response.DiscountAmount = discountAmount
	response.FinalAmount = finalAmount
	return response, nil
}

func calculateDiscount(totalAmount float64, promo *model.Promo) float64 {
	switch promo.DiscountType {
	case "percentage":
		discount := totalAmount * promo.DiscountValue
		if promo.MaxDiscount > 0 && discount > promo.MaxDiscount {
			discount = promo.MaxDiscount
		}
		return discount
	case "fixed":
		return promo.DiscountValue
	default:
		return 0
	}
}
