package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Opt-in only; never derived from NOMAD_ADDR so a developer's cluster is not targeted by accident.
func liveNomadServer(t *testing.T, readOnly bool) *server.MCPServer {
	t.Helper()
	address := os.Getenv("NOMAD_INTEGRATION_ADDR")
	if address == "" {
		t.Skip("set NOMAD_INTEGRATION_ADDR to a disposable Nomad dev agent to run live tests")
	}
	client, err := utils.NewNomadClient(address, os.Getenv("NOMAD_INTEGRATION_TOKEN"))
	if err != nil {
		t.Fatalf("connect to Nomad: %v", err)
	}
	s := server.NewMCPServer("live", "test", server.WithToolCapabilities(true), server.WithOutputSchemaValidation())
	registerToolsWithPolicy(s, client, log.New(io.Discard, "", 0), readOnly)
	return s
}

func liveObject(t *testing.T, result *mcp.CallToolResult) map[string]any {
	t.Helper()
	object, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("expected structured object, got %T", result.StructuredContent)
	}
	return object
}

func mustSucceed(t *testing.T, s *server.MCPServer, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()
	result := callToolResult(t, s, name, arguments)
	if result.IsError {
		t.Fatalf("%s failed: %s", name, firstText(t, result))
	}
	return &result
}

func TestLiveNomad_VariablesCASLifecycle(t *testing.T) {
	s := liveNomadServer(t, false)
	path := fmt.Sprintf("mcp-it/%d", time.Now().UnixNano())
	t.Cleanup(func() { callToolResult(t, s, "delete_variable", map[string]any{"path": path}) })

	mustSucceed(t, s, "create_variable", map[string]any{"path": path, "key": "k", "value": "v1", "cas": 0})
	if result := callToolResult(t, s, "create_variable", map[string]any{"path": path, "key": "k", "value": "v2", "cas": 0}); !result.IsError {
		t.Fatalf("create with cas=0 must fail when the variable exists")
	}

	variable := liveObject(t, mustSucceed(t, s, "get_variable", map[string]any{"path": path}))
	items, _ := variable["Items"].(map[string]any)
	if items["k"] != "v1" {
		t.Fatalf("expected Items.k=v1, got %v", variable["Items"])
	}
	modifyIndex, ok := variable["ModifyIndex"].(json.Number)
	if !ok {
		t.Fatalf("expected numeric ModifyIndex, got %T", variable["ModifyIndex"])
	}
	stale, _ := modifyIndex.Int64()

	listed := liveObject(t, mustSucceed(t, s, "list_variables", map[string]any{"prefix": path}))
	if entries, _ := listed["items"].([]any); len(entries) != 1 {
		t.Fatalf("expected one listed variable, got %v", listed["items"])
	}

	if result := callToolResult(t, s, "delete_variable", map[string]any{"path": path, "cas": stale + 1000}); !result.IsError {
		t.Fatalf("delete with a stale cas must fail")
	}
	mustSucceed(t, s, "delete_variable", map[string]any{"path": path, "cas": json.Number(modifyIndex.String())})
	if result := callToolResult(t, s, "get_variable", map[string]any{"path": path}); !result.IsError {
		t.Fatalf("variable should be gone after conditional delete")
	}
}

func TestLiveNomad_JobHistoryAndPurge(t *testing.T) {
	s := liveNomadServer(t, false)
	jobID := fmt.Sprintf("mcp-it-%d", time.Now().UnixNano())
	t.Cleanup(func() { callToolResult(t, s, "stop_job", map[string]any{"job_id": jobID, "purge": true}) })

	spec := fmt.Sprintf(`job %q {
  type        = "batch"
  datacenters = ["*"]
  group "g" {
    task "t" {
      driver = "raw_exec"
      config {
        command = "true"
      }
    }
  }
}`, jobID)
	mustSucceed(t, s, "run_job", map[string]any{"job_spec": spec, "detach": true})

	job := liveObject(t, mustSucceed(t, s, "get_job", map[string]any{"job_id": jobID}))
	if job["ID"] != jobID {
		t.Fatalf("expected job %s, got %v", jobID, job["ID"])
	}

	response := mustJSONRPCResponse(t, s.HandleMessage(t.Context(), mustJSON(t, map[string]any{
		"jsonrpc": "2.0", "id": 3, "method": "resources/read",
		"params": map[string]any{"uri": "nomad://jobs/" + jobID + "/history"},
	})))
	history, ok := response.Result.(mcp.ReadResourceResult)
	if !ok || len(history.Contents) != 1 {
		t.Fatalf("unexpected history result %T", response.Result)
	}
	text, _ := history.Contents[0].(mcp.TextResourceContents)
	if !strings.Contains(text.Text, jobID) {
		t.Fatalf("history must decode the Versions envelope, got %q", text.Text)
	}

	mustSucceed(t, s, "stop_job", map[string]any{"job_id": jobID, "purge": true})
}

func TestLiveNomad_ReadOnlyDiscovery(t *testing.T) {
	s := liveNomadServer(t, true)
	listed := callToolsList(t, s)
	if hasTool(listed.Tools, "create_variable") || hasTool(listed.Tools, "run_job") {
		t.Fatalf("read-only mode must not expose write tools")
	}
	namespaces := liveObject(t, mustSucceed(t, s, "list_namespaces", map[string]any{}))
	if entries, _ := namespaces["items"].([]any); len(entries) == 0 {
		t.Fatalf("expected at least the default namespace")
	}
}
