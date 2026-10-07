package middleware

import (
	"mime"
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

const AuthBodyLimit = 16 << 10

// RequestSafety adds a narrow auth ceiling before binding/credential work.
// Fiber retains the separately configured global ceiling for all requests.
func RequestSafety() fiber.Handler {
	return func(c fiber.Ctx) error {
		path := securityPath(c)
		if strings.HasPrefix(path, "/api/v1/auth/") && c.Method() == fiber.MethodPost && len(c.Body()) > AuthBodyLimit {
			return response.Fail(c, 413, "Request body is too large.", "PAYLOAD_TOO_LARGE")
		}
		jsonBody := c.Method() == fiber.MethodPost && (path == "/api/v1/auth/login" || path == "/api/v1/auth/register") || c.Method() == fiber.MethodPatch && path == "/api/v1/users/me"
		if jsonBody {
			media, _, err := mime.ParseMediaType(c.Get("Content-Type"))
			if err != nil || media != "application/json" {
				return response.Fail(c, 415, "JSON content type required.", "UNSUPPORTED_MEDIA_TYPE")
			}
		}
		return c.Next()
	}
}
