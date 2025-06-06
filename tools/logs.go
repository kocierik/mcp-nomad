package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterLogTools registers all log-related tools
func RegisterLogTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// Get allocation logs tool
	getAllocationLogsTool := mcp.NewTool("get_allocation_logs",
		mcp.WithDescription("Get logs from a specific task in an allocation"),
		mcp.WithString("allocation_id",
			mcp.Required(),
			mcp.Description("The ID of the allocation"),
		),
		mcp.WithString("task",
			mcp.Required(),
			mcp.Description("The name of the task"),
		),
		mcp.WithString("type",
			mcp.Description("The type of logs to retrieve (stdout or stderr, default: stdout)"),
			mcp.Enum("stdout", "stderr"),
		),
		mcp.WithBoolean("follow",
			mcp.Description("Whether to follow/tail the logs (default: false)"),
		),
		mcp.WithNumber("tail",
			mcp.Description("Number of lines to show from the end (default: 100, 0 means use default)"),
		),
		mcp.WithNumber("offset",
			mcp.Description("The offset to start reading from (ignored if tail is specified)"),
		),
	)
	s.AddTool(getAllocationLogsTool, GetAllocationLogsHandler(nomadClient, logger))
}

// GetAllocationLogsHandler returns a handler for getting allocation logs
func GetAllocationLogsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arguments, ok := request.Params.Arguments.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Invalid arguments"), nil
		}

		allocID, ok := arguments["allocation_id"].(string)
		if !ok || allocID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		task, ok := arguments["task"].(string)
		if !ok || task == "" {
			return mcp.NewToolResultError("task is required"), nil
		}

		logType := "stdout"
		if lt, ok := arguments["type"].(string); ok && lt != "" {
			logType = lt
		}

		follow := false
		if f, ok := arguments["follow"].(bool); ok {
			follow = f
		}

		tail := int64(0)
		if t, ok := arguments["tail"].(float64); ok {
			tail = int64(t)
		}

		offset := int64(0)
		if o, ok := arguments["offset"].(float64); ok {
			offset = int64(o)
		}

		logs, err := client.GetAllocationLogs(allocID, task, logType, follow, tail, offset)
		if err != nil {
			logger.Printf("Error getting allocation logs: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get allocation logs", err), nil
		}

		result := map[string]string{
			"logs": logs,
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format logs", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
