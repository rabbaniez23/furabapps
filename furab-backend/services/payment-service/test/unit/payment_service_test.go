// Package unit contains unit tests for payment-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/payment-service/internal/model"
	"furab-backend/services/payment-service/internal/service"
	"furab-backend/services/payment-service/test/unit/mock"
)

func newTestService(t *testing.T) (
	service.PaymentService,
	*mock.MockPaymentRepository,
	*mock.MockPricingClient,
	*mock.MockPromoClient,
	*mock.MockWalletClient,
	*mock.MockSettlementClient,
) {
	t.Helper()

	repo := mock.NewMockPaymentRepository()
	pricing := &mock.MockPricingClient{Amount: 100000}
	promo := &mock.MockPromoClient{FinalAmount: 80000}
	wallet := &mock.MockWalletClient{}
	settlement := &mock.MockSettlementClient{}

	svc := service.NewPaymentService(repo, pricing, promo, wallet, settlement)
	return svc, repo, pricing, promo, wallet, settlement
}

func validInitiateRequest() *model.InitiatePaymentRequest {
	return &model.InitiatePaymentRequest{
		OrderID:        "ORD-100",
		UserID:         "USR-1",
		PaymentMethod:  "wallet",
		PaymentDetail:  "wallet-topup",
		PromoCode:      "HEMAT20",
		IdempotencyKey: "idem-init-1",
	}
}

func TestInitiatePayment_Success(t *testing.T) {
	svc, _, pricing, promo, wallet, _ := newTestService(t)

	p, err := svc.InitiatePayment(context.Background(), validInitiateRequest())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID == "" || p.TransactionReference == "" || p.TransactionTime.IsZero() {
		t.Fatalf("required output fields should be set")
	}
	if p.PaymentStatus != model.StatusAuthorized {
		t.Fatalf("expected authorized, got %s", p.PaymentStatus)
	}
	if p.MethodID != "wallet" || p.PaymentDetail != "wallet-topup" {
		t.Fatalf("expected payment method and detail to be preserved, got %s / %s", p.MethodID, p.PaymentDetail)
	}
	if p.TransactionReference != "TXN-ORD-100" {
		t.Fatalf("expected transaction reference TXN-ORD-100, got %s", p.TransactionReference)
	}
	if pricing.Calls != 1 || promo.Calls != 1 || wallet.Locks != 1 {
		t.Fatalf("orchestration call counts mismatch pricing=%d promo=%d lock=%d", pricing.Calls, promo.Calls, wallet.Locks)
	}
}

func TestInitiatePayment_InvalidRequest(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)

	_, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request")
	}
}

func TestInitiatePayment_MissingPaymentMethod(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)

	_, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{OrderID: "ORD-ERR", UserID: "USR-ERR", PaymentDetail: "detail-only"})
	if err == nil {
		t.Fatal("expected error for missing payment method")
	}
}

func TestInitiatePayment_MissingPaymentDetail(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)

	_, err := svc.InitiatePayment(context.Background(), &model.InitiatePaymentRequest{OrderID: "ORD-ERR", UserID: "USR-ERR", PaymentMethod: "wallet"})
	if err == nil {
		t.Fatal("expected error for missing payment detail")
	}
}

func TestInitiatePayment_Idempotency(t *testing.T) {
	ctx := context.Background()
	svc, repo, _, _, wallet, _ := newTestService(t)

	// First call - create payment
	req := &model.InitiatePaymentRequest{
		OrderID:        "ORD-IDEM",
		UserID:         "USR-IDEM",
		PaymentMethod:  "wallet",
		PaymentDetail:  "idem-detail",
		IdempotencyKey: "unique-idem-key-123",
	}
	first, err := svc.InitiatePayment(ctx, req)
	if err != nil {
		t.Fatalf("first initiate failed: %v", err)
	}

	initialWalletLocks := wallet.Locks
	initialPaymentCount := len(repo.Payments)

	// Second call with same idempotency key - should return same payment WITHOUT calling wallet service again
	second, err := svc.InitiatePayment(ctx, req)
	if err != nil {
		t.Fatalf("second initiate failed: %v", err)
	}

	// Verify both calls return same payment
	if first.ID != second.ID {
		t.Fatalf("idempotency key should return same payment: first=%s, second=%s", first.ID, second.ID)
	}

	// Verify no duplicate payment was created
	if len(repo.Payments) != initialPaymentCount {
		t.Fatalf("expected no new payment created on second call: before=%d, after=%d", initialPaymentCount, len(repo.Payments))
	}

	// Verify wallet.LockBalance was NOT called again (idempotency should skip it)
	// Should remain at initialWalletLocks, not increase to initialWalletLocks+1
	if wallet.Locks != initialWalletLocks {
		t.Fatalf("idempotency should NOT call wallet lock again: expected %d, got %d", initialWalletLocks, wallet.Locks)
	}

	// Verify payment data is consistent
	if second.PaymentStatus != model.StatusAuthorized || second.FinalAmount == 0 {
		t.Fatalf("expected authorized payment with valid amount, got status=%s amount=%f", second.PaymentStatus, second.FinalAmount)
	}

	// Verify transaction reference and other details are preserved
	if second.TransactionReference != first.TransactionReference {
		t.Fatalf("transaction reference should be same for idempotent calls: first=%s, second=%s", first.TransactionReference, second.TransactionReference)
	}
}


