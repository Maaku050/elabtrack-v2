package handlers

import (
	"net/url"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

const refreshCookieName = "elabtrack_v2_refresh"
const refreshCookiePath = "/api/v1/auth"

// browserSessionPolicy has no independent environment flags capable of
// disabling production security. HTTP attributes belong at this boundary.
type browserSessionPolicy struct {
	secure, enabled bool
	origins         map[string]struct{}
}

func newBrowserSessionPolicy(env config.Environment, cfg config.SecurityConfig) browserSessionPolicy {
	p := browserSessionPolicy{secure: env != config.Development && env != config.Test, enabled: env == config.Development || env == config.Test || env == config.Production, origins: map[string]struct{}{}}
	if len(cfg.AllowedOrigins) == 0 {
		p.enabled = false
	}
	for _, origin := range cfg.AllowedOrigins {
		u, err := url.Parse(origin)
		if err != nil || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || !(u.Scheme == "https" || u.Scheme == "http") || env == config.Production && u.Scheme != "https" {
			p.enabled = false
			continue
		}
		p.origins[origin] = struct{}{}
	}
	return p
}
func (p browserSessionPolicy) allowed(c fiber.Ctx) bool {
	_, trusted := p.origins[c.Get("Origin")]
	return p.enabled && c.Method() == fiber.MethodPost && trusted
}
func (p browserSessionPolicy) cookie(raw string, expires time.Time) *fiber.Cookie {
	return &fiber.Cookie{Name: refreshCookieName, Value: raw, Path: refreshCookiePath,
		HTTPOnly: true, Secure: p.secure, SameSite: "Lax", Expires: expires}
}
func (p browserSessionPolicy) issue(c fiber.Ctx, pair auth.TokenPairDTO) {
	c.Cookie(p.cookie(pair.RefreshToken, pair.RefreshExpiresAt))
}
func (p browserSessionPolicy) clear(c fiber.Ctx) {
	cookie := p.cookie("", time.Unix(1, 0).UTC())
	cookie.MaxAge = -1
	c.Cookie(cookie)
}

// BrowserSessionDTO is an explicit access/user allowlist. The raw secret and
// persistence metadata cannot reach normal JSON even through future refactors.
type BrowserSessionDTO struct {
	AccessToken string           `json:"access_token"`
	ExpiresAt   time.Time        `json:"expires_at"`
	TokenType   string           `json:"token_type"`
	User        auth.AuthUserDTO `json:"user"`
}

func browserSession(pair auth.TokenPairDTO) BrowserSessionDTO {
	return BrowserSessionDTO{pair.AccessToken, pair.ExpiresAt, pair.TokenType, pair.User}
}

// beginBrowserSession prevents credential response caching and checks CSRF origin.
func (h *AuthHandler) beginBrowserSession(c fiber.Ctx) bool {
	c.Set("Cache-Control", "no-store")
	return h.browser.allowed(c)
}
func originDenied(c fiber.Ctx) error {
	return response.Forbidden(c, "You do not have access to this resource.")
}
