package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type structuredResultKind int

const (
	structuredNone structuredResultKind = iota
	structuredObject
	structuredArrayObjects
	structuredArrayStrings
	structuredScalar
)

type toolPolicyEntry struct {
	readOnly       bool
	structuredKind structuredResultKind
}

var toolPolicyCatalog = map[string]toolPolicyEntry{
	"list_acl_tokens":        {readOnly: true, structuredKind: structuredArrayObjects},
	"get_acl_token":          {readOnly: true, structuredKind: structuredObject},
	"create_acl_token":       {readOnly: false, structuredKind: structuredObject},
	"delete_acl_token":       {readOnly: false, structuredKind: structuredNone},
	"list_acl_policies":      {readOnly: true, structuredKind: structuredArrayObjects},
	"get_acl_policy":         {readOnly: true, structuredKind: structuredObject},
	"create_acl_policy":      {readOnly: false, structuredKind: structuredObject},
	"delete_acl_policy":      {readOnly: false, structuredKind: structuredNone},
	"list_acl_roles":         {readOnly: true, structuredKind: structuredArrayObjects},
	"get_acl_role":           {readOnly: true, structuredKind: structuredObject},
	"create_acl_role":        {readOnly: false, structuredKind: structuredObject},
	"delete_acl_role":        {readOnly: false, structuredKind: structuredNone},
	"bootstrap_acl_token":    {readOnly: false, structuredKind: structuredObject},
	"list_allocations":       {readOnly: true, structuredKind: structuredArrayObjects},
	"get_allocation":         {readOnly: true, structuredKind: structuredObject},
	"stop_allocation":        {readOnly: false, structuredKind: structuredNone},
	"get_cluster_leader":     {readOnly: true, structuredKind: structuredArrayObjects},
	"list_cluster_peers":     {readOnly: true, structuredKind: structuredArrayStrings},
	"list_regions":           {readOnly: true, structuredKind: structuredArrayStrings},
	"list_deployments":       {readOnly: true, structuredKind: structuredArrayObjects},
	"get_deployment":         {readOnly: true, structuredKind: structuredObject},
	"list_jobs":              {readOnly: true, structuredKind: structuredArrayObjects},
	"get_job":                {readOnly: true, structuredKind: structuredObject},
	"run_job":                {readOnly: false, structuredKind: structuredObject},
	"stop_job":               {readOnly: false, structuredKind: structuredObject},
	"scale_job":              {readOnly: false, structuredKind: structuredObject},
	"get_job_allocations":    {readOnly: true, structuredKind: structuredArrayObjects},
	"get_job_evaluations":    {readOnly: true, structuredKind: structuredArrayObjects},
	"get_evaluation":         {readOnly: true, structuredKind: structuredObject},
	"get_job_deployments":    {readOnly: true, structuredKind: structuredArrayObjects},
	"get_job_summary":        {readOnly: true, structuredKind: structuredObject},
	"get_job_services":       {readOnly: true, structuredKind: structuredArrayObjects},
	"get_allocation_logs":    {readOnly: true, structuredKind: structuredNone},
	"list_namespaces":        {readOnly: true, structuredKind: structuredArrayObjects},
	"create_namespace":       {readOnly: false, structuredKind: structuredObject},
	"delete_namespace":       {readOnly: false, structuredKind: structuredObject},
	"list_nodes":             {readOnly: true, structuredKind: structuredArrayObjects},
	"get_node":               {readOnly: true, structuredKind: structuredObject},
	"drain_node":             {readOnly: false, structuredKind: structuredObject},
	"eligibility_node":       {readOnly: false, structuredKind: structuredObject},
	"list_sentinel_policies": {readOnly: true, structuredKind: structuredArrayObjects},
	"get_sentinel_policy":    {readOnly: true, structuredKind: structuredObject},
	"create_sentinel_policy": {readOnly: false, structuredKind: structuredObject},
	"delete_sentinel_policy": {readOnly: false, structuredKind: structuredObject},
	"list_variables":         {readOnly: true, structuredKind: structuredArrayObjects},
	"get_variable":           {readOnly: true, structuredKind: structuredObject},
	"create_variable":        {readOnly: false, structuredKind: structuredObject},
	"delete_variable":        {readOnly: false, structuredKind: structuredObject},
	"list_volumes":           {readOnly: true, structuredKind: structuredArrayObjects},
	"get_volume":             {readOnly: true, structuredKind: structuredObject},
	"delete_volume":          {readOnly: false, structuredKind: structuredNone},
}

