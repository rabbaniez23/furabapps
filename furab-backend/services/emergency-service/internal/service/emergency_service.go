// Package service implements the business logic for emergency-service.
package service

import "context"

// EmergencyService defines the interface for emergency-service business logic.
type EmergencyService interface {

	// TriggerSOS implements the business logic for TriggerSOS.
	TriggerSOS(ctx context.Context) error

	// GetEmergencyContacts implements the business logic for GetEmergencyContacts.
	GetEmergencyContacts(ctx context.Context) error

	// UpdateContacts implements the business logic for UpdateContacts.
	UpdateContacts(ctx context.Context) error
}

// emergencyServiceImpl is the concrete implementation of EmergencyService.
type emergencyServiceImpl struct {
	// TODO: add repository and event publisher dependencies
}

// NewEmergencyService creates a new EmergencyService.
func NewEmergencyService() EmergencyService {
	return &emergencyServiceImpl{}
}
