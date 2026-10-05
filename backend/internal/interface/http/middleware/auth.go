package middleware

import (
	"strings"

	"github.com/fullstacktemplate/backend/internal/application"
	"github.com/fullstacktemplate/backend/internal/interface/http/response"
	"github.com/fullstacktemplate/backend/internal/shared/constants"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Auth returns a Fiber middleware that validates a Bearer access token via
// the application.TokenIssuer port and stores the resulting claims in the
// Fiber context for downstream handlers.
func Auth(issuer application.TokenIssuer) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get(constants.HeaderAuthorization)
		if header == "" {
			return response.Unauthorized(c, "Missing authorization header.")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			return response.Unauthorized(c, "Invalid authorization header.")
		}

		claims, err := issuer.VerifyAccessToken(c.Context(), parts[1])
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired access token.")
		}

		c.Locals(constants.CtxClaims, claims)
		c.Locals(constants.CtxUserID, claims.UserID)
		return c.Next()
	}
}

// RequireRole returns a middleware that allows only the given roles.
// It must be installed after Auth so claims are available.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals(constants.CtxClaims).(application.Claims)
		if !ok {
			return response.Unauthorized(c, "Authentication required.")
		}
		if _, ok := allowed[claims.Role]; !ok {
			return response.Forbidden(c, "Insufficient permissions.")
		}
		return c.Next()
	}
}

// ClaimsFromContext returns the authenticated claims from Fiber locals.
func ClaimsFromContext(c fiber.Ctx) (application.Claims, bool) {
	claims, ok := c.Locals(constants.CtxClaims).(application.Claims)
	return claims, ok
}

// UserIDFromContext returns the authenticated user id from Fiber locals.
func UserIDFromContext(c fiber.Ctx) (uuid.UUID, bool) {
	id, ok := c.Locals(constants.CtxUserID).(uuid.UUID)
	return id, ok
}
