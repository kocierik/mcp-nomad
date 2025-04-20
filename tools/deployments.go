// File: tools/deployments.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListDeploymentsHandler returns a handler for the list_deployments tool
func ListDeploymentsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		deployments, err := client.ListDeployments(namespace)
		if err != nil {
			logger.Printf("Error listing deployments: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list deployments", err), nil
		}

		deploymentsJSON, err := json.MarshalIndent(deployments, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format deployment list", err), nil
		}

		return mcp.NewToolResultText(string(deploymentsJSON)), nil
	}
}

// GetDeploymentHandler returns a handler for the get_deployment tool
func GetDeploymentHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		deploymentID, ok := request.Params.Arguments["deployment_id"].(string)
		if !ok || deploymentID == "" {
			return mcp.NewToolResultError("Deployment ID is required"), nil
		}

		deployment, err := client.GetDeployment(deploymentID)
		if err != nil {
			logger.Printf("Error getting deployment %s: %v", deploymentID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get deployment %s", deploymentID), err), nil
		}

		deploymentJSON, err := json.MarshalIndent(deployment, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format deployment details", err), nil
		}

		return mcp.NewToolResultText(string(deploymentJSON)), nil
	}
}
