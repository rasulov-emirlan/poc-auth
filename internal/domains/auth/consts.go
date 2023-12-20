package auth

import "errors"

var (
	ErrEmailNotFound = errors.New("email not found")
	ErrEmailTaken    = errors.New("email taken")
)
