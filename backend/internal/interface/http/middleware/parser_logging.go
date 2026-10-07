package middleware

import (
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/observability"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type parserErrorKey struct{}

// ParserErrorLogging wraps Fiber's existing parser callback. Keep its error
// mapping, middleware traversal and response rendering, but emit completion
// only after it has written the final response. No raw parser error is logged.
func ParserErrorLogging(app *fiber.App, next func(*fasthttp.RequestCtx, error), log *logger.Logger) func(*fasthttp.RequestCtx, error) {
	return func(ctx *fasthttp.RequestCtx, err error) {
		start := time.Now()
		ctx.SetUserValue(parserErrorKey{}, true)
		next(ctx, err)
		if log == nil {
			return
		}
		c := app.AcquireCtx(ctx)
		defer app.ReleaseCtx(c)
		// ReleaseCtx clears Fiber's standard user context. Preserve the ID already
		// emitted by the parser response, rather than generating a second one.
		id := string(ctx.Response.Header.Peek("X-Request-ID"))
		c.SetContext(observability.WithRequestID(c.Context(), id))
		completeRequest(log, c, start)
	}
}
