package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/nomad-mcp-server/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// ListJobsHandler returns a handler for the list_jobs tool
func ListJobsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		status := ""
		if s, ok := request.Params.Arguments["status"].(string); ok {
			status = s
		}

		jobs, err := client.ListJobs(namespace, status)
		if err != nil {
			logger.Printf("Error listing jobs: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list jobs", err), nil
		}

		jobsJSON, err := json.MarshalIndent(jobs, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format job list", err), nil
		}

		return mcp.NewToolResultText(string(jobsJSON)), nil
	}
}

// GetJobHandler returns a handler for the get_job tool
func GetJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("Job ID is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		job, err := client.GetJob(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job %s: %v", jobID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to get job %s", jobID), err), nil
		}

		jobJSON, err := json.MarshalIndent(job, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format job details", err), nil
		}

		return mcp.NewToolResultText(string(jobJSON)), nil
	}
}

// RunJobHandler returns a handler for the run_job tool
func RunJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobSpec, ok := request.Params.Arguments["job_spec"].(string)
		if !ok || jobSpec == "" {
			return mcp.NewToolResultError("Job specification is required"), nil
		}

		detach := false
		if d, ok := request.Params.Arguments["detach"].(bool); ok {
			detach = d
		}

		result, err := client.RunJob(jobSpec, detach)
		if err != nil {
			logger.Printf("Error running job: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to run job", err), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format response", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// StopJobHandler returns a handler for the stop_job tool
func StopJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("Job ID is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		purge := false
		if p, ok := request.Params.Arguments["purge"].(bool); ok {
			purge = p
		}

		result, err := client.StopJob(jobID, namespace, purge)
		if err != nil {
			logger.Printf("Error stopping job %s: %v", jobID, err)
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("Failed to stop job %s", jobID), err), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format response", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}
