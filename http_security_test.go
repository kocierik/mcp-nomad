package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestParseAllowedOrigins_RejectsInvalidEntries(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		errPart string
	}{
		{name: "empty entry", env: "http://localhost:8080, ,http://127.0.0.1:8080", errPart: "empty"},
		{name: "malformed", env: "not-a-url", errPart: "invalid"},
		{name: "path", env: "http://localhost:8080/path", errPart: "path"},
		{name: "query", env: "http://localhost:8080?x=1", errPart: "query"},
		{name: "userinfo", env: "http://user@localhost:8080", errPart: "userinfo"},
		{name: "empty query", env: "http://localhost:8080?", errPart: "query"},
		{name: "fragment", env: "http://localhost:8080#fragment", errPart: "fragment"},
		{name: "empty fragment", env: "http://localhost:8080#", errPart: "fragment"},
		{name: "invalid port", env: "http://localhost:70000", errPart: "port"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseAllowedOrigins(tt.env)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !containsInsensitive(err.Error(), tt.errPart) {
				t.Fatalf("expected error to contain %q, got %q", tt.errPart, err.Error())
			}
		})
	}
}

func TestBuildDefaultAllowedOrigins(t *testing.T) {
	origins, err := buildAllowedOrigins("", "8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !origins["http://localhost:8080"] {
		t.Fatalf("expected localhost default origin")
	}
	if !origins["http://127.0.0.1:8080"] {
		t.Fatalf("expected loopback ip default origin")
	}
	if !origins["http://[::1]:8080"] {
		t.Fatalf("expected IPv6 loopback default origin")
	}
}

func TestOriginCanonicalDefaultPorts(t *testing.T) {
	for _, testCase := range []struct{ input, expected string }{
		{"https://example.com", "https://example.com:443"},
		{"https://example.com:443", "https://example.com:443"},
		{"http://localhost", "http://localhost:80"},
		{"http://[::1]", "http://[::1]:80"},
	} {
		t.Run(testCase.input, func(t *testing.T) {
			actual, err := normalizeOrigin(testCase.input)
			if err != nil || actual != testCase.expected {
				t.Fatalf("normalizeOrigin(%q) = %q, %v; want %q", testCase.input, actual, err, testCase.expected)
			}
		})
	}
}

func TestHTTPRejectsAmbiguousOriginHeaders(t *testing.T) {
	for _, origins := range [][]string{
		{"http://localhost:8080", "http://attacker.example:8080"},
		{"http://localhost:8080", "http://localhost:8080"},
		{""},
	} {
		request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/mcp", nil)
		request.Header.Set("Authorization", "Bearer mcp-secret")
		for _, origin := range origins {
			request.Header.Add("Origin", origin)
		}
		response := httptest.NewRecorder()
		handler := newHTTPAuthOriginMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}), httpSecurityConfig{
			tokenHash:      sha256Hex("mcp-secret"),
			allowedOrigins: map[string]bool{"http://localhost:8080": true},
		})
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusForbidden {
			t.Fatalf("origins %q: expected 403, got %d", origins, response.Code)
		}
	}
}

func TestHTTPAuthOriginMiddleware(t *testing.T) {
	cfg := httpSecurityConfig{
		tokenHash:      sha256Hex("mcp-secret"),
		allowedOrigins: map[string]bool{"http://localhost:8080": true},
	}

	h := newHTTPAuthOriginMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}), cfg)

	tests := []struct {
		name               string
		authValues         []string
		origin             string
		expectedStatusCode int
	}{
		{name: "missing auth no origin", expectedStatusCode: http.StatusUnauthorized},
		{name: "raw token", authValues: []string{"mcp-secret"}, expectedStatusCode: http.StatusUnauthorized},
		{name: "basic token", authValues: []string{"Basic abc"}, expectedStatusCode: http.StatusUnauthorized},
		{name: "wrong bearer", authValues: []string{"Bearer wrong"}, expectedStatusCode: http.StatusUnauthorized},
		{name: "multiple auth headers", authValues: []string{"Bearer mcp-secret", "Bearer mcp-secret"}, expectedStatusCode: http.StatusUnauthorized},
		{name: "malicious prefix", authValues: []string{"Bearer mcp-secret"}, origin: "http://localhost:8080.evil.example", expectedStatusCode: http.StatusForbidden},
		{name: "allowed origin", authValues: []string{"Bearer mcp-secret"}, origin: "http://localhost:8080", expectedStatusCode: http.StatusNoContent},
		{name: "no origin allowed when auth ok", authValues: []string{"Bearer mcp-secret"}, expectedStatusCode: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.test/mcp", nil)
			for _, v := range tt.authValues {
				req.Header.Add("Authorization", v)
			}
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatusCode {
				t.Fatalf("expected %d, got %d", tt.expectedStatusCode, rr.Code)
			}
			if tt.expectedStatusCode == http.StatusUnauthorized {
				if rr.Header().Get("WWW-Authenticate") != "Bearer" {
					t.Fatalf("expected WWW-Authenticate Bearer header")
				}
			}
		})
	}
}

