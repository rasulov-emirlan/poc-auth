package auth

import "errors"

var (
	ErrEmailNotFound               = errors.New("email not found")
	ErrEmailTaken                  = errors.New("email taken")
	ErrPasswordTooShort            = errors.New("password must be at least 8 characters long")
	ErrFirstnameOrLastnameTooShort = errors.New("firstname and lastname must be at least 1 character long")
)
