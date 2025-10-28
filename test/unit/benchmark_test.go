package unit

import (
	"testing"

	"github.com/kocierik/mcp-nomad/test/mocks"
	"github.com/kocierik/mcp-nomad/test/testdata"
	"github.com/kocierik/mcp-nomad/types"
)

// BenchmarkMockClientListJobs benchmarks the mock client ListJobs method directly
func BenchmarkMockClientListJobs(b *testing.B) {
	mockClient := &mocks.MockNomadClient{}
	mockClient.ListJobsFunc = func(namespace, status string) ([]types.JobSummary, error) {
		return testdata.SampleJobs, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockClient.ListJobs("default", "")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMockClientGetJob benchmarks the mock client GetJob method directly
func BenchmarkMockClientGetJob(b *testing.B) {
	mockClient := &mocks.MockNomadClient{}
	mockClient.GetJobFunc = func(jobID, namespace string) (types.Job, error) {
		return types.Job{
			ID:   jobID,
			Name: jobID,
			Type: "service",
		}, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockClient.GetJob("test-job-1", "default")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMockClientRunJob benchmarks the mock client RunJob method directly
func BenchmarkMockClientRunJob(b *testing.B) {
	mockClient := &mocks.MockNomadClient{}
	mockClient.RunJobFunc = func(jobSpec string, detach bool) (map[string]interface{}, error) {
		return map[string]interface{}{
			"EvalID":         "eval-123",
			"JobModifyIndex": 1,
		}, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockClient.RunJob(testdata.SampleJobSpecs["simple"], false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMockClientListNodes benchmarks the mock client ListNodes method directly
func BenchmarkMockClientListNodes(b *testing.B) {
	mockClient := &mocks.MockNomadClient{}
	mockClient.ListNodesFunc = func(status string) ([]types.NodeSummary, error) {
		return testdata.SampleNodes, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockClient.ListNodes("")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMockClientGetAllocationLogs benchmarks the mock client GetAllocationLogs method directly
func BenchmarkMockClientGetAllocationLogs(b *testing.B) {
	mockClient := &mocks.MockNomadClient{}
	mockClient.GetAllocationLogsFunc = func(allocID, task, logType string, follow bool, tail, offset int64) (string, error) {
		return testdata.SampleLogs["nginx_stdout"], nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mockClient.GetAllocationLogs("alloc-1", "nginx", "stdout", false, 0, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}
