package model

import "time"

// Merchant represents a shop/store owner in the system.
type Merchant struct {
	MerchantID        string    `json:"merchant_id"`
	UserID            string    `json:"user_id"`
	NamaToko          string    `json:"nama_toko"`
	Alamat            string    `json:"alamat"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	Kategori          string    `json:"kategori"`
	JamBuka           string    `json:"jam_buka"`
	JamTutup          string    `json:"jam_tutup"`
	StatusOperasional string    `json:"status_operasional"` // "open" or "closed"
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
