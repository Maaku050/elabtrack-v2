package terms

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestPublicationValidationAndContentIdentity(t *testing.T) {
	actor := uuid.New()
	body := "SYNTHETIC TEST DOCUMENT\nLiteral <script> text, not markup.\n"
	v, err := NewVersion(" test-1 ", " Test only ", body, actor)
	if err != nil || v.Version != "test-1" || v.Title != "Test only" || v.Body != body || v.ContentHash != fmt.Sprintf("%x", sha256.Sum256([]byte(body))) {
		t.Fatal("content identity not retained")
	}
	for _, c := range []struct {
		version, title, body string
		actor                uuid.UUID
	}{
		{"", "title", "text", actor}, {"invalid version", "title", "text", actor}, {"v1", " ", "text", actor}, {"v1", strings.Repeat("x", 201), "text", actor},
		{"v1", "title", " ", actor}, {"v1", "title", strings.Repeat("x", MaxBodyBytes+1), actor}, {"v1", "title", "text", uuid.Nil}, {"v1", "title", "x\x00y", actor}, {"v1", "title", string([]byte{255}), actor},
	} {
		if _, err := NewVersion(c.version, c.title, c.body, c.actor); err == nil {
			t.Fatal("invalid publication accepted")
		}
	}
}
func TestConsentNeverInferredFromMissingOrHistoricalEvidence(t *testing.T) {
	v := &Version{ID: uuid.New()}
	old := &Acceptance{TermsVersionID: uuid.New()}
	current := &Acceptance{TermsVersionID: v.ID}
	for _, c := range []struct {
		v                 *Version
		a                 *Acceptance
		previous          bool
		state             string
		allowed, required bool
	}{
		{nil, nil, false, "unpublished", false, false}, {nil, old, true, "unpublished", false, false},
		{v, nil, false, "required", false, true}, {v, nil, true, "updated", false, true},
		{v, old, true, "updated", false, true}, {v, current, true, "accepted", true, false},
	} {
		s := StatusFor(c.v, c.a, c.previous)
		if s.State != c.state || s.CanInitiateBorrowing != c.allowed || s.AcceptanceRequired != c.required {
			t.Fatal("incorrect consent state")
		}
	}
}
