package middleware

import (
	"net/netip"
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

type clientInfoKey struct{}
type clientInfo struct {
	ip    string
	https bool
}

// ClientInfo is the sole forwarding-header boundary. Fiber proxy parsing stays
// disabled: Host/Scheme are never inferred from arbitrary alternate headers.
func ClientInfo(cfg config.SecurityConfig) fiber.Handler {
	return func(c fiber.Ctx) error {
		peer, _ := netip.AddrFromSlice(c.RequestCtx().RemoteIP())
		peer = peer.Unmap()
		// Reject duplicated XFP; concatenate all XFF field lines before parsing.
		var forwarded []string
		proto := ""
		protoCount := 0
		for key, value := range c.RequestCtx().Request.Header.All() {
			if strings.EqualFold(string(key), "X-Forwarded-For") {
				forwarded = append(forwarded, string(value))
			}
			if strings.EqualFold(string(key), "X-Forwarded-Proto") {
				proto = string(value)
				protoCount++
			}
		}
		if protoCount != 1 {
			proto = ""
		}
		info := resolveClient(peer, strings.Join(forwarded, ","), proto, c.RequestCtx().IsTLS(), cfg.TrustedProxies)
		c.Locals(clientInfoKey{}, info)
		return c.Next()
	}
}

func trustedPeer(ip netip.Addr, proxies []netip.Prefix) bool {
	for _, prefix := range proxies {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func resolveClient(peer netip.Addr, xff, proto string, tls bool, proxies []netip.Prefix) clientInfo {
	info := clientInfo{ip: "unknown", https: tls}
	if !peer.IsValid() {
		return info
	}
	peer = peer.Unmap()
	info.ip = peer.String()
	if !trustedPeer(peer, proxies) {
		return info
	}
	// A trusted immediate proxy must overwrite XFP from authoritative transport.
	info.https = tls || proto == "https"
	if xff == "" || len(xff) > 1024 {
		return info
	}
	parts := strings.Split(xff, ",")
	if len(parts) > 16 {
		return info
	}
	ips := make([]netip.Addr, len(parts))
	for i, part := range parts {
		ip, err := netip.ParseAddr(strings.TrimSpace(part))
		if err != nil || ip.Zone() != "" || ip.IsUnspecified() || ip.IsMulticast() {
			return info
		}
		ips[i] = ip.Unmap()
	}
	// Walk back through explicitly trusted proxy hops; never take an attacker-
	// supplied leftmost value beyond the first untrusted hop. All-trusted falls
	// back to the peer, as does any malformed/ambiguous forwarding chain.
	for i := len(ips) - 1; i >= 0; i-- {
		if !trustedPeer(ips[i], proxies) {
			info.ip = ips[i].String()
			break
		}
	}
	return info
}

// ClientIP is an abuse-control signal, never an authenticated account identity.
// Without the middleware, fail back to the actual socket, never a header.
func ClientIP(c fiber.Ctx) string {
	if info, ok := c.Locals(clientInfoKey{}).(clientInfo); ok {
		return info.ip
	}
	return c.RequestCtx().RemoteIP().String()
}
func authoritativeHTTPS(c fiber.Ctx) bool {
	if info, ok := c.Locals(clientInfoKey{}).(clientInfo); ok {
		return info.https
	}
	return c.RequestCtx().IsTLS()
}
