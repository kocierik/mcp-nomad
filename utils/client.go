// File: utils/client.go
// Package utils provides utility functions and types for interacting with the Nomad API.
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kocierik/mcp-nomad/types"
)

// NomadClient handles interactions with the Nomad API.
// It provides methods for managing jobs, deployments, namespaces, nodes, allocations,
// variables, volumes, and ACL tokens.
type NomadClient struct {
	address          string
	token            string
	httpClient       *http.Client
	DefaultTailLines int // Default number of lines to show when tailing logs
}

// NewNomadClient creates a new Nomad client with the specified address and token.
// It validates the connection to the Nomad server before returning.
//
// Example:
//
//	client, err := NewNomadClient("http://localhost:4646", "your-token")
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewNomadClient(address, token string) (*NomadClient, error) {
	// Validate the address
	if address == "" {
		return nil, fmt.Errorf("nomad address is required")
	}

	// Create the client
	client := &NomadClient{
		address: address,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		DefaultTailLines: 100, // Default to showing last 100 lines
	}

	// Test the connection
	_, err := client.makeRequest("GET", "status/leader", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Nomad server: %v", err)
	}

	return client, nil
}

// SetToken sets the ACL token for the client
func (c *NomadClient) SetToken(token string) {
	c.token = token
}

// GetToken returns the current ACL token
func (c *NomadClient) GetToken() string {
	return c.token
}

// SetDefaultTailLines sets the default number of lines to show when tailing logs
func (c *NomadClient) SetDefaultTailLines(lines int) error {
	if lines <= 0 {
		return fmt.Errorf("number of lines must be positive")
	}
	c.DefaultTailLines = lines
	return nil
}

// GetDefaultTailLines returns the current default number of lines for log tailing
func (c *NomadClient) GetDefaultTailLines() int {
	return c.DefaultTailLines
}

