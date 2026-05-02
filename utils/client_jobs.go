package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListJobs lists jobs in the specified namespace
func (c *NomadClient) ListJobs(ctx context.Context, namespace, status string) ([]types.JobSummary, error) {
	path := "jobs"

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
	if status != "" {
		queryParams["status"] = status
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var jobs []types.JobSummary
	if err := json.Unmarshal(respBody, &jobs); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return jobs, nil
}

// GetJob retrieves a specific job by ID
func (c *NomadClient) GetJob(ctx context.Context, jobID, namespace string) (types.Job, error) {
	path := fmt.Sprintf("job/%s", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return types.Job{}, err
	}

	var job types.Job
	if err := json.Unmarshal(respBody, &job); err != nil {
		return types.Job{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return job, nil
}

// RunJob submits a job to Nomad
func (c *NomadClient) RunJob(ctx context.Context, jobSpec string, detach bool) (map[string]interface{}, error) {
	// Try to parse as JSON first
	var jobData interface{}
	if err := json.Unmarshal([]byte(jobSpec), &jobData); err != nil {
		// If not JSON, assume it's HCL and use Nomad's HCL parser endpoint
		path := "jobs/parse"
		parseRequest := map[string]string{
			"JobHCL": jobSpec,
		}

		// First parse the HCL to validate and convert to JSON
		parseResp, err := c.makeRequest(ctx, "POST", path, nil, parseRequest)
		if err != nil {
			return nil, fmt.Errorf("error parsing HCL job spec: %v", err)
		}

		// Unmarshal the parse response into a map
		var parsedJob map[string]interface{}
		if err := json.Unmarshal(parseResp, &parsedJob); err != nil {
			return nil, fmt.Errorf("error unmarshaling parsed job spec: %v", err)
		}

		// Use the parsed job data
		jobData = parsedJob
	}

	// Wrap the job data in a Job field as required by the Nomad API
	jobRequest := map[string]interface{}{
		"Job": jobData,
	}

	queryParams := map[string]string{}
	if detach {
		queryParams["detach"] = "true"
	}

	respBody, err := c.makeRequest(ctx, "POST", "jobs", queryParams, jobRequest)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

// StopJob stops a job
func (c *NomadClient) StopJob(ctx context.Context, jobID, namespace string, purge bool) (map[string]interface{}, error) {
	path := fmt.Sprintf("job/%s", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
	if purge {
		queryParams["purge"] = "true"
	}

	respBody, err := c.makeRequest(ctx, "DELETE", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

// GetJobVersions returns the versions of a job
func (c *NomadClient) GetJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error) {
	path := fmt.Sprintf("/v1/job/%s/versions", jobID)
	if namespace != "" {
		path = fmt.Sprintf("%s?namespace=%s", path, namespace)
	}

	var versions []types.Job
	err := c.get(ctx, path, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// GetJobSubmission retrieves the original job submission
func (c *NomadClient) GetJobSubmission(ctx context.Context, jobID, namespace string) (string, error) {
	path := fmt.Sprintf("job/%s/submission", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// ListJobVersions lists all versions of a job
func (c *NomadClient) ListJobVersions(ctx context.Context, jobID, namespace string) ([]types.Job, error) {
	path := fmt.Sprintf("job/%s/versions", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var versions []types.Job
	if err := json.Unmarshal(respBody, &versions); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return versions, nil
}

// ListJobAllocations lists all allocations for a job
func (c *NomadClient) ListJobAllocations(ctx context.Context, jobID, namespace string) ([]types.Allocation, error) {
	path := fmt.Sprintf("job/%s/allocations", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var allocations []types.Allocation
	if err := json.Unmarshal(respBody, &allocations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return allocations, nil
}

// ListJobEvaluations lists all evaluations for a job
func (c *NomadClient) ListJobEvaluations(ctx context.Context, jobID, namespace string) ([]types.Evaluation, error) {
	path := fmt.Sprintf("job/%s/evaluations", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var evaluations []types.Evaluation
	if err := json.Unmarshal(respBody, &evaluations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return evaluations, nil
}

// UpdateJob updates an existing job
func (c *NomadClient) UpdateJob(ctx context.Context, job types.Job, enforceIndex bool) error {
	path := "jobs"
	if enforceIndex {
		path = fmt.Sprintf("%s?enforce_index=true", path)
	}

	_, err := c.makeRequest(ctx, "POST", path, nil, job)
	return err
}

// DispatchJob dispatches a parameterized job
func (c *NomadClient) DispatchJob(ctx context.Context, jobID string, payload map[string]interface{}, meta map[string]string) (string, error) {
	path := fmt.Sprintf("job/%s/dispatch", jobID)

	request := map[string]interface{}{
		"Payload": payload,
		"Meta":    meta,
	}

	respBody, err := c.makeRequest(ctx, "POST", path, nil, request)
	if err != nil {
		return "", err
	}

	var response struct {
		DispatchedJobID string `json:"DispatchedJobID"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	return response.DispatchedJobID, nil
}

// RevertJob reverts a job to a specific version
func (c *NomadClient) RevertJob(ctx context.Context, jobID string, version int, enforceIndex bool) error {
	path := fmt.Sprintf("job/%s/revert", jobID)
	if enforceIndex {
		path = fmt.Sprintf("%s?enforce_index=true", path)
	}

	request := map[string]interface{}{
		"JobVersion": version,
	}

	_, err := c.makeRequest(ctx, "POST", path, nil, request)
	return err
}

// SetJobStability sets the stability of a job
func (c *NomadClient) SetJobStability(ctx context.Context, jobID string, version int, stable bool) error {
	path := fmt.Sprintf("job/%s/stability", jobID)

	request := map[string]interface{}{
		"JobVersion": version,
		"Stable":     stable,
	}

	_, err := c.makeRequest(ctx, "POST", path, nil, request)
	return err
}

// CreateJobEvaluation forces a new evaluation for a job
func (c *NomadClient) CreateJobEvaluation(ctx context.Context, jobID string) (string, error) {
	path := fmt.Sprintf("job/%s/evaluate", jobID)

	respBody, err := c.makeRequest(ctx, "POST", path, nil, nil)
	if err != nil {
		return "", err
	}

	var response struct {
		EvalID string `json:"EvalID"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", err
	}

	return response.EvalID, nil
}

// CreateJobPlan creates a plan for a job
func (c *NomadClient) CreateJobPlan(ctx context.Context, job types.Job) (types.JobPlan, error) {
	path := "job/plan"

	respBody, err := c.makeRequest(ctx, "POST", path, nil, job)
	if err != nil {
		return types.JobPlan{}, err
	}

	var plan types.JobPlan
	if err := json.Unmarshal(respBody, &plan); err != nil {
		return types.JobPlan{}, err
	}

	return plan, nil
}

// ForceNewPeriodicInstance forces a new instance of a periodic job
func (c *NomadClient) ForceNewPeriodicInstance(ctx context.Context, jobID string) error {
	path := fmt.Sprintf("job/%s/periodic/force", jobID)

	_, err := c.makeRequest(ctx, "POST", path, nil, nil)
	return err
}

// GetJobScaleStatus retrieves the scale status of a job
func (c *NomadClient) GetJobScaleStatus(ctx context.Context, jobID, namespace string) (types.JobScaleStatus, error) {
	path := fmt.Sprintf("job/%s/scale", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return types.JobScaleStatus{}, err
	}

	var status types.JobScaleStatus
	if err := json.Unmarshal(respBody, &status); err != nil {
		return types.JobScaleStatus{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return status, nil
}

// ScaleTaskGroup scales a task group
func (c *NomadClient) ScaleTaskGroup(ctx context.Context, jobID, group string, count int, namespace string) error {
	path := fmt.Sprintf("job/%s/scale", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	request := map[string]interface{}{
		"Count": count,
		"Target": map[string]interface{}{
			"Group": group,
		},
	}

	_, err := c.makeRequest(ctx, "POST", path, queryParams, request)
	return err
}

// ListJobServices lists all services for a job
func (c *NomadClient) ListJobServices(ctx context.Context, jobID, namespace string) ([]types.Service, error) {
	path := fmt.Sprintf("job/%s/services", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var services []types.Service
	if err := json.Unmarshal(respBody, &services); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return services, nil
}

// GetJobSummary retrieves a summary of a job
func (c *NomadClient) GetJobSummary(ctx context.Context, jobID, namespace string) (types.JobSummary, error) {
	path := fmt.Sprintf("job/%s/summary", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return types.JobSummary{}, err
	}

	var response struct {
		ID          string                       `json:"ID"`
		Namespace   string                       `json:"Namespace"`
		Summary     map[string]types.TaskSummary `json:"Summary"`
		Children    *types.JobChildrenSummary    `json:"Children"`
		CreateIndex int                          `json:"CreateIndex"`
		ModifyIndex int                          `json:"ModifyIndex"`
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return types.JobSummary{}, err
	}

	return types.JobSummary{
		ID:          response.ID,
		Summary:     response.Summary,
		Children:    response.Children,
		CreateIndex: response.CreateIndex,
		ModifyIndex: response.ModifyIndex,
	}, nil
}
