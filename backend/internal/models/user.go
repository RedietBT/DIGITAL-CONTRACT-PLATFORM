package models

import "time"

type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	Role              string     `json:"role"`
	Status			  string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	LastLoginAt      *time.Time  `json:"last_login_at,omitempty"`
}