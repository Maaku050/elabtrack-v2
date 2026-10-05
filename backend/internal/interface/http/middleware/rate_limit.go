package middleware

import (
	"time"

	"github.com/fullstacktemplate/backend/internal/config"
	"github.com/fullstacktemplate/backend/internal/interface/http/response"
	"github.com/fullstacktemplate/backend/internal/shared/constants"
	"github.com/gofiber/fiber/v3"
	fiberlimiter "github.com/gofiber/fiber/v3/middleware/limiter"
)

// RateLimit returns a Fiber middleware that limits requests per IP.
// On limit exceeded it returns a 429 with the standard response shape.
func RateLimit(cfg config.SecurityConfig) fiber.Handler {
	return fiberlimiter.New(fiberlimiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: cfg.RateLimitWindow,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return response.Fail(c, 429, "Too many requests. Please slow down.", constants.CodeRateLimited)
		},
		Next: func(c fiber.Ctx) bool {
			// Don't rate-limit the health check.
			return c.Path() == "/api/v1/health"
		},
	})
}

// rateLimitWindow is exported for tests that need to know the configured window.
func rateLimitWindow(cfg config.SecurityConfig) time.Duration { return cfg.RateLimitWindow }
