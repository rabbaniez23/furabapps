package model

import "time"

// Review represents a textual review given by a user to a driver or merchant.
type Review struct {
	ReviewID   string    `json:"review_id"`
	UserID     string    `json:"user_id"`
	TargetID   string    `json:"target_id"`
	TargetType string    `json:"target_type"` // "driver" or "merchant"
	OrderID    string    `json:"order_id"`
	RatingID   *string   `json:"rating_id,omitempty"` // optional field
	Comment    string    `json:"comment"`
	Status     string    `json:"status"` // "active", "flagged", "removed"
	CreatedAt  time.Time `json:"created_at"`
}

// ReviewReport represents a report filed against a specific review.
type ReviewReport struct {
	ReportID   string    `json:"report_id"`
	ReviewID   string    `json:"review_id"`
	ReportedBy string    `json:"reported_by"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
}
