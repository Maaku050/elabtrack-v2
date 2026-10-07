package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type panicKey struct{}

// Recovery sits inside the completion logger. Any panic becomes INTERNAL_ERROR,
// even a panic carrying a framework/domain error. Values are never formatted.
func Recovery(l *logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if value := recover(); value != nil {
				c.Locals(panicKey{}, true)
				if l != nil {
					fields := append(requestFields(c), zap.String("event", "http.panic_recovered"), zap.String("panic_type", fmt.Sprintf("%T", value)), zap.ByteString("stack", debug.Stack()))
					l.Error("panic recovered", fields...)
				}
				err = shared.ErrInternal
			}
		}()
		return c.Next()
	}
}

// ErrorHandler covers returned domain/application/framework errors and parser
// failures; the same mapper is used by direct handler response.Error calls.
func ErrorHandler(l *logger.Logger, env config.Environment) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		SetSecurityHeaders(c, env)
		result := response.Error(c, err)
		logFailure(l, c)
		return result
	}
}
