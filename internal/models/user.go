package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Role         UserRole  `json:"role"`
	Balance      int64     `json:"balance"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
