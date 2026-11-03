package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type fakeDynamicResourcesNomad struct {
	lastGetJobNamespace         string
	lastGetJobVersionsNamespace string
}

func (f *fakeDynamicResourcesNomad) GetJob(ctx context.Context, jobID, namespace string) (types.Job, error) {
	f.lastGetJobNamespace = namespace
	return types.Job{ID: jobID, Name: "job-" + jobID}, nil
}

func (f *fakeDynamicResourcesNomad) GetJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error) {
	f.lastGetJobVersionsNamespace = namespace
	return []types.Job{{ID: jobID + "-v1"}}, nil
}

func (f *fakeDynamicResourcesNomad) GetNode(ctx context.Context, nodeID string) (types.Node, error) {
	return types.Node{ID: nodeID}, nil
}

func (f *fakeDynamicResourcesNomad) GetAllocation(ctx context.Context, allocID string) (types.Allocation, error) {
	return types.Allocation{ID: allocID, TaskStates: map[string]types.TaskState{"task": {}}}, nil
}

func (f *fakeDynamicResourcesNomad) GetAllocationLogs(ctx context.Context, allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
	return "allocation logs", nil
}

func readResource(t *testing.T, s *server.MCPServer, uri string) mcp.ReadResourceResult {
	t.Helper()

	request := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":%q}}`, uri)
	response := s.HandleMessage(t.Context(), []byte(request))

	jsonRPCResponse, ok := response.(mcp.JSONRPCResponse)
	if !ok {
		if jsonRPCError, ok := response.(mcp.JSONRPCError); ok {
			t.Fatalf("resources/read returned JSONRPC error: code=%d message=%s", jsonRPCError.Error.Code, jsonRPCError.Error.Message)
		}
		t.Fatalf("unexpected response type: %T", response)
	}

	result, ok := jsonRPCResponse.Result.(mcp.ReadResourceResult)
	if !ok {
		t.Fatalf("unexpected resources/read result type: %T", jsonRPCResponse.Result)
	}

	return result
}

func mustTextContent(t *testing.T, result mcp.ReadResourceResult) mcp.TextResourceContents {
	t.Helper()

	if len(result.Contents) != 1 {
		t.Fatalf("expected exactly one resource content item, got %d", len(result.Contents))
	}

	content, ok := result.Contents[0].(mcp.TextResourceContents)
	if !ok {
		t.Fatalf("unexpected resource content type: %T", result.Contents[0])
	}

	return content
}

func TestRegisterResourcesWithOptions_ReadsStaticDocsAndSystemInfo(t *testing.T) {
	fake := &fakeDynamicResourcesNomad{}
	s := server.NewMCPServer("test", "test", server.WithResourceCapabilities(true, true))
	logger := log.New(io.Discard, "", 0)

	startup := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	RegisterResourcesWithOptions(s, fake, logger, ResourceOptions{
		ReadmeContent:  "embedded readme",
		LicenseContent: "embedded license",
		ServerVersion:  "test-version",
		BuildTime:      "build-time",
		GoVersion:      "go1.26.2",
		StartTime:      startup,
	})

	readme := mustTextContent(t, readResource(t, s, "docs://readme"))
	if readme.Text != "embedded readme" {
		t.Fatalf("unexpected README content: %q", readme.Text)
	}

	license := mustTextContent(t, readResource(t, s, "docs://license"))
	if license.Text != "embedded license" {
		t.Fatalf("unexpected LICENSE content: %q", license.Text)
	}

	system1 := mustTextContent(t, readResource(t, s, "system://info"))
	system2 := mustTextContent(t, readResource(t, s, "system://info"))

	var systemInfo1 map[string]interface{}
	if err := json.Unmarshal([]byte(system1.Text), &systemInfo1); err != nil {
		t.Fatalf("unmarshal system info #1: %v", err)
	}

	var systemInfo2 map[string]interface{}
	if err := json.Unmarshal([]byte(system2.Text), &systemInfo2); err != nil {
		t.Fatalf("unmarshal system info #2: %v", err)
	}

	if got, _ := systemInfo1["server_version"].(string); got != "test-version" {
		t.Fatalf("server_version mismatch: got %q", got)
	}

	if got, _ := systemInfo1["build_time"].(string); got != "build-time" {
		t.Fatalf("build_time mismatch: got %q", got)
	}

	if got, _ := systemInfo1["go_version"].(string); got != "go1.26.2" {
		t.Fatalf("go_version mismatch: got %q", got)
	}

	start1, _ := systemInfo1["start_time"].(string)
	start2, _ := systemInfo2["start_time"].(string)
	if start1 != startup.Format(time.RFC3339) {
		t.Fatalf("start_time mismatch: got %q", start1)
	}
	if start2 != start1 {
		t.Fatalf("start_time changed between reads: %q != %q", start2, start1)
	}
}

