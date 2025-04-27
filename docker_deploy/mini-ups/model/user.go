package model

import "time"

// User respresents a user in our system.
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"unique_index" json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
