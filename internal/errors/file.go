package errors

import "errors"

var (
	ErrEmptyFile     = errors.New("file is empty")
	ErrFileRead      = errors.New("error reading file")
	ErrFileOpen      = errors.New("error opening file")
	ErrInvalidFormat = errors.New("invalid file formate")
)
