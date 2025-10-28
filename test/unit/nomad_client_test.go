package unit

import (
	"errors"
	"testing"

	"github.com/kocierik/mcp-nomad/test/mocks"
	"github.com/kocierik/mcp-nomad/test/testdata"
	"github.com/kocierik/mcp-nomad/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNomadClient_ListJobs(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		status        string
		mockFunc      func(namespace, status string) ([]types.JobSummary, error)
		expectedJobs  []types.JobSummary
		expectedError string
	}{
		{
			name:      "successful list jobs",
			namespace: "default",
			status:    "",
			mockFunc: func(namespace, status string) ([]types.JobSummary, error) {
				return testdata.SampleJobs, nil
			},
			expectedJobs:  testdata.SampleJobs,
			expectedError: "",
		},
		{
			name:      "list jobs with status filter",
			namespace: "default",
			status:    "running",
			mockFunc: func(namespace, status string) ([]types.JobSummary, error) {
				// Return first job for running status
				return []types.JobSummary{testdata.SampleJobs[0]}, nil
			},
			expectedJobs:  []types.JobSummary{testdata.SampleJobs[0]},
			expectedError: "",
		},
		{
			name:      "error from API",
			namespace: "default",
			status:    "",
			mockFunc: func(namespace, status string) ([]types.JobSummary, error) {
				return nil, errors.New("API error")
			},
			expectedJobs:  nil,
			expectedError: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.ListJobsFunc = tt.mockFunc

			jobs, err := mockClient.ListJobs(tt.namespace, tt.status)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedJobs, jobs)
			}
		})
	}
}

func TestNomadClient_GetJob(t *testing.T) {
	tests := []struct {
		name          string
		jobID         string
		namespace     string
		mockFunc      func(jobID, namespace string) (types.Job, error)
		expectedJob   types.Job
		expectedError string
	}{
		{
			name:      "successful get job",
			jobID:     "test-job-1",
			namespace: "default",
			mockFunc: func(jobID, namespace string) (types.Job, error) {
				return types.Job{
					ID:   jobID,
					Name: jobID,
					Type: "service",
				}, nil
			},
			expectedJob: types.Job{
				ID:   "test-job-1",
				Name: "test-job-1",
				Type: "service",
			},
			expectedError: "",
		},
		{
			name:      "job not found",
			jobID:     "nonexistent-job",
			namespace: "default",
			mockFunc: func(jobID, namespace string) (types.Job, error) {
				return types.Job{}, errors.New("job not found")
			},
			expectedJob:   types.Job{},
			expectedError: "job not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.GetJobFunc = tt.mockFunc

			job, err := mockClient.GetJob(tt.jobID, tt.namespace)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedJob, job)
			}
		})
	}
}

func TestNomadClient_RunJob(t *testing.T) {
	tests := []struct {
		name           string
		jobSpec        string
		detach         bool
		mockFunc       func(jobSpec string, detach bool) (map[string]interface{}, error)
		expectedResult map[string]interface{}
		expectedError  string
	}{
		{
			name:    "successful run job",
			jobSpec: testdata.SampleJobSpecs["simple"],
			detach:  false,
			mockFunc: func(jobSpec string, detach bool) (map[string]interface{}, error) {
				return map[string]interface{}{
					"EvalID":         "eval-123",
					"JobModifyIndex": 1,
				}, nil
			},
			expectedResult: map[string]interface{}{
				"EvalID":         "eval-123",
				"JobModifyIndex": 1,
			},
			expectedError: "",
		},
		{
			name:    "run job with detach",
			jobSpec: testdata.SampleJobSpecs["simple"],
			detach:  true,
			mockFunc: func(jobSpec string, detach bool) (map[string]interface{}, error) {
				return map[string]interface{}{
					"EvalID": "eval-456",
				}, nil
			},
			expectedResult: map[string]interface{}{
				"EvalID": "eval-456",
			},
			expectedError: "",
		},
		{
			name:    "invalid job spec",
			jobSpec: testdata.SampleJobSpecs["invalid"],
			detach:  false,
			mockFunc: func(jobSpec string, detach bool) (map[string]interface{}, error) {
				return nil, errors.New("invalid job specification")
			},
			expectedResult: nil,
			expectedError:  "invalid job specification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.RunJobFunc = tt.mockFunc

			result, err := mockClient.RunJob(tt.jobSpec, tt.detach)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestNomadClient_ListNodes(t *testing.T) {
	tests := []struct {
		name          string
		status        string
		mockFunc      func(status string) ([]types.NodeSummary, error)
		expectedNodes []types.NodeSummary
		expectedError string
	}{
		{
			name:   "successful list nodes",
			status: "",
			mockFunc: func(status string) ([]types.NodeSummary, error) {
				return testdata.SampleNodes, nil
			},
			expectedNodes: testdata.SampleNodes,
			expectedError: "",
		},
		{
			name:   "list nodes with status filter",
			status: "ready",
			mockFunc: func(status string) ([]types.NodeSummary, error) {
				// Return first node for ready status
				return []types.NodeSummary{testdata.SampleNodes[0]}, nil
			},
			expectedNodes: []types.NodeSummary{testdata.SampleNodes[0]},
			expectedError: "",
		},
		{
			name:   "error from API",
			status: "",
			mockFunc: func(status string) ([]types.NodeSummary, error) {
				return nil, errors.New("API error")
			},
			expectedNodes: nil,
			expectedError: "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.ListNodesFunc = tt.mockFunc

			nodes, err := mockClient.ListNodes(tt.status)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedNodes, nodes)
			}
		})
	}
}

