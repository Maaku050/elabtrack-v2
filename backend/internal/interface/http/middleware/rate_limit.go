package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/constants"
	"github.com/gofiber/fiber/v3"
	fiberlimiter "github.com/gofiber/fiber/v3/middleware/limiter"
)

// RateLimit uses separate process-local, expiring fixed-window IP buckets.
// It is not distributed enforcement or account identity. Successful and failed
// attempts both count; no request body/token/email enters a key.
func RateLimit(cfg config.SecurityConfig) fiber.Handler {
	window := cfg.RateLimitWindow
	if window < time.Second {
		window = time.Minute
	}
	policy := func(max, fallback int) fiber.Handler {
		if max <= 0 {
			max = fallback
		}
		return fiberlimiter.New(fiberlimiter.Config{
			Max: max, Expiration: window, DisableHeaders: true,
			KeyGenerator: ClientIP,
			LimitReached: func(c fiber.Ctx) error {
				c.Set("Cache-Control", "no-store")
				// A full window is a conservative delay; do not expose bucket counts.
				c.Set("Retry-After", strconv.FormatInt(int64((window+time.Second-1)/time.Second), 10))
				return response.Fail(c, 429, "Too many requests. Please slow down.", constants.CodeRateLimited)
			},
		})
	}
	general, login, refresh, register := policy(cfg.RateLimitMax, 120), policy(cfg.LoginRateLimitMax, 10), policy(cfg.RefreshRateLimitMax, 60), policy(cfg.RegisterRateLimitMax, 5)
	return func(c fiber.Ctx) error {
		path := securityPath(c)
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}
		if c.Method() == fiber.MethodPost {
			switch path {
			case "/api/v1/auth/login", "/api/v1/auth/activate":
				return login(c)
			case "/api/v1/auth/refresh":
				return refresh(c)
			case "/api/v1/auth/register":
				return register(c)
			}
		}
		return general(c)
	}
}

// Match Fiber's case-insensitive/trailing-slash routing in perimeter checks.
func securityPath(c fiber.Ctx) string { return strings.ToLower(strings.TrimRight(c.Path(), "/")) }
