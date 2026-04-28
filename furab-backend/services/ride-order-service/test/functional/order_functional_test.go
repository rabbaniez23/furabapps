//go:build functional
// +build functional

// Package functional contains functional tests for the ride order service.
// Functional tests access a real PostgreSQL database.
// Run with: go test ./test/functional/... -v -tags=functional
//
// Prerequisites:
//   - PostgreSQL running on localhost:5432
//   - Database "ride_order_service_test" created
//   - Environment variable DB_HOST, DB_PORT, DB_USER, DB_PASSWORD set (or use defaults)
package functional

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"furab-backend/services/ride-order-service/internal/model"
	"furab-backend/services/ride-order-service/internal/repository"
	"furab-backend/services/ride-order-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	testDB   *sql.DB
	testRepo repository.OrderRepository
	testSvc  service.OrderService
)

// TestMain sets up the test database and service before running functional tests.
// It creates the necessary table, runs all tests, and cleans up afterward.
func TestMain(m *testing.M) {
	// Setup database connection
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "furab")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "furab_secret")
	dbName := getEnvOrDefault("DB_NAME", "ride_order_service_test")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	testDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	// Wait for DB to be ready
	for i := 0; i < 30; i++ {
		err = testDB.Ping()
		if err == nil {
			break
		}
		log.Printf("Waiting for database... (%d/30)", i+1)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("Database is not ready: %v", err)
	}

	// Setup schema
	if err := setupSchema(); err != nil {
		log.Fatalf("Failed to setup schema: %v", err)
	}

	// Initialize repository and service (no event publisher for functional tests)
	testRepo = repository.NewPostgresOrderRepository(testDB)
	testSvc = service.NewOrderService(testRepo, nil) // nil publisher for testing

	// Run tests
	code := m.Run()

	// Cleanup
	teardownSchema()
	os.Exit(code)
}

