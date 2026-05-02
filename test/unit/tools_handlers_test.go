package unit

import (
	"context"
	"io"
	"log"
	"testing"

	"github.com/kocierik/mcp-nomad/test/mocks"
	"github.com/kocierik/mcp-nomad/tools"
	"github.com/kocierik/mcp-nomad/types"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestListVariablesHandler_usesEffectiveToolNamespace(t *testing.T) {
	t.Run("NOMAD_NAMESPACE when tool omits namespace", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "from-env")

		var got string
		mock := &mocks.MockNomadClient{}
		mock.ListVariablesFunc = func(_ context.Context, namespace string, _ string, _ string, _ int, _ string) ([]types.Variable, error) {
			got = namespace
			return []types.Variable{{Path: "a", Value: `{}`}}, nil
		}

		h := tools.ListVariablesHandler(mock, testLogger())
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}

		res, err := h(context.Background(), req)
		require.NoError(t, err)
		require.False(t, res.IsError)
		require.NotEmpty(t, res.Content)
		assert.Equal(t, "from-env", got)
	})

	t.Run("explicit namespace overrides env", func(t *testing.T) {
		t.Setenv("NOMAD_NAMESPACE", "env-ns")

		var got string
		mock := &mocks.MockNomadClient{}
		mock.ListVariablesFunc = func(_ context.Context, namespace string, _ string, _ string, _ int, _ string) ([]types.Variable, error) {
			got = namespace
			return nil, nil
		}

		h := tools.ListVariablesHandler(mock, testLogger())
		req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
			"namespace": "explicit",
		}}}

		_, err := h(context.Background(), req)
		require.NoError(t, err)
		assert.Equal(t, "explicit", got)
	})
}

func TestListDeploymentsHandler_usesEffectiveToolNamespace(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "edge")

	var got string
	mock := &mocks.MockNomadClient{}
	mock.ListDeploymentsFunc = func(_ context.Context, namespace string) ([]types.DeploymentSummary, error) {
		got = namespace
		return nil, nil
	}

	h := tools.ListDeploymentsHandler(mock, testLogger())
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}}
	_, err := h(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "edge", got)
}

func TestScaleJobHandler_usesEffectiveToolNamespace(t *testing.T) {
	t.Setenv("NOMAD_NAMESPACE", "scale-ns")

	var got string
	mock := &mocks.MockNomadClient{}
	mock.ScaleTaskGroupFunc = func(_ context.Context, jobID string, group string, count int, namespace string) error {
		got = namespace
		return nil
	}

	h := tools.ScaleJobHandler(mock, testLogger())
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"job_id":    "job1",
		"group":     "web",
		"count":     float64(2),
		"namespace": nil,
	}}}

	res, err := h(context.Background(), req)
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, "scale-ns", got)
}

func TestGetJobHandler_InvalidArguments_IsErrorResult(t *testing.T) {
	t.Parallel()

	mock := &mocks.MockNomadClient{}
	h := tools.GetJobHandler(mock, testLogger())

	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: "not-a-map"}}

	res, err := h(context.Background(), req)
	require.NoError(t, err)
	require.True(t, res.IsError)
}

func TestGetAllocationHandler_returnsJSONFromClient(t *testing.T) {
	t.Parallel()

	want := types.Allocation{ID: "alloc-abc", JobID: "j1"}
	mock := &mocks.MockNomadClient{}
	mock.GetAllocationFunc = func(_ context.Context, allocID string) (types.Allocation, error) {
		assert.Equal(t, "alloc-abc", allocID)
		return want, nil
	}

	h := tools.GetAllocationHandler(mock, testLogger())
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"allocation_id": "alloc-abc",
	}}}

	res, err := h(context.Background(), req)
	require.NoError(t, err)
	require.False(t, res.IsError)

	require.Len(t, res.Content, 1)
	text, ok := res.Content[0].(mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, text.Text, `"ID": "alloc-abc"`)
}

func TestStopAllocationHandler_callsClientStop(t *testing.T) {
	t.Parallel()

	var got string
	mock := &mocks.MockNomadClient{}
	mock.StopAllocationFunc = func(_ context.Context, allocID string) error {
		got = allocID
		return nil
	}

	h := tools.StopAllocationHandler(mock, testLogger())
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"allocation_id": "a1",
	}}}

	res, err := h(context.Background(), req)
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, "a1", got)

	text, ok := res.Content[0].(mcp.TextContent)
	require.True(t, ok)
	assert.Contains(t, text.Text, "stopped successfully")
}

func TestListAllocationsHandler_passesNamespaceAndJob(t *testing.T) {
	t.Parallel()

	var gotNs, gotJob string
	mock := &mocks.MockNomadClient{}
	mock.ListAllocationsFunc = func(_ context.Context, namespace, jobID string) ([]types.Allocation, error) {
		gotNs, gotJob = namespace, jobID
		return []types.Allocation{}, nil
	}

	h := tools.ListAllocationsHandler(mock, testLogger())
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"namespace": "apps",
		"job_id":    "demo",
	}}}

	res, err := h(context.Background(), req)
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, "apps", gotNs)
	assert.Equal(t, "demo", gotJob)
}
