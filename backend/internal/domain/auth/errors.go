package auth

import "errors"

// IsTokenNotFound reports whether err is the token-not-found domain error.
func IsTokenNotFound(err error) bool { return errors.Is(err, ErrTokenNotFound) }
