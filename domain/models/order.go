package models

import "time"

type User struct {
	ID        string `json:"id"`        // User ID
	FirstName string `json:"firstName"` // User's first name
	LastName  string `json:"lastName"`  // User's last name
	Email     string `json:"email"`     // User's email address
}

// BulkOrderEvent represents the event to create a bulk order
type BulkOrderEvent struct {
	FilePath    string    `json:"filePath"`    // Path to the uploaded file
	User        User      `json:"user"`        // User who triggered the request
	RequestTime time.Time `json:"requestTime"` // Timestamp of the request
}

type Order struct {
	HubID     string    `json:"hub_id"`
	SkuId     string    `json:"sku_id"`
	TenantID  string    `json:"tenant_id"`
	SellerID  string    `json:"seller_id"`
	Qty       int       `json:"qty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
