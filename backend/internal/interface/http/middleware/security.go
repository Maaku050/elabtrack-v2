package middleware

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

// SecurityHeaders returns a Fiber middleware that sets common security
// headers (X-Content-Type-Options, X-Frame-Options, Referrer-Policy, etc.).
func SecurityHeaders(env config.Environment) fiber.Handler {
	return func(c fiber.Ctx) error { SetSecurityHeaders(c, env); return c.Next() }
}

// SetSecurityHeaders also runs in the error boundary for parser errors that
// occur before middleware. Document CSP belongs to the SPA serving layer.
func SetSecurityHeaders(c fiber.Ctx, env config.Environment) {
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
	if env == config.Production && authoritativeHTTPS(c) {
		c.Set("Strict-Transport-Security", "max-age=31536000")
	}
}
