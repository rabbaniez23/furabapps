package model

import "time"

// Rating represents a single rating entry
type Rating struct {
	RatingID   string    `json:"rating_id"`
	ReviewerID string    `json:"reviewer_id"`
	OrderID    string    `json:"order_id"`
	TargetType string    `json:"target_type"` // "driver" or "merchant"
	TargetID   string    `json:"target_id"`
	Score      int       `json:"score"`
	CreatedAt  time.Time `json:"created_at"`
}

// RatingSummary represents the aggregated statistics for a target
type RatingSummary struct {
	TargetType   string       `json:"target_type"`
	TargetID     string       `json:"target_id"`
	AverageScore float64      `json:"average_score"`
	TotalCount   int          `json:"total_count"`
	Distribution map[int]int  `json:"distribution"`
	LastUpdated  time.Time    `json:"last_updated"`
}
