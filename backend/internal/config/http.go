package config

import (
	"errors"
	"net"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
)

// CanonicalOrigin accepts only serialized HTTP(S) origins. Configuration may
// canonicalize case/default ports; incoming browser Origin must match exactly.
func CanonicalOrigin(raw string) (string, error) {
	invalid := errors.New("invalid HTTP origin")
	u, err := url.Parse(raw)
	if err != nil || strings.TrimSpace(raw) != raw || u.Opaque != "" || u.Hostname() == "" || u.User != nil || u.Path != "" || u.RawPath != "" || u.ForceQuery || u.RawQuery != "" || u.Fragment != "" || strings.ContainsAny(raw, "?#") || (u.Scheme != "http" && u.Scheme != "https") {
		return "", invalid
	}
	host := strings.ToLower(u.Hostname())
	if ip, err := netip.ParseAddr(host); err == nil {
		if ip.Zone() != "" {
			return "", invalid
		}
		host = ip.String()
	} else {
		if len(host) > 253 {
			return "", invalid
		}
		for _, label := range strings.Split(host, ".") {
			if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
				return "", invalid
			}
			for _, c := range label {
				if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-') {
					return "", invalid
				}
			}
		}
		// Reject noncanonical/invalid numeric IPv4 spellings browsers reinterpret.
		last := host[strings.LastIndex(host, ".")+1:]
		digits := true
		for _, c := range last {
			if c < '0' || c > '9' {
				digits = false
			}
		}
		if _, err := strconv.ParseUint(last, 0, 64); err == nil || digits {
			return "", invalid
		}
	}
	port := u.Port()
	if port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return "", invalid
		}
		port = strconv.Itoa(n)
		if u.Scheme == "http" && n == 80 || u.Scheme == "https" && n == 443 {
			port = ""
		}
	} else if strings.HasSuffix(u.Host, ":") {
		return "", invalid
	}
	if port != "" {
		host = net.JoinHostPort(host, port)
	} else if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return u.Scheme + "://" + host, nil
}

func readTrustedProxies(r *reader) []netip.Prefix {
	raw := r.value("TRUSTED_PROXIES", "")
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var prefixes []netip.Prefix
	for _, value := range strings.Split(raw, ",") {
		value = strings.TrimSpace(value)
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			ip, ipErr := netip.ParseAddr(value)
			if ipErr != nil || ip.Zone() != "" {
				r.fail("TRUSTED_PROXIES", "must contain literal proxy IP addresses or CIDRs")
				continue
			}
			ip = ip.Unmap()
			prefix = netip.PrefixFrom(ip, ip.BitLen())
		}
		if prefix.Bits() == 0 || prefix.Addr().Is4In6() || prefix.Addr().IsMulticast() || prefix.Addr().Zone() != "" || prefix.Addr().IsUnspecified() && prefix.Bits() == prefix.Addr().BitLen() {
			r.fail("TRUSTED_PROXIES", "must not trust all addresses, multicast, scoped or unspecified proxy addresses")
			continue
		}
		prefixes = append(prefixes, prefix.Masked())
	}
	return prefixes
}
