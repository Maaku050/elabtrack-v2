package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// RandomHex returns n random bytes encoded as hex (2n chars).
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("random hex: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Ptr returns a pointer to v. Useful for optional DTO fields.
func Ptr[T any](v T) *T { return &v }

// NormalizeEmail lowercases and trims an email address.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Coalesce returns the first non-empty string.
func Coalesce(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
