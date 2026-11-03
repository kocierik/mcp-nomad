package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmbeddedRootAssetsAvailableOutsideWorkingDirectory(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	if got := readmeAsset(); !strings.Contains(got, "mcp-nomad") {
		t.Fatalf("embedded README content missing expected marker")
	}

	if got := licenseAsset(); !strings.Contains(strings.ToLower(got), "mit license") {
		t.Fatalf("embedded LICENSE content missing expected marker")
	}

	if _, err := os.Stat(filepath.Join(tempDir, "README.md")); !os.IsNotExist(err) {
		t.Fatalf("test sanity check failed: README.md unexpectedly exists in temp dir")
	}
}

func TestVersionInjectionViaLdflags(t *testing.T) {
	expected := os.Getenv("EXPECT_INJECTED_VERSION")
	if expected == "" {
		t.Skip("set EXPECT_INJECTED_VERSION when running with -ldflags injection check")
	}

	if version != expected {
		t.Fatalf("version mismatch: got %q, expected %q", version, expected)
	}
}
