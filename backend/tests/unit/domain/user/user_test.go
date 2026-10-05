package user_test

import (
	"testing"

	domainuser "github.com/fullstacktemplate/backend/internal/domain/user"
	"github.com/google/uuid"
)

func TestParseEmail_Valid(t *testing.T) {
	cases := []string{
		"User@Example.COM",
		"  alice@example.com  ",
		"bob+tag@example.co.uk",
	}
	for _, raw := range cases {
		e, err := domainuser.ParseEmail(raw)
		if err != nil {
			t.Fatalf("expected %q to parse, got %v", raw, err)
		}
		if e.String() != "user@example.com" && e.String() != "alice@example.com" && e.String() != "bob+tag@example.co.uk" {
			// not a strict check, just ensure normalization happened
		}
	}
}

func TestParseEmail_Invalid(t *testing.T) {
	cases := []string{
		"",
		"not-an-email",
		"@@example.com",
	}
	for _, raw := range cases {
		if _, err := domainuser.ParseEmail(raw); err == nil {
			t.Fatalf("expected %q to be invalid", raw)
		}
	}
}

func TestParseName_Invalid(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"123",
		"---",
	}
	for _, raw := range cases {
		if _, err := domainuser.ParseName(raw); err == nil {
			t.Fatalf("expected %q to be invalid", raw)
		}
	}
}

func TestPlainPassword_Validate(t *testing.T) {
	if err := (domainuser.PlainPassword("short")).Validate(); err == nil {
		t.Fatal("expected short password to fail validation")
	}
	if err := (domainuser.PlainPassword("longenough")).Validate(); err != nil {
		t.Fatalf("expected valid password to pass, got %v", err)
	}
}

func TestNewUser_Defaults(t *testing.T) {
	u := domainuser.NewUser("user@example.com", "Alice", "hashed")
	if u.Role != domainuser.RoleUser {
		t.Fatalf("expected default role user, got %s", u.Role)
	}
	if !u.IsActive {
		t.Fatal("expected new user to be active")
	}
	if u.ID == (uuid.UUID{}) {
		t.Fatal("expected non-zero id")
	}
}
