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

// RegisterVariableTools registers all variable-related tools
func RegisterVariableTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List variables tool
	listVariablesTool := mcp.NewTool("list_variables",
		mcp.WithDescription("List all variables in Nomad"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list variables from (default: default)"),
		),
		mcp.WithString("prefix",
			mcp.Description("Optional prefix to filter variables"),
		),
		mcp.WithString("next_token",
			mcp.Description("Token for pagination"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Number of variables per page"),
		),
		mcp.WithString("filter",
			mcp.Description("Expression to filter results"),
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
		mcp.WithString("namespace",
			mcp.Description("The namespace of the variable (default: default)"),
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
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key for the variable"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value for the variable"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the variable (default: default)"),
		),
		mcp.WithNumber("cas",
			mcp.Description("Check-and-set value for optimistic concurrency control"),
		),
		mcp.WithString("lock_operation",
			mcp.Description("Lock operation to perform (acquire, release)"),
			mcp.Enum("acquire", "release"),
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
		mcp.WithString("namespace",
			mcp.Description("The namespace of the variable (default: default)"),
		),
		mcp.WithNumber("cas",
			mcp.Description("Check-and-set value for optimistic concurrency control"),
		),
	)
	s.AddTool(deleteVariableTool, DeleteVariableHandler(nomadClient, logger))
}

// ListVariablesHandler returns a handler for listing variables
func ListVariablesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		namespace := "default"
		if ns, ok := arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		prefix := ""
		if p, ok := arguments["prefix"].(string); ok {
			prefix = p
		}

		nextToken := ""
		if nt, ok := arguments["next_token"].(string); ok {
			nextToken = nt
		}

		perPage := 0
		if pp, ok := arguments["per_page"].(float64); ok {
			perPage = int(pp)
		}

		filter := ""
		if f, ok := arguments["filter"].(string); ok {
			filter = f
		}

		variables, err := client.ListVariables(namespace, prefix, nextToken, perPage, filter)
		if err != nil {
			logger.Printf("Error listing variables: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list variables", err), nil
		}

		variablesJSON, err := json.MarshalIndent(variables, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format variables", err), nil
		}

		return mcp.NewToolResultText(string(variablesJSON)), nil
	}
}

// GetVariableHandler returns a handler for getting variable details
func GetVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		path, ok := arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		namespace := "default"
		if ns, ok := arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		variable, err := client.GetVariable(path, namespace)
		if err != nil {
			logger.Printf("Error getting variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get variable", err), nil
		}

		variableJSON, err := json.MarshalIndent(variable, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format variable", err), nil
		}

		return mcp.NewToolResultText(string(variableJSON)), nil
	}
}

// CreateVariableHandler returns a handler for creating a variable
func CreateVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		path, ok := arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		key, ok := arguments["key"].(string)
		if !ok || key == "" {
			return mcp.NewToolResultError("key is required"), nil
		}

		value, ok := arguments["value"].(string)
		if !ok || value == "" {
			return mcp.NewToolResultError("value is required"), nil
		}

		namespace := "default"
		if ns, ok := arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		cas := 0
		if c, ok := arguments["cas"].(float64); ok && c > 0 {
			cas = int(c)
		}

		lockOp := ""
		if l, ok := arguments["lock_operation"].(string); ok && l != "" {
			lockOp = l
		}

		// Create a map with the key-value pair
		items := map[string]string{
			key: value,
		}

		// Create the variable value with the required structure
		variableValue := map[string]interface{}{
			"Items": items,
		}

		// Add CAS if provided
		if cas > 0 {
			variableValue["CAS"] = cas
		}

		// Add lock operation if provided
		if lockOp != "" {
			variableValue["LockOperation"] = lockOp
		}

		// Convert to JSON string
		jsonValue, err := json.Marshal(variableValue)
		if err != nil {
			logger.Printf("Error marshaling variable value: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to format variable value", err), nil
		}

		variable := types.Variable{
			Path:  path,
			Value: string(jsonValue),
		}

		err = client.CreateVariable(variable, namespace, cas, lockOp)
		if err != nil {
			logger.Printf("Error creating variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create variable", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Variable created at path: %s with key: %s", path, key),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// DeleteVariableHandler returns a handler for deleting a variable
func DeleteVariableHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		path, ok := arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		namespace := "default"
		if ns, ok := arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		cas := 0
		if c, ok := arguments["cas"].(float64); ok && c > 0 {
			cas = int(c)
		}

		err := client.DeleteVariable(path, namespace, cas)
		if err != nil {
			logger.Printf("Error deleting variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete variable", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Variable deleted at path: %s", path),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
