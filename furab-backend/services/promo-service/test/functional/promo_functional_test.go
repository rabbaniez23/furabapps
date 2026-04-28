//go:build functional
// +build functional

// Package functional contains functional tests for promo-service.
// These tests access a real database.
// Run with: go test ./test/functional/... -v -tags=functional
package functional

import (
	"testing"
)

// TestFunctional_Promo_CreateAndGet tests basic CRUD with real database.
func TestFunctional_Promo_CreateAndGet(t *testing.T) {
	// TODO: Setup test database connection
	// TODO: Create entity and verify retrieval
	t.Skip("TODO: Implement with real database")
}

// TestFunctional_Promo_FullFlow tests the complete lifecycle.
func TestFunctional_Promo_FullFlow(t *testing.T) {
	// TODO: Test full business flow with real database
	t.Skip("TODO: Implement with real database")
}
