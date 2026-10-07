package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func proxyPrefixes() []netip.Prefix {
	return []netip.Prefix{netip.MustParsePrefix("10.20.0.0/24"), netip.MustParsePrefix("2001:db8:abcd::/64")}
}

// Exercise Fiber's actual handler with an explicit socket peer; app.Test uses
// an unspecified fake peer and cannot establish proxy/IP trust by RemoteAddr.
func perimeterRequest(app *fiber.App, peer, method, path string, headers map[string]string, body string) (int, string, map[string]string) {
	var req fasthttp.Request
	req.SetRequestURI("http://api.example.invalid" + path)
	req.Header.SetMethod(method)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.SetBodyString(body)
	var ctx fasthttp.RequestCtx
	ctx.Init(&req, &net.TCPAddr{IP: net.ParseIP(peer), Port: 1234}, nil)
	app.Handler()(&ctx)
	responseHeaders := map[string]string{}
	for key, value := range ctx.Response.Header.All() {
		responseHeaders[string(key)] = string(value)
	}
	return ctx.Response.StatusCode(), string(ctx.Response.Body()), responseHeaders
}

func TestClientIPTrust(t *testing.T) {
	for _, tt := range []struct {
		name, peer, xff, proto, want string
		https                        bool
	}{
		{"direct", "198.51.100.7", "", "", "198.51.100.7", false},
		{"spoofed untrusted", "198.51.100.7", "203.0.113.8", "https", "198.51.100.7", false},
		{"trusted", "10.20.0.2", "203.0.113.8", "https", "203.0.113.8", true},
		{"untrusted hop stops spoofed prefix", "10.20.0.2", "192.0.2.55, 198.51.100.7", "http", "198.51.100.7", false},
		{"trusted chain", "10.20.0.2", "203.0.113.8, 10.20.0.3", "https", "203.0.113.8", true},
		{"malformed entire chain", "10.20.0.2", "bad, 203.0.113.8", "http", "10.20.0.2", false},
		{"empty hop", "10.20.0.2", "203.0.113.8,", "http", "10.20.0.2", false},
		{"port rejected", "10.20.0.2", "203.0.113.8:9", "http", "10.20.0.2", false},
		{"zone rejected", "10.20.0.2", "fe80::1%eth0", "http", "10.20.0.2", false},
		{"unspecified rejected", "10.20.0.2", "0.0.0.0", "http", "10.20.0.2", false},
		{"multicast rejected", "10.20.0.2", "224.0.0.1", "http", "10.20.0.2", false},
		{"all trusted", "10.20.0.2", "10.20.0.3, 10.20.0.4", "http", "10.20.0.2", false},
		{"ipv6", "2001:db8:abcd::2", "2001:db8:1::7", "https", "2001:db8:1::7", true},
		{"mapped canonical", "10.20.0.2", "::ffff:203.0.113.8", "https", "203.0.113.8", true},
		{"ambiguous scheme", "10.20.0.2", "203.0.113.8", "https, http", "203.0.113.8", false},
		{"too many hops", "10.20.0.2", strings.Repeat("203.0.113.8,", 16) + "203.0.113.8", "http", "10.20.0.2", false},
		{"oversized header", "10.20.0.2", strings.Repeat("x", 1025), "http", "10.20.0.2", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(ClientInfo(config.SecurityConfig{TrustedProxies: proxyPrefixes()}))
			app.Get("/probe", func(c fiber.Ctx) error {
				return c.JSON(fiber.Map{"ip": ClientIP(c), "https": authoritativeHTTPS(c), "scheme": c.Scheme(), "host": c.Hostname()})
			})
			_, body, _ := perimeterRequest(app, tt.peer, "GET", "/probe", map[string]string{"X-Forwarded-For": tt.xff, "X-Forwarded-Proto": tt.proto, "X-Real-IP": "192.0.2.1", "Forwarded": "for=192.0.2.1;proto=https", "X-Forwarded-Host": "attacker.invalid", "X-Forwarded-Ssl": "on"}, "")
			var got struct {
				IP     string `json:"ip"`
				HTTPS  bool   `json:"https"`
				Scheme string `json:"scheme"`
				Host   string `json:"host"`
			}
			if json.Unmarshal([]byte(body), &got) != nil || got.IP != tt.want || got.HTTPS != tt.https || got.Scheme != "http" || got.Host != "api.example.invalid" {
				t.Fatal("forwarding trust boundary failed")
			}
		})
	}
	if got := resolveClient(netip.MustParseAddr("10.20.0.2"), "203.0.113.8", "https", false, nil); got.ip != "10.20.0.2" || got.https {
		t.Fatal("empty allowlist trusted proxy")
	}
	if !resolveClient(netip.MustParseAddr("198.51.100.7"), "", "", true, nil).https {
		t.Fatal("actual TLS lost")
	}
}

