package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kocierik/mcp-nomad/types"
)

// ListDeployments lists all deployments
func (c *NomadClient) ListDeployments(ctx context.Context, namespace string) ([]types.DeploymentSummary, error) {
	path := "deployments"

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var deployments []types.DeploymentSummary
	if err := json.Unmarshal(respBody, &deployments); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployments, nil
}

// GetDeployment retrieves a specific deployment
func (c *NomadClient) GetDeployment(ctx context.Context, deploymentID string) (types.Deployment, error) {
	path := fmt.Sprintf("deployment/%s", deploymentID)

	respBody, err := c.makeRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		return types.Deployment{}, err
	}

	var deployment types.Deployment
	if err := json.Unmarshal(respBody, &deployment); err != nil {
		return types.Deployment{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployment, nil
}

// ListJobDeployments lists all deployments for a job
func (c *NomadClient) ListJobDeployments(ctx context.Context, jobID, namespace string) ([]types.JobDeployment, error) {
	path := fmt.Sprintf("job/%s/deployments", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var deployments []types.JobDeployment
	if err := json.Unmarshal(respBody, &deployments); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployments, nil
}

// GetJobDeployment retrieves the most recent deployment for a job
func (c *NomadClient) GetJobDeployment(ctx context.Context, jobID, namespace string) (types.JobDeployment, error) {
	path := fmt.Sprintf("job/%s/deployment", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest(ctx, "GET", path, queryParams, nil)
	if err != nil {
		return types.JobDeployment{}, err
	}

	var deployment types.JobDeployment
	if err := json.Unmarshal(respBody, &deployment); err != nil {
		return types.JobDeployment{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployment, nil
}
