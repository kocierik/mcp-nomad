// File: tools/nodes.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/kocierik/nomad-mcp-server/utils"
)

// ListNodesHandler returns a handler for the list_nodes tool
func ListNodesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		status := ""
		if s, ok := request.Params.Arguments["status"].(string); ok {
			status = s
		}

		nodes, err := client.ListNodes(status)
		if err != nil {
			logger.Printf("Error listing nodes: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list nodes", err), nil
		}

		nodesJSON, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format node list", err), nil
		}

		return mcp.NewToolResultText(string(nodesJSON)), nil
	}
}

// GetNodeHandler returns a handler for the get_node tool
func GetNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, ok := request.Params.Arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("Node ID is required"), nil
		}

		node, err := client.GetNode(nodeID)
		if err != nil {
			logger.Printf("Error getting node %s: %v", nodeID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get node %s", nodeID), err), nil
		}

		nodeJSON, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format node details", err), nil
		}

		return mcp.NewToolResultText(string(nodeJSON)), nil
	}
}

// DrainNodeHandler returns a handler for the drain_node tool
func DrainNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, ok := request.Params.Arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("Node ID is required"), nil
		}

		enable, ok := request.Params.Arguments["enable"].(bool)
		if !ok {
			return mcp.NewToolResultError("Enable parameter is required"), nil
		}

		deadline := -1.0
		if d, ok := request.Params.Arguments["deadline"].(float64); ok {
			deadline = d
		}

		result, err := client.DrainNode(nodeID, enable, int64(deadline))
		if err != nil {
			logger.Printf("Error setting drain mode for node %s: %v", nodeID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to set drain mode for node %s", nodeID), err), nil
		}

		status := "enabled"
		if !enable {
			status = "disabled"
		}

		return mcp.NewToolResultText(fmt.Sprintf("Drain mode %s for node %s. %s", status, nodeID, result)), nil
	}
}

// EligibilityNodeHandler returns a handler for the eligibility_node tool
func EligibilityNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodeID, ok := request.Params.Arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("Node ID is required"), nil
		}

		eligible, ok := request.Params.Arguments["eligible"].(string)
		if !ok {
			return mcp.NewToolResultError("Eligible parameter must be 'eligible' or 'ineligible'"), nil
		}

		_, err := client.EligibilityNode(nodeID, eligible)
		if err != nil {
			logger.Printf("Error setting eligibility for node %s: %v", nodeID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to set eligibility for node %s", nodeID), err), nil
		}

		status := "enabled"
		if eligible == "ineligible" {
			status = "disabled"
		}

		return mcp.NewToolResultText(fmt.Sprintf("Eligibility %s for node %s.", status, nodeID)), nil
	}
}
