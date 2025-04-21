package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterVariableTools registers all variable-related tools
func RegisterVariableTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List variables tool
	listVariablesTool := mcp.NewTool("list_variables",
		mcp.WithDescription("List all variables in Nomad"),
		mcp.WithString("prefix",
			mcp.Description("Optional prefix to filter variables"),
		),
	)
	s.AddTool(listVariablesTool, ListVariablesHandler(nomadClient, logger))

	// Get variable tool
	getVariableTool := mcp.NewTool("get_variable",
		mcp.WithDescription("Get variable details by path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path of the variable to retrieve"),
		),
	)
	s.AddTool(getVariableTool, GetVariableHandler(nomadClient, logger))

	// Create variable tool
	createVariableTool := mcp.NewTool("create_variable",
		mcp.WithDescription("Create or update a variable"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path where to create the variable"),
		),
		mcp.WithObject("items",
			mcp.Required(),
			mcp.Description("The key-value pairs to store in the variable"),
		),
	)
	s.AddTool(createVariableTool, CreateVariableHandler(nomadClient, logger))

	// Delete variable tool
	deleteVariableTool := mcp.NewTool("delete_variable",
		mcp.WithDescription("Delete a variable"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("The path of the variable to delete"),
		),
	)
	s.AddTool(deleteVariableTool, DeleteVariableHandler(nomadClient, logger))
}

// ListVariablesHandler returns a handler for listing variables
func ListVariablesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path := "vars"
		prefix, ok := request.Params.Arguments["prefix"].(string)
		if ok && prefix != "" {
			path = fmt.Sprintf("vars/%s", prefix)
		}

		body, err := client.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error listing variables: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list variables", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// GetVariableHandler returns a handler for getting variable details
func GetVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, ok := request.Params.Arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("variable path is required"), nil
		}

		body, err := client.MakeRequest("GET", fmt.Sprintf("var/%s", path), nil, nil)
		if err != nil {
			logger.Printf("Error getting variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get variable", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// CreateVariableHandler returns a handler for creating a variable
func CreateVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, ok := request.Params.Arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("variable path is required"), nil
		}

		items, ok := request.Params.Arguments["items"].(map[string]interface{})
		if !ok || len(items) == 0 {
			return mcp.NewToolResultError("variable items are required"), nil
		}

		variable := map[string]interface{}{
			"Path":  path,
			"Items": items,
		}

		body, err := client.MakeRequest("PUT", fmt.Sprintf("var/%s", path), nil, variable)
		if err != nil {
			logger.Printf("Error creating variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create variable", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// DeleteVariableHandler returns a handler for deleting a variable
func DeleteVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		path, ok := request.Params.Arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("variable path is required"), nil
		}

		_, err := client.MakeRequest("DELETE", fmt.Sprintf("var/%s", path), nil, nil)
		if err != nil {
			logger.Printf("Error deleting variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete variable", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Variable %s deleted successfully", path)), nil
	}
}
