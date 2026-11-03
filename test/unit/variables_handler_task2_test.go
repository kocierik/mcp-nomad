package unit

import (
	"context"
	"encoding/json"
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

type deleteVariableCASCapableClient struct {
	*mocks.MockNomadClient
	DeleteVariableCASFunc func(context.Context, string, string, *uint64) error
}

func (c *deleteVariableCASCapableClient) DeleteVariableCAS(ctx context.Context, path, namespace string, cas *uint64) error {
	if c.DeleteVariableCASFunc != nil {
		return c.DeleteVariableCASFunc(ctx, path, namespace, cas)
	}
	return nil
}

func textFromResult(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	require.NotNil(t, result)
	require.NotEmpty(t, result.Content)
	text, ok := result.Content[0].(mcp.TextContent)
	require.True(t, ok)
	return text.Text
}

func variableHandlerRequest(cas interface{}, lockOperation string) mcp.CallToolRequest {
	args := map[string]interface{}{
		"path":      "app/config",
		"key":       "k",
		"value":     "v",
		"namespace": "apps",
	}
	if cas != nil {
		args["cas"] = cas
	}
	if lockOperation != "" {
		args["lock_operation"] = lockOperation
	}
	return mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: args}}
}

func TestCreateVariableHandler_CASExplicitZero(t *testing.T) {
	t.Parallel()

	var gotCAS int
	var gotVar types.Variable
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
		gotCAS = cas
		gotVar = variable
		assert.Equal(t, "apps", namespace)
		assert.Equal(t, "", lockOperation)
		return nil
	}

	h := tools.CreateVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), variableHandlerRequest(json.Number("0"), ""))
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, 0, gotCAS)
	require.NotNil(t, gotVar.CAS)
	assert.Equal(t, uint64(0), *gotVar.CAS)
}

func TestCreateVariableHandler_CASLargeIntegerPreserved(t *testing.T) {
	t.Parallel()

	var gotCAS int
	var gotVar types.Variable
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
		gotCAS = cas
		gotVar = variable
		return nil
	}

	h := tools.CreateVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), variableHandlerRequest(json.Number("9007199254740993"), ""))
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, int64(9007199254740993), int64(gotCAS))
	require.NotNil(t, gotVar.CAS)
	assert.Equal(t, uint64(9007199254740993), *gotVar.CAS)
}

func TestCreateVariableHandler_InvalidCASRejected(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cas  interface{}
	}{
		{name: "fraction", cas: 1.5},
		{name: "negative", cas: -1.0},
		{name: "overflow", cas: json.Number("18446744073709551616")},
		{name: "malformed", cas: json.Number("not-a-number")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			called := false
			mock := &mocks.MockNomadClient{}
			mock.CreateVariableFunc = func(_ context.Context, _ types.Variable, _ string, _ int, _ string) error {
				called = true
				return nil
			}

			h := tools.CreateVariableHandler(mock, log.New(io.Discard, "", 0))
			res, err := h(context.Background(), variableHandlerRequest(tt.cas, ""))
			require.NoError(t, err)
			require.True(t, res.IsError)
			assert.Contains(t, textFromResult(t, res), "cas")
			assert.False(t, called)
		})
	}
}

func TestCreateVariableHandler_RejectsUnsupportedLockOperation(t *testing.T) {
	t.Parallel()

	called := false
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, _ types.Variable, _ string, _ int, _ string) error {
		called = true
		return nil
	}

	h := tools.CreateVariableHandler(mock, testLogger())
	res, err := h(context.Background(), variableHandlerRequest(nil, "acquire"))
	require.NoError(t, err)
	require.True(t, res.IsError)
	assert.Contains(t, textFromResult(t, res), "lock_operation")
	assert.False(t, called)
}

func TestCreateVariableHandler_CASOmitted(t *testing.T) {
	t.Parallel()

	var gotCAS int
	var gotVar types.Variable
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
		gotCAS = cas
		gotVar = variable
		assert.Equal(t, "apps", namespace)
		assert.Equal(t, "", lockOperation)
		return nil
	}

	h := tools.CreateVariableHandler(mock, testLogger())
	res, err := h(context.Background(), variableHandlerRequest(nil, ""))
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, 0, gotCAS)
	assert.Nil(t, gotVar.CAS)
}

func TestCreateVariableHandler_CASPositive(t *testing.T) {
	t.Parallel()

	var gotCAS int
	var gotVar types.Variable
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, variable types.Variable, namespace string, cas int, lockOperation string) error {
		gotCAS = cas
		gotVar = variable
		return nil
	}

	h := tools.CreateVariableHandler(mock, testLogger())
	res, err := h(context.Background(), variableHandlerRequest(json.Number("7"), ""))
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, 7, gotCAS)
	require.NotNil(t, gotVar.CAS)
	assert.Equal(t, uint64(7), *gotVar.CAS)
}

