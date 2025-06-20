// File: tools/namespaces.go
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

// RegisterNamespaceTools registers all namespace-related tools
func RegisterNamespaceTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List namespaces tool
	listNamespacesTool := mcp.NewTool("list_namespaces",
		mcp.WithDescription("List all namespaces in Nomad"),
	)
	s.AddTool(listNamespacesTool, ListNamespacesHandler(nomadClient, logger))

	// Create namespace tool
	createNamespaceTool := mcp.NewTool("create_namespace",
		mcp.WithDescription("Create a new namespace"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the namespace to create"),
		),
		mcp.WithString("description",
			mcp.Description("Description of the namespace"),
		),
	)
	s.AddTool(createNamespaceTool, CreateNamespaceHandler(nomadClient, logger))

	// Delete namespace tool
	deleteNamespaceTool := mcp.NewTool("delete_namespace",
		mcp.WithDescription("Delete a namespace"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the namespace to delete"),
		),
	)
	s.AddTool(deleteNamespaceTool, DeleteNamespaceHandler(nomadClient, logger))
}

// ListNamespacesHandler returns a handler for listing namespaces
func ListNamespacesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespaces, err := client.ListNamespaces()
		if err != nil {
			logger.Printf("Error listing namespaces: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list namespaces", err), nil
		}

		namespacesJSON, err := json.MarshalIndent(namespaces, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format namespaces", err), nil
		}

		return mcp.NewToolResultText(string(namespacesJSON)), nil
	}
}

// CreateNamespaceHandler returns a handler for creating a namespace
func CreateNamespaceHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		namespace := types.Namespace{
			Name:        name,
			Description: description,
		}

		err := client.CreateNamespace(namespace)
		if err != nil {
			logger.Printf("Error creating namespace: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create namespace", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Successfully created namespace %s", name),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// DeleteNamespaceHandler returns a handler for deleting a namespace
func DeleteNamespaceHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		name, ok := arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}

		err := client.DeleteNamespace(name)
		if err != nil {
			logger.Printf("Error deleting namespace: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete namespace", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Successfully deleted namespace %s", name),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