func TestCapturePayment_Success(t *testing.T) {
	ctx := context.Background()
	svc, repo, _, _, wallet, settlement := newTestService(t)

	p, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-CAP", UserID: "USR-CAP", PaymentMethod: "wallet", PaymentDetail: "cap-detail", IdempotencyKey: "idem-cap"})
	if err != nil {
		t.Fatalf("unexpected initiate error: %v", err)
	}

	got, err := svc.CapturePayment(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected capture error: %v", err)
	}
	if got.PaymentStatus != model.StatusCaptured {
		t.Fatalf("expected captured status, got %s", got.PaymentStatus)
	}
	if wallet.Deducts != 1 {
		t.Fatalf("expected wallet deduct once, got %d", wallet.Deducts)
	}
	if settlement.Calls != 1 {
		t.Fatalf("expected settlement call once, got %d", settlement.Calls)
	}
	if repo.LastUpdatedID != p.ID {
		t.Fatalf("expected repo update for payment %s, got %s", p.ID, repo.LastUpdatedID)
	}
}

func TestCapturePayment_InvalidState(t *testing.T) {
	ctx := context.Background()
	svc, repo, _, _, _, _ := newTestService(t)
	repo.Payments["PAY-1"] = &model.Payment{ID: "PAY-1", PaymentStatus: model.StatusPending}

	_, err := svc.CapturePayment(ctx, "PAY-1")
	if err == nil {
		t.Fatal("expected error for non-authorized payment")
	}
}

func TestCancelPayment_Success(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, wallet, _ := newTestService(t)

	p, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-CAN", UserID: "USR-CAN", PaymentMethod: "wallet", PaymentDetail: "cancel-detail", IdempotencyKey: "idem-can"})
	if err != nil {
		t.Fatalf("unexpected initiate error: %v", err)
	}

	got, err := svc.CancelPayment(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected cancel error: %v", err)
	}
	if got.PaymentStatus != model.StatusCancelled {
		t.Fatalf("expected cancelled status, got %s", got.PaymentStatus)
	}
	if wallet.Unlocks != 1 {
		t.Fatalf("expected wallet unlock once, got %d", wallet.Unlocks)
	}
}

func TestCancelPayment_InvalidState(t *testing.T) {
	ctx := context.Background()
	svc, repo, _, _, _, _ := newTestService(t)
	repo.Payments["PAY-2"] = &model.Payment{ID: "PAY-2", PaymentStatus: model.StatusCaptured}

	_, err := svc.CancelPayment(ctx, "PAY-2")
	if err == nil {
		t.Fatal("expected error when cancelling captured payment")
	}
}

func TestRefundPayment_Success(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, wallet, _ := newTestService(t)

	p, err := svc.InitiatePayment(ctx, &model.InitiatePaymentRequest{OrderID: "ORD-REF", UserID: "USR-REF", PaymentMethod: "wallet", PaymentDetail: "refund-detail", IdempotencyKey: "idem-ref"})
	if err != nil {
		t.Fatalf("unexpected initiate error: %v", err)
	}
	_, err = svc.CapturePayment(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected capture error: %v", err)
	}

	refunded, err := svc.RefundPayment(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected refund error: %v", err)
	}
	if refunded.PaymentStatus != model.StatusRefunded {
		t.Fatalf("expected refunded, got %s", refunded.PaymentStatus)
	}
	if wallet.Credits != 1 {
		t.Fatalf("expected wallet credit once, got %d", wallet.Credits)
	}
}

func TestRefundPayment_InvalidState(t *testing.T) {
	ctx := context.Background()
	svc, repo, _, _, _, _ := newTestService(t)
	repo.Payments["PAY-REF-ERR"] = &model.Payment{ID: "PAY-REF-ERR", PaymentStatus: model.StatusAuthorized}

	_, err := svc.RefundPayment(ctx, "PAY-REF-ERR")
	if err == nil {
		t.Fatal("expected error when refunding non-captured payment")
	}
}

func TestGetPayment_Success(t *testing.T) {
	svc, repo, _, _, _, _ := newTestService(t)
	repo.Payments["PAY-GET"] = &model.Payment{ID: "PAY-GET", PaymentStatus: model.StatusAuthorized}

	payment, err := svc.GetPayment(context.Background(), "PAY-GET")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if payment == nil || payment.ID != "PAY-GET" {
		t.Fatalf("expected payment PAY-GET, got %v", payment)
	}
}

func TestGetPayment_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService(t)

	_, err := svc.GetPayment(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing payment")
	}
	if !errors.Is(err, service.ErrPaymentNotFound) {
		t.Fatalf("expected ErrPaymentNotFound, got %v", err)
	}
}
