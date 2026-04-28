// Package repository provides data access layer for emergency-service.
package repository

import "context"

// EmergencyRepository defines the interface for emergency-service data access.
type EmergencyRepository interface {

	// TriggerSOS performs the TriggerSOS operation.
	TriggerSOS(ctx context.Context) error

	// GetEmergencyContacts performs the GetEmergencyContacts operation.
	GetEmergencyContacts(ctx context.Context) error

	// UpdateContacts performs the UpdateContacts operation.
	UpdateContacts(ctx context.Context) error
}

// postgresEmergencyRepository implements EmergencyRepository using PostgreSQL.
type postgresEmergencyRepository struct {
	// TODO: add *sql.DB field
}

// NewPostgresEmergencyRepository creates a new PostgreSQL-based repository.
func NewPostgresEmergencyRepository() EmergencyRepository {
	return &postgresEmergencyRepository{}
}
