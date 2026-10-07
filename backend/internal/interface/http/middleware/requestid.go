package middleware

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

// RequestID always generates server-owned UUIDv4 correlation. All inbound IDs
// are ignored (including valid-looking values); they are not trusted identity.
func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error { response.EnsureRequestID(c); return c.Next() }
}