func TestCreateVariableHandler_CASLargeIntegerPreservedFromWireJSON(t *testing.T) {
	t.Parallel()

	var request mcp.CallToolRequest
	require.NoError(t, json.Unmarshal([]byte(`{
		"params": {
			"arguments": {
				"path": "app/config",
				"key": "k",
				"value": "v",
				"namespace": "apps",
				"cas": 9007199254740993
			}
		}
	}`), &request))

	var gotCAS int
	var gotVar types.Variable
	mock := &mocks.MockNomadClient{}
	mock.CreateVariableFunc = func(_ context.Context, variable types.Variable, _ string, cas int, _ string) error {
		gotCAS = cas
		gotVar = variable
		return nil
	}

	h := tools.CreateVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), request)
	require.NoError(t, err)
	require.False(t, res.IsError)
	require.NotNil(t, gotVar.CAS)
	assert.Equal(t, uint64(9007199254740993), *gotVar.CAS)
	assert.Equal(t, int64(9007199254740993), int64(gotCAS))
}

func TestDeleteVariableHandler_InvalidCASRejectedWithoutClientCall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cas  interface{}
	}{
		{name: "malformed string", cas: "abc"},
		{name: "fraction", cas: 1.5},
		{name: "negative", cas: -1.0},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			calledDelete := false
			calledDeleteCAS := false

			mock := &deleteVariableCASCapableClient{MockNomadClient: &mocks.MockNomadClient{}}
			mock.DeleteVariableFunc = func(_ context.Context, _ string, _ string, _ int) error {
				calledDelete = true
				return nil
			}
			mock.DeleteVariableCASFunc = func(_ context.Context, _ string, _ string, _ *uint64) error {
				calledDeleteCAS = true
				return nil
			}

			h := tools.DeleteVariableHandler(mock, log.New(io.Discard, "", 0))
			res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
				"path":      "app/config",
				"namespace": "apps",
				"cas":       tt.cas,
			}}})
			require.NoError(t, err)
			require.True(t, res.IsError)
			assert.Contains(t, textFromResult(t, res), "cas")
			assert.False(t, calledDelete)
			assert.False(t, calledDeleteCAS)
		})
	}
}

func TestDeleteVariableHandler_ExplicitCASZeroUsesCASAwareClient(t *testing.T) {
	t.Parallel()

	var gotCAS *uint64
	var gotPath string
	var gotNamespace string

	mock := &deleteVariableCASCapableClient{MockNomadClient: &mocks.MockNomadClient{}}
	mock.DeleteVariableCASFunc = func(_ context.Context, path, namespace string, cas *uint64) error {
		gotPath = path
		gotNamespace = namespace
		gotCAS = cas
		return nil
	}

	h := tools.DeleteVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"path":      "app/config",
		"namespace": "apps",
		"cas":       json.Number("0"),
	}}})
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, "app/config", gotPath)
	assert.Equal(t, "apps", gotNamespace)
	require.NotNil(t, gotCAS)
	assert.Equal(t, uint64(0), *gotCAS)
}

func TestDeleteVariableHandler_ExplicitCASPositiveUsesCASAwareClient(t *testing.T) {
	t.Parallel()

	var gotCAS *uint64
	mock := &deleteVariableCASCapableClient{MockNomadClient: &mocks.MockNomadClient{}}
	mock.DeleteVariableCASFunc = func(_ context.Context, _ string, _ string, cas *uint64) error {
		gotCAS = cas
		return nil
	}

	h := tools.DeleteVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"path":      "app/config",
		"namespace": "apps",
		"cas":       json.Number("7"),
	}}})
	require.NoError(t, err)
	require.False(t, res.IsError)
	require.NotNil(t, gotCAS)
	assert.Equal(t, uint64(7), *gotCAS)
}

func TestDeleteVariableHandler_ExplicitCASZeroRejectsLegacyClient(t *testing.T) {
	t.Parallel()

	called := false
	mock := &mocks.MockNomadClient{}
	mock.DeleteVariableFunc = func(_ context.Context, _ string, _ string, _ int) error {
		called = true
		return nil
	}

	handler := tools.DeleteVariableHandler(mock, log.New(io.Discard, "", 0))
	result, err := handler(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"path": "app/config",
		"cas":  json.Number("0"),
	}}})
	require.NoError(t, err)
	require.True(t, result.IsError)
	require.Contains(t, textFromResult(t, result), "cas")
	require.False(t, called)
}

func TestDeleteVariableHandler_WithoutCASAwareClientFallsBackToLegacyDelete(t *testing.T) {
	t.Parallel()

	var gotCAS int
	var gotPath string
	var gotNamespace string

	mock := &mocks.MockNomadClient{}
	mock.DeleteVariableFunc = func(_ context.Context, path, namespace string, cas int) error {
		gotPath = path
		gotNamespace = namespace
		gotCAS = cas
		return nil
	}

	h := tools.DeleteVariableHandler(mock, log.New(io.Discard, "", 0))
	res, err := h(context.Background(), mcp.CallToolRequest{Params: mcp.CallToolParams{Arguments: map[string]interface{}{
		"path":      "app/config",
		"namespace": "apps",
		"cas":       json.Number("7"),
	}}})
	require.NoError(t, err)
	require.False(t, res.IsError)
	assert.Equal(t, "app/config", gotPath)
	assert.Equal(t, "apps", gotNamespace)
	assert.Equal(t, 7, gotCAS)
}
