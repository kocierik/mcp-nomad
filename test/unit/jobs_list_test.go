package unit

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kocierik/mcp-nomad/test/mocks"
	"github.com/kocierik/mcp-nomad/tools"
	"github.com/kocierik/mcp-nomad/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// list_jobs must build its listing from the single /v1/jobs call: no per-job
// GetJob/GetJobSummary round-trips (an N+1 that timed out on large namespaces).
func TestListJobsHandler_singleCallAndFilters(t *testing.T) {
	mock := &mocks.MockNomadClient{}
	mock.ListJobsFunc = func(_ context.Context, _ string, _ string) ([]types.JobSummary, error) {
		return []types.JobSummary{
			{ID: "web", Name: "web", Type: "service", Status: "running"},
			{ID: "web-old", Name: "web-old", Type: "service", Status: "dead"},
			{ID: "batch-1", Name: "batch-1", Type: "batch", Status: "running"},
		}, nil
	}
	mock.GetJobFunc = func(context.Context, string, string) (types.Job, error) {
		t.Fatal("list_jobs must not call GetJob per job")
		return types.Job{}, nil
	}
	mock.GetJobSummaryFunc = func(context.Context, string, string) (types.JobSummary, error) {
		t.Fatal("list_jobs must not call GetJobSummary per job")
		return types.JobSummary{}, nil
	}

	call := func(args map[string]interface{}) []map[string]interface{} {
		h := tools.ListJobsHandler(mock, testLogger())
		res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
		require.NoError(t, err)
		require.False(t, res.IsError)
		var out []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &out))
		return out
	}

	assert.Len(t, call(map[string]interface{}{}), 3)
	running := call(map[string]interface{}{"status": "running"})
	require.Len(t, running, 2)
	assert.Equal(t, "web", running[0]["Name"])
	assert.Len(t, call(map[string]interface{}{"prefix": "web"}), 2)
}
