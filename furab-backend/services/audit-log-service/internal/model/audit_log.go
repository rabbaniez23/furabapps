package model

import "time"

// AuditLog represents a record of an action performed in the system.
type AuditLog struct {
	LogID       string                 `json:"log_id"`
	ServiceName string                 `json:"service_name"`
	ActorID     string                 `json:"actor_id"`
	ActorType   string                 `json:"actor_type"`
	Action      string                 `json:"action"`
	TargetID    string                 `json:"target_id"`
	TargetType  string                 `json:"target_type"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}
