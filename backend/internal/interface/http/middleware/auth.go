package middleware

import (
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/constants"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type identityKey struct{}
type principalKey struct{}

// Auth is the single protected-request boundary: verify bearer identity, then
// resolve the current permitted account. JWT email/role never authorize a route.
func Auth(issuer application.TokenIssuer, accounts appauth.AccountResolver) fiber.Handler {
	return func(c fiber.Ctx) error {
		parts := strings.Fields(c.Get(constants.HeaderAuthorization))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Unauthorized(c, "Authentication required.")
		}
		if issuer == nil || accounts == nil {
			return response.Internal(c, "Unable to process request.")
		}
		claims, err := issuer.VerifyAccessToken(c.Context(), parts[1])
		if err != nil || claims.UserID == uuid.Nil {
			return response.Unauthorized(c, "Authentication required.")
		}
		identity := appauth.Identity{UserID: claims.UserID}
		principal, err := accounts.ResolveCurrentAccount(c.Context(), identity)
		if err != nil {
			return response.Error(c, err)
		}
		// Defense against a miswired resolver; only the verified account may proceed.
		if principal.ID != identity.UserID {
			return response.Unauthorized(c, "Authentication required.")
		}
		if !principal.IsActive || !principal.Role.Valid() {
			return response.Forbidden(c, "You do not have access to this resource.")
		}
		c.Locals(identityKey{}, identity)
		c.Locals(principalKey{}, principal)
		return c.Next()
	}
}

// RequireRole uses current database state and only temporary known roles.
// Install after Auth. Unknown configured/current roles cannot authorize access.
func RequireRole(roles ...domainuser.Role) fiber.Handler {
	allowed := map[domainuser.Role]bool{}
	for _, role := range roles {
		if role.Valid() {
			allowed[role] = true
		}
	}
	return func(c fiber.Ctx) error {
		principal, ok := PrincipalFromContext(c)
		if !ok {
			return response.Unauthorized(c, "Authentication required.")
		}
		if !principal.IsActive || !principal.Role.Valid() || !allowed[principal.Role] {
			return response.Forbidden(c, "You do not have access to this resource.")
		}
		return c.Next()
	}
}

func PrincipalFromContext(c fiber.Ctx) (appauth.Principal, bool) {
	principal, ok := c.Locals(principalKey{}).(appauth.Principal)
	return principal, ok && principal.ID != uuid.Nil
}
func IdentityFromContext(c fiber.Ctx) (appauth.Identity, bool) {
	identity, ok := c.Locals(identityKey{}).(appauth.Identity)
	return identity, ok && identity.UserID != uuid.Nil
}
