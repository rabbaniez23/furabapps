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
	dbHost := getEnvOrDefault("DB_HOST", "127.0.0.1")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "furab")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "furab_secret")
	dbName := getEnvOrDefault("DB_NAME", "pricing_service")

	// Step 1: Auto-create database
	adminDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort)
	adminDB, err := sql.Open("pgx", adminDSN)
	if err != nil {
		log.Printf("functional tests skipped: failed to connect admin database: %v", err)
		os.Exit(0)
	}
	for i := 0; i < 30; i++ {
		if err = adminDB.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Printf("functional tests skipped: database not ready: %v", err)
		os.Exit(0)
	}
	_, _ = adminDB.Exec("CREATE DATABASE " + dbName)
	adminDB.Close()

	// Step 2: Connect to target database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("functional tests skipped: failed to connect database: %v", err)
		os.Exit(0)
	}
	defer testDB.Close()

	for i := 0; i < 30; i++ {
		if err = testDB.Ping(); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Printf("functional tests skipped: database not ready: %v", err)
		os.Exit(0)
	}

	// Setup schema and seed data
	setupSchema()

	testRepo = repository.NewPostgresPriceRepository(testDB)
	testSvc = service.NewPriceService(testRepo, client.NewDummyOrderClient(), client.NewDummyLocationClient())

	code := m.Run()

	teardownSchema()
	os.Exit(code)
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func setupSchema() {
	query := `
		CREATE TABLE IF NOT EXISTS pricing_rules (
			rule_id VARCHAR(36) PRIMARY KEY,
			type VARCHAR(50) NOT NULL UNIQUE,
			value DOUBLE PRECISION NOT NULL,
			description TEXT
		);
		DELETE FROM pricing_rules;
		INSERT INTO pricing_rules (rule_id, type, value, description) VALUES
		('r1', 'delivery', 2500, 'Delivery fee per km'),
		('r2', 'service', 0.1, 'Service fee 10%');
	`
	_, err := testDB.Exec(query)
	if err != nil {
		log.Fatalf("Failed to setup schema: %v", err)
	}
}

func teardownSchema() {
	testDB.Exec("DROP TABLE IF EXISTS pricing_rules")
}

// TestFunctional_Price_Calculate tests end-to-end pricing rule reads and calculation.
func TestFunctional_Price_Calculate(t *testing.T) {
	ctx := context.Background()

	// Dummy order client will return items worth some money.
	// Dummy location client will return some distance.
	// We just ensure the DB rules are successfully fetched and applied.

	orderID := "order-func-123"
	res, err := testSvc.CalculatePrice(ctx, orderID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.OrderID != orderID {
		t.Errorf("expected order ID %s, got %s", orderID, res.OrderID)
	}

	if res.DeliveryFee <= 0 {
		t.Errorf("expected delivery fee > 0, got %f", res.DeliveryFee)
	}

	if res.ServiceFee <= 0 {
		t.Errorf("expected service fee > 0, got %f", res.ServiceFee)
	}

	if res.TotalAmount <= 0 {
		t.Errorf("expected total amount > 0, got %f", res.TotalAmount)
	}

	t.Logf("Success CalculatePrice: Total=%f, Delivery=%f, Service=%f", res.TotalAmount, res.DeliveryFee, res.ServiceFee)
}

func TestFunctional_Price_MissingRule(t *testing.T) {
	ctx := context.Background()
	testDB.Exec("DELETE FROM pricing_rules WHERE type = 'delivery'")

	_, err := testSvc.CalculatePrice(ctx, "order-missing")
	if err == nil {
		t.Errorf("expected error due to missing delivery rule")
	}

	// Restore for other tests if any
	setupSchema()
}
