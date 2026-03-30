package errors

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict")
	ErrInvalidID  = errors.New("invalid id")
)
