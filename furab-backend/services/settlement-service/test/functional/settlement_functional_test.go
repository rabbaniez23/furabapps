//go:build functional
// +build functional

// Package functional contains functional tests for settlement-service.
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

	"furab-backend/services/settlement-service/internal/repository"
	"furab-backend/services/settlement-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB   *sql.DB
	testRepo repository.SettlementRepository
	testSvc  service.SettlementService
)

func TestMain(m *testing.M) {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "furab")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "furab_secret")
	dbName := getEnvOrDefault("DB_NAME", "settlement_service_test")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

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

	testRepo = repository.NewPostgresSettlementRepository(testDB)
	testSvc = service.NewSettlementService(testRepo, nil, nil, nil)

	code := m.Run()
	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// TestFunctional_Settlement_ProcessSuccess validates successful split distribution flow.
func TestFunctional_Settlement_ProcessSuccess(t *testing.T) {
	_ = context.Background()
	t.Skip("TODO: implement using real settlements data and downstream service stubs")
}

// TestFunctional_Settlement_Idempotency validates duplicate payment trigger handling.
func TestFunctional_Settlement_Idempotency(t *testing.T) {
	_ = context.Background()
	t.Skip("TODO: implement idempotency verification using same payment_id")
}
