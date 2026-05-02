package utils

import "strings"

// CanonicalAuthorizationBearer strips RFC 9665-ish "Bearer" prefix only (case-insensitive).
// Raw tokens stored as the whole Authorization header value remain supported — no prefix required.
func CanonicalAuthorizationBearer(authorizationHeader string) string {
	s := strings.TrimSpace(authorizationHeader)
	if s == "" {
		return ""
	}
	// Common form: Bearer <token>
	firstSpace := strings.IndexByte(s, ' ')
	if firstSpace <= 0 {
		return s
	}
	scheme := s[:firstSpace]
	if !strings.EqualFold(scheme, "Bearer") {
		return s // custom scheme — pass through unchanged (avoid breaking non-Bearer setups)
	}
	return strings.TrimSpace(s[firstSpace+1:])
}