func limiterApp(cfg config.SecurityConfig) *fiber.App {
	app := fiber.New()
	app.Use(ClientInfo(cfg))
	app.Use(RateLimit(cfg))
	app.Use(func(c fiber.Ctx) error { return c.SendStatus(204) })
	return app
}
func TestDifferentiatedRateLimits(t *testing.T) {
	cfg := config.SecurityConfig{RateLimitMax: 2, LoginRateLimitMax: 3, RefreshRateLimitMax: 12, RegisterRateLimitMax: 2, RateLimitWindow: time.Hour, TrustedProxies: proxyPrefixes()}
	app := limiterApp(cfg)
	for i := 0; i < 3; i++ {
		status, _, _ := perimeterRequest(app, "198.51.100.7", "POST", "/api/v1/auth/login", nil, "")
		if status != 204 {
			t.Fatal("login below threshold denied")
		}
	}
	status, body, headers := perimeterRequest(app, "198.51.100.7", "POST", "/API/V1/AUTH/LOGIN/", map[string]string{"X-Forwarded-For": "192.0.2.123"}, "")
	if status != 429 || headers["Retry-After"] != "3600" || headers["Cache-Control"] != "no-store" || strings.Contains(body, "198.51") || headers["X-Ratelimit-Limit"] != "" {
		t.Fatal("unsafe or unstable limit response")
	}
	var envelope struct {
		Success bool
		Error   struct{ Code string }
	}
	if json.Unmarshal([]byte(body), &envelope) != nil || envelope.Success || envelope.Error.Code != "RATE_LIMITED" {
		t.Fatal("429 envelope")
	}
	status, _, _ = perimeterRequest(app, "198.51.100.8", "POST", "/api/v1/auth/login", nil, "")
	if status != 204 {
		t.Fatal("different direct keys share limit")
	}
	for _, client := range []string{"203.0.113.8", "203.0.113.9"} {
		for i := 0; i < 3; i++ {
			status, _, _ = perimeterRequest(app, "10.20.0.2", "POST", "/api/v1/auth/login", map[string]string{"X-Forwarded-For": client}, "")
			if status != 204 {
				t.Fatal("trusted forwarded keys share limit")
			}
		}
	}
	for i := 0; i < 12; i++ {
		status, _, _ = perimeterRequest(app, "198.51.100.7", "POST", "/api/v1/auth/refresh", nil, "")
		if status != 204 {
			t.Fatal("refresh ordinary budget denied")
		}
	}
	status, _, _ = perimeterRequest(app, "198.51.100.7", "POST", "/api/v1/auth/refresh", nil, "")
	if status != 429 {
		t.Fatal("refresh has no abuse limit")
	}
	for _, path := range []string{"/api/v1/auth/register", "/api/v1/auth/logout", "/api/v1/users/me"} {
		isolated := limiterApp(cfg)
		for i := 0; i < 2; i++ {
			status, _, _ = perimeterRequest(isolated, "198.51.100.7", "POST", path, nil, "")
			if status != 204 {
				t.Fatal("ordinary requests denied")
			}
		}
		status, _, _ = perimeterRequest(isolated, "198.51.100.7", "POST", path, nil, "")
		if status != 429 {
			t.Fatal("public/general budget missing")
		}
	}
	for i := 0; i < 5; i++ {
		status, _, _ = perimeterRequest(app, "198.51.100.7", "GET", "/api/v1/health", nil, "")
		want := 204
		if i >= 2 {
			want = 429
		}
		if status != want {
			t.Fatal("public health must share the general budget")
		}
	}
}

func TestRateLimitConcurrentEnforcement(t *testing.T) {
	app := limiterApp(config.SecurityConfig{LoginRateLimitMax: 10, RateLimitWindow: time.Hour})
	var wg sync.WaitGroup
	results := make(chan int, 32)
	for i := 0; i < 32; i++ {
		wg.Go(func() {
			status, _, _ := perimeterRequest(app, "198.51.100.7", "POST", "/api/v1/auth/login", nil, "")
			results <- status
		})
	}
	wg.Wait()
	close(results)
	allowed, denied := 0, 0
	for status := range results {
		if status == 204 {
			allowed++
		} else if status == 429 {
			denied++
		} else {
			t.Fatal("unexpected status")
		}
	}
	if allowed != 10 || denied != 22 {
		t.Fatal("concurrent requests escaped limit")
	}
}

