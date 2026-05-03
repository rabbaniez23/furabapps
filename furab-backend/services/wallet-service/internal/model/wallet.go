// Package model defines the domain models for wallet-service.
package model

import "time"

type TransactionType string

const (
	TypeHold    TransactionType = "hold"
	TypeRelease TransactionType = "release"
	TypeDebit   TransactionType = "debit"
	TypeCredit  TransactionType = "credit"
	TypeRefund  TransactionType = "refund"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusSuccess   TransactionStatus = "SUCCESS"
	StatusFailed    TransactionStatus = "FAILED"
	StatusCancelled TransactionStatus = "CANCELLED"
)

// Wallet represents the Wallet model in wallet-service.
type Wallet struct {
	ID        string    `json:"wallet_id"`
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction represents the Transaction model in wallet-service.
type Transaction struct {
	ID             string            `json:"transaction_id"`
	WalletID       string            `json:"wallet_id"`
	ReferenceID    string            `json:"reference_id"`
	Type           TransactionType   `json:"type"`
	Amount         float64           `json:"amount"`
	Status         TransactionStatus `json:"status"`
	CurrentBalance float64           `json:"current_balance"`
	CreatedAt      time.Time         `json:"created_at"`
}

type WalletResult struct {
	Status         TransactionStatus `json:"status"`
	CurrentBalance float64           `json:"current_balance"`
	TransactionID  string            `json:"transaction_id"`
}
