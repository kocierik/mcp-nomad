package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kocierik/mcp-nomad/utils"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterJobTools registers all job-related tools
func RegisterJobTools(s *server.MCPServer, nomadClient *utils.NomadClient, logger *log.Logger) {
	// List jobs tool
	listJobsTool := mcp.NewTool("list_jobs",
		mcp.WithDescription("List all jobs in Nomad"),
		mcp.WithString("namespace",
			mcp.Description("The namespace to list jobs from (default: default)"),
		),
		mcp.WithString("status",
			mcp.Description("Filter jobs by status (pending, running, dead)"),
			mcp.Enum("pending", "running", "dead", ""),
		),
	)
	s.AddTool(listJobsTool, ListJobsHandler(nomadClient, logger))

	// Get job tool
	getJobTool := mcp.NewTool("get_job",
		mcp.WithDescription("Get job details by ID"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to retrieve"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobTool, GetJobHandler(nomadClient, logger))

	// Run job tool
	runJobTool := mcp.NewTool("run_job",
		mcp.WithDescription("Run a new job or update an existing job"),
		mcp.WithString("job_spec",
			mcp.Required(),
			mcp.Description("The job specification in HCL or JSON format"),
		),
		mcp.WithBoolean("detach",
			mcp.Description("Return immediately instead of monitoring deployment"),
		),
	)
	s.AddTool(runJobTool, RunJobHandler(nomadClient, logger))

	// Stop job tool
	stopJobTool := mcp.NewTool("stop_job",
		mcp.WithDescription("Stop a running job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to stop"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
		mcp.WithBoolean("purge",
			mcp.Description("Purge the job from Nomad instead of just stopping it"),
		),
	)
	s.AddTool(stopJobTool, StopJobHandler(nomadClient, logger))

	// Scale job tool
	scaleJobTool := mcp.NewTool("scale_job",
		mcp.WithDescription("Scale a job's task group"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to scale"),
		),
		mcp.WithString("group",
			mcp.Required(),
			mcp.Description("The task group to scale"),
		),
		mcp.WithNumber("count",
			mcp.Required(),
			mcp.Description("The new count for the task group"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(scaleJobTool, ScaleJobHandler(nomadClient, logger))

	// Get job allocations tool
	getJobAllocationsTool := mcp.NewTool("get_job_allocations",
		mcp.WithDescription("Get allocations for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to get allocations for"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobAllocationsTool, GetJobAllocationsHandler(nomadClient, logger))

	// Get job evaluations tool
	getJobEvaluationsTool := mcp.NewTool("get_job_evaluations",
		mcp.WithDescription("Get evaluations for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to get evaluations for"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobEvaluationsTool, GetJobEvaluationsHandler(nomadClient, logger))
}

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
