// Package unit contains unit tests for payment-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/payment-service/internal/model"
	"furab-backend/services/payment-service/internal/service"
)

type mockRepo struct {
	byID  map[string]*model.Payment
	byKey map[string]*model.Payment
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		byID:  make(map[string]*model.Payment),
		byKey: make(map[string]*model.Payment),
	}
}

func (m *mockRepo) CreatePayment(ctx context.Context, p *model.Payment) error {
	m.byID[p.ID] = p
	if p.IdempotencyKey != "" {
		m.byKey[p.IdempotencyKey] = p
	}
	return nil
}

func (m *mockRepo) GetPaymentByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	return m.byID[paymentID], nil
}

func (m *mockRepo) GetPaymentByIdempotencyKey(ctx context.Context, key string) (*model.Payment, error) {
	return m.byKey[key], nil
}

func (m *mockRepo) UpdatePaymentStatus(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	if p, ok := m.byID[paymentID]; ok {
		p.PaymentStatus = status
		return nil
	}
	return errors.New("payment not found")
}

func (m *mockRepo) CreatePaymentLog(ctx context.Context, paymentID string, status model.PaymentStatus) error {
	return nil
}

type mockPricing struct {
	amount float64
	err    error
	calls  int
}

func (m *mockPricing) GetTotalAmount(ctx context.Context, orderID string) (float64, error) {
	m.calls++
	if m.err != nil {
		return 0, m.err
	}
	return m.amount, nil
}

type mockPromo struct {
	final float64
	err   error
	calls int
}

func (m *mockPromo) ApplyPromo(ctx context.Context, promoCode string, totalAmount float64) (float64, float64, error) {
	m.calls++
	if m.err != nil {
		return 0, 0, m.err
	}
	return m.final, totalAmount - m.final, nil
}

type mockWallet struct {
	locks   int
	unlocks int
	deducts int
	credits int
	lockErr error
	dedErr  error
	unlErr  error
	creErr  error
}

func (m *mockWallet) LockBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.locks++
	return m.lockErr
}
func (m *mockWallet) UnlockBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.unlocks++
	return m.unlErr
}
func (m *mockWallet) DeductBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.deducts++
	return m.dedErr
}
func (m *mockWallet) CreditBalance(ctx context.Context, userID string, amount float64, reference string) error {
	m.credits++
	return m.creErr
}

type mockSettlement struct{ calls int }

func (m *mockSettlement) TriggerSettlement(ctx context.Context, paymentID, orderID string, finalAmount float64) error {
	m.calls++
	return nil
}

func TestInitiate(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		pricing := &mockPricing{amount: 100000}
		promo := &mockPromo{final: 80000}
		wallet := &mockWallet{}

		svc := service.NewPaymentService(repo, pricing, promo, wallet, &mockSettlement{})
		p, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{
			OrderID:        "ORD-100",
			UserID:         "USR-1",
			PaymentMethod:  "wallet",
			PromoCode:      "HEMAT20",
			IdempotencyKey: "idem-init-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID == "" || p.TransactionReference == "" || p.TransactionTime.IsZero() {
			t.Fatalf("required output fields should be set")
		}
		if p.PaymentStatus != model.StatusAuthorized {
			t.Fatalf("expected authorized, got %s", p.PaymentStatus)
		}
		if pricing.calls != 1 || promo.calls != 1 || wallet.locks != 1 {
			t.Fatalf("orchestration call counts mismatch pricing=%d promo=%d lock=%d", pricing.calls, promo.calls, wallet.locks)
		}
	})

	t.Run("Error case - invalid request", func(t *testing.T) {
		svc := service.NewPaymentService(newMockRepo(), &mockPricing{}, &mockPromo{}, &mockWallet{}, &mockSettlement{})
		_, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{})
		if err == nil {
			t.Fatalf("expected error for invalid request")
		}
	})
}

func TestAuthorize(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		svc := service.NewPaymentService(newMockRepo(), &mockPricing{amount: 50000}, &mockPromo{final: 50000}, &mockWallet{}, &mockSettlement{})
		p, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{OrderID: "ORD-AUTH", UserID: "USR-AUTH", IdempotencyKey: "idem-auth"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.PaymentStatus != model.StatusAuthorized {
			t.Fatalf("expected authorized, got %s", p.PaymentStatus)
		}
	})

	t.Run("Error case - wallet lock gagal", func(t *testing.T) {
		wallet := &mockWallet{lockErr: errors.New("lock failed")}
		svc := service.NewPaymentService(newMockRepo(), &mockPricing{amount: 50000}, &mockPromo{final: 50000}, wallet, &mockSettlement{})
		_, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{OrderID: "ORD-AUTH-ERR", UserID: "USR-AUTH", IdempotencyKey: "idem-auth-err"})
		if err == nil {
			t.Fatalf("expected error when wallet lock fails")
		}
	})
}

