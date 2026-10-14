package models

import (
	"time"
	"uuid"
)

type User struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Balance      float64   `json:"balance"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}
