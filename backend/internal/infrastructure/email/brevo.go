package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"io"
	"net/http"
	"time"
)

type Brevo struct {
	key, email, name string
	client           *http.Client
}

func NewBrevo(key, senderEmail, senderName string) *Brevo {
	return &Brevo{key: key, email: senderEmail, name: senderName, client: &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}
func (b *Brevo) SendActivation(ctx context.Context, m d.Mail) (string, error) {
	if b.key == "" || b.email == "" {
		return "", d.ErrMailUnavailable
	}
	person := func(email, name string) map[string]string { return map[string]string{"email": email, "name": name} }
	payload := map[string]any{"sender": person(b.email, b.name), "to": []any{person(m.To, m.Name)}, "subject": "Activate your eLabTrack account", "textContent": "An FSMO administrator created your eLabTrack account. Create your own eLabTrack password using this single-use link:\n\n" + m.ActivationURL + "\n\nThis is separate from your email account password. If this was unexpected, contact FSMO. The link expires; request another from an administrator if needed."}
	raw, _ := json.Marshal(payload)
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(raw))
	if e != nil {
		return "", d.ErrMailUnknown
	}
	req.Header.Set("api-key", b.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, e := b.client.Do(req)
	if e != nil {
		return "", d.ErrMailUnknown
	}
	defer resp.Body.Close()
	body, e := io.ReadAll(io.LimitReader(resp.Body, 8193))
	if e != nil || len(body) > 8192 {
		return "", d.ErrMailUnknown
	}
	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("email provider rejected submission")
	}
	var result struct {
		ID string `json:"messageId"`
	}
	if json.Unmarshal(body, &result) != nil || result.ID == "" || len(result.ID) > 256 {
		return "", d.ErrMailUnknown
	}
	return result.ID, nil
}
