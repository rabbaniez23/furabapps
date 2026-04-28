//go:build functional
// +build functional

// Package functional contains functional tests for user-service.
// These tests access a real database.
// Run with: go test ./test/functional/... -v -tags=functional
package functional

import (
	"testing"
)

// TestFunctional_User_CreateAndGet tests basic CRUD with real database.
func TestFunctional_User_CreateAndGet(t *testing.T) {
	// TODO: Setup test database connection
	// TODO: Create entity and verify retrieval
	t.Skip("TODO: Implement with real database")
}

// TestFunctional_User_FullFlow tests the complete lifecycle.
func TestFunctional_User_FullFlow(t *testing.T) {
	// TODO: Test full business flow with real database
	t.Skip("TODO: Implement with real database")
}
