package domain

import "errors"

// Erros do domínio.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrPasswordHash       = errors.New("failed to hash password")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