func applyToolPolicy(s *server.MCPServer, readOnly bool) {
	listed := s.ListTools()
	if len(listed) == 0 {
		return
	}

	adapted := make([]server.ServerTool, 0, len(listed))
	for name, entry := range listed {
		policy, ok := toolPolicyCatalog[name]
		if !ok {
			continue
		}
		if readOnly && !policy.readOnly {
			continue
		}

		tool := entry.Tool
		mcp.WithReadOnlyHintAnnotation(policy.readOnly)(&tool)
		mcp.WithDestructiveHintAnnotation(!policy.readOnly)(&tool)
		mcp.WithIdempotentHintAnnotation(policy.readOnly)(&tool)

		if schema, ok := policy.outputSchema(); ok {
			tool.RawOutputSchema = nil
			tool.OutputSchema = schema
		}

		handler := entry.Handler
		if policy.structuredKind != structuredNone {
			handler = adaptStructuredToolResult(policy.structuredKind, handler)
		}

		adapted = append(adapted, server.ServerTool{Tool: tool, Handler: handler})
	}

	sort.Slice(adapted, func(i, j int) bool {
		return adapted[i].Tool.Name < adapted[j].Tool.Name
	})
	s.SetTools(adapted...)
}

func adaptStructuredToolResult(kind structuredResultKind, next server.ToolHandlerFunc) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := next(ctx, request)
		if err != nil || result == nil || result.IsError {
			return result, err
		}

		text, ok := firstTextResult(result.Content)
		if !ok {
			return result, nil
		}

		structured, parseErr := parseStructuredResultText(kind, text)
		if parseErr != nil {
			return mcp.NewToolResultError(fmt.Sprintf("tool returned invalid JSON for structured adaptation: %v", parseErr)), nil
		}

		return mcp.NewToolResultStructured(structured, text), nil
	}
}

func firstTextResult(content []mcp.Content) (string, bool) {
	for _, item := range content {
		text, ok := item.(mcp.TextContent)
		if ok {
			return text.Text, true
		}
	}
	return "", false
}

func parseStructuredResultText(kind structuredResultKind, text string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()

	var parsed any
	if err := decoder.Decode(&parsed); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("unexpected trailing JSON tokens")
		}
		return nil, err
	}

	switch kind {
	case structuredObject:
		object, ok := parsed.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected JSON object")
		}
		return object, nil
	case structuredArrayObjects:
		if parsed == nil {
			return map[string]any{"items": []any{}}, nil
		}
		items, ok := parsed.([]any)
		if !ok {
			return nil, fmt.Errorf("expected JSON array")
		}
		for _, item := range items {
			if _, ok := item.(map[string]any); !ok {
				return nil, fmt.Errorf("expected array of JSON objects")
			}
		}
		return map[string]any{"items": items}, nil
	case structuredArrayStrings:
		if parsed == nil {
			return map[string]any{"items": []string{}}, nil
		}
		items, ok := parsed.([]any)
		if !ok {
			return nil, fmt.Errorf("expected JSON array")
		}
		for _, item := range items {
			if _, ok := item.(string); !ok {
				return nil, fmt.Errorf("expected array of JSON strings")
			}
		}
		return map[string]any{"items": items}, nil
	case structuredScalar:
		if _, ok := parsed.(map[string]any); ok {
			return nil, fmt.Errorf("expected scalar JSON value")
		}
		if _, ok := parsed.([]any); ok {
			return nil, fmt.Errorf("expected scalar JSON value")
		}
		return map[string]any{"value": parsed}, nil
	default:
		return nil, fmt.Errorf("unsupported structured result kind")
	}
}

func (p toolPolicyEntry) outputSchema() (mcp.ToolOutputSchema, bool) {
	switch p.structuredKind {
	case structuredObject:
		return mcp.ToolOutputSchema{
			Type:                 "object",
			Properties:           map[string]any{},
			Required:             []string{},
			AdditionalProperties: true,
		}, true
	case structuredArrayObjects:
		return mcp.ToolOutputSchema{
			Type: "object",
			Properties: map[string]any{
				"items": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "object"},
				},
			},
			Required:             []string{"items"},
			AdditionalProperties: false,
		}, true
	case structuredArrayStrings:
		return mcp.ToolOutputSchema{
			Type: "object",
			Properties: map[string]any{
				"items": map[string]any{
					"type":  "array",
					"items": map[string]any{"type": "string"},
				},
			},
			Required:             []string{"items"},
			AdditionalProperties: false,
		}, true
	case structuredScalar:
		return mcp.ToolOutputSchema{
			Type: "object",
			Properties: map[string]any{
				"value": map[string]any{},
			},
			Required:             []string{"value"},
			AdditionalProperties: false,
		}, true
	default:
		return mcp.ToolOutputSchema{}, false
	}
}
