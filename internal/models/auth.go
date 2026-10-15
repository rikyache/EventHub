package models

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (r RegisterRequest) IsValid() error {
	if strings.TrimSpace(r.Username) == "" {
		return errors.New("username is required")
	}

	if len([]rune(r.Username)) < 3 || len([]rune(r.Username)) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return errors.New("invalid email format")
	}

	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(r.Surname) == "" {
		return errors.New("surname is required")
	}

	return nil
}
