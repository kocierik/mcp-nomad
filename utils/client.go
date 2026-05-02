// Package utils provides utility functions and types for interacting with the Nomad HTTP API.
// The Nomad REST client intentionally avoids importing github.com/hashicorp/nomad (very large
// dependency tree): behavior is aligned with Nomad CLI environment variables such as NOMAD_REGION,
// TLS env vars consumed in buildTLSConfig, and REST paths under /v1/.
package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NomadClient handles interactions with the Nomad API.
// It provides methods for managing jobs, deployments, namespaces, nodes, allocations,
// variables, volumes, and ACL tokens.
type NomadClient struct {
	address          string
	token            string
	httpClient       *http.Client
	DefaultTailLines int // Default number of lines to show when tailing logs
}

// NewNomadClient creates a new Nomad client with the specified address and token.
// It validates the connection to the Nomad server before returning.
//
// Example:
//
//	client, err := NewNomadClient("http://localhost:4646", "your-token")
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewNomadClient(address, token string) (*NomadClient, error) {
	// Validate the address
	if address == "" {
		return nil, fmt.Errorf("nomad address is required")
	}

	// Create the client
	client := &NomadClient{
		address: address,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: buildTLSConfig(),
			},
		},
		DefaultTailLines: 100, // Default to showing last 100 lines
	}

	// Test the connection
	_, err := client.makeRequest(context.Background(), "GET", "status/leader", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Nomad server: %w", err)
	}

	return client, nil
}

// SetToken sets the ACL token for the client
func (c *NomadClient) SetToken(token string) {
	c.token = token
}

// GetToken returns the current ACL token
func (c *NomadClient) GetToken() string {
	return c.token
}

// SetDefaultTailLines sets the default number of lines to show when tailing logs
func (c *NomadClient) SetDefaultTailLines(lines int) error {
	if lines <= 0 {
		return fmt.Errorf("number of lines must be positive")
	}
	c.DefaultTailLines = lines
	return nil
}

// GetDefaultTailLines returns the current default number of lines for log tailing
func (c *NomadClient) GetDefaultTailLines() int {
	return c.DefaultTailLines
}

// buildTLSConfig constructs a *tls.Config from the standard NOMAD_* TLS environment
// variables, matching the behavior of the official Nomad CLI and Go SDK.
func buildTLSConfig() *tls.Config {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if caFile := os.Getenv("NOMAD_CACERT"); caFile != "" {
		if caPEM, ok := readCACertPEM(caFile); ok {
			pool := x509.NewCertPool()
			if pool.AppendCertsFromPEM(caPEM) {
				cfg.RootCAs = pool
			}
		}
	}

	if os.Getenv("NOMAD_SKIP_VERIFY") == "true" {
		cfg.InsecureSkipVerify = true
	}

	if name := os.Getenv("NOMAD_TLS_SERVER_NAME"); name != "" {
		cfg.ServerName = name
	}

	return cfg
}

// readCACertPEM reads the PEM at path from NOMAD_CACERT using os.Root so the open is
// confined to the parent directory (same file as Nomad CLI, without path components in "file name").
func readCACertPEM(caPath string) ([]byte, bool) {
	p := filepath.Clean(strings.TrimSpace(caPath))
	if p == "" || p == "." {
		return nil, false
	}
	var err error
	if !filepath.IsAbs(p) {
		p, err = filepath.Abs(p)
		if err != nil {
			return nil, false
		}
	}
	dir := filepath.Dir(p)
	base := filepath.Base(p)
	if base == "" || base == "." || base == ".." {
		return nil, false
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, false
	}
	defer root.Close()
	f, err := root.Open(base)
	if err != nil {
		return nil, false
	}
	defer func() { _ = f.Close() }()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, false
	}
	return b, true
}
