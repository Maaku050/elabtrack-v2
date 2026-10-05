package shared

import "errors"

// Sentinel errors shared across domains. Domain layers declare their own
// context-specific errors; these are the generic cross-cutting ones.
var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation error")
	ErrInternal     = errors.New("internal server error")
	ErrInvalidInput = errors.New("invalid input")
)
