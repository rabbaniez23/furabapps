// Package unit contains unit tests for settlement-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/settlement-service/internal/model"
	"furab-backend/services/settlement-service/internal/service"
)

type mockRepo struct {
	byPayment map[string]*model.Settlement
}

func newMockRepo() *mockRepo {
	return &mockRepo{byPayment: make(map[string]*model.Settlement)}
}

func (m *mockRepo) CreateSettlement(ctx context.Context, s *model.Settlement) error {
	m.byPayment[s.PaymentID] = s
	return nil
}
func (m *mockRepo) GetSettlementByPaymentID(ctx context.Context, paymentID string) (*model.Settlement, error) {
	return m.byPayment[paymentID], nil
}
func (m *mockRepo) UpdateSettlementStatus(ctx context.Context, settlementID string, status model.SettlementStatus) error {
	for _, s := range m.byPayment {
		if s.ID == settlementID {
			s.Status = status
			return nil
		}
	}
	return errors.New("settlement not found")
}

type mockWallet struct {
	fail       bool
	creditCall int
}

func (m *mockWallet) CreditBalance(ctx context.Context, walletID string, amount float64, referenceID string) error {
	m.creditCall++
	if m.fail {
		return errors.New("wallet service error")
	}
	return nil
}

type mockDriver struct{ walletID string }

func (m *mockDriver) GetDriverWalletIDByOrderID(ctx context.Context, orderID string) (string, error) {
	if m.walletID == "" {
		return "", errors.New("driver wallet not found")
	}
	return m.walletID, nil
}

type mockMerchant struct{ walletID string }

func (m *mockMerchant) GetMerchantWalletIDByOrderID(ctx context.Context, orderID string) (string, error) {
	if m.walletID == "" {
		return "", errors.New("merchant wallet not found")
	}
	return m.walletID, nil
}

func TestProcessSettlement(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		wallet := &mockWallet{}
		drv := &mockDriver{walletID: "W-DRV-1"}
		mer := &mockMerchant{walletID: "W-MER-1"}
		svc := service.NewSettlementService(repo, wallet, drv, mer)

		res, err := svc.ProcessSettlement(context.Background(), &model.ProcessSettlementRequest{
			PaymentID:   "PAY-101",
			OrderID:     "ORD-500",
			TotalAmount: 100000,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != "Success" {
			t.Fatalf("expected Success, got %s", res.Status)
		}
		if res.DriverAmount != 80000 || res.MerchantAmount != 15000 || res.PlatformFee != 5000 {
			t.Fatalf("split mismatch: %+v", res)
		}
		if wallet.creditCall != 2 {
			t.Fatalf("expected 2 wallet credits, got %d", wallet.creditCall)
		}
	})

	t.Run("Error case - wallet credit failed", func(t *testing.T) {
		repo := newMockRepo()
		wallet := &mockWallet{fail: true}
		drv := &mockDriver{walletID: "W-DRV-1"}
		mer := &mockMerchant{walletID: "W-MER-1"}
		svc := service.NewSettlementService(repo, wallet, drv, mer)

		res, err := svc.ProcessSettlement(context.Background(), &model.ProcessSettlementRequest{
			PaymentID:   "PAY-102",
			OrderID:     "ORD-501",
			TotalAmount: 50000,
		})
		if err == nil {
			t.Fatalf("expected error")
		}
		if res == nil || res.Status != "Failed" {
			t.Fatalf("expected failed response, got %+v", res)
		}
		stored := repo.byPayment["PAY-102"]
		if stored == nil || stored.Status != model.StatusFailed {
			t.Fatalf("expected stored failed status")
		}
	})

	t.Run("Idempotency - same payment_id not processed twice", func(t *testing.T) {
		repo := newMockRepo()
		repo.byPayment["PAY-200"] = &model.Settlement{
			ID:             "SET-200",
			PaymentID:      "PAY-200",
			OrderID:        "ORD-700",
			TotalAmount:    100000,
			DriverAmount:   80000,
			MerchantAmount: 15000,
			PlatformFee:    5000,
			Status:         model.StatusSuccess,
		}
		wallet := &mockWallet{}
		svc := service.NewSettlementService(repo, wallet, &mockDriver{walletID: "W-DRV"}, &mockMerchant{walletID: "W-MER"})

		res, err := svc.ProcessSettlement(context.Background(), &model.ProcessSettlementRequest{
			PaymentID:   "PAY-200",
			OrderID:     "ORD-700",
			TotalAmount: 100000,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != "Success" {
			t.Fatalf("expected Success, got %s", res.Status)
		}
		if wallet.creditCall != 0 {
			t.Fatalf("idempotent replay must not credit again")
		}
	})
}
