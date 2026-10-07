package bootstrap

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// Exercise the real Fiber/fasthttp parser callback, which app.Test does not
// render for oversized bodies. integration/runtime.py covers actual sockets.
func TestParserErrorFinalCompletion(t *testing.T) {
	var oversized fasthttp.RequestHeader
	headerError := oversized.Read(bufio.NewReaderSize(strings.NewReader("POST / HTTP/1.1\r\nX-Large: "+strings.Repeat("x", 10000)+"\r\n\r\n"), 128))
	for _, tt := range []struct {
		name   string
		err    error
		status int
		code   string
	}{
		{"body", fasthttp.ErrBodyTooLarge, 413, "PAYLOAD_TOO_LARGE"},
		{"headers", headerError, 431, "BAD_REQUEST"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			core, entries := observer.New(zap.DebugLevel)
			log := &logger.Logger{Logger: zap.New(core)}
			app := newServer(httpSecurityConfig(t), nil, log)
			app.Handler() // Initialize Fiber's underlying server without a listener.
			server := app.Server()
			server.ErrorHandler = middleware.ParserErrorLogging(app, server.ErrorHandler, log)
			var ctx fasthttp.RequestCtx
			ctx.Request.Header.SetMethod("POST")
			ctx.Request.SetRequestURI("/api/v1/parser-probe-not-a-route")
			server.ErrorHandler(&ctx, tt.err)
			var envelope map[string]any
			if json.Unmarshal(ctx.Response.Body(), &envelope) != nil || ctx.Response.StatusCode() != tt.status {
				t.Fatal("unsafe parser response")
			}
			er := envelope["error"].(map[string]any)
			id := string(ctx.Response.Header.Peek("X-Request-ID"))
			if er["code"] != tt.code || er["requestId"] != id || id == "" {
				t.Fatal("parser correlation/contract")
			}
			completed := entries.FilterMessage("http request completed").All()
			if len(completed) != 1 || completed[0].ContextMap()["status"] != int64(tt.status) || completed[0].ContextMap()["request_id"] != id {
				t.Fatal("parser completion must report final status once")
			}
		})
	}
}
