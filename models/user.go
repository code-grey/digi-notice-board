package models

import "time"

// User represents a user in the system.
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"unique" json:"username"`
    Email     string    `gorm:"unique" json:"email"`
    Password  string    `json:"-"` // never expose the password in responses
    Role      string    `json:"role"` // e.g., "user" or "admin"
    CreatedAt time.Time `json:"created_at"`
}

