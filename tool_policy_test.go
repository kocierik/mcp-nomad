package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/kocierik/mcp-nomad/test/mocks"
	"github.com/kocierik/mcp-nomad/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestToolPolicy_ReadOnlyFiltersWritesAndPreventsWriteCalls(t *testing.T) {
	calledRunJob := false
	mock := &mocks.MockNomadClient{}
	mock.ListJobsFunc = func(context.Context, string, string) ([]types.JobSummary, error) {
		return []types.JobSummary{{ID: "job-1", Name: "job-1"}}, nil
	}
	mock.RunJobFunc = func(context.Context, string, bool) (map[string]interface{}, error) {
		calledRunJob = true
		return map[string]interface{}{"ok": true}, nil
	}

	s := server.NewMCPServer("test", "test", server.WithToolCapabilities(true))
	registerToolsWithPolicy(s, mock, log.New(io.Discard, "", 0), true)

	listed := callToolsList(t, s)
	if !hasTool(listed.Tools, "list_jobs") {
		t.Fatalf("expected list_jobs to be discoverable in read-only mode")
	}
	if !hasTool(listed.Tools, "list_regions") {
		t.Fatalf("expected list_regions to be discoverable in read-only mode")
	}
	if !hasTool(listed.Tools, "get_allocation_logs") {
		t.Fatalf("expected get_allocation_logs to be discoverable in read-only mode")
	}
	if hasTool(listed.Tools, "run_job") {
		t.Fatalf("run_job must be absent in read-only mode")
	}
	if hasTool(listed.Tools, "create_namespace") {
		t.Fatalf("create_namespace must be absent in read-only mode")
	}

	resp := s.HandleMessage(t.Context(), mustJSON(t, map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      "run_job",
			"arguments": map[string]any{"job_spec": "{}"},
		},
	}))

	if _, ok := resp.(mcp.JSONRPCError); !ok {
		t.Fatalf("expected JSONRPCError for unregistered write tool in read-only mode, got %T", resp)
	}
	if calledRunJob {
		t.Fatalf("run_job handler should not be called when read-only policy removes the tool")
	}
}

func TestToolPolicy_CatalogCoversRegisteredTools(t *testing.T) {
	server := server.NewMCPServer("test", "test")
	registerTools(server, &mocks.MockNomadClient{}, log.New(io.Discard, "", 0))
	registered := server.ListTools()
	for name := range registered {
		if _, exists := toolPolicyCatalog[name]; !exists {
			t.Errorf("registered tool %q has no policy", name)
		}
	}
	for name := range toolPolicyCatalog {
		if _, exists := registered[name]; !exists {
			t.Errorf("policy %q has no registered tool", name)
		}
	}
}

func TestToolPolicy_FiltersUnknownToolsAndSetsAnnotations(t *testing.T) {
	s := server.NewMCPServer("test", "test", server.WithToolCapabilities(true))
	registerToolsWithPolicy(s, &mocks.MockNomadClient{}, log.New(io.Discard, "", 0), false)
	s.AddTool(mcp.NewTool("unknown_tool"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("{}"), nil
	})

	applyToolPolicy(s, false)
	listed := callToolsList(t, s)

	if hasTool(listed.Tools, "unknown_tool") {
		t.Fatalf("unknown_tool should be excluded by explicit policy catalog")
	}
	response := s.HandleMessage(t.Context(), mustJSON(t, map[string]any{
		"jsonrpc": "2.0", "id": 8, "method": "tools/call",
		"params": map[string]any{"name": "unknown_tool", "arguments": map[string]any{}},
	}))
	if _, ok := response.(mcp.JSONRPCError); !ok {
		t.Fatalf("unknown tool must be uncallable, got %T", response)
	}

	listJobsTool := mustFindTool(t, listed.Tools, "list_jobs")
	if listJobsTool.Annotations.ReadOnlyHint == nil || !*listJobsTool.Annotations.ReadOnlyHint {
		t.Fatalf("list_jobs should be annotated readOnly=true")
	}
	if listJobsTool.Annotations.DestructiveHint == nil || *listJobsTool.Annotations.DestructiveHint {
		t.Fatalf("list_jobs should be annotated destructive=false")
	}
	if listJobsTool.Annotations.IdempotentHint == nil || !*listJobsTool.Annotations.IdempotentHint {
		t.Fatalf("list_jobs should be annotated idempotent=true")
	}

	runJobTool := mustFindTool(t, listed.Tools, "run_job")
	if runJobTool.Annotations.ReadOnlyHint == nil || *runJobTool.Annotations.ReadOnlyHint {
		t.Fatalf("run_job should be annotated readOnly=false")
	}
	if runJobTool.Annotations.DestructiveHint == nil || !*runJobTool.Annotations.DestructiveHint {
		t.Fatalf("run_job should be annotated destructive=true")
	}
	if runJobTool.Annotations.IdempotentHint == nil || *runJobTool.Annotations.IdempotentHint {
		t.Fatalf("run_job should be annotated idempotent=false")
	}
}

