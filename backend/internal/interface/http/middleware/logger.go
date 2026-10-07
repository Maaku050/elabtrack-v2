package middleware

import (
	"io"
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
func Logger(outputs ...io.Writer) fiber.Handler {
	var stream io.Writer = os.Stdout
	if len(outputs) > 0 {
		stream = outputs[0]
	}
	return fiberlogger.New(fiberlogger.Config{
		Format:     "${time} ${method} ${path} ${status} ${latency} client_ip=${client_ip} req_id=${respHeader:X-Request-ID}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "UTC",
		Stream:     stream,
		CustomTags: map[string]fiberlogger.LogFunc{
			"client_ip": func(output fiberlogger.Buffer, c fiber.Ctx, _ *fiberlogger.Data, _ string) (int, error) {
				return output.WriteString(ClientIP(c))
			},
		},
	})
}
