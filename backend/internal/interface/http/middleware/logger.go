package middleware

import (
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// requestFields is an allowlist: route templates, never raw URLs/query values,
// bodies, headers, cookies, claims, account names/emails or config dumps.
func requestFields(c fiber.Ctx) []zap.Field {
	route := "unmatched"
	if actual := c.Route(); actual != nil && actual.Method != "USE" && actual.Path != "/" {
		route = actual.Path
	}
	return []zap.Field{zap.String("request_id", response.EnsureRequestID(c)), zap.String("method", c.Method()), zap.String("route", route), zap.String("client_ip", ClientIP(c))}
}

// Logger completes error rendering before observing status and duration.
// Zap's injectable core supports tests and the selected console/JSON encoder.
func Logger(l *logger.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			err = c.App().ErrorHandler(c, err)
		}
		if l == nil {
			return err
		}
		status := c.Response().StatusCode()
		out := response.OutcomeFromContext(c)
		fields := append(requestFields(c), zap.String("event", "http.request_completed"), zap.Int("status", status), zap.Float64("duration_ms", float64(time.Since(start))/float64(time.Millisecond)))
		if out.Code != "" {
			fields = append(fields, zap.String("error_code", out.Code))
		}
		event := ""
		route := c.Route()
		if route != nil && c.Method() == "POST" {
			switch route.Path {
			case "/api/v1/auth/login":
				if status < 300 {
					event = "auth.login_succeeded"
				} else {
					event = "auth.login_failed"
				}
			case "/api/v1/auth/refresh":
				if status < 300 {
					event = "auth.refresh_succeeded"
				} else {
					event = "auth.refresh_failed"
				}
			case "/api/v1/auth/logout":
				if status < 300 {
					event = "auth.logout_completed"
				} else {
					event = "auth.logout_failed"
				}
			}
		}
		if event != "" {
			fields = append(fields, zap.String("security_event", event))
		}
		switch {
		case status >= 500 || out.Unexpected:
			l.Error("http request completed", fields...)
		case status == 403 || status == 429:
			l.Warn("http request completed", fields...)
		case status >= 400 && (status != 401 || event == ""):
			l.Debug("http request completed", fields...)
		default:
			l.Info("http request completed", fields...)
		}
		// response.Error writes directly; it still receives one safe failure record.
		logFailure(l, c)
		return err
	}
}

type failureLoggedKey struct{}

func logFailure(l *logger.Logger, c fiber.Ctx) {
	out := response.OutcomeFromContext(c)
	if l == nil || !out.Unexpected || c.Locals(failureLoggedKey{}) == true || c.Locals(panicKey{}) == true {
		return
	}
	c.Locals(failureLoggedKey{}, true)
	fields := append(requestFields(c), zap.String("event", "http.request_failed"), zap.Int("status", c.Response().StatusCode()), zap.String("error_code", out.Code), zap.String("error_class", out.ErrorClass), zap.String("operation", out.Operation))
	l.Error("request error", fields...)
}
