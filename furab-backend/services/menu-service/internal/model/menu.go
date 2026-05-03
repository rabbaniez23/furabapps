package model

import "time"

// Menu represents a menu item offered by a merchant.
type Menu struct {
	MenuID      string    `json:"menu_id"`
	MerchantID  string    `json:"merchant_id"`
	NamaMenu    string    `json:"nama_menu"`
	Harga       float64   `json:"harga"`
	Kategori    string    `json:"kategori"`
	Deskripsi   string    `json:"deskripsi"`
	Stok        int       `json:"stok"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
