package email

import (
	"context"
	"encoding/json"
	"errors"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"io"
	"net/http"
	"strings"
	"testing"
)

type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestBrevoSubmission(t *testing.T) {
	b := NewBrevo("fixture-key", "sender@example.invalid", "FSMO TEST")
	calls := 0
	b.client.Transport = transport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != "POST" || r.URL.String() != "https://api.brevo.com/v3/smtp/email" || r.Header.Get("api-key") != "fixture-key" {
			t.Fatal("provider contract")
		}
		var v map[string]any
		if json.NewDecoder(r.Body).Decode(&v) != nil || !strings.Contains(v["textContent"].(string), "#token=fixture") {
			t.Fatal("activation message")
		}
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader(`{"messageId":"test-only"}`)), Header: http.Header{}}, nil
	})
	id, e := b.SendActivation(context.Background(), d.Mail{To: "recipient@example.invalid", Name: "Test", ActivationURL: "http://localhost/activate#token=fixture"})
	if e != nil || id != "test-only" || calls != 1 {
		t.Fatal("acceptance result")
	}
	b.key = ""
	if _, e = b.SendActivation(context.Background(), d.Mail{}); !errors.Is(e, d.ErrMailUnavailable) || calls != 1 {
		t.Fatal("unconfigured adapter made a call")
	}
}
func TestBrevoUnknownAndRejected(t *testing.T) {
	for _, tc := range []struct {
		code        int
		body        string
		wantUnknown bool
	}{{201, `{}`, true}, {201, strings.Repeat("a", 8193), true}, {401, `{"message":"sensitive"}`, false}, {503, `{}`, false}} {
		b := NewBrevo("test", "sender@example.invalid", "test")
		b.client.Transport = transport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: tc.code, Body: io.NopCloser(strings.NewReader(tc.body)), Header: http.Header{}}, nil
		})
		_, e := b.SendActivation(context.Background(), d.Mail{})
		if e == nil || errors.Is(e, d.ErrMailUnknown) != tc.wantUnknown || strings.Contains(e.Error(), "sensitive") {
			t.Fatal("unsafe provider outcome")
		}
	}
}
