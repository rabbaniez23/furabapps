// Package unit contains unit tests for wallet-service.
// Unit tests do NOT access any database or external service.
package unit

import (
	"context"
	"errors"
	"testing"

	"furab-backend/services/wallet-service/internal/model"
	"furab-backend/services/wallet-service/internal/service"
)

type mockRepo struct {
	walletByUser map[string]*model.Wallet
	txByRefType  map[string]*model.Transaction
	failUpdate   bool
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		walletByUser: make(map[string]*model.Wallet),
		txByRefType:  make(map[string]*model.Transaction),
	}
}

func key(ref string, typ model.TransactionType) string {
	return string(typ) + "::" + ref
}

func (m *mockRepo) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	w, ok := m.walletByUser[userID]
	if !ok {
		return nil, errors.New("wallet not found")
	}
	return w, nil
}

func (m *mockRepo) UpdateBalance(ctx context.Context, walletID string, newBalance float64) error {
	if m.failUpdate {
		return errors.New("update failed")
	}
	for _, w := range m.walletByUser {
		if w.ID == walletID {
			w.Balance = newBalance
			return nil
		}
	}
	return errors.New("wallet not found")
}

func (m *mockRepo) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	if tx.ReferenceID != "" {
		m.txByRefType[key(tx.ReferenceID, tx.Type)] = tx
	}
	return nil
}

func (m *mockRepo) GetTransactionByReference(ctx context.Context, referenceID string, typ model.TransactionType) (*model.Transaction, error) {
	if tx, ok := m.txByRefType[key(referenceID, typ)]; ok {
		return tx, nil
	}
	return nil, nil
}

func TestHold(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 100000}
		svc := service.NewWalletService(repo)

		res, err := svc.HoldBalance(context.Background(), "U1", 30000, "REF-H1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != model.StatusSuccess || res.CurrentBalance != 70000 {
			t.Fatalf("unexpected hold result: %+v", res)
		}
	})

	t.Run("Error case - insufficient balance", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 10000}
		svc := service.NewWalletService(repo)

		_, err := svc.HoldBalance(context.Background(), "U1", 20000, "REF-H2")
		if err == nil {
			t.Fatalf("expected insufficient balance error")
		}
	})
}

func TestRelease(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 70000}
		svc := service.NewWalletService(repo)
		res, err := svc.ReleaseBalance(context.Background(), "U1", 30000, "REF-R1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.CurrentBalance != 100000 {
			t.Fatalf("expected balance 100000, got %v", res.CurrentBalance)
		}
	})
}

func TestDebit(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 100000}
		svc := service.NewWalletService(repo)
		res, err := svc.DebitBalance(context.Background(), "U1", 20000, "REF-D1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.CurrentBalance != 80000 {
			t.Fatalf("expected balance 80000, got %v", res.CurrentBalance)
		}
	})
}

func TestCredit(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 80000}
		svc := service.NewWalletService(repo)
		res, err := svc.CreditBalance(context.Background(), "U1", 20000, "REF-C1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.CurrentBalance != 100000 {
			t.Fatalf("expected balance 100000, got %v", res.CurrentBalance)
		}
	})
}

func TestRefund(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 50000}
		svc := service.NewWalletService(repo)
		res, err := svc.Refund(context.Background(), "U1", 10000, "REF-RF1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.Status != model.StatusSuccess || res.CurrentBalance != 60000 {
			t.Fatalf("unexpected refund result: %+v", res)
		}
	})
}

func TestIdempotency(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		repo := newMockRepo()
		repo.walletByUser["U1"] = &model.Wallet{ID: "W1", UserID: "U1", Balance: 100000}
		svc := service.NewWalletService(repo)

		first, err := svc.DebitBalance(context.Background(), "U1", 10000, "REF-IDEMP")
		if err != nil {
			t.Fatalf("unexpected first error: %v", err)
		}
		second, err := svc.DebitBalance(context.Background(), "U1", 10000, "REF-IDEMP")
		if err != nil {
			t.Fatalf("unexpected second error: %v", err)
		}
		if second.TransactionID != first.TransactionID {
			t.Fatalf("expected same transaction on idempotent retry")
		}
		if repo.walletByUser["U1"].Balance != 90000 {
			t.Fatalf("balance should not be double deducted, got %v", repo.walletByUser["U1"].Balance)
		}
	})
}
