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

func TestListAllocationsHandler_filtersAndCompact(t *testing.T) {
	mock := &mocks.MockNomadClient{}
	mock.ListAllocationsFunc = func(context.Context, string, string) ([]types.Allocation, error) {
		return []types.Allocation{
			{ID: "a1", NodeID: "node-aaaa", ClientStatus: "pending", JobID: "j1", CreateTime: 1700000000000000000,
				TaskStates: map[string]types.TaskState{"t": {State: "pending", Events: []types.TaskEvent{
					{Type: "Template", DisplayMessage: "Missing: vault.read(secret/x)"}}}}},
			{ID: "a2", NodeID: "node-aaaa", ClientStatus: "running", JobID: "j2"},
			{ID: "a3", NodeID: "node-bbbb", ClientStatus: "pending", JobID: "j3"},
		}, nil
	}
	h := tools.ListAllocationsHandler(mock, testLogger())
	call := func(args map[string]interface{}) string {
		res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}})
		require.NoError(t, err)
		require.False(t, res.IsError)
		return res.Content[0].(mcp.TextContent).Text
	}

	var full []types.Allocation
	require.NoError(t, json.Unmarshal([]byte(call(map[string]interface{}{})), &full))
	assert.Len(t, full, 3)

	require.NoError(t, json.Unmarshal([]byte(call(map[string]interface{}{"node_id": "node-aaaa", "client_status": "pending"})), &full))
	require.Len(t, full, 1)
	assert.Equal(t, "a1", full[0].ID)

	// compact accepts a real bool and the string "true" (clients with a stale schema send strings)
	for _, v := range []interface{}{true, "true"} {
		var compact []map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(call(map[string]interface{}{"client_status": "pending", "compact": v})), &compact))
		require.Len(t, compact, 2)
		assert.Equal(t, "a1", compact[0]["id"])
		tasks := compact[0]["tasks"].(map[string]interface{})
		assert.Equal(t, "Template: Missing: vault.read(secret/x)", tasks["t"].(map[string]interface{})["last_event"])
	}
}

func TestGetEvaluationHandler(t *testing.T) {
	mock := &mocks.MockNomadClient{}
	mock.GetEvaluationFunc = func(_ context.Context, id, ns string) (types.Evaluation, error) {
		assert.Equal(t, "ev1", id)
		assert.Equal(t, "ns1", ns)
		return types.Evaluation{ID: id, FailedTGAllocs: map[string]interface{}{"grp": map[string]interface{}{"NodesEvaluated": 3}}}, nil
	}
	h := tools.GetEvaluationHandler(mock, testLogger())
	res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{"evaluation_id": "ev1", "namespace": "ns1"}}})
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Contains(t, res.Content[0].(mcp.TextContent).Text, "NodesEvaluated")

	res, err = h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{}}})
	require.NoError(t, err)
	assert.True(t, res.IsError)
}
