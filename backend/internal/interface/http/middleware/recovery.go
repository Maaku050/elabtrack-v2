package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
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
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e interface{}) {
			if l != nil {
				l.Error("panic recovered",
					zap.String("panic_type", fmt.Sprintf("%T", e)),
					zap.ByteString("stack", debug.Stack()),
					zap.String("client_ip", ClientIP(c)),
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
func ErrorHandler(l *logger.Logger, env config.Environment) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		if err == nil {
			return c.Next()
		}
		SetSecurityHeaders(c, env)
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) && fiberErr != nil {
			status := fiberErr.Code
			if status < 400 || status > 599 {
				status = 500
			}
			// Custom framework messages may embed infrastructure/user material.
			return response.Fail(c, status, http.StatusText(status), "FIBER_ERROR")
		}
		if l != nil {
			l.Error("request error",
				zap.String("error_type", fmt.Sprintf("%T", err)),
				zap.String("client_ip", ClientIP(c)),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("request_id", c.GetRespHeader("X-Request-ID")),
			)
		}
		return response.Error(c, err)
	}
}
