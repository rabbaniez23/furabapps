// Package mock provides mock implementations for payment-service testing.
package mock

import (
	"context"
	"errors"

	"furab-backend/services/payment-service/internal/model"
)

type MockPaymentRepository struct {
	Payments      map[string]*model.Payment
	PaymentsByKey map[string]*model.Payment
	CreateErr     error
	GetByIDErr    error
	GetByKeyErr   error
	UpdateErr     error
	LogErr        error
	LastUpdatedID string
}

func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{
		Payments:      make(map[string]*model.Payment),
		PaymentsByKey: make(map[string]*model.Payment),
	}
}

func (m *MockPaymentRepository) CreatePayment(ctx context.Context, p *model.Payment) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	m.Payments[p.ID] = p
	if p.IdempotencyKey != "" {
		m.PaymentsByKey[p.IdempotencyKey] = p
	}
	return nil
}

func (m *MockPaymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	return m.Payments[paymentID], nil
}

func (m *MockPaymentRepository) GetPaymentByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error) {
	if m.GetByKeyErr != nil {
		return nil, m.GetByKeyErr
	}
	return m.PaymentsByKey[key], nil
}

func (m *MockPaymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	p, ok := m.Payments[paymentID]
	if !ok {
		return errors.New("payment not found")
	}
	p.PaymentStatus = status
	m.LastUpdatedID = paymentID
	return nil
}

func (m *MockPaymentRepository) CreatePaymentLog(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	if m.LogErr != nil {
		return m.LogErr
	}
	return nil
}

type MockPricingClient struct {
	Amount float64
	Err    error
	Calls  int
}

func (m *MockPricingClient) GetTotalAmount(ctx context.Context, orderID string) (float64, error) {
	m.Calls++
	if m.Err != nil {
		return 0, m.Err
	}
	return m.Amount, nil
}

type MockPromoClient struct {
	FinalAmount float64
	Err         error
	Calls       int
}

func (m *MockPromoClient) ApplyPromo(ctx context.Context, promoCode string, totalAmount float64) (float64, float64, error) {
	m.Calls++
	if m.Err != nil {
		return 0, 0, m.Err
	}
	return m.FinalAmount, totalAmount - m.FinalAmount, nil
}

type MockWalletClient struct {
	Locks   int
	Unlocks int
	Deducts int
	Credits int
	LockErr error
	UnlErr  error
	DedErr  error
	CreErr  error
}

func (m *MockWalletClient) LockBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.Locks++
	return m.LockErr
}

func (m *MockWalletClient) UnlockBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.Unlocks++
	return m.UnlErr
}

func (m *MockWalletClient) DeductBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.Deducts++
	return m.DedErr
}

func (m *MockWalletClient) CreditBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.Credits++
	return m.CreErr
}

type MockSettlementClient struct {
	Calls int
	Err   error
}

func (m *MockSettlementClient) TriggerSettlement(ctx context.Context, paymentID, orderID string, finalAmount float64) error {
	m.Calls++
	return m.Err
}
