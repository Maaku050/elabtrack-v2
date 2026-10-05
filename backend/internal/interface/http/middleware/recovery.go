package middleware

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	fiberrecover "github.com/gofiber/fiber/v3/middleware/recover"
	"go.uber.org/zap"
)

// Recovery returns a Fiber middleware that recovers from panics so a single
// panicking handler never crashes the process. The panic is logged with the
// request id; the client receives a generic 500 response.
func Recovery(l *logger.Logger) fiber.Handler {
	return fiberrecover.New(fiberrecover.Config{
		EnableStackTrace: false,
		StackTraceHandler: func(c fiber.Ctx, e interface{}) {
			if l != nil {
				l.Error("panic recovered",
					zap.Any("error", e),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.String("request_id", c.GetRespHeader("X-Request-ID")),
				)
			}
		},
	})
}

// ErrorHandler is the centralized Fiber error handler. It maps errors
// returned from handlers into the standard API response shape.
func ErrorHandler(l *logger.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		if err == nil {
			return c.Next()
		}
		// fiber.Err* are sent with their intended status by Fiber; route them
		// through our response shape for consistency.
		if fiberErr, ok := err.(*fiber.Error); ok {
			return response.Fail(c, fiberErr.Code, fiberErr.Message, "FIBER_ERROR")
		}
		if l != nil {
			l.Error("request error",
				zap.Error(err),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("request_id", c.GetRespHeader("X-Request-ID")),
			)
		}
		return response.Error(c, err)
	}
}
