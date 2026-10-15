package service

import "errors"

var (
	ErrInvalidUsername    = errors.New("username must be between 3 and 30 characters")
	ErrInvalidName        = errors.New("name is required")
	ErrInvalidSurname     = errors.New("surname is required")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("password must be between 8 and 72 bytes")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
