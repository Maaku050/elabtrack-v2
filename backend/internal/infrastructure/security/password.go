package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implements application.PasswordHasher using bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher constructs a bcrypt hasher with the given cost.
// Use bcrypt.DefaultCost (10) unless you have a specific reason otherwise.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// Hash returns the bcrypt hash of the plaintext password.
func (h *BcryptHasher) Hash(plain string) (string, error) {
	if plain == "" {
		return "", fmt.Errorf("password must not be empty")
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash: %w", err)
	}
	return string(b), nil
}

// Compare returns nil if the plaintext matches the hash.
func (h *BcryptHasher) Compare(hashed, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)); err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
