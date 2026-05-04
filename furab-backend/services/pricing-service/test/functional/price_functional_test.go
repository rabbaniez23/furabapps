//go:build functional
// +build functional

// Package functional contains functional tests for pricing-service.
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

	"furab-backend/services/pricing-service/internal/client"
	"furab-backend/services/pricing-service/internal/repository"
	"furab-backend/services/pricing-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB   *sql.DB
	testRepo repository.PriceRepository
	testSvc  service.PriceService
)

func TestMain(m *testing.M) {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "furab")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "furab_secret")
	dbName := getEnvOrDefault("DB_NAME", "pricing_service_test")

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

	testRepo = repository.NewPostgresPriceRepository(testDB)
	testSvc = service.NewPriceService(testRepo, client.NewDummyOrderClient(), client.NewDummyLocationClient())

	code := m.Run()
	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// TestFunctional_Price_Calculate tests end-to-end pricing rule reads and calculation.
func TestFunctional_Price_Calculate(t *testing.T) {
	_ = context.Background()
	_ = testSvc
	t.Skip("TODO: implement with real pricing_rules data")
}
