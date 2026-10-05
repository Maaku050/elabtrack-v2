package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
)

// Logger returns a Fiber request-logger middleware.
// It logs method, path, status, latency and request id without ever
// logging request bodies, passwords, or tokens.
//
// For zap-backed structured logging, wire fiberzap directly in bootstrap
// when JSON logs are required in production.
func Logger() fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		Format:     "${time} ${method} ${path} ${status} ${latency} req_id=${respHeader:X-Request-ID}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
		Stream:     os.Stdout,
	})
}