func TestToolPolicy_StructuredArrayObjectAndSchemaValidation(t *testing.T) {
	mock := &mocks.MockNomadClient{}
	mock.ListJobsFunc = func(context.Context, string, string) ([]types.JobSummary, error) {
		return []types.JobSummary{{ID: "job-1", Name: "api"}}, nil
	}
	mock.GetJobFunc = func(context.Context, string, string) (types.Job, error) {
		return types.Job{ID: "job-1", Name: "api"}, nil
	}

	s := server.NewMCPServer(
		"test",
		"test",
		server.WithToolCapabilities(true),
		server.WithOutputSchemaValidation(),
	)
	registerToolsWithPolicy(s, mock, log.New(io.Discard, "", 0), false)

	listed := callToolsList(t, s)
	listJobsTool := mustFindTool(t, listed.Tools, "list_jobs")
	if listJobsTool.OutputSchema.Type != "object" {
		t.Fatalf("list_jobs output schema should be object envelope, got %q", listJobsTool.OutputSchema.Type)
	}
	itemsSchema, ok := listJobsTool.OutputSchema.Properties["items"].(map[string]any)
	if !ok {
		t.Fatalf("list_jobs output schema should expose items array")
	}
	if got, _ := itemsSchema["type"].(string); got != "array" {
		t.Fatalf("list_jobs output schema items type should be array, got %v", itemsSchema["type"])
	}

	listJobsResult := callToolResult(t, s, "list_jobs", map[string]any{})
	if listJobsResult.IsError {
		t.Fatalf("list_jobs should succeed")
	}
	listJobsStructured, ok := listJobsResult.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("list_jobs structuredContent should be object envelope, got %T", listJobsResult.StructuredContent)
	}
	items, ok := listJobsStructured["items"].([]any)
	if !ok {
		t.Fatalf("list_jobs structuredContent.items should be []any, got %T", listJobsStructured["items"])
	}
	if len(items) != 1 {
		t.Fatalf("expected one item in list_jobs structured content, got %d", len(items))
	}
	text := firstText(t, listJobsResult)
	if !strings.Contains(text, "\"ID\": \"job-1\"") {
		t.Fatalf("expected text fallback to preserve original handler JSON, got %q", text)
	}

	getJobResult := callToolResult(t, s, "get_job", map[string]any{"job_id": "job-1"})
	if getJobResult.IsError {
		t.Fatalf("get_job should succeed")
	}
	getJobStructured, ok := getJobResult.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("get_job structuredContent should be object, got %T", getJobResult.StructuredContent)
	}
	if got, _ := getJobStructured["ID"].(string); got != "job-1" {
		t.Fatalf("get_job structured ID mismatch: got %q", got)
	}
}

func TestToolPolicy_InvalidJSONResultBecomesToolError(t *testing.T) {
	s := server.NewMCPServer(
		"test",
		"test",
		server.WithToolCapabilities(true),
		server.WithOutputSchemaValidation(),
	)
	s.AddTool(mcp.NewTool("list_jobs"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("not-json"), nil
	})

	applyToolPolicy(s, false)
	result := callToolResult(t, s, "list_jobs", map[string]any{})
	if !result.IsError {
		t.Fatalf("classified JSON tool with invalid JSON text must return isError=true")
	}
}

func TestToolPolicy_RejectsTrailingJSON(t *testing.T) {
	for _, text := range []string{`{} {}`, `[] []`, `{} trailing`} {
		if _, err := parseStructuredResultText(structuredObject, text); err == nil {
			t.Errorf("accepted trailing JSON in %q", text)
		}
	}
}

