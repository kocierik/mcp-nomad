// File: tools/namespaces.go
package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kocierik/nomad-mcp-server/types"
	"github.com/kocierik/nomad-mcp-server/utils"
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

// ListNamespacesHandler returns a handler for the list_namespaces tool
func ListNamespacesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespaces, err := client.ListNamespaces()
		if err != nil {
			logger.Printf("Error listing namespaces: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list namespaces", err), nil
		}

		namespacesJSON, err := json.MarshalIndent(namespaces, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format namespace list", err), nil
		}

		return mcp.NewToolResultText(string(namespacesJSON)), nil
	}
}

// CreateNamespaceHandler returns a handler for the create_namespace tool
func CreateNamespaceHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("Namespace name is required"), nil
		}

		description := ""
		if desc, ok := request.Params.Arguments["description"].(string); ok {
			description = desc
		}

		namespace := types.Namespace{
			Name:        name,
			Description: description,
		}

		err := client.CreateNamespace(namespace)
		if err != nil {
			logger.Printf("Error creating namespace %s: %v", name, err)
			return mcp.NewToolResultErrorFromErr("Failed to create namespace", err), nil
		}

		return mcp.NewToolResultText("Namespace '" + name + "' created successfully"), nil
	}
}

// DeleteNamespaceHandler returns a handler for the delete_namespace tool
func DeleteNamespaceHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("Namespace name is required"), nil
		}

		err := client.DeleteNamespace(name)
		if err != nil {
			logger.Printf("Error deleting namespace %s: %v", name, err)
			return mcp.NewToolResultErrorFromErr("Failed to delete namespace", err), nil
		}

		return mcp.NewToolResultText("Namespace '" + name + "' deleted successfully"), nil
	}
}
