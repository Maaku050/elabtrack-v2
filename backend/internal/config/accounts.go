package config

import (
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

type AccountsConfig struct {
	Policy                            d.Policy
	BrevoKey, SenderEmail, SenderName string
}

func parseAccounts(r *reader, env Environment, frontend string) AccountsConfig {
	c := AccountsConfig{Policy: d.Policy{ActivationTTL: r.duration("ACTIVATION_TTL", 24*time.Hour), ResendCooldown: r.duration("ACTIVATION_RESEND_COOLDOWN", time.Minute), ActivationURL: r.value("ACTIVATION_URL", frontend+"/activate")}, BrevoKey: r.value("BREVO_API_KEY", ""), SenderEmail: r.value("BREVO_SENDER_EMAIL", ""), SenderName: r.value("BREVO_SENDER_NAME", "FSMO eLabTrack")}
	if c.Policy.ActivationTTL < 15*time.Minute || c.Policy.ActivationTTL > 72*time.Hour {
		r.fail("ACTIVATION_TTL", "must be between 15m and 72h")
	}
	if c.Policy.ResendCooldown < time.Minute || c.Policy.ResendCooldown > time.Hour {
		r.fail("ACTIVATION_RESEND_COOLDOWN", "must be between 1m and 1h")
	}
	u, e := url.Parse(c.Policy.ActivationURL)
	if e != nil || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "/activate" || u.Scheme+"://"+u.Host != frontend {
		r.fail("ACTIVATION_URL", "must be /activate on FRONTEND_URL without credentials, query or fragment")
	}
	for _, domain := range strings.Split(r.value("STUDENT_EMAIL_DOMAINS", ""), ",") {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		valid := len(domain) <= 253 && strings.Contains(domain, ".")
		for _, part := range strings.Split(domain, ".") {
			if len(part) == 0 || len(part) > 63 || part[0] == '-' || part[len(part)-1] == '-' {
				valid = false
			}
			for _, ch := range part {
				if !(ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '-') {
					valid = false
				}
			}
		}
		if !valid {
			r.fail("STUDENT_EMAIL_DOMAINS", "must contain exact ASCII email domain names")
		}
		c.Policy.StudentDomains = append(c.Policy.StudentDomains, domain)
	}
	if c.SenderEmail != "" {
		a, e := mail.ParseAddress(c.SenderEmail)
		if e != nil || a.Address != c.SenderEmail {
			r.fail("BREVO_SENDER_EMAIL", "must be a single sender email address")
		}
	}
	if len(c.SenderName) > 100 || strings.ContainsAny(c.SenderName, "\r\n") {
		r.fail("BREVO_SENDER_NAME", "must be at most 100 bytes without line breaks")
	}
	// Missing provider credentials are a readiness state, not an API startup failure.
	return c
}
