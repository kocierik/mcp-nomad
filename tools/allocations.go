package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListAllocationsHandler returns a handler for listing allocations
func ListAllocationsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		jobID := ""
		if jid, ok := request.Params.Arguments["job_id"].(string); ok {
			jobID = jid
		}

		// List allocations using the Nomad API
		path := "allocations"
		queryParams := make(map[string]string)
		if namespace != "default" {
			queryParams["namespace"] = namespace
		}
		if jobID != "" {
			queryParams["job_id"] = jobID
		}

		body, err := client.MakeRequest("GET", path, queryParams, nil)
		if err != nil {
			logger.Printf("Error listing allocations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list allocations", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// GetAllocationHandler returns a handler for getting allocation details
func GetAllocationHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		allocationID, ok := request.Params.Arguments["allocation_id"].(string)
		if !ok || allocationID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		// Get allocation using the Nomad API
		path := fmt.Sprintf("allocation/%s", allocationID)
		body, err := client.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting allocation: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get allocation", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// StopAllocationHandler returns a handler for stopping an allocation
func StopAllocationHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		allocationID, ok := request.Params.Arguments["allocation_id"].(string)
		if !ok || allocationID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		// Stop allocation using the Nomad API
		path := fmt.Sprintf("allocation/%s/stop", allocationID)
		_, err := client.MakeRequest("POST", path, nil, nil)
		if err != nil {
			logger.Printf("Error stopping allocation: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to stop allocation", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Allocation %s stopped successfully", allocationID)), nil
	}
}