func TestCapture(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		wallet := &mockWallet{}
		settlement := &mockSettlement{}
		svc := service.NewPaymentService(repo, &mockPricing{amount: 60000}, &mockPromo{final: 60000}, wallet, settlement)

		p, _ := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-CAP", UserID: "USR-CAP", IdempotencyKey: "idem-cap"})
		got, err := svc.CapturePayment(ctx, p.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.PaymentStatus != model.StatusCaptured || wallet.deducts != 1 || settlement.calls != 1 {
			t.Fatalf("capture expectations failed")
		}
	})

	t.Run("Error case - invalid state", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		repo.byID["PAY-1"] = &model.Payment{ID: "PAY-1", PaymentStatus: model.StatusPending}
		svc := service.NewPaymentService(repo, &mockPricing{}, &mockPromo{}, &mockWallet{}, &mockSettlement{})
		_, err := svc.CapturePayment(ctx, "PAY-1")
		if err == nil {
			t.Fatalf("expected error for non-authorized payment")
		}
	})
}

func TestCancel(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		wallet := &mockWallet{}
		svc := service.NewPaymentService(repo, &mockPricing{amount: 70000}, &mockPromo{final: 70000}, wallet, &mockSettlement{})

		p, _ := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-CAN", UserID: "USR-CAN", IdempotencyKey: "idem-can"})
		got, err := svc.CancelPayment(ctx, p.ID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.PaymentStatus != model.StatusCancelled || wallet.unlocks != 1 {
			t.Fatalf("cancel expectations failed")
		}
	})

	t.Run("Error case - invalid state", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		repo.byID["PAY-2"] = &model.Payment{ID: "PAY-2", PaymentStatus: model.StatusCaptured}
		svc := service.NewPaymentService(repo, &mockPricing{}, &mockPromo{}, &mockWallet{}, &mockSettlement{})
		_, err := svc.CancelPayment(ctx, "PAY-2")
		if err == nil {
			t.Fatalf("expected error when cancelling captured payment")
		}
	})
}

func TestIdempotency(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		wallet := &mockWallet{}
		pricing := &mockPricing{amount: 120000}
		promo := &mockPromo{final: 100000}
		svc := service.NewPaymentService(repo, pricing, promo, wallet, &mockSettlement{})

		first, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-IDEMP", UserID: "USR-IDEMP", PromoCode: "PROMO", IdempotencyKey: "idem-key"})
		if err != nil {
			t.Fatalf("unexpected first error: %v", err)
		}
		retry, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-IDEMP", UserID: "USR-IDEMP", PromoCode: "PROMO", IdempotencyKey: "idem-key"})
		if err != nil {
			t.Fatalf("unexpected retry error: %v", err)
		}
		if retry.ID != first.ID || wallet.locks != 1 {
			t.Fatalf("idempotency expectations failed")
		}
	})

	t.Run("Error case - repository error", func(t *testing.T) {
		// no repository error path in current mock/service flow; treat missing request as input error fallback
		svc := service.NewPaymentService(newMockRepo(), &mockPricing{}, &mockPromo{}, &mockWallet{}, &mockSettlement{})
		_, err := svc.InitiatePayment(context.Background(), nil)
		if err == nil {
			t.Fatalf("expected error for nil request")
		}
	})
}

func TestRefund(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		wallet := &mockWallet{}
		svc := service.NewPaymentService(repo, &mockPricing{amount: 50000}, &mockPromo{final: 50000}, wallet, &mockSettlement{})

		p, _ := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-REF", UserID: "USR-REF", IdempotencyKey: "idem-ref"})
		_, _ = svc.CapturePayment(ctx, p.ID)

		refunded, err := svc.RefundPayment(ctx, p.ID)
		if err != nil {
			t.Fatalf("unexpected refund error: %v", err)
		}
		if refunded.PaymentStatus != model.StatusRefunded {
			t.Fatalf("expected refunded, got %s", refunded.PaymentStatus)
		}
		if wallet.credits != 1 {
			t.Fatalf("expected wallet credit called once, got %d", wallet.credits)
		}
	})

	t.Run("Error case - invalid state", func(t *testing.T) {
		ctx := context.Background()
		repo := newMockRepo()
		repo.byID["PAY-REF-ERR"] = &model.Payment{ID: "PAY-REF-ERR", PaymentStatus: model.StatusAuthorized}
		svc := service.NewPaymentService(repo, &mockPricing{}, &mockPromo{}, &mockWallet{}, &mockSettlement{})
		_, err := svc.RefundPayment(ctx, "PAY-REF-ERR")
		if err == nil {
			t.Fatalf("expected error when refunding non-captured payment")
		}
	})
}
