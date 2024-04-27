package errors

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidProvider = errors.New("invalid provider, valid providers are google and github")
	ErrInvalidToken    = errors.New("invalid token")
)
