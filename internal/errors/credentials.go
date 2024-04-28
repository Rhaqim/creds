package errors

import "errors"

var (
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidUser            = errors.New("invalid user")
	ErrInvalidOrganization    = errors.New("invalid organization")
	ErrInvalidCredential      = errors.New("invalid credential")
	ErrInvalidEnvironment     = errors.New("invalid environment")
	ErrInvalidFile            = errors.New("invalid file")
	ErrInvalidCredentialField = errors.New("invalid credential field")
)
