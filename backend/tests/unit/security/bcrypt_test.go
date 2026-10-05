package security_test

import (
	"testing"

	"github.com/fullstacktemplate/backend/internal/infrastructure/security"
)

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	h := security.NewBcryptHasher(0) // default cost

	hashed, err := h.Hash("password123")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if hashed == "password123" {
		t.Fatal("hash should differ from plaintext")
	}

	if err := h.Compare(hashed, "password123"); err != nil {
		t.Fatalf("expected compare to pass: %v", err)
	}
	if err := h.Compare(hashed, "wrong"); err == nil {
		t.Fatal("expected compare to fail for wrong password")
	}
}

func TestBcryptHasher_EmptyRejected(t *testing.T) {
	h := security.NewBcryptHasher(0)
	if _, err := h.Hash(""); err == nil {
		t.Fatal("expected empty password to be rejected")
	}
}
