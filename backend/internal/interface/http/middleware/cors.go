package middleware

import (
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

// CORS returns a CORS middleware configured from the application config.
// Only the configured allowlist of origins is permitted.
func CORS(cfg config.SecurityConfig) fiber.Handler {
	allowed := map[string]bool{}
	valid := len(cfg.AllowedOrigins) > 0
	for _, origin := range cfg.AllowedOrigins {
		canonical, err := config.CanonicalOrigin(origin)
		if err != nil || canonical != origin {
			valid = false
		}
		allowed[origin] = true
	}
	methods := map[string]bool{"GET": true, "POST": true, "PATCH": true, "OPTIONS": true}
	headers := map[string]bool{"content-type": true, "accept": true, "authorization": true, "x-request-id": true}
	return func(c fiber.Ctx) error {
		c.Vary("Origin")
		origin := c.Get("Origin")
		if origin == "" {
			return c.Next()
		} // Auth handlers separately deny missing Origin.
		if !valid || !allowed[origin] {
			return response.Forbidden(c, "Origin is not allowed.")
		}
		if c.Method() == fiber.MethodOptions && c.Get("Access-Control-Request-Method") != "" {
			c.Vary("Access-Control-Request-Method", "Access-Control-Request-Headers")
			if !methods[c.Get("Access-Control-Request-Method")] {
				return response.Forbidden(c, "CORS request is not allowed.")
			}
			if requested := c.Get("Access-Control-Request-Headers"); requested != "" {
				for _, header := range strings.Split(requested, ",") {
					if !headers[strings.ToLower(strings.TrimSpace(header))] {
						return response.Forbidden(c, "CORS request is not allowed.")
					}
				}
			}
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, X-Request-ID")
			c.Set("Access-Control-Max-Age", "300")
			return c.SendStatus(204)
		}
		c.Set("Access-Control-Allow-Origin", origin)
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Expose-Headers", "X-Request-ID, Retry-After")
		return c.Next()
	}
}
