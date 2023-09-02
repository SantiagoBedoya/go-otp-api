package service

import "errors"

var (
	ErrEmailInUse      = errors.New("email is already in use")
	ErrInvalidPassword = errors.New("password does not match")
	ErrUserNotFound    = errors.New("user is not found")
)