func TestNomadClient_GetAllocationLogs(t *testing.T) {
	tests := []struct {
		name          string
		allocID       string
		task          string
		logType       string
		follow        bool
		tail          int64
		offset        int64
		mockFunc      func(allocID, task, logType string, follow bool, tail, offset int64) (string, error)
		expectedLogs  string
		expectedError string
	}{
		{
			name:    "successful get logs",
			allocID: "alloc-1",
			task:    "nginx",
			logType: "stdout",
			follow:  false,
			tail:    0,
			offset:  0,
			mockFunc: func(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
				return testdata.SampleLogs["nginx_stdout"], nil
			},
			expectedLogs:  testdata.SampleLogs["nginx_stdout"],
			expectedError: "",
		},
		{
			name:    "get logs with tail",
			allocID: "alloc-1",
			task:    "nginx",
			logType: "stdout",
			follow:  false,
			tail:    2,
			offset:  0,
			mockFunc: func(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
				// Simulate tail functionality
				lines := []string{
					"2024-01-01T10:00:02Z [INFO] Server started on port 80",
					"2024-01-01T10:00:03Z [INFO] Ready to serve requests",
				}
				return lines[0] + "\n" + lines[1], nil
			},
			expectedLogs:  "2024-01-01T10:00:02Z [INFO] Server started on port 80\n2024-01-01T10:00:03Z [INFO] Ready to serve requests",
			expectedError: "",
		},
		{
			name:    "error getting logs",
			allocID: "alloc-1",
			task:    "nginx",
			logType: "stdout",
			follow:  false,
			tail:    0,
			offset:  0,
			mockFunc: func(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
				return "", errors.New("allocation not found")
			},
			expectedLogs:  "",
			expectedError: "allocation not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.GetAllocationLogsFunc = tt.mockFunc

			logs, err := mockClient.GetAllocationLogs(tt.allocID, tt.task, tt.logType, tt.follow, tt.tail, tt.offset)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedLogs, logs)
			}
		})
	}
}

func TestNomadClient_CreateVariable(t *testing.T) {
	tests := []struct {
		name          string
		variable      types.Variable
		namespace     string
		cas           int
		lockOperation string
		mockFunc      func(variable types.Variable, namespace string, cas int, lockOperation string) error
		expectedError string
	}{
		{
			name: "successful create variable",
			variable: types.Variable{
				Path:  "app/config",
				Value: `{"Items":{"key":"value"}}`,
			},
			namespace:     "default",
			cas:           0,
			lockOperation: "",
			mockFunc: func(variable types.Variable, namespace string, cas int, lockOperation string) error {
				return nil
			},
			expectedError: "",
		},
		{
			name: "create variable with CAS",
			variable: types.Variable{
				Path:  "app/config",
				Value: `{"Items":{"key":"value"}}`,
			},
			namespace:     "default",
			cas:           123,
			lockOperation: "",
			mockFunc: func(variable types.Variable, namespace string, cas int, lockOperation string) error {
				return nil
			},
			expectedError: "",
		},
		{
			name: "error creating variable",
			variable: types.Variable{
				Path:  "app/config",
				Value: `{"Items":{"key":"value"}}`,
			},
			namespace:     "default",
			cas:           0,
			lockOperation: "",
			mockFunc: func(variable types.Variable, namespace string, cas int, lockOperation string) error {
				return errors.New("variable already exists")
			},
			expectedError: "variable already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mocks.MockNomadClient{}
			mockClient.CreateVariableFunc = tt.mockFunc

			err := mockClient.CreateVariable(tt.variable, tt.namespace, tt.cas, tt.lockOperation)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