// setupSchema creates the ride_orders table for testing.
func setupSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS ride_orders (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) NOT NULL,
			driver_id VARCHAR(36),
			pickup_lat DOUBLE PRECISION NOT NULL,
			pickup_lng DOUBLE PRECISION NOT NULL,
			pickup_address TEXT NOT NULL,
			dropoff_lat DOUBLE PRECISION NOT NULL,
			dropoff_lng DOUBLE PRECISION NOT NULL,
			dropoff_address TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
			fare DOUBLE PRECISION NOT NULL DEFAULT 0,
			distance DOUBLE PRECISION NOT NULL DEFAULT 0,
			estimated_duration INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_ride_orders_user_id ON ride_orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_ride_orders_driver_id ON ride_orders(driver_id);
		CREATE INDEX IF NOT EXISTS idx_ride_orders_status ON ride_orders(status);
	`
	_, err := testDB.Exec(query)
	return err
}

// teardownSchema drops the test table.
func teardownSchema() {
	testDB.Exec("DROP TABLE IF EXISTS ride_orders")
}

// cleanupOrders removes all orders from the test table.
func cleanupOrders() {
	testDB.Exec("DELETE FROM ride_orders")
}

// getEnvOrDefault reads an environment variable or returns a default.
func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// --- Functional Test Cases ---

// TestFunctional_CreateAndGetOrder tests the full flow of creating and retrieving an order.
func TestFunctional_CreateAndGetOrder(t *testing.T) {
	cleanupOrders()

	ctx := context.Background()

	// Create a ride order
	req := &model.CreateRideOrderRequest{
		UserID: "func-user-001",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Monas, Jakarta Pusat",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.1751,
			Longitude: 106.8650,
			Address:   "Ancol, Jakarta Utara",
		},
	}

	order, err := testSvc.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Verify created order
	if order.ID == "" {
		t.Fatal("expected non-empty order ID")
	}
	if order.Status != model.RideStatusPending {
		t.Errorf("expected status PENDING, got: %s", order.Status)
	}
	if order.UserID != req.UserID {
		t.Errorf("expected user ID %s, got: %s", req.UserID, order.UserID)
	}

	// Retrieve the order
	fetched, err := testSvc.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("failed to get order: %v", err)
	}

	if fetched.ID != order.ID {
		t.Errorf("expected order ID %s, got: %s", order.ID, fetched.ID)
	}
	if fetched.PickupLocation.Address != req.PickupLocation.Address {
		t.Errorf("expected pickup address %s, got: %s",
			req.PickupLocation.Address, fetched.PickupLocation.Address)
	}
	if fetched.Fare != order.Fare {
		t.Errorf("expected fare %.2f, got: %.2f", order.Fare, fetched.Fare)
	}
}

// TestFunctional_FullRideFlow tests the complete ride lifecycle:
// Create → Assign → Start → Complete
func TestFunctional_FullRideFlow(t *testing.T) {
	cleanupOrders()

	ctx := context.Background()

	// Step 1: Create order
	req := &model.CreateRideOrderRequest{
		UserID: "func-user-002",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Sudirman, Jakarta",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.2600,
			Longitude: 106.7810,
			Address:   "Senayan, Jakarta",
		},
	}

	order, err := testSvc.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("step 1 - create: %v", err)
	}
	t.Logf("Created order: %s (status: %s)", order.ID, order.Status)

	// Step 2: Assign driver
	assigned, err := testSvc.AssignDriver(ctx, order.ID, "driver-func-001")
	if err != nil {
		t.Fatalf("step 2 - assign driver: %v", err)
	}
	if assigned.Status != model.RideStatusAssigned {
		t.Errorf("expected ASSIGNED, got: %s", assigned.Status)
	}
	if assigned.DriverID != "driver-func-001" {
		t.Errorf("expected driver ID driver-func-001, got: %s", assigned.DriverID)
	}
	t.Logf("Assigned driver: %s (status: %s)", assigned.DriverID, assigned.Status)

	// Step 3: Start ride
	started, err := testSvc.StartRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("step 3 - start ride: %v", err)
	}
	if started.Status != model.RideStatusStarted {
		t.Errorf("expected STARTED, got: %s", started.Status)
	}
	t.Logf("Ride started (status: %s)", started.Status)

	// Step 4: Complete ride
	completed, err := testSvc.CompleteRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("step 4 - complete ride: %v", err)
	}
	if completed.Status != model.RideStatusCompleted {
		t.Errorf("expected COMPLETED, got: %s", completed.Status)
	}
	t.Logf("Ride completed (status: %s, fare: Rp %.0f)", completed.Status, completed.Fare)

	// Verify final state in database
	final, err := testSvc.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("failed to get final order: %v", err)
	}
	if final.Status != model.RideStatusCompleted {
		t.Errorf("expected final status COMPLETED, got: %s", final.Status)
	}
}

// TestFunctional_CancelRide tests cancelling a ride order.
func TestFunctional_CancelRide(t *testing.T) {
	cleanupOrders()

	ctx := context.Background()

	// Create order
	req := &model.CreateRideOrderRequest{
		UserID: "func-user-003",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Thamrin, Jakarta",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.3000,
			Longitude: 106.8500,
			Address:   "Kuningan, Jakarta",
		},
	}

	order, err := testSvc.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Cancel the ride
	cancelled, err := testSvc.CancelRide(ctx, order.ID)
	if err != nil {
		t.Fatalf("failed to cancel: %v", err)
	}
	if cancelled.Status != model.RideStatusCancelled {
		t.Errorf("expected CANCELLED, got: %s", cancelled.Status)
	}

	// Verify in database
	final, err := testSvc.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	if final.Status != model.RideStatusCancelled {
		t.Errorf("expected CANCELLED in DB, got: %s", final.Status)
	}
}

// TestFunctional_InvalidTransition tests that invalid status transitions are rejected.
// Trying to complete a PENDING order (skipping assign and start) should fail.
func TestFunctional_InvalidTransition(t *testing.T) {
	cleanupOrders()

	ctx := context.Background()

	// Create order
	req := &model.CreateRideOrderRequest{
		UserID: "func-user-004",
		PickupLocation: model.Location{
			Latitude:  -6.2088,
			Longitude: 106.8456,
			Address:   "Blok M, Jakarta",
		},
		DropoffLocation: model.Location{
			Latitude:  -6.3500,
			Longitude: 106.8300,
			Address:   "Pondok Indah, Jakarta",
		},
	}

	order, err := testSvc.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	// Try to complete directly (should fail - must go through ASSIGNED → STARTED first)
	_, err = testSvc.CompleteRide(ctx, order.ID)
	if err == nil {
		t.Fatal("expected error when completing PENDING order")
	}
	t.Logf("Correctly rejected invalid transition: %v", err)

	// Try to start directly (should fail - must be ASSIGNED first)
	_, err = testSvc.StartRide(ctx, order.ID)
	if err == nil {
		t.Fatal("expected error when starting PENDING order")
	}
	t.Logf("Correctly rejected invalid transition: %v", err)
}

// TestFunctional_GetUserOrders tests retrieving multiple orders for a user.
func TestFunctional_GetUserOrders(t *testing.T) {
	cleanupOrders()

	ctx := context.Background()
	userID := "func-user-005"

	// Create 3 orders for this user
	for i := 0; i < 3; i++ {
		req := &model.CreateRideOrderRequest{
			UserID: userID,
			PickupLocation: model.Location{
				Latitude:  -6.2088 + float64(i)*0.01,
				Longitude: 106.8456,
				Address:   fmt.Sprintf("Pickup Location %d", i+1),
			},
			DropoffLocation: model.Location{
				Latitude:  -6.3000 + float64(i)*0.01,
				Longitude: 106.8500,
				Address:   fmt.Sprintf("Dropoff Location %d", i+1),
			},
		}

		_, err := testSvc.CreateOrder(ctx, req)
		if err != nil {
			t.Fatalf("failed to create order %d: %v", i+1, err)
		}
	}

	// Get all orders for this user
	orders, total, err := testSvc.GetUserOrders(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("failed to get user orders: %v", err)
	}

	if total != 3 {
		t.Errorf("expected total 3, got: %d", total)
	}
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got: %d", len(orders))
	}

	// Test pagination
	ordersPage, _, err := testSvc.GetUserOrders(ctx, userID, 2, 0)
	if err != nil {
		t.Fatalf("failed to get paginated orders: %v", err)
	}
	if len(ordersPage) != 2 {
		t.Errorf("expected 2 orders in page, got: %d", len(ordersPage))
	}
}
