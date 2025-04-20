// File: utils/client.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kocierik/nomad-mcp-server/types"
)

// NomadClient handles interactions with the Nomad API
type NomadClient struct {
	address    string
	token      string
	httpClient *http.Client
}

// NewNomadClient creates a new Nomad client
func NewNomadClient() (*NomadClient, error) {
	address := os.Getenv("NOMAD_ADDR")
	if address == "" {
		address = "http://localhost:4646"
	}

	// Get Nomad token from environment if available
	token := os.Getenv("NOMAD_TOKEN")

	return &NomadClient{
		address: address,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// makeRequest is a helper function to make HTTP requests to the Nomad API
func (c *NomadClient) makeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/%s", c.address, path)

	// Add query parameters
	if len(queryParams) > 0 {
		queryParts := make([]string, 0, len(queryParams))
		for key, value := range queryParams {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(queryParts, "&"))
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-Nomad-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ListJobs lists jobs in the specified namespace
func (c *NomadClient) ListJobs(namespace, status string) ([]types.JobSummary, error) {
	path := "jobs"
	if namespace != "" && namespace != "default" {
		path = fmt.Sprintf("namespace/%s/jobs", namespace)
	}

	queryParams := make(map[string]string)
	if status != "" {
		queryParams["status"] = status
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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
func (c *NomadClient) GetJob(jobID, namespace string) (types.Job, error) {
	path := fmt.Sprintf("job/%s", jobID)
	if namespace != "" && namespace != "default" {
		path = fmt.Sprintf("namespace/%s/job/%s", namespace, jobID)
	}

	respBody, err := c.makeRequest("GET", path, nil, nil)
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
func (c *NomadClient) RunJob(jobSpec string, detach bool) (map[string]interface{}, error) {
	// Try to parse as JSON first
	var jobData interface{}
	if err := json.Unmarshal([]byte(jobSpec), &jobData); err != nil {
		// If not JSON, assume it's HCL and use Nomad's HCL parser endpoint
		path := "jobs/parse"
		parseRequest := map[string]string{
			"JobHCL": jobSpec,
		}

		// First parse the HCL to validate and convert to JSON
		parseResp, err := c.makeRequest("POST", path, nil, parseRequest)
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

	respBody, err := c.makeRequest("POST", "jobs", queryParams, jobRequest)
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
func (c *NomadClient) StopJob(jobID, namespace string, purge bool) (map[string]interface{}, error) {
	path := fmt.Sprintf("job/%s", jobID)
	if namespace != "" && namespace != "default" {
		path = fmt.Sprintf("namespace/%s/job/%s", namespace, jobID)
	}

	queryParams := map[string]string{}
	if purge {
		queryParams["purge"] = "true"
	}

	respBody, err := c.makeRequest("DELETE", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

// ListDeployments lists all deployments
func (c *NomadClient) ListDeployments(namespace string) ([]types.DeploymentSummary, error) {
	path := "deployments"
	if namespace != "" && namespace != "default" {
		path = fmt.Sprintf("namespace/%s/deployments", namespace)
	}

	respBody, err := c.makeRequest("GET", path, nil, nil)
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
func (c *NomadClient) GetDeployment(deploymentID string) (types.Deployment, error) {
	path := fmt.Sprintf("deployment/%s", deploymentID)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.Deployment{}, err
	}

	var deployment types.Deployment
	if err := json.Unmarshal(respBody, &deployment); err != nil {
		return types.Deployment{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployment, nil
}

// ListNamespaces lists all namespaces
func (c *NomadClient) ListNamespaces() ([]types.Namespace, error) {
	respBody, err := c.makeRequest("GET", "namespaces", nil, nil)
	if err != nil {
		return nil, err
	}

	var namespaces []types.Namespace
	if err := json.Unmarshal(respBody, &namespaces); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return namespaces, nil
}

// CreateNamespace creates a new namespace
func (c *NomadClient) CreateNamespace(namespace types.Namespace) error {
	_, err := c.makeRequest("POST", "namespace", nil, namespace)
	return err
}

// DeleteNamespace deletes a namespace
func (c *NomadClient) DeleteNamespace(name string) error {
	path := fmt.Sprintf("namespace/%s", name)
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// ListNodes lists all nodes in the cluster
func (c *NomadClient) ListNodes(status string) ([]types.NodeSummary, error) {
	queryParams := make(map[string]string)
	if status != "" {
		queryParams["status"] = status
	}

	respBody, err := c.makeRequest("GET", "nodes", queryParams, nil)
	if err != nil {
		return nil, err
	}

	var nodes []types.NodeSummary
	if err := json.Unmarshal(respBody, &nodes); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return nodes, nil
}

// GetNode retrieves a specific node by ID
func (c *NomadClient) GetNode(nodeID string) (types.Node, error) {
	path := fmt.Sprintf("node/%s", nodeID)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.Node{}, err
	}

	var node types.Node
	if err := json.Unmarshal(respBody, &node); err != nil {
		return types.Node{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return node, nil
}

// DrainNode enables or disables drain mode for a node
func (c *NomadClient) DrainNode(nodeID string, enable bool, deadline int64) (string, error) {
	path := fmt.Sprintf("node/%s/drain", nodeID)

	drainSpec := map[string]interface{}{
		"DrainSpec": map[string]interface{}{
			"Deadline":         deadline,
			"IgnoreSystemJobs": false,
		},
		"Meta": map[string]string{
			"reason": "Initiated via API",
		},
	}

	if !enable {
		drainSpec = map[string]interface{}{
			"DrainSpec": nil,
			"Meta": map[string]string{
				"reason": "Drain disabled via API",
			},
		}
	}

	respBody, err := c.makeRequest("POST", path, nil, drainSpec)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if enable {
		if deadline > 0 {
			return fmt.Sprintf("Node drain enabled with deadline %d seconds", deadline), nil
		}
		return "Node drain enabled with no deadline", nil
	}
	return "Node drain disabled", nil
}

// MakeRequest is a helper function to make HTTP requests to the Nomad API
func (c *NomadClient) MakeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	return c.makeRequest(method, path, queryParams, body)
}
