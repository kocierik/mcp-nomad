package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNomadHTTPError_Error(t *testing.T) {
	err := NewNomadHTTPError(404, "GET", "job/example", []byte(`{"msg":"missing"}`))
	require.ErrorContains(t, err, "nomad API error GET job/example: HTTP 404")
	require.ErrorContains(t, err, "missing")

	var nh *NomadHTTPError
	require.True(t, errors.As(err, &nh))
	require.Equal(t, 404, nh.Status())
	require.Equal(t, "GET", nh.Method)
}

func TestNewNomadHTTPError_truncated(t *testing.T) {
	long := strings.Repeat("a", 600)
	err := NewNomadHTTPError(500, "GET", "x", []byte(long))
	require.True(t, errors.As(err, new(*NomadHTTPError)))
	require.ErrorContains(t, err, "…")
}

func TestNomadHTTPError_binaryBody(t *testing.T) {
	err := NewNomadHTTPError(500, "POST", "jobs", []byte{0xff, 0xfe, 0xfd})
	require.Contains(t, err.Error(), "[invalid UTF-8 response]")
}

func TestCanonicalAuthorizationBearer(t *testing.T) {
	require.Equal(t, "abc.xyz", CanonicalAuthorizationBearer("Bearer abc.xyz"))
	require.Equal(t, "abc.xyz", CanonicalAuthorizationBearer("bearer abc.xyz"))
	require.Equal(t, "abc.xyz", CanonicalAuthorizationBearer("  Bearer abc.xyz  "))
	require.Equal(t, "raw-token", CanonicalAuthorizationBearer("raw-token"))
	require.Equal(t, "", CanonicalAuthorizationBearer("   "))
	require.Equal(t, "Basic xxx", CanonicalAuthorizationBearer("Basic xxx"))
}
