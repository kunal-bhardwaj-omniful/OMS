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
