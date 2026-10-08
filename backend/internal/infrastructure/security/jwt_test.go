package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func testIssuer() (*JWTIssuer, *domainuser.User, time.Time) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	j := NewJWTIssuer(config.JWTConfig{Secret: "synthetic-test-signing-material-only", Issuer: "elabtrack-v2", AccessTTL: time.Minute})
	j.now = func() time.Time { return now }
	return j, &domainuser.User{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Email: "synthetic@example.invalid", Role: domainuser.RoleAdmin}, now
}
func TestAccessJWTValidation(t *testing.T) {
	j, u, now := testIssuer()
	raw, err := j.IssueAccessToken(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	got, err := j.VerifyAccessToken(context.Background(), raw)
	if err != nil || got.UserID != u.ID || got.Role != string(domainuser.RoleAdmin) {
		t.Fatal("valid issuance rejected")
	}
	base := accessClaims{UserID: u.ID, Email: u.Email, Role: "admin", Purpose: accessPurpose, RegisteredClaims: jwt.RegisteredClaims{
		Issuer: j.issuer, Subject: u.ID.String(), Audience: jwt.ClaimStrings{accessAudience},
		IssuedAt: jwt.NewNumericDate(now), NotBefore: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
	}}
	cases := []struct {
		name   string
		mutate func(*accessClaims)
		method jwt.SigningMethod
		key    any
	}{
		{name: "wrong signature", key: []byte("different-key")},
		{name: "HS384", method: jwt.SigningMethodHS384}, {name: "HS512", method: jwt.SigningMethodHS512},
		{name: "none", method: jwt.SigningMethodNone, key: jwt.UnsafeAllowNoneSignatureType},
		{name: "expired", mutate: func(c *accessClaims) { c.ExpiresAt = jwt.NewNumericDate(now.Add(-time.Second)) }},
		{name: "exact expiry", mutate: func(c *accessClaims) { c.ExpiresAt = jwt.NewNumericDate(now) }},
		{name: "missing expiry", mutate: func(c *accessClaims) { c.ExpiresAt = nil }},
		{name: "wrong issuer", mutate: func(c *accessClaims) { c.Issuer = "another-issuer" }},
		{name: "missing issuer", mutate: func(c *accessClaims) { c.Issuer = "" }},
		{name: "wrong audience", mutate: func(c *accessClaims) { c.Audience = jwt.ClaimStrings{"another-api"} }},
		{name: "missing audience", mutate: func(c *accessClaims) { c.Audience = nil }},
		{name: "wrong purpose", mutate: func(c *accessClaims) { c.Purpose = "refresh" }},
		{name: "missing purpose", mutate: func(c *accessClaims) { c.Purpose = "" }},
		{name: "missing subject", mutate: func(c *accessClaims) { c.Subject = "" }},
		{name: "mismatched subject", mutate: func(c *accessClaims) { c.Subject = uuid.NewString() }},
		{name: "malformed subject", mutate: func(c *accessClaims) { c.Subject = "not-a-uuid" }},
		{name: "zero uid", mutate: func(c *accessClaims) { c.UserID = uuid.Nil }},
		{name: "future iat", mutate: func(c *accessClaims) { c.IssuedAt = jwt.NewNumericDate(now.Add(time.Second)) }},
		{name: "missing iat", mutate: func(c *accessClaims) { c.IssuedAt = nil }},
		{name: "future nbf", mutate: func(c *accessClaims) { c.NotBefore = jwt.NewNumericDate(now.Add(time.Second)) }},
		{name: "missing nbf", mutate: func(c *accessClaims) { c.NotBefore = nil }},
		{name: "nbf before iat", mutate: func(c *accessClaims) { c.NotBefore = jwt.NewNumericDate(now.Add(-time.Second)) }},
		{name: "expiry before issuance", mutate: func(c *accessClaims) {
			c.IssuedAt = jwt.NewNumericDate(now.Add(-time.Second))
			c.NotBefore = c.IssuedAt
			c.ExpiresAt = jwt.NewNumericDate(now.Add(-2 * time.Second))
		}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := base
			if tt.mutate != nil {
				tt.mutate(&c)
			}
			method := tt.method
			if method == nil {
				method = jwt.SigningMethodHS256
			}
			key := tt.key
			if key == nil {
				key = j.secret
			}
			bad, err := jwt.NewWithClaims(method, c).SignedString(key)
			if err != nil {
				t.Fatal(err)
			}
			claims, err := j.VerifyAccessToken(context.Background(), bad)
			if err == nil || claims.UserID != uuid.Nil {
				t.Fatal("invalid credential accepted")
			}
			if err.Error() != errInvalidAccessToken.Error() || strings.Contains(err.Error(), bad) {
				t.Fatal("unsafe error")
			}
		})
	}
	t.Run("opaque refresh is not access", func(t *testing.T) {
		refresh, err := j.GenerateRefreshToken()
		if err != nil {
			t.Fatal(err)
		}
		if _, err = j.VerifyAccessToken(context.Background(), refresh); err == nil {
			t.Fatal("refresh accepted as access")
		}
	})
	t.Run("missing uid", func(t *testing.T) {
		c := jwt.MapClaims{"iss": j.issuer, "aud": accessAudience, "sub": u.ID.String(), "iat": now.Unix(), "nbf": now.Unix(), "exp": now.Add(time.Minute).Unix(), "purpose": accessPurpose}
		bad, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(j.secret)
		if _, err := j.VerifyAccessToken(context.Background(), bad); err == nil {
			t.Fatal("missing uid accepted")
		}
	})
}
func TestOpaqueRefreshAndDigest(t *testing.T) {
	j, _, _ := testIssuer()
	previous := ""
	for range 4 {
		raw, err := j.GenerateRefreshToken()
		if err != nil {
			t.Fatal(err)
		}
		bytes, err := hex.DecodeString(raw)
		if err != nil || len(bytes) != 32 || raw != strings.ToLower(raw) || raw == previous {
			t.Fatal("opaque format/entropy source contract failed")
		}
		hash, err := (SHA256RefreshHasher{}).Hash(raw)
		if err != nil {
			t.Fatal(err)
		}
		expected := sha256.Sum256([]byte(raw))
		if hash == raw || hash != hex.EncodeToString(expected[:]) {
			t.Fatal("digest contract failed")
		}
		previous = raw
	}
	for _, raw := range []string{"", "credential-sentinel", strings.Repeat("A", 64), strings.Repeat("g", 64), strings.Repeat("a", 65)} {
		hash, err := (SHA256RefreshHasher{}).Hash(raw)
		if err == nil || hash != "" || strings.Contains(err.Error(), raw) && raw != "" {
			t.Fatal("invalid credential/hash error contract")
		}
	}
}
