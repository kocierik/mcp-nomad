package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kocierik/mcp-nomad/test/testdata"
	"github.com/kocierik/mcp-nomad/types"
	"github.com/kocierik/mcp-nomad/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockNomadServer creates a mock Nomad API server for integration testing
type MockNomadServer struct {
	server *httptest.Server
}

func NewMockNomadServer() *MockNomadServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/job/", func(w http.ResponseWriter, r *http.Request) {
		jobID := r.URL.Path[len("/v1/job/"):]

		if r.Method == "GET" {
			job := types.Job{
				ID:   jobID,
				Name: jobID,
				Type: "service",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(job)
		} else if r.Method == "DELETE" {
			response := map[string]interface{}{
				"EvalID": "eval-456",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	})

	// Job submission endpoint
	mux.HandleFunc("/v1/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			jobs := testdata.SampleJobs
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jobs)
		} else if r.Method == "POST" {
			// Handle job submission (HCL parsing)
			// For testing, we'll just return a successful response
			response := map[string]interface{}{
				"EvalID":         "eval-123",
				"JobModifyIndex": 1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	})

	// Job parsing endpoint
	mux.HandleFunc("/v1/jobs/parse", func(w http.ResponseWriter, r *http.Request) {
		// Mock HCL parsing - return a simple job structure
		parsedJob := map[string]interface{}{
			"ID":          "test-job",
			"Name":        "test-job",
			"Type":        "service",
			"Datacenters": []string{"dc1"},
			"TaskGroups": []map[string]interface{}{
				{
					"Name":  "web",
					"Count": 2,
					"Tasks": []map[string]interface{}{
						{
							"Name":   "nginx",
							"Driver": "docker",
							"Config": map[string]interface{}{
								"image": "nginx:latest",
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(parsedJob)
	})

	// Nodes endpoints
	mux.HandleFunc("/v1/nodes", func(w http.ResponseWriter, r *http.Request) {
		nodes := testdata.SampleNodes
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodes)
	})

	mux.HandleFunc("/v1/node/", func(w http.ResponseWriter, r *http.Request) {
		nodeID := r.URL.Path[len("/v1/node/"):]

		if r.Method == "GET" {
			node := types.Node{
				ID:     nodeID,
				Name:   nodeID,
				Status: "ready",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(node)
		} else if r.Method == "POST" {
			// Handle drain/eligibility operations
			// Check if it's an eligibility request by path
			if strings.Contains(r.URL.Path, "/eligibility") {
				// Extract node ID from path (remove /eligibility suffix)
				actualNodeID := strings.TrimSuffix(nodeID, "/eligibility")
				// For eligibility, return a NodeSummary
				nodeSummary := types.NodeSummary{
					ID:     actualNodeID,
					Name:   actualNodeID,
					Status: "ready",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(nodeSummary)
			} else {
				// For drain operations
				response := map[string]interface{}{
					"EvalID": "eval-789",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	})

	// Namespaces endpoints
	mux.HandleFunc("/v1/namespaces", func(w http.ResponseWriter, r *http.Request) {
		namespaces := testdata.SampleNamespaces
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(namespaces)
	})

	mux.HandleFunc("/v1/namespace/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			response := map[string]interface{}{
				"CreateIndex": 1,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
		}
	})

	// Allocations endpoints
	mux.HandleFunc("/v1/allocations", func(w http.ResponseWriter, r *http.Request) {
		allocations := testdata.SampleAllocations
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allocations)
	})

	mux.HandleFunc("/v1/allocation/", func(w http.ResponseWriter, r *http.Request) {
		allocID := r.URL.Path[len("/v1/allocation/"):]
		allocation := types.Allocation{
			ID:   allocID,
			Name: "test-allocation",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allocation)
	})

	// Logs endpoint
	mux.HandleFunc("/v1/client/fs/logs/", func(w http.ResponseWriter, r *http.Request) {
		logs := testdata.SampleLogs["nginx_stdout"]
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(logs))
	})

	// Variables endpoints
	mux.HandleFunc("/v1/vars", func(w http.ResponseWriter, r *http.Request) {
		variables := testdata.SampleVariables
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(variables)
	})

	mux.HandleFunc("/v1/var/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			variable := testdata.SampleVariables[0]
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(variable)
		} else if r.Method == "PUT" {
			w.WriteHeader(http.StatusOK)
		} else if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
		}
	})

	// ACL endpoints
	mux.HandleFunc("/v1/acl/tokens", func(w http.ResponseWriter, r *http.Request) {
		tokens := testdata.SampleACLTokens
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	})

	mux.HandleFunc("/v1/acl/token/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			token := testdata.SampleACLTokens[0]
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(token)
		} else if r.Method == "POST" {
			token := testdata.SampleACLTokens[0]
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(token)
		} else if r.Method == "DELETE" {
			w.WriteHeader(http.StatusOK)
		}
	})

	// Cluster endpoints
	mux.HandleFunc("/v1/operator/raft/configuration", func(w http.ResponseWriter, r *http.Request) {
		data := testdata.SampleClusterData["leader"]
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	mux.HandleFunc("/v1/regions", func(w http.ResponseWriter, r *http.Request) {
		data := testdata.SampleClusterData["regions"]
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	// Status endpoint for connection testing
	mux.HandleFunc("/v1/status/leader", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"Leader": "127.0.0.1:4647",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	server := httptest.NewServer(mux)

	return &MockNomadServer{
		server: server,
	}
}

func (m *MockNomadServer) URL() string {
	return m.server.URL
}

func (m *MockNomadServer) Close() {
	m.server.Close()
}

func TestNomadClientIntegration(t *testing.T) {
	// Start mock server
	mockServer := NewMockNomadServer()
	defer mockServer.Close()

	// Create client
	client, err := utils.NewNomadClient(mockServer.URL(), "")
	require.NoError(t, err)

	t.Run("ListJobs", func(t *testing.T) {
		jobs, err := client.ListJobs("default", "")
		require.NoError(t, err)
		assert.Len(t, jobs, 2)
		assert.Equal(t, "test-job-1", jobs[0].ID)
		assert.Equal(t, "test-job-2", jobs[1].ID)
	})

	t.Run("GetJob", func(t *testing.T) {
		job, err := client.GetJob("test-job-1", "default")
		require.NoError(t, err)
		assert.Equal(t, "test-job-1", job.ID)
		assert.Equal(t, "service", job.Type)
	})

	t.Run("RunJob", func(t *testing.T) {
		result, err := client.RunJob(testdata.SampleJobSpecs["simple"], false)
		require.NoError(t, err)
		assert.Contains(t, result, "EvalID")
		assert.Equal(t, "eval-123", result["EvalID"])
	})

	t.Run("StopJob", func(t *testing.T) {
		result, err := client.StopJob("test-job-1", "default", false)
		require.NoError(t, err)
		assert.Contains(t, result, "EvalID")
		assert.Equal(t, "eval-456", result["EvalID"])
	})

	t.Run("ListNodes", func(t *testing.T) {
		nodes, err := client.ListNodes("")
		require.NoError(t, err)
		assert.Len(t, nodes, 2)
		assert.Equal(t, "node-1", nodes[0].ID)
		assert.Equal(t, "node-2", nodes[1].ID)
	})

	t.Run("GetNode", func(t *testing.T) {
		node, err := client.GetNode("node-1")
		require.NoError(t, err)
		assert.Equal(t, "node-1", node.ID)
		assert.Equal(t, "ready", node.Status)
	})

	t.Run("DrainNode", func(t *testing.T) {
		result, err := client.DrainNode("node-1", true, 300)
		require.NoError(t, err)
		assert.Contains(t, result, "drain enabled")
	})

	t.Run("EligibilityNode", func(t *testing.T) {
		node, err := client.EligibilityNode("node-1", "eligible")
		require.NoError(t, err)
		assert.Equal(t, "node-1", node.ID)
	})

	t.Run("ListNamespaces", func(t *testing.T) {
		namespaces, err := client.ListNamespaces()
		require.NoError(t, err)
		assert.Len(t, namespaces, 2)
		assert.Equal(t, "default", namespaces[0].Name)
		assert.Equal(t, "production", namespaces[1].Name)
	})

	t.Run("CreateNamespace", func(t *testing.T) {
		namespace := types.Namespace{
			Name:        "test-namespace",
			Description: "Test namespace",
		}
		err := client.CreateNamespace(namespace)
		require.NoError(t, err)
	})

	t.Run("DeleteNamespace", func(t *testing.T) {
		err := client.DeleteNamespace("test-namespace")
		require.NoError(t, err)
	})

	t.Run("ListAllocations", func(t *testing.T) {
		allocations, err := client.ListAllocations()
		require.NoError(t, err)
		assert.Len(t, allocations, 1)
		assert.Equal(t, "alloc-1", allocations[0].ID)
	})

	t.Run("GetAllocation", func(t *testing.T) {
		allocation, err := client.GetAllocation("alloc-1")
		require.NoError(t, err)
		assert.Equal(t, "alloc-1", allocation.ID)
		assert.Equal(t, "test-allocation", allocation.Name)
	})

	t.Run("GetAllocationLogs", func(t *testing.T) {
		logs, err := client.GetAllocationLogs("alloc-1", "nginx", "stdout", false, 0, 0)
		require.NoError(t, err)
		assert.Contains(t, logs, "Starting nginx")
		assert.Contains(t, logs, "Server started on port 80")
	})

	t.Run("ListVariables", func(t *testing.T) {
		variables, err := client.ListVariables("default", "", "", 0, "")
		require.NoError(t, err)
		assert.Len(t, variables, 1)
		assert.Equal(t, "app/config", variables[0].Path)
	})

	t.Run("GetVariable", func(t *testing.T) {
		variable, err := client.GetVariable("app/config", "default")
		require.NoError(t, err)
		assert.Equal(t, "app/config", variable.Path)
		assert.Contains(t, variable.Value, "database_url")
	})

	t.Run("CreateVariable", func(t *testing.T) {
		variable := types.Variable{
			Path:  "test/config",
			Value: `{"Items":{"key":"value"}}`,
		}
		err := client.CreateVariable(variable, "default", 0, "")
		require.NoError(t, err)
	})

	t.Run("DeleteVariable", func(t *testing.T) {
		err := client.DeleteVariable("test/config", "default", 0)
		require.NoError(t, err)
	})

	t.Run("ListACLTokens", func(t *testing.T) {
		tokens, err := client.ListACLTokens()
		require.NoError(t, err)
		assert.Len(t, tokens, 1)
		assert.Equal(t, "token-1", tokens[0].AccessorID)
	})

	t.Run("GetACLToken", func(t *testing.T) {
		token, err := client.GetACLToken("token-1")
		require.NoError(t, err)
		assert.Equal(t, "token-1", token.AccessorID)
		assert.Equal(t, "test-token", token.Name)
	})

	t.Run("CreateACLToken", func(t *testing.T) {
		token := types.ACLToken{
			Name:     "test-token",
			Type:     "client",
			Policies: []string{"read-only"},
		}
		newToken, err := client.CreateACLToken(token)
		require.NoError(t, err)
		assert.Equal(t, "test-token", newToken.Name)
	})

	t.Run("DeleteACLToken", func(t *testing.T) {
		err := client.DeleteACLToken("token-1")
		require.NoError(t, err)
	})

	t.Run("GetClusterLeader", func(t *testing.T) {
		data, err := client.GetClusterLeader()
		require.NoError(t, err)
		assert.Contains(t, string(data), "server-1")
	})

	t.Run("ListClusterPeers", func(t *testing.T) {
		data, err := client.ListClusterPeers()
		require.NoError(t, err)
		assert.Contains(t, string(data), "server-1")
	})

	t.Run("ListRegions", func(t *testing.T) {
		data, err := client.ListRegions()
		require.NoError(t, err)
		assert.Contains(t, string(data), "global")
	})
}

func TestNomadClientErrorHandling(t *testing.T) {
	// Test with invalid server URL
	_, err := utils.NewNomadClient("http://invalid-url:9999", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to Nomad server")

	// Test with empty address
	_, err = utils.NewNomadClient("", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nomad address is required")
}
