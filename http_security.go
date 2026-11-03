package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type httpSecurityConfig struct {
	tokenHash      [32]byte
	allowedOrigins map[string]bool
}

func sha256Hex(value string) [32]byte {
	return sha256.Sum256([]byte(value))
}

func parseAllowedOrigins(originsEnv string) (map[string]bool, error) {
	allowed := map[string]bool{}
	trimmed := strings.TrimSpace(originsEnv)
	if trimmed == "" {
		return allowed, nil
	}

	for _, raw := range strings.Split(originsEnv, ",") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			return nil, fmt.Errorf("allowed origins contain empty entry")
		}
		origin, err := normalizeOrigin(entry)
		if err != nil {
			return nil, err
		}
		allowed[origin] = true
	}

	return allowed, nil
}

func buildAllowedOrigins(originsEnv string, port string) (map[string]bool, error) {
	allowed, err := parseAllowedOrigins(originsEnv)
	if err != nil {
		return nil, err
	}
	for _, host := range []string{"localhost", "127.0.0.1", "::1"} {
		origin, err := normalizeOrigin("http://" + net.JoinHostPort(host, port))
		if err != nil {
			return nil, err
		}
		allowed[origin] = true
	}
	return allowed, nil
}

func newHTTPAuthOriginMiddleware(next http.Handler, cfg httpSecurityConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationValues := r.Header.Values("Authorization")
		if len(authorizationValues) != 1 {
			writeUnauthorized(w)
			return
		}

		providedToken, ok := bearerTokenFromAuthorization(authorizationValues[0])
		if !ok {
			writeUnauthorized(w)
			return
		}

		providedTokenHash := sha256.Sum256([]byte(providedToken))
		if subtle.ConstantTimeCompare(cfg.tokenHash[:], providedTokenHash[:]) != 1 {
			writeUnauthorized(w)
			return
		}

		originValues := r.Header.Values("Origin")
		if len(originValues) > 1 {
			http.Error(w, "Invalid origin", http.StatusForbidden)
			return
		}
		if len(originValues) == 1 {
			normalized, err := normalizeOrigin(originValues[0])
			if err != nil || !cfg.allowedOrigins[normalized] {
				http.Error(w, "Invalid origin", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func bearerTokenFromAuthorization(authorization string) (string, bool) {
	scheme, token, found := strings.Cut(strings.TrimSpace(authorization), " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t") {
		return "", false
	}
	return token, true
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func normalizeOrigin(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u == nil {
		return "", fmt.Errorf("invalid origin")
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid origin")
	}
	if u.User != nil {
		return "", fmt.Errorf("origin must not include userinfo")
	}
	if u.Path != "" {
		return "", fmt.Errorf("origin must not include path")
	}
	if u.RawQuery != "" || u.ForceQuery {
		return "", fmt.Errorf("origin must not include query")
	}
	if u.Fragment != "" || strings.Contains(raw, "#") {
		return "", fmt.Errorf("origin must not include fragment")
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("origin must use http or https")
	}
	host := u.Hostname()
	if host == "" {
		return "", fmt.Errorf("origin host is required")
	}
	port := u.Port()
	if port == "" {
		if strings.HasSuffix(u.Host, ":") {
			return "", fmt.Errorf("invalid origin port")
		}
		port = "80"
		if scheme == "https" {
			port = "443"
		}
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return "", fmt.Errorf("invalid origin port")
	}
	return fmt.Sprintf("%s://%s", scheme, net.JoinHostPort(strings.ToLower(host), strconv.Itoa(portNumber))), nil
}

func advertisedListenHost(listenAddress string) string {
	host := strings.TrimSpace(listenAddress)
	if host == "" {
		return "localhost"
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsUnspecified() {
		return "localhost"
	}
	return host
}
