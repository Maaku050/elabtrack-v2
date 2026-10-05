package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// SHA256RefreshHasher digests a uniformly random 256-bit credential. Unlike
// passwords, these credentials have no low-entropy dictionary to slow down.
type SHA256RefreshHasher struct{}

func (SHA256RefreshHasher) Hash(raw string) (string, error) {
	if len(raw) != 64 {
		return "", errors.New("invalid refresh credential")
	}
	for _, c := range raw {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return "", errors.New("invalid refresh credential")
		}
	}
	digest := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(digest[:]), nil
}
