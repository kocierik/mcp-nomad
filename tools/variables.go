package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type createVariableArguments struct {
	Path          string  `json:"path"`
	Key           string  `json:"key"`
	Value         string  `json:"value"`
	Namespace     string  `json:"namespace,omitempty"`
	CAS           *uint64 `json:"cas,omitempty"`
	LockOperation *string `json:"lock_operation,omitempty"`
}

type deleteVariableArguments struct {
	Path      string  `json:"path"`
	Namespace string  `json:"namespace,omitempty"`
	CAS       *uint64 `json:"cas,omitempty"`
}

type variableDeleteCASAPI interface {
	DeleteVariableCAS(ctx context.Context, path, namespace string, cas *uint64) error
}

// RegisterVariableTools registers all variable-related tools
func RegisterVariableTools(s *server.MCPServer, nomadClient utils.VariableAPI, logger *log.Logger) {
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
func ListVariablesHandler(client utils.VariableAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		namespace := utils.EffectiveToolNamespace(arguments)

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

		variables, err := client.ListVariables(ctx, namespace, prefix, nextToken, perPage, filter)
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
func GetVariableHandler(client utils.VariableAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		path, ok := arguments["path"].(string)
		if !ok || path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		namespace := utils.EffectiveToolNamespace(arguments)

		variable, err := client.GetVariable(ctx, path, namespace)
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
func CreateVariableHandler(client utils.VariableAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, err := parseCreateVariableArguments(request)
		if err != nil {
			if strings.Contains(err.Error(), "cas") {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		if arguments.Path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		if arguments.Key == "" {
			return mcp.NewToolResultError("key is required"), nil
		}

		if arguments.Value == "" {
			return mcp.NewToolResultError("value is required"), nil
		}

		namespace := utils.EffectiveToolNamespace(map[string]interface{}{"namespace": arguments.Namespace})

		cas, explicitCAS, err := parseVariableCASArgument(arguments.CAS)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if arguments.LockOperation != nil && strings.TrimSpace(*arguments.LockOperation) != "" {
			return mcp.NewToolResultError("lock_operation is unsupported: lease payload is not available"), nil
		}

		items := map[string]string{
			arguments.Key: arguments.Value,
		}

		variable := types.Variable{
			Path:  arguments.Path,
			Items: items,
			CAS:   explicitCAS,
		}

		err = client.CreateVariable(ctx, variable, namespace, cas, "")
		if err != nil {
			logger.Printf("Error creating variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to create variable", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Variable created at path: %s with key: %s", arguments.Path, arguments.Key),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

func parseCreateVariableArguments(request mcp.CallToolRequest) (createVariableArguments, error) {
	var arguments createVariableArguments
	if err := decodeToolArguments(request, &arguments); err != nil {
		if rawArgumentsContainKey(request, "cas") {
			return createVariableArguments{}, fmt.Errorf("cas must be a non-negative integer")
		}
		return createVariableArguments{}, err
	}
	return arguments, nil
}

func parseDeleteVariableArguments(request mcp.CallToolRequest) (deleteVariableArguments, error) {
	var arguments deleteVariableArguments
	if err := decodeToolArguments(request, &arguments); err != nil {
		if rawArgumentsContainKey(request, "cas") {
			return deleteVariableArguments{}, fmt.Errorf("cas must be a non-negative integer")
		}
		return deleteVariableArguments{}, err
	}
	return arguments, nil
}

func rawArgumentsContainKey(request mcp.CallToolRequest, key string) bool {
	if argumentsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
		_, found := argumentsMap[key]
		if found {
			return true
		}
	}

	rawArguments := request.Params.RawArguments
	if len(rawArguments) == 0 {
		if request.Params.Arguments == nil {
			return false
		}
		serializedArguments, err := json.Marshal(request.Params.Arguments)
		if err != nil {
			return false
		}
		rawArguments = serializedArguments
	}

	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(rawArguments, &rawMap); err != nil {
		return false
	}
	_, found := rawMap[key]
	return found
}

func decodeToolArguments(request mcp.CallToolRequest, target interface{}) error {
	rawArguments := request.Params.RawArguments
	if len(rawArguments) == 0 {
		if request.Params.Arguments == nil {
			rawArguments = []byte("{}")
		} else {
			serializedArguments, err := json.Marshal(request.Params.Arguments)
			if err != nil {
				return err
			}
			rawArguments = serializedArguments
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(rawArguments))
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("unexpected trailing arguments")
		}
		return err
	}

	return nil
}

func parseVariableCASArgument(casValue *uint64) (cas int, explicit *uint64, err error) {
	if casValue == nil {
		return 0, nil, nil
	}
	if *casValue > uint64(math.MaxInt) {
		return 0, nil, fmt.Errorf("cas exceeds supported range")
	}

	casCopy := *casValue
	return int(casCopy), &casCopy, nil
}

// DeleteVariableHandler returns a handler for deleting a variable
func DeleteVariableHandler(client utils.VariableAPI, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, err := parseDeleteVariableArguments(request)
		if err != nil {
			if strings.Contains(err.Error(), "cas") {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		if arguments.Path == "" {
			return mcp.NewToolResultError("path is required"), nil
		}

		namespace := utils.EffectiveToolNamespace(map[string]interface{}{"namespace": arguments.Namespace})

		cas, explicitCAS, err := parseVariableCASArgument(arguments.CAS)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if arguments.CAS != nil {
			if casAPI, ok := client.(variableDeleteCASAPI); ok {
				err = casAPI.DeleteVariableCAS(ctx, arguments.Path, namespace, explicitCAS)
			} else if cas == 0 {
				return mcp.NewToolResultError("client does not support explicit cas zero for deletion"), nil
			} else {
				err = client.DeleteVariable(ctx, arguments.Path, namespace, cas)
			}
		} else {
			err = client.DeleteVariable(ctx, arguments.Path, namespace, 0)
		}
		if err != nil {
			logger.Printf("Error deleting variable: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to delete variable", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Variable deleted at path: %s", arguments.Path),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
