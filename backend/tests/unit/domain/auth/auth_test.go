package auth_test

import (
	"strings"
	"testing"
	"time"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/google/uuid"
)

func TestRefreshToken_ValidAt(t *testing.T) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	valid := domainauth.NewRefreshToken(strings.Repeat("a", 64), uuid.New(), now, now.Add(time.Hour))
	if !valid.ValidAt(now) {
		t.Fatal("fresh session denied")
	}
	if valid.ValidAt(valid.ExpiresAt) || valid.ValidAt(valid.ExpiresAt.Add(time.Nanosecond)) {
		t.Fatal("expiry boundary must deny")
	}
	valid.RevokedAt = &now
	if valid.ValidAt(now) {
		t.Fatal("revoked session accepted")
	}
}
