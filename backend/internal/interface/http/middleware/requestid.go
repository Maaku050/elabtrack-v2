package middleware

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/constants"
	"github.com/gofiber/fiber/v3"
	fiberrequestid "github.com/gofiber/fiber/v3/middleware/requestid"
)

// RequestID returns a Fiber middleware that ensures every request has a
// unique X-Request-ID. If the client sends one it is preserved; otherwise a
// new secure request ID is generated. The id is exposed on the response header and in
// Fiber request context for logging via requestid.FromContext.
func RequestID() fiber.Handler {
	return fiberrequestid.New(fiberrequestid.Config{
		Header:    constants.HeaderRequestID,
		Generator: nil, // use Fiber v3 secure-token generator
	})
}
