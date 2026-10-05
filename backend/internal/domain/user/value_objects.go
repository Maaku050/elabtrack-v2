package user

import (
	"net/mail"
	"strings"
	"unicode"
)

// Email is a value object representing a validated email address.
type Email string

// ParseEmail validates and normalizes an email address.
func ParseEmail(raw string) (Email, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil {
		return "", ErrInvalidEmail
	}
	return Email(strings.ToLower(addr.Address)), nil
}

// String returns the email as a string.
func (e Email) String() string { return string(e) }

// Name is a value object representing a non-empty display name.
type Name string

// ParseName validates and normalizes a display name.
func ParseName(raw string) (Name, error) {
	t := strings.TrimSpace(raw)
	if len(t) < 1 || len(t) > 100 {
		return "", ErrInvalidName
	}
	if !hasLetter(t) {
		return "", ErrInvalidName
	}
	return Name(t), nil
}

func hasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// String returns the name as a string.
func (n Name) String() string { return string(n) }

// PlainPassword is a value object for an unhashed password prior to hashing.
type PlainPassword string

// Validate enforces the minimum password strength.
func (p PlainPassword) Validate() error {
	if len(p) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}