func TestToolPolicy_LargeNumbersPreservedUsingJSONNumber(t *testing.T) {
	s := server.NewMCPServer(
		"test",
		"test",
		server.WithToolCapabilities(true),
		server.WithOutputSchemaValidation(),
	)
	s.AddTool(mcp.NewTool("get_job"), func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("{\"big\":18446744073709551615}"), nil
	})

	applyToolPolicy(s, false)
	result := callToolResult(t, s, "get_job", map[string]any{"job_id": "job-1"})
	if result.IsError {
		t.Fatalf("expected successful structured conversion")
	}
	structured, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("structured content should be object, got %T", result.StructuredContent)
	}
	big, ok := structured["big"].(json.Number)
	if !ok {
		t.Fatalf("expected big to be json.Number, got %T (%v)", structured["big"], structured["big"])
	}
	if big.String() != "18446744073709551615" {
		t.Fatalf("big number changed: got %s", big.String())
	}
}

func TestToolPolicy_AllocationLogsRemainTextOnly(t *testing.T) {
	mock := &mocks.MockNomadClient{}
	mock.GetAllocationLogsFunc = func(context.Context, string, string, string, bool, int64, int64) (string, error) {
		return "hello logs", nil
	}

	s := server.NewMCPServer("test", "test", server.WithToolCapabilities(true), server.WithOutputSchemaValidation())
	registerToolsWithPolicy(s, mock, log.New(io.Discard, "", 0), false)

	result := callToolResult(t, s, "get_allocation_logs", map[string]any{
		"allocation_id": "alloc-1",
		"task":          "web",
	})
	if result.IsError {
		t.Fatalf("get_allocation_logs should succeed")
	}
	if result.StructuredContent != nil {
		t.Fatalf("get_allocation_logs should remain text-only")
	}
	if !strings.Contains(firstText(t, result), "hello logs") {
		t.Fatalf("expected text fallback to be unchanged")
	}
}

func callToolsList(t *testing.T, s *server.MCPServer) mcp.ListToolsResult {
	t.Helper()
	resp := s.HandleMessage(t.Context(), mustJSON(t, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]any{},
	}))

	jsonResp := mustJSONRPCResponse(t, resp)
	result, ok := jsonResp.Result.(mcp.ListToolsResult)
	if !ok {
		t.Fatalf("unexpected tools/list result type: %T", jsonResp.Result)
	}
	return result
}

func callToolResult(t *testing.T, s *server.MCPServer, name string, arguments map[string]any) mcp.CallToolResult {
	t.Helper()
	resp := s.HandleMessage(t.Context(), mustJSON(t, map[string]any{
		"jsonrpc": "2.0",
		"id":      7,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      name,
			"arguments": arguments,
		},
	}))

	jsonResp := mustJSONRPCResponse(t, resp)
	result, ok := jsonResp.Result.(*mcp.CallToolResult)
	if !ok {
		t.Fatalf("unexpected tools/call result type: %T", jsonResp.Result)
	}
	if result == nil {
		t.Fatalf("tools/call result pointer is nil")
	}
	return *result
}

func firstText(t *testing.T, result mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatalf("tool result content is empty")
	}
	text, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("first content is %T, expected mcp.TextContent", result.Content[0])
	}
	return text.Text
}

func hasTool(tools []mcp.Tool, name string) bool {
	for _, tool := range tools {
		if tool.Name == name {
			return true
		}
	}
	return false
}

func mustFindTool(t *testing.T, tools []mcp.Tool, name string) mcp.Tool {
	t.Helper()
	for _, tool := range tools {
		if tool.Name == name {
			return tool
		}
	}
	t.Fatalf("tool %q not found", name)
	return mcp.Tool{}
}

func mustJSONRPCResponse(t *testing.T, msg mcp.JSONRPCMessage) mcp.JSONRPCResponse {
	t.Helper()
	resp, ok := msg.(mcp.JSONRPCResponse)
	if ok {
		return resp
	}
	if errResp, ok := msg.(mcp.JSONRPCError); ok {
		t.Fatalf("unexpected JSONRPC error: code=%d message=%s", errResp.Error.Code, errResp.Error.Message)
	}
	t.Fatalf("unexpected response type: %T", msg)
	return mcp.JSONRPCResponse{}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return b
}
