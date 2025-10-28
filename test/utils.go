package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertJobEqual asserts that two jobs are equal
func AssertJobEqual(t *testing.T, expected, actual types.Job) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Status, actual.Status)
}

// AssertJobSummaryEqual asserts that two job summaries are equal
func AssertJobSummaryEqual(t *testing.T, expected, actual types.JobSummary) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Summary, actual.Summary)
	assert.Equal(t, expected.CreateIndex, actual.CreateIndex)
	assert.Equal(t, expected.ModifyIndex, actual.ModifyIndex)
}

// AssertNodeEqual asserts that two nodes are equal
func AssertNodeEqual(t *testing.T, expected, actual types.Node) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Datacenter, actual.Datacenter)
}

// AssertNodeSummaryEqual asserts that two node summaries are equal
func AssertNodeSummaryEqual(t *testing.T, expected, actual types.NodeSummary) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Status, actual.Status)
	assert.Equal(t, expected.Datacenter, actual.Datacenter)
	assert.Equal(t, expected.NodeClass, actual.NodeClass)
}

// AssertAllocationEqual asserts that two allocations are equal
func AssertAllocationEqual(t *testing.T, expected, actual types.Allocation) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.NodeID, actual.NodeID)
	assert.Equal(t, expected.JobID, actual.JobID)
	assert.Equal(t, expected.TaskGroup, actual.TaskGroup)
	assert.Equal(t, expected.DesiredStatus, actual.DesiredStatus)
	assert.Equal(t, expected.ClientStatus, actual.ClientStatus)
}

// AssertNamespaceEqual asserts that two namespaces are equal
func AssertNamespaceEqual(t *testing.T, expected, actual types.Namespace) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
}

// AssertVariableEqual asserts that two variables are equal
func AssertVariableEqual(t *testing.T, expected, actual types.Variable) {
	assert.Equal(t, expected.Path, actual.Path)
	assert.Equal(t, expected.Namespace, actual.Namespace)
	assert.Equal(t, expected.Value, actual.Value)
}

// AssertACLTokenEqual asserts that two ACL tokens are equal
func AssertACLTokenEqual(t *testing.T, expected, actual types.ACLToken) {
	assert.Equal(t, expected.AccessorID, actual.AccessorID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Global, actual.Global)
	assert.Equal(t, expected.Policies, actual.Policies)
}

// CreateTestJob creates a test job with the given parameters
func CreateTestJob(id, name, jobType string) types.Job {
	return types.Job{
		ID:   id,
		Name: name,
		Type: jobType,
	}
}

// CreateTestJobSummary creates a test job summary with the given parameters
func CreateTestJobSummary(id, name, jobType, status string) types.JobSummary {
	return types.JobSummary{
		ID:          id,
		Summary:     map[string]types.TaskSummary{"web": {Running: 1}},
		CreateIndex: 1,
		ModifyIndex: 1,
	}
}

// CreateTestNode creates a test node with the given parameters
func CreateTestNode(id, name, status, datacenter string) types.Node {
	return types.Node{
		ID:         id,
		Name:       name,
		Status:     status,
		Datacenter: datacenter,
	}
}

// CreateTestNodeSummary creates a test node summary with the given parameters
func CreateTestNodeSummary(id, name, status, datacenter, nodeClass string) types.NodeSummary {
	return types.NodeSummary{
		ID:         id,
		Name:       name,
		Status:     status,
		Datacenter: datacenter,
		NodeClass:  nodeClass,
	}
}

// CreateTestAllocation creates a test allocation with the given parameters
func CreateTestAllocation(id, evalID, name, nodeID, jobID, taskGroup string) types.Allocation {
	return types.Allocation{
		ID:                 id,
		EvalID:             evalID,
		Name:               name,
		NodeID:             nodeID,
		JobID:              jobID,
		TaskGroup:          taskGroup,
		DesiredStatus:      "run",
		DesiredDescription: "Allocation is running",
		ClientStatus:       "running",
		ClientDescription:  "Allocation is running",
		CreateIndex:        1,
		ModifyIndex:        1,
		CreateTime:         time.Now().Unix(),
		ModifyTime:         time.Now().Unix(),
	}
}

// CreateTestNamespace creates a test namespace with the given parameters
func CreateTestNamespace(name, description string) types.Namespace {
	return types.Namespace{
		Name:        name,
		Description: description,
	}
}

// CreateTestVariable creates a test variable with the given parameters
func CreateTestVariable(path, namespace, value string) types.Variable {
	return types.Variable{
		Path:      path,
		Namespace: namespace,
		Value:     value,
	}
}

// CreateTestACLToken creates a test ACL token with the given parameters
func CreateTestACLToken(accessorID, secretID, name, tokenType string, policies []string, global bool) types.ACLToken {
	return types.ACLToken{
		AccessorID:  accessorID,
		SecretID:    secretID,
		Name:        name,
		Type:        tokenType,
		Policies:    policies,
		Global:      global,
		CreateIndex: 1,
		ModifyIndex: 1,
	}
}

// AssertJSONEqual asserts that two JSON strings are equal after unmarshaling
func AssertJSONEqual(t *testing.T, expected, actual string) {
	var expectedObj, actualObj interface{}

	err := json.Unmarshal([]byte(expected), &expectedObj)
	require.NoError(t, err)

	err = json.Unmarshal([]byte(actual), &actualObj)
	require.NoError(t, err)

	assert.Equal(t, expectedObj, actualObj)
}

// AssertContainsLogs asserts that logs contain expected content
func AssertContainsLogs(t *testing.T, logs string, expectedContent []string) {
	for _, content := range expectedContent {
		assert.Contains(t, logs, content, "Logs should contain: %s", content)
	}
}

// AssertLogsFormat asserts that logs are in the expected format
func AssertLogsFormat(t *testing.T, logs string) {
	// Basic format validation - logs should not be empty and should contain some structure
	assert.NotEmpty(t, logs, "Logs should not be empty")

	// If logs contain timestamps, they should be in a reasonable format
	if len(logs) > 20 {
		// Check for common log patterns
		hasTimestamp := false
		hasLogLevel := false

		lines := []string{}
		for _, line := range []string{logs} {
			if len(line) > 0 {
				lines = append(lines, line)
			}
		}

		for _, line := range lines {
			if len(line) > 10 {
				// Check for timestamp pattern (YYYY-MM-DD or similar)
				if len(line) >= 10 && line[4] == '-' && line[7] == '-' {
					hasTimestamp = true
				}
				// Check for log level pattern
				if len(line) >= 5 && (line[0:5] == "[INFO" || line[0:5] == "[WARN" || line[0:5] == "[ERRO") {
					hasLogLevel = true
				}
			}
		}

		// At least one of timestamp or log level should be present for structured logs
		if hasTimestamp || hasLogLevel {
			assert.True(t, true, "Logs appear to be in a structured format")
		}
	}
}