func TestStreamableHTTP_InitializeAndToolsCall_UsesFixedNomadToken(t *testing.T) {
	var mu sync.Mutex
	capturedNomadToken := ""

	nomad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/status/leader":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`"127.0.0.1:4646"`))
		case "/v1/jobs":
			mu.Lock()
			capturedNomadToken = r.Header.Get("X-Nomad-Token")
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer nomad.Close()

	nomadClient, err := utils.NewNomadClient(nomad.URL, "nomad-static-token")
	if err != nil {
		t.Fatalf("failed to create Nomad client: %v", err)
	}

	mcpServer := server.NewMCPServer("Nomad MCP", "test", server.WithToolCapabilities(true))
	registerTools(mcpServer, nomadClient, log.New(io.Discard, "", 0))

	handler := newHTTPAuthOriginMiddleware(
		server.NewStreamableHTTPServer(mcpServer, server.WithStateLess(true)),
		httpSecurityConfig{
			tokenHash:      sha256Hex("client-bearer"),
			allowedOrigins: map[string]bool{"http://inspector.example:3000": true},
		},
	)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	initResp := postRPC(t, ts.URL+"/mcp", map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  mcp.MethodInitialize,
		"params":  map[string]any{},
	}, "Bearer client-bearer", "http://inspector.example:3000")
	defer initResp.Body.Close()
	if initResp.StatusCode != http.StatusOK {
		t.Fatalf("expected initialize status 200, got %d", initResp.StatusCode)
	}

	callResp := postRPC(t, ts.URL+"/mcp", map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  mcp.MethodToolsCall,
		"params": map[string]any{
			"name":      "list_jobs",
			"arguments": map[string]any{},
		},
	}, "Bearer client-bearer", "http://inspector.example:3000")
	defer callResp.Body.Close()
	if callResp.StatusCode != http.StatusOK {
		t.Fatalf("expected tools/call status 200, got %d", callResp.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if capturedNomadToken != "nomad-static-token" {
		t.Fatalf("expected Nomad token %q, got %q", "nomad-static-token", capturedNomadToken)
	}
	if capturedNomadToken == "client-bearer" {
		t.Fatalf("request bearer token was forwarded to Nomad")
	}
}

func TestSSEIngress_Secured(t *testing.T) {
	mcpServer := server.NewMCPServer("Nomad MCP", "test")
	handler := newHTTPAuthOriginMiddleware(
		server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080")),
		httpSecurityConfig{
			tokenHash:      sha256Hex("mcp-secret"),
			allowedOrigins: map[string]bool{"http://localhost:8080": true},
		},
	)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	noAuthReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/sse", nil)
	noAuthResp, err := http.DefaultClient.Do(noAuthReq)
	if err != nil {
		t.Fatalf("failed no-auth request: %v", err)
	}
	defer noAuthResp.Body.Close()
	if noAuthResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", noAuthResp.StatusCode)
	}

	badOriginReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/sse", nil)
	badOriginReq.Header.Set("Authorization", "Bearer mcp-secret")
	badOriginReq.Header.Set("Origin", "http://localhost:8080.attacker")
	badOriginResp, err := http.DefaultClient.Do(badOriginReq)
	if err != nil {
		t.Fatalf("failed bad-origin request: %v", err)
	}
	defer badOriginResp.Body.Close()
	if badOriginResp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", badOriginResp.StatusCode)
	}

	goodReq, _ := http.NewRequest(http.MethodGet, ts.URL+"/sse", nil)
	goodReq.Header.Set("Authorization", "Bearer mcp-secret")
	goodReq.Header.Set("Origin", "http://localhost:8080")
	goodResp, err := http.DefaultClient.Do(goodReq)
	if err != nil {
		t.Fatalf("failed good request: %v", err)
	}
	defer goodResp.Body.Close()
	if goodResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", goodResp.StatusCode)
	}
}

func postRPC(t *testing.T, endpoint string, payload map[string]any, auth string, origin string) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	return resp
}

func containsInsensitive(s string, substr string) bool {
	return bytes.Contains(bytes.ToLower([]byte(s)), bytes.ToLower([]byte(substr)))
}
