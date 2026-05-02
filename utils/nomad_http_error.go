package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// MaxNomadHTTPErrorBodyBytes caps retained body bytes from Nomad error responses.
const MaxNomadHTTPErrorBodyBytes = 512

// NomadHTTPError describes a Nomad HTTP API failure (typically status >= 400).
// The stored body slice is capped; use errors.As for structured handling instead of dumping raw payloads to operators.
type NomadHTTPError struct {
	StatusCode int
	Method     string
	Path       string // API path fragment under /v1/ (normalized).
	body       []byte // truncated UTF-8 or raw capped bytes before validation
	truncated  bool   // original response exceeded MaxNomadHTTPErrorBodyBytes
}

// NewNomadHTTPError records a Nomad HTTP error. Copy at most MaxNomadHTTPErrorBodyBytes from respBody.
func NewNomadHTTPError(statusCode int, method, normalizedPath string, respBody []byte) *NomadHTTPError {
	trunc := len(respBody) > MaxNomadHTTPErrorBodyBytes
	store := respBody
	if trunc {
		store = append([]byte(nil), respBody[:MaxNomadHTTPErrorBodyBytes]...)
	} else if len(respBody) > 0 {
		store = append([]byte(nil), respBody...)
	}
	return &NomadHTTPError{
		StatusCode: statusCode,
		Method:     method,
		Path:       normalizedPath,
		body:       store,
		truncated:  trunc,
	}
}

func (e *NomadHTTPError) Error() string {
	if e == nil {
		return ""
	}
	snip := sanitizeErrorBodySnippet(e.body, e.truncated)
	if snip == "" {
		return fmt.Sprintf("nomad API error %s %s: HTTP %d", e.Method, e.Path, e.StatusCode)
	}
	return fmt.Sprintf("nomad API error %s %s: HTTP %d (%s)", e.Method, e.Path, e.StatusCode, snip)
}

// Status returns the HTTP status code from Nomad (e.g. 404).
func (e *NomadHTTPError) Status() int {
	if e == nil {
		return 0
	}
	return e.StatusCode
}

// Snippet returns a one-line sanitized excerpt of the capped response body.
func (e *NomadHTTPError) Snippet() string {
	if e == nil {
		return ""
	}
	return sanitizeErrorBodySnippet(e.body, e.truncated)
}

func sanitizeErrorBodySnippet(b []byte, truncated bool) string {
	if len(b) == 0 {
		if truncated {
			return "(response truncated)"
		}
		return ""
	}
	if !utf8.Valid(b) {
		return "[invalid UTF-8 response]"
	}
	s := strings.ReplaceAll(string(b), "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.Join(strings.Fields(s), " ")
	if truncated {
		return s + "…"
	}
	return s
}
