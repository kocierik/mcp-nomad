// File: tools/nodes.go
package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterNodeTools registers all node-related tools
func RegisterNodeTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List nodes tool
	listNodesTool := mcp.NewTool("list_nodes",
		mcp.WithDescription("List all nodes in the Nomad cluster"),
		mcp.WithString("status",
			mcp.Description("Filter nodes by status"),
			mcp.Enum("ready", "down", ""),
		),
	)
	s.AddTool(listNodesTool, ListNodesHandler(nomadClient, logger))

	// Get node tool
	getNodeTool := mcp.NewTool("get_node",
		mcp.WithDescription("Get details for a specific node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to retrieve"),
		),
	)
	s.AddTool(getNodeTool, GetNodeHandler(nomadClient, logger))

	// Drain node tool
	drainNodeTool := mcp.NewTool("drain_node",
		mcp.WithDescription("Enable or disable drain mode for a node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to drain"),
		),
		mcp.WithBoolean("enable",
			mcp.Required(),
			mcp.Description("Enable or disable drain mode"),
		),
		mcp.WithNumber("deadline",
			mcp.Description("Deadline in seconds for the drain operation (default: -1, no deadline)"),
		),
	)
	s.AddTool(drainNodeTool, DrainNodeHandler(nomadClient, logger))

	// Eligibility node tool
	eligibilityNodeTool := mcp.NewTool("eligibility_node",
		mcp.WithDescription("Set eligibility for a node"),
		mcp.WithString("node_id",
			mcp.Required(),
			mcp.Description("The ID of the node to set eligibility for"),
		),
		mcp.WithString("eligible",
			mcp.Required(),
			mcp.Description("The eligibility status to set (eligible or ineligible)"),
		),
	)
	s.AddTool(eligibilityNodeTool, EligibilityNodeHandler(nomadClient, logger))
}

// ListNodesHandler returns a handler for listing nodes
func ListNodesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		status := ""
		if s, ok := arguments["status"].(string); ok && s != "" {
			status = s
		}

		nodes, err := client.ListNodes(status)
		if err != nil {
			logger.Printf("Error listing nodes: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list nodes", err), nil
		}

		nodesJSON, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format nodes", err), nil
		}

		return mcp.NewToolResultText(string(nodesJSON)), nil
	}
}

// GetNodeHandler returns a handler for getting node details
func GetNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		nodeID, ok := arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("node_id is required"), nil
		}

		node, err := client.GetNode(nodeID)
		if err != nil {
			logger.Printf("Error getting node: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get node", err), nil
		}

		nodeJSON, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format node", err), nil
		}

		return mcp.NewToolResultText(string(nodeJSON)), nil
	}
}

// DrainNodeHandler returns a handler for draining a node
func DrainNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		nodeID, ok := arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("node_id is required"), nil
		}

		enable := true
		if e, ok := arguments["enable"].(bool); ok {
			enable = e
		}

		deadline := int64(0)
		if d, ok := arguments["deadline"].(float64); ok {
			deadline = int64(d)
		}

		result, err := client.DrainNode(nodeID, enable, deadline)
		if err != nil {
			logger.Printf("Error draining node: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to drain node", err), nil
		}

		response := map[string]string{
			"message": result,
		}

		responseJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format response", err), nil
		}

		return mcp.NewToolResultText(string(responseJSON)), nil
	}
}

// EligibilityNodeHandler returns a handler for setting node eligibility
func EligibilityNodeHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		nodeID, ok := arguments["node_id"].(string)
		if !ok || nodeID == "" {
			return mcp.NewToolResultError("node_id is required"), nil
		}

		eligible, ok := arguments["eligible"].(string)
		if !ok || eligible == "" {
			return mcp.NewToolResultError("eligible is required"), nil
		}

		node, err := client.EligibilityNode(nodeID, eligible)
		if err != nil {
			logger.Printf("Error setting node eligibility: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to set node eligibility", err), nil
		}

		nodeJSON, err := json.MarshalIndent(node, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format node", err), nil
		}

		return mcp.NewToolResultText(string(nodeJSON)), nil
	}
}
