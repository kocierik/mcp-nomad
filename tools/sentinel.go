package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterSentinelTools registers all Sentinel-related tools with the MCP server
func RegisterSentinelTools(s *server.MCPServer, client *utils.NomadClient, logger *log.Logger) {
	// List policies tool
	listPoliciesTool := mcp.NewTool("list_sentinel_policies",
		mcp.WithDescription("List all Sentinel policies"),
	)
	s.AddTool(listPoliciesTool, ListSentinelPoliciesHandler(client, logger))

	// Get policy tool
	getPolicyTool := mcp.NewTool("get_sentinel_policy",
		mcp.WithDescription("Get a specific Sentinel policy by name"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the policy to retrieve"),
		),
	)
	s.AddTool(getPolicyTool, GetSentinelPolicyHandler(client, logger))

	// Create policy tool
	createPolicyTool := mcp.NewTool("create_sentinel_policy",
		mcp.WithDescription("Create a new Sentinel policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the policy"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the policy"),
		),
		mcp.WithString("scope",
			mcp.Required(),
			mcp.Description("The scope of the policy (e.g., submit-job)"),
		),
		mcp.WithString("enforcement_level",
			mcp.Required(),
			mcp.Description("The enforcement level (advisory, soft-mandatory, hard-mandatory)"),
			mcp.Enum("advisory", "soft-mandatory", "hard-mandatory"),
		),
		mcp.WithString("policy",
			mcp.Required(),
			mcp.Description("The Sentinel policy code"),
		),
	)
	s.AddTool(createPolicyTool, CreateSentinelPolicyHandler(client, logger))

	// Delete policy tool
	deletePolicyTool := mcp.NewTool("delete_sentinel_policy",
		mcp.WithDescription("Delete a Sentinel policy"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the policy to delete"),
		),
	)
	s.AddTool(deletePolicyTool, DeleteSentinelPolicyHandler(client, logger))
}

// ListSentinelPoliciesHandler returns a handler for listing Sentinel policies
func ListSentinelPoliciesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		policies, err := client.ListSentinelPolicies()
		if err != nil {
			logger.Printf("Error listing Sentinel policies: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list Sentinel policies", err), nil
		}

		policiesJSON, err := json.MarshalIndent(policies, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format Sentinel policies", err), nil
		}

		return mcp.NewToolResultText(string(policiesJSON)), nil
	}
}

// GetSentinelPolicyHandler returns a handler for getting Sentinel policy details
func GetSentinelPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		policy, err := client.GetSentinelPolicy(name)
		if err != nil {
			logger.Printf("Error getting Sentinel policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get Sentinel policy", err), nil
		}

		policyJSON, err := json.MarshalIndent(policy, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format Sentinel policy", err), nil
		}

		return mcp.NewToolResultText(string(policyJSON)), nil
	}
}

// CreateSentinelPolicyHandler returns a handler for creating a Sentinel policy
func CreateSentinelPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		description := ""
		if d, ok := arguments["description"].(string); ok {
			description = d
		}

		policy := types.SentinelPolicy{
			Name:        name,
			Description: description,
		}

		err := client.CreateSentinelPolicy(policy)
		if err != nil {
			logger.Printf("Error creating Sentinel policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create Sentinel policy", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Successfully created Sentinel policy %s", name),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// DeleteSentinelPolicyHandler returns a handler for deleting a Sentinel policy
func DeleteSentinelPolicyHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		err := client.DeleteSentinelPolicy(name)
		if err != nil {
			logger.Printf("Error deleting Sentinel policy: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete Sentinel policy", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Successfully deleted Sentinel policy %s", name),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
