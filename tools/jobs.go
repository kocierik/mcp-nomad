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

	// Get job deployments tool
	getJobDeploymentsTool := mcp.NewTool("get_job_deployments",
		mcp.WithDescription("Get deployments for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to get deployments for"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobDeploymentsTool, GetJobDeploymentsHandler(nomadClient, logger))

	// Get job summary tool
	getJobSummaryTool := mcp.NewTool("get_job_summary",
		mcp.WithDescription("Get summary for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to get summary for"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobSummaryTool, GetJobSummaryHandler(nomadClient, logger))

	// Get job services tool
	getJobServicesTool := mcp.NewTool("get_job_services",
		mcp.WithDescription("Get services for a job"),
		mcp.WithString("job_id",
			mcp.Required(),
			mcp.Description("The ID of the job to get services for"),
		),
		mcp.WithString("namespace",
			mcp.Description("The namespace of the job (default: default)"),
		),
	)
	s.AddTool(getJobServicesTool, GetJobServicesHandler(nomadClient, logger))
}

// ListJobsHandler returns a handler for listing jobs
func ListJobsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		statusFilter := ""
		if s, ok := request.Params.Arguments["status"].(string); ok && s != "" {
			statusFilter = s
		}

		initialJobStubs, err := client.ListJobs(namespace, statusFilter)
		if err != nil {
			logger.Printf("Error listing initial jobs: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to list jobs", err), nil
		}

		type EnhancedJobDetail struct {
			ID                string                   `json:"ID"`
			ParentID          string                   `json:"ParentID"`
			Name              string                   `json:"Name"`
			Type              string                   `json:"Type"`
			Priority          int                      `json:"Priority"`
			Status            string                   `json:"Status"`
			StatusDescription string                   `json:"StatusDescription"`
			JobSummary        *types.JobSummaryDetails `json:"JobSummary"`
			CreateIndex       int                      `json:"CreateIndex"`
			ModifyIndex       int                      `json:"ModifyIndex"`
			JobModifyIndex    int                      `json:"JobModifyIndex"`
		}

		var detailedJobs []EnhancedJobDetail

		for _, stub := range initialJobStubs {
			jobID := stub.ID

			fullJob, errJob := client.GetJob(jobID, namespace)
			if errJob != nil {
				logger.Printf("Error getting full details for job %s in namespace %s: %v. Skipping this job.", jobID, namespace, errJob)
				continue
			}

			item := EnhancedJobDetail{
				ID:                fullJob.ID,
				ParentID:          fullJob.ParentID,
				Name:              fullJob.Name,
				Type:              fullJob.Type,
				Priority:          fullJob.Priority,
				Status:            fullJob.Status,
				StatusDescription: "",
				CreateIndex:       fullJob.CreateIndex,
				ModifyIndex:       fullJob.ModifyIndex,
				JobModifyIndex:    fullJob.JobModifyIndex,
				JobSummary:        nil,
			}

			basicSummaryValue, errSummary := client.GetJobSummary(jobID, namespace)
			if errSummary == nil {
				detailedSummaryForOutput := types.JobSummaryDetails{
					JobID:       fullJob.ID,
					Namespace:   namespace,
					Summary:     basicSummaryValue.Summary,
					Children:    basicSummaryValue.Children,
					CreateIndex: basicSummaryValue.CreateIndex,
					ModifyIndex: basicSummaryValue.ModifyIndex,
				}
				item.JobSummary = &detailedSummaryForOutput
			} else {
				logger.Printf("Error getting summary for job %s in namespace %s: %v. JobSummary will be null.", jobID, namespace, errSummary)
			}

			detailedJobs = append(detailedJobs, item)
		}

		jobsJSON, err := json.MarshalIndent(detailedJobs, "", "  ")
		if err != nil {
			logger.Printf("Error marshalling detailed job list: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to format detailed job list", err), nil
		}

		return mcp.NewToolResultText(string(jobsJSON)), nil
	}
}

// GetJobHandler returns a handler for getting job details
func GetJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		job, err := client.GetJob(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job", err), nil
		}

		jobJSON, err := json.MarshalIndent(job, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format job", err), nil
		}

		return mcp.NewToolResultText(string(jobJSON)), nil
	}
}

// RunJobHandler returns a handler for running a job
func RunJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobSpec, ok := request.Params.Arguments["job_spec"].(string)
		if !ok || jobSpec == "" {
			return mcp.NewToolResultError("job_spec is required"), nil
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
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

// StopJobHandler returns a handler for stopping a job
func StopJobHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
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
			logger.Printf("Error stopping job: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to stop job", err), nil
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
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
			return mcp.NewToolResultError("count is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		err := client.ScaleTaskGroup(jobID, group, int(count), namespace)
		if err != nil {
			logger.Printf("Error scaling job: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to scale job", err), nil
		}

		result := map[string]string{
			"message": fmt.Sprintf("Successfully scaled job %s task group %s to %d", jobID, group, int(count)),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format result", err), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
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

		allocations, err := client.ListJobAllocations(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job allocations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job allocations", err), nil
		}

		allocationsJSON, err := json.MarshalIndent(allocations, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format allocations", err), nil
		}

		return mcp.NewToolResultText(string(allocationsJSON)), nil
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

		evaluations, err := client.ListJobEvaluations(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job evaluations: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job evaluations", err), nil
		}

		evaluationsJSON, err := json.MarshalIndent(evaluations, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format evaluations", err), nil
		}

		return mcp.NewToolResultText(string(evaluationsJSON)), nil
	}
}

// GetJobDeploymentsHandler returns a handler for getting job deployments
func GetJobDeploymentsHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		deployments, err := client.ListJobDeployments(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job deployments: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job deployments", err), nil
		}

		deploymentsJSON, err := json.MarshalIndent(deployments, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format deployments", err), nil
		}

		return mcp.NewToolResultText(string(deploymentsJSON)), nil
	}
}

// GetJobSummaryHandler returns a handler for getting job summary
func GetJobSummaryHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		summary, err := client.GetJobSummary(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job summary: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job summary", err), nil
		}

		summaryJSON, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format job summary", err), nil
		}

		return mcp.NewToolResultText(string(summaryJSON)), nil
	}
}

// GetJobServicesHandler returns a handler for getting job services
func GetJobServicesHandler(client *utils.NomadClient, logger *log.Logger) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		jobID, ok := request.Params.Arguments["job_id"].(string)
		if !ok || jobID == "" {
			return mcp.NewToolResultError("job_id is required"), nil
		}

		namespace := "default"
		if ns, ok := request.Params.Arguments["namespace"].(string); ok && ns != "" {
			namespace = ns
		}

		services, err := client.ListJobServices(jobID, namespace)
		if err != nil {
			logger.Printf("Error getting job services: %v", err)
			return mcp.NewToolResultErrorFromErr("Failed to get job services", err), nil
		}

		servicesJSON, err := json.MarshalIndent(services, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format job services", err), nil
		}

		return mcp.NewToolResultText(string(servicesJSON)), nil
	}
}
