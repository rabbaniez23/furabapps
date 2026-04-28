//go:build functional
// +build functional

// Package functional contains functional tests for auth-service.
// These tests access a real database.
// Run with: go test ./test/functional/... -v -tags=functional
package functional

import (
	"testing"
)

// TestFunctional_Auth_CreateAndGet tests basic CRUD with real database.
func TestFunctional_Auth_CreateAndGet(t *testing.T) {
	// TODO: Setup test database connection
	// TODO: Create entity and verify retrieval
	t.Skip("TODO: Implement with real database")
}

// TestFunctional_Auth_FullFlow tests the complete lifecycle.
func TestFunctional_Auth_FullFlow(t *testing.T) {
	// TODO: Test full business flow with real database
	t.Skip("TODO: Implement with real database")
}
