package test

import (
	"os"
	"testing"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	NomadAddr       string
	NomadToken      string
	SkipIntegration bool
}

// GetTestConfig returns test configuration from environment variables
func GetTestConfig() *TestConfig {
	return &TestConfig{
		NomadAddr:       getEnvOrDefault("NOMAD_ADDR", "http://localhost:4646"),
		NomadToken:      getEnvOrDefault("NOMAD_TOKEN", ""),
		SkipIntegration: getEnvOrDefault("SKIP_INTEGRATION", "false") == "true",
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SkipIfIntegrationSkipped skips the test if integration tests are disabled
func SkipIfIntegrationSkipped(t *testing.T) {
	config := GetTestConfig()
	if config.SkipIntegration {
		t.Skip("Integration tests are disabled (set SKIP_INTEGRATION=false to enable)")
	}
}

// SkipIfNoNomadServer skips the test if no Nomad server is available
func SkipIfNoNomadServer(t *testing.T) {
	config := GetTestConfig()
	if config.NomadAddr == "" {
		t.Skip("No Nomad server configured (set NOMAD_ADDR to enable)")
	}
}

// SkipIfNoNomadToken skips the test if no Nomad token is provided
func SkipIfNoNomadToken(t *testing.T) {
	config := GetTestConfig()
	if config.NomadToken == "" {
		t.Skip("No Nomad token provided (set NOMAD_TOKEN to enable)")
	}
}
