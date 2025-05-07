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
	// Get task logs tool
	getTaskLogsTool := mcp.NewTool("get_task_logs",
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
	s.AddTool(getTaskLogsTool, GetTaskLogsHandler(nomadClient, logger))
}

// GetTaskLogsHandler returns a handler for getting task logs
func GetTaskLogsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		allocID, ok := request.Params.Arguments["allocation_id"].(string)
		if !ok || allocID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		task, ok := request.Params.Arguments["task"].(string)
		if !ok || task == "" {
			return mcp.NewToolResultError("task is required"), nil
		}

		logType := "stdout"
		if lt, ok := request.Params.Arguments["log_type"].(string); ok && lt != "" {
			logType = lt
		}

		follow := false
		if f, ok := request.Params.Arguments["follow"].(bool); ok {
			follow = f
		}

		tail := int64(0)
		if t, ok := request.Params.Arguments["tail"].(float64); ok {
			tail = int64(t)
		}

		offset := int64(0)
		if o, ok := request.Params.Arguments["offset"].(float64); ok {
			offset = int64(o)
		}

		logs, err := client.GetTaskLogs(allocID, task, logType, follow, tail, offset)
		if err != nil {
			logger.Printf("Error getting task logs: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get task logs", err), nil
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

// GetAllocationLogsHandler returns a handler for getting allocation logs
func GetAllocationLogsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		allocID, ok := request.Params.Arguments["allocation_id"].(string)
		if !ok || allocID == "" {
			return mcp.NewToolResultError("allocation_id is required"), nil
		}

		logs, err := client.GetAllocationLogs(allocID)
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
