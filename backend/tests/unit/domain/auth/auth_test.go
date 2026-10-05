package auth_test

import (
	"testing"
	"time"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/google/uuid"
)

func TestRefreshToken_IsValid(t *testing.T) {
	userID := uuid.New()
	future := time.Now().Add(time.Hour).UTC()
	past := time.Now().Add(-time.Hour).UTC()

	valid := domainauth.NewRefreshToken("tok", userID, future)
	if !valid.IsValid() {
		t.Fatal("expected fresh token to be valid")
	}

	expired := domainauth.NewRefreshToken("tok", userID, past)
	if expired.IsValid() {
		t.Fatal("expected expired token to be invalid")
	}

	revoked := domainauth.NewRefreshToken("tok", userID, future)
	revoked.Revoke()
	if revoked.IsValid() {
		t.Fatal("expected revoked token to be invalid")
	}
	if !revoked.IsRevoked() {
		t.Fatal("expected revoked flag to be set")
	}
}
