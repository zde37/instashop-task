package pkg

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrDatabase           = errors.New("database error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInvalidInput       = errors.New("invalid input")
	ErrOrderNotPending    = errors.New("order is not in pending status")
	ErrUnauthorized       = errors.New("unauthorized action")
)
