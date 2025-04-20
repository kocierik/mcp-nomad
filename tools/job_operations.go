package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// ScaleJobHandler returns a handler for scaling a job
func ScaleJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		group, ok := request.Params.Arguments["group"].(string)
		if !ok || group == "" {
			return mcp.NewToolResultError("group is required"), nil
		}

		count, ok := request.Params.Arguments["count"].(float64)
		if !ok {
			return mcp.NewToolResultError("count is required and must be a number"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		path := fmt.Sprintf("job/%s/scale", jobID)
		if namespace != "default" {
			path = fmt.Sprintf("namespace/%s/job/%s/scale", namespace, jobID)
		}

		scaleRequest := map[string]interface{}{
			"Count": count,
			"Target": map[string]interface{}{
				"Group": group,
			},
			"Meta": map[string]string{
				"reason": "Scaled via API",
			},
		}

		body, err := client.MakeRequest("POST", path, nil, scaleRequest)
		if err != nil {
			logger.Printf("Error scaling job: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to scale job", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// GetJobAllocationsHandler returns a handler for getting job allocations
func GetJobAllocationsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		path := fmt.Sprintf("job/%s/allocations", jobID)
		if namespace != "default" {
			path = fmt.Sprintf("namespace/%s/job/%s/allocations", namespace, jobID)
		}

		body, err := client.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting job allocations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job allocations", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}

// GetJobEvaluationsHandler returns a handler for getting job evaluations
func GetJobEvaluationsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		path := fmt.Sprintf("job/%s/evaluations", jobID)
		if namespace != "default" {
			path = fmt.Sprintf("namespace/%s/job/%s/evaluations", namespace, jobID)
		}

		body, err := client.MakeRequest("GET", path, nil, nil)
		if err != nil {
			logger.Printf("Error getting job evaluations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job evaluations", err), nil
		}

		return mcp.NewToolResultText(string(body)), nil
	}
}
