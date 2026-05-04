//go:build functional
// +build functional

// Package functional contains functional tests for wallet-service.
// Functional tests access a real PostgreSQL database.
// Run with: go test ./test/functional/... -v -tags=functional
package functional

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"furab-backend/services/wallet-service/internal/repository"
	"furab-backend/services/wallet-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB   *sql.DB
	testRepo repository.WalletRepository
	testSvc  service.WalletService
)

// TestMain sets up and tears down test infrastructure.
func TestMain(m *testing.M) {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "furab")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "furab_secret")
	dbName := getEnvOrDefault("DB_NAME", "wallet_service_test")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("functional tests skipped: failed to connect database: %v", err)
		os.Exit(0)
	}
	defer testDB.Close()

	for i := 0; i < 10; i++ {
		if err = testDB.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Printf("functional tests skipped: database not ready: %v", err)
		os.Exit(0)
	}

	testRepo = repository.NewPostgresWalletRepository(testDB)
	testSvc = service.NewWalletService(testRepo)

	code := m.Run()
	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// TestFunctional_HoldReleaseFlow validates hold then release flow end-to-end.
func TestFunctional_HoldReleaseFlow(t *testing.T) {
	_ = context.Background()
	t.Skip("TODO: implement using real wallet and transaction records")
}

// TestFunctional_DebitCreditRefundFlow validates debit/credit/refund flow end-to-end.
func TestFunctional_DebitCreditRefundFlow(t *testing.T) {
	_ = context.Background()
	t.Skip("TODO: implement using real wallet and transaction records")
}
