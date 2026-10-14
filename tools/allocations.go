package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllocationTools registers all allocation-related tools
func RegisterAllocationTools(s *server.MCPServer, nomadClient utils.AllocationAPI, logger *log.Logger) {
	// List allocations tool
	listAllocationsTool := mcp.NewTool("list_allocations",
		mcp.WithDescription("List all allocations in Nomad"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list allocations from (default: default)"),
		),
		mcp.WithString("job_id",
			mcp.Description("If set, list allocations via GET /v1/job/:job_id/allocations (namespace respected); otherwise GET /v1/allocations"),
		),
		mcp.WithString("node_id",
			mcp.Description("Only allocations placed on this node (full ID or unique prefix). Filtering is done client-side"),
		),
		mcp.WithString("client_status",
			mcp.Description("Only allocations with this ClientStatus"),
			mcp.Enum("pending", "running", "complete", "failed", "lost", "unknown", ""),
		),
		mcp.WithBoolean("compact",
			mcp.Description("Return a slim projection (ID, job, group, node, statuses, times, per-task state + last event message) instead of full allocation objects. Strongly recommended with namespace=\"*\""),
		),
	)
	s.AddTool(listAllocationsTool, ListAllocationsHandler(nomadClient, logger))

	// Get allocation tool
	getAllocationTool := mcp.NewTool("get_allocation",
		mcp.WithDescription("Get allocation details by ID"),
		mcp.WithString("allocation_id",
			mcp.Required(),
			mcp.Description("The ID of the allocation to retrieve"),
		),
	)
	s.AddTool(getAllocationTool, GetAllocationHandler(nomadClient, logger))

	// Stop allocation tool
	stopAllocationTool := mcp.NewTool("stop_allocation",
		mcp.WithDescription("Stop a running allocation"),
		mcp.WithString("allocation_id",
			mcp.Required(),
			mcp.Description("The ID of the allocation to stop"),
		),
	)
	s.AddTool(stopAllocationTool, StopAllocationHandler(nomadClient, logger))
}

// ListAllocationsHandler returns a handler for listing allocations
func ListAllocationsHandler(client utils.AllocationAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		namespace := utils.EffectiveToolNamespace(arguments)
		jobID := ""
		if j, ok := arguments["job_id"].(string); ok {
			jobID = strings.TrimSpace(j)
		}

		allocations, err := client.ListAllocations(ctx, namespace, jobID)
		if err != nil {
			logger.Printf("Error listing allocations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list allocations", err), nil
		}

		nodeID, _ := arguments["node_id"].(string)
		nodeID = strings.TrimSpace(nodeID)
		clientStatus, _ := arguments["client_status"].(string)
		// Clients with a stale cached schema send booleans as strings ("true").
		compact := false
		switch v := arguments["compact"].(type) {
		case bool:
			compact = v
		case string:
			compact = strings.EqualFold(strings.TrimSpace(v), "true")
		}

		filtered := allocations[:0:0]
		for _, a := range allocations {
			if nodeID != "" && !strings.HasPrefix(a.NodeID, nodeID) {
				continue
			}
			if clientStatus != "" && a.ClientStatus != clientStatus {
				continue
			}
			filtered = append(filtered, a)
		}

		var payload interface{} = filtered
		if compact {
			payload = compactAllocations(filtered)
		}
		allocationsJSON, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format allocations", err), nil
		}

		return mcp.NewToolResultText(string(allocationsJSON)), nil
	}
}

// GetAllocationHandler returns a handler for getting allocation details
func GetAllocationHandler(client utils.AllocationAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		allocID, ok := arguments["allocation_id"].(string)
		if !ok || allocID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		allocation, err := client.GetAllocation(ctx, allocID)
		if err != nil {
			logger.Printf("Error getting allocation: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get allocation", err), nil
		}

		allocationJSON, err := json.MarshalIndent(allocation, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format allocation", err), nil
		}

		return mcp.NewToolResultText(string(allocationJSON)), nil
	}
}

// StopAllocationHandler returns a handler for stopping an allocation
func StopAllocationHandler(client utils.AllocationAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		allocationID, ok := arguments["allocation_id"].(string)
		if !ok || allocationID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		err := client.StopAllocation(ctx, allocationID)
		if err != nil {
			logger.Printf("Error stopping allocation: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to stop allocation", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Allocation %s stopped successfully", allocationID)), nil
	}
}

type compactTask struct {
	State     string `json:"state"`
	Failed    bool   `json:"failed"`
	LastEvent string `json:"last_event,omitempty"`
}

type compactAlloc struct {
	ID           string                 `json:"id"`
	Namespace    string                 `json:"namespace,omitempty"`
	JobID        string                 `json:"job_id"`
	TaskGroup    string                 `json:"task_group"`
	NodeID       string                 `json:"node_id"`
	NodeName     string                 `json:"node_name,omitempty"`
	ClientStatus string                 `json:"client_status"`
	DesiredState string                 `json:"desired_status"`
	Description  string                 `json:"client_description,omitempty"`
	EvalID       string                 `json:"eval_id"`
	Created      string                 `json:"created"`
	Modified     string                 `json:"modified"`
	Tasks        map[string]compactTask `json:"tasks,omitempty"`
}

// compactAllocations projects allocations onto the fields needed to triage them
// (full objects with event histories are ~7KB each, ~4.5MB for namespace "*").
func compactAllocations(in []types.Allocation) []compactAlloc {
	out := make([]compactAlloc, 0, len(in))
	for _, a := range in {
		c := compactAlloc{
			ID: a.ID, Namespace: a.Namespace, JobID: a.JobID, TaskGroup: a.TaskGroup,
			NodeID: a.NodeID, NodeName: a.NodeName, ClientStatus: a.ClientStatus,
			DesiredState: a.DesiredStatus, Description: a.ClientDescription, EvalID: a.EvalID,
			Created:  time.Unix(0, a.CreateTime).UTC().Format(time.RFC3339),
			Modified: time.Unix(0, a.ModifyTime).UTC().Format(time.RFC3339),
		}
		if len(a.TaskStates) > 0 {
			c.Tasks = make(map[string]compactTask, len(a.TaskStates))
			for name, ts := range a.TaskStates {
				ct := compactTask{State: ts.State, Failed: ts.Failed}
				if n := len(ts.Events); n > 0 {
					e := ts.Events[n-1]
					msg := e.DisplayMessage
					if msg == "" {
						msg = e.Message
					}
					ct.LastEvent = e.Type + ": " + msg
				}
				c.Tasks[name] = ct
			}
		}
		out = append(out, c)
	}
	return out
}
