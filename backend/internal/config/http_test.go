package config

import (
	"strings"
	"testing"
	"time"
)

func TestTrustedProxyConfiguration(t *testing.T) {
	cfg, err := Parse(productionValues())
	if err != nil || len(cfg.Security.TrustedProxies) != 0 {
		t.Fatal("default must use socket peer")
	}
	values := productionValues()
	values["TRUSTED_PROXIES"] = "10.20.0.2, 192.0.2.129/24, ::1, 2001:db8::/64"
	cfg, err = Parse(values)
	if err != nil || len(cfg.Security.TrustedProxies) != 4 || cfg.Security.TrustedProxies[1].String() != "192.0.2.0/24" {
		t.Fatal("literal proxy configuration not parsed")
	}
	for _, raw := range []string{"*", "private", "proxy.example.invalid", "10.0.0.1:8080", "10.0.0.0/99", "0.0.0.0/0", "::/0", "0.0.0.0", "::", "224.0.0.1", "fe80::1%eth0", "::ffff:10.0.0.0/120", "10.0.0.1,"} {
		t.Run(raw, func(t *testing.T) {
			values := productionValues()
			values["TRUSTED_PROXIES"] = raw
			_, err := Parse(values)
			if err == nil || !strings.Contains(err.Error(), "TRUSTED_PROXIES") {
				t.Fatal("invalid proxy configuration accepted")
			}
			if strings.Contains(err.Error(), raw) {
				t.Fatal("configuration value echoed")
			}
		})
	}
}

func TestCanonicalOriginPolicy(t *testing.T) {
	for raw, want := range map[string]string{"https://APP.example.invalid:443": "https://app.example.invalid", "http://localhost:080": "http://localhost", "http://[::1]:5173": "http://[::1]:5173", "https://127.0.0.1": "https://127.0.0.1"} {
		got, err := CanonicalOrigin(raw)
		if err != nil || got != want {
			t.Fatalf("canonical origin: %s", raw)
		}
	}
	for _, raw := range []string{"*", "null", "https://*.example.invalid", "https://app.example.invalid/", "https://app.example.invalid/path", "https://user:pass@app.example.invalid", "https://app.example.invalid?", "https://app.example.invalid#", "https://app.example.invalid:0", "https://app.example.invalid:65536", "https://app.example.invalid:", "ftp://app.example.invalid", " https://app.example.invalid", "https://app.example.invalid ", "https://app.example.invalid\\evil", "https://bad_host.invalid", "https://app..invalid", "https://app.example.invalid.", "http://127.1", "https://[fe80::1%25eth0]"} {
		t.Run(raw, func(t *testing.T) {
			if _, err := CanonicalOrigin(raw); err == nil {
				t.Fatal("invalid origin accepted")
			}
		})
	}
	values := productionValues()
	values["FRONTEND_URL"] = " https://APP.example.invalid:443 "
	values["ALLOWED_ORIGINS"] = " https://APP.example.invalid:443, https://other.example.invalid "
	cfg, err := Parse(values)
	if err != nil || cfg.Security.FrontendURL != "https://app.example.invalid" || cfg.Security.AllowedOrigins[0] != cfg.Security.FrontendURL {
		t.Fatal("configuration normalization inconsistent")
	}
	values["ALLOWED_ORIGINS"] = "https://other.example.invalid"
	if _, err := Parse(values); err == nil {
		t.Fatal("frontend excluded from credentials allowlist")
	}
}

func TestHTTPPerimeterSettings(t *testing.T) {
	cfg, err := Parse(map[string]string{})
	if err != nil || cfg.App.IdleTimeout != time.Minute || cfg.Security.LoginRateLimitMax != 10 || cfg.Security.RefreshRateLimitMax != 60 || cfg.Security.RegisterRateLimitMax != 5 {
		t.Fatal("perimeter defaults incorrect")
	}
	for key, value := range map[string]string{"APP_IDLE_TIMEOUT": "0s", "LOGIN_RATE_LIMIT_MAX": "0", "REFRESH_RATE_LIMIT_MAX": "bad", "REGISTER_RATE_LIMIT_MAX": "-1", "RATE_LIMIT_WINDOW": "1ms"} {
		t.Run(key, func(t *testing.T) {
			if _, err := Parse(map[string]string{key: value}); err == nil || !strings.Contains(err.Error(), key) {
				t.Fatal("invalid perimeter setting accepted")
			}
		})
	}
	if _, err := Parse(map[string]string{"RATE_LIMIT_WINDOW": "2h"}); err == nil {
		t.Fatal("unbounded limiter window accepted")
	}
}
