// Package model defines the domain models for menu-service.
package model

import "time"

// Menu represents the Menu model in menu-service.
type Menu struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Menu-specific fields
}

// MenuItem represents the MenuItem model in menu-service.
type MenuItem struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add MenuItem-specific fields
}

// Category represents the Category model in menu-service.
type Category struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Add Category-specific fields
}

