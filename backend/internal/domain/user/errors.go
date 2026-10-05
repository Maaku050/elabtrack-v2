package user

import "errors"

// Sentinel errors specific to the user domain are declared in entity.go.
// This file is reserved for additional user-domain error helpers if needed.

// IsUserNotFound reports whether err is the user-not-found domain error.
func IsUserNotFound(err error) bool { return errors.Is(err, ErrUserNotFound) }
