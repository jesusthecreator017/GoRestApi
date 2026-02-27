package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

const RealIPKey contextKey = "real_ip"

// Headers to check for IP in prio order
var realIPHeaders = []string{
	"CF-Connecting-IP", // CloudFlare
	"X-Real-IP",        // Nginx
	"X-Forwaded-For",   // Most other proxies/load balancers
}

func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if isTrustedProxy(req) {
			if ip := extractIP(req); ip != "" {
				// Clone the request and override the RemoteAddr
				ctx := context.WithValue(req.Context(), RealIPKey, ip)
				req.WithContext(ctx)

				req = req.Clone(req.Context())
				req.RemoteAddr = ip
			}
		}
		next.ServeHTTP(w, req)
	})
}

func extractIP(req *http.Request) string {
	for _, header := range realIPHeaders {
		value := req.Header.Get(header)
		if value == "" {
			continue
		}

		// X-Forwarded-For can contain a comma-separated list of IPs
		// The first one is the original client IP
		ip := strings.TrimSpace(strings.SplitN(value, ",", 2)[0])

		if isValidIP(ip) {
			return ip
		}
	}

	return ""
}

// Figure out whether the IP is valid
func isValidIP(ip string) bool {
	// Parse the IP
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}

	// Reject private/loopback addresses if you only want public IPs
	// Remove this block if you want to allow private IPs
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(parsed) {
			return false
		}
	}

	return true
}

func isTrustedProxy(req *http.Request) bool {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return false
	}

	trustedProxies := []string{
		"10.00.00.1",
		"192.168.1.1",
	}

	for _, proxy := range trustedProxies {
		if host == proxy {
			return true
		}
	}

	return false
}
