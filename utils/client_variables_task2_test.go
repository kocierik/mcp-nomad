package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kocierik/mcp-nomad/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVariable_DecodesOfficialVariablePayload(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/var/app/config", r.URL.Path)
		require.Equal(t, "apps", r.URL.Query().Get("namespace"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Path":"app/config","Namespace":"apps","ModifyIndex":42,"Items":{"key":"value"}}`))
	}))
	defer server.Close()

	client := &NomadClient{address: server.URL, httpClient: server.Client()}

	v, err := client.GetVariable(context.Background(), "app/config", "apps")
	require.NoError(t, err)
	assert.Equal(t, "app/config", v.Path)
	assert.Equal(t, "apps", v.Namespace)
	assert.Equal(t, 42, v.ModifyIndex)
	assert.Equal(t, map[string]string{"key": "value"}, v.Items)
}

func TestCreateVariable_CASQueryBehaviorAndBody(t *testing.T) {
	t.Parallel()

	makeClient := func(t *testing.T, assertRequest func(*http.Request, map[string]interface{})) *NomadClient {
		t.Helper()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assertRequest(r, body)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)
		return &NomadClient{address: server.URL, httpClient: server.Client()}
	}

	t.Run("omitted cas when legacy caller passes zero", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request, body map[string]interface{}) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/v1/var/app/config", r.URL.Path)
			require.Equal(t, "apps", r.URL.Query().Get("namespace"))
			_, hasCASQuery := r.URL.Query()["cas"]
			assert.False(t, hasCASQuery)
			_, hasCASField := body["CAS"]
			assert.False(t, hasCASField)
			assert.Equal(t, map[string]interface{}{"key": "value"}, body["Items"])
		})

		err := client.CreateVariable(context.Background(), types.Variable{
			Path:  "app/config",
			Value: `{"Items":{"key":"value"}}`,
		}, "apps", 0, "")
		require.NoError(t, err)
	})

	t.Run("explicit cas zero can be requested", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request, body map[string]interface{}) {
			require.Equal(t, "0", r.URL.Query().Get("cas"))
			_, hasCASField := body["CAS"]
			assert.False(t, hasCASField)
			assert.Equal(t, map[string]interface{}{"k": "v"}, body["Items"])
		})

		cas := uint64(0)
		err := client.CreateVariable(context.Background(), types.Variable{
			Path:  "app/config",
			Items: map[string]string{"k": "v"},
			CAS:   &cas,
		}, "apps", 0, "")
		require.NoError(t, err)
	})

	t.Run("positive cas goes to query", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request, body map[string]interface{}) {
			require.Equal(t, "7", r.URL.Query().Get("cas"))
			_, hasCASField := body["CAS"]
			assert.False(t, hasCASField)
			assert.Equal(t, map[string]interface{}{"a": "b"}, body["Items"])
		})

		err := client.CreateVariable(context.Background(), types.Variable{
			Path:  "app/config",
			Items: map[string]string{"a": "b"},
		}, "apps", 7, "")
		require.NoError(t, err)
	})

	t.Run("legacy value with unsupported keys only forwards official payload fields", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request, body map[string]interface{}) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "5", r.URL.Query().Get("cas"))
			assert.Equal(t, map[string]interface{}{"k": "v"}, body["Items"])
			assert.Equal(t, "app/config", body["Path"])
			assert.Equal(t, "apps", body["Namespace"])
			_, hasCAS := body["CAS"]
			_, hasLock := body["LockOperation"]
			assert.False(t, hasCAS)
			assert.False(t, hasLock)
		})

		err := client.CreateVariable(context.Background(), types.Variable{
			Path:  "app/config",
			Value: `{"Items":{"k":"v"},"CAS":11,"LockOperation":"acquire","Path":"override","Namespace":"override"}`,
		}, "apps", 5, "")
		require.NoError(t, err)
	})
}

func TestCreateVariable_RejectsLockOperationBeforeHTTP(t *testing.T) {
	t.Parallel()

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &NomadClient{address: server.URL, httpClient: server.Client()}
	err := client.CreateVariable(context.Background(), types.Variable{Path: "app/config", Items: map[string]string{"key": "value"}}, "apps", 0, "acquire")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported lock operation")
	assert.False(t, called)
}

func TestDeleteVariableCAS_QueryBehavior(t *testing.T) {
	t.Parallel()

	makeClient := func(t *testing.T, assertRequest func(*http.Request)) *NomadClient {
		t.Helper()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertRequest(r)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)
		return &NomadClient{address: server.URL, httpClient: server.Client()}
	}

	t.Run("explicit zero is preserved in query", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/var/app/config", r.URL.Path)
			require.Equal(t, "apps", r.URL.Query().Get("namespace"))
			require.Equal(t, "0", r.URL.Query().Get("cas"))
		})

		cas := uint64(0)
		err := client.DeleteVariableCAS(context.Background(), "app/config", "apps", &cas)
		require.NoError(t, err)
	})

	t.Run("positive cas is preserved in query", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request) {
			require.Equal(t, "12", r.URL.Query().Get("cas"))
		})

		cas := uint64(12)
		err := client.DeleteVariableCAS(context.Background(), "app/config", "apps", &cas)
		require.NoError(t, err)
	})

	t.Run("omitted cas does not set query", func(t *testing.T) {
		client := makeClient(t, func(r *http.Request) {
			_, hasCAS := r.URL.Query()["cas"]
			assert.False(t, hasCAS)
		})

		err := client.DeleteVariableCAS(context.Background(), "app/config", "apps", nil)
		require.NoError(t, err)
	})
}
