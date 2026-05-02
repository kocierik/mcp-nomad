package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

		allocationsJSON, err := json.MarshalIndent(allocations, "", "  ")
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
