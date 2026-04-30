// Package repository provides data access layer for pricing-service.
package repository

import (
	"context"
	"errors"

	"furab-backend/services/pricing-service/internal/model"
)

var (
	ErrPriceRuleNotFound = errors.New("price rule not found")
)

// PriceRepository defines the interface for pricing-service pricing rule access.
type PriceRepository interface {
	GetPricingRules(ctx context.Context) ([]model.PriceRule, error)
	GetPricingRuleByType(ctx context.Context, ruleType string) (*model.PriceRule, error)
}

type inMemoryPriceRepository struct {
	rules []model.PriceRule
}

// NewInMemoryPriceRepository creates a pricing repository with fixed rules.
func NewInMemoryPriceRepository() PriceRepository {
	return &inMemoryPriceRepository{
		rules: []model.PriceRule{
			{RuleID: "delivery-per-km", Type: "delivery", Value: 5000, Description: "Delivery fee per kilometer"},
			{RuleID: "service-percent", Type: "service", Value: 0.05, Description: "Service fee percentage"},
			{RuleID: "tax-percent", Type: "tax", Value: 0.0, Description: "Tax percentage"},
		},
	}
}

func (r *inMemoryPriceRepository) GetPricingRules(ctx context.Context) ([]model.PriceRule, error) {
	return r.rules, nil
}

func (r *inMemoryPriceRepository) GetPricingRuleByType(ctx context.Context, ruleType string) (*model.PriceRule, error) {
	for _, rule := range r.rules {
		if rule.Type == ruleType {
			copy := rule
			return &copy, nil
		}
	}

	return nil, ErrPriceRuleNotFound
}

type postgresPriceRepository struct {
}

// NewPostgresPriceRepository creates a new PostgreSQL-based repository.
func NewPostgresPriceRepository() PriceRepository {
	return &postgresPriceRepository{}
}

func (*postgresPriceRepository) GetPricingRules(ctx context.Context) ([]model.PriceRule, error) {
	return nil, errors.New("postgres price repository not implemented")
}

func (*postgresPriceRepository) GetPricingRuleByType(ctx context.Context, ruleType string) (*model.PriceRule, error) {
	return nil, errors.New("postgres price repository not implemented")
}