func TestExplicitCredentialedCORS(t *testing.T) {
	app := fiber.New()
	app.Use(CORS(config.SecurityConfig{AllowedOrigins: []string{"http://localhost:5173", "https://app.example.invalid"}}))
	app.Use(func(c fiber.Ctx) error { return c.SendStatus(200) })
	for _, origin := range []string{"http://localhost:5173", "https://app.example.invalid"} {
		status, _, h := perimeterRequest(app, "198.51.100.7", "OPTIONS", "/api/v1/auth/refresh", map[string]string{"Origin": origin, "Access-Control-Request-Method": "POST", "Access-Control-Request-Headers": "content-type, AUTHORIZATION, X-Request-ID"}, "")
		if status != 204 || h["Access-Control-Allow-Origin"] != origin || h["Access-Control-Allow-Credentials"] != "true" || h["Access-Control-Max-Age"] != "300" || h["Access-Control-Allow-Methods"] != "GET, POST, PATCH, OPTIONS" || !strings.Contains(h["Vary"], "Access-Control-Request-Headers") {
			t.Fatal("credentialed preflight invalid")
		}
	}
	for _, origin := range []string{"https://attacker.invalid", "null", "*", "https://app.example.invalid.attacker.invalid", "https://app.example.invalid/", "HTTPS://APP.example.invalid", " https://app.example.invalid", "https://app.example.invalid:443", "http://app.example.invalid"} {
		status, _, h := perimeterRequest(app, "198.51.100.7", "POST", "/api/v1/auth/refresh", map[string]string{"Origin": origin}, "")
		if status != 403 || h["Access-Control-Allow-Origin"] != "" || h["Access-Control-Allow-Credentials"] != "" {
			t.Fatal("arbitrary or noncanonical origin accepted")
		}
	}
	for _, headers := range []map[string]string{
		{"Access-Control-Request-Method": "DELETE"}, {"Access-Control-Request-Method": "PUT"}, {"Access-Control-Request-Method": "TRACE"}, {"Access-Control-Request-Method": "POST", "Access-Control-Request-Headers": "X-Admin"},
	} {
		headers["Origin"] = "https://app.example.invalid"
		status, _, h := perimeterRequest(app, "198.51.100.7", "OPTIONS", "/api/v1/auth/refresh", headers, "")
		if status != 403 || h["Access-Control-Allow-Origin"] != "" {
			t.Fatal("unneeded method/header preflight allowed")
		}
	}
	status, _, h := perimeterRequest(app, "198.51.100.7", "GET", "/api/v1/health", nil, "")
	if status != 200 || h["Access-Control-Allow-Origin"] != "" || h["Vary"] != "Origin" {
		t.Fatal("nonbrowser CORS behavior")
	}
}

func TestSecurityHeadersTransportAuthority(t *testing.T) {
	for _, tt := range []struct {
		name        string
		env         config.Environment
		peer, proto string
		wantHSTS    bool
	}{
		{"development", config.Development, "10.20.0.2", "https", false},
		{"test", config.Test, "10.20.0.2", "https", false},
		{"production HTTP", config.Production, "10.20.0.2", "http", false},
		{"production spoof", config.Production, "198.51.100.7", "https", false},
		{"production HTTPS", config.Production, "10.20.0.2", "https", true},
		{"ambiguous production", config.Production, "10.20.0.2", "https, http", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(ClientInfo(config.SecurityConfig{TrustedProxies: proxyPrefixes()}))
			app.Use(SecurityHeaders(tt.env))
			app.Get("/", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"ok": true}) })
			_, _, h := perimeterRequest(app, tt.peer, "GET", "/", map[string]string{"X-Forwarded-Proto": tt.proto}, "")
			if h["X-Content-Type-Options"] != "nosniff" || h["X-Frame-Options"] != "DENY" || h["Referrer-Policy"] != "no-referrer" || h["Permissions-Policy"] == "" || (h["Strict-Transport-Security"] != "") != tt.wantHSTS || h["Content-Security-Policy"] != "" || h["Cross-Origin-Resource-Policy"] != "" {
				t.Fatal("security headers wrong for transport")
			}
		})
	}
}

func TestRepeatedForwardingHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.RequestCtx().Request.Header.Add("X-Forwarded-For", "198.51.100.7")
		c.RequestCtx().Request.Header.Add("X-Forwarded-Proto", "https")
		return c.Next()
	})
	app.Use(ClientInfo(config.SecurityConfig{TrustedProxies: proxyPrefixes()}))
	app.Get("/", func(c fiber.Ctx) error { return c.JSON(fiber.Map{"ip": ClientIP(c), "https": authoritativeHTTPS(c)}) })
	_, body, _ := perimeterRequest(app, "10.20.0.2", "GET", "/", map[string]string{"X-Forwarded-For": "192.0.2.55", "X-Forwarded-Proto": "http"}, "")
	var got struct {
		IP    string `json:"ip"`
		HTTPS bool   `json:"https"`
	}
	if json.Unmarshal([]byte(body), &got) != nil || got.IP != "198.51.100.7" || got.HTTPS {
		t.Fatal("duplicate forwarding fields became authoritative")
	}
}

func TestRequestLogSecretSafety(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := &logger.Logger{Logger: zap.New(core)}
	app := fiber.New()
	app.Use(ClientInfo(config.SecurityConfig{TrustedProxies: proxyPrefixes()}))
	app.Use(Logger(log))
	app.Post("/api/v1/auth/login", func(c fiber.Ctx) error { return c.SendStatus(204) })
	const secret = "synthetic-private-credential-sentinel"
	status, _, _ := perimeterRequest(app, "10.20.0.2", "POST", "/api/v1/auth/login", map[string]string{"Authorization": "Bearer " + secret, "Cookie": "elabtrack_v2_refresh=" + secret, "X-Forwarded-For": "203.0.113.8"}, `{"password":"`+secret+`"}`)
	if status != 204 || logs.Len() != 1 || strings.Contains(fmt.Sprint(logs.All()), secret) || logs.All()[0].ContextMap()["client_ip"] != "203.0.113.8" {
		t.Fatal("request log leaked credentials or ignored shared IP")
	}
}
