package repository

import "errors"

var (
	ErrEmailTaken        = errors.New("email already taken")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrUserNotFound      = errors.New("user not found")
	ErrTicketUnavailable = errors.New("ticket unavailable")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