// makeRequest is a helper function to make HTTP requests to the Nomad API
func (c *NomadClient) makeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	baseURL := fmt.Sprintf("%s/v1/%s", c.address, path)

	// Create url.Values for proper query parameter encoding
	query := url.Values{}
	for key, value := range queryParams {
		query.Set(key, value)
	}

	// Add query parameters to URL
	if len(query) > 0 {
		baseURL = fmt.Sprintf("%s?%s", baseURL, query.Encode())
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add ACL token to headers if available
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

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
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

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
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

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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

// DrainNode enables or disables drain mode for a node
func (c *NomadClient) EligibilityNode(nodeID string, eligible string) (types.NodeSummary, error) {
	path := fmt.Sprintf("node/%s/eligibility", nodeID)

	eligibilitySpec := map[string]interface{}{
		"Eligibility": eligible,
	}

	respBody, err := c.makeRequest("POST", path, nil, eligibilitySpec)
	if err != nil {
		return types.NodeSummary{}, err
	}

	var nodes types.NodeSummary
	if err := json.Unmarshal(respBody, &nodes); err != nil {
		return types.NodeSummary{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return nodes, nil
}

// MakeRequest is a helper function to make HTTP requests to the Nomad API
func (c *NomadClient) MakeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	return c.makeRequest(method, path, queryParams, body)
}

// ListVolumes lists all host volumes
func (c *NomadClient) ListVolumes(nodeID string, pluginID string, nextToken string, perPage int, filter string) ([]types.Volume, error) {
	path := "volumes"
	query := url.Values{}
	if nodeID != "" {
		query.Set("node_id", nodeID)
	}
	if pluginID != "" {
		query.Set("plugin_id", pluginID)
	}
	if nextToken != "" {
		query.Set("next_token", nextToken)
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	var volumes []types.Volume
	if err := c.get(path+"?"+query.Encode(), &volumes); err != nil {
		return nil, fmt.Errorf("error listing volumes: %v", err)
	}

	return volumes, nil
}

// GetVolume retrieves a specific host volume
func (c *NomadClient) GetVolume(volumeID string) (*types.Volume, error) {
	path := fmt.Sprintf("/v1/volume/host/%s", volumeID)
	var volume types.Volume
	if err := c.get(path, &volume); err != nil {
		return nil, fmt.Errorf("error getting volume: %v", err)
	}

	return &volume, nil
}

// DeleteVolume deletes a host volume
func (c *NomadClient) DeleteVolume(volumeID string) error {
	path := fmt.Sprintf("/v1/volume/host/%s/delete", volumeID)
	if err := c.delete(path); err != nil {
		return fmt.Errorf("error deleting volume: %v", err)
	}

	return nil
}

// ListVolumeClaims lists all volume claims
func (c *NomadClient) ListVolumeClaims(namespace string, claimID string, jobID string, taskGroup string, volumeName string, nextToken string, perPage int) ([]types.VolumeClaim, error) {
	path := "volumes/"
	query := url.Values{}
	query.Set("namespace", namespace)

	if claimID != "" {
		query.Set("claim_id", claimID)
	}
	if jobID != "" {
		query.Set("job_id", jobID)
	}
	if taskGroup != "" {
		query.Set("task_group", taskGroup)
	}
	if volumeName != "" {
		query.Set("volume_name", volumeName)
	}
	if nextToken != "" {
		query.Set("next_token", nextToken)
	}
	if perPage > 0 {
		query.Set("per_page", strconv.Itoa(perPage))
	}

	var claims []types.VolumeClaim
	if err := c.get(path+"?"+query.Encode(), &claims); err != nil {
		return nil, fmt.Errorf("error listing volume claims: %v", err)
	}

	return claims, nil
}

// DeleteVolumeClaim deletes a volume claim
func (c *NomadClient) DeleteVolumeClaim(claimID string) error {
	path := fmt.Sprintf("/v1/volumes/claim/%s", claimID)
	if err := c.delete(path); err != nil {
		return fmt.Errorf("error deleting volume claim: %v", err)
	}

	return nil
}

// ListACLTokens lists all ACL tokens
func (c *NomadClient) ListACLTokens() ([]types.ACLToken, error) {
	respBody, err := c.makeRequest("GET", "acl/tokens", nil, nil)
	if err != nil {
		return nil, err
	}

	var tokens []types.ACLToken
	if err := json.Unmarshal(respBody, &tokens); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return tokens, nil
}

// GetACLToken retrieves a specific ACL token by accessor ID
func (c *NomadClient) GetACLToken(accessorID string) (types.ACLToken, error) {
	path := fmt.Sprintf("acl/token/%s", accessorID)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.ACLToken{}, err
	}

	var token types.ACLToken
	if err := json.Unmarshal(respBody, &token); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return token, nil
}

// CreateACLToken creates a new ACL token
func (c *NomadClient) CreateACLToken(token types.ACLToken) (types.ACLToken, error) {
	respBody, err := c.makeRequest("POST", "acl/token", nil, token)
	if err != nil {
		return types.ACLToken{}, err
	}

	var newToken types.ACLToken
	if err := json.Unmarshal(respBody, &newToken); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return newToken, nil
}

// DeleteACLToken deletes an ACL token
func (c *NomadClient) DeleteACLToken(accessorID string) error {
	path := fmt.Sprintf("acl/token/%s", accessorID)
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// ListACLPolicies lists all ACL policies
func (c *NomadClient) ListACLPolicies() ([]types.ACLPolicy, error) {
	respBody, err := c.makeRequest("GET", "acl/policies", nil, nil)
	if err != nil {
		return nil, err
	}

	var policies []types.ACLPolicy
	if err := json.Unmarshal(respBody, &policies); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policies, nil
}

// GetACLPolicy retrieves a specific ACL policy by name
func (c *NomadClient) GetACLPolicy(name string) (types.ACLPolicy, error) {
	path := fmt.Sprintf("acl/policy/%s", name)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.ACLPolicy{}, err
	}

	var policy types.ACLPolicy
	if err := json.Unmarshal(respBody, &policy); err != nil {
		return types.ACLPolicy{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policy, nil
}

// CreateACLPolicy creates a new ACL policy OK

func (c *NomadClient) CreateACLPolicy(policy types.ACLPolicy) error {
	path := fmt.Sprintf("acl/policy/%s", policy.Name)

	_, err := c.makeRequest("POST", path, nil, policy)
	return err
}

// DeleteACLPolicy deletes an ACL policy
func (c *NomadClient) DeleteACLPolicy(name string) error {
	path := fmt.Sprintf("acl/policy/%s", name)
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// ListACLRoles lists all ACL roles
func (c *NomadClient) ListACLRoles() ([]types.ACLRole, error) {
	respBody, err := c.makeRequest("GET", "acl/roles", nil, nil)
	if err != nil {
		return nil, err
	}

	var roles []types.ACLRole
	if err := json.Unmarshal(respBody, &roles); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return roles, nil
}

// GetACLRole retrieves a specific ACL role by ID
func (c *NomadClient) GetACLRole(id string) (types.ACLRole, error) {
	path := fmt.Sprintf("acl/role/%s", id)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.ACLRole{}, err
	}

	var role types.ACLRole
	if err := json.Unmarshal(respBody, &role); err != nil {
		return types.ACLRole{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return role, nil
}

// CreateACLRole creates a new ACL role
func (c *NomadClient) CreateACLRole(role types.ACLRole) (types.ACLRole, error) {
	respBody, err := c.makeRequest("POST", "acl/role", nil, role)
	if err != nil {
		return types.ACLRole{}, err
	}

	var newRole types.ACLRole
	if err := json.Unmarshal(respBody, &newRole); err != nil {
		return types.ACLRole{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return newRole, nil
}

// DeleteACLRole deletes an ACL role
func (c *NomadClient) DeleteACLRole(id string) error {
	path := fmt.Sprintf("acl/role/%s", id)
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// Helper methods for HTTP requests
func (c *NomadClient) get(path string, result interface{}) error {
	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, result)
}

func (c *NomadClient) delete(path string) error {
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// VolumeClaim represents a volume claim in Nomad
type VolumeClaim struct {
	AllocID       string `json:"AllocID"`
	CreateIndex   int    `json:"CreateIndex"`
	ID            string `json:"ID"`
	JobID         string `json:"JobID"`
	ModifyIndex   int    `json:"ModifyIndex"`
	Namespace     string `json:"Namespace"`
	TaskGroupName string `json:"TaskGroupName"`
	VolumeID      string `json:"VolumeID"`
	VolumeName    string `json:"VolumeName"`
}

// BootstrapACLToken bootstraps the ACL system and returns the initial management token
func (c *NomadClient) BootstrapACLToken() (types.ACLToken, error) {
	respBody, err := c.makeRequest("POST", "acl/bootstrap", nil, nil)
	if err != nil {
		return types.ACLToken{}, err
	}

	var token types.ACLToken
	if err := json.Unmarshal(respBody, &token); err != nil {
		return types.ACLToken{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return token, nil
}

// GetAllocationLogs retrieves logs from a specific task in an allocation
func (c *NomadClient) GetAllocationLogs(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
	if allocID == "" {
		return "", fmt.Errorf("allocation ID is required")
	}
	if task == "" {
		return "", fmt.Errorf("task name is required")
	}

	// Set default log type if not specified
	if logType == "" {
		logType = "stdout"
	}

	// Build query parameters
	queryParams := map[string]string{
		"task":   task,
		"type":   logType,
		"follow": fmt.Sprintf("%v", follow),
		"plain":  "true",
	}

	// If tail is specified, we want to read from the end
	if tail > 0 {
		queryParams["origin"] = "end"
		// Estimate bytes needed for tail lines (assume average 200 bytes per line)
		estimatedBytes := tail * 200
		queryParams["offset"] = fmt.Sprintf("%d", estimatedBytes)
	} else if offset > 0 {
		queryParams["offset"] = fmt.Sprintf("%d", offset)
	}

	// Make request to Nomad API
	path := fmt.Sprintf("client/fs/logs/%s", allocID)
	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get allocation logs: %v", err)
	}

	// If tail was specified, we need to process the response to get the correct number of lines
	if tail > 0 {
		lines := strings.Split(string(respBody), "\n")
		if len(lines) > int(tail) {
			// Take only the last 'tail' lines
			lines = lines[len(lines)-int(tail):]
		}
		return strings.Join(lines, "\n"), nil
	}

	return string(respBody), nil
}

// GetClusterLeader return the info of the cluster leader
func (c *NomadClient) GetClusterLeader() ([]byte, error) {
	respBody, err := c.makeRequest("GET", "operator/raft/configuration", nil, nil)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// ListClusterPeers return the list of the cluster nodes
func (c *NomadClient) ListClusterPeers() ([]byte, error) {
	respBody, err := c.makeRequest("GET", "operator/raft/configuration", nil, nil)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// ListRegions return the regions listed
func (c *NomadClient) ListRegions() ([]byte, error) {
	respBody, err := c.MakeRequest("GET", "regions", nil, nil) // Check Nomad API for correct endpoint
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// GetJobVersions returns the versions of a job
func (c *NomadClient) GetJobVersions(jobID, namespace string) ([]types.Job, error) {
	path := fmt.Sprintf("/v1/job/%s/versions", jobID)
	if namespace != "" {
		path = fmt.Sprintf("%s?namespace=%s", path, namespace)
	}

	var versions []types.Job
	err := c.get(path, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// GetAllocation returns the details of an allocation
func (c *NomadClient) GetAllocation(allocID string) (types.Allocation, error) {
	path := fmt.Sprintf("allocation/%s", allocID)

	var alloc types.Allocation
	err := c.get(path, &alloc)
	if err != nil {
		return types.Allocation{}, err
	}

	return alloc, nil
}

// ListSentinelPolicies lists all Sentinel policies
func (c *NomadClient) ListSentinelPolicies() ([]types.SentinelPolicy, error) {
	respBody, err := c.makeRequest("GET", "sentinel/policies", nil, nil)
	if err != nil {
		return nil, err
	}

	var policies []types.SentinelPolicy
	if err := json.Unmarshal(respBody, &policies); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policies, nil
}

// GetSentinelPolicy retrieves a specific Sentinel policy by name
func (c *NomadClient) GetSentinelPolicy(name string) (types.SentinelPolicy, error) {
	path := fmt.Sprintf("sentinel/policy/%s", name)

	respBody, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return types.SentinelPolicy{}, err
	}

	var policy types.SentinelPolicy
	if err := json.Unmarshal(respBody, &policy); err != nil {
		return types.SentinelPolicy{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return policy, nil
}

// CreateSentinelPolicy creates a new Sentinel policy
func (c *NomadClient) CreateSentinelPolicy(policy types.SentinelPolicy) error {
	path := fmt.Sprintf("sentinel/policy/%s", policy.Name)
	_, err := c.makeRequest("POST", path, nil, policy)
	return err
}

// DeleteSentinelPolicy deletes a Sentinel policy
func (c *NomadClient) DeleteSentinelPolicy(name string) error {
	path := fmt.Sprintf("sentinel/policy/%s", name)
	_, err := c.makeRequest("DELETE", path, nil, nil)
	return err
}

// GetJobSubmission retrieves the original job submission
func (c *NomadClient) GetJobSubmission(jobID, namespace string) (string, error) {
	path := fmt.Sprintf("job/%s/submission", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// ListJobVersions lists all versions of a job
func (c *NomadClient) ListJobVersions(jobID, namespace string) ([]types.Job, error) {
	path := fmt.Sprintf("job/%s/versions", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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
func (c *NomadClient) ListJobAllocations(jobID, namespace string) ([]types.Allocation, error) {
	path := fmt.Sprintf("job/%s/allocations", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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
func (c *NomadClient) ListJobEvaluations(jobID, namespace string) ([]types.Evaluation, error) {
	path := fmt.Sprintf("job/%s/evaluations", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var evaluations []types.Evaluation
	if err := json.Unmarshal(respBody, &evaluations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return evaluations, nil
}

// ListJobDeployments lists all deployments for a job
func (c *NomadClient) ListJobDeployments(jobID, namespace string) ([]types.JobDeployment, error) {
	path := fmt.Sprintf("job/%s/deployments", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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
func (c *NomadClient) GetJobDeployment(jobID, namespace string) (types.JobDeployment, error) {
	path := fmt.Sprintf("job/%s/deployment", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return types.JobDeployment{}, err
	}

	var deployment types.JobDeployment
	if err := json.Unmarshal(respBody, &deployment); err != nil {
		return types.JobDeployment{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return deployment, nil
}

// GetJobSummary retrieves a summary of a job
func (c *NomadClient) GetJobSummary(jobID, namespace string) (types.JobSummary, error) {
	path := fmt.Sprintf("job/%s/summary", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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

// UpdateJob updates an existing job
func (c *NomadClient) UpdateJob(job types.Job, enforceIndex bool) error {
	path := "jobs"
	if enforceIndex {
		path = fmt.Sprintf("%s?enforce_index=true", path)
	}

	_, err := c.makeRequest("POST", path, nil, job)
	return err
}

// DispatchJob dispatches a parameterized job
func (c *NomadClient) DispatchJob(jobID string, payload map[string]interface{}, meta map[string]string) (string, error) {
	path := fmt.Sprintf("job/%s/dispatch", jobID)

	request := map[string]interface{}{
		"Payload": payload,
		"Meta":    meta,
	}

	respBody, err := c.makeRequest("POST", path, nil, request)
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
func (c *NomadClient) RevertJob(jobID string, version int, enforceIndex bool) error {
	path := fmt.Sprintf("job/%s/revert", jobID)
	if enforceIndex {
		path = fmt.Sprintf("%s?enforce_index=true", path)
	}

	request := map[string]interface{}{
		"JobVersion": version,
	}

	_, err := c.makeRequest("POST", path, nil, request)
	return err
}

// SetJobStability sets the stability of a job
func (c *NomadClient) SetJobStability(jobID string, version int, stable bool) error {
	path := fmt.Sprintf("job/%s/stability", jobID)

	request := map[string]interface{}{
		"JobVersion": version,
		"Stable":     stable,
	}

	_, err := c.makeRequest("POST", path, nil, request)
	return err
}

// CreateJobEvaluation forces a new evaluation for a job
func (c *NomadClient) CreateJobEvaluation(jobID string) (string, error) {
	path := fmt.Sprintf("job/%s/evaluate", jobID)

	respBody, err := c.makeRequest("POST", path, nil, nil)
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
func (c *NomadClient) CreateJobPlan(job types.Job) (types.JobPlan, error) {
	path := "job/plan"

	respBody, err := c.makeRequest("POST", path, nil, job)
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
func (c *NomadClient) ForceNewPeriodicInstance(jobID string) error {
	path := fmt.Sprintf("job/%s/periodic/force", jobID)

	_, err := c.makeRequest("POST", path, nil, nil)
	return err
}

// GetJobScaleStatus retrieves the scale status of a job
func (c *NomadClient) GetJobScaleStatus(jobID, namespace string) (types.JobScaleStatus, error) {
	path := fmt.Sprintf("job/%s/scale", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
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
func (c *NomadClient) ScaleTaskGroup(jobID, group string, count int, namespace string) error {
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

	_, err := c.makeRequest("POST", path, queryParams, request)
	return err
}

// ListJobServices lists all services for a job
func (c *NomadClient) ListJobServices(jobID, namespace string) ([]types.Service, error) {
	path := fmt.Sprintf("job/%s/services", jobID)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var services []types.Service
	if err := json.Unmarshal(respBody, &services); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return services, nil
}

// ListVariables lists variables in the specified namespace
func (c *NomadClient) ListVariables(namespace, prefix string, nextToken string, perPage int, filter string) ([]types.Variable, error) {
	path := "vars"

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
	if prefix != "" {
		queryParams["prefix"] = prefix
	}
	if nextToken != "" {
		queryParams["next_token"] = nextToken
	}
	if perPage > 0 {
		queryParams["per_page"] = strconv.Itoa(perPage)
	}
	if filter != "" {
		queryParams["filter"] = filter
	}

	respBody, err := c.makeRequest("GET", path, queryParams, nil)
	if err != nil {
		return nil, err
	}

	var variables []types.Variable
	if err := json.Unmarshal(respBody, &variables); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return variables, nil
}

// GetVariable retrieves a specific variable by path
func (c *NomadClient) GetVariable(path, namespace string) (types.Variable, error) {
	apiPath := fmt.Sprintf("var/%s", path)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	respBody, err := c.makeRequest("GET", apiPath, queryParams, nil)
	if err != nil {
		return types.Variable{}, err
	}

	var variable types.Variable
	if err := json.Unmarshal(respBody, &variable); err != nil {
		return types.Variable{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return variable, nil
}

// CreateVariable creates a new variable
func (c *NomadClient) CreateVariable(variable types.Variable, namespace string, cas int, lockOperation string) error {
	apiPath := fmt.Sprintf("var/%s", variable.Path)

	// Parse the Value string into a map to use as request body
	var requestBody map[string]interface{}
	if err := json.Unmarshal([]byte(variable.Value), &requestBody); err != nil {
		return fmt.Errorf("failed to parse variable value: %v", err)
	}

	// Add CAS if provided
	if cas > 0 {
		requestBody["CAS"] = cas
	}

	// Add lock operation if provided
	if lockOperation != "" {
		requestBody["LockOperation"] = lockOperation
	}

	// Add namespace as query parameter if provided
	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}

	_, err := c.makeRequest("PUT", apiPath, queryParams, requestBody)
	return err
}

// DeleteVariable deletes a variable by path
func (c *NomadClient) DeleteVariable(path, namespace string, cas int) error {
	apiPath := fmt.Sprintf("var/%s", path)

	queryParams := make(map[string]string)
	if namespace != "" && namespace != "default" {
		queryParams["namespace"] = namespace
	}
	if cas > 0 {
		queryParams["cas"] = strconv.Itoa(cas)
	}

	_, err := c.makeRequest("DELETE", apiPath, queryParams, nil)
	return err
}

// ListAllocations lists all allocations in the cluster
func (c *NomadClient) ListAllocations() ([]types.Allocation, error) {
	respBody, err := c.makeRequest("GET", "allocations", nil, nil)
	if err != nil {
		return nil, err
	}

	var allocations []types.Allocation
	if err := json.Unmarshal(respBody, &allocations); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return allocations, nil
}