func TestRegisterResourcesWithOptions_UsesEffectiveNamespaceForJobResources(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "env-namespace")

	fake := &fakeDynamicResourcesNomad{}
	s := server.NewMCPServer("test", "test", server.WithResourceCapabilities(true, true))
	logger := log.New(io.Discard, "", 0)

	RegisterResourcesWithOptions(s, fake, logger, ResourceOptions{})

	_ = readResource(t, s, "nomad://jobs/job-a/spec?namespace=query-namespace")
	if fake.lastGetJobNamespace != "query-namespace" {
		t.Fatalf("GetJob namespace mismatch: got %q", fake.lastGetJobNamespace)
	}

	_ = readResource(t, s, "nomad://jobs/job-a/history")
	if fake.lastGetJobVersionsNamespace != "env-namespace" {
		t.Fatalf("GetJobVersions namespace mismatch for env fallback: got %q", fake.lastGetJobVersionsNamespace)
	}

	_ = readResource(t, s, "nomad://jobs/job-a/history?namespace=history-ns")
	if fake.lastGetJobVersionsNamespace != "history-ns" {
		t.Fatalf("GetJobVersions namespace mismatch for query override: got %q", fake.lastGetJobVersionsNamespace)
	}
}

func TestExtractIDFromURI_ValidatesNomadShape(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		prefix string
		suffix string
		want   string
	}{
		{name: "valid spec", uri: "nomad://jobs/job-1/spec", prefix: "jobs/", suffix: "/spec", want: "job-1"},
		{name: "valid with query", uri: "nomad://jobs/job-1/spec?namespace=apps", prefix: "jobs/", suffix: "/spec", want: "job-1"},
		{name: "invalid scheme", uri: "http://jobs/job-1/spec", prefix: "jobs/", suffix: "/spec", want: ""},
		{name: "invalid host", uri: "nomad://nodes/job-1/spec", prefix: "jobs/", suffix: "/spec", want: ""},
		{name: "invalid suffix", uri: "nomad://jobs/job-1/history", prefix: "jobs/", suffix: "/spec", want: ""},
		{name: "missing id", uri: "nomad://jobs//spec?namespace=apps", prefix: "jobs/", suffix: "/spec", want: ""},
		{name: "extra path segments", uri: "nomad://jobs/group/job-1/spec", prefix: "jobs/", suffix: "/spec", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIDFromURI(tt.uri, tt.prefix, tt.suffix); got != tt.want {
				t.Fatalf("extractIDFromURI(%q) = %q, want %q", tt.uri, got, tt.want)
			}
		})
	}
}

func TestRegisterResourcesWithOptions_DoesNotReadStaticFilesFromWorkingDirectory(t *testing.T) {
	fake := &fakeDynamicResourcesNomad{}
	s := server.NewMCPServer("test", "test", server.WithResourceCapabilities(true, true))
	logger := log.New(io.Discard, "", 0)

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

	RegisterResourcesWithOptions(s, fake, logger, ResourceOptions{
		ReadmeContent:  "embedded readme",
		LicenseContent: "embedded license",
	})

	readme := mustTextContent(t, readResource(t, s, "docs://readme"))
	if readme.Text != "embedded readme" {
		t.Fatalf("unexpected README content: %q", readme.Text)
	}
}
